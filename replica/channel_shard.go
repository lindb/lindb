// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package replica

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./channel_shard.go -destination=./channel_shard_mock.go -package=replica

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
	// Write writes the data into the channel, ErrCanceled is returned when the channel is canceled before
	// data is written successfully.
	// Concurrent safe.
	Write(metric *protoMetricsV1.Metric) error

	SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)

	Stop()
}

// channel implements Channel.
type channel struct {
	// context to close channel
	ctx context.Context

	database string
	shardID  models.ShardID

	newWriteStreamFn func(
		ctx context.Context,
		target models.Node,
		database string, shardState *models.ShardState,
		fct rpc.ClientStreamFactory,
	) (rpc.WriteStream, error)

	fct rpc.ClientStreamFactory
	// channel to convert multiple goroutine writeTask to single goroutine writeTask to FanOutQueue
	ch chan []byte

	chunk Chunk // buffer current writeTask metric for compress

	running *atomic.Bool

	// last flush time
	lastFlushTime time.Time
	// interval for check flush
	checkFlushInterval time.Duration
	// interval for flush
	flushInterval time.Duration

	lock4write sync.Mutex

	logger *logger.Logger
}

// newChannel returns a new channel with specific attribution.
func newChannel(
	ctx context.Context,
	database string,
	shardID models.ShardID,
	fct rpc.ClientStreamFactory,
) Channel {
	c := &channel{
		ctx:                ctx,
		database:           database,
		shardID:            shardID,
		newWriteStreamFn:   rpc.NewWriteStream,
		fct:                fct,
		ch:                 make(chan []byte, 2),
		running:            atomic.NewBool(false),
		checkFlushInterval: time.Second,
		chunk:              newChunk(defaultBufferSize), //TODO add config
		lastFlushTime:      time.Now(),
		logger:             logger.GetLogger("replica", "ShardChannel"),
	}

	return c
}

// Write writes the data into the channel, ErrCanceled is returned when the ctx is canceled before
// data is wrote successfully.
// Concurrent safe.
func (c *channel) Write(metric *protoMetricsV1.Metric) error {
	c.lock4write.Lock()
	defer c.lock4write.Unlock()

	if !c.running.Load() {
		return fmt.Errorf("shard write channle is not running")
	}
	c.chunk.Append(metric)

	if c.chunk.IsFull() {
		data, err := c.chunk.MarshalBinary()
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}
		select {
		case c.ch <- data:
			return nil
		case <-c.ctx.Done():
			return ErrCanceled
		}
	}
	return nil
}

func (c *channel) SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	if c.running.CAS(false, true) {
		//TODO check target exist???
		target := liveNodes[shardState.Leader]
		go c.writeTask(shardState, &target)
		c.logger.Info("start shard write channel successfully", logger.String("db", c.database),
			logger.Any("shardID", c.shardID))
	}
}

func (c *channel) Stop() {
	c.running.Store(false)
}

func (c *channel) writePendingBeforeClose() {
	// flush chunk pending data if chunk not empty
	if !c.chunk.IsEmpty() {
		// flush chunk pending data if chunk not empty
		c.flushChunk()
	}
	close(c.ch)
}

func (c *channel) checkFlush() {
	now := time.Now()
	if now.After(c.lastFlushTime.Add(c.flushInterval)) {
		c.lock4write.Lock()
		defer c.lock4write.Unlock()

		if !c.chunk.IsEmpty() {
			c.flushChunk()
		}
		c.lastFlushTime = now
	}
}

// flushChunk flushes the chunk data and appends data into queue
func (c *channel) flushChunk() {
	data, err := c.chunk.MarshalBinary()
	if err != nil {
		c.logger.Error("chunk marshal err", logger.Error(err))
		return
	}
	if len(data) == 0 {
		return
	}
	select {
	case c.ch <- data:
	case <-c.ctx.Done():
		c.logger.Warn("task has already canceled")
	}
}

// writeTask consumes data from chan, then appends the data into queue
func (c *channel) writeTask(shardState models.ShardState, target models.Node) {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(c.checkFlushInterval)
	defer ticker.Stop()

	var stream rpc.WriteStream
	var err error
	defer func() {
		if stream != nil {
			if err := stream.Close(); err != nil {
				c.logger.Error("close write stream err when exit write task", logger.Error(err))
			}
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.ch:
			if stream == nil {
				stream, err = c.newWriteStreamFn(c.ctx, target, c.database, &shardState, c.fct)
				if err != nil {
					//TODO do retry, add max retry count?
					//c.ch <- data
					continue
				}
			}
			if err := stream.Send(data); err == nil {
				putMarshalBlock(&data)
			} else {
				c.logger.Error("send write request err", logger.Error(err))
				if err == io.EOF {
					if err0 := stream.Close(); err0 != nil {
						c.logger.Error("close write stream err, when do write request", logger.Error(err))
					}
					stream = nil
				}
				//TODO do retry
				//c.ch <- data
			}
		case <-ticker.C:
			// check
			c.checkFlush()
		}
	}
}

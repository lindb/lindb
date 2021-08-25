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
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./shard_channel.go -destination=./shard_channel_mock.go -package=replica

var newSenderFn = newSender

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
	// Write writes the data into the channel, ErrCanceled is returned when the channel is canceled before
	// data is wrote successfully.
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
	// channel to convert multiple goroutine writeTask to single goroutine writeTask to FanOutQueue
	ch chan []byte

	chunk Chunk // buffer current writeTask metric for compress

	running *atomic.Bool
	sender  Sender

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
		sender:             newSenderFn(ctx, database, shardID, fct),
		database:           database,
		shardID:            shardID,
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
	c.sender.SyncShardState(shardState, liveNodes)
	if c.running.CAS(false, true) {
		go c.writeTask()
		c.logger.Info("start shard write channel successfully", logger.String("db", c.database),
			logger.Any("shardID", c.shardID))
	}
}

func (c *channel) Stop() {
	c.running.Store(false)
}

func (c *channel) writePendingBeforeClose() {
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		// try to drain data from chan
		for data := range c.ch {
			err := c.sender.Send(data)
			if err != nil {
				c.logger.Error("append to queue err", logger.Error(err))
			}
		}
		wait.Done()
	}()
	// flush chunk pending data if chunk not empty
	if !c.chunk.IsEmpty() {
		// flush chunk pending data if chunk not empty
		c.flushChunk()
	}
	close(c.ch)
	wait.Wait()
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
func (c *channel) writeTask() {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(c.checkFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.ch:
			err := c.sender.Send(data)
			if err != nil {
				c.logger.Error("append to queue err", logger.Error(err))
				//TODO do retry, add max retry count?
				c.ch <- data
			}
		case <-ticker.C:
			// check
			c.checkFlush()
		}
	}
}

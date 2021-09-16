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
	"io"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source=./channel_family.go -destination=./channel_family_mock.go -package=replica

type FamilyChannel interface {
	// Write writes the data into the channel,
	// ErrCanceled is returned when the channel is canceled before data is written successfully.
	// Concurrent safe.
	Write(ctx context.Context, rows []metric.BrokerRow) error
	Stop()
}

type familyChannel struct {
	// context to close channel
	ctx        context.Context
	cancel     context.CancelFunc
	database   string
	shardID    models.ShardID
	familyTime int64

	newWriteStreamFn func(
		ctx context.Context,
		target models.Node,
		database string, shardState *models.ShardState, familyTime int64,
		fct rpc.ClientStreamFactory,
	) (rpc.WriteStream, error)

	fct rpc.ClientStreamFactory

	// channel to convert multiple goroutine writeTask to single goroutine writeTask to FanOutQueue
	ch    chan *compressedChunk
	chunk Chunk // buffer current writeTask metric for compress

	lastFlushTime      time.Time     // last flush time
	checkFlushInterval time.Duration // interval for check flush
	batchTimout        time.Duration // interval for flush

	lock4write sync.Mutex
	logger     *logger.Logger
}

func newFamilyChannel(
	ctx context.Context,
	cfg config.Write,
	database string,
	shardID models.ShardID,
	familyTime int64,
	fct rpc.ClientStreamFactory,
	shardState models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) FamilyChannel {
	c, cancel := context.WithCancel(ctx)
	fc := &familyChannel{
		ctx:                c,
		cancel:             cancel,
		database:           database,
		shardID:            shardID,
		familyTime:         familyTime,
		fct:                fct,
		newWriteStreamFn:   rpc.NewWriteStream,
		ch:                 make(chan *compressedChunk, 2),
		checkFlushInterval: time.Second,
		batchTimout:        cfg.BatchTimeout.Duration(),
		chunk:              newChunk(cfg.BatchBlockSize),
		lastFlushTime:      time.Now(),
		logger:             logger.GetLogger("replica", "FamilyChannel"),
	}

	//TODO check target exist???

	target := liveNodes[shardState.Leader]
	go fc.writeTask(shardState, &target)

	return fc
}

// Write writes the data into the channel, ErrCanceled is returned when the ctx is canceled before
// data is written successfully.
// Concurrent safe.
func (fc *familyChannel) Write(ctx context.Context, rows []metric.BrokerRow) error {
	fc.lock4write.Lock()
	defer fc.lock4write.Unlock()

	for idx := 0; idx < len(rows); idx++ {
		if _, err := rows[idx].WriteTo(fc.chunk); err != nil {
			return err
		}

		if err := fc.flushChunkOnFull(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (fc *familyChannel) flushChunkOnFull(ctx context.Context) error {
	if !fc.chunk.IsFull() {
		return nil
	}
	compressed, err := fc.chunk.Compress()
	if err != nil {
		return err
	}

	select {
	case fc.ch <- compressed:
		return nil
	case <-ctx.Done(): // timeout of http ingestion api
		return ErrIngestTimeout
	case <-fc.ctx.Done():
		return ErrFamilyChannelCanceled
	}
}

// writeTask consumes data from chan, then appends the data into queue
func (fc *familyChannel) writeTask(shardState models.ShardState, target models.Node) {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(fc.checkFlushInterval)
	defer ticker.Stop()

	var stream rpc.WriteStream
	var err error
	defer func() {
		if stream != nil {
			if err := stream.Close(); err != nil {
				fc.logger.Error("close write stream err when exit write task", logger.Error(err))
			}
		}
	}()

	for {
		select {
		case <-fc.ctx.Done():
			return
		case compressed := <-fc.ch:
			if compressed == nil {
				// close chan
				continue
			}
			if stream == nil {
				//TODO need set transport.defaultMaxStreamsClient???
				stream, err = fc.newWriteStreamFn(fc.ctx, target, fc.database, &shardState, fc.familyTime, fc.fct)
				if err != nil {
					//TODO do retry, add max retry count?
					//c.ch <- data
					continue
				}
			}
			if err := stream.Send(*compressed); err == nil {
				compressed.Release()
			} else {
				fc.logger.Error(
					"failed writing compressed chunk to storage",
					logger.String("target", target.Indicator()),
					logger.String("database", fc.database),
					logger.Error(err))
				if err == io.EOF {
					if closeError := stream.Close(); closeError != nil {
						fc.logger.Error("failed closing write stream",
							logger.String("target", target.Indicator()),
							logger.Error(closeError))
					}
					stream = nil
				}
				//TODO do retry
				//c.ch <- data
			}
		case <-ticker.C:
			// check
			fc.checkFlush()
		}
	}
}

func (fc *familyChannel) writePendingBeforeClose() {
	// flush chunk pending data if chunk not empty
	if !fc.chunk.IsEmpty() {
		// flush chunk pending data if chunk not empty
		fc.flushChunk()
	}
	close(fc.ch)
}

func (fc *familyChannel) checkFlush() {
	now := time.Now()
	if now.After(fc.lastFlushTime.Add(fc.batchTimout)) {
		fc.lock4write.Lock()
		defer fc.lock4write.Unlock()

		if !fc.chunk.IsEmpty() {
			fc.flushChunk()
		}
		fc.lastFlushTime = now
	}
}

func (fc *familyChannel) Stop() {
	fc.cancel()
}

// flushChunk flushes the chunk data and appends data into queue
func (fc *familyChannel) flushChunk() {
	compressed, err := fc.chunk.Compress()
	if err != nil {
		fc.logger.Error("chunk marshal err", logger.Error(err))
		return
	}
	if compressed == nil || len(*compressed) == 0 {
		return
	}
	select {
	case fc.ch <- compressed:
	case <-fc.ctx.Done():
		fc.logger.Warn("writer is canceled")
	}
}

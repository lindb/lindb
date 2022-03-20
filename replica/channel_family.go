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
	"github.com/lindb/lindb/pkg/fasttime"
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
	// leaderChanged notifies family channel need change leader send stream
	leaderChanged(shardState models.ShardState,
		liveNodes map[models.NodeID]models.StatefulNode)
	// Stop stops the family channel.
	Stop()
	FamilyTime() int64
	isExpire(ahead, behind int64) bool
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

	fct           rpc.ClientStreamFactory
	shardState    models.ShardState
	liveNodes     map[models.NodeID]models.StatefulNode
	currentTarget models.Node

	// channel to convert multiple goroutine writeTask to single goroutine writeTask to FanOutQueue
	ch                  chan *compressedChunk
	leaderChangedSignal chan struct{}
	chunk               Chunk // buffer current writeTask metric for compress

	lastFlushTime      time.Time     // last flush time
	checkFlushInterval time.Duration // interval for check flush
	batchTimout        time.Duration // interval for flush

	lock4write sync.Mutex
	lock4meta  sync.Mutex
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
		ctx:                 c,
		cancel:              cancel,
		database:            database,
		shardID:             shardID,
		familyTime:          familyTime,
		fct:                 fct,
		shardState:          shardState,
		liveNodes:           liveNodes,
		newWriteStreamFn:    rpc.NewWriteStream,
		ch:                  make(chan *compressedChunk, 2),
		leaderChangedSignal: make(chan struct{}),
		checkFlushInterval:  time.Second,
		batchTimout:         cfg.BatchTimeout.Duration(),
		chunk:               newChunk(cfg.BatchBlockSize),
		lastFlushTime:       time.Now(),
		logger:              logger.GetLogger("replica", "FamilyChannel"),
	}
	go fc.writeTask()

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

// leaderChanged notifies family channel need change leader send stream
func (fc *familyChannel) leaderChanged(
	shardState models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) {
	fc.lock4meta.Lock()
	fc.shardState = shardState
	fc.liveNodes = liveNodes
	fc.lock4meta.Unlock()
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
func (fc *familyChannel) writeTask() {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(fc.checkFlushInterval)
	defer ticker.Stop()

	retryBuffers := make([]*compressedChunk, 0)
	retry := func(compressed *compressedChunk) {
		// TODO add config
		if len(retryBuffers) > 100 {
			fc.logger.Error("too many retry messages, drop current message")
		} else {
			retryBuffers = append(retryBuffers, compressed)
		}
	}
	send := func(stream rpc.WriteStream, compressed *compressedChunk) bool {
		if stream == nil {
			return false
		}
		if err := stream.Send(*compressed); err != nil {
			fc.logger.Error(
				"failed writing compressed chunk to storage",
				logger.String("target", fc.currentTarget.Indicator()),
				logger.String("database", fc.database),
				logger.Error(err))
			if err == io.EOF {
				if closeError := stream.Close(); closeError != nil {
					fc.logger.Error("failed closing write stream",
						logger.String("target", fc.currentTarget.Indicator()),
						logger.Error(closeError))
				}
			}
			// retry if err
			retry(compressed)
			return false
		}
		compressed.Release()
		return true
	}

	var stream rpc.WriteStream
	defer func() {
		if stream != nil {
			if err := stream.Close(); err != nil {
				fc.logger.Error("close write stream err when exit write task", logger.Error(err))
			}
		}
	}()
	var err error

	for {
		select {
		case <-fc.ctx.Done():
			return
		case <-fc.leaderChangedSignal:
			fc.logger.Info("shard leader changed, need switch send stream",
				logger.String("oldTarget", fc.currentTarget.Indicator()),
				logger.String("database", fc.database))
			if stream != nil {
				// if stream isn't nil, need close old stream first.
				if err = stream.Close(); err != nil {
					fc.logger.Error("close write stream err when leader changed", logger.Error(err))
				}
				stream = nil
			}
		case compressed := <-fc.ch:
			if compressed == nil {
				// close chan
				continue
			}
			if stream == nil {
				fc.lock4meta.Lock()
				leader := fc.liveNodes[fc.shardState.Leader]
				shardState := fc.shardState
				fc.currentTarget = &leader
				fc.lock4meta.Unlock()
				stream, err = fc.newWriteStreamFn(fc.ctx, fc.currentTarget, fc.database, &shardState, fc.familyTime, fc.fct)
				if err != nil {
					retry(compressed)
					continue
				}
			}
			if send(stream, compressed) {
				// if send ok, do pending retry message
				if len(retryBuffers) > 0 {
					messages := retryBuffers
					retryBuffers = make([]*compressedChunk, 0)
					fail := false
					for _, msg := range messages {
						if fail {
							retry(msg)
							continue
						}
						if !send(stream, msg) {
							fail = true
						}
					}
				}
			} else {
				stream = nil
			}
		case <-ticker.C:
			// check
			fc.checkFlush()
		}
	}

	// TODO write pending check after stop???
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

func (fc *familyChannel) isExpire(ahead, _ int64) bool {
	now := fasttime.UnixMilliseconds()
	// add 15 minute buffer
	if ahead > 0 && fc.lastFlushTime.Unix()*1000+ahead+15*time.Minute.Milliseconds() > now {
		return false
	}
	return true
}

func (fc *familyChannel) FamilyTime() int64 {
	return fc.familyTime
}

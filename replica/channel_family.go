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
	"runtime/pprof"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source=./channel_family.go -destination=./channel_family_mock.go -package=replica

// FamilyChannel represents family write shardChannel.
type FamilyChannel interface {
	// Write writes the data into the shardChannel,
	// ErrCanceled is returned when the shardChannel is canceled before data is written successfully.
	// Concurrent safe.
	Write(ctx context.Context, rows []metric.BrokerRow) error
	// leaderChanged notifies family shardChannel need change leader send stream
	leaderChanged(shardState models.ShardState,
		liveNodes map[models.NodeID]models.StatefulNode)
	// Stop stops the family shardChannel.
	Stop(timeout int64)
	// FamilyTime returns the family time of current shardChannel.
	FamilyTime() int64
	// isExpire returns if current family is expired.
	isExpire(ahead, behind int64) bool
}

// familyChannel implements FamilyChannel interface.
type familyChannel struct {
	// context to close shardChannel
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

	// shardChannel to convert multiple goroutine writeTask to single goroutine writeTask to FanOutQueue
	ch                  chan *compressedChunk
	leaderChangedSignal chan struct{}
	stoppedSignal       chan struct{}
	stoppingSignal      chan struct{}
	chunk               Chunk // buffer current writeTask metric for compress

	lastFlushTime      *atomic.Int64 // last flush time
	checkFlushInterval time.Duration // interval for check flush
	batchTimeout       time.Duration // interval for flush
	maxRetryBuf        int

	lock4write sync.Mutex
	lock4meta  sync.Mutex

	statistics *metrics.BrokerFamilyWriteStatistics
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
		leaderChangedSignal: make(chan struct{}, 1),
		stoppedSignal:       make(chan struct{}, 1),
		stoppingSignal:      make(chan struct{}, 1),
		checkFlushInterval:  time.Second,
		batchTimeout:        cfg.BatchTimeout.Duration(),
		maxRetryBuf:         100, // TODO add config
		chunk:               newChunk(cfg.BatchBlockSize),
		lastFlushTime:       atomic.NewInt64(timeutil.Now()),
		statistics:          metrics.NewBrokerFamilyWriteStatistics(database),
		logger:              logger.GetLogger("replica", "FamilyChannel"),
	}

	fc.statistics.ActiveWriteFamilies.Incr()

	go func() {
		channelFamilyLabels := pprof.Labels("database", database,
			"shard", shardID.String(), "family", timeutil.FormatTimestamp(familyTime, timeutil.DataTimeFormat2))
		pprof.Do(c, channelFamilyLabels, fc.writeTask)
	}()

	return fc
}

// Write writes the data into the shardChannel, ErrCanceled is returned when the ctx is canceled before
// data is written successfully.
// Concurrent safe.
func (fc *familyChannel) Write(ctx context.Context, rows []metric.BrokerRow) error {
	total := len(rows)
	success := 0

	fc.lock4write.Lock()
	defer func() {
		if total > 0 {
			fc.statistics.BatchMetrics.Add(float64(success))
			fc.statistics.BatchMetricFailures.Add(float64(total - success))
		}
		fc.lock4write.Unlock()
	}()

	for idx := 0; idx < total; idx++ {
		if _, err := rows[idx].WriteTo(fc.chunk); err != nil {
			return err
		}

		if err := fc.flushChunkOnFull(ctx); err != nil {
			return err
		}
		success++
	}

	return nil
}

// leaderChanged notifies family shardChannel need change leader send stream
func (fc *familyChannel) leaderChanged(
	shardState models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) {
	fc.lock4meta.Lock()
	fc.shardState = shardState
	fc.liveNodes = liveNodes
	fc.lock4meta.Unlock()

	fc.leaderChangedSignal <- struct{}{}
	fc.statistics.LeaderChanged.Incr()
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
	case <-ctx.Done(): // timeout of http ingestion api
		return ErrIngestTimeout
	case <-fc.ctx.Done():
		return ErrFamilyChannelCanceled
	case fc.ch <- compressed:
		fc.lastFlushTime.Store(timeutil.Now())
		return nil
	}
}

// writeTask consumes data from chan, then appends the data into queue
func (fc *familyChannel) writeTask(_ context.Context) {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(fc.checkFlushInterval)
	defer ticker.Stop()

	retryBuffers := make([]*compressedChunk, 0)
	retry := func(compressed *compressedChunk) {
		if len(retryBuffers) > fc.maxRetryBuf {
			fc.logger.Error("too many retry messages, drop current message")
			fc.statistics.RetryDrop.Incr()
		} else {
			retryBuffers = append(retryBuffers, compressed)
			fc.statistics.Retry.Incr()
		}
	}
	var stream rpc.WriteStream
	send := func(compressed *compressedChunk) bool {
		if compressed == nil {
			return true
		}
		if len(*compressed) == 0 {
			compressed.Release()
			return true
		}
		if stream == nil {
			fc.lock4meta.Lock()
			leader := fc.liveNodes[fc.shardState.Leader]
			shardState := fc.shardState
			fc.currentTarget = &leader
			fc.lock4meta.Unlock()
			s, err := fc.newWriteStreamFn(fc.ctx, fc.currentTarget, fc.database, &shardState, fc.familyTime, fc.fct)
			if err != nil {
				fc.statistics.CreateStreamFailures.Incr()
				retry(compressed)
				return false
			}
			fc.statistics.CreateStream.Incr()
			stream = s
		}
		if err := stream.Send(*compressed); err != nil {
			fc.statistics.SendFailure.Incr()
			fc.logger.Error(
				"failed writing compressed chunk to storage",
				logger.String("target", fc.currentTarget.Indicator()),
				logger.String("database", fc.database),
				logger.Error(err))
			if err == io.EOF {
				if closeError := stream.Close(); closeError != nil {
					fc.statistics.CloseStreamFailures.Incr()
					fc.logger.Error("failed closing write stream",
						logger.String("target", fc.currentTarget.Indicator()),
						logger.Error(closeError))
				} else {
					fc.statistics.CloseStream.Incr()
				}
				stream = nil
			}
			// retry if err
			retry(compressed)
			return false
		}
		fc.statistics.SendSuccess.Incr()
		fc.statistics.SendSize.Add(float64(len(*compressed)))
		fc.statistics.PendingSend.Decr()
		compressed.Release()
		return true
	}

	defer func() {
		if stream != nil {
			if err := stream.Close(); err != nil {
				fc.statistics.CloseStreamFailures.Incr()
				fc.logger.Error("close write stream err when exit write task", logger.Error(err))
			} else {
				fc.statistics.CloseStream.Incr()
			}
		}
	}()

	// send pending in buffer before stop channel.
	sendBeforeStop := func() {
		defer func() {
			fc.stoppedSignal <- struct{}{}
		}()
		sendLastMsg := func(compressed *compressedChunk) {
			if !send(compressed) {
				fc.logger.Error("send message failure before close channel, message lost")
			}
		}
		// flush chunk pending data if chunk not empty
		if !fc.chunk.IsEmpty() {
			// flush chunk pending data if chunk not empty
			compressed, err0 := fc.chunk.Compress()
			if err0 != nil {
				fc.logger.Error("compress chunk err when send last chunk data", logger.Error(err0))
			} else {
				sendLastMsg(compressed)
			}
		}
		// try to write pending data
		for compressed := range fc.ch {
			sendLastMsg(compressed)
		}
	}
	var err error
	for {
		select {
		case <-fc.stoppingSignal:
			sendBeforeStop()
			return
		case <-fc.ctx.Done():
			sendBeforeStop()
			return
		case <-fc.leaderChangedSignal:
			if stream != nil {
				fc.logger.Info("shard leader changed, need switch send stream",
					logger.String("oldTarget", fc.currentTarget.Indicator()),
					logger.String("database", fc.database))
				// if stream isn't nil, need close old stream first.
				if err = stream.Close(); err != nil {
					fc.logger.Error("close write stream err when leader changed", logger.Error(err))
				}
				stream = nil
			}
		case compressed := <-fc.ch:
			if send(compressed) {
				// if send ok, retry pending message
				if len(retryBuffers) > 0 {
					messages := retryBuffers
					retryBuffers = make([]*compressedChunk, 0)
					for _, msg := range messages {
						if !send(msg) {
							retry(msg)
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
}

// checkFlush checks if channel needs to flush data.
func (fc *familyChannel) checkFlush() {
	now := timeutil.Now()
	if now-fc.lastFlushTime.Load() >= fc.batchTimeout.Milliseconds() {
		fc.lock4write.Lock()
		defer fc.lock4write.Unlock()

		if !fc.chunk.IsEmpty() {
			fc.flushChunk()
			fc.lastFlushTime.Store(now)
		}
	}
}

// Stop stops current write family shardChannel.
func (fc *familyChannel) Stop(timeout int64) {
	close(fc.stoppingSignal)
	close(fc.ch)

	ticker := time.NewTicker(time.Duration(time.Millisecond.Nanoseconds() * timeout))
	select {
	case <-fc.stoppedSignal:
	case <-ticker.C:
	}
	fc.cancel()

	fc.statistics.ActiveWriteFamilies.Decr()
}

// flushChunk flushes the chunk data and appends data into queue
func (fc *familyChannel) flushChunk() {
	compressed, err := fc.chunk.Compress()
	if err != nil {
		fc.logger.Error("compress chunk err", logger.Error(err))
		return
	}
	if compressed == nil || len(*compressed) == 0 {
		return
	}
	select {
	case fc.ch <- compressed:
		fc.statistics.PendingSend.Incr()
	case <-fc.ctx.Done():
		fc.logger.Warn("writer is canceled")
	}
}

// isExpire returns if current family is expired.
func (fc *familyChannel) isExpire(ahead, _ int64) bool {
	now := timeutil.Now()
	fc.logger.Info("family channel expire check",
		logger.String("database", fc.database),
		logger.Any("shard", fc.shardID),
		logger.Int64("head", ahead),
		logger.String("family", timeutil.FormatTimestamp(fc.lastFlushTime.Load(), timeutil.DataTimeFormat2)))
	// add 15 minute buffer
	return fc.lastFlushTime.Load()+ahead+15*time.Minute.Milliseconds() < now
}

// FamilyTime returns the family time of current shardChannel.
func (fc *familyChannel) FamilyTime() int64 {
	return fc.familyTime
}

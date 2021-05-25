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

package tsdb

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/mem"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./data_flush_checker.go -destination=./data_flush_checker_mock.go -package=tsdb

var (
	// can be modified in runtime
	memoryUsageCheckInterval = *atomic.NewDuration(time.Second)
)

// DataFlushChecker represents the memory database flush checker.
// There are 4 flush policies of the Engine as below:
// 1. FullFlush
//    highest priority, triggered by external API from the users.
//    this action will blocks any other flush checkers.
// 2. GlobalMemoryUsageChecker
//    This checker will check the global memory usage of the host periodically,
//    when the metric is above MemoryHighWaterMark, a `watermarkFlusher` will be spawned
//    whose responsibility is to flush the biggest shard until memory is lower than  MemoryLowWaterMark.
// 3. ShardMemoryUsageChecker
//    This checker will check each shard's memory usage periodically,
//    If this shard is above ShardMemoryUsedThreshold. it will be flushed to disk.
// 4. DatabaseMetaFlusher
//    It is a simple checker which flush the meta of database to disk periodically.
//
// a). Each shard or database is restricted to flush by one goroutine at the same time via CAS operation;
// b). The flush workers runs concurrently;
// c). All unit will be flushed when closing;
type DataFlushChecker interface {
	// Start starts the checker goroutine in background
	Start()
	// Stop stops the background check goroutine
	Stop()

	// requestFlushJob requests a flush job for the spec shard
	requestFlushJob(shard Shard, global bool)
}

// flushRequest represents the shard flush job request
type flushRequest struct {
	shard  Shard
	global bool // above high memory watermark
}

// dataFlushChecker implements DataFlushCheck interface
type dataFlushChecker struct {
	ctx    context.Context
	cancel context.CancelFunc

	shardInFlushing      sync.Map
	flushRequestCh       chan *flushRequest          // shard to flush
	flushInFlight        atomic.Int32                // current pending in flushing
	isWatermarkFlushing  atomic.Bool                 // this flag symbols if has goroutine in high water-mark flushing
	memoryStatGetterFunc monitoring.MemoryStatGetter // used for mocking
}

// newDataFlushChecker creates the data flush checker
func newDataFlushChecker(ctx context.Context) DataFlushChecker {
	c, cancel := context.WithCancel(ctx)
	return &dataFlushChecker{
		ctx:                  c,
		cancel:               cancel,
		flushRequestCh:       make(chan *flushRequest),
		memoryStatGetterFunc: mem.VirtualMemory,
	}
}

// Start starts the checker goroutine in background
func (fc *dataFlushChecker) Start() {
	go fc.startCheckDataFlush()
}

// Stop stops the background check goroutine
func (fc *dataFlushChecker) Stop() {
	fc.cancel()
}

// startCheckDataFlush starts check memory usage for each memory database under shard
func (fc *dataFlushChecker) startCheckDataFlush() {
	// 1. start timer
	timer := time.NewTimer(memoryUsageCheckInterval.Load())
	defer timer.Stop()

	// 2. start some flush workers
	for i := 0; i < constants.FlushConcurrency; i++ {
		go fc.flushWorker()
	}

	for {
		select {
		case <-fc.ctx.Done():
			return
		case <-timer.C:
			// check each shard if need do flush job
			GetShardManager().WalkEntry(func(shard Shard) {
				if shard.NeedFlush() {
					fc.requestFlushJob(shard, false)
				}
			})
			if fc.flushInFlight.Load() == 0 {
				// check Global memory is above than the high watermark
				stat, _ := fc.memoryStatGetterFunc()
				if stat.UsedPercent > constants.MemoryHighWaterMark &&
					!fc.isWatermarkFlushing.Load() {
					// memory is higher than the high-watermark
					// restrict watermarkFlusher concurrency thread-safe
					fc.flushBiggestMemoryUsageShard()
				}
			}
			// reset check interval
			timer.Reset(memoryUsageCheckInterval.Load())
		}
	}
}

// requestFlushJob requests a flush job for the spec shard
func (fc *dataFlushChecker) requestFlushJob(shard Shard, global bool) {
	_, ok := fc.shardInFlushing.Load(shard.ShardInfo())
	if ok {
		// if shard is in flushing queue, returns it
		return
	}
	fc.shardInFlushing.Store(shard.ShardInfo(), shard)
	select {
	case <-fc.ctx.Done():
		return
	case fc.flushRequestCh <- &flushRequest{shard: shard, global: global}:
		// add count of flush in flight
		fc.flushInFlight.Inc()
	}
}

// flushWorker consumes the flush request from chan
func (fc *dataFlushChecker) flushWorker() {
	for {
		select {
		case <-fc.ctx.Done():
			return
		case request := <-fc.flushRequestCh:
			// do flush job
			fc.doFlush(request)
		}
	}
}

// doFlush does the flush job for the spec shard
func (fc *dataFlushChecker) doFlush(request *flushRequest) {
	shard := request.shard
	global := request.global
	shardInfo := shard.ShardInfo()
	defer func() {
		if global {
			fc.isWatermarkFlushing.Store(false)
		}
		fc.flushInFlight.Dec()
		// delete shard from flushing queue
		fc.shardInFlushing.Delete(shardInfo)
	}()

	if global {
		fc.isWatermarkFlushing.Store(true)
	}
	if err := shard.Flush(); err != nil {
		//TODO add metric
		engineLogger.Error("flush shard memory database error",
			logger.String("shard", shardInfo), logger.Error(err))
	}
}

// flushBiggestMemoryUsageShard picks the biggest memory usage shard to flush
func (fc *dataFlushChecker) flushBiggestMemoryUsageShard() {
	var (
		biggestShard   Shard
		biggestMemSize int32
	)
	GetShardManager().WalkEntry(func(shard Shard) {
		// skip shard in flushing
		if shard.IsFlushing() {
			return
		}

		//FIXME(stone1100)
		theShardSize := int32(1024)
		//shard.MemoryDatabase().MemSize()
		if theShardSize > biggestMemSize {
			// pick a shard that has biggest memory size
			biggestMemSize = theShardSize
			biggestShard = shard
		}
	})
	if biggestMemSize == 0 {
		return
	}
	// request flush job
	fc.requestFlushJob(biggestShard, true)
}

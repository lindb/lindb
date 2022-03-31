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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
)

//go:generate mockgen -source=./data_flush_checker.go -destination=./data_flush_checker_mock.go -package=tsdb

var (
	// can be modified in runtime
	memoryUsageCheckInterval = *atomic.NewDuration(time.Second)
	ignoreMemorySize         = ltoml.Size(4 * 1024 * 1024)
)

// DataFlushChecker represents the memory database flush checker.
// There are 4 flush policies of the Engine as below:
// 1. FullFlush
//    the highest priority, triggered by external API from the users.
//    this action will block any other flush checkers.
// 2. GlobalMemoryUsageChecker
//    This checker will check the global memory usage of the host periodically,
//    when the metric is above MemoryHighWaterMark, a `watermarkFlusher` will be spawned
//    whose responsibility is to flush the biggest family until memory is lower than  MemoryLowWaterMark.
// 3. FamilyMemoryUsageChecker
//    This checker will check each family's memory usage periodically,
//    If this family is above FamilyMemoryUsedThreshold. it will be flushed to disk.
// 4. DatabaseMetaFlusher
//    It is a simple checker which flush the meta of database to disk periodically.
//
// a). Each family or database is restricted to flush by one goroutine at the same time via CAS operation;
// b). The flush workers runs concurrently;
// c). All unit will be flushed when closing;
type DataFlushChecker interface {
	// Start starts the checker goroutine in background.
	Start()
	// Stop stops the background check goroutine.
	Stop()

	// requestFlushJob requests a flush job for the spec shard/families.
	requestFlushJob(request *flushRequest)
}

// flushRequest represents the families flush job request
type flushRequest struct {
	shard    Shard
	families []DataFamily
	global   bool // above high memory watermark
}

// dataFlushChecker implements DataFlushCheck interface
type dataFlushChecker struct {
	ctx    context.Context
	cancel context.CancelFunc

	shardInFlushing      sync.Map
	flushRequestCh       chan *flushRequest          // family to flush
	flushInFlight        atomic.Int32                // current pending in flushing
	isWatermarkFlushing  atomic.Bool                 // this flag symbols if has goroutine in high water-mark flushing
	memoryStatGetterFunc monitoring.MemoryStatGetter // used for mocking
	running              *atomic.Bool
	logger               *logger.Logger
}

// newDataFlushChecker creates the data flush checker
func newDataFlushChecker(ctx context.Context) DataFlushChecker {
	c, cancel := context.WithCancel(ctx)
	return &dataFlushChecker{
		ctx:                  c,
		cancel:               cancel,
		flushRequestCh:       make(chan *flushRequest),
		memoryStatGetterFunc: mem.VirtualMemory,
		running:              atomic.NewBool(false),
		logger:               engineLogger,
	}
}

// Start starts the checker goroutine in background
func (fc *dataFlushChecker) Start() {
	if fc.running.CAS(false, true) {
		go fc.startCheckDataFlush()
	}
}

// Stop stops the background check goroutine
func (fc *dataFlushChecker) Stop() {
	if fc.running.CAS(true, false) {
		fc.cancel()
	}
}

// startCheckDataFlush starts check memory usage for each memory database under family
func (fc *dataFlushChecker) startCheckDataFlush() {
	// 1. start timer
	timer := time.NewTimer(memoryUsageCheckInterval.Load())
	defer timer.Stop()

	// 2. start some flush workers
	for i := 0; i < config.GlobalStorageConfig().TSDB.FlushConcurrency; i++ {
		go fc.flushWorker()
	}
	fc.logger.Info("DataFlusher Checker is running",
		logger.Int32("workers", int32(config.GlobalStorageConfig().TSDB.FlushConcurrency)))

	for {
		select {
		case <-fc.ctx.Done():
			return
		case <-timer.C:
			needFlushShards := make(map[string]*flushRequest)
			// check each family if it needs to do flush job
			GetFamilyManager().WalkEntry(func(family DataFamily) {
				if family.NeedFlush() {
					shard := family.Shard()
					needFlushShard, ok := needFlushShards[shard.Indicator()]
					if !ok {
						needFlushShards[shard.Indicator()] = &flushRequest{
							shard:    shard,
							families: []DataFamily{family},
							global:   false,
						}
					} else {
						needFlushShard.families = append(needFlushShard.families, family)
					}
				}
			})

			for _, request := range needFlushShards {
				fc.requestFlushJob(request)
			}

			if len(needFlushShards) == 0 && !fc.isWatermarkFlushing.Load() && fc.flushInFlight.Load() == 0 {
				// check Global memory is above than the high watermark, if no shard need flush
				stat, _ := fc.memoryStatGetterFunc()
				maxMemUsageLimit := config.GlobalStorageConfig().TSDB.MaxMemUsageBeforeFlush * 100
				if stat.UsedPercent > maxMemUsageLimit {
					// memory is higher than the high-watermark
					// restrict watermarkFlusher concurrency thread-safe
					fc.logger.Info("memory is higher than the high watermark, need pick biggest memory usage family to flush",
						logger.Any("memoryUsed", stat.UsedPercent),
						logger.Any("limit", maxMemUsageLimit))
					fc.flushBiggestMemoryUsageFamily()
				}
			}
			// reset check interval
			timer.Reset(memoryUsageCheckInterval.Load())
		}
	}
}

// requestFlushJob requests a flush job for the spec shard/families.
func (fc *dataFlushChecker) requestFlushJob(request *flushRequest) {
	if !fc.running.Load() {
		return
	}
	_, ok := fc.shardInFlushing.Load(request.shard.Indicator())
	if ok {
		// if shard is in flushing queue, returns it
		return
	}
	fc.shardInFlushing.Store(request.shard.Indicator(), request)
	select {
	case <-fc.ctx.Done():
		return
	case fc.flushRequestCh <- request:
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
			if request != nil {
				// do flush job
				fc.doFlush(request)
			}
		}
	}
}

// doFlush does the flush job for the spec family.
func (fc *dataFlushChecker) doFlush(request *flushRequest) {
	indicator := request.shard.Indicator()
	defer func() {
		// after flush, try garbage collect(write buffer)
		request.shard.BufferManager().GarbageCollect()

		if request.global {
			fc.isWatermarkFlushing.Store(false)
		}
		fc.flushInFlight.Dec()
		// delete family from flushing queue
		fc.shardInFlushing.Delete(indicator)
	}()

	if request.global {
		fc.isWatermarkFlushing.Store(true)
	}
	// TODO
	_ = request.shard.Flush()

	// TODO add flush timeout?
	for _, family := range request.families {
		if err := family.Flush(); err != nil {
			engineLogger.Error("flush family memory database error",
				logger.String("family", family.Indicator()), logger.Error(err))
		}
	}
}

// flushBiggestMemoryUsageFamily picks the biggest memory usage family to flush
func (fc *dataFlushChecker) flushBiggestMemoryUsageFamily() {
	var (
		biggestFamily DataFamily
		// ignore family whose memdb size is smaller than 4MB
		// without this threshold, when tsdb is under insufficient system resources,
		// watermark flushing may create a large number of small L0 files，
		// which is not helpful for reducing system memory usage
		biggestMemSize = ignoreMemorySize
	)
	GetFamilyManager().WalkEntry(func(family DataFamily) {
		// skip family in flushing
		if family.IsFlushing() {
			return
		}
		thisFamilyMemDBSize := ltoml.Size(family.MemDBSize())
		if thisFamilyMemDBSize > biggestMemSize {
			biggestMemSize = thisFamilyMemDBSize
			biggestFamily = family
		}
	})
	// no available family for flushing
	if biggestFamily == nil {
		return
	}
	// request flush job
	fc.requestFlushJob(&flushRequest{
		shard:    biggestFamily.Shard(),
		families: []DataFamily{biggestFamily},
		global:   true,
	})
}

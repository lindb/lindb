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
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/mem"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
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

// flushShard represents the shard flush job request
type flushShard struct {
	shard    Shard
	families []DataFamily
}

// flushRequest represents the families flush job request
type flushRequest struct {
	db     Database
	shards map[models.ShardID]*flushShard
	global bool // above high memory watermark
}

// dataFlushChecker implements DataFlushCheck interface
type dataFlushChecker struct {
	ctx    context.Context
	cancel context.CancelFunc

	dbInFlushing         sync.Map           // database name => flush request
	flushRequestCh       chan *flushRequest // family to flush
	flushInFlight        atomic.Int32       // current pending in flushing
	isWatermarkFlushing  atomic.Bool        // this flag symbols if it has goroutine in high water-mark flushing
	running              *atomic.Bool
	memoryStatGetterFunc monitoring.MemoryStatGetter // used for mocking
	logger               *logger.Logger
}

// newDataFlushChecker creates the data flush checker
func newDataFlushChecker(ctx context.Context) DataFlushChecker {
	c, cancel := context.WithCancel(ctx)
	return &dataFlushChecker{
		ctx:                  c,
		cancel:               cancel,
		flushRequestCh:       make(chan *flushRequest, 8),
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
			fc.check()
			// reset check interval
			timer.Reset(memoryUsageCheckInterval.Load())
		}
	}
}

// check finds family which need flush data.
func (fc *dataFlushChecker) check() {
	needFlushDBs := make(map[string]*flushRequest)
	// check each family if it needs to do flush job
	GetFamilyManager().WalkEntry(func(family DataFamily) {
		if family.NeedFlush() {
			shard := family.Shard()
			dbName := shard.Database().Name()
			needFlushDB, ok := needFlushDBs[dbName]
			if !ok {
				needFlushDB = &flushRequest{
					db:     shard.Database(),
					shards: make(map[models.ShardID]*flushShard),
					global: false,
				}
				needFlushDBs[dbName] = needFlushDB
			}
			needFlushShard, ok := needFlushDB.shards[shard.ShardID()]
			if !ok {
				needFlushShard = &flushShard{
					shard: shard,
				}
				needFlushDB.shards[shard.ShardID()] = needFlushShard
			}
			needFlushShard.families = append(needFlushShard.families, family)
			metrics.ShardStatistics.FlushInFlight.WithTagValues(dbName, strconv.Itoa(int(shard.ShardID()))).Incr()
		}
	})

	for _, request := range needFlushDBs {
		fc.requestFlushJob(request)
	}

	if len(needFlushDBs) == 0 && !fc.isWatermarkFlushing.Load() && fc.flushInFlight.Load() == 0 {
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
}

// requestFlushJob requests a flush job for the spec shard/families.
func (fc *dataFlushChecker) requestFlushJob(request *flushRequest) {
	if !fc.running.Load() {
		return
	}
	_, ok := fc.dbInFlushing.Load(request.db.Name())
	if ok {
		// if shard is in flushing queue, returns it
		return
	}
	select {
	case <-fc.ctx.Done():
		return
	case fc.flushRequestCh <- request:
		fc.dbInFlushing.Store(request.db.Name(), request)
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
	indicator := request.db.Name()
	defer func() {
		if request.global {
			fc.isWatermarkFlushing.Store(false)
		}
		fc.flushInFlight.Dec()
		// delete family from flushing queue
		fc.dbInFlushing.Delete(indicator)
	}()

	if request.global {
		fc.isWatermarkFlushing.Store(true)
	}
	// flush data step:
	// 1. flush database metadata(metric/tag/field) if it needs
	// 2. flush index database for each shard if it needs
	// 3. flush family data
	if err := request.db.FlushMeta(); err != nil {
		engineLogger.Error("flush database metadata error",
			logger.String("database", request.db.Name()), logger.Error(err))
		return
	}
	// wait metadata flush job completed, maybe other goroutine is flushing.
	request.db.WaitFlushMetaCompleted()
	// flush each shard
	for shardID := range request.shards {
		shardReq := request.shards[shardID]
		fc.flushShard(shardReq)
	}
}

// flushShard flushes index data and family metric data.
func (fc *dataFlushChecker) flushShard(request *flushShard) {
	// after flush, try garbage collect(write buffer)
	defer request.shard.BufferManager().GarbageCollect()

	if err := request.shard.FlushIndex(); err != nil {
		engineLogger.Error("flush shard index memory database error",
			logger.String("shard", request.shard.Indicator()), logger.Error(err))
		return
	}
	// wait index flush job completed, maybe other goroutine is flushing.
	request.shard.WaitFlushIndexCompleted()

	// TODO add flush timeout?
	for _, family := range request.families {
		if err := family.Flush(); err != nil {
			engineLogger.Error("flush family memory database error",
				logger.String("family", family.Indicator()), logger.Error(err))
		}
		metrics.ShardStatistics.FlushInFlight.
			WithTagValues(request.shard.Database().Name(), strconv.Itoa(int(request.shard.ShardID()))).Decr()
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
	shard := biggestFamily.Shard()
	// request flush job
	fc.requestFlushJob(&flushRequest{
		db: shard.Database(),
		shards: map[models.ShardID]*flushShard{
			shard.ShardID(): {
				shard:    biggestFamily.Shard(),
				families: []DataFamily{biggestFamily},
			},
		},
		global: true,
	})
}

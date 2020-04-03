package tsdb

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
)

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

	shardInFlushing      map[string]Shard
	flushRequestCh       chan *flushRequest          // shard to flush
	flushInFlight        atomic.Int32                // current pending in flushing
	isWatermarkFlushing  atomic.Bool                 // this flag symbols if has goroutine in high water-mark flushing
	memoryStatGetterFunc monitoring.MemoryStatGetter // used for mocking
	mutex                sync.RWMutex
}

// newDataFlushChecker creates the data flush checker
func newDataFlushChecker(ctx context.Context) DataFlushChecker {
	c, cancel := context.WithCancel(ctx)
	return &dataFlushChecker{
		ctx:                  c,
		cancel:               cancel,
		flushRequestCh:       make(chan *flushRequest),
		shardInFlushing:      make(map[string]Shard),
		memoryStatGetterFunc: monitoring.GetMemoryStat,
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
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	shard := request.shard
	global := request.global
	shardInfo := shard.ShardInfo()
	_, ok := fc.shardInFlushing[shardInfo]
	if !ok {
		// if shard is not in flushing queue, then does flush job
		defer func() {
			if global {
				fc.isWatermarkFlushing.Store(false)
			}
			fc.flushInFlight.Dec()
			// delete shard from flushing queue
			delete(fc.shardInFlushing, shardInfo)
		}()

		if global {
			fc.isWatermarkFlushing.Store(true)
		}
		if err := shard.Flush(); err != nil {
			//TODO add metric
			engineLogger.Error("flush shard memory database error",
				logger.String("shard", shardInfo), logger.Error(err))
		}
		fc.shardInFlushing[shard.ShardInfo()] = shard
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

		theShardSize := shard.MemoryDatabase().MemSize()
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

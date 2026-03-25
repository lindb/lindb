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

package flush

import (
	"sort"
	"time"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/metrics"
)

// FlushReason indicates why a flush was triggered.
type FlushReason int

const (
	FlushReasonNone         FlushReason = iota
	FlushReasonSize                     // mutable memdb exceeds size threshold
	FlushReasonTTL                      // mutable memdb exceeds time-to-live threshold
	FlushReasonGlobalMem                // system-wide memory pressure exceeds threshold
	FlushReasonSegmentRange             // segment time is outside the allowed write window
)

// CheckerConfig holds configuration for the flush checker.
type CheckerConfig struct {
	MaxMemDBSizeBytes        int64         // global default: max mutable memdb size in bytes before flush
	MutableMemDBTTLNano      int64         // global default: max mutable memdb age in nanoseconds before flush
	MaxMemUsageBeforeFlush   float64       // system memory usage ratio [0,1] that triggers global flush sweep
	TargetMemUsageAfterFlush float64       // target memory usage ratio [0,1] after global flush sweep completes
	FlushConcurrency         int           // maximum number of concurrent flush jobs
	CheckInterval            time.Duration // how often the checker evaluates segments; default 1s
}

// Checker is the unified flush scheduler for all three engine types.
// It runs a background goroutine that periodically evaluates registered
// FlushableSegments and triggers flush jobs based on size, TTL, global
// memory pressure, or segment time-range violations.
type Checker interface {
	// Start starts the background flush scheduling loop.
	Start()
	// Stop stops the background flush scheduling loop and waits for completion.
	Stop()
}

// checker is the concrete implementation of Checker.
type checker struct {
	cfg      CheckerConfig
	tracker  MemDBTracker
	memUsage MemoryUsageProvider

	semaphore chan struct{} // concurrency limiter for flush goroutines

	stopCh chan struct{}
	doneCh chan struct{}

	logger logger.Logger
}

// NewChecker creates a new Checker with the given config and memory usage provider.
func NewChecker(cfg CheckerConfig, memUsage MemoryUsageProvider) Checker {
	concurrency := cfg.FlushConcurrency
	if concurrency <= 0 {
		concurrency = 1
	}
	return &checker{
		cfg:       cfg,
		tracker:   GetMemDBTracker(),
		memUsage:  memUsage,
		semaphore: make(chan struct{}, concurrency),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		logger:    logger.GetLogger("Flush", "Checker"),
	}
}

// Start implements Checker.
func (c *checker) Start() {
	go c.run()
}

// Stop implements Checker.
func (c *checker) Stop() {
	close(c.stopCh)
	<-c.doneCh
}

// run is the background scheduling goroutine.
func (c *checker) run() {
	defer close(c.doneCh)

	interval := c.cfg.CheckInterval
	if interval <= 0 {
		interval = time.Minute // safe fallback
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.check()
		}
	}
}

// segmentWithMemInfo bundles a segment with its cached memdb info for one check cycle.
type segmentWithMemInfo struct {
	seg  FlushableSegment
	info *MemDBInfo
}

// check performs one scheduling cycle: updates metrics, handles global memory pressure,
// then evaluates per-segment flush conditions.
func (c *checker) check() {
	segments := c.tracker.GetAll()

	// collect active segments (those with non-nil MemDBInfo)
	active := make([]segmentWithMemInfo, 0, len(segments))
	var totalBytes int64
	for _, seg := range segments {
		info := seg.MutableMemDBInfo()
		if info == nil {
			continue
		}
		active = append(active, segmentWithMemInfo{seg: seg, info: info})
		totalBytes += info.MemSize
	}

	// update gauges
	metrics.FlushCheckerStatistics.ActiveMemDBCount.Update(float64(len(active)))
	metrics.FlushCheckerStatistics.TotalMemDBBytes.Update(float64(totalBytes))

	// step 1: check global memory pressure
	if len(active) > 0 {
		usage, totalBytes, err := c.memUsage.MemoryUsage()
		if err != nil {
			c.logger.Warn("failed to get memory usage for flush check", logger.Error(err))
		} else if usage > c.cfg.MaxMemUsageBeforeFlush {
			c.handleGlobalMemPressure(active, usage, totalBytes)
			return
		}
	}

	// step 2: evaluate per-segment flush conditions
	nowMs := time.Now().UnixMilli()
	for _, item := range active {
		reason := c.evalFlushReason(item.seg, item.info, nowMs)
		if reason == FlushReasonNone {
			continue
		}
		c.triggerFlush(item.seg, reason)
	}
}

// handleGlobalMemPressure flushes segments in descending MemSize order.
// Because flush is asynchronous, we cannot re-query live memory after each trigger.
// Instead we use a static estimate: accumulate the MemSize of segments already
// queued for flush and stop when (usedBytes - willFreeBytes) / totalBytes <= target.
func (c *checker) handleGlobalMemPressure(active []segmentWithMemInfo, usageRatio float64, totalBytes uint64) {
	if totalBytes == 0 {
		return
	}
	// sort by MemSize descending so we free the most memory first
	sort.Slice(active, func(i, j int) bool {
		return active[i].info.MemSize > active[j].info.MemSize
	})

	usedBytes := int64(float64(totalBytes) * usageRatio)
	targetUsedBytes := int64(float64(totalBytes) * c.cfg.TargetMemUsageAfterFlush)
	var willFreeBytes int64

	for _, item := range active {
		// stop once the estimated post-flush usage would be at or below target
		if usedBytes-willFreeBytes <= targetUsedBytes {
			break
		}
		c.triggerFlush(item.seg, FlushReasonGlobalMem)
		willFreeBytes += item.info.MemSize
	}
}

// evalFlushReason evaluates whether a segment should be flushed and returns the reason.
func (c *checker) evalFlushReason(seg FlushableSegment, info *MemDBInfo, nowMs int64) FlushReason {
	meta := seg.SegmentMeta()

	// size threshold: use per-db config if set, otherwise fall back to global default
	sizeThreshold := meta.SizeThresholdBytes
	if sizeThreshold <= 0 {
		sizeThreshold = c.cfg.MaxMemDBSizeBytes
	}
	if sizeThreshold > 0 && info.MemSize >= sizeThreshold {
		return FlushReasonSize
	}

	// TTL threshold: use per-db config if set (convert ms to ns), otherwise global default (ns)
	ttlNano := meta.TimeThresholdNano
	if ttlNano <= 0 {
		ttlNano = c.cfg.MutableMemDBTTLNano
	}
	if ttlNano > 0 && info.Uptime >= time.Duration(ttlNano) {
		return FlushReasonTTL
	}

	// segment range check
	delay := meta.SegmentOutRangeDelay // ms
	segTimeMs := info.SegmentTime
	if meta.Ahead > 0 && segTimeMs > nowMs+meta.Ahead+delay {
		return FlushReasonSegmentRange
	}

	return FlushReasonNone
}

// triggerFlush attempts to dispatch a flush job for the given segment.
// It is non-blocking: if the concurrency semaphore is full, the flush is skipped
// and counted; the segment will be re-evaluated on the next tick.
func (c *checker) triggerFlush(seg FlushableSegment, reason FlushReason) {
	select {
	case c.semaphore <- struct{}{}:
		// acquired slot
	default:
		// concurrency limit reached, skip this cycle
		metrics.FlushCheckerStatistics.FlushSkipped.Incr()
		return
	}

	metrics.FlushCheckerStatistics.FlushInFlight.Incr()
	c.recordFlushReason(reason)

	go func() {
		defer func() {
			<-c.semaphore
			metrics.FlushCheckerStatistics.FlushInFlight.Decr()
		}()

		if err := seg.Flush(); err != nil {
			metrics.FlushCheckerStatistics.FlushErrors.Incr()
			c.logger.Error("flush segment error", logger.Error(err))
		}
	}()
}

// recordFlushReason increments the appropriate reason counter.
func (c *checker) recordFlushReason(reason FlushReason) {
	switch reason {
	case FlushReasonSize:
		metrics.FlushCheckerStatistics.FlushBySize.Incr()
	case FlushReasonTTL:
		metrics.FlushCheckerStatistics.FlushByTTL.Incr()
	case FlushReasonGlobalMem:
		metrics.FlushCheckerStatistics.FlushByGlobal.Incr()
	case FlushReasonSegmentRange:
		metrics.FlushCheckerStatistics.FlushByRange.Incr()
	}
}

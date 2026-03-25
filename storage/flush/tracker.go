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
	"fmt"
	"sync"
	"time"

	"github.com/lindb/common/pkg/timeutil"
)

// MakeSegmentKey builds the unique string key for a segment.
// Format: "{dbName}/{shardID}/{segmentTime}", e.g. "mydb/0/2006-01-02 15:04:05"
func MakeSegmentKey(dbName, shardID string, segmentTimeMs int64) string {
	return fmt.Sprintf("%s/%s/%s", dbName, shardID, timeutil.FormatTimestamp(segmentTimeMs, timeutil.DataTimeFormat2))
}

var (
	tracker MemDBTracker
	once    sync.Once
)

// GetMemDBTracker returns the singleton MemDBTracker instance.
func GetMemDBTracker() MemDBTracker {
	once.Do(func() {
		tracker = NewMemDBTracker()
	})
	return tracker
}

// FlushableSegment defines the interface for segments that support flush scheduling.
// Any segment that contains mutable in-memory data should implement this interface.
type FlushableSegment interface {
	// Flush flushes the mutable in-memory data to persistent storage.
	Flush() error
	// MutableMemDBInfo returns information about the current mutable memory database.
	// Returns nil if there is no flushable content (e.g., already flushed or empty).
	MutableMemDBInfo() *MemDBInfo
	// SegmentMeta returns metadata about the segment used for flush decision-making.
	SegmentMeta() SegmentMeta
	// SegmentKey returns a globally unique key for this segment.
	// Format: "{dbName}/{shardID}/{segmentTime}", e.g. "mydb/0/2006-01-02 15:04:05"
	SegmentKey() string
}

// MemDBInfo holds runtime information about a mutable in-memory database.
type MemDBInfo struct {
	MemSize     int64         // memory size in bytes
	CreatedTime int64         // segment construction time, unix nanoseconds
	SegmentTime int64         // segment's covering start time, unix milliseconds
	NumOfRows   int           // number of rows (metric=series, log=logID count, trace=rowsWritten)
	Uptime      time.Duration // elapsed time since segment creation
}

// SegmentMeta holds metadata about a segment used for flush scheduling decisions.
type SegmentMeta struct {
	DatabaseName         string // name of the owning database
	ShardID              string // shard identifier
	SizeThresholdBytes   int64  // per-db size threshold in bytes; 0 means use global default
	TimeThresholdNano    int64  // per-db time threshold in nanoseconds; 0 means use global default
	Ahead                int64  // allowed write-ahead window, milliseconds
	Behind               int64  // allowed write-behind window, milliseconds
	SegmentOutRangeDelay int64  // grace period after segment goes out of write range, milliseconds
}

// MemDBTracker tracks all flushable segments that contain mutable in-memory data.
// It provides a thread-safe registry for the flush checker to discover segments.
type MemDBTracker interface {
	// Register adds a flushable segment to the tracker.
	Register(seg FlushableSegment)
	// Unregister removes a flushable segment from the tracker.
	Unregister(seg FlushableSegment)
	// GetAll returns a snapshot of all currently registered segments.
	// The returned slice is safe to iterate without holding any locks.
	GetAll() []FlushableSegment
}

// memDBTracker is the default implementation of MemDBTracker.
type memDBTracker struct {
	segments map[string]FlushableSegment // key = SegmentKey()
	mu       sync.RWMutex
}

// NewMemDBTracker creates a new MemDBTracker instance.
func NewMemDBTracker() MemDBTracker {
	return &memDBTracker{
		segments: make(map[string]FlushableSegment),
	}
}

// Register implements MemDBTracker.
func (t *memDBTracker) Register(seg FlushableSegment) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.segments[seg.SegmentKey()] = seg
}

// Unregister implements MemDBTracker.
func (t *memDBTracker) Unregister(seg FlushableSegment) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.segments, seg.SegmentKey())
}

// GetAll implements MemDBTracker.
// Returns a snapshot so callers do not need to hold any lock during iteration.
func (t *memDBTracker) GetAll() []FlushableSegment {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]FlushableSegment, 0, len(t.segments))
	for _, seg := range t.segments {
		result = append(result, seg)
	}
	return result
}

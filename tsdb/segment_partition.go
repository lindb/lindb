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
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/option"

	"io"
	"sort"
	"sync"
)

//go:generate mockgen -source=./segment_partition.go -destination=./segment_partition_mock.go -package=tsdb

type SegmentPartition interface {
	// GetOrCreateSegment returns segment, if not exist create a new segment.
	GetOrCreateSegment(familyTime int64) (Segment, error)
	// GetSegments returns segment list
	GetSegments() []Segment
	// FlushIndex concurrently writes all segment index data to disk.
	FlushIndex() error
	// WaitFlushIndexCompleted waits for the concurrent FlushIndex operation to complete.
	WaitFlushIndexCompleted()
	// TTL expires the data of each dataSegment base on time to live, which runs in parallel.
	TTL() error
	// EvictSegment evicts dataSegment which long term no read operation.
	EvictSegment()
	// Recover segment if existed, recover been invoked when shard init.
	Recover() error
	// Closer releases all segment's resource, such as flush data, spawned goroutines etc.
	io.Closer
}

type segmentPartition struct {
	shard          Shard
	segments       []Segment
	intervals      option.Intervals
	isFlushing     atomic.Bool // restrict flusher concurrency
	flushCondition *sync.Cond  // flush condition
	logger         logger.Logger
}

func newSegmentPartition(shard Shard, intervals option.Intervals) SegmentPartition {
	return &segmentPartition{
		shard:          shard,
		segments:       make([]Segment, 0, 1024),
		intervals:      intervals,
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		logger:         logger.GetLogger("TSDB", "SegmentPartition"),
	}
}

func (sp *segmentPartition) GetOrCreateSegment(familyTime int64) (Segment, error) {
	timestamp := getMonthTimestampFunc(familyTime)
	for _, segment := range sp.segments {
		if segment.GetTimestamp() == timestamp {
			return segment, nil
		}
	}
	segment, err := newSegmentFunc(sp.shard, timestamp, sp.intervals)
	if err != nil {
		return nil, err
	}
	sp.segments = append(sp.segments, segment)
	sort.Slice(sp.segments, func(i, j int) bool {
		return sp.segments[i].GetTimestamp() < sp.segments[j].GetTimestamp()
	})
	return segment, nil
}

func (sp *segmentPartition) Recover() error {
	indexDir := shardIndexPath(sp.shard.Database().Name(), sp.shard.ShardID())
	if !fileExist(indexDir) {
		return nil
	}
	dirs, err := listDirName(indexDir)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		timestamp, err := parseTimestamp(dir, "200601")
		if err != nil {
			return err
		}
		segment, err := newSegmentFunc(sp.shard, timestamp, sp.intervals)
		if err != nil {
			return err
		}
		sp.segments = append(sp.segments, segment)
	}
	sort.Slice(sp.segments, func(i, j int) bool {
		return sp.segments[i].GetTimestamp() < sp.segments[j].GetTimestamp()
	})
	return nil
}

func (sp *segmentPartition) GetSegments() []Segment {
	return sp.segments
}

func (sp *segmentPartition) concurrentDo(handle func(Segment)) {
	segments := sp.segments
	if len(segments) == 0 {
		return
	}
	g := sync.WaitGroup{}
	g.Add(len(segments))
	for _, segment := range segments {
		segment := segment
		go func() {
			defer g.Done()
			handle(segment)
		}()
	}
	g.Wait()
}

func (sp *segmentPartition) FlushIndex() error {
	// another flush process is running
	if !sp.isFlushing.CompareAndSwap(false, true) {
		return nil
	}
	// 1. mark flush job doing
	defer func() {
		sp.flushCondition.L.Lock()
		sp.isFlushing.Store(false)
		sp.flushCondition.L.Unlock()
		// mark flush job complete, notify
		sp.flushCondition.Broadcast()
	}()

	sp.concurrentDo(func(segment Segment) {
		// index flush
		if err := segment.FlushIndex(); err != nil {
			sp.logger.Error("failed to flush indexDB ",
				logger.String("database", sp.shard.Database().Name()),
				logger.Any("shardID", sp.shard.ShardID()),
				logger.String("segmentName", segment.GetName()),
				logger.Error(err))
		} else {
			sp.logger.Info("flush indexDB successfully",
				logger.String("database", sp.shard.Database().Name()),
				logger.Any("shardID", sp.shard.ShardID()),
				logger.String("segmentName", segment.GetName()),
			)
		}
	})

	return nil
}

func (sp *segmentPartition) WaitFlushIndexCompleted() {
	sp.flushCondition.L.Lock()
	if sp.isFlushing.Load() {
		sp.flushCondition.Wait()
	}
	sp.flushCondition.L.Unlock()
}

func (sp *segmentPartition) TTL() error {
	segments := sp.segments
	if len(segments) == 0 {
		return nil
	}

	sp.concurrentDo(func(segment Segment) {
		if err := segment.TTL(); err != nil {
			sp.logger.Warn("do segment ttl failure",
				logger.String("database", sp.shard.Database().Name()),
				logger.Any("shardID", sp.shard.ShardID()),
				logger.String("segmentName", segment.GetName()),
				logger.Error(err))
		}
	})

	return nil
}

func (sp *segmentPartition) EvictSegment() {
	segments := sp.segments
	if len(segments) == 0 {
		return
	}
	sp.concurrentDo(func(segment Segment) {
		segment.EvictSegment()
	})
}

func (sp *segmentPartition) Close() error {
	segments := sp.segments
	if len(segments) == 0 {
		return nil
	}
	sp.concurrentDo(func(segment Segment) {
		if err := segment.Close(); err != nil {
			sp.logger.Warn("do segment close failure",
				logger.String("database", sp.shard.Database().Name()),
				logger.Any("shardID", sp.shard.ShardID()),
				logger.String("segmentName", segment.GetName()),
				logger.Error(err))
		}
	})
	return nil
}

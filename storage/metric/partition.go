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

package metric

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"sync"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
)

type partition struct {
	base.Partition
	shard *Shard

	segments map[int]store.Segment

	kvStore kv.Store

	mutex sync.Mutex
}

func NewPartition(shard *Shard, partitionTime int64, interval timeutil.Interval) (store.Partition, error) {
	calc := interval.Calculator()
	partitionName := calc.GetSegment(partitionTime)
	dir := store.PartitionPath(shard.Database().Name(), shard.ShardID(), interval, partitionName)

	storeOption := kv.DefaultStoreOption()
	intervals := shard.Database().GetOption().Option.Intervals
	if shard.CurrentInterval() == interval && len(intervals) > 1 {
		// if interval == writeable interval and database set auto rollup intervals
		sort.Sort(intervals) // need sort interval
		var rollup []timeutil.Interval
		for _, rollupInterval := range intervals {
			rollup = append(rollup, rollupInterval.Interval)
		}
		storeOption.Rollup = rollup[1:]
		storeOption.Source = interval
	}
	kvStore, err := kv.GetStoreManager().CreateStore(path.Join(dir, "data"), storeOption)
	if err != nil {
		return nil, fmt.Errorf("create kv store for partition error:%s", err)
	}
	p := &partition{
		Partition: base.Partition{
			Timestamp: partitionTime,
			Dir:       dir,
			Interval:  interval,
		},
		shard:    shard,
		kvStore:  kvStore,
		segments: make(map[int]store.Segment),
	}

	segments, err := fileutil.ListDir(p.Dir)
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		if segment == "data" {
			continue
		}
		segmentSlot, err := strconv.Atoi(segment)
		if err != nil {
			// TODO: add metric
			continue
		}
		// create data family
		segmentTime := store.MinuteIntervalCalc.CalcFamilyStartTime(partitionTime, segmentSlot)
		fmt.Printf("create trace segment: %d\n", segmentTime)
		p.GetOrCreateSegment(segmentTime)
	}

	return p, nil
}

// Close implements store.Partition.
func (p *partition) Close() error {
	for _, segment := range p.segments {
		segment.Close()
	}
	return nil
}

// GetOrCreateSegment implements store.Partition.
func (p *partition) GetOrCreateSegment(timestamp int64) (store.Segment, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	segmentKey := p.Interval.Calculator().CalcFamily(timestamp, p.Interval.Calculator().CalcSegmentTime(timestamp))

	if segment, ok := p.segments[segmentKey]; ok {
		return segment, nil
	}

	segment, err := NewSegment(timestamp, p)
	if err != nil {
		return nil, err
	}
	p.segments[segmentKey] = segment
	return segment, nil
}

// GetSegments implements store.Partition.
func (p *partition) GetSegments(timeRange timeutil.TimeRange) (segments []store.Segment) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, segment := range p.segments {
		if (&timeRange).Overlap(segment.SegmentTimeRange()) {
			segments = append(segments, segment)
		}
	}
	return
}

// Shard implements store.Partition.
func (p *partition) Shard() store.Shard {
	return p.shard
}

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
	"fmt"
	"sync"

	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./interval_segment.go -destination=./interval_segment_mock.go -package=tsdb

// IntervalSegment represents a interval segment, there are some segments in a shard.
type IntervalSegment interface {
	// GetOrCreateSegment creates new segment if not exist, if exist return it
	GetOrCreateSegment(segmentName string) (Segment, error)
	// getDataFamilies returns data family list by time range, return nil if not match
	getDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// Close closes interval segment, release resource
	Close()
}

// intervalSegment implements IntervalSegment interface
type intervalSegment struct {
	shard    Shard
	interval timeutil.Interval
	segments sync.Map

	mutex sync.Mutex
}

// newIntervalSegment create interval segment based on interval/type/path etc.
func newIntervalSegment(shard Shard, interval timeutil.Interval) (segment IntervalSegment, err error) {
	path := ShardSegmentPath(shard.Database().Name(), shard.ShardID(), interval)
	err = mkDirIfNotExist(path)
	if err != nil {
		return segment, err
	}
	intervalSegment := &intervalSegment{
		shard:    shard,
		interval: interval,
	}

	defer func() {
		// if create or init segment fail, need close segment
		if err != nil {
			intervalSegment.Close()
		}
	}()

	// load segments if exist
	// TODO too many kv store load???
	segmentNames, err := listDir(path)
	if err != nil {
		return segment, err
	}
	for _, segmentName := range segmentNames {
		seg, err0 := newSegmentFunc(shard, segmentName, intervalSegment.interval)
		if err0 != nil {
			return nil, fmt.Errorf("create segmenet error: %s", err)
		}
		intervalSegment.segments.Store(segmentName, seg)
	}

	// set segment
	segment = intervalSegment
	return segment, nil
}

// GetOrCreateSegment creates new segment if not exist, if exist return it
func (s *intervalSegment) GetOrCreateSegment(segmentName string) (Segment, error) {
	segment, ok := s.getSegment(segmentName)
	if !ok {
		// double check, make sure only create segment once
		s.mutex.Lock()
		defer s.mutex.Unlock()
		segment, ok = s.getSegment(segmentName)
		if !ok {
			seg, err := newSegmentFunc(s.shard, segmentName, s.interval)
			if err != nil {
				return nil, fmt.Errorf("create segmenet error: %s", err)
			}
			s.segments.Store(segmentName, seg)
			return seg, nil
		}
	}
	return segment, nil
}

// getDataFamilies returns data family list by time range, return nil if not match
func (s *intervalSegment) getDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	intervalCalc := s.interval.Calculator()
	segmentQueryTimeRange := &timeutil.TimeRange{
		Start: intervalCalc.CalcSegmentTime(timeRange.Start), // need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		End:   timeRange.End,
	}
	s.segments.Range(func(k, v interface{}) bool {
		segment, ok := v.(Segment)
		if ok {
			baseTime := segment.BaseTime()
			if segmentQueryTimeRange.Contains(baseTime) {
				familyQueryTimeRange := segmentQueryTimeRange.Intersect(timeRange)
				// read lock in segment
				families := segment.getDataFamilies(familyQueryTimeRange)
				if len(families) > 0 {
					result = append(result, families...)
				}
			}
		}
		return true
	})
	return result
}

// Close closes interval segment, release resource
func (s *intervalSegment) Close() {
	s.segments.Range(func(k, v interface{}) bool {
		seg, ok := v.(Segment)
		if ok {
			seg.Close()
		}
		return true
	})
}

// getSegment returns segment by name
func (s *intervalSegment) getSegment(segmentName string) (Segment, bool) {
	segment, _ := s.segments.Load(segmentName)
	seg, ok := segment.(Segment)
	return seg, ok
}

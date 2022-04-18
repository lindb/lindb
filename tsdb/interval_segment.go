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
	"path"
	"sync"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./interval_segment.go -destination=./interval_segment_mock.go -package=tsdb

// IntervalSegment represents a interval segment, there are some segments in a shard.
type IntervalSegment interface {
	// GetOrCreateSegment creates new segment if not exist, if exist return it
	GetOrCreateSegment(segmentName string) (Segment, error)
	// Close closes interval segment, release resource
	Close()
	// getDataFamilies returns data family list by time range, return nil if not match
	getDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// TTL expires segment base on time to live.
	TTL() error
}

// intervalSegment implements IntervalSegment interface
type intervalSegment struct {
	dir      string
	shard    Shard
	interval option.Interval
	segments map[string]Segment

	mutex sync.Mutex

	logger *logger.Logger
}

// newIntervalSegment create interval segment based on interval/type/path etc.
func newIntervalSegment(shard Shard, interval option.Interval) (segment IntervalSegment, err error) {
	dir := ShardSegmentPath(shard.Database().Name(), shard.ShardID(), interval.Interval)
	err = mkDirIfNotExist(dir)
	if err != nil {
		return nil, err
	}
	intervalSegment := &intervalSegment{
		dir:      dir,
		shard:    shard,
		interval: interval,
		segments: make(map[string]Segment),
		logger:   logger.GetLogger("TSDB", "IntervalSegment"),
	}

	defer func() {
		// if create or init segment fail, need close segment
		if err != nil {
			intervalSegment.Close()
		}
	}()

	// load segments if exist
	// TODO too many kv store load???
	segmentNames, err := listDir(dir)
	if err != nil {
		return nil, err
	}
	for _, segmentName := range segmentNames {
		seg, err0 := newSegmentFunc(shard, segmentName, intervalSegment.interval.Interval)
		if err0 != nil {
			return nil, fmt.Errorf("create segmenet error: %s", err0)
		}
		intervalSegment.segments[segmentName] = seg
	}

	// set segment
	segment = intervalSegment
	return segment, nil
}

// GetOrCreateSegment creates new segment if not exist, if exist return it
func (s *intervalSegment) GetOrCreateSegment(segmentName string) (Segment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	segment, ok := s.segments[segmentName]
	if !ok {
		seg, err := newSegmentFunc(s.shard, segmentName, s.interval.Interval)
		if err != nil {
			return nil, fmt.Errorf("create segmenet error: %s", err)
		}
		s.segments[segmentName] = seg
		return seg, nil
	}
	return segment, nil
}

// getDataFamilies returns data family list by time range, return nil if not match
func (s *intervalSegment) getDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	intervalCalc := s.interval.Interval.Calculator()
	segmentQueryTimeRange := &timeutil.TimeRange{
		Start: intervalCalc.CalcSegmentTime(timeRange.Start), // need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		End:   timeRange.End,
	}
	var segments []Segment
	s.mutex.Lock()
	for k := range s.segments {
		segments = append(segments, s.segments[k])
	}
	s.mutex.Unlock()

	for i := range segments {
		segment := segments[i]
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
	return result
}

// Close closes interval segment, release resource
func (s *intervalSegment) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k := range s.segments {
		segment := s.segments[k]
		segment.Close()
	}
}

// TTL expires segment base on time to live.
func (s *intervalSegment) TTL() error {
	segmentNames, err := listDir(s.dir)
	if err != nil {
		return err
	}
	now := timeutil.Now()
	expireInterval := s.interval.Retention.Int64()
	for _, segmentName := range segmentNames {
		calc := s.interval.Interval.Calculator()
		baseTime, err := calc.ParseSegmentTime(segmentName)
		if err != nil {
			s.logger.Warn("parse segment time from path failure",
				logger.String("path", s.dir), logger.String("segment", segmentName),
				logger.Error(err))
			continue
		}
		// add 2 hours buffer, for some cases stop write.
		if now-baseTime > expireInterval+2*timeutil.OneHour {
			s.dropSegment(segmentName)
		}
	}
	return nil
}

// dropSegment drops segment's data.
func (s *intervalSegment) dropSegment(segmentName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	segment, ok := s.segments[segmentName]
	if ok {
		segment.Close()
	}
	delete(s.segments, segmentName)

	if err := removeDir(path.Join(s.dir, segmentName)); err != nil {
		s.logger.Warn("remove segment dir failure",
			logger.String("path", s.dir), logger.String("segment", segmentName),
			logger.Error(err))
		return
	}
	s.logger.Info("do segment ttl successfully",
		logger.String("path", s.dir), logger.String("segment", segmentName))
}

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
	// GetDataFamilies returns data family list by time range, return nil if not match
	GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// Close closes interval segment, release resource
	Close()
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
	return &intervalSegment{
		dir:      dir,
		shard:    shard,
		interval: interval,
		segments: make(map[string]Segment),
		logger:   logger.GetLogger("TSDB", "IntervalSegment"),
	}, nil
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

// GetDataFamilies returns data family list by time range, return nil if not match
func (s *intervalSegment) GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	now := timeutil.Now()
	intervalCalc := s.interval.Interval.Calculator()
	segmentQueryTimeRange := &timeutil.TimeRange{
		// need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		Start: intervalCalc.CalcSegmentTime(timeRange.Start),
		End:   timeRange.End,
	}

	expireInterval := s.interval.Retention.Int64()
	if err := s.walkSegment(func(segmentName string, segmentTime int64) {
		if now-segmentTime >= expireInterval {
			// segment is expired, need to ignore
			return
		}
		if !segmentQueryTimeRange.Contains(segmentTime) {
			// segment time not in time range of query
			return
		}

		segment, err := s.getOrLoadSegment(segmentName)
		if err != nil {
			// TODO add metric
			// ignore err
			return
		}

		familyQueryTimeRange := segmentQueryTimeRange.Intersect(timeRange)
		// read lock in segment
		families := segment.GetDataFamilies(familyQueryTimeRange)
		if len(families) > 0 {
			result = append(result, families...)
		}
	}); err != nil {
		s.logger.Warn("list segment failure when get data families",
			logger.String("path", s.dir), logger.Error(err))
		return nil
	}
	return result
}

// getOrLoadSegment returns segment for current interval.
// 1. return segment if it's exist in memory cache;
// 2. return segment if it's exist in storage.
func (s *intervalSegment) getOrLoadSegment(segmentName string) (Segment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	segment, ok := s.segments[segmentName]
	if !ok {
		var err error
		segment, err = newSegmentFunc(s.shard, segmentName, s.interval.Interval)
		if err != nil {
			return nil, err
		}
		s.segments[segmentName] = segment
	}
	return segment, nil
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
	now := timeutil.Now()
	expireInterval := s.interval.Retention.Int64()

	return s.walkSegment(func(segmentName string, segmentTime int64) {
		// add 2 hours buffer, for some cases stop write.
		if now-segmentTime > expireInterval+2*timeutil.OneHour {
			s.dropSegment(segmentName)
		}
	})
}

// walkSegment lists all segment under current interval segment dir.
func (s *intervalSegment) walkSegment(fn func(segmentName string, segmentTime int64)) error {
	segmentNames, err := listDir(s.dir)
	if err != nil {
		return err
	}
	for _, segmentName := range segmentNames {
		calc := s.interval.Interval.Calculator()
		baseTime, err := calc.ParseSegmentTime(segmentName)
		if err != nil {
			s.logger.Warn("parse segment time from path failure",
				logger.String("path", s.dir), logger.String("segment", segmentName),
				logger.Error(err))
			continue
		}
		fn(segmentName, baseTime)
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

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

	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./interval_data_segment.go -destination=./interval_data_segment_mock.go -package=tsdb

// IntervalDataSegment represents a interval dataSegment, there are some segments in a shard.
type IntervalDataSegment interface {
	// GetOrCreateSegment creates new dataSegment if not exist, if exist return it
	GetOrCreateSegment(segmentName string) (DataSegment, error)
	// GetDataFamilies returns data family list by time range, return nil if not match
	GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily
	// Close closes interval dataSegment, release resource
	Close()
	// TTL expires dataSegment base on time to live.
	TTL() error
	// EvictSegment evicts dataSegment which long term no read operation.
	EvictSegment()
}

// intervalDataSegment implements IntervalDataSegment interface
type intervalDataSegment struct {
	dir      string
	shard    Shard
	interval option.Interval
	segments map[string]DataSegment

	mutex sync.Mutex

	logger logger.Logger
}

// newIntervalDataSegment create interval dataSegment based on interval/type/path etc.
func newIntervalDataSegment(shard Shard, interval option.Interval) (segment IntervalDataSegment, err error) {
	dir := ShardIntervalSegmentPath(shard.Database().Name(), shard.ShardID(), interval.Interval)
	err = mkDirIfNotExist(dir)
	if err != nil {
		return nil, err
	}
	return &intervalDataSegment{
		dir:      dir,
		shard:    shard,
		interval: interval,
		segments: make(map[string]DataSegment),
		logger:   logger.GetLogger("TSDB", "IntervalDataSegment"),
	}, nil
}

// GetOrCreateSegment creates new dataSegment if not exist, if exist return it
func (s *intervalDataSegment) GetOrCreateSegment(segmentName string) (DataSegment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if segment, ok := s.segments[segmentName]; ok {
		return segment, nil
	}
	seg, err := newDataSegmentFunc(s.shard, segmentName, s.interval.Interval)
	if err != nil {
		return nil, fmt.Errorf("create segmenet error: %s", err)
	}
	s.segments[segmentName] = seg
	return seg, nil
}

// GetDataFamilies returns data family list by time range, return nil if not match
func (s *intervalDataSegment) GetDataFamilies(timeRange timeutil.TimeRange) []DataFamily {
	var result []DataFamily
	now := commontimeutil.Now()
	intervalCalc := s.interval.Interval.Calculator()
	segmentQueryTimeRange := &timeutil.TimeRange{
		// need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		Start: intervalCalc.CalcSegmentTime(timeRange.Start),
		End:   timeRange.End,
	}

	expireInterval := s.interval.Retention.Int64()
	if err := s.walkSegment(func(segmentName string, segmentTime int64) {
		if now-segmentTime >= expireInterval {
			// dataSegment is expired, need to ignore
			return
		}
		if !segmentQueryTimeRange.Contains(segmentTime) {
			// dataSegment time not in time range of query
			return
		}

		segment, err := s.getOrLoadSegment(segmentName)
		if err != nil {
			// TODO: add metric
			// ignore err
			s.logger.Info("get or load dataSegment failure",
				logger.String("path", s.dir), logger.String("dataSegment", segmentName), logger.Error(err))
			return
		}

		familyQueryTimeRange := segmentQueryTimeRange.Intersect(timeRange)
		// read lock in dataSegment
		families := segment.GetDataFamilies(familyQueryTimeRange)
		if len(families) > 0 {
			result = append(result, families...)
		}
	}); err != nil {
		s.logger.Warn("list dataSegment failure when get data families",
			logger.String("path", s.dir), logger.Error(err))
		return nil
	}
	return result
}

// getOrLoadSegment returns dataSegment for current interval.
// 1. return dataSegment if it's exist in memory cache;
// 2. return dataSegment if it's exist in storage.
func (s *intervalDataSegment) getOrLoadSegment(segmentName string) (DataSegment, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if segment, ok := s.segments[segmentName]; ok {
		return segment, nil
	}
	var err error
	segment, err := newDataSegmentFunc(s.shard, segmentName, s.interval.Interval)
	if err != nil {
		return nil, err
	}
	s.segments[segmentName] = segment
	return segment, nil
}

// Close closes interval dataSegment, release resource
func (s *intervalDataSegment) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k := range s.segments {
		segment := s.segments[k]
		segment.Close()
	}
}

// TTL expires dataSegment base on time to live.
func (s *intervalDataSegment) TTL() error {
	now := commontimeutil.Now()
	expireInterval := s.interval.Retention.Int64()

	return s.walkSegment(func(segmentName string, segmentTime int64) {
		// add 2 hours buffer, for some cases stop write.
		if now-segmentTime > expireInterval+2*commontimeutil.OneHour {
			s.dropSegment(segmentName)
		}
	})
}

// EvictSegment evicts dataSegment which long term no read operation.
func (s *intervalDataSegment) EvictSegment() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for segmentName := range s.segments {
		segment := s.segments[segmentName]
		// add 2 hours buffer, for some cases stop write.
		if segment.NeedEvict() {
			segment.Close()
			delete(s.segments, segmentName)

			s.logger.Info("evict dataSegment complete",
				logger.String("path", s.dir), logger.String("dataSegment", segmentName))
		}
	}
}

// walkSegment lists all dataSegment under current interval dataSegment dir.
func (s *intervalDataSegment) walkSegment(fn func(segmentName string, segmentTime int64)) error {
	fmt.Println(s.dir)
	segmentNames, err := listDir(s.dir)
	if err != nil {
		return err
	}
	for _, segmentName := range segmentNames {
		calc := s.interval.Interval.Calculator()
		baseTime, err := calc.ParseSegmentTime(segmentName)
		if err != nil {
			s.logger.Warn("parse dataSegment time from path failure",
				logger.String("path", s.dir), logger.String("dataSegment", segmentName),
				logger.Error(err))
			continue
		}
		fn(segmentName, baseTime)
	}
	return nil
}

// dropSegment drops dataSegment's data.
func (s *intervalDataSegment) dropSegment(segmentName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if segment, ok := s.segments[segmentName]; ok {
		segment.Close()
	}
	delete(s.segments, segmentName)

	if err := removeDir(path.Join(s.dir, segmentName)); err != nil {
		s.logger.Warn("remove dataSegment dir failure",
			logger.String("path", s.dir), logger.String("dataSegment", segmentName),
			logger.Error(err))
		return
	}
	s.logger.Info("do dataSegment ttl successfully",
		logger.String("path", s.dir), logger.String("dataSegment", segmentName))
}

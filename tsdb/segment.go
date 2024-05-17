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

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"

	"io"
	"path"
	"strconv"
	"time"
)

//go:generate mockgen -source=./segment.go -destination=./segment_mock.go -package=tsdb

type Segment interface {
	// GetName returns segment name like "202401".
	GetName() string
	// IndexDB returns metric index database.
	IndexDB() index.MetricIndexDatabase
	// GetOrCreateDataFamily returns data family, if not exist create a new data family.
	GetOrCreateDataFamily(familyTime int64) (DataFamily, error)
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// GetTimestamp returns the start of the month as a timestamp.
	GetTimestamp() int64
	// FlushIndex flush index data to disk
	FlushIndex() error
	// TTL expires the data of each dataSegment base on time to live.
	TTL() error
	// EvictSegment evicts dataSegment which long term no read operation.
	EvictSegment()
	// Closer releases segment's resource, such as flush data, spawned goroutines etc.
	io.Closer
}

type segment struct {
	name                string
	shard               Shard
	indexDB             index.MetricIndexDatabase
	writableDataSegment IntervalDataSegment
	rollupTargets       map[timeutil.Interval]IntervalDataSegment
	interval            timeutil.Interval
	timestamp           int64
	statistics          *metrics.SegmentStatistics
	logger              logger.Logger
}

func newSegment(shard Shard, timestamp int64, intervals option.Intervals) (Segment, error) {
	segmentName := formatTimestamp(timestamp, "200601")
	dir := path.Join(shardIndexPath(shard.Database().Name(), shard.ShardID()), segmentName)
	indexDB, err := newIndexDBFunc(dir, shard.Database().MetaDB())
	if err != nil {
		return nil, err
	}

	segment := &segment{
		name:          segmentName,
		shard:         shard,
		indexDB:       indexDB,
		rollupTargets: make(map[timeutil.Interval]IntervalDataSegment),
		timestamp:     timestamp,
		statistics:    metrics.NewSegmentStatistics(shard.Database().Name(), strconv.Itoa(int(shard.ShardID())), segmentName),
		logger:        logger.GetLogger("TSDB", "Segment"),
	}

	for idx, targetInterval := range intervals {
		// new dataSegment for rollup
		intervalDataSegment, err := newIntervalDataSegmentFunc(shard, targetInterval)
		if err != nil {
			return nil, err
		}
		if idx == 0 {
			segment.interval = targetInterval.Interval
			// the smallest interval for writing
			segment.writableDataSegment = intervalDataSegment
		}
		// set rollup target dataSegment
		segment.rollupTargets[targetInterval.Interval] = intervalDataSegment
	}

	return segment, nil
}

func (s *segment) GetName() string {
	return s.name
}

func (s *segment) GetOrCreateDataFamily(familyTime int64) (DataFamily, error) {
	segmentName := s.interval.Calculator().GetSegment(familyTime)
	// source dataSegment
	dataSegment, err := s.writableDataSegment.GetOrCreateSegment(segmentName)
	if err != nil {
		return nil, err
	}
	// build rollup target dataSegment if set auto rollup interval
	for interval, rollupSegment := range s.rollupTargets {
		_, err = rollupSegment.GetOrCreateSegment(interval.Calculator().GetSegment(familyTime))
		if err != nil {
			return nil, err
		}
	}
	family, err := dataSegment.GetOrCreateDataFamily(familyTime)
	if err != nil {
		return nil, err
	}
	return family, nil
}

func (s *segment) GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily {
	// first check query interval is writable interval.
	if s.interval.Type() == intervalType || len(s.rollupTargets) == 1 {
		// if no rollup, need to use current writable interval.
		return s.writableDataSegment.GetDataFamilies(timeRange)
	}
	// then find family from rollup targets
	for interval, rollupSegment := range s.rollupTargets {
		if interval.Type() == intervalType {
			return rollupSegment.GetDataFamilies(timeRange)
		}
	}
	return nil
}

func (s *segment) Close() error {
	for _, rollupSegment := range s.rollupTargets {
		rollupSegment.Close()
	}
	return nil
}

func (s *segment) IndexDB() index.MetricIndexDatabase {
	return s.indexDB
}

func (s *segment) GetTimestamp() int64 {
	return s.timestamp
}

func (s *segment) FlushIndex() error {
	startTime := time.Now()
	defer func() {
		s.statistics.IndexDBFlushDuration.UpdateSince(startTime)
	}()

	ch := make(chan error, 1)
	s.indexDB.Notify(&index.FlushNotifier{
		Callback: func(err error) {
			ch <- err
		},
	})
	err := <-ch
	if err != nil {
		s.statistics.IndexDBFlushFailures.Incr()
	}

	return err
}

func (s *segment) TTL() error {
	for _, rollupSegment := range s.rollupTargets {
		if err := rollupSegment.TTL(); err != nil {
			s.logger.Warn("do segment ttl failure",
				logger.String("database", s.shard.Database().Name()),
				logger.Any("shardID", s.shard.ShardID()),
				logger.String("segmentName", s.name),
				logger.Error(err),
			)
		}
	}
	return nil
}

func (s *segment) EvictSegment() {
	for _, rollupSegment := range s.rollupTargets {
		rollupSegment.EvictSegment()
	}
}

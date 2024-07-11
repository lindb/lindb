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
	"io"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/memdb"
)

//go:generate mockgen -source=./shard.go -destination=./shard_mock.go -package=tsdb

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
	// Database returns the database.
	Database() Database
	// ShardID returns the shard id.
	ShardID() models.ShardID
	// CurrentInterval returns current interval for metric write.
	CurrentInterval() timeutil.Interval
	// Indicator returns the unique shard info.
	Indicator() string
	// GetOrCrateDataFamily returns data family, if not exist create a new data family.
	GetOrCrateDataFamily(familyTime int64) (DataFamily, error)
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// IndexDB returns the metric index database, include inverted/forward index.
	IndexDB() index.MetricIndexDatabase
	// MemIndexDB returns memory index database.
	MemIndexDB() memdb.IndexDatabase
	// BufferManager returns write temp memory manager.
	BufferManager() memdb.BufferManager
	// FlushIndex flushes index data to disk.
	FlushIndex() error
	// WaitFlushIndexCompleted waits flush index job completed.
	WaitFlushIndexCompleted()
	// initIndexDatabase initializes index database
	initIndexDatabase() error
	// TTL expires the data of each segment base on time to live.
	TTL()
	// EvictSegment evicts segment which long term no read operation.
	EvictSegment()
	// notifyLimitsChange notifies the limits changed.
	notifyLimitsChange()
	// Closer releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer
}

// shard implements Shard interface
type shard struct {
	db     Database
	option *option.DatabaseOption

	bufferMgr memdb.BufferManager
	// segments keeps all rollup target interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	rollupTargets  map[timeutil.Interval]IntervalSegment
	segment        IntervalSegment // smallest interval for writing data
	flushCondition *sync.Cond      // flush condition

	limits *models.Limits // NOTE: limits only update in write goroutine
	logger logger.Logger

	statistics *metrics.ShardStatistics

	indexDB    index.MetricIndexDatabase
	memIndexDB memdb.IndexDatabase

	indicator string // => db/shard
	// write accept time range
	interval timeutil.Interval
	id       models.ShardID

	isFlushing    atomic.Bool // restrict flusher concurrency
	limitsChanged atomic.Bool
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if create fail.
func newShard(
	db Database,
	shardID models.ShardID,
) (s Shard, err error) {
	shardPath := shardPath(db.Name(), shardID)
	err = mkDirIfNotExist(shardPath)
	if err != nil {
		return nil, err
	}
	dbOption := db.GetOption()
	createdShard := &shard{
		db:             db,
		indicator:      shardIndicator(db.Name(), shardID),
		id:             shardID,
		option:         dbOption,
		bufferMgr:      memdb.NewBufferManager(shardTempBufferPath(db.Name(), shardID)),
		rollupTargets:  make(map[timeutil.Interval]IntervalSegment),
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewShardStatistics(db.Name(), strconv.Itoa(int(shardID))),
		logger:         logger.GetLogger("TSDB", "Shard"),
	}
	// try cleanup history dirty write buffer
	createdShard.bufferMgr.Cleanup()

	// sort intervals
	sort.Sort(dbOption.Intervals)

	createdShard.interval = dbOption.Intervals[0].Interval

	for idx, targetInterval := range dbOption.Intervals {
		// new segment for rollup
		var segment IntervalSegment
		segment, err = newIntervalSegmentFunc(createdShard, targetInterval)
		if err != nil {
			return nil, err
		}
		if idx == 0 {
			// the smallest interval for writing
			createdShard.segment = segment
		}
		// set rollup target segment
		createdShard.rollupTargets[targetInterval.Interval] = segment
	}

	defer func() {
		if err == nil {
			return
		}
		if err0 := createdShard.Close(); err0 != nil {
			engineLogger.Error("close shard error when create shard fail",
				logger.String("database", createdShard.db.Name()),
				logger.Any("shardID", createdShard.id), logger.Error(err0))
		}
	}()

	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	// init datatbase limits
	createdShard.limits = db.GetLimits()

	createdShard.memIndexDB = memdb.NewIndexDatabase(db.MemMetaDB(), createdShard.indexDB)

	return createdShard, nil
}

// Database returns the database.
func (s *shard) Database() Database { return s.db }

// ShardID returns the shard id.
func (s *shard) ShardID() models.ShardID { return s.id }

// Indicator returns the unique shard info.
func (s *shard) Indicator() string { return s.indicator }

// notifyLimitsChange notifies the limits changed.
func (s *shard) notifyLimitsChange() {
	s.limitsChanged.Store(true)
}

// CurrentInterval returns current interval for metric  write.
func (s *shard) CurrentInterval() timeutil.Interval { return s.interval }

// BufferManager returns write temp memory manager.
func (s *shard) BufferManager() memdb.BufferManager {
	return s.bufferMgr
}

// IndexDB returns the metric index database, include inverted/forward index.
func (s *shard) IndexDB() index.MetricIndexDatabase {
	return s.indexDB
}

// MemIndexDB returns memory index database.
func (s *shard) MemIndexDB() memdb.IndexDatabase {
	return s.memIndexDB
}

func (s *shard) GetOrCrateDataFamily(familyTime int64) (DataFamily, error) {
	segmentName := s.interval.Calculator().GetSegment(familyTime)
	// source segment
	segment, err := s.segment.GetOrCreateSegment(segmentName)
	if err != nil {
		return nil, err
	}
	// build rollup target segment if set auto rollup interval
	for interval, rollupSegment := range s.rollupTargets {
		_, err = rollupSegment.GetOrCreateSegment(interval.Calculator().GetSegment(familyTime))
		if err != nil {
			return nil, err
		}
	}
	family, err := segment.GetOrCreateDataFamily(familyTime)
	if err != nil {
		return nil, err
	}
	return family, nil
}

func (s *shard) GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily {
	// first check query interval is writable interval.
	if s.interval.Type() == intervalType || len(s.rollupTargets) == 1 {
		// if no rollup, need to use current writable interval.
		return s.segment.GetDataFamilies(timeRange)
	}
	// then find family from rollup targets
	for interval, rollupSegment := range s.rollupTargets {
		if interval.Type() == intervalType {
			return rollupSegment.GetDataFamilies(timeRange)
		}
	}
	return nil
}

func (s *shard) Close() error {
	// finally, cleanup temp buffer.
	defer s.bufferMgr.Cleanup()
	// wait previous flush job completed
	s.WaitFlushIndexCompleted()

	if s.memIndexDB != nil {
		// need flush index data
		if err := s.flushIndex(); err != nil {
			return err
		}
		s.memIndexDB.Close()
	}
	if s.indexDB != nil {
		// flush index db in database level
		if err := s.indexDB.Close(); err != nil {
			return err
		}
	}
	// close segment/flush family data
	s.segment.Close()
	for _, rollupSegment := range s.rollupTargets {
		rollupSegment.Close()
	}
	return nil
}

// FlushIndex flushes index data to disk
func (s *shard) FlushIndex() (err error) {
	// another flush process is running
	if !s.isFlushing.CompareAndSwap(false, true) {
		return nil
	}
	// 1. mark flush job doing
	startTime := time.Now()
	defer func() {
		s.flushCondition.L.Lock()
		s.isFlushing.Store(false)
		s.flushCondition.L.Unlock()
		// mark flush job complete, notify
		s.flushCondition.Broadcast()
		s.statistics.IndexDBFlushDuration.UpdateSince(startTime)
	}()
	// index flush
	if err = s.flushIndex(); err != nil {
		s.statistics.IndexDBFlushFailures.Incr()
		s.logger.Error("failed to flush indexDB ",
			logger.String("database", s.db.Name()),
			logger.Any("shardID", s.id),
			logger.Error(err))
		return err
	}
	s.logger.Info("flush indexDB successfully",
		logger.String("database", s.db.Name()),
		logger.Any("shardID", s.id),
	)

	return nil
}

func (s *shard) flushIndex() error {
	ch := make(chan error, 1)
	s.memIndexDB.Notify(&memdb.FlushEvent{
		Callback: func(err error) {
			ch <- err
		},
	})
	if err := <-ch; err != nil {
		return err
	}
	return nil
}

// WaitFlushIndexCompleted waits flush index job completed.
func (s *shard) WaitFlushIndexCompleted() {
	s.flushCondition.L.Lock()
	if s.isFlushing.Load() {
		s.flushCondition.Wait()
	}
	s.flushCondition.L.Unlock()
}

// TTL expires the data of each segment base on time to live.
func (s *shard) TTL() {
	for interval, rollupSegment := range s.rollupTargets {
		if err := rollupSegment.TTL(); err != nil {
			s.logger.Warn("do segment ttl failure",
				logger.String("database", s.db.Name()),
				logger.Any("shardID", s.id),
				logger.String("segment", interval.Type().String()),
				logger.Error(err),
			)
		}
	}
}

// EvictSegment evicts segment which long term no read operation.
func (s *shard) EvictSegment() {
	for _, rollupSegment := range s.rollupTargets {
		rollupSegment.EvictSegment()
	}
}

// initIndexDatabase initializes the index database
func (s *shard) initIndexDatabase() error {
	var err error
	s.indexDB, err = newIndexDBFunc(shardIndexPath(s.db.Name(), s.ShardID()), s.db.MetaDB())
	if err != nil {
		return err
	}
	return nil
}

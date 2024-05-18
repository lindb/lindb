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

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/memdb"

	"fmt"
	"io"
	"sort"
	"strconv"
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
	// GetOrCreateDataFamily returns data family, if not exist create a new data family.
	GetOrCreateDataFamily(familyTime int64) (DataFamily, error)
	// GetDataFamilies returns data family list by interval type and time range, return nil if not match
	GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily
	// IndexDB returns the metric index database, include inverted/forward index.
	IndexDB() index.MetricIndexDatabase
	// GetIndexDB returns the metric index database as family time
	GetIndexDB(familyTime int64) index.MetricIndexDatabase
	// BufferManager returns write temp memory manager.
	BufferManager() memdb.BufferManager
	// FlushIndex flushes index data to disk.
	FlushIndex() error
	// WaitFlushIndexCompleted waits flush index job completed.
	WaitFlushIndexCompleted()
	// TTL expires the data of each dataSegment base on time to live.
	TTL()
	// EvictSegment evicts dataSegment which long term no read operation.
	EvictSegment()
	// notifyLimitsChange notifies the limits changed.
	notifyLimitsChange()
	// Closer releases shard's resource, such as flush data, spawned goroutines etc.
	io.Closer
}

// shard implements Shard interface
type shard struct {
	db        Database
	indicator string // => db/shard
	id        models.ShardID
	option    *option.DatabaseOption

	bufferMgr memdb.BufferManager

	limits        *models.Limits // NOTE: limits only update in write goroutine
	limitsChanged atomic.Bool
	logger        logger.Logger

	statistics *metrics.ShardStatistics

	segmentPartition SegmentPartition
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
		db:         db,
		indicator:  shardIndicator(db.Name(), shardID),
		id:         shardID,
		option:     dbOption,
		bufferMgr:  memdb.NewBufferManager(shardTempBufferPath(db.Name(), shardID)),
		statistics: metrics.NewShardStatistics(db.Name(), strconv.Itoa(int(shardID))),
		logger:     logger.GetLogger("TSDB", "Shard"),
	}

	// try cleanup history dirty write buffer
	createdShard.bufferMgr.Cleanup()
	// sort intervals
	sort.Sort(dbOption.Intervals)

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

	createdShard.segmentPartition = newSegmentPartitionFunc(createdShard, dbOption.Intervals)
	if err0 := createdShard.segmentPartition.Recover(); err0 != nil {
		return nil, err0
	}

	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	// init database limits
	createdShard.limits = db.GetLimits()
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
func (s *shard) CurrentInterval() timeutil.Interval {
	return s.option.Intervals[0].Interval
}

// BufferManager returns write temp memory manager.
func (s *shard) BufferManager() memdb.BufferManager {
	return s.bufferMgr
}

// IndexDB returns the metric index database, include inverted/forward index.
func (s *shard) IndexDB() index.MetricIndexDatabase {
	return nil
}

func (s *shard) GetIndexDB(familyTime int64) index.MetricIndexDatabase {
	segment, _ := s.segmentPartition.GetOrCreateSegment(familyTime)
	return segment.IndexDB()
}

func (s *shard) GetOrCreateDataFamily(familyTime int64) (DataFamily, error) {
	segment, err := s.segmentPartition.GetOrCreateSegment(familyTime)
	if err != nil {
		return nil, err
	}
	return segment.GetOrCreateDataFamily(familyTime)
}

func (s *shard) GetDataFamilies(intervalType timeutil.IntervalType, timeRange timeutil.TimeRange) []DataFamily {
	var dataFamilies []DataFamily
	for _, segment := range s.segmentPartition.GetSegments() {
		result := segment.GetDataFamilies(intervalType, timeRange)
		dataFamilies = append(dataFamilies, result...)
	}
	return dataFamilies
}

func (s *shard) Close() error {
	// finally, cleanup temp buffer.
	defer s.bufferMgr.Cleanup()

	// wait previous flush job completed
	s.segmentPartition.WaitFlushIndexCompleted()

	// flush index
	if err := s.segmentPartition.FlushIndex(); err != nil {
		s.statistics.IndexDBFlushFailures.Incr()
		s.logger.Error("failed to flush indexDB ",
			logger.String("database", s.db.Name()),
			logger.Any("shardID", s.id),
			logger.Error(err))
	}

	// close index
	if err := s.segmentPartition.Close(); err != nil {
		s.statistics.IndexDBFlushFailures.Incr()
		s.logger.Error("failed to flush indexDB ",
			logger.String("database", s.db.Name()),
			logger.Any("shardID", s.id),
			logger.Error(err))
	}

	return nil
}

// FlushIndex flushes index data to disk
func (s *shard) FlushIndex() (err error) {
	if err := s.segmentPartition.FlushIndex(); err != nil {
		s.statistics.IndexDBFlushFailures.Incr()
		s.logger.Error("failed to flush indexDB ",
			logger.String("database", s.db.Name()),
			logger.Any("shardID", s.id),
			logger.Error(err))
	}
	return nil
}

// WaitFlushIndexCompleted waits flush index job completed.
func (s *shard) WaitFlushIndexCompleted() {
	s.segmentPartition.WaitFlushIndexCompleted()
}

// TTL expires the data of each segment base on time to live.
func (s *shard) TTL() {
	if err := s.segmentPartition.TTL(); err != nil {
		s.logger.Warn("do segment ttl failure",
			logger.String("database", s.db.Name()),
			logger.Any("shardID", s.id),
			logger.Error(err))
	}
}

// EvictSegment evicts segment which long term no read operation.
func (s *shard) EvictSegment() {
	s.segmentPartition.EvictSegment()
}

// initIndexDatabase initializes the index database
func (s *shard) initIndexDatabase() error {
	return nil
}

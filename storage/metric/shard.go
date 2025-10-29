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
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/metric/memdb"
	"github.com/lindb/lindb/storage/store"
)

// Shard implements Shard interface
type Shard struct {
	base.Shard

	bufferMgr memdb.BufferManager
	// segments keeps all rollup target interval segments,
	// includes one smallest interval segment for writing data, and rollup interval segments
	rollupPartitions map[timeutil.Interval]*store.Partitions

	flushCondition *sync.Cond // flush condition

	statistics *metrics.ShardStatistics

	indexDB    index.MetricIndexDatabase
	memIndexDB memdb.IndexDatabase

	indicator string // => db/shard
	// write accept time range
	interval timeutil.Interval

	isFlushing atomic.Bool // restrict flusher concurrency

	logger logger.Logger
}

// newShard creates shard instance, if shard path exist then load shard data for init.
// return error if create fail.
func newShard(
	db *Database,
	shardID models.ShardID,
) (s store.Shard, err error) {
	shardPath := store.ShardPath(db.Name(), shardID)
	err = fileutil.MkDirIfNotExist(shardPath)
	if err != nil {
		return nil, err
	}
	dbOption := db.GetOption().Option
	createdShard := &Shard{
		Shard: base.Shard{
			ID: shardID,
			DB: db,
		},
		indicator:        shardIndicator(db.Name(), shardID),
		bufferMgr:        memdb.NewBufferManager(shardTempBufferPath(db.Name(), shardID)),
		rollupPartitions: make(map[timeutil.Interval]*store.Partitions),
		isFlushing:       *atomic.NewBool(false),
		flushCondition:   sync.NewCond(&sync.Mutex{}),
		statistics:       metrics.NewShardStatistics(db.Name(), strconv.Itoa(int(shardID))),
		logger:           logger.GetLogger("Metric", "Shard"),
	}
	// try cleanup history dirty write buffer
	createdShard.bufferMgr.Cleanup()

	// sort intervals
	sort.Sort(dbOption.Intervals)

	createdShard.interval = dbOption.Intervals[0].Interval

	defer func() {
		if err == nil {
			return
		}
		if err0 := createdShard.Close(); err0 != nil {
			createdShard.logger.Error("close shard error when create shard fail",
				logger.String("database", createdShard.Database().Name()),
				logger.Any("shardID", createdShard.ShardID()), logger.Error(err0))
		}
	}()

	for idx, targetInterval := range dbOption.Intervals {
		// new partitions for rollup
		partitions := store.NewPartitions(targetInterval.Interval)
		if idx == 0 {
			// the smallest interval for writing
			if err = partitions.Load(store.PartitionsPath(db.Name(), shardID, targetInterval.Interval), func(timestamp int64) (*store.LazyPartition, error) {
				return store.NewLazyPartition(timestamp, createdShard.createPartition), nil
			}); err != nil {
				break
			}

			createdShard.Partitions = partitions
			createdShard.CalcPartitionTimeFn = targetInterval.Interval.Calculator().CalcSegmentTime
		}
		// set rollup partitions
		createdShard.rollupPartitions[targetInterval.Interval] = partitions
	}

	if err = createdShard.initIndexDatabase(); err != nil {
		return nil, fmt.Errorf("create index database for shard[%d] error: %s", shardID, err)
	}
	createdShard.memIndexDB = memdb.NewIndexDatabase(db.MemMetaDB(), createdShard.indexDB)

	createdShard.CreatePartitionFn = createdShard.createPartition
	return createdShard, nil
}

func (s *Shard) createPartition(timestamp int64) (store.Partition, error) {
	return NewPartition(s, timestamp, s.interval)
}

// Indicator returns the unique shard info.
func (s *Shard) Indicator() string { return s.indicator }

// CurrentInterval returns current interval for metric  write.
func (s *Shard) CurrentInterval() timeutil.Interval { return s.interval }

// BufferManager returns write temp memory manager.
func (s *Shard) BufferManager() memdb.BufferManager {
	return s.bufferMgr
}

// IndexDB returns the metric index database, include inverted/forward index.
func (s *Shard) IndexDB() index.MetricIndexDatabase {
	return s.indexDB
}

// MemIndexDB returns memory index database.
func (s *Shard) MemIndexDB() memdb.IndexDatabase {
	return s.memIndexDB
}

func (s *Shard) Close() error {
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
	for _, rollupSegment := range s.rollupPartitions {
		partitions := rollupSegment.GetPartitions()
		for _, partition := range partitions {
			if partition.Loaded() {

				p, err := partition.Get()
				if err != nil {
					s.logger.Warn("load partition fail when close shard", logger.Error(err))
					continue
				}
				p.Close()
			}
		}
	}
	return nil
}

// FlushIndex flushes index data to disk
func (s *Shard) FlushIndex() (err error) {
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
			logger.String("database", s.Database().Name()),
			logger.Any("shardID", s.ShardID()),
			logger.Error(err))
		return err
	}
	s.logger.Info("flush indexDB successfully",
		logger.String("database", s.Database().Name()),
		logger.Any("shardID", s.ShardID()),
	)

	return nil
}

func (s *Shard) flushIndex() error {
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
func (s *Shard) WaitFlushIndexCompleted() {
	s.flushCondition.L.Lock()
	if s.isFlushing.Load() {
		s.flushCondition.Wait()
	}
	s.flushCondition.L.Unlock()
}

// TTL expires the data of each segment base on time to live.
func (s *Shard) TTL() {
	// for interval, rollupSegment := range s.rollupTargets {
	// 	if err := rollupSegment.TTL(); err != nil {
	// 		s.logger.Warn("do segment ttl failure",
	// 			logger.String("database", s.db.Name()),
	// 			logger.Any("shardID", s.id),
	// 			logger.String("segment", interval.Type().String()),
	// 			logger.Error(err),
	// 		)
	// 	}
	// }
}

// EvictSegment evicts segment which long term no read operation.
func (s *Shard) EvictSegment() {
	// for _, rollupSegment := range s.rollupTargets {
	// 	rollupSegment.EvictSegment()
	// }
}

// initIndexDatabase initializes the index database
func (s *Shard) initIndexDatabase() error {
	var err error
	db := s.Database().(*Database)
	s.indexDB, err = index.NewMetricIndexDatabase(shardIndexPath(db.Name(), s.ShardID()), db.MetaDB())
	if err != nil {
		return err
	}
	return nil
}

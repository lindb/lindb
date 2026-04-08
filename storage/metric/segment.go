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
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	lmetrics "github.com/lindb/arrow/pkg/metrics"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/flush"
	"github.com/lindb/lindb/storage/metric/memdb"
	"github.com/lindb/lindb/storage/metric/tblstore/metricsdata"
	"github.com/lindb/lindb/storage/store"
)

type Segment struct {
	base.Segment
	partition *partition

	kvFamily kv.Family

	intervalCalc timeutil.IntervalCalculator
	interval     timeutil.Interval

	immutableMemDB memdb.MemoryDatabase
	mutableMemDB   memdb.MemoryDatabase

	mutex sync.Mutex

	lastReadTime *atomic.Int64

	isFlushing     atomic.Bool
	flushCondition sync.WaitGroup
	lastFlushTime  int64

	createdTime int64 // unix nanoseconds, set at construction time

	statistics *metrics.FamilyStatistics
	logger     logger.Logger
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	interval := partition.PartitionInterval()
	intervalCalc := interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)
	familySlot := intervalCalc.CalcFamily(timestamp, segmentTime)
	family := fmt.Sprintf("%d", familySlot)
	segmentPath := filepath.Join(partition.Path(), family)
	segmentStartTime := intervalCalc.CalcFamilyStartTime(segmentTime, familySlot)
	timeRange := timeutil.TimeRange{
		Start: segmentStartTime,
		End:   intervalCalc.CalcFamilyEndTime(segmentStartTime),
	}
	// FIXME: close kv store if load segment fail
	kvFamily := partition.kvStore.GetFamily(family)
	if kvFamily == nil {
		// create kv family
		var err error
		familyOption := kv.FamilyOption{
			CompactThreshold: 0,
			NumOfPoints:      timeRange.NumOfPoints(interval),
			Merger:           string(metricsdata.MetricDataMerger),
		}
		kvFamily, err = partition.kvStore.CreateFamily(family, familyOption)
		if err != nil {
			return nil, err
		}
	}
	shard := partition.shard
	db := shard.Database().(*Database)
	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeRange,
			Interval:  interval,
			Path:      segmentPath,
			WALs:      make(map[models.NodeID]store.WriteAheadLog),
		},
		partition:    partition,
		kvFamily:     kvFamily,
		interval:     interval,
		intervalCalc: intervalCalc,
		lastReadTime: atomic.NewInt64(fasttime.UnixMilliseconds()),
		createdTime:  time.Now().UnixNano(),
		statistics:   metrics.NewFamilyStatistics(db.Name(), shard.ShardID().String()),
		logger:       logger.GetLogger("Metric", "Segment"),
	}

	seg.CreateWriteAheadLog = func(p string) (store.WriteAheadLog, error) {
		return store.CreateWriteAheadLog(p, seg)
	}

	snapshot := kvFamily.GetSnapshot()
	defer snapshot.Close()

	ackSequences := snapshot.GetCurrent().GetSequences()
	if err := seg.LoadWALs(ackSequences); err != nil {
		return nil, err
	}

	// register with flush tracker so the checker can schedule flush for this segment
	flush.GetMemDBTracker().Register(seg)

	return seg, nil
}

func (s *Segment) Partition() store.Partition {
	return s.partition
}

// MutableMemDBInfo implements flush.FlushableSegment.
// Returns nil if there is no mutable memory database to flush.
func (s *Segment) MutableMemDBInfo() *flush.MemDBInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mutableMemDB == nil {
		return nil
	}
	created := time.Unix(0, s.createdTime)
	return &flush.MemDBInfo{
		MemSize:     s.mutableMemDB.MemSize(),
		CreatedTime: s.createdTime,
		SegmentTime: s.TimeRange.Start,
		NumOfRows:   int(s.mutableMemDB.NumOfSeries()),
		Uptime:      time.Since(created),
	}
}

// SegmentKey implements flush.FlushableSegment.
// Returns a globally unique key: "{dbName}/{shardID}/{segmentTime}".
func (s *Segment) SegmentKey() string {
	shard := s.partition.shard
	db := shard.Database().(*Database)
	return flush.MakeSegmentKey(db.Name(), shard.ShardID().String(), s.TimeRange.Start)
}

// SegmentMeta implements flush.FlushableSegment.
// Returns metadata for flush scheduling decisions.
func (s *Segment) SegmentMeta() flush.SegmentMeta {
	shard := s.partition.shard
	db := shard.Database().(*Database)
	dbOpt := db.GetOption().Option

	ahead, behind := dbOpt.GetAcceptWritableRange()

	var sizeThresholdBytes int64
	if dbOpt.Data.SizeThreshold > 0 {
		sizeThresholdBytes = dbOpt.Data.SizeThreshold * 1024 * 1024 // MB → bytes
	}
	var timeThresholdNano int64
	if dbOpt.Data.TimeThreshold > 0 {
		timeThresholdNano = dbOpt.Data.TimeThreshold * int64(time.Millisecond) // ms → ns
	}

	return flush.SegmentMeta{
		DatabaseName:         db.Name(),
		ShardID:              shard.ShardID().String(),
		SizeThresholdBytes:   sizeThresholdBytes,
		TimeThresholdNano:    timeThresholdNano,
		Ahead:                ahead,
		Behind:               behind,
		SegmentOutRangeDelay: dbOpt.Data.SegmentOutRangeDelay,
	}
}

// GetOrCreateMemoryDatabase returns memory database by given segment time.
func (s *Segment) GetOrCreateMemoryDatabase(segmentTime int64) (memdb.MemoryDatabase, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mutableMemDB == nil {
		shard := s.partition.shard
		newDB, err := memdb.NewMemoryDatabase(&memdb.MemoryDatabaseCfg{
			SegmentTime:   segmentTime,
			IntervalCalc:  s.intervalCalc,
			Interval:      s.interval,
			Name:          shard.Database().Name(),
			IndexDatabase: shard.MemIndexDB(),
			BufferMgr:     shard.BufferManager(),
		})
		if err != nil {
			return nil, err
		}
		s.mutableMemDB = newDB
		s.statistics.ActiveMemDBs.Incr()
	}
	return s.mutableMemDB, nil
}

// Flush implements store.Segment.
func (s *Segment) Flush() error {
	if s.isFlushing.CompareAndSwap(false, true) {
		defer func() {
			// mark flush job complete, notify
			s.flushCondition.Done()
			s.isFlushing.Store(false)
		}()

		// 1. mark flush job doing
		s.flushCondition.Add(1)

		startTime := time.Now()

		// add lock when switch memory database
		s.mutex.Lock()
		if s.immutableMemDB != nil || s.mutableMemDB == nil || s.mutableMemDB.NumOfSeries() == 0 {
			// if immutable memory database not nil or no data need flush, return it
			s.mutex.Unlock()
			return nil
		}
		waitingFlushMemDB := s.mutableMemDB
		s.immutableMemDB = waitingFlushMemDB
		s.mutableMemDB = nil
		// mark mutable memory database nil, write data will be created
		waitingFlushMemDB.MarkReadOnly()

		immutableSeq := make(map[models.NodeID]int64)
		for leader, wal := range s.WALs {
			immutableSeq[leader] = wal.SwitchSequence(leader)
		}
		s.mutex.Unlock()

		if err := s.flushMemoryDatabase(immutableSeq, waitingFlushMemDB); err != nil {
			return err
		}

		// flush success, mark immutable memory database nil
		s.mutex.Lock()
		s.immutableMemDB = nil
		// save persisted sequence, ack replica sequence in flushMemoryDatabase func
		for leader, wal := range s.WALs {
			wal.PersistSequence(leader)
		}

		s.mutex.Unlock()

		endTime := time.Now()
		s.lastFlushTime = endTime.UnixMilli()
		s.logger.Info("flush memory database successfully",
			logger.String("segment", s.Path),
			logger.String("flush-duration", endTime.Sub(startTime).String()),
			logger.Int64("segmentTime", s.TimeRange.Start),
			logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
	}

	// another flush process is running
	return nil
}

// Write implements store.Segment.
// It parses the payload as Arrow IPC (produced by MetricBuilder.Bytes()).
func (s *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	reader, err := lmetrics.NewReader(msg)
	if err != nil {
		return 0, err
	}
	defer reader.Release()
	numRows := reader.NumRows()
	if numRows == 0 {
		return 0, nil
	}
	db, dbErr := s.GetOrCreateMemoryDatabase(s.TimeRange.Start)
	if dbErr != nil {
		s.statistics.WriteMetricFailures.Add(float64(numRows))
		return 0, dbErr
	}
	db.AcquireWrite()
	defer func() {
		s.statistics.WriteBatches.Incr()
		db.CompleteWrite()
	}()

	for row := range numRows {
		if writeErr := db.WriteArrow(reader, row); writeErr != nil {
			s.statistics.WriteMetricFailures.Incr()
			s.logger.Error("failed writing arrow row",
				logger.String("segment", s.Path), logger.Error(writeErr))
		} else {
			s.statistics.WriteMetrics.Incr()
		}
	}
	return numRows, nil
}

// Close implements store.Segment.
func (s *Segment) Close() error {
	// unregister from flush tracker before closing so the checker stops scheduling this segment
	flush.GetMemDBTracker().Unregister(s)

	s.logger.Info("starting close data segment", logger.String("segment", s.Path))
	start := time.Now()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.flushCondition.Wait()

	if s.immutableMemDB != nil {
		immutableSequence := make(map[models.NodeID]int64)
		for leader, wal := range s.WALs {
			immutableSequence[leader] = wal.GetImmutableSequence(leader)
		}
		if err := s.flushMemoryDatabase(immutableSequence, s.immutableMemDB); err != nil {
			return err
		}
	}
	if s.mutableMemDB != nil {
		sequences := make(map[models.NodeID]int64)
		for leader, wal := range s.WALs {
			sequences[leader] = wal.GetSequence(leader)
		}
		if err := s.flushMemoryDatabase(sequences, s.mutableMemDB); err != nil {
			return err
		}
	}

	// FIXME:
	// GetFamilyManager().RemoveFamily(f)
	s.statistics.ActiveFamilies.Decr()

	s.logger.Info("close data segment complete", logger.String("segment", s.Path), logger.Any("cost", time.Since(start)))
	return nil
}

// Filter filters the data based on metric/version/seriesIDs,
// if it finds data then returns the FilterResultSet, else returns nil
func (s *Segment) Filter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	s.lastReadTime.Store(fasttime.UnixMilliseconds())
	memRS, err := s.memoryFilter(ctx)
	if !errors.Is(err, constants.ErrNotFound) && err != nil {
		// FIXME: ignore not found??
		return nil, err
	}
	fileRS, err := s.fileFilter(ctx)
	if !errors.Is(err, constants.ErrNotFound) && err != nil {
		return nil, err
	}
	resultSet = append(resultSet, memRS...)
	resultSet = append(resultSet, fileRS...)
	return
}

func (s *Segment) memoryFilter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	memFilter := func(memDB memdb.MemoryDatabase) error {
		rs, err := memDB.Filter(ctx)
		if err != nil {
			return err
		}
		resultSet = append(resultSet, rs...)
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mutableMemDB != nil {
		if err := memFilter(s.mutableMemDB); err != nil {
			return nil, err
		}
	}
	if s.immutableMemDB != nil {
		if err := memFilter(s.immutableMemDB); err != nil {
			return nil, err
		}
	}
	return
}

func (s *Segment) fileFilter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	snapShot := s.kvFamily.GetSnapshot()
	defer func() {
		if err != nil || len(resultSet) == 0 {
			// if not find metrics data or has error, close snapshot directly
			snapShot.Close()
		}
	}()
	metricKey := uint32(ctx.MetricID)
	readers, err := snapShot.FindReaders(metricKey)
	if err != nil {
		s.logger.Error("filter data family error", logger.Error(err))
		return nil, err
	}
	querySlotRange := s.interval.CalcSlotRange(s.TimeRange.Start, ctx.TimeRange)
	var metricReaders []metricsdata.MetricReader
	for _, reader := range readers {
		value, err0 := reader.Get(metricKey)
		// metric data not found
		if err0 != nil {
			continue
		}
		r, err := metricsdata.NewReader(reader.Path(), value)
		if err != nil {
			return nil, err
		}
		storageTimeRange := r.GetTimeRange()
		if storageTimeRange.Overlap(querySlotRange) {
			metricReaders = append(metricReaders, r)
		}
	}
	if len(metricReaders) == 0 {
		return nil, nil
	}
	filter := metricsdata.NewFilter(s.TimeRange.Start, s.interval, querySlotRange, snapShot, metricReaders)
	return filter.Filter(ctx.SeriesIDs, ctx.Fields)
}

// flushMemoryDatabase flushes memory database to disk.
func (s *Segment) flushMemoryDatabase(sequences map[models.NodeID]int64, memDB memdb.MemoryDatabase) error {
	startTime := time.Now()
	flusher := s.kvFamily.NewFlusher()
	defer func() {
		flusher.Release()
		s.statistics.MemDBFlushDuration.UpdateSince(startTime)
	}()

	for leader, seq := range sequences {
		if seq < 0 {
			// skip invalid sequence
			continue
		}
		flusher.Sequence(int32(leader), seq)
	}

	dataFlusher, err := metricsdata.NewFlusher(flusher)
	if err != nil {
		return err
	}
	// flush family data
	if err := memDB.FlushFamilyTo(dataFlusher); err != nil {
		s.logger.Error("failed to flush memory database",
			logger.String("segment", s.Path),
			logger.Int64("memDBSize", memDB.MemSize()))
		s.statistics.MemDBFlushFailures.Incr()
		return err
	}

	s.statistics.ActiveMemDBs.Decr()

	if err := memDB.Close(); err != nil {
		// ignore close memory database err, if not maybe write duplicate data into file storage
		s.logger.Warn("failed to close memory database",
			logger.String("segment", s.Path),
			logger.Int64("memDBSize", memDB.MemSize()))
		return nil
	}

	return nil
}

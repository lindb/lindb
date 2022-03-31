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
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./family.go -destination=./family_mock.go -package=tsdb

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// Indicator returns data family indicator's string.
	Indicator() string
	// Shard returns shard.
	Shard() Shard
	// Interval returns the interval data family's interval
	Interval() timeutil.Interval
	// FamilyTime returns the current family's time.
	FamilyTime() int64
	// TimeRange returns the data family's base time range
	TimeRange() timeutil.TimeRange
	// Family returns the raw kv family
	Family() kv.Family
	// WriteRows writes metric rows with same family in batch.
	WriteRows(rows []metric.StorageRow) error
	ValidateSequence(leader int32, seq int64) bool
	CommitSequence(leader int32, seq int64)
	AckSequence(leader int32, fn func(seq int64))

	NeedFlush() bool
	IsFlushing() bool
	Flush() error
	MemDBSize() int64

	// DataFilter filters data under data family based on query condition
	flow.DataFilter
	io.Closer
}

// dataFamily represents a wrapper of kv store's family with basic info
type dataFamily struct {
	indicator    string // database + shard + family time
	shard        Shard
	interval     timeutil.Interval
	intervalCalc timeutil.IntervalCalculator
	familyTime   int64
	timeRange    timeutil.TimeRange
	family       kv.Family

	mutableMemDB   memdb.MemoryDatabase
	immutableMemDB memdb.MemoryDatabase

	// leader => seq
	seq          map[int32]atomic.Int64
	immutableSeq map[int32]int64
	persistSeq   map[int32]atomic.Int64

	callbacks map[int32][]func(seq int64) // leader => callback

	isFlushing     atomic.Bool    // restrict flusher concurrency
	flushCondition sync.WaitGroup // flush condition

	mutex sync.Mutex

	statistics struct {
		writeBatches        *linmetric.BoundCounter
		writeMetrics        *linmetric.BoundCounter
		writeMetricFailures *linmetric.BoundCounter
		writeFields         *linmetric.BoundCounter
		memdbTotalSize      *linmetric.BoundGauge
		memdbNumber         *linmetric.BoundGauge
		memFlushTimer       *linmetric.BoundHistogram
		indexFlushTimer     *linmetric.BoundHistogram
	}

	logger *logger.Logger
}

// newDataFamily creates a data family storage unit
func newDataFamily(
	shard Shard,
	interval timeutil.Interval,
	timeRange timeutil.TimeRange,
	familyTime int64,
	family kv.Family,
) DataFamily {
	f := &dataFamily{
		shard:        shard,
		interval:     interval,
		intervalCalc: interval.Calculator(),
		timeRange:    timeRange,
		familyTime:   familyTime,
		family:       family,
		seq:          make(map[int32]atomic.Int64),
		persistSeq:   make(map[int32]atomic.Int64),
		callbacks:    make(map[int32][]func(seq int64)),

		logger: logger.GetLogger("TSDB", "family"),
	}
	// get current persist write sequence
	snapshot := family.GetSnapshot()
	defer snapshot.Close()

	sequences := snapshot.GetCurrent().GetSequences()
	for leader, seq := range sequences {
		sequence := *atomic.NewInt64(seq)
		f.seq[leader] = sequence
		f.persistSeq[leader] = sequence
	}

	dbName := shard.Database().Name()
	shardIDStr := strconv.Itoa(int(shard.ShardID()))

	f.statistics.writeBatches = writeBatchesVec.WithTagValues(dbName, shardIDStr)
	f.statistics.writeMetrics = writeMetricsVec.WithTagValues(dbName, shardIDStr)
	f.statistics.writeMetricFailures = metricMetricFailuresVec.WithTagValues(dbName, shardIDStr)
	f.statistics.writeFields = writeFieldsVec.WithTagValues(dbName, shardIDStr)
	f.statistics.memdbTotalSize = memdbTotalSizeVec.WithTagValues(dbName, shardIDStr)
	f.statistics.memdbNumber = memdbNumberVec.WithTagValues(dbName, shardIDStr)
	f.statistics.memFlushTimer = memFlushTimerVec.WithTagValues(dbName, shardIDStr)
	f.statistics.indexFlushTimer = indexFlushTimerVec.WithTagValues(dbName, shardIDStr)

	f.indicator = fmt.Sprintf("%s/%s/%d", dbName, shardIDStr, familyTime)

	// add data family into global family manager
	GetFamilyManager().AddFamily(f)
	return f
}

// Indicator returns data family indicator's string.
func (f *dataFamily) Indicator() string {
	return f.indicator
}

func (f *dataFamily) Shard() Shard {
	return f.shard
}

// Interval returns the data family's interval
func (f *dataFamily) Interval() timeutil.Interval {
	return f.interval
}

// TimeRange returns the data family's base time range
func (f *dataFamily) TimeRange() timeutil.TimeRange {
	return f.timeRange
}

// Family returns the kv store's family
func (f *dataFamily) Family() kv.Family {
	return f.family
}

func (f *dataFamily) FamilyTime() int64 {
	return f.familyTime
}

func (f *dataFamily) NeedFlush() bool {
	if f.IsFlushing() {
		return false
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.immutableMemDB != nil {
		// check immutable memory database, make sure it is nil
		return false
	}
	if f.mutableMemDB == nil || f.mutableMemDB.Size() <= 0 {
		// no data
		return false
	}

	// check memory database's uptime
	ttl := config.GlobalStorageConfig().TSDB.MutableMemDBTTL.Duration()
	if f.mutableMemDB.Uptime() >= ttl {
		f.logger.Info("memory database is expired, need do flush job",
			logger.String("family", f.indicator),
			logger.String("uptime", f.mutableMemDB.Uptime().String()),
			logger.String("mutable-memdb-ttl", ttl.String()),
		)
		return true
	}

	// check memory database's heap size
	maxMemDBSize := int64(config.GlobalStorageConfig().TSDB.MaxMemDBSize) // TODO need cfg
	if f.mutableMemDB.MemSize() >= maxMemDBSize {
		f.logger.Info("memory database is above memory threshold, need do flush job",
			logger.String("family", f.indicator),
			logger.String("uptime", f.mutableMemDB.Uptime().String()),
			logger.String("memdb-size", ltoml.Size(f.mutableMemDB.MemSize()).String()),
			logger.Int64("max-memdb-size", maxMemDBSize),
		)
		return true
	}

	// TODO need change metric
	// f.statistics.memdbNumber.Update(float64(len(s.families.Entries())))
	// f.statistics.memdbTotalSize.Update(float64(s.MemDBTotalSize()))

	return false
}

func (f *dataFamily) IsFlushing() bool {
	return f.isFlushing.Load()
}

func (f *dataFamily) Flush() error {
	if f.isFlushing.CAS(false, true) {
		defer func() {
			// TODO add commit kv meta after ack successfully
			// mark flush job complete, notify
			f.flushCondition.Done()
			f.isFlushing.Store(false)
		}()

		// 1. mark flush job doing
		f.flushCondition.Add(1)

		startTime := time.Now()
		// TODO flush index first????

		// add lock when switch memory database
		f.mutex.Lock()
		if f.immutableMemDB != nil || f.mutableMemDB == nil || f.mutableMemDB.Size() == 0 {
			// if immutable memory database not nil or no data need flush, return it
			f.mutex.Unlock()
			return nil
		}
		waitingFlushMemDB := f.mutableMemDB
		f.immutableMemDB = waitingFlushMemDB
		f.mutableMemDB = nil // mark mutable memory database nil, write data will be created
		immutableSeq := make(map[int32]int64)
		for leader, seq := range f.seq {
			immutableSeq[leader] = seq.Load()
		}
		f.immutableSeq = immutableSeq
		f.mutex.Unlock()

		if err := f.flushMemoryDatabase(immutableSeq, waitingFlushMemDB); err != nil {
			return err
		}

		// flush success, mark immutable memory database nil
		f.mutex.Lock()
		f.immutableMemDB = nil
		f.immutableSeq = nil
		for leader, seq := range immutableSeq {
			f.seq[leader] = *atomic.NewInt64(seq)
		}
		for leader, fns := range f.callbacks {
			seq, ok := f.seq[leader]
			if ok {
				s := seq.Load()
				for _, fn := range fns {
					fn(s)
				}
			}
		}
		f.mutex.Unlock()

		endTime := time.Now()
		f.logger.Info("flush memory database successfully",
			logger.String("family", f.indicator),
			logger.String("flush-duration", endTime.Sub(startTime).String()),
			logger.Int64("familyTime", f.familyTime),
			logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
		f.statistics.memFlushTimer.UpdateDuration(endTime.Sub(startTime))
	}

	// another flush process is running
	return nil
}

func (f *dataFamily) MemDBSize() int64 {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.mutableMemDB != nil {
		return f.mutableMemDB.MemSize()
	}
	return 0
}

// Filter filters the data based on metric/version/seriesIDs,
// if it finds data then returns the FilterResultSet, else returns nil
func (f *dataFamily) Filter(executeCtx *flow.ShardExecuteContext) (resultSet []flow.FilterResultSet, err error) {
	memRS, err := f.memoryFilter(executeCtx)
	if err != nil {
		return nil, err
	}
	fileRS, err := f.fileFilter(executeCtx)
	if err != nil {
		return nil, err
	}
	resultSet = append(resultSet, memRS...)
	resultSet = append(resultSet, fileRS...)
	return
}

func (f *dataFamily) memoryFilter(shardExecuteContext *flow.ShardExecuteContext) (resultSet []flow.FilterResultSet, err error) {
	memFilter := func(memDB memdb.MemoryDatabase) error {
		rs, err := memDB.Filter(shardExecuteContext)
		if err != nil {
			return err
		}
		resultSet = append(resultSet, rs...)
		return nil
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.mutableMemDB != nil {
		if err := memFilter(f.mutableMemDB); err != nil {
			return nil, err
		}
	}
	if f.immutableMemDB != nil {
		if err := memFilter(f.immutableMemDB); err != nil {
			return nil, err
		}
	}
	return
}

func (f *dataFamily) fileFilter(shardExecuteContext *flow.ShardExecuteContext) (resultSet []flow.FilterResultSet, err error) {
	snapShot := f.family.GetSnapshot()
	defer func() {
		if err != nil || len(resultSet) == 0 {
			// if not find metrics data or has error, close snapshot directly
			snapShot.Close()
		}
	}()
	metricKey := uint32(shardExecuteContext.StorageExecuteCtx.MetricID)
	readers, err := snapShot.FindReaders(metricKey)
	if err != nil {
		engineLogger.Error("filter data family error", logger.Error(err))
		return
	}
	// TODO need check time range???
	var metricReaders []metricsdata.MetricReader
	for _, reader := range readers {
		value, err0 := reader.Get(metricKey)
		// metric data not found
		if err0 != nil {
			continue
		}
		r, err := newReaderFunc(reader.Path(), value)
		if err != nil {
			return nil, err
		}
		metricReaders = append(metricReaders, r)
	}
	if len(metricReaders) == 0 {
		return
	}
	filter := newFilterFunc(f.timeRange.Start, snapShot, metricReaders)
	return filter.Filter(shardExecuteContext.SeriesIDsAfterFiltering, shardExecuteContext.StorageExecuteCtx.Fields)
}

func (f *dataFamily) WriteRows(rows []metric.StorageRow) error {
	defer f.statistics.writeBatches.Incr()

	if len(rows) == 0 {
		return nil
	}

	db, err := f.GetOrCreateMemoryDatabase(f.familyTime)
	if err != nil {
		// all rows are dropped
		f.statistics.writeMetricFailures.Add(float64(len(rows)))
		return err
	}
	db.AcquireWrite()
	defer db.CompleteWrite()

	releaseFunc := db.WithLock()
	defer releaseFunc()

	for idx := range rows {
		row := rows[idx]
		if !row.Writable {
			f.statistics.writeMetricFailures.Incr()
			continue
		}
		row.SlotIndex = uint16(f.intervalCalc.CalcSlot(
			row.Timestamp(),
			f.familyTime,
			f.interval.Int64()),
		)
		if err = db.WriteRow(&row); err == nil {
			f.statistics.writeMetrics.Incr()
			f.statistics.writeFields.Add(float64(len(row.FieldIDs)))
		} else {
			f.statistics.writeMetricFailures.Incr()
			f.logger.Error("failed writing row", logger.Error(err))
		}
	}
	// check memory database size in background flush checker job
	return nil
}

func (f *dataFamily) ValidateSequence(leader int32, seq int64) bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	seqForLeader, ok := f.seq[leader]
	if !ok {
		return true
	}
	return seq > seqForLeader.Load()
}

func (f *dataFamily) CommitSequence(leader int32, seq int64) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	seqForLeader := f.seq[leader]
	seqForLeader.Store(seq)
	f.seq[leader] = seqForLeader
}

func (f *dataFamily) AckSequence(leader int32, fn func(seq int64)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.callbacks[leader] = append(f.callbacks[leader], fn)

	seqForLeader, ok := f.seq[leader]
	if ok {
		// invoke ack sequence after register function, maybe some cases lost ack index.
		fn(seqForLeader.Load())
	}
}

// GetOrCreateMemoryDatabase returns memory database by given family time.
func (f *dataFamily) GetOrCreateMemoryDatabase(familyTime int64) (memdb.MemoryDatabase, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.mutableMemDB == nil {
		newDB, err := newMemoryDBFunc(memdb.MemoryDatabaseCfg{
			FamilyTime: familyTime,
			Name:       f.shard.Database().Name(),
			BufferMgr:  f.shard.BufferManager(),
		})
		if err != nil {
			return nil, err
		}
		f.mutableMemDB = newDB
	}
	return f.mutableMemDB, nil
}

func (f *dataFamily) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.flushCondition.Wait()

	if f.immutableMemDB != nil {
		if err := f.flushMemoryDatabase(f.immutableSeq, f.immutableMemDB); err != nil {
			return err
		}
	}
	if f.mutableMemDB != nil {
		sequences := make(map[int32]int64)
		for leader, seq := range f.seq {
			sequences[leader] = seq.Load()
		}
		if err := f.flushMemoryDatabase(sequences, f.mutableMemDB); err != nil {
			return err
		}
	}

	GetFamilyManager().RemoveFamily(f)
	return nil
}

func (f *dataFamily) flushMemoryDatabase(sequences map[int32]int64, memDB memdb.MemoryDatabase) error {
	flusher := f.family.NewFlusher()
	defer flusher.Release()

	for leader, seq := range sequences {
		flusher.Sequence(leader, seq)
	}

	dataFlusher, err := metricsdata.NewFlusher(flusher)
	if err != nil {
		return err
	}
	// flush family data
	if err := memDB.FlushFamilyTo(dataFlusher); err != nil {
		f.logger.Error("failed to flush memory database",
			logger.String("family", f.indicator),
			logger.Int64("memDBSize", memDB.MemSize()))
		return err
	}

	// invoke sequence ack callback
	for leader, seq := range sequences {
		callbacks, ok := f.callbacks[leader]
		if ok {
			for _, fn := range callbacks {
				fn(seq)
			}
		}
	}

	if err := memDB.Close(); err != nil {
		// ignore close memory database err, if not maybe write duplicate data into file storage
		f.logger.Warn("failed to close memory database",
			logger.String("family", f.indicator),
			logger.Int64("memDBSize", memDB.MemSize()))
		return nil
	}
	return nil
}

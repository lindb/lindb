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
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./family.go -destination=./family_mock.go -package=tsdb

// for testing
var (
	newReaderFunc = metricsdata.NewReader
	newFilterFunc = metricsdata.NewFilter
)

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// Indicator returns data family indicator's string.
	Indicator() string
	Shard() Shard
	// Interval returns the interval data family's interval
	Interval() timeutil.Interval
	// TimeRange returns the data family's base time range
	TimeRange() timeutil.TimeRange
	// Family returns the raw kv family
	Family() kv.Family
	// WriteRows writes metric rows with same family in batch.
	WriteRows(rows []metric.StorageRow) error
	ValidateSequence(seq int64) bool
	CommitSequence(seq int64)
	AckSequence(fn func(seq int64))

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

	seq          atomic.Int64
	immutableSeq atomic.Int64
	persistSeq   atomic.Int64

	callbacks []func(seq int64)

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
	snapshot := family.GetSnapshot()
	defer func() {
		snapshot.Close()
	}()
	// get current persist write sequence
	seq := snapshot.GetCurrent().GetSequence()

	f := &dataFamily{
		shard:        shard,
		interval:     interval,
		intervalCalc: interval.Calculator(),
		timeRange:    timeRange,
		familyTime:   familyTime,
		family:       family,
		seq:          *atomic.NewInt64(seq),
		immutableSeq: *atomic.NewInt64(seq),
		persistSeq:   *atomic.NewInt64(seq),

		logger: logger.GetLogger("tsdb", "family"),
	}
	dbName := shard.DatabaseName()
	shardIDStr := strconv.Itoa(int(shard.ShardID()))

	f.statistics.writeBatches = writeBatchesVec.WithTagValues(dbName, shardIDStr)
	f.statistics.writeMetrics = writeMetricsVec.WithTagValues(dbName, shardIDStr)
	f.statistics.writeMetricFailures = writeMetricFailuresVec.WithTagValues(dbName, shardIDStr)
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

	//TODO need change metric
	//f.statistics.memdbNumber.Update(float64(len(s.families.Entries())))
	//f.statistics.memdbTotalSize.Update(float64(s.MemDBTotalSize()))

	return false
}

func (f *dataFamily) IsFlushing() bool {
	return f.isFlushing.Load()
}

func (f *dataFamily) Flush() error {
	if f.isFlushing.CAS(false, true) {
		defer func() {
			//TODO add commit kv meta after ack successfully
			// mark flush job complete, notify
			f.flushCondition.Done()
			f.isFlushing.Store(false)
		}()

		// 1. mark flush job doing
		f.flushCondition.Add(1)

		startTime := time.Now()
		//TODO flush index first????

		// add lock when switch memory database
		f.mutex.Lock()
		if f.immutableMemDB != nil || f.mutableMemDB == nil || f.mutableMemDB.Size() <= 0 {
			// if immutable memory database not nil or no data need flush, return it
			f.mutex.Unlock()
			return nil
		}
		waitingFlushMemDB := f.mutableMemDB
		f.immutableMemDB = waitingFlushMemDB
		f.mutableMemDB = nil // mark mutable memory database nil, write data will be created
		f.immutableSeq.Store(f.seq.Load())
		f.mutex.Unlock()

		if err := f.flushMemoryDatabase(f.immutableSeq.Load(), waitingFlushMemDB); err != nil {
			return err
		}

		// flush success, mark immutable memory database nil
		var fns []func(seq int64)
		f.mutex.Lock()
		f.immutableMemDB = nil
		f.persistSeq.Store(f.immutableSeq.Load())
		// copy it
		fns = append(fns, f.callbacks...)
		f.mutex.Unlock()

		// invoke sequence ack callback
		seq := f.persistSeq.Load()
		for _, fn := range fns {
			fn(seq)
		}

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
// if finds data then returns the FilterResultSet, else returns nil
func (f *dataFamily) Filter(metricID uint32,
	seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
	fields field.Metas,
) (resultSet []flow.FilterResultSet, err error) {
	memRS, err := f.memoryFilter(metricID, seriesIDs, timeRange, fields)
	if err != nil {
		return nil, err
	}
	fileRS, err := f.fileFilter(metricID, seriesIDs, timeRange, fields)
	if err != nil {
		return nil, err
	}
	resultSet = append(resultSet, memRS...)
	resultSet = append(resultSet, fileRS...)
	return
}

func (f *dataFamily) memoryFilter(metricID uint32,
	seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
	fields field.Metas,
) (resultSet []flow.FilterResultSet, err error) {

	memFilter := func(memDB memdb.MemoryDatabase) error {
		rs, err := memDB.Filter(metricID, seriesIDs, timeRange, fields)
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

func (f *dataFamily) fileFilter(metricID uint32,
	seriesIDs *roaring.Bitmap, _ timeutil.TimeRange,
	fields field.Metas,
) (resultSet []flow.FilterResultSet, err error) {
	snapShot := f.family.GetSnapshot()
	defer func() {
		if err != nil || len(resultSet) == 0 {
			// if not find metrics data or has err, close snapshot directly
			snapShot.Close()
		}
	}()
	readers, err := snapShot.FindReaders(metricID)
	if err != nil {
		engineLogger.Error("filter data family error", logger.Error(err))
		return
	}
	//TODO need check time range???
	var metricReaders []metricsdata.MetricReader
	for _, reader := range readers {
		value, err := reader.Get(metricID)
		// metric data not found
		if err != nil {
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
	return filter.Filter(seriesIDs, fields)
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
		if !rows[idx].Writable {
			f.statistics.writeMetricFailures.Incr()
			continue
		}
		rows[idx].SlotIndex = uint16(f.intervalCalc.CalcSlot(
			rows[idx].Timestamp(),
			f.familyTime,
			f.interval.Int64()),
		)
		if err = db.WriteRow(&rows[idx]); err == nil {
			f.statistics.writeMetrics.Incr()
			f.statistics.writeFields.Add(float64(len(rows[idx].FieldIDs)))
		} else {
			f.statistics.writeMetricFailures.Incr()
			f.logger.Error("failed writing row", logger.Error(err))
		}
	}
	// check memory database size in background flush checker job
	return nil
}

func (f *dataFamily) ValidateSequence(seq int64) bool {
	return seq > f.seq.Load()
}

func (f *dataFamily) CommitSequence(seq int64) {
	f.seq.Store(seq)
}

func (f *dataFamily) AckSequence(fn func(seq int64)) {
	f.mutex.Lock()
	f.callbacks = append(f.callbacks, fn)
	f.mutex.Unlock()
}

// GetOrCreateMemoryDatabase returns memory database by given family time.
func (f *dataFamily) GetOrCreateMemoryDatabase(familyTime int64) (memdb.MemoryDatabase, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.mutableMemDB == nil {
		newDB, err := newMemoryDBFunc(memdb.MemoryDatabaseCfg{
			FamilyTime: familyTime,
			Name:       f.shard.DatabaseName(),
			TempPath:   filepath.Join(f.shard.Path(), filepath.Join(tempDir, fmt.Sprintf("%d", timeutil.Now()))),
		})
		if err != nil {
			return nil, err
		}
		f.mutableMemDB = newDB
	}
	return f.mutableMemDB, nil
}

func (f *dataFamily) Close() error {
	f.flushCondition.Wait()

	if f.immutableMemDB != nil {
		if err := f.flushMemoryDatabase(f.immutableSeq.Load(), f.immutableMemDB); err != nil {
			return err
		}
	}
	if f.mutableMemDB != nil {
		if err := f.flushMemoryDatabase(f.seq.Load(), f.mutableMemDB); err != nil {
			return err
		}
	}

	GetFamilyManager().RemoveFamily(f)
	return nil
}

func (f *dataFamily) flushMemoryDatabase(seq int64, memDB memdb.MemoryDatabase) error {
	flusher := f.family.NewFlusher()
	flusher.Sequence(seq)

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
	for _, fn := range f.callbacks {
		fn(seq)
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

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

package memdb

import (
	"bytes"
	"io"
	"math"
	"sync"
	"time"
	"unsafe"

	"go.uber.org/atomic"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("TSDB", "MemDB")

type nilPointer struct{}

var nilPointerSize *nilPointer

const (
	MetricStoreEntry = 8 + 4 + // ns+name hash(uint64) + metric store index(int)
		int64(unsafe.Sizeof(metricStore{})) + // metric store struct size
		2 + 2 + // metric slot range
		int64(unsafe.Sizeof(roaring.Bitmap{})) + // series ids
		int64(unsafe.Sizeof([][]uint32{})) // series ids
	HashSeriesMappingEntry  = 8 + 4              // tags hash + memory series id
	SeriesMappingEntry      = 4 * math.MaxUint16 // global series + memory series
	FieldMetaEntry          = int64(unsafe.Sizeof(field.Meta{}))
	FieldStoreEntry         = int64(unsafe.Sizeof(fieldStore{}))
	IntMapValuesEntry       = int64(unsafe.Sizeof([]uint16{})) + 2*math.MaxUint16
	NilPointerEntry         = int64(unsafe.Sizeof(nilPointerSize))
	IntMapStructValuesEntry = IntMapValuesEntry + math.MaxUint16*NilPointerEntry
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// MarkReadOnly marks memory database cannot writable.
	MarkReadOnly()
	// IsReadOnly returns memory database if it is readonly.
	IsReadOnly() bool
	// AcquireWrite acquires writing data points
	AcquireWrite()
	// WithLock retrieves the lock of memdb, and returns the release function
	WithLock() (release func())
	// WriteRow must be called after WithLock
	// Used for batch write
	WriteRow(row *metric.StorageRow) error
	// CompleteWrite completes writing data points
	CompleteWrite()
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int64
	// DataFilter filters the data based on condition
	flow.DataFilter
	// Closer closes the memory database resource
	io.Closer
	// FamilyTime returns the family time of this memdb
	FamilyTime() int64
	// Uptime returns duration since created
	Uptime() time.Duration
	// NumOfMetrics returns the number of metrics.
	NumOfMetrics() int
	// NumOfSeries returns the number of series.
	NumOfSeries() int
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	FamilyTime    int64
	Name          string
	BufferMgr     BufferManager
	MetaNotifier  func(notifier index.Notifier)
	IndexNotifier func(notifier index.Notifier)
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID uint32

	timeutil.SlotRange // start/end time slot, metric level flush context
	fieldIdx           int
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	cfg         *MemoryDatabaseCfg
	allocSize   atomic.Int64 // allocated size
	numOfSeries atomic.Int32 // num of series

	familyTime int64
	name       string

	metricStore      map[uint64]int    // ns+metirc name hash -> metric store index
	metricIndexStore *imap.IntMap[int] // metric id => metric store index
	stores           []mStoreINTF      // all metric stores

	timeSeriesStores []tStoreINTF   // time series id(memory unique) => field store
	sequence         *atomic.Uint32 // time series id generate sequence

	buf         DataPointBuffer
	createdTime int64
	statistics  *metrics.MemDBStatistics

	writeCondition sync.WaitGroup
	lock           sync.RWMutex // lock of create metric store

	readonly atomic.Bool
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg *MemoryDatabaseCfg) (MemoryDatabase, error) {
	buf, err := cfg.BufferMgr.AllocBuffer(cfg.FamilyTime)
	if err != nil {
		return nil, err
	}
	db := &memoryDatabase{
		cfg:              cfg,
		familyTime:       cfg.FamilyTime,
		name:             cfg.Name,
		buf:              buf,
		metricStore:      make(map[uint64]int),
		metricIndexStore: imap.NewIntMap[int](),
		timeSeriesStores: make([]tStoreINTF, math.MaxUint8), // pre-alloc all field's bucket
		sequence:         atomic.NewUint32(0),
		allocSize:        *atomic.NewInt64(0),
		createdTime:      fasttime.UnixNano(),
		statistics:       metrics.NewMemDBStatistics(cfg.Name),
	}
	for i := 0; i < math.MaxUint8; i++ {
		db.timeSeriesStores[i] = newTimeSeriesStore()
	}
	return db, nil
}

// MarkReadOnly marks memory database cannot writable.
func (md *memoryDatabase) MarkReadOnly() {
	md.readonly.Store(true)
}

// IsReadOnly returns memory database if it is readonly.
func (md *memoryDatabase) IsReadOnly() bool {
	return md.readonly.Load()
}

func (md *memoryDatabase) FamilyTime() int64 { return md.familyTime }

// getOrCreateMetricStore returns metric store, if not exist creates a new store.
func (md *memoryDatabase) getOrCreateMetricStore(row *metric.StorageRow) (mStore mStoreINTF, metricIdx int, created bool) {
	md.lock.Lock()
	defer md.lock.Unlock()
	hash := row.NameHash()
	metricIdx, ok := md.metricStore[hash]
	if ok {
		mStore = md.stores[metricIdx]
		return mStore, metricIdx, false
	}
	created = true
	metricIdx = len(md.stores)
	// not found need create new metric store
	mStore = newMetricStore()
	md.metricStore[hash] = metricIdx
	md.stores = append(md.stores, mStore)

	md.allocSize.Add(MetricStoreEntry)
	return
}

// AcquireWrite acquires writing data points
func (md *memoryDatabase) AcquireWrite() {
	md.writeCondition.Add(1)
}

// CompleteWrite completes writing data points
func (md *memoryDatabase) CompleteWrite() {
	md.writeCondition.Done()
}

func (md *memoryDatabase) WithLock() (release func()) {
	md.lock.Lock()
	return md.lock.Unlock
}

func (md *memoryDatabase) WriteRow(row *metric.StorageRow) error {
	mStore, metricIdx, created := md.getOrCreateMetricStore(row)
	if created {
		// notify metric metadata update, without lock
		notifier := index.GetMetaNotifier()
		notifier.Namespace = row.NamespaceStr()
		notifier.MetricName = string(row.Name())
		notifier.Callback = func(metricID uint32, err error) {
			if err != nil {
				memDBLogger.Error("generate metric id failure", logger.String("metric", string(row.Name())), logger.Error(err))
				return
			}
			// build index goroutine callback, with lock
			md.lock.Lock()
			size := len(md.metricIndexStore.Values())
			md.metricIndexStore.Put(metricID, metricIdx)

			if len(md.metricIndexStore.Values())-size > 0 {
				md.allocSize.Add(IntMapValuesEntry)
				md.allocSize.Add(SeriesMappingEntry)
			}
			md.lock.Unlock()
		}
		md.cfg.MetaNotifier(notifier)
	}

	var size int
	defer func() {
		md.allocSize.Add(int64(size))
	}()

	var timeSeriesID uint32 // unique id under memory database
	var newSeries bool

	md.lock.Lock()
	timeSeriesID = mStore.GenTStore(row.TagsHash(), func() uint32 {
		newSeries = true
		seriesIdx := md.sequence.Inc()

		// heap size
		md.allocSize.Add(HashSeriesMappingEntry)
		return seriesIdx
	})
	md.lock.Unlock()

	if newSeries {
		// notify time series update
		notifier := index.GetMetaNotifier()
		notifier.Namespace = row.NamespaceStr()
		notifier.MetricName = string(row.Name())
		notifier.TagHash = row.TagsHash()
		if row.TagsLen() > 0 {
			it := row.NewKeyValueIterator()
			for it.HasNext() {
				notifier.Tags = append(notifier.Tags, tag.NewTag(bytes.Clone(it.NextKey()), bytes.Clone(it.NextValue())))
			}
		}
		notifier.Callback = func(seriesID uint32, err error) {
			if err != nil {
				memDBLogger.Error("generate time series id failure", logger.String("metric", string(row.Name())), logger.Error(err))
				return
			}
			md.lock.Lock()
			newValueBucket := mStore.IndexTStore(seriesID, timeSeriesID)
			if newValueBucket {
				md.allocSize.Add(IntMapValuesEntry)
				md.allocSize.Add(SeriesMappingEntry)
			}
			md.lock.Unlock()
		}
		md.cfg.IndexNotifier(notifier)
		md.numOfSeries.Inc()
	}

	written := false
	afterWrite := func(writtenLinFieldSize int) {
		row.Fields++
		size += writtenLinFieldSize
		written = true
	}

	simpleFieldItr := row.NewSimpleFieldIterator()
	for simpleFieldItr.HasNext() {
		writtenLinFieldSize, err := md.writeLinField(
			row,
			simpleFieldItr.NextName(),
			simpleFieldItr.NextType(),
			simpleFieldItr.NextValue(),
			mStore, timeSeriesID,
		)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	compoundFieldItr, ok := row.NewCompoundFieldIterator()

	var (
		err                 error
		writtenLinFieldSize int
	)
	if !ok {
		goto End
	}

	// write histogram_min
	if compoundFieldItr.Min() > 0 {
		writtenLinFieldSize, err = md.writeLinField(
			row, compoundFieldItr.HistogramMinFieldName(),
			field.MinField, compoundFieldItr.Min(),
			mStore, timeSeriesID)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	// write histogram_max
	if compoundFieldItr.Max() > 0 {
		writtenLinFieldSize, err = md.writeLinField(
			row, compoundFieldItr.HistogramMaxFieldName(),
			field.MaxField, compoundFieldItr.Max(),
			mStore, timeSeriesID)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}
	// write histogram_sum
	writtenLinFieldSize, err = md.writeLinField(
		row, compoundFieldItr.HistogramSumFieldName(),
		field.SumField, compoundFieldItr.Sum(),
		mStore, timeSeriesID)
	if err != nil {
		return err
	}
	afterWrite(writtenLinFieldSize)
	// write histogram_count
	writtenLinFieldSize, err = md.writeLinField(
		row, compoundFieldItr.HistogramCountFieldName(),
		field.SumField, compoundFieldItr.Count(),
		mStore, timeSeriesID)
	if err != nil {
		return err
	}
	afterWrite(writtenLinFieldSize)

	// write __bucket_${boundary}
	// assume that length of ExplicitBounds equals to Values
	// data must be valid before write
	for compoundFieldItr.HasNextBucket() {
		writtenLinFieldSize, err = md.writeLinField(
			row, compoundFieldItr.BucketName(),
			field.HistogramField, compoundFieldItr.NextValue(),
			mStore, timeSeriesID)
		if err != nil {
			return err
		}
		afterWrite(writtenLinFieldSize)
	}

End:
	if written {
		mStore.SetSlot(row.SlotIndex)
	}
	return nil
}

func (md *memoryDatabase) writeLinField(
	row *metric.StorageRow,
	fName field.Name, fType field.Type, fValue float64,
	mStore mStoreINTF, ts uint32,
) (writtenSize int, err error) {
	var fm field.Meta
	var created bool

	md.lock.Lock()
	fm, created = mStore.GenField(fName, fType)
	md.lock.Unlock()

	if created {
		fieldNotifier := index.GetFieldNotifier()
		fieldNotifier.Field = fm
		fieldNotifier.Namespace = row.NamespaceStr()
		fieldNotifier.MetricName = string(row.Name())
		fieldNotifier.Callback = func(fieldID field.ID, err error) {
			if err != nil {
				memDBLogger.Error("generate field id failure", logger.String("metric", string(row.Name())),
					logger.String("field", fm.Name.String()), logger.Error(err))
				return
			}
			md.lock.Lock()
			defer md.lock.Unlock()

			mStore.UpdateFieldMeta(fieldID, fm)
		}
		md.cfg.MetaNotifier(fieldNotifier)

		md.allocSize.Add(FieldMetaEntry)
		md.allocSize.Add(int64(len(fName)))
	}

	tsStore := md.timeSeriesStores[fm.Index]
	writtenSize, err = tsStore.Write(ts, fType, row.SlotIndex, fValue, func() (fStoreINTF, error) {
		buf, err0 := md.buf.AllocPage()
		if err0 != nil {
			md.statistics.AllocatePageFailures.Incr()
			return nil, err0
		}
		md.statistics.AllocatedPages.Incr()
		fStore := newFieldStore(buf)

		md.allocSize.Add(FieldStoreEntry) // field store size
		return fStore, nil
	})
	if err != nil {
		return writtenSize, err
	}
	return writtenSize, nil
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder.
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher) error {
	// waiting current writing complete
	md.writeCondition.Wait()

	flushCtx := &flushContext{}
	if err := md.metricIndexStore.WalkEntry(func(metricID uint32, storeIndex int) error {
		flushCtx.metricID = metricID
		mStore := md.stores[storeIndex]
		if err := mStore.FlushMetricsDataTo(flusher, flushCtx, func(memSeriesID uint32, fields field.Metas) error {
			for _, fm := range fields {
				tsStores := md.timeSeriesStores[fm.Index]
				fStore, ok := tsStores.Get(memSeriesID)
				if ok {
					if err := fStore.FlushFieldTo(flusher, fm, flushCtx); err != nil {
						return err
					}
				} else {
					// must flush nil data for metric has multi-field.
					// because each series need fill all field data in order.
					_ = flusher.FlushField(nil)
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return flusher.Close()
}

// Filter filters the data based on metric/seriesIDs,
// if it finds data then returns the flow.FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(shardExecuteContext *flow.ShardExecuteContext) (rs []flow.FilterResultSet, err error) {
	md.lock.RLock()
	defer md.lock.RUnlock()

	mStoreIdx, ok := md.metricIndexStore.Get(uint32(shardExecuteContext.StorageExecuteCtx.MetricID))
	if !ok {
		return
	}
	mStore := md.stores[mStoreIdx]
	querySlotRange := shardExecuteContext.StorageExecuteCtx.CalcSourceSlotRange(md.familyTime)

	storageSlotRange := mStore.GetSlotRange()
	if !storageSlotRange.Overlap(querySlotRange) {
		return
	}
	return mStore.Filter(shardExecuteContext, md)
}

// MemSize returns the time series database memory size
func (md *memoryDatabase) MemSize() int64 {
	return md.allocSize.Load()
}

// Close releases resources for current memory database.
func (md *memoryDatabase) Close() error {
	md.buf.Release()
	return nil
}

func (md *memoryDatabase) Uptime() time.Duration {
	return time.Duration(fasttime.UnixNano() - md.createdTime)
}

// NumOfMetrics returns the number of metrics.
func (md *memoryDatabase) NumOfMetrics() int {
	md.lock.RLock()
	defer md.lock.RUnlock()

	return len(md.metricStore)
}

// NumOfSeries returns the number of series.
func (md *memoryDatabase) NumOfSeries() int {
	return int(md.numOfSeries.Load())
}

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
	"io"
	"math"
	"sync"
	"time"
	"unsafe"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
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
	IntMapValuesEntry       = int64(unsafe.Sizeof([]uint16{})) + 2*math.MaxUint16
	NilPointerEntry         = int64(unsafe.Sizeof(nilPointerSize))
	IntMapStructValuesEntry = IntMapValuesEntry + math.MaxUint16*NilPointerEntry
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
// NOTE: only one goroutine does writing operator.
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
	// MemSize returns the memory-size of memory database.
	MemSize() int64
	// DataFilter filters the data based on condition
	flow.DataFilter
	// Closer closes the memory database resource
	io.Closer
	// FamilyTime returns the family time of this memdb
	FamilyTime() int64
	// CreatedTime returns created timestamp of family's memory database.
	CreatedTime() int64
	// Uptime returns duration since created
	Uptime() time.Duration
	// NumOfSeries returns the number of series.
	NumOfSeries() int
	// MemTimeSeriesIDs returns all memory time series ids under current database.
	MemTimeSeriesIDs() *roaring.Bitmap
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	IntervalCalc  timeutil.IntervalCalculator
	BufferMgr     BufferManager
	IndexDatabase IndexDatabase
	Name          string
	Interval      timeutil.Interval
	FamilyTime    int64
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	cfg     *MemoryDatabaseCfg
	indexDB IndexDatabase
	// time series stores structure:
	// field memory index => time series store
	// time series store: time series id(memory unique) => field write buffer(temp mmap)
	fieldWriteStores   sync.Map // field index => (memory time series id => data point write buffer)
	fieldCompressStore sync.Map // field index => (memory time series id => compact buffer)
	timeSeriesIDs      *roaring.Bitmap

	statistics *metrics.MemDBStatistics

	name           string
	writeCondition sync.WaitGroup

	familyTime  int64
	createdTime int64 // create time(ns)

	readonly atomic.Bool
	lock     sync.RWMutex // lock of create metric store
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg *MemoryDatabaseCfg) (MemoryDatabase, error) {
	db := &memoryDatabase{
		cfg:           cfg,
		indexDB:       cfg.IndexDatabase,
		familyTime:    cfg.FamilyTime,
		name:          cfg.Name,
		timeSeriesIDs: roaring.New(),
		createdTime:   fasttime.UnixNano(),
		statistics:    metrics.NewMemDBStatistics(cfg.Name),
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
	var (
		isNewSeries bool
		memSeriesID uint32 // unique id under memory database
	)

	timeSeriesIndex := md.indexDB.GetOrCreateTimeSeriesIndex(row)
	mStore, newMetric := md.indexDB.GetMetadataDatabase().GetOrCreateMetricMeta(row)

	tagsHash := row.TagsHash()

	// generate memory level unique time series id
	memSeriesID, isNewSeries = timeSeriesIndex.GenMemTimeSeriesID(tagsHash, md.indexDB.GenMemSeriesID)

	if isNewSeries {
		row.MemSeriesID = memSeriesID
		// notify index worker does index building
		md.indexDB.Notify(row)
	} else {
		row.Done()
	}
	slotIndex := uint16(md.cfg.IntervalCalc.CalcSlot(
		row.Timestamp(),
		md.familyTime,
		md.cfg.Interval.Int64()),
	)

	defer func() {
		if newMetric || len(row.Fields) > 0 {
			// notify meta worker does build metadata
			md.indexDB.GetMetadataDatabase().Notify(row)
		} else {
			row.Done()
		}

		timeSeriesIndex.StoreTimeRange(md.createdTime, slotIndex)
		md.timeSeriesIDs.Add(memSeriesID)
	}()

	simpleFieldItr := row.NewSimpleFieldIterator()
	for simpleFieldItr.HasNext() {
		if err := md.writeLinField(
			mStore, memSeriesID, row,
			slotIndex,
			simpleFieldItr.NextName(),
			simpleFieldItr.NextType(),
			simpleFieldItr.NextValue(),
		); err != nil {
			return err
		}
	}

	// write compound fields
	if err := md.writeCompoundField(row, mStore, memSeriesID, slotIndex); err != nil {
		return err
	}

	return nil
}

func (md *memoryDatabase) writeCompoundField(row *metric.StorageRow,
	mStore mStoreINTF, memSeriesID uint32, slotIndex uint16,
) error {
	compoundFieldItr, ok := row.NewCompoundFieldIterator()
	if !ok {
		return nil
	}
	// write histogram_min
	if err := md.writeLinField(
		mStore, memSeriesID, row, slotIndex, compoundFieldItr.HistogramMinFieldName(),
		field.MinField, compoundFieldItr.Min()); err != nil {
		return err
	}
	// write histogram_max
	if err := md.writeLinField(
		mStore, memSeriesID, row, slotIndex, compoundFieldItr.HistogramMaxFieldName(),
		field.MaxField, compoundFieldItr.Max()); err != nil {
		return err
	}
	sum := compoundFieldItr.Sum()
	// write histogram_sum
	if err := md.writeLinField(
		mStore, memSeriesID, row, slotIndex, compoundFieldItr.HistogramSumFieldName(),
		field.SumField, sum); err != nil {
		return err
	}
	// write histogram_count
	if err := md.writeLinField(
		mStore, memSeriesID, row, slotIndex, compoundFieldItr.HistogramCountFieldName(),
		field.SumField, compoundFieldItr.Count()); err != nil {
		return err
	}

	// write __bucket_${boundary}
	// assume that length of ExplicitBounds equals to Values
	// data must be valid before write
	for compoundFieldItr.HasNextBucket() {
		bucketValue := compoundFieldItr.NextValue()
		if bucketValue > 0 {
			if err := md.writeLinField(
				mStore, memSeriesID, row, slotIndex, compoundFieldItr.BucketName(),
				field.HistogramField, bucketValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func (md *memoryDatabase) getFieldWriteBuffer(fieldIndex uint8) (DataPointBuffer, error) {
	buf, ok := md.fieldWriteStores.Load(fieldIndex)
	if ok {
		return buf.(DataPointBuffer), nil
	}

	// alloc a new data point buffer
	newBuf, err := md.cfg.BufferMgr.AllocBuffer(md.cfg.FamilyTime)
	if err != nil {
		return nil, err
	}
	// cache data point buffer
	md.fieldWriteStores.Store(fieldIndex, newBuf)
	return newBuf, nil
}

func (md *memoryDatabase) getFieldCompressBuffer(memSeriesID uint32, fieldIndex uint8) []byte {
	store, ok := md.fieldCompressStore.Load(fieldIndex)
	if !ok {
		return nil
	}
	return (store.(CompressStore)).GetCompressBuffer(memSeriesID)
}

func (md *memoryDatabase) storeFieldComressBuffer(memSeriesID uint32, fieldIndex uint8, buf []byte) {
	var store CompressStore
	storeObj, ok := md.fieldCompressStore.Load(fieldIndex)
	if !ok {
		store = NewCompressStore()
		md.fieldCompressStore.Store(fieldIndex, store)
	} else {
		store = storeObj.(CompressStore)
	}
	store.StoreCompressBuffer(memSeriesID, buf)
}

func (md *memoryDatabase) writeLinField(
	mStore mStoreINTF,
	memSeriesID uint32, row *metric.StorageRow, slotIndex uint16,
	fName field.Name, fType field.Type, fValue float64,
) (err error) {
	var fm field.Meta
	fm, isNew := mStore.GenField(fName, fType)
	if isNew {
		row.Fields = append(row.Fields, fm)
	}
	var buf DataPointBuffer

	buf, err = md.getFieldWriteBuffer(fm.Index)
	if err != nil {
		return err
	}
	page, err := buf.GetOrCreatePage(memSeriesID)
	if err != nil {
		return err
	}

	// write data into buffer
	write(md, page, memSeriesID, fm.Index, fType, slotIndex, fValue)

	// record write metric field statistics
	row.WrittenFields++
	return nil
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder.
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher) error {
	// waiting current writing complete
	md.writeCondition.Wait()

	metaDB := md.indexDB.GetMetadataDatabase()
	metricIDs := metaDB.GetMetricIDs()
	metricIDsIt := metricIDs.Iterator()
	for metricIDsIt.HasNext() {
		metricID := metricIDsIt.Next()
		memMetricID, ok := metaDB.GetMemMetricID(metricID)
		if !ok {
			continue // flush next metric if memory metric meta not exist
		}
		mStore, ok := metaDB.GetMetricMeta(memMetricID)
		if !ok {
			continue // flush next metric if memory metric meta not exist
		}
		// shard level metric time series index, shared multi data families
		timeSeriesIndex, ok := md.indexDB.GetTimeSeriesIndex(memMetricID)
		if !ok {
			continue // flush next metric if time series index not exist
		}
		slotRange, ok := timeSeriesIndex.GetTimeRange(md.createdTime)
		if !ok {
			continue // flush next metric if not time range
		}
		timeSeriesIDs := timeSeriesIndex.MemTimeSeriesIDs()
		curMetricMemTimeSeriesIDs := roaring.FastAnd(timeSeriesIDs, md.timeSeriesIDs)
		if curMetricMemTimeSeriesIDs.IsEmpty() {
			continue // flush next metric if current metric no data written
		}
		var needFlushFields field.Metas // current memory database's fields
		var buffers []DataPointBuffer
		allFields := mStore.GetFields()
		for idx := range allFields {
			f := allFields[idx]
			if !f.Persisted {
				// ignore if field meta not persist
				continue
			}
			buf, ok := md.fieldWriteStores.Load(f.Index)
			if ok {
				buffer := buf.(DataPointBuffer)
				buffers = append(buffers, buffer)
				needFlushFields = append(needFlushFields, f)
			}
		}
		if len(buffers) == 0 {
			continue // flush next metric if temp buffers of field not exist
		}
		// prepare for flushing metric
		flusher.PrepareMetric(metricID, needFlushFields)
		// flush time series of metric
		if err := timeSeriesIndex.FlushMetricsDataTo(flusher, func(memSeriesID uint32) error {
			for idx, buf := range buffers {
				buf, ok := buf.GetPage(memSeriesID)
				if ok {
					// flush field data
					if err := flushFieldTo(md, memSeriesID, buf, *slotRange, flusher, idx, needFlushFields[idx]); err != nil {
						return err
					}
				} else {
					// TEST: need test
					// NOTE: must flush nil data for metric has multi-field.
					// because each series need fill all field data in order.
					_ = flusher.FlushField(nil)
				}
			}
			return nil
		}); err != nil {
			return err
		}

		if err := flusher.CommitMetric(*slotRange); err != nil {
			return err
		}
	}
	return flusher.Close()
}

// Filter filters the data based on metric/seriesIDs,
// if it finds data then returns the flow.FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(shardExecuteContext *flow.ShardExecuteContext) (rs []flow.FilterResultSet, err error) {
	memMetricID, ok := md.indexDB.GetMetadataDatabase().GetMemMetricID(uint32(shardExecuteContext.StorageExecuteCtx.MetricID))
	if !ok {
		// metric not found
		return
	}
	timeSeriesIndex, ok := md.indexDB.GetTimeSeriesIndex(memMetricID)
	if !ok {
		// time series not found
		return
	}
	storageSlotRange, ok := timeSeriesIndex.GetTimeRange(md.createdTime)
	if !ok {
		// no data(time range not exist)
		return
	}
	querySlotRange := shardExecuteContext.StorageExecuteCtx.CalcSourceSlotRange(md.familyTime)
	if !storageSlotRange.Overlap(querySlotRange) {
		// time range not match
		return
	}
	return md.filter(shardExecuteContext, memMetricID, storageSlotRange, timeSeriesIndex)
}

// MemSize returns the time series database memory size.
func (md *memoryDatabase) MemSize() (memSize int64) {
	md.fieldWriteStores.Range(func(key, value any) bool {
		memSize += (value.(DataPointBuffer)).BufferSize()
		return true
	})
	md.fieldCompressStore.Range(func(key, value any) bool {
		memSize += (value.(CompressStore)).MemSize()
		return true
	})
	return memSize
}

// CreatedTime returns created timestamp of family's memory database.
func (md *memoryDatabase) CreatedTime() int64 {
	return md.createdTime
}

// Close releases resources for current memory database.
func (md *memoryDatabase) Close() error {
	md.fieldWriteStores.Range(func(key, value any) bool {
		(value.(DataPointBuffer)).Release()
		return true
	})
	md.indexDB.Cleanup(md)
	return nil
}

func (md *memoryDatabase) Uptime() time.Duration {
	return time.Duration(fasttime.UnixNano() - md.createdTime)
}

// MemTimeSeriesIDs returns all memory time series ids under current database.
// NOTE: after database flush invoke.
func (md *memoryDatabase) MemTimeSeriesIDs() *roaring.Bitmap {
	return md.timeSeriesIDs
}

// NumOfSeries returns the number of series.
func (md *memoryDatabase) NumOfSeries() int {
	md.lock.RLock()
	defer md.lock.RUnlock()

	return int(md.timeSeriesIDs.GetCardinality())
}

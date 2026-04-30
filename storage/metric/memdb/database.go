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
	"fmt"
	"io"
	"math"
	"sync"
	"time"
	"unsafe"

	"github.com/cespare/xxhash/v2"
	lmetrics "github.com/lindb/arrow/pkg/metrics"
	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/storage/metric/tblstore/metricsdata"
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
	// WriteArrow writes a single metric row from an Arrow IPC reader directly,
	// without an intermediate StorageRow conversion.
	WriteArrow(reader *lmetrics.Reader, row int) error
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
	SegmentTime   int64
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
	return &memoryDatabase{
		cfg:           cfg,
		indexDB:       cfg.IndexDatabase,
		familyTime:    cfg.SegmentTime,
		name:          cfg.Name,
		timeSeriesIDs: roaring.New(),
		createdTime:   fasttime.UnixNano(),
		statistics:    metrics.NewMemDBStatistics(cfg.Name),
	}, nil
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

func (md *memoryDatabase) getFieldWriteBuffer(fm field.Meta, fType field.Type) (DataPointBuffer, error) {
	buf, ok := md.fieldWriteStores.Load(fm.Index)
	if ok {
		return buf.(DataPointBuffer), nil
	}

	// alloc a new data point buffer
	newBuf, err := md.cfg.BufferMgr.AllocBuffer(md.cfg.SegmentTime)
	if err != nil {
		md.statistics.AllocatePageFailures.Incr()
		return nil, err
	}

	md.statistics.AllocatedPages.Incr()
	// cache data point buffer
	md.fieldWriteStores.Store(fm.Index, newBuf)
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

// WriteArrow writes a single metric row from an Arrow IPC reader directly,
// without an intermediate StorageRow conversion.
func (md *memoryDatabase) WriteArrow(reader *lmetrics.Reader, row int) error {
	// compute nameHash: xxhash(namespace + name), matching BrokerRowProtoConverter.hashOfName()
	namespace := reader.Namespace(row)
	if namespace == "" {
		namespace = constants.DefaultNamespace
	}
	name := reader.Name(row)
	nameHash := computeNameHash(namespace, name)

	// get or create time series index and metric meta store
	timeSeriesIndex := md.indexDB.GetOrCreateTimeSeriesIndexByHash(nameHash)
	mStore, newMetric := md.indexDB.GetMetadataDatabase().GetOrCreateMetricMetaByHash(nameHash)

	// route by attribute hash (same label-set → same memory series)
	tagsHash := reader.AttrHash(row)
	memSeriesID, isNewSeries := timeSeriesIndex.GenMemTimeSeriesID(tagsHash, md.indexDB.GenMemSeriesID)

	if isNewSeries {
		// collect attributes for index event
		var attrs []struct{ key, value string }
		reader.Attributes(row, func(k, v string) {
			attrs = append(attrs, struct{ key, value string }{k, v})
		})
		// asynchronously build the persistent series index for this new series
		md.indexDB.Notify(&arrowIndexEvent{
			nameHash:    nameHash,
			memSeriesID: memSeriesID,
			namespace:   namespace,
			name:        name,
			attrHash:    tagsHash,
			attrs:       attrs,
		})
	}

	// compute slot index: timestamp is in nanoseconds, convert to ms first
	timestampMs := reader.Timestamp(row) / 1_000_000
	slotIndex := uint16(md.cfg.IntervalCalc.CalcSlot(timestampMs, md.familyTime, md.cfg.Interval.Int64()))

	// build namespace/name byte slices for metadata notification
	nsBytes := []byte(namespace)
	if len(nsBytes) == 0 {
		nsBytes = []byte("default")
	}
	nameBytes := []byte(name)

	var fieldMetas []field.Meta

	// write simple fields from Arrow reader using the direct (no-StorageRow) path
	reader.Fields(row, func(fname string, kind model.AggregationKind, value float64) {
		fType := aggKindToFieldType(kind)
		fm, isNew := md.writeLinFieldDirect(mStore, memSeriesID, slotIndex,
			field.Name(fname), fType, value)
		if isNew {
			fieldMetas = append(fieldMetas, fm)
		}
	})

	// write exemplars from Arrow reader
	reader.Exemplars(row, func(e *model.Exemplar) {
		_ = md.writeExemplarFieldArrow(mStore, memSeriesID, slotIndex,
			"__exemplar__", field.ExemplarField, e.TraceID, e.SpanID, e.Duration)
	})

	// notify metadata worker for persistent field ID assignment
	if newMetric || len(fieldMetas) > 0 {
		md.indexDB.GetMetadataDatabase().Notify(&arrowMetaEvent{
			nameHash:   nameHash,
			namespace:  nsBytes,
			name:       nameBytes,
			fieldMetas: fieldMetas,
		})
	}

	timeSeriesIndex.StoreTimeRange(md.createdTime, slotIndex)
	md.timeSeriesIDs.Add(memSeriesID)
	return nil
}

// computeNameHash returns xxhash(namespace + name), matching BrokerRowProtoConverter.hashOfName().
func computeNameHash(namespace, name string) uint64 {
	// reuse a small stack buffer to avoid heap allocation for common short strings
	// namespace can be empty; hash is then just xxhash(name)
	buf := make([]byte, 0, len(namespace)+len(name))
	buf = append(buf, namespace...)
	buf = append(buf, name...)
	return xxhash.Sum64(buf)
}

// aggKindToFieldType maps model.AggregationKind to the LinDB field.Type.
func aggKindToFieldType(kind model.AggregationKind) field.Type {
	switch kind {
	case model.AggregationSum:
		return field.SumField
	case model.AggregationMin:
		return field.MinField
	case model.AggregationMax:
		return field.MaxField
	case model.AggregationLast:
		return field.LastField
	case model.AggregationFirst:
		return field.FirstField
	default:
		return field.SumField
	}
}

// writeExemplarFieldArrow writes an exemplar field from raw byte slices (Arrow path).
func (md *memoryDatabase) writeExemplarFieldArrow(
	mStore mStoreINTF,
	memSeriesID uint32, slotIndex uint16,
	fName field.Name, fType field.Type,
	traceID, spanID []byte, duration int64,
) (err error) {
	var fm field.Meta
	fm, _ = mStore.GenField(fName, fType)
	var buf DataPointBuffer
	buf, err = md.getFieldWriteBuffer(fm, fType)
	if err != nil {
		return err
	}
	page, err := buf.GetOrCreateExemplarPage(memSeriesID)
	if err != nil {
		return err
	}
	page.write(slotIndex, traceID, spanID, duration)
	return nil
}

// writeLinFieldDirect writes a single field value without a StorageRow, used by the Arrow write path.
// It returns the field meta and whether the field was newly created.
func (md *memoryDatabase) writeLinFieldDirect(
	mStore mStoreINTF,
	memSeriesID uint32, slotIndex uint16,
	fName field.Name, fType field.Type, fValue float64,
) (fm field.Meta, isNew bool) {
	fm, isNew = mStore.GenField(fName, fType)
	buf, err := md.getFieldWriteBuffer(fm, fType)
	if err != nil {
		return fm, isNew
	}
	page, err := buf.GetOrCreatePage(memSeriesID)
	if err != nil {
		return fm, isNew
	}
	// write data into buffer using the shared write helper
	write(md, page, memSeriesID, fm.Index, fType, slotIndex, fValue)
	return fm, isNew
}

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
		var fieldWritten bool
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
				fieldWritten = false
				fm := needFlushFields[idx]
				if fm.Type.IsExemplar() {
					// flush exemplar field
					page, ok := buf.GetExemplarPage(memSeriesID)
					if ok {
						if err := page.flush(flusher); err != nil {
							return err
						}
						fieldWritten = true
					}
				} else {
					// flush normal field
					page, ok := buf.GetPage(memSeriesID)
					if ok {
						if err := flushFieldTo(md, memSeriesID, page, *slotRange, flusher, idx, fm); err != nil {
							return err
						}
						fieldWritten = true
					}
				}

				if !fieldWritten {
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
func (md *memoryDatabase) Filter(metricScanCtx *flow.MetricScanContext) (rs []flow.FilterResultSet, err error) {
	memMetricID, ok := md.indexDB.GetMetadataDatabase().GetMemMetricID(uint32(metricScanCtx.MetricID))
	if !ok {
		// metric not found
		fmt.Println("metric not found")
		return
	}
	timeSeriesIndex, ok := md.indexDB.GetTimeSeriesIndex(memMetricID)
	if !ok {
		// time series not found
		fmt.Println("series not found")
		return
	}
	storageSlotRange, ok := timeSeriesIndex.GetTimeRange(md.createdTime)
	if !ok {
		// no data(time range not exist)
		fmt.Println("time range not exist")
		return
	}
	querySlotRange := md.cfg.Interval.CalcSlotRange(md.familyTime, metricScanCtx.TimeRange)
	if !storageSlotRange.Overlap(querySlotRange) {
		// time range not match
		fmt.Println("time range not match")
		return
	}
	slotRange := storageSlotRange.Intersect(querySlotRange)
	return md.filter(metricScanCtx, memMetricID, slotRange, timeSeriesIndex)
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

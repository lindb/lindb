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
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("tsdb", "MemDB")

var (
	memDBScope                    = linmetric.NewScope("lindb.tsdb.memdb")
	unknownFieldTypeCounterVec    = memDBScope.NewDeltaCounterVec("unknown_field_type_counter", "db")
	generateFieldIDFailCounterVec = memDBScope.NewDeltaCounterVec("generate_field_id_fails", "db")
	writeMetricsCounterVec        = memDBScope.NewDeltaCounterVec("write_metrics", "db")
	writeMetricsFailure           = memDBScope.NewDeltaCounterVec("write_metric_failures", "db")
	writeFieldsCounterVec         = memDBScope.NewDeltaCounterVec("write_fields", "db")
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// AcquireWrite acquires writing data points
	AcquireWrite()
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(
		namespace, metricName string,
		metricID, seriesID uint32,
		slotIndex uint16,
		simpleFields []*protoMetricsV1.SimpleField,
		compoundField *protoMetricsV1.CompoundField,
	) (err error)
	// CompleteWrite completes writing data points
	CompleteWrite()
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int32
	// DataFilter filters the data based on condition
	flow.DataFilter
	// Closer closes the memory database resource
	io.Closer
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	FamilyTime int64
	Name       string
	Metadata   metadb.Metadata
	TempPath   string
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID uint32

	timeutil.SlotRange // start/end time slot, metric level flush context
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	familyTime int64
	name       string
	metadata   metadb.Metadata // metadata for assign metric id/field id

	mStores *MetricBucketStore // metric id => mStoreINTF
	buf     DataPointBuffer

	writeCondition sync.WaitGroup
	rwMutex        sync.RWMutex // lock of create metric store

	allocSize                atomic.Int32 // allocated size
	writeMetricsCounter      *linmetric.BoundDeltaCounter
	writeMetricFailures      *linmetric.BoundDeltaCounter
	writeFieldsCounter       *linmetric.BoundDeltaCounter
	generatedFieldIDFailures *linmetric.BoundDeltaCounter
	gotUnknownFields         *linmetric.BoundDeltaCounter
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg MemoryDatabaseCfg) (MemoryDatabase, error) {
	buf, err := newDataPointBuffer(cfg.TempPath)
	if err != nil {
		return nil, err
	}
	return &memoryDatabase{
		familyTime:               cfg.FamilyTime,
		name:                     cfg.Name,
		metadata:                 cfg.Metadata,
		buf:                      buf,
		mStores:                  NewMetricBucketStore(),
		allocSize:                *atomic.NewInt32(0),
		writeMetricsCounter:      writeMetricsCounterVec.WithTagValues(cfg.Name),
		writeMetricFailures:      writeMetricsFailure.WithTagValues(cfg.Name),
		writeFieldsCounter:       writeFieldsCounterVec.WithTagValues(cfg.Name),
		generatedFieldIDFailures: generateFieldIDFailCounterVec.WithTagValues(cfg.Name),
		gotUnknownFields:         unknownFieldTypeCounterVec.WithTagValues(cfg.Name),
	}, err
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricID uint32) (mStore mStoreINTF) {
	mStore, ok := md.mStores.Get(metricID)
	if !ok {
		// not found need create new metric store
		mStore = newMetricStore()
		md.allocSize.Add(emptyMStoreSize)
		md.mStores.Put(metricID, mStore)
	}
	// found metric store in current memory database
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

// Write writes metric-point to database.
func (md *memoryDatabase) Write(
	namespace, metricName string,
	metricID, seriesID uint32,
	slotIndex uint16,
	simpleFields []*protoMetricsV1.SimpleField,
	compoundField *protoMetricsV1.CompoundField,
) (err error) {
	md.rwMutex.Lock()
	defer md.rwMutex.Unlock()

	mStore := md.getOrCreateMStore(metricID)

	tStore, size := mStore.GetOrCreateTStore(seriesID)
	written := false

	if compoundField != nil {
		writtenSize, err := md.writeCompoundField(namespace, metricName, slotIndex,
			mStore, tStore, compoundField)
		if err != nil {
			md.writeMetricFailures.Incr()
			return err
		}
		size += writtenSize
		written = true
	}

	for _, SimpleField := range simpleFields {
		if protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED == SimpleField.Type {
			md.gotUnknownFields.Incr()
			continue
		}
		var (
			fieldType    field.Type
			isCumulative bool
		)
		switch SimpleField.Type {
		case protoMetricsV1.SimpleFieldType_DELTA_SUM:
			fieldType = field.SumField
		case protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM:
			fieldType = field.SumField
			isCumulative = true
		case protoMetricsV1.SimpleFieldType_GAUGE:
			fieldType = field.GaugeField
		default:
			md.gotUnknownFields.Incr()
			continue
		}
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			SimpleField.Name, fieldType, SimpleField.Value,
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return err
		}
		size += writtenLinFieldSize
		written = true
	}

	if written {
		mStore.SetSlot(slotIndex)
	}
	md.writeMetricsCounter.Incr()
	md.allocSize.Add(int32(size))
	return nil
}

func (md *memoryDatabase) writeCompoundField(
	namespace, metricName string,
	slotIndex uint16,
	mStore mStoreINTF, tStore tStoreINTF,
	compoundField *protoMetricsV1.CompoundField,
) (writtenSize int, err error) {
	isCumulative := false
	switch compoundField.Type {
	case protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM:
		isCumulative = true
	case protoMetricsV1.CompoundFieldType_DELTA_HISTOGRAM:
	default:
		md.gotUnknownFields.Incr()
		return 0, nil
	}
	// write histogram_min
	if compoundField.Min > 0 {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			field.HistogramConverter.MinFieldName, field.MinField, compoundField.Min,
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return writtenSize, err
		}
		writtenSize += writtenLinFieldSize
	}
	// write histogram_max
	if compoundField.Max > 0 {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			field.HistogramConverter.MaxFieldName, field.MaxField, compoundField.Max,
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return writtenSize, err
		}
		writtenSize += writtenLinFieldSize
	}
	// write histogram_count
	if compoundField.Max > 0 {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			field.HistogramConverter.CountFieldName, field.SumField, compoundField.Count,
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return writtenSize, err
		}
		writtenSize += writtenLinFieldSize
	}
	// write histogram_sum
	writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
		field.HistogramConverter.SumFieldName, field.SumField, compoundField.Sum,
		mStore, tStore, isCumulative,
	)
	if err != nil {
		return writtenSize, err
	}
	writtenSize += writtenLinFieldSize
	// write histogram_data
	// assume that length of ExplicitBounds equals to Values
	// data must be valid before write
	for idx := range compoundField.ExplicitBounds {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			field.HistogramConverter.BucketName(compoundField.ExplicitBounds[idx]),
			field.HistogramField, compoundField.Values[idx],
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return writtenSize, err
		}
		writtenSize += writtenLinFieldSize
	}
	return writtenSize, err
}

func (md *memoryDatabase) writeLinField(
	namespace, metricName string,
	slotIndex uint16,
	fieldName string, fieldType field.Type, fieldValue float64,
	mStore mStoreINTF, tStore tStoreINTF,
	isCumulativeField bool,
) (writtenSize int, err error) {
	fieldID, err := md.metadata.MetadataDatabase().GenFieldID(
		namespace, metricName, field.Name(fieldName), fieldType)
	if err != nil {
		md.generatedFieldIDFailures.Incr()
		md.writeMetricFailures.Incr()
		// ignore generate field-id error
		return 0, nil
	}
	fStore, ok := tStore.GetFStore(fieldID)
	if !ok {
		buf, err := md.buf.AllocPage()
		if err != nil {
			md.writeMetricFailures.Incr()
			return 0, err
		}
		if isCumulativeField {
			fStore = newCumulativeSumFieldStore(buf, fieldID)
		} else {
			fStore = newFieldStore(buf, fieldID)
		}
		writtenSize += tStore.InsertFStore(fStore)
		// if write data success, add field into metric level for cache
		mStore.AddField(fieldID, fieldType)
	}
	writtenSize += fStore.Write(fieldType, slotIndex, fieldValue)
	md.writeFieldsCounter.Incr()
	return writtenSize, nil
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder.
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher) error {
	// waiting current writing complete
	md.writeCondition.Wait()

	if err := md.mStores.WalkEntry(func(key uint32, value mStoreINTF) error {
		if err := value.FlushMetricsDataTo(flusher, flushContext{
			metricID: key,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return flusher.Commit()
}

// Filter filters the data based on metric/seriesIDs,
// if finds data then returns the flow.FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(metricID uint32,
	seriesIDs *roaring.Bitmap,
	timeRange timeutil.TimeRange,
	fields field.Metas,
) ([]flow.FilterResultSet, error) {
	md.rwMutex.RLock()
	defer md.rwMutex.RUnlock()

	mStore, ok := md.mStores.Get(metricID)
	if !ok {
		return nil, nil
	}
	//TODO filter slot range
	return mStore.Filter(md.familyTime, seriesIDs, fields)
}

// MemSize returns the time series database memory size
func (md *memoryDatabase) MemSize() int32 {
	return md.allocSize.Load()
}

// Close closes memory data point buffer
func (md *memoryDatabase) Close() error {
	return md.buf.Close()
}

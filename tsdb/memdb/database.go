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
	"strconv"
	"sync"
	"time"

	"github.com/lindb/roaring"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/monitoring"
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
	getUnknownFieldTypeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mem_get_unknown_field_type",
			Help: "Get unknown field type when write data.",
		},
		[]string{"db"},
	)
	generateFieldIDFailCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mem_generate_field_id_fail",
			Help: "Generate field id fail when write data.",
		},
		[]string{"db"},
	)
	writeDataPointCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mem_write_data_points",
			Help: "Write data points.",
		},
		[]string{"db"},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(getUnknownFieldTypeCounter)
	monitoring.StorageRegistry.MustRegister(generateFieldIDFailCounter)
	monitoring.StorageRegistry.MustRegister(writeDataPointCounter)
}

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
	reportTicker             time.Ticker
	writtenDataPoints        atomic.Int64
	generatedFieldIDFailures atomic.Int64
	gotUnknownFields         atomic.Int64
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
		reportTicker:             *time.NewTicker(time.Second * 10),
		writtenDataPoints:        *atomic.NewInt64(0),
		generatedFieldIDFailures: *atomic.NewInt64(0),
		gotUnknownFields:         *atomic.NewInt64(0),
	}, err
}

func (md *memoryDatabase) statsReporter() {
	// todo: use otel sdk reporter
	writeDataPointC := writeDataPointCounter.WithLabelValues(md.name)
	generateFieldIDFailC := generateFieldIDFailCounter.WithLabelValues(md.name)
	getUnknownFieldTypeC := getUnknownFieldTypeCounter.WithLabelValues(md.name)
	var (
		lastWrittenDataPoints        int64 = 0
		lastGeneratedFieldIDFailures int64 = 0
		lastGotUnknownFields         int64 = 0
	)
	for range md.reportTicker.C {
		writtenDataPoints := md.writtenDataPoints.Load()
		generatedFieldIDFailures := md.generatedFieldIDFailures.Load()
		gotUnknownFields := md.gotUnknownFields.Load()
		writeDataPointC.Add(float64(writtenDataPoints - lastWrittenDataPoints))
		generateFieldIDFailC.Add(float64(generatedFieldIDFailures - lastGeneratedFieldIDFailures))
		getUnknownFieldTypeC.Add(float64(gotUnknownFields - lastGotUnknownFields))

		lastWrittenDataPoints = writtenDataPoints
		lastGeneratedFieldIDFailures = generatedFieldIDFailures
		lastGotUnknownFields = gotUnknownFields
	}
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
			return err
		}
		size += writtenSize
		written = true
	}

	for _, SimpleField := range simpleFields {
		if protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED == SimpleField.Type {
			md.gotUnknownFields.Inc()
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
			md.gotUnknownFields.Inc()
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
		md.gotUnknownFields.Inc()
		return 0, nil
	}
	// write histogram_min
	if compoundField.Min > 0 {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			"HistogramMin", field.MinField, compoundField.Min,
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
			"HistogramMax", field.MaxField, compoundField.Max,
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
			"HistogramCount", field.SumField, compoundField.Count,
			mStore, tStore, isCumulative,
		)
		if err != nil {
			return writtenSize, err
		}
		writtenSize += writtenLinFieldSize
	}
	// write histogram_sum
	writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
		"HistogramSum", field.SumField, compoundField.Sum,
		mStore, tStore, isCumulative,
	)
	if err != nil {
		return writtenSize, err
	}
	writtenSize += writtenLinFieldSize
	// write histogram_data
	for idx := range compoundField.Values {
		writtenLinFieldSize, err := md.writeLinField(namespace, metricName, slotIndex,
			"Histogram"+strconv.Itoa(idx), field.HistogramField, compoundField.Values[idx],
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
		md.generatedFieldIDFailures.Inc()
		// ignore generate field-id error
		return 0, nil
	}
	md.writtenDataPoints.Inc()
	fStore, ok := tStore.GetFStore(fieldID)
	if !ok {
		buf, err := md.buf.AllocPage()
		if err != nil {
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
	seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
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
	md.reportTicker.Stop()
	return md.buf.Close()
}

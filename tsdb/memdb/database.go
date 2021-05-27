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
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
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
	Write(namespace, metricName string, metricID, seriesID uint32, slotIndex uint16, fields []*pb.Field) (err error)
	// CompleteWrite completes writing data points
	CompleteWrite()
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int32
	// flow.DataFilter filters the data based on condition
	flow.DataFilter
	// io.Closer closes the memory database resource
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

	allocSize atomic.Int32 // allocated size

	writeCondition sync.WaitGroup
	rwMutex        sync.RWMutex // lock of create metric store

	writeDataPointCounter      prometheus.Counter
	generateFieldIDFailCounter prometheus.Counter
	getUnknownFieldTypeCounter prometheus.Counter
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg MemoryDatabaseCfg) (MemoryDatabase, error) {
	buf, err := newDataPointBuffer(cfg.TempPath)
	if err != nil {
		return nil, err
	}
	return &memoryDatabase{
		familyTime:                 cfg.FamilyTime,
		name:                       cfg.Name,
		metadata:                   cfg.Metadata,
		buf:                        buf,
		mStores:                    NewMetricBucketStore(),
		allocSize:                  *atomic.NewInt32(0),
		writeDataPointCounter:      writeDataPointCounter.WithLabelValues(cfg.Name),
		generateFieldIDFailCounter: generateFieldIDFailCounter.WithLabelValues(cfg.Name),
		getUnknownFieldTypeCounter: getUnknownFieldTypeCounter.WithLabelValues(cfg.Name),
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
func (md *memoryDatabase) Write(namespace, metricName string,
	metricID, seriesID uint32,
	slotIndex uint16, fields []*pb.Field,
) (err error) {
	md.rwMutex.Lock()
	defer md.rwMutex.Unlock()

	mStore := md.getOrCreateMStore(metricID)

	tStore, size := mStore.GetOrCreateTStore(seriesID)
	written := false

	for _, f := range fields {
		fieldType := getFieldType(f)
		if fieldType == field.Unknown {
			md.getUnknownFieldTypeCounter.Inc()
			continue
		}
		fieldID, err := md.metadata.MetadataDatabase().GenFieldID(namespace, metricName, field.Name(f.Name), fieldType)
		if err != nil {
			md.generateFieldIDFailCounter.Inc()
			continue
		}
		md.writeDataPointCounter.Inc()
		pStore, ok := tStore.GetFStore(fieldID)
		if !ok {
			buf, err := md.buf.AllocPage()
			if err != nil {
				return err
			}
			pStore = newFieldStore(buf, fieldID)
			size += tStore.InsertFStore(pStore)
		}
		size += pStore.Write(fieldType, slotIndex, f.Value)

		// if write data success, add field into metric level for cache
		mStore.AddField(fieldID, fieldType)
		written = true
	}
	if written {
		mStore.SetSlot(slotIndex)
	}
	md.allocSize.Add(int32(size))
	return nil
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
	return md.buf.Close()
}

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
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	commontimeutil "github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestMemoryDatabase_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	indexDB := NewMockIndexDatabase(ctrl)
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		FamilyTime:    10,
		BufferMgr:     bufferMgr,
		IndexDatabase: indexDB,
	}
	mdINTF, err := NewMemoryDatabase(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	assert.Equal(t, int64(10), mdINTF.FamilyTime())
	assert.False(t, mdINTF.IsReadOnly())
	assert.Zero(t, mdINTF.NumOfSeries())
	assert.Zero(t, mdINTF.MemSize())
	l := mdINTF.WithLock()
	l()
	mdINTF.MarkReadOnly()
	assert.True(t, mdINTF.IsReadOnly())
	md := mdINTF.(*memoryDatabase)
	indexDB.EXPECT().Cleanup(md)
	err = mdINTF.Close()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.True(t, mdINTF.Uptime() > 0)
}

func TestDatabase_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := "./db_write"
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(1), nil).AnyTimes()
	metaDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any()).Return(field.ID(2), nil).AnyTimes()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	indexDB.EXPECT().GenSeriesID(gomock.Any(), gomock.Any()).Return(uint32(100), nil).AnyTimes()
	memMetaDB := NewMetadataDatabase(&models.DatabaseConfig{}, metaDB)
	memIndexDB := NewIndexDatabase(memMetaDB, indexDB)
	bufferMgr := NewBufferManager(path.Join(name, "buf"))
	interval := timeutil.Interval(10_000)
	now, _ := commontimeutil.ParseTimestamp("2023-01-01 22:23:00", commontimeutil.DataTimeFormat2)
	familyTime := interval.Calculator().CalcFamilyTime(now)
	cfg := &MemoryDatabaseCfg{
		FamilyTime:    familyTime,
		BufferMgr:     bufferMgr,
		IndexDatabase: memIndexDB,
		Interval:      interval,
		IntervalCalc:  interval.Calculator(),
	}
	db, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Timestamp: now,
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}

	row := protoToStorageRow(m)
	assert.NoError(t, db.WriteRow(row))
	// write old time, compress data
	m.Timestamp = now - 20_1000
	row = protoToStorageRow(m)
	assert.NoError(t, db.WriteRow(row))
	m.Tags = []*protoMetricsV1.KeyValue{{Key: "key2", Value: "value2"}}
	m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
		Name: "f2", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10,
	})
	assert.NoError(t, db.WriteRow(protoToStorageRow(m)))
	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test2",
		Namespace: "ns",
		Timestamp: now,
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f4", Type: protoMetricsV1.SimpleFieldType_LAST, Value: 10},
		},
		CompoundField: &protoMetricsV1.CompoundField{
			Min:            10,
			Max:            10,
			Sum:            10,
			Count:          10,
			ExplicitBounds: []float64{1, 1, 1, 1, 1, math.Inf(1) + 1},
			Values:         []float64{1, 1, 1, 1, 1, 1},
		},
	})
	assert.NoError(t, db.WriteRow(row))
	assert.NotZero(t, db.MemSize())
	assert.False(t, db.MemTimeSeriesIDs().IsEmpty())

	// wait meta/index update
	time.Sleep(500 * time.Millisecond)
	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			MetricID: 1,
			Fields:   field.Metas{{Name: "f1", Type: field.SumField}},
			Query: &stmt.Query{
				TimeRange:       timeutil.TimeRange{Start: now - 200, End: now + 200},
				StorageInterval: interval,
			},
		},
		SeriesIDsAfterFiltering: roaring.BitmapOf(0, 100, 2),
	}
	db.MarkReadOnly()
	rs, err := db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.True(t, strings.HasSuffix(rs[0].Identifier(), "/memory/readonly"))
	assert.Equal(t, cfg.FamilyTime, rs[0].FamilyTime())
	assert.True(t, rs[0].SlotRange().End > 0)
	assert.Equal(t, []uint32{100}, rs[0].SeriesIDs().ToArray())
	assert.Nil(t, rs[0].Load(&flow.DataLoadContext{
		SeriesIDHighKey: 10,
	}))
	assert.Nil(t, rs[0].Load(&flow.DataLoadContext{
		SeriesIDHighKey:       0,
		LowSeriesIDsContainer: roaring.BitmapOf(1000).GetContainer(0),
	}))
	loader := rs[0].Load(&flow.DataLoadContext{
		SeriesIDHighKey:       0,
		LowSeriesIDsContainer: roaring.BitmapOf(0, 1).GetContainer(0),
	})
	assert.NotNil(t, loader)
	loader.Load(&flow.DataLoadContext{
		SeriesIDHighKey:       0,
		LowSeriesIDsContainer: roaring.BitmapOf(0, 1).GetContainer(0),
		LowSeriesIDs:          []uint16{0, 1},
		DownSampling:          func(slotRange timeutil.SlotRange, seriesIdx uint16, fieldIdx int, getter encoding.TSDValueGetter) {},
	})
	rs[0].Close()
	fm := ctx.StorageExecuteCtx.Fields
	// field not match
	ctx.StorageExecuteCtx.Fields = field.Metas{{Name: "xxx"}}
	ctx.SeriesIDsAfterFiltering = roaring.BitmapOf(100)
	rs, err = db.Filter(ctx)
	assert.Error(t, err)
	assert.Len(t, rs, 0)
	ctx.StorageExecuteCtx.Fields = fm
	// series not match
	ctx.SeriesIDsAfterFiltering = roaring.BitmapOf(1000)
	rs, err = db.Filter(ctx)
	assert.Error(t, err)
	assert.Len(t, rs, 0)
	ctx.SeriesIDsAfterFiltering = roaring.BitmapOf(0, 100, 2)
	// timerange not match
	ctx.StorageExecuteCtx.Query.TimeRange = timeutil.TimeRange{
		Start: now - 1000_000,
		End:   now - 500_000,
	}
	rs, err = db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 0)

	ctx.StorageExecuteCtx.MetricID = 10
	rs, err = db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 0)

	kvStore, err := kv.GetStoreManager().CreateStore(path.Join(name, "data"), kv.StoreOption{Levels: 2})
	assert.NoError(t, err)
	family, err := kvStore.CreateFamily("10", kv.FamilyOption{
		Merger: string(metricsdata.MetricDataMerger),
	})
	assert.NoError(t, err)
	kvFlusher := family.NewFlusher()
	defer kvFlusher.Release()

	flusher, err := metricsdata.NewFlusher(kvFlusher)
	assert.NoError(t, err)

	assert.NoError(t, db.FlushFamilyTo(flusher))

	snapshot := family.GetSnapshot()
	defer snapshot.Close()
	c := 0
	assert.NoError(t, snapshot.Load(1, func(value []byte) error {
		c++
		r, err := metricsdata.NewReader("dd", value)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		return nil
	}))
	assert.Equal(t, 1, c)

	assert.NoError(t, db.Close())
}

func TestMemoryDatabase_AcquireWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	indexDB := NewMockIndexDatabase(ctrl)
	metaDB := NewMockMetadataDatabase(ctrl)
	indexDB.EXPECT().GetMetadataDatabase().Return(metaDB).AnyTimes()
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		BufferMgr:     bufferMgr,
		IndexDatabase: indexDB,
	}
	buf, err := newDataPointBuffer(filepath.Join(t.TempDir(), "db_dir"))
	assert.NoError(t, err)

	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(buf, nil).AnyTimes()

	mdINTF, err := NewMemoryDatabase(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	mdINTF.AcquireWrite()
	a := time.After(100 * time.Millisecond)
	go func() {
		<-a
		mdINTF.CompleteWrite()
	}()
	flusher := metricsdata.NewMockFlusher(ctrl)
	metaDB.EXPECT().GetMetricIDs().Return(roaring.New())
	flusher.EXPECT().Close().Return(nil)
	err = mdINTF.FlushFamilyTo(flusher)
	assert.NoError(t, err)
}

func TestDatabase_Filter_NoData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metaDB := NewMockMetadataDatabase(ctrl)
	indexDB := NewMockIndexDatabase(ctrl)
	indexDB.EXPECT().GetMetadataDatabase().Return(metaDB).AnyTimes()
	metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(10), true).AnyTimes()
	db := &memoryDatabase{
		indexDB: indexDB,
	}
	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			MetricID: 10,
		},
	}
	indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(nil, false)
	rs, err := db.Filter(ctx)
	assert.NoError(t, err)
	assert.Empty(t, rs)

	timeSeriesIndex := NewMockTimeSeriesIndex(ctrl)
	indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
	timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(nil, false)
	rs, err = db.Filter(ctx)
	assert.NoError(t, err)
	assert.Empty(t, rs)
}

func TestMemoryDatabase_Write_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := "./db_write_error"
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(1), nil).AnyTimes()
	metaDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any()).Return(field.ID(2), nil).AnyTimes()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	indexDB.EXPECT().GenSeriesID(gomock.Any(), gomock.Any()).Return(uint32(100), nil).AnyTimes()
	memMetaDB := NewMetadataDatabase(&models.DatabaseConfig{}, metaDB)
	memIndexDB := NewIndexDatabase(memMetaDB, indexDB)
	interval := timeutil.Interval(10_000)
	bufferMgr := NewMockBufferManager(ctrl)
	buf := NewMockDataPointBuffer(ctrl)
	now, _ := commontimeutil.ParseTimestamp("2023-01-01 22:23:00", commontimeutil.DataTimeFormat2)
	familyTime := interval.Calculator().CalcFamilyTime(now)
	cfg := MemoryDatabaseCfg{
		FamilyTime:    familyTime,
		BufferMgr:     bufferMgr,
		IndexDatabase: memIndexDB,
		IntervalCalc:  interval.Calculator(),
		Interval:      interval,
	}
	db, err := NewMemoryDatabase(&cfg)
	assert.NoError(t, err)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Timestamp: now,
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}
	// write sample field error
	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))

	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(buf, nil).AnyTimes()

	m.SimpleFields = nil
	m.CompoundField = &protoMetricsV1.CompoundField{
		Min:            10,
		Max:            10,
		Sum:            10,
		Count:          10,
		ExplicitBounds: []float64{1, 1, 1, 1, 1, math.Inf(1) + 1},
		Values:         []float64{1, 1, 1, 1, 1, 1},
	}
	// write histogram min error
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	// write histogram max error
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(make([]byte, 128), nil)
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	// write histogram sum error
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(make([]byte, 128), nil).MaxTimes(2)
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	// write histogram count error
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(make([]byte, 128), nil).MaxTimes(3)
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(make([]byte, 128), nil).MaxTimes(4)
	buf.EXPECT().GetOrCreatePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
}

func TestMemoryDatabase_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := NewMockIndexDatabase(ctrl)
	metaDB := NewMockMetadataDatabase(ctrl)
	mStore := NewMockmStoreINTF(ctrl)
	timeSeriesIndex := NewMockTimeSeriesIndex(ctrl)
	buf := NewMockDataPointBuffer(ctrl)
	indexDB.EXPECT().GetMetadataDatabase().Return(metaDB).AnyTimes()
	db := &memoryDatabase{
		indexDB:       indexDB,
		timeSeriesIDs: roaring.BitmapOf(1, 2),
	}
	db.fieldWriteStores.Store(uint8(1), buf)
	flusher := metricsdata.NewMockFlusher(ctrl)

	cases := []struct {
		prepare func()
		name    string
		wantErr bool
	}{
		{
			name: "metric ids not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.New())
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "memory metric id not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), false)
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "memory metric meta not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(nil, false)
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "time series index not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(nil, false)
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "time range not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(nil, false)
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "time series not match",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(&timeutil.SlotRange{}, true)
				timeSeriesIndex.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(100))
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "field data not found",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(&timeutil.SlotRange{}, true)
				timeSeriesIndex.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(1))
				mStore.EXPECT().GetFields().Return(nil)
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "flush field not persist",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(&timeutil.SlotRange{}, true)
				timeSeriesIndex.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(1))
				mStore.EXPECT().GetFields().Return(field.Metas{{Name: "test", Index: 1}})
				flusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "flush field data err",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(&timeutil.SlotRange{}, true)
				timeSeriesIndex.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(1))
				mStore.EXPECT().GetFields().Return(field.Metas{{Name: "test", Persisted: true, Index: 1}})
				flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any())
				timeSeriesIndex.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						tableFlusher metricsdata.Flusher,
						flushFields func(memSeriesID uint32) error,
					) error {
						return flushFields(100)
					})
				flusher.EXPECT().GetEncoder(gomock.Any()).Return(encoding.NewTSDEncoder(100))
				buf.EXPECT().GetPage(gomock.Any()).Return(make([]byte, pageSize), true)
				flusher.EXPECT().FlushField(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "commit metric data err",
			prepare: func() {
				metaDB.EXPECT().GetMetricIDs().Return(roaring.BitmapOf(1))
				metaDB.EXPECT().GetMemMetricID(gomock.Any()).Return(uint64(0), true)
				metaDB.EXPECT().GetMetricMeta(gomock.Any()).Return(mStore, true)
				indexDB.EXPECT().GetTimeSeriesIndex(gomock.Any()).Return(timeSeriesIndex, true)
				timeSeriesIndex.EXPECT().GetTimeRange(gomock.Any()).Return(&timeutil.SlotRange{}, true)
				timeSeriesIndex.EXPECT().MemTimeSeriesIDs().Return(roaring.BitmapOf(1))
				mStore.EXPECT().GetFields().Return(field.Metas{{Name: "test", Persisted: true, Index: 1}})
				flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any())
				timeSeriesIndex.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(nil)
				flusher.EXPECT().CommitMetric(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.prepare != nil {
				c.prepare()
			}
			err := db.FlushFamilyTo(flusher)
			if c.wantErr != (err != nil) {
				t.Fatalf("run %s fail", c.name)
			}
		})
	}
}

func protoToStorageRow(m *protoMetricsV1.Metric) *metric.StorageRow {
	var ml protoMetricsV1.MetricList
	ml.Metrics = append(ml.Metrics, m)
	var buf bytes.Buffer
	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	_, err := converter.MarshalProtoMetricListV1To(ml, &buf)
	if err != nil {
		panic(err)
	}

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()[0]
}

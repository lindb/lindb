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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/imap"
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

	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		FamilyTime: 10,
		BufferMgr:  bufferMgr,
	}
	buf := NewMockDataPointBuffer(ctrl)
	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(buf, nil)
	mdINTF, err := NewMemoryDatabase(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	assert.Equal(t, int64(10), mdINTF.FamilyTime())
	assert.False(t, mdINTF.IsReadOnly())
	assert.Zero(t, mdINTF.NumOfMetrics())
	assert.Zero(t, mdINTF.NumOfSeries())
	assert.Zero(t, mdINTF.MemSize())
	l := mdINTF.WithLock()
	l()
	mdINTF.MarkReadOnly()
	assert.True(t, mdINTF.IsReadOnly())
	buf.EXPECT().Release()
	err = mdINTF.Close()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.True(t, mdINTF.Uptime() > 0)

	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(nil, fmt.Errorf("err"))
	mdINTF, err = NewMemoryDatabase(&cfg)
	assert.Error(t, err)
	assert.Nil(t, mdINTF)
}

func TestDatabase_Write(t *testing.T) {
	name := "./db_write"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	metaDB, err := index.NewMetricMetaDatabase(path.Join(name, "meta"))
	assert.NoError(t, err)
	indexDB, err := index.NewMetricIndexDatabase(path.Join(name, "index"), metaDB)
	assert.NoError(t, err)
	bufferMgr := NewBufferManager(path.Join(name, "buf"))
	cfg := &MemoryDatabaseCfg{
		FamilyTime:    100,
		BufferMgr:     bufferMgr,
		MetaNotifier:  metaDB.Notify,
		IndexNotifier: indexDB.Notify,
	}
	db, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}

	row := protoToStorageRow(m)
	row.SlotIndex = 10
	assert.NoError(t, db.WriteRow(row))
	m.Tags = []*protoMetricsV1.KeyValue{{Key: "key2", Value: "value2"}}
	m.SimpleFields = append(m.SimpleFields, &protoMetricsV1.SimpleField{
		Name: "f2", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10,
	})
	assert.NoError(t, db.WriteRow(protoToStorageRow(m)))
	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test2",
		Namespace: "ns",
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

	// wait meta/index update
	time.Sleep(500 * time.Millisecond)
	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			MetricID: 0,
			Fields:   field.Metas{{Name: "f1", Type: field.SumField}},
			Query: &stmt.Query{
				StorageInterval: timeutil.Interval(10 * time.Second),
			},
		},
		SeriesIDsAfterFiltering: roaring.BitmapOf(0, 1, 2),
	}
	db.MarkReadOnly()
	rs, err := db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.True(t, strings.HasSuffix(rs[0].Identifier(), "/memory/readonly"))
	assert.Equal(t, int64(100), rs[0].FamilyTime())
	assert.True(t, rs[0].SlotRange().End > 0)
	assert.Equal(t, []uint32{0, 1}, rs[0].SeriesIDs().ToArray())
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
	ctx.SeriesIDsAfterFiltering = roaring.BitmapOf(100)
	rs, err = db.Filter(ctx)
	assert.Error(t, err)
	assert.Len(t, rs, 0)
	ctx.SeriesIDsAfterFiltering = roaring.BitmapOf(0, 1, 2)
	// timerange not match
	ctx.StorageExecuteCtx.Query.TimeRange = timeutil.TimeRange{
		Start: time.Now().Add(time.Hour * 24).UnixMilli(),
		End:   time.Now().Add(time.Hour * 25).UnixMilli(),
	}
	rs, err = db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 0)

	ctx.StorageExecuteCtx.MetricID = 10
	rs, err = db.Filter(ctx)
	assert.NoError(t, err)
	assert.Len(t, rs, 0)

	// flush meta
	ch := make(chan error, 1)
	metaDB.Notify(&index.FlushNotifier{
		Callback: func(err error) {
			ch <- err
		},
	})
	err = <-ch
	assert.NoError(t, err)

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
	assert.NoError(t, snapshot.Load(0, func(value []byte) error {
		c++
		r, err := metricsdata.NewReader("dd", value)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		return nil
	}))
	assert.Equal(t, 1, c)
}

func TestMemoryDatabase_AcquireWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
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
	flusher.EXPECT().Close().Return(nil)
	err = mdINTF.FlushFamilyTo(flusher)
	assert.NoError(t, err)
}

func TestMemoryDatabase_Notify_Error(t *testing.T) {
	name := "./db_write"
	defer func() {
		_ = os.RemoveAll(name)
	}()
	bufferMgr := NewBufferManager(path.Join(name, "buf"))
	var wait sync.WaitGroup
	wait.Add(3)
	cfg := &MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
		MetaNotifier: func(notifier index.Notifier) {
			wait.Done()
			switch n := notifier.(type) {
			case *index.MetaNotifier:
				n.Callback(0, fmt.Errorf("err"))
			case *index.FieldNotifier:
				n.Callback(0, fmt.Errorf("err"))
			}
		},
		IndexNotifier: func(notifier index.Notifier) {
			wait.Done()
			n := notifier.(*index.MetaNotifier)
			n.Callback(0, fmt.Errorf("err"))
		},
	}
	db, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}
	assert.NoError(t, db.WriteRow(protoToStorageRow(m)))
	wait.Wait()
}

func TestMemoryDatabase_Write_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()
	bufferMgr := NewMockBufferManager(ctrl)
	buf := NewMockDataPointBuffer(ctrl)
	bufferMgr.EXPECT().AllocBuffer(gomock.Any()).Return(buf, nil)
	cfg := MemoryDatabaseCfg{
		BufferMgr:     bufferMgr,
		MetaNotifier:  func(notifier index.Notifier) {},
		IndexNotifier: func(notifier index.Notifier) {},
	}
	db, err := NewMemoryDatabase(&cfg)
	assert.NoError(t, err)
	m := &protoMetricsV1.Metric{
		Name: "test1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}
	// write sample field error
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
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
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	m.CompoundField.Min = 0
	// write histogram max error
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	m.CompoundField.Max = 0
	// write histogram sum error
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	// write histogram count error
	buf.EXPECT().AllocPage().Return(make([]byte, 128), nil)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
	buf.EXPECT().AllocPage().Return(make([]byte, 128), nil).MaxTimes(2)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	assert.Error(t, db.WriteRow(protoToStorageRow(m)))
}

func TestMemoryDatabase_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStore := NewMockmStoreINTF(ctrl)
	flusher := metricsdata.NewMockFlusher(ctrl)
	metricIndex := imap.NewIntMap[int]()
	metricIndex.Put(0, 0)
	timeSeries := newTimeSeriesStore()
	fStore := NewMockfStoreINTF(ctrl)
	fStore.EXPECT().Capacity().Return(10).MaxTimes(2)
	fStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any())
	_, _ = timeSeries.Write(0, field.SumField, 0, 0, func() (fStoreINTF, error) {
		return fStore, nil
	})
	db := &memoryDatabase{
		metricIndexStore: metricIndex,
		stores:           []mStoreINTF{mStore},
		timeSeriesStores: []tStoreINTF{timeSeries},
	}
	t.Run("flush metric store error", func(t *testing.T) {
		mStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		assert.Error(t, db.FlushFamilyTo(flusher))
	})
	t.Run("flush field error", func(t *testing.T) {
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		mStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(
			tableFlusher metricsdata.Flusher,
			flushCtx *flushContext,
			flushFields func(memSeriesID uint32, fields field.Metas) error,
		) error {
			return flushFields(0, field.Metas{{Index: 0}})
		})
		assert.Error(t, db.FlushFamilyTo(flusher))
	})
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
	return &br.Rows()[0]
}

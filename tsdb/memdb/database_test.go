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
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
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
	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil)
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	assert.Equal(t, int64(10), mdINTF.FamilyTime())
	buf.EXPECT().Release()
	err = mdINTF.Close()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.True(t, mdINTF.Uptime() > 0)

	bufferMgr.EXPECT().AllocBuffer().Return(nil, fmt.Errorf("err"))
	mdINTF, err = NewMemoryDatabase(cfg)
	assert.Error(t, err)
	assert.Nil(t, mdINTF)
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

	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil).AnyTimes()

	mdINTF, err := NewMemoryDatabase(cfg)
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

func protoToStorageRow(m *protoMetricsV1.Metric) *metric.StorageRow {
	var ml protoMetricsV1.MetricList
	ml.Metrics = append(ml.Metrics, m)
	var buf bytes.Buffer
	converter := metric.NewProtoConverter()
	_, _ = converter.MarshalProtoMetricListV1To(ml, &buf)

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return &br.Rows()[0]
}

func TestMemoryDatabase_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
	}
	buf, err := newDataPointBuffer(filepath.Join(t.TempDir(), "db_dir"))
	assert.NoError(t, err)

	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil).AnyTimes()
	// mock
	mockMStore := NewMockmStoreINTF(ctrl)
	tStore := NewMocktStoreINTF(ctrl)
	tStore.EXPECT().Capacity().Return(100).AnyTimes()
	fStore := NewMockfStoreINTF(ctrl)
	fStore.EXPECT().Capacity().Return(100).AnyTimes()
	mockMStore.EXPECT().Capacity().Return(100).AnyTimes()
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, false).AnyTimes()
	// build memory-database
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)
	assert.Zero(t, md.MemSize())

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	gomock.InOrder(
		tStore.EXPECT().GetFStore(gomock.Any()).Return(fStore, true),
		fStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetSlot(gomock.Any()).Times(1),
	)

	row := protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	})
	row.MetricID = 1
	row.SeriesID = 10
	row.SlotIndex = 1
	row.FieldIDs = []field.ID{10}
	assert.NoError(t, md.WriteRow(row))
	assert.NotZero(t, md.Size())

	// case 2: new metric store
	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	})
	row.MetricID = 20
	row.SeriesID = 20
	row.SlotIndex = 1
	row.FieldIDs = []field.ID{10}
	assert.NoError(t, md.WriteRow(row))

	// case 3: create new field store
	gomock.InOrder(
		tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false),
		tStore.EXPECT().InsertFStore(gomock.Any()),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetSlot(gomock.Any()),
	)
	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f4", Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: 10},
		},
	})
	row.MetricID = 1
	row.SeriesID = 10
	row.SlotIndex = 15
	row.FieldIDs = []field.ID{10}
	assert.NoError(t, md.WriteRow(row))
	assert.True(t, md.MemSize() > 0)

	// case4, write histogram field
	tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false).AnyTimes()
	tStore.EXPECT().InsertFStore(gomock.Any()).AnyTimes()
	mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()).AnyTimes()
	mockMStore.EXPECT().SetSlot(gomock.Any()).AnyTimes()
	releaseLock := md.WithLock()

	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f4", Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: 10},
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
	row.MetricID = 1
	row.SeriesID = 10
	row.SlotIndex = 15
	row.FieldIDs = []field.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	releaseLock()
	assert.NoError(t, md.WriteRow(row))
	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_Write_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
	}

	// mock
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Capacity().Return(100).AnyTimes()
	tStore := NewMocktStoreINTF(ctrl)
	tStore.EXPECT().Capacity().Return(100).AnyTimes()
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, false).AnyTimes()
	buf := NewMockDataPointBuffer(ctrl)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil).AnyTimes()
	// build memory-database
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false)

	row := protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f4", Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: 10},
		},
	})
	row.MetricID = 1
	row.SeriesID = 10
	row.SlotIndex = 15
	row.FieldIDs = []field.ID{10}
	assert.Error(t, md.WriteRow(row))

	buf.EXPECT().Release()
	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bufferMgr := NewMockBufferManager(ctrl)
	cfg := MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
	}
	buf, err := newDataPointBuffer(filepath.Join(t.TempDir(), "db_dir"))
	assert.NoError(t, err)

	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil).AnyTimes()
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)
	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().CommitMetric(gomock.Any()).Return(nil).AnyTimes()
	flusher.EXPECT().Close().Return(nil).AnyTimes()
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	md.mStores.Put(uint32(3333), mockMStore)

	// case 1: flusher ok
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(nil)
	err = md.FlushFamilyTo(flusher)
	assert.NoError(t, err)
	// case 2: flusher err
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = md.FlushFamilyTo(flusher)
	assert.Error(t, err)

	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bufferMgr := NewMockBufferManager(ctrl)
	now := timeutil.Now()
	cfg := MemoryDatabaseCfg{
		BufferMgr:  bufferMgr,
		FamilyTime: now,
	}
	buf, err := newDataPointBuffer(filepath.Join(t.TempDir(), "db_dir"))
	assert.NoError(t, err)

	bufferMgr.EXPECT().AllocBuffer().Return(buf, nil).AnyTimes()
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)

	// case 1: family not found
	rs, err := md.Filter(&flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			MetricID: metric.ID(3333),
			Query: &stmtpkg.Query{
				Interval:  timeutil.Interval(timeutil.OneMinute),
				TimeRange: timeutil.TimeRange{},
			},
			Fields: field.Metas{{ID: 1}},
		},
	})
	assert.NoError(t, err)
	assert.Nil(t, rs)
	// case 2: metric store not found
	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			MetricID: metric.ID(3333),
			Query: &stmtpkg.Query{
				Interval:  timeutil.Interval(timeutil.OneMinute),
				TimeRange: timeutil.TimeRange{Start: now - 10, End: now + 20},
			},
			Fields: field.Metas{{ID: 1}},
		},
	}
	rs, err = md.Filter(ctx)
	assert.NoError(t, err)
	assert.Nil(t, rs)
	// case 3: filter success
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Filter(gomock.Any(), gomock.Any()).Return([]flow.FilterResultSet{}, nil)
	md.mStores.Put(uint32(3333), mockMStore)
	mockMStore.EXPECT().GetSlotRange().Return(&timeutil.SlotRange{Start: 0, End: 60})
	rs, err = md.Filter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	// case 4: slot range not match
	mockMStore.EXPECT().GetSlotRange().Return(&timeutil.SlotRange{Start: 600, End: 600})
	rs, err = md.Filter(ctx)
	assert.NoError(t, err)
	assert.Nil(t, rs)

	err = md.Close()
	assert.NoError(t, err)
}

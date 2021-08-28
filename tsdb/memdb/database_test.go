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
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

const testDBPath = "test_db"

var cfg = MemoryDatabaseCfg{
	TempPath: testDBPath,
}

func TestMemoryDatabase_New(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
	}()

	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	err = mdINTF.Close()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.True(t, mdINTF.Uptime() > 0)

	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	mdINTF, err = NewMemoryDatabase(cfg)
	assert.Error(t, err)
	assert.Nil(t, mdINTF)
}

func TestMemoryDatabase_AcquireWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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

func TestMemoryDatabase_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()
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
	err = md.Write(&MetricPoint{
		MetricID:  1,
		SeriesID:  10,
		SlotIndex: 1,
		FieldIDs:  []field.ID{10},
		Proto: &protoMetricsV1.Metric{
			Name:      "test1",
			Namespace: "ns",
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
			},
		}})
	assert.NoError(t, err)
	// case 2: field type unknown
	err = md.Write(&MetricPoint{
		MetricID:  1,
		SeriesID:  10,
		SlotIndex: 1,
		FieldIDs:  []field.ID{10},
		Proto: &protoMetricsV1.Metric{
			Name:      "test1",
			Namespace: "ns",
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED, Value: 10},
			},
		}})
	assert.NoError(t, err)
	// case 3: new metric store
	err = md.Write(
		&MetricPoint{
			MetricID:  20,
			SeriesID:  20,
			SlotIndex: 1,
			FieldIDs:  []field.ID{10},
			Proto: &protoMetricsV1.Metric{
				Name:      "test1",
				Namespace: "ns",
				SimpleFields: []*protoMetricsV1.SimpleField{
					{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
				},
			}})
	assert.NoError(t, err)
	// case 4: create new field store
	gomock.InOrder(
		tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false),
		tStore.EXPECT().InsertFStore(gomock.Any()),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetSlot(gomock.Any()),
	)
	err = md.Write(
		&MetricPoint{
			MetricID:  1,
			SeriesID:  10,
			SlotIndex: 15,
			FieldIDs:  []field.ID{10},
			Proto: &protoMetricsV1.Metric{
				Name:      "test1",
				Namespace: "ns",
				SimpleFields: []*protoMetricsV1.SimpleField{
					{Name: "f4", Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: 10},
				},
			}})
	assert.NoError(t, err)
	assert.True(t, md.MemSize() > 0)
	// case5, write histogram field
	tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false).AnyTimes()
	tStore.EXPECT().InsertFStore(gomock.Any()).AnyTimes()
	mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()).AnyTimes()
	mockMStore.EXPECT().SetSlot(gomock.Any()).AnyTimes()
	releaseLock := md.WithLock()
	err = md.WriteWithoutLock(
		&MetricPoint{
			MetricID:  1,
			SeriesID:  10,
			SlotIndex: 15,
			FieldIDs:  []field.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Proto: &protoMetricsV1.Metric{
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
					ExplicitBounds: []float64{1, 1, 1, 1, math.Inf(1) + 1},
					Values:         []float64{1, 1, 1, 1, 1, 1},
					Type:           protoMetricsV1.CompoundFieldType_DELTA_HISTOGRAM,
				},
			}})

	releaseLock()
	assert.NoError(t, err)
	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_Write_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()

	// mock
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Capacity().Return(100).AnyTimes()
	tStore := NewMocktStoreINTF(ctrl)
	tStore.EXPECT().Capacity().Return(100).AnyTimes()
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, false).AnyTimes()
	// build memory-database
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	buf := NewMockDataPointBuffer(ctrl)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	md := mdINTF.(*memoryDatabase)
	md.buf = buf

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false)
	err = md.Write(
		&MetricPoint{
			MetricID:  1,
			SeriesID:  10,
			SlotIndex: 15,
			FieldIDs:  []field.ID{10},
			Proto: &protoMetricsV1.Metric{
				Name:      "test1",
				Namespace: "ns",
				SimpleFields: []*protoMetricsV1.SimpleField{
					{Name: "f4", Type: protoMetricsV1.SimpleFieldType_GAUGE, Value: 10},
				},
			}})
	assert.Error(t, err)

	buf.EXPECT().Close().Return(nil)
	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)

	// case 1: family not found
	rs, err := md.Filter(uint32(3333), nil, timeutil.TimeRange{}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.Nil(t, rs)
	now := timeutil.Now()
	// case 2: metric store not found
	rs, err = md.Filter(0, nil, timeutil.TimeRange{Start: now - 10, End: now + 20}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.Nil(t, rs)
	// case 3: filter success
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return([]flow.FilterResultSet{}, nil)
	md.mStores.Put(uint32(3333), mockMStore)
	rs, err = md.Filter(uint32(3333), nil, timeutil.TimeRange{Start: now - 10, End: now + 20}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	err = md.Close()
	assert.NoError(t, err)
}

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
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestTimeSeriesIndex_TimeRange(t *testing.T) {
	idx := NewTimeSeriesIndex()
	idx.StoreTimeRange(10, 20)
	timeRange, ok := idx.GetTimeRange(10)
	assert.True(t, ok)
	assert.Equal(t, timeRange, &timeutil.SlotRange{Start: 20, End: 20})

	idx.ClearTimeRange(10)
	timeRange, ok = idx.GetTimeRange(10)
	assert.False(t, ok)
	assert.Nil(t, timeRange)
}

func TestTimeSeriesIndex_FlushMetricsData(t *testing.T) {
	idx := NewTimeSeriesIndex()
	idx.IndexTimeSeries(10, 10)
	err := idx.FlushMetricsDataTo(nil, func(memSeriesID uint32) error { return fmt.Errorf("err") })
	assert.Error(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().FlushSeries(gomock.Any()).Return(fmt.Errorf("err"))
	err = idx.FlushMetricsDataTo(flusher, func(memSeriesID uint32) error { return nil })
	assert.Error(t, err)

	err = idx.(*timeSeriesIndex).flushMetricsDataTo(roaring.BitmapOf(2), nil, nil)
	assert.NoError(t, err)
}

func TestTimeSeriesIndex_GenMemTimeSeriesID(t *testing.T) {
	idx := NewTimeSeriesIndex()
	id, isNew := idx.GenMemTimeSeriesID(10, func() uint32 {
		return 100
	})
	assert.Equal(t, uint32(100), id)
	assert.True(t, isNew)

	id, isNew = idx.(*timeSeriesIndex).genMemTimeSeriesID(10, nil)
	assert.Equal(t, uint32(100), id)
	assert.False(t, isNew)
}

func TestTimeSeriesIndex_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("no data under time series index", func(t *testing.T) {
		idx := NewTimeSeriesIndex()
		idx.Load(nil, 0, timeutil.NewSlotRange(0, 100), []*fieldEntry{})
	})

	t.Run("no field data", func(t *testing.T) {
		idx := NewTimeSeriesIndex()
		// add data
		idx.IndexTimeSeries(10, 100)
		seriesIDs := roaring.BitmapOf(10)
		ctx := &flow.DataLoadContext{
			MinSeriesID:  10,
			MaxSeriesID:  1000,
			LowSeriesIDs: seriesIDs.GetContainerAtIndex(0).ToArray(),
		}
		idx.Load(ctx, 0, timeutil.NewSlotRange(0, 100), []*fieldEntry{})
	})

	t.Run("load data", func(t *testing.T) {
		idx := NewTimeSeriesIndex()
		// add data
		idx.IndexTimeSeries(10, 100)
		seriesIDs := roaring.BitmapOf(10)
		ctx := &flow.DataLoadContext{
			MinSeriesID:  10,
			MaxSeriesID:  1000,
			LowSeriesIDs: seriesIDs.GetContainerAtIndex(0).ToArray(),
			DownSampling: func(slotRange timeutil.SlotRange, seriesIdx uint16,
				fieldIdx int, getter encoding.TSDValueGetter) {
			},
			Decoder: encoding.GetTSDDecoder(),
		}
		compress := NewCompressStore()
		encoder := encoding.NewTSDEncoder(10)
		// case 1: encode with err
		encoder.AppendTime(bit.One)
		encoder.AppendValue(uint64(10))
		data, err := encoder.Bytes()
		assert.NoError(t, err)
		compress.StoreCompressBuffer(100, data)
		idx.Load(ctx, 0, timeutil.NewSlotRange(0, 100), []*fieldEntry{
			{
				compressBuf: compress,
			},
		})
		pageBuf := NewMockDataPointBuffer(ctrl)
		pageBuf.EXPECT().GetPage(gomock.Any()).Return([]byte{1, 2, 3}, true)
		idx.Load(ctx, 0, timeutil.NewSlotRange(0, 100), []*fieldEntry{
			{
				compressBuf: compress,
				pageBuf:     pageBuf,
			},
		})
	})
}

func TestTimeSeriesIndex_GC(t *testing.T) {
	index := NewTimeSeriesIndex()
	id := uint32(0)
	newID := func() uint32 {
		id++
		return id
	}
	index.GenMemTimeSeriesID(10, newID)
	index.GenMemTimeSeriesID(20, newID)
	index.GenMemTimeSeriesID(30, newID)
	index.GenMemTimeSeriesID(40, newID)

	index.IndexTimeSeries(100, 1)
	index.IndexTimeSeries(200, 2)
	index.IndexTimeSeries(300, 3)
	index.IndexTimeSeries(400, 4)

	index1 := index.(*timeSeriesIndex)
	index1.ExpireTimeSeriesIDs(roaring.BitmapOf(1, 2), 1000)
	index1.ExpireTimeSeriesIDs(roaring.BitmapOf(4), 100)

	index1.GC(500) // gc 4
	assert.Equal(t, roaring.BitmapOf(100, 200, 300, 400), index1.ids.Keys())
	assert.Equal(t, 4, index1.NumOfSeries())
	index1.GC(5000) // gc 1,2
	assert.Equal(t, 1, index1.NumOfSeries())
	assert.Equal(t, roaring.BitmapOf(300), index1.ids.Keys())
	index1.ExpireTimeSeriesIDs(roaring.BitmapOf(3), 1000)
	index.GenMemTimeSeriesID(30, newID)
	index1.GC(5000) // ignore gc 3
	assert.Equal(t, 1, index1.NumOfSeries())

	index1.ExpireTimeSeriesIDs(roaring.BitmapOf(3), 1000)
	index1.GC(5000) // gc 3
	assert.Equal(t, 0, index1.NumOfSeries())
	assert.True(t, index1.ids.Keys().IsEmpty())
}

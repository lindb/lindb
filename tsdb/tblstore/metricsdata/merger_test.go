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

package metricsdata

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func Test_NewMerger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := kv.NewMockFlusher(ctrl)
	flusher.EXPECT().StreamWriter().Return(nil, io.ErrClosedPipe)
	_, err := NewMerger(flusher)
	assert.Error(t, err)
}

func Test_Compact(t *testing.T) {
	flusher := kv.NewNopFlusher()
	mergerIntf, err := NewMerger(flusher)
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		assertMergeReentrant(t, flusher, mergerIntf)
	}
}

func assertMergeReentrant(t *testing.T, flusher kv.Flusher, mergerIntf kv.Merger) {
	err := mergerIntf.Merge(
		1,
		[][]byte{
			mockRealMetricBlock([]uint32{1, 2, 4}, 11, 15),
			mockRealMetricBlock([]uint32{2, 20}, 16, 20),
			mockRealMetricBlock([]uint32{2, 30}, 5, 10),
		})
	assert.Nil(t, err)

	r, err := NewReader("test", flusher.(*kv.NopFlusher).Bytes())
	assert.Nil(t, err)
	assert.Len(t, r.GetFields(), 2)
	assert.EqualValues(t, r.GetFields(), field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
	})
	for _, seriesID := range []uint32{1, 2, 4, 20, 30} {
		assert.True(t, r.GetSeriesIDs().Contains(seriesID))
	}

	assert.Equal(t, timeutil.SlotRange{Start: uint16(5), End: uint16(20)}, r.GetTimeRange())

	container := r.GetSeriesIDs().GetContainer(0)
	found := 0
	ctx := &flow.DataLoadContext{
		SeriesIDHighKey:       0,
		LowSeriesIDsContainer: container,
		ShardExecuteCtx: &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				Fields: r.GetFields(),
			},
		},
		DownSampling: func(slotRange timeutil.SlotRange, seriesIdx uint16, fieldIdx int, fieldData []byte) {
			found++
		},
	}
	loader := r.Load(ctx)
	// not exist
	ctx.LowSeriesIDsContainer = roaring.BitmapOf(0).GetContainerAtIndex(0)
	ctx.Grouping()
	loader.Load(ctx)
	assert.Equal(t, 0, found)
	// 11-15
	ctx.LowSeriesIDsContainer = roaring.BitmapOf(1).GetContainerAtIndex(0)
	ctx.Grouping()
	loader.Load(ctx)
	assert.Equal(t, 2, found)
}

func mockRealMetricBlock(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(10, field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
	})
	encoder := encoding.NewTSDEncoder(start)
	for i := start; i <= end; i++ {
		encoder.AppendTime(true)
		encoder.AppendValue(math.Float64bits(float64(i)))
	}
	data, _ := encoder.BytesWithoutTime()
	for _, seriesID := range seriesIDs {
		_ = flusher.FlushField(data)
		_ = flusher.FlushField(data)
		_ = flusher.FlushSeries(seriesID)
	}
	_ = flusher.CommitMetric(timeutil.SlotRange{Start: start, End: end})
	return nopKVFlusher.Bytes()
}

func TestMerger_Compact_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := NewMockFlusher(ctrl)
	seriesMerger := NewMockSeriesMerger(ctrl)
	nopFlusher := kv.NewNopFlusher()
	merge, _ := NewMerger(nopFlusher)
	m := merge.(*merger)
	m.dataFlusher = flusher
	m.seriesMerger = seriesMerger
	// case 1: new metricReader err
	err := merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, nopFlusher.Bytes())
	// case 2: series merge err
	flusher.EXPECT().PrepareMetric(uint32(1),
		field.Metas{{ID: 2, Type: field.SumField}, {ID: 10, Type: field.MinField}}).AnyTimes()

	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("err"))
	err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
			mockMetricMergeBlock([]uint32{2, 30}, 5, 5),
		})
	assert.Error(t, err)
	assert.Nil(t, nopFlusher.Bytes())
	// case 3: merge success
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(4)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushSeries(uint32(2)),
		flusher.EXPECT().FlushSeries(uint32(4)),
		flusher.EXPECT().FlushSeries(uint32(20)),
		flusher.EXPECT().CommitMetric(timeutil.SlotRange{Start: 10, End: 15}).Return(nil),
	)
	err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
		})
	assert.NoError(t, err)
	assert.False(t, len(nopFlusher.Bytes()) > 0) // data flush is mock
	// case 4: flush metric err
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)).Return(nil),
		flusher.EXPECT().CommitMetric(timeutil.SlotRange{Start: 10, End: 10}).Return(fmt.Errorf("err")),
	)
	err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1}, 10, 10),
		})
	assert.Error(t, err)
	assert.Nil(t, nopFlusher.Bytes())
}

func TestMerger_Rollup_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rollup := kv.NewMockRollup(ctrl)
	flusher := NewMockFlusher(ctrl)
	seriesMerger := NewMockSeriesMerger(ctrl)
	nopFlusher := kv.NewNopFlusher()
	merge, _ := NewMerger(nopFlusher)
	merge.Init(map[string]interface{}{kv.RollupContext: rollup})

	m := merge.(*merger)
	m.dataFlusher = flusher
	m.seriesMerger = seriesMerger
	// case 1: rollup merge success
	flusher.EXPECT().PrepareMetric(uint32(1),
		field.Metas{{ID: 2, Type: field.SumField}, {ID: 10, Type: field.MinField}}).AnyTimes()
	rollup.EXPECT().IntervalRatio().Return(uint16(10))
	rollup.EXPECT().GetTimestamp(uint16(10)).Return(int64(100))
	rollup.EXPECT().CalcSlot(int64(100)).Return(uint16(0))
	rollup.EXPECT().GetTimestamp(uint16(15)).Return(int64(150))
	rollup.EXPECT().CalcSlot(int64(150)).Return(uint16(0))
	rollup.EXPECT().BaseSlot().Return(uint16(10))
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(4)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushSeries(uint32(2)),
		flusher.EXPECT().FlushSeries(uint32(4)),
		flusher.EXPECT().FlushSeries(uint32(20)),
		flusher.EXPECT().CommitMetric(timeutil.SlotRange{Start: 0, End: 0}).Return(nil),
	)
	err := merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
		})
	assert.NoError(t, err)
	assert.False(t, len(nopFlusher.Bytes()) > 0) // data flush is mock
}

func mockMetricMergeBlock(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(10, field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
	})
	for _, seriesID := range seriesIDs {
		_ = flusher.FlushField([]byte{1, 2, 3})
		_ = flusher.FlushField([]byte{1, 2, 3})
		_ = flusher.FlushSeries(seriesID)
	}
	_ = flusher.CommitMetric(timeutil.SlotRange{Start: start, End: end})
	return nopKVFlusher.Bytes()
}

func mockMetricMergeBlockOneField(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(10, field.Metas{
		{ID: 2, Type: field.SumField},
	})
	for _, seriesID := range seriesIDs {
		_ = flusher.FlushField([]byte{1, 2, 3})
		_ = flusher.FlushSeries(seriesID)
	}
	_ = flusher.CommitMetric(timeutil.SlotRange{Start: start, End: end})
	return nopKVFlusher.Bytes()
}

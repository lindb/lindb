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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"
)

func TestMerger_Compact_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := NewMockFlusher(ctrl)
	seriesMerger := NewMockSeriesMerger(ctrl)
	merge := NewMerger()
	m := merge.(*merger)
	m.dataFlusher = flusher
	m.seriesMerger = seriesMerger
	// case 1: new metricReader err
	data, err := merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 2: series merge err
	flusher.EXPECT().FlushFieldMetas(field.Metas{{ID: 2, Type: field.SumField}, {ID: 10, Type: field.MinField}}).AnyTimes()
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("err"))
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
			mockMetricMergeBlock([]uint32{2, 30}, 5, 5),
		})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 3: merge success
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(4)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushSeries(uint32(2)),
		flusher.EXPECT().FlushSeries(uint32(4)),
		flusher.EXPECT().FlushSeries(uint32(20)),
		flusher.EXPECT().FlushMetric(uint32(1), uint16(10), uint16(15)).Return(nil),
	)
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
		})
	assert.NoError(t, err)
	assert.False(t, len(data) > 0) // data flush is mock
	// case 4: flush metric err
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushMetric(uint32(1), uint16(10), uint16(10)).Return(fmt.Errorf("err")),
	)
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1}, 10, 10),
		})
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestMerger_Rollup_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rollup := kv.NewMockRollup(ctrl)
	flusher := NewMockFlusher(ctrl)
	seriesMerger := NewMockSeriesMerger(ctrl)
	merge := NewMerger()
	merge.Init(map[string]interface{}{kv.RollupContext: rollup})

	m := merge.(*merger)
	m.dataFlusher = flusher
	m.seriesMerger = seriesMerger
	// case 1: rollup merge success
	flusher.EXPECT().FlushFieldMetas(field.Metas{{ID: 2, Type: field.SumField}, {ID: 10, Type: field.MinField}}).AnyTimes()
	rollup.EXPECT().IntervalRatio().Return(uint16(10))
	rollup.EXPECT().GetTimestamp(uint16(10)).Return(int64(100))
	rollup.EXPECT().CalcSlot(int64(100)).Return(uint16(0))
	rollup.EXPECT().GetTimestamp(uint16(15)).Return(int64(150))
	rollup.EXPECT().CalcSlot(int64(150)).Return(uint16(0))
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(4)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushSeries(uint32(2)),
		flusher.EXPECT().FlushSeries(uint32(4)),
		flusher.EXPECT().FlushSeries(uint32(20)),
		flusher.EXPECT().FlushMetric(uint32(1), uint16(0), uint16(0)).Return(nil),
	)
	data, err := merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
		})
	assert.NoError(t, err)
	assert.False(t, len(data) > 0) // data flush is mock
}

func mockMetricMergeBlock(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
	})
	for _, seriesID := range seriesIDs {
		flusher.FlushField([]byte{1, 2, 3})
		flusher.FlushField([]byte{1, 2, 3})
		flusher.FlushSeries(seriesID)
	}
	_ = flusher.FlushMetric(uint32(10), start, end)
	return nopKVFlusher.Bytes()
}

func mockMetricMergeBlockOneField(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 2, Type: field.SumField},
	})
	for _, seriesID := range seriesIDs {
		flusher.FlushField([]byte{1, 2, 3})
		flusher.FlushSeries(seriesID)
	}
	_ = flusher.FlushMetric(uint32(10), start, end)
	return nopKVFlusher.Bytes()
}

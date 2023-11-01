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

package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

func TestGroupingContext_BuildSingleTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]tag.KeyID{1}, map[tag.KeyID][]GroupingScanner{1: {scanner}})
	storageSeriesIDs := roaring.BitmapOf(1, 2, 3, 10)
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(storageSeriesIDs.GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	querySeriesIDs := roaring.BitmapOf(1, 2, 6, 10)
	lowSeriesContainer := querySeriesIDs.GetContainerAtIndex(0)
	// found series id 1,2,10, tag value id: 10,20
	dataLoadCtx := &DataLoadContext{
		SeriesIDHighKey:       1,
		LowSeriesIDsContainer: lowSeriesContainer,
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				DownSamplingSpecs:   aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
				GroupingTagValueIDs: make([]*roaring.Bitmap, 1),
				Query:               &stmt.Query{GroupBy: []string{"a"}},
			},
		},
		IsGrouping: true,
	}
	dataLoadCtx.Grouping()
	ctx.BuildGroup(dataLoadCtx)
	rs := dataLoadCtx.GroupingSeriesAgg
	assert.Len(t, rs, 2)
	// found series id: 1,2,10, tag value id: 10,20,10, index of refs: 0,1,0
	refMap := make(map[uint16]int)

	dataLoadCtx.IterateLowSeriesIDs(roaring.FastAnd(querySeriesIDs, storageSeriesIDs).GetContainer(0),
		func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
			refIdx := dataLoadCtx.GroupingSeriesAggRefs[seriesIdxFromQuery]
			refMap[refIdx]++
		})
	assert.Len(t, refMap, 2)
	assert.Equal(t, 2, refMap[0])
	assert.Equal(t, 1, refMap[1])
	assert.Equal(t, uint16(0), dataLoadCtx.GroupingSeriesAggRefs[1-dataLoadCtx.MinSeriesID])
	assert.Equal(t, uint16(1), dataLoadCtx.GroupingSeriesAggRefs[2-dataLoadCtx.MinSeriesID])
	assert.Equal(t, uint16(0), dataLoadCtx.GroupingSeriesAggRefs[10-dataLoadCtx.MinSeriesID])

	// high key not found
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).Return(nil, nil)
	dataLoadCtx.GroupingSeriesAgg = nil
	ctx.BuildGroup(dataLoadCtx)
	rs = dataLoadCtx.GroupingSeriesAgg
	assert.Empty(t, rs)
}

func TestGroupingContext_BuildMultiTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]tag.KeyID{1, 2}, map[tag.KeyID][]GroupingScanner{1: {scanner}, 2: {scanner}})
	storageSeriesIDs := roaring.BitmapOf(1, 2, 3, 10)
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(storageSeriesIDs.GetContainerAtIndex(0), []uint32{10, 20, 30, 10}).MaxTimes(2)
	querySeriesIDs := roaring.BitmapOf(1, 2, 6, 10)
	lowSeriesContainer := querySeriesIDs.GetContainerAtIndex(0)
	// found series id 1,2,10, tag value id: 10,20
	dataLoadCtx := &DataLoadContext{
		SeriesIDHighKey:       1,
		LowSeriesIDsContainer: lowSeriesContainer,
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				DownSamplingSpecs:   aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
				GroupingTagValueIDs: make([]*roaring.Bitmap, 2),
				Query:               &stmt.Query{GroupBy: []string{"a", "b"}},
			},
		},
		IsGrouping: true,
	}
	dataLoadCtx.Grouping()
	ctx.BuildGroup(dataLoadCtx)
	rs := dataLoadCtx.GroupingSeriesAgg
	assert.Len(t, rs, 2)
	// found series id: 1,2,10, tag value id: 10,20,10, index of refs: 0,1,0
	refMap := make(map[uint16]int)

	dataLoadCtx.IterateLowSeriesIDs(roaring.FastAnd(querySeriesIDs, storageSeriesIDs).GetContainer(0),
		func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
			refIdx := dataLoadCtx.GroupingSeriesAggRefs[seriesIdxFromQuery]
			refMap[refIdx]++
		})
	assert.Len(t, refMap, 2)
	assert.Equal(t, 2, refMap[0])
	assert.Equal(t, 1, refMap[1])
	assert.Equal(t, uint16(0), dataLoadCtx.GroupingSeriesAggRefs[1-dataLoadCtx.MinSeriesID])
	assert.Equal(t, uint16(1), dataLoadCtx.GroupingSeriesAggRefs[2-dataLoadCtx.MinSeriesID])
	assert.Equal(t, uint16(0), dataLoadCtx.GroupingSeriesAggRefs[10-dataLoadCtx.MinSeriesID])
}

func TestGroupingContext_ScanTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]tag.KeyID{1}, map[tag.KeyID][]GroupingScanner{1: {scanner}})
	// case 1: get tag value ids
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(roaring.BitmapOf(1, 2, 3, 10).GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	result := ctx.ScanTagValueIDs(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Equal(t, []uint32{10, 20}, result[0].ToArray())
	// case 2: empty tag value
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).Return(nil, nil)
	result = ctx.ScanTagValueIDs(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Equal(t, roaring.New(), result[0])
}

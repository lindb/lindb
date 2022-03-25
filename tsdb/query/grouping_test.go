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

package query

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func TestGroupingContext_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := flow.NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]tag.KeyID{1}, map[tag.KeyID][]flow.GroupingScanner{1: {scanner}})
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(roaring.BitmapOf(1, 2, 3, 10).GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	lowSeriesContainer := roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0)
	// found series id 1,2,10, tag value id: 10,20
	dataLoadCtx := &flow.DataLoadContext{
		SeriesIDHighKey:       1,
		LowSeriesIDsContainer: lowSeriesContainer,
		ShardExecuteCtx: &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				DownSamplingSpecs:   aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
				GroupingTagValueIDs: make([]*roaring.Bitmap, 1),
			},
		},
	}
	dataLoadCtx.Grouping()
	ctx.BuildGroup(dataLoadCtx)
	rs := dataLoadCtx.GroupingSeriesAgg
	assert.Len(t, rs, 2)

	// high key not found
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).Return(nil, nil)
	dataLoadCtx.GroupingSeriesAgg = nil
	ctx.BuildGroup(dataLoadCtx)
	rs = dataLoadCtx.GroupingSeriesAgg
	assert.Empty(t, rs)
}

func TestGroupingContext_ScanTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := flow.NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]tag.KeyID{1}, map[tag.KeyID][]flow.GroupingScanner{1: {scanner}})
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

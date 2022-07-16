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

package operator

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestDataLoader_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rs := flow.NewMockFilterResultSet(ctrl)
	ctx := &flow.DataLoadContext{
		PendingDataLoadTasks: atomic.NewInt32(0),
		ShardExecuteCtx: &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				Query: &stmt.Query{
					Interval:      1,
					IntervalRatio: 1.0,
				},
				DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(1, 2, 3),
		},
	}
	ctx.PrepareAggregatorWithoutGrouping()
	agg := aggregation.NewMockSeriesAggregator(ctrl)
	ctx.WithoutGroupingSeriesAgg.Aggregator = agg
	segment := &flow.TimeSegmentResultSet{FilterRS: []flow.FilterResultSet{rs}}

	t.Run("series not in rs", func(t *testing.T) {
		op := NewDataLoad(ctx, segment, rs)
		rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(4, 5))
		assert.NoError(t, op.Execute())
	})
	t.Run("cannot found loader", func(t *testing.T) {
		op := NewDataLoad(ctx, segment, rs)
		rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2, 4, 5))
		rs.EXPECT().Load(gomock.Any()).Return(nil)
		assert.NoError(t, op.Execute())
	})
	t.Run("load data", func(t *testing.T) {
		segment.IntervalRatio = 1
		loader := flow.NewMockDataLoader(ctrl)
		rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2))
		rs.EXPECT().Load(gomock.Any()).Return(loader)
		agg.EXPECT().GetAggregator(gomock.Any()).Return(nil, false)
		fAgg := aggregation.NewMockFieldAggregator(ctrl)
		agg.EXPECT().GetAggregator(gomock.Any()).Return(fAgg, true)
		getter := encoding.NewMockTSDValueGetter(ctrl)
		getter.EXPECT().GetValue(gomock.Any()).Return(5.0, true).AnyTimes()
		loader.EXPECT().Load(gomock.Any()).Do(func(ctx *flow.DataLoadContext) {
			ctx.DownSampling(timeutil.SlotRange{Start: 5, End: 5}, 0, 0, getter)
			ctx.DownSampling(timeutil.SlotRange{Start: 5, End: 5}, 0, 0, getter)
		})
		op := NewDataLoad(ctx, segment, rs)
		assert.NoError(t, op.Execute())
	})
}

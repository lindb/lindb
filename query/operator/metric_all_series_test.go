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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestMetricAllSeries_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	indexSegment := index.NewMockMetricIndexSegment(ctrl)
	shard.EXPECT().IndexSegment().Return(indexSegment).AnyTimes()

	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{
				Interval:      1,
				IntervalRatio: 1.0,
			},
			DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
		},
	}

	indexSegment.EXPECT().GetGroupingContext(gomock.Any()).Return(nil).AnyTimes()
	indexSegment.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	indexSegment.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(nil, constants.ErrNotFound)
	indexSegment.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(roaring.BitmapOf(3, 5), nil)

	t.Run("get series ids failure", func(t *testing.T) {
		op := NewMetricAllSeries(ctx, shard)
		indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any()).Return(nil, fmt.Errorf("err")).AnyTimes()
		assert.Error(t, op.Execute())
	})
	t.Run("series ids not found", func(t *testing.T) {
		ctx.SeriesIDsAfterFiltering = roaring.New()
		op := NewMetricAllSeries(ctx, shard)
		indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any()).Return(nil, constants.ErrNotFound).AnyTimes()
		assert.Error(t, op.Execute())
		assert.Equal(t, roaring.New(), ctx.SeriesIDsAfterFiltering)
	})
	t.Run("found series ids", func(t *testing.T) {
		ctx.SeriesIDsAfterFiltering = roaring.New()
		op := NewMetricAllSeries(ctx, shard)
		indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any()).Return(roaring.BitmapOf(3, 5), nil).AnyTimes()
		assert.NoError(t, op.Execute())
		assert.Equal(t, roaring.BitmapOf(0, 3, 5), ctx.SeriesIDsAfterFiltering)
	})
}

func TestMetricAllSeries_Stats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexSegment := index.NewMockMetricIndexSegment(ctrl)
	shard.EXPECT().IndexSegment().Return(indexSegment).AnyTimes()

	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{
				Interval:      1,
				IntervalRatio: 1.0,
			},
			DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
		},
		SeriesIDsAfterFiltering: roaring.BitmapOf(1, 2),
	}

	op := NewMetricAllSeries(ctx, shard)
	assert.Equal(t, "All Series", op.Identifier())

	op1 := op.(TrackableOperator)
	assert.NotNil(t, op1.Stats())
}

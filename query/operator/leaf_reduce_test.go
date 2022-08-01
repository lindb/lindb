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
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestLeafReducer_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataLoadCtx := &flow.DataLoadContext{
		PendingDataLoadTasks: atomic.NewInt32(0),
	}
	t.Run("pending data load task", func(t *testing.T) {
		dataLoadCtx.PendingDataLoadTasks.Inc()
		op := NewLeafReduce(nil, dataLoadCtx)
		assert.NoError(t, op.Execute())
	})
	t.Run("all data load task completed", func(t *testing.T) {
		agg := aggregation.NewMockSeriesAggregator(ctrl)
		it := series.NewMockIterator(ctrl)
		it.EXPECT().FieldName().Return(field.Name("f"))
		agg.EXPECT().Reset()
		agg.EXPECT().ResultSet().Return(it)
		dataLoadCtx.ShardExecuteCtx = &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				Query: &stmt.Query{
					Interval:      1,
					IntervalRatio: 1.0,
				},
				Stats:             &models.LeafNodeStats{},
				DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
			},
		}
		dataLoadCtx.WithoutGroupingSeriesAgg = &flow.GroupingSeriesAgg{
			Aggregator: agg,
		}
		dataLoadCtx.PendingDataLoadTasks.Store(0)
		op := NewLeafReduce(&context.LeafExecuteContext{
			ReduceCtx: context.NewLeafReduceContext(dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx, nil),
		}, dataLoadCtx)
		assert.NoError(t, op.Execute())
	})
}

func TestLeafReducer_Identifier(t *testing.T) {
	assert.Equal(t, "Reduce", NewLeafReduce(nil, nil).Identifier())
}

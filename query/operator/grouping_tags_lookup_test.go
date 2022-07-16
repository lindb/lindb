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

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestGroupingTagsLookup_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	seriesIDs := roaring.BitmapOf(1, 2, 3)
	ctx := &flow.ShardExecuteContext{
		SeriesIDsAfterFiltering: seriesIDs,
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query:             &stmt.Query{},
			Stats:             models.NewStorageStats(),
			DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
		},
	}
	dataLoadCtx := &flow.DataLoadContext{
		ShardExecuteCtx:       ctx,
		LowSeriesIDsContainer: seriesIDs.GetContainerAtIndex(0),
		IsGrouping:            false,
	}
	t.Run("no group", func(t *testing.T) {
		op := NewGroupingTagsLookup(dataLoadCtx)
		assert.NoError(t, op.Execute())
	})
	t.Run("has grouping", func(t *testing.T) {
		groupingCtx := flow.NewMockGroupingContext(ctrl)
		groupingCtx.EXPECT().BuildGroup(gomock.Any()).AnyTimes()
		ctx.GroupingContext = groupingCtx
		op := NewGroupingTagsLookup(dataLoadCtx)
		assert.NoError(t, op.Execute())
	})
}

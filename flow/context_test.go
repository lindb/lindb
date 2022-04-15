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
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

func TestStorageExecuteContext_collectGroupingTagValueIDs(t *testing.T) {
	ctx := &StorageExecuteContext{
		GroupingTagValueIDs: make([]*roaring.Bitmap, 2),
	}
	ctx.collectGroupingTagValueIDs([]uint32{1, 4})
	ctx.collectGroupingTagValueIDs([]uint32{2, 5})
	ctx.collectGroupingTagValueIDs([]uint32{8, 10})
	assert.Equal(t, roaring.BitmapOf(1, 2, 8), ctx.GroupingTagValueIDs[0])
	assert.Equal(t, roaring.BitmapOf(4, 5, 10), ctx.GroupingTagValueIDs[1])
}

func TestStorageExecuteContext(t *testing.T) {
	assert.True(t, (&StorageExecuteContext{Query: &stmt.Query{Condition: &stmt.FieldExpr{}}}).HasWhereCondition())
	assert.False(t, (&StorageExecuteContext{Query: &stmt.Query{}}).HasWhereCondition())
}

func TestStorageExecuteContext_HasGroupingTagValueIDs(t *testing.T) {
	ctx := &StorageExecuteContext{
		GroupingTagValueIDs: make([]*roaring.Bitmap, 2),
	}
	assert.False(t, ctx.HasGroupingTagValueIDs())
	ctx = &StorageExecuteContext{
		GroupingTagValueIDs: []*roaring.Bitmap{nil, roaring.BitmapOf(1), nil},
	}
	assert.True(t, ctx.HasGroupingTagValueIDs())
}

func TestStorageExecuteContext_SortFields(t *testing.T) {
	ctx := &StorageExecuteContext{
		Fields: field.Metas{{ID: 4}, {ID: 1}, {ID: 3}},
	}
	ctx.SortFields()
	assert.Equal(t, field.Metas{{ID: 1}, {ID: 3}, {ID: 4}}, ctx.Fields)
}

func TestStorageExecuteContext_QueryStats(t *testing.T) {
	assert.Nil(t, (&StorageExecuteContext{}).QueryStats())
	assert.NotNil(t, (&StorageExecuteContext{Stats: models.NewStorageStats()}).QueryStats())
}

func TestStorageExecuteContext_Release(t *testing.T) {
	ctx := &StorageExecuteContext{
		TaskCtx: NewTaskContextWithTimeout(context.TODO(), time.Second),
	}
	ctx.Release()
	ctx.ShardContexts = make([]*ShardExecuteContext, 2)
	ctx.Release()
	ctx.ShardContexts[0] = &ShardExecuteContext{}
	ctx.Release()
}

func TestTimeSegmentContext_AddFilterResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := NewTimeSegmentContext()
	rs1 := NewMockFilterResultSet(ctrl)
	rs1.EXPECT().FamilyTime().Return(int64(20)).AnyTimes()
	rs1.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	rs1.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1)).AnyTimes()
	rs2 := NewMockFilterResultSet(ctrl)
	rs2.EXPECT().FamilyTime().Return(int64(10)).AnyTimes()
	rs2.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	rs2.EXPECT().SeriesIDs().Return(roaring.BitmapOf(2)).AnyTimes()
	rs3 := NewMockFilterResultSet(ctrl)
	rs3.EXPECT().FamilyTime().Return(int64(20)).AnyTimes()
	rs3.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	rs3.EXPECT().SeriesIDs().Return(roaring.BitmapOf(3)).AnyTimes()

	ctx.AddFilterResultSet(timeutil.Interval(10), rs1)
	ctx.AddFilterResultSet(timeutil.Interval(10), rs2)
	ctx.AddFilterResultSet(timeutil.Interval(10), rs3)

	assert.Equal(t, roaring.BitmapOf(1, 2, 3), ctx.SeriesIDs)
	segments := ctx.GetTimeSegments()
	assert.Len(t, segments, 2)
	assert.Equal(t, int64(10), segments[0].FamilyTime)
	assert.Equal(t, int64(20), segments[1].FamilyTime)

	rs1.EXPECT().Close()
	rs2.EXPECT().Close()
	rs3.EXPECT().Close()
	ctx.Release()
}

func TestShardExecuteContext(t *testing.T) {
	ctx := NewShardExecuteContext(&StorageExecuteContext{})
	assert.True(t, ctx.IsSeriesIDsEmpty())
	ctx.TimeSegmentContext = &TimeSegmentContext{SeriesIDs: roaring.BitmapOf(1, 2, 3)}
	assert.False(t, ctx.IsSeriesIDsEmpty())
	ctx.Release()
}

func TestGroupingSeriesAgg_reduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
	agg := &GroupingSeriesAgg{Aggregator: seriesAgg}
	c := 0
	seriesAgg.EXPECT().Reset()
	agg.reduce(func(it series.GroupedIterator) {
		c++
	})
	assert.Equal(t, 1, c)

	c = 0
	seriesAgg.EXPECT().Reset()
	agg = &GroupingSeriesAgg{Aggregators: aggregation.FieldAggregates{seriesAgg}}
	agg.reduce(func(it series.GroupedIterator) {
		c++
	})
	assert.Equal(t, 1, c)
}

func TestDataLoadContext_PrepareAggregatorWithoutGrouping(t *testing.T) {
	ctx := &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Fields:            field.Metas{{ID: 1}},
				DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
				Query:             &stmt.Query{},
			},
		},
	}
	ctx.PrepareAggregatorWithoutGrouping()
	assert.NotNil(t, ctx.WithoutGroupingSeriesAgg.Aggregator)
	assert.Nil(t, ctx.WithoutGroupingSeriesAgg.Aggregators)

	ctx = &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Fields: field.Metas{{ID: 1}, {ID: 2}},
				DownSamplingSpecs: aggregation.AggregatorSpecs{
					aggregation.NewAggregatorSpec("a", field.SumField),
					aggregation.NewAggregatorSpec("b", field.SumField),
				},
				Query: &stmt.Query{},
			},
		},
	}
	ctx.PrepareAggregatorWithoutGrouping()
	assert.Nil(t, ctx.WithoutGroupingSeriesAgg.Aggregator)
	assert.NotNil(t, ctx.WithoutGroupingSeriesAgg.Aggregators)
}

func TestDataLoadContext_NewSeriesAggregator(t *testing.T) {
	ctx := &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Fields:              field.Metas{{ID: 1}},
				DownSamplingSpecs:   aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
				Query:               &stmt.Query{},
				GroupByTagKeyIDs:    []tag.KeyID{1},
				GroupingTagValueIDs: make([]*roaring.Bitmap, 1),
			},
		},
	}
	idx := ctx.NewSeriesAggregator(string([]byte{1, 0, 0, 0}))
	assert.Equal(t, uint16(0), idx)
	idx = ctx.NewSeriesAggregator(string([]byte{2, 0, 0, 0}))
	assert.Equal(t, uint16(1), idx)
	assert.NotNil(t, ctx.GroupingSeriesAgg[0].Aggregator)
	assert.Nil(t, ctx.GroupingSeriesAgg[0].Aggregators)

	ctx = &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Fields: field.Metas{{ID: 1}, {ID: 2}},
				DownSamplingSpecs: aggregation.AggregatorSpecs{
					aggregation.NewAggregatorSpec("a", field.SumField),
					aggregation.NewAggregatorSpec("b", field.SumField),
				},
				Query: &stmt.Query{},
			},
		},
	}
	idx = ctx.NewSeriesAggregator("")
	assert.Equal(t, uint16(0), idx)
	assert.Nil(t, ctx.GroupingSeriesAgg[0].Aggregator)
	assert.NotNil(t, ctx.GroupingSeriesAgg[0].Aggregators)
}

func TestDataLoadContext_HasGroupingData(t *testing.T) {
	ctx := &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query: &stmt.Query{},
			},
		},
	}
	assert.True(t, ctx.HasGroupingData())
	ctx = &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query: &stmt.Query{GroupBy: []string{"ip"}},
			},
		},
	}
	assert.False(t, ctx.HasGroupingData())
	ctx.GroupingSeriesAgg = []*GroupingSeriesAgg{nil}
	assert.True(t, ctx.HasGroupingData())
}

func TestDataLoadContext_IterateLowSeriesIDs(t *testing.T) {
	querySeriesIDs := roaring.BitmapOf(5, 11, 13)
	storageSeriesIDs := roaring.BitmapOf(1, 3, 5, 7, 9, 11, 13, 15)
	ctx := &DataLoadContext{
		LowSeriesIDsContainer: querySeriesIDs.GetContainer(0),
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query: &stmt.Query{GroupBy: []string{"ip"}},
			},
		},
	}
	ctx.Grouping()
	findSeriesIDs := roaring.New()
	storageLowSeriesContainer := storageSeriesIDs.GetContainer(0)
	storageLowSeriesIDs := storageLowSeriesContainer.ToArray()
	ctx.IterateLowSeriesIDs(storageLowSeriesContainer, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
		findSeriesIDs.Add(uint32(storageLowSeriesIDs[seriesIdxFromStorage]))
	})
	assert.Equal(t, querySeriesIDs, findSeriesIDs)
}

func TestDataLoadContext_GetSeriesAggregator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agg := aggregation.NewMockSeriesAggregator(ctrl)
	agg.EXPECT().Reset().AnyTimes()
	ctx := &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query: &stmt.Query{GroupBy: []string{"ip"}},
			},
		},
		GroupingSeriesAggRefs: []uint16{1},
		GroupingSeriesAgg:     []*GroupingSeriesAgg{{Aggregator: agg}, {Aggregator: agg}},
	}
	aggregator := ctx.GetSeriesAggregator(0, 0)
	assert.NotNil(t, aggregator)
	ctx.Reduce(func(it series.GroupedIterator) {})

	ctx = &DataLoadContext{
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query:  &stmt.Query{},
				Fields: field.Metas{{ID: 1}, {ID: 2}},
			},
		},
		WithoutGroupingSeriesAgg: &GroupingSeriesAgg{Aggregators: aggregation.FieldAggregates{agg, agg}},
	}
	aggregator = ctx.GetSeriesAggregator(0, 1)
	assert.NotNil(t, aggregator)
	ctx.Reduce(func(it series.GroupedIterator) {})
}

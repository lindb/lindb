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
//
package storagequery

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
)

func TestBaseQueryTask_Run(t *testing.T) {
	task := &baseQueryTask{}
	assert.Nil(t, task.Run())
}

func TestQueryStatTask_Run(t *testing.T) {
	task := &queryStatTask{}
	task.BeforeRun()
	task.AfterRun()
}

func TestSeriesIDsSearchTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newSeriesSearchFunc = newSeriesSearch
		ctrl.Finish()
	}()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	ctx := &flow.ShardExecuteContext{
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{},
			Stats: models.NewStorageStats(),
		},
		SeriesIDsAfterFiltering: roaring.New(),
	}
	task := newSeriesIDsSearchTask(ctx, shard)
	// case 1: search err
	indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: no group by add series ids without tags
	indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(roaring.New(), nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(series.IDWithoutTags), ctx.SeriesIDsAfterFiltering)
	ctx.SeriesIDsAfterFiltering.Clear()
	// case 3: group by tag
	indexDB.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).Return(roaring.New(), nil)
	ctx.StorageExecuteCtx.Query = &stmt.Query{GroupBy: []string{"host"}}
	task = newSeriesIDsSearchTask(ctx, shard)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), ctx.SeriesIDsAfterFiltering.GetCardinality())
	// case 4: has condition, return err
	q, _ := sql.Parse("select f from cpu where ip<>'1.1.1.1'")
	query := q.(*stmt.Query)
	seriesSearch := NewMockSeriesSearch(ctrl)
	newSeriesSearchFunc = func(filter series.Filter, filterResult map[string]*flow.TagFilterResult, condition stmt.Expr) SeriesSearch {
		return seriesSearch
	}
	seriesSearch.EXPECT().Search().Return(nil, fmt.Errorf("err"))
	ctx.StorageExecuteCtx.Query = query
	task = newSeriesIDsSearchTask(ctx, shard)
	err = task.Run()
	assert.Error(t, err)
	// case 5: has condition, return series ids
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), ctx.SeriesIDsAfterFiltering)
	ctx.SeriesIDsAfterFiltering.Clear()
	// case 6: explain
	q, _ = sql.Parse("explain select f from cpu where ip<>'1.1.1.1'")
	query = q.(*stmt.Query)
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil)
	shard.EXPECT().ShardID().Return(models.ShardID(10))
	ctx.StorageExecuteCtx.Query = query
	task = newSeriesIDsSearchTask(ctx, shard)
	task = newSeriesIDsSearchTask(ctx, shard)
	err = task.Run()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), ctx.SeriesIDsAfterFiltering)
}

func TestFamilyFilterTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resultSet := flow.NewMockFilterResultSet(ctrl)
	resultSet.EXPECT().Identifier().Return("memory").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	ctx := &flow.ShardExecuteContext{
		SeriesIDsAfterFiltering: seriesIDs,
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{},
			Stats: models.NewStorageStats(),
		},
		TimeSegmentContext: flow.NewTimeSegmentContext(),
	}
	task := newFamilyFilterTask(ctx, shard)
	// case 1: get empty family
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil)
	err := task.Run()
	assert.NoError(t, err)
	assert.Empty(t, ctx.TimeSegmentContext.TimeSegments)
	// case 2: family filter err
	family := tsdb.NewMockDataFamily(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family}).AnyTimes()
	family.EXPECT().Filter(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = task.Run()
	assert.Error(t, err)
	// case 3: family filter not found err
	family.EXPECT().Filter(gomock.Any()).Return(nil, constants.ErrNotFound)
	err = task.Run()
	assert.NoError(t, err)
	// case 4: get data
	family.EXPECT().Interval().Return(timeutil.Interval(10000))
	resultSet.EXPECT().FamilyTime().Return(int64(10))
	resultSet.EXPECT().SeriesIDs().Return(roaring.New())
	resultSet.EXPECT().SlotRange().Return(timeutil.SlotRange{}).MaxTimes(3)
	family.EXPECT().Filter(gomock.Any()).Return([]flow.FilterResultSet{resultSet}, nil)
	err = task.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, ctx.TimeSegmentContext.TimeSegments)
	// case 5: explain
	family.EXPECT().Interval().Return(timeutil.Interval(10000))
	resultSet.EXPECT().FamilyTime().Return(int64(10))
	resultSet.EXPECT().SeriesIDs().Return(roaring.New())
	resultSet.EXPECT().FamilyTime().Return(int64(10)).MaxTimes(2)
	ctx.StorageExecuteCtx.Query.Explain = true
	task = newFamilyFilterTask(ctx, shard)
	family.EXPECT().Filter(gomock.Any()).Return([]flow.FilterResultSet{resultSet}, nil)
	shard.EXPECT().ShardID().Return(models.ShardID(10))
	err = task.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, ctx.TimeSegmentContext.TimeSegments)
}

func TestGroupingContextFindTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	seriesIDs := roaring.BitmapOf(1, 2, 3)
	ctx := &flow.ShardExecuteContext{
		SeriesIDsAfterFiltering: seriesIDs,
		StorageExecuteCtx: &flow.StorageExecuteContext{
			Query: &stmt.Query{},
			Stats: models.NewStorageStats(),
		},
	}
	task := newGroupingContextFindTask(ctx, shard)
	// case 1: get grouping context err
	indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(fmt.Errorf("err"))
	err := task.Run()
	assert.Error(t, err)
	// case 2: get grouping context
	indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil)
	err = task.Run()
	assert.NoError(t, err)
	// case 3: explain
	indexDB.EXPECT().GetGroupingContext(gomock.Any()).Return(nil)
	ctx.StorageExecuteCtx.Query.Explain = true
	task = newGroupingContextFindTask(ctx, shard)
	shard.EXPECT().ShardID().Return(models.ShardID(10))
	err = task.Run()
	assert.NoError(t, err)
}

func TestBuildGroupTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
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
	task := newBuildGroupTask(shard, dataLoadCtx)
	// case 1: no group
	err := task.Run()
	assert.NoError(t, err)
	// case 2: has grouping
	groupingCtx := flow.NewMockGroupingContext(ctrl)
	groupingCtx.EXPECT().BuildGroup(gomock.Any()).AnyTimes()
	ctx.GroupingContext = groupingCtx
	task = newBuildGroupTask(shard, dataLoadCtx)
	err = task.Run()
	assert.NoError(t, err)
	// case 3: explain
	ctx.StorageExecuteCtx.Query.Explain = true
	task = newBuildGroupTask(shard, dataLoadCtx)
	shard.EXPECT().ShardID().Return(models.ShardID(10))
	err = task.Run()
	assert.NoError(t, err)
}

func TestDataLoadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := tsdb.NewMockShard(ctrl)
	qf := flow.NewMockStorageQueryFlow(ctrl)
	qf.EXPECT().Reduce(gomock.Any())
	rs := flow.NewMockFilterResultSet(ctrl)
	rs.EXPECT().Identifier().Return(fmt.Sprintf("shard/%d/segment/day/20210703", 10))
	rs.EXPECT().Identifier().Return("memory").AnyTimes()
	ctx := &flow.DataLoadContext{
		ShardExecuteCtx: &flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				Query: &stmt.Query{
					Interval:      1,
					IntervalRatio: 1.0,
				},
				Stats:             models.NewStorageStats(),
				DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(1, 2, 3),
		},
	}
	ctx.PrepareAggregatorWithoutGrouping()
	agg := aggregation.NewMockSeriesAggregator(ctrl)
	ctx.WithoutGroupingSeriesAgg.Aggregator = agg
	segment := &flow.TimeSegmentResultSet{FilterRS: []flow.FilterResultSet{rs}}
	task := newDataLoadTask(shard, qf, ctx, 0, segment)
	loader := flow.NewMockDataLoader(ctrl)
	rs.EXPECT().Load(gomock.Any()).Return(nil)
	rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(10, 20))
	// case 1: load data
	err := task.Run()
	assert.NoError(t, err)
	// case 2: explain
	ctx.ShardExecuteCtx.StorageExecuteCtx.Query.Explain = true
	task = newDataLoadTask(shard, qf, ctx, 0, segment)
	shard.EXPECT().ShardID().Return(models.ShardID(10)).AnyTimes()
	rs.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2)).AnyTimes()
	rs.EXPECT().Load(gomock.Any()).Return(loader).AnyTimes()
	agg.EXPECT().GetAggregator(gomock.Any()).Return(nil, false)
	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	agg.EXPECT().Reset()
	agg.EXPECT().GetAggregator(gomock.Any()).Return(fAgg, true)
	getter := encoding.NewMockTSDValueGetter(ctrl)
	getter.EXPECT().GetValue(gomock.Any()).Return(5.0, true).AnyTimes()
	loader.EXPECT().Load(gomock.Any()).Do(func(ctx *flow.DataLoadContext) {
		ctx.DownSampling(timeutil.SlotRange{Start: 5, End: 5}, 0, 0, getter)
		ctx.DownSampling(timeutil.SlotRange{Start: 5, End: 5}, 0, 0, getter)
	})
	err = task.Run()
	assert.NoError(t, err)
	err = task.Run()
	assert.NoError(t, err)
}

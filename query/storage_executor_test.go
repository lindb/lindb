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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

type mockQueryFlow struct {
}

func (m *mockQueryFlow) ReduceTagValues(_ int, _ map[uint32]string) {
}

func (m *mockQueryFlow) Prepare(_ timeutil.Interval, _ int, _ timeutil.TimeRange, _ aggregation.AggregatorSpecs) {
}

func (m *mockQueryFlow) Filtering(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) Grouping(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) Load(task concurrent.Task) {
	task()
}

func (m *mockQueryFlow) Reduce(_ string, _ series.GroupedIterator) {
}

func (m *mockQueryFlow) Complete(_ error) {
}

func newMockQueryFlow() flow.StorageQueryFlow {
	return &mockQueryFlow{}
}

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	exeCtx := parallel.NewMockExecuteContext(ctrl)
	exeCtx.EXPECT().Complete(gomock.Any()).AnyTimes()

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().GetOption().Return(option.DatabaseOption{Interval: "10s"}).AnyTimes()
	mockDatabase.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.Interval(timeutil.OneSecond)}

	// case 1: query shards is empty
	exec := newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext(nil, query))
	queryFlow.EXPECT().Complete(errNoShardID)
	exec.Execute()

	// case 2: shards of engine is empty
	mockDatabase.EXPECT().NumOfShards().Return(0)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(errNoShardInDatabase)
	exec.Execute()

	// case 3: num. of shard not match
	mockDatabase.EXPECT().NumOfShards().Return(2)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(errShardNotMatch)
	exec.Execute()

	// case 4: shard not found
	mockDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, false).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(errShardNotFound)
	exec.Execute()
	// case 4: shard not match
	mockDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, false)
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, true).MaxTimes(2)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(errShardNumNotMatch)
	exec.Execute()

	// case 6: normal case
	q, _ := sql.Parse("select f from cpu")
	query = q.(*stmt.Query)
	mockDB1 := newMockDatabase(ctrl)
	mockDB1.EXPECT().GetOption().Return(option.DatabaseOption{Interval: "10s"})
	exec = newStorageExecutor(queryFlow, mockDB1, newStorageExecuteContext([]int32{1, 2, 3}, query))
	gomock.InOrder(
		queryFlow.EXPECT().Prepare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		queryFlow.EXPECT().Filtering(gomock.Any()).MaxTimes(3*2), //memory db and shard
	)
	exec.Execute()
}

func TestStorageExecute_Plan_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStorageExecutePlanFunc = newStorageExecutePlan
		ctrl.Finish()
	}()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)

	mockDatabase := newMockDatabase(ctrl)
	plan := NewMockPlan(ctrl)
	newStorageExecutePlanFunc = func(namespace string, metadata metadb.Metadata, query *stmt.Query) Plan {
		return plan
	}
	plan.EXPECT().Plan().Return(fmt.Errorf("err"))

	// find metric name err
	q, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	query := q.(*stmt.Query)
	exec := newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(fmt.Errorf("err"))
	exec.Execute()
}

func TestStorageExecutor_TagSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		newTagSearchFunc = newTagSearch
	}()
	tagSearch := NewMockTagSearch(ctrl)
	newTagSearchFunc = func(namespace, metricName string, condition stmt.Expr, metadata metadb.Metadata) TagSearch {
		return tagSearch
	}
	mockDatabase := newMockDatabase(ctrl)
	qFlow := flow.NewMockStorageQueryFlow(ctrl)
	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)

	// case 1: tag search err
	exec := newStorageExecutor(qFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	tagSearch.EXPECT().Filter().Return(nil, fmt.Errorf("err"))
	qFlow.EXPECT().Complete(fmt.Errorf("err"))
	exec.Execute()
	// case 2: tag search not result
	exec = newStorageExecutor(qFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	tagSearch.EXPECT().Filter().Return(nil, nil)
	qFlow.EXPECT().Complete(constants.ErrNotFound)
	exec.Execute()
}

func TestStorageExecute_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newSeriesSearchFunc = newSeriesSearch
		newTagSearchFunc = newTagSearch
		ctrl.Finish()
	}()

	tagSearch := NewMockTagSearch(ctrl)
	newTagSearchFunc = func(namespace, metricName string, condition stmt.Expr, metadata metadb.Metadata) TagSearch {
		return tagSearch
	}
	tagSearch.EXPECT().Filter().Return(map[string]*tagFilterResult{
		"host": {tagValueIDs: roaring.BitmapOf(1, 2)},
	}, nil).AnyTimes()
	seriesSearch := NewMockSeriesSearch(ctrl)
	newSeriesSearchFunc = func(filter series.Filter, filterResult map[string]*tagFilterResult, condition stmt.Expr) SeriesSearch {
		return seriesSearch
	}
	queryFlow := newMockQueryFlow()

	metadata := metadb.NewMockMetadata(ctrl)
	metadataIndex := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataIndex).AnyTimes()
	metadataIndex.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), "host").Return(uint32(10), nil).AnyTimes()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().GetOption().Return(option.DatabaseOption{Interval: "10s"}).AnyTimes()

	index := indexdb.NewMockIndexDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().CurrentInterval().Return(timeutil.Interval(10000)).AnyTimes()
	shard.EXPECT().IndexDatabase().Return(index).AnyTimes()

	// mock data
	mockDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true).AnyTimes()
	mockDatabase.EXPECT().GetShard(int32(2)).Return(shard, true).AnyTimes()
	mockDatabase.EXPECT().GetShard(int32(3)).Return(shard, true).AnyTimes()
	mockDatabase.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataIndex.EXPECT().GetMetricID(gomock.Any(), "cpu").Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetField(gomock.Any(), gomock.Any(), field.Name("f")).
		Return(field.Meta{ID: 10, Type: field.SumField}, nil).AnyTimes()
	shard.EXPECT().IndexDatabase().Return(nil).AnyTimes()

	// case 1: series search err
	q, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	query := q.(*stmt.Query)

	seriesSearch.EXPECT().Search().Return(nil, fmt.Errorf("err")).Times(3)
	exec := newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	exec.Execute()
	// case 2: normal case without filter
	q, _ = sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	query = q.(*stmt.Query)

	index.EXPECT().GetSeriesIDsForMetric(gomock.Any(), gomock.Any()).DoAndReturn(func(a, b string) (*roaring.Bitmap, error) {
		return roaring.BitmapOf(1, 2, 3), nil
	}).AnyTimes()
	filterRS := flow.NewMockFilterResultSet(ctrl)
	filterRS.EXPECT().Identifier().Return("memory").AnyTimes()
	filterRS.EXPECT().FamilyTime().Return(int64(10)).AnyTimes()
	filterRS.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	filterRS.EXPECT().Load(gomock.Any(), gomock.Any()).MaxTimes(3)
	filterRS.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2, 3)).MaxTimes(3)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	exec.Execute()
	// case 3: normal case with filter
	q, _ = sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	query = q.(*stmt.Query)

	filterRS.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	filterRS.EXPECT().Load(gomock.Any(), gomock.Any()).MaxTimes(3)
	filterRS.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2, 3)).MaxTimes(3)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	exec.Execute()

	//case 4: filter data err
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, fmt.Errorf("err")).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	exec.Execute()

	// case 5: filter result is nil
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).MaxTimes(3)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	exec.Execute()

	// case 6: filter shard data err
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).MaxTimes(3)
	family := tsdb.NewMockDataFamily(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]tsdb.DataFamily{family}).MaxTimes(3)
	family.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("err")).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	exec.Execute()

	// case 7: group by
	q, _ = sql.Parse("select f from cpu where host='1.1.1.1' group by host")
	query = q.(*stmt.Query)

	filterRS.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	filterRS.EXPECT().Load(gomock.Any(), gomock.Any()).MaxTimes(3)
	filterRS.EXPECT().SeriesIDs().Return(roaring.BitmapOf(1, 2, 3)).MaxTimes(3)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	shard.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]flow.FilterResultSet{filterRS}, nil).MaxTimes(3)
	index.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")).MaxTimes(3)
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1, 2, 3}, query))
	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).Times(3)
	exec.Execute()
}

func TestStorageExecutor_Execute_GroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newBuildGroupTaskFunc = newBuildGroupTask
		newDataLoadTaskFunc = newDataLoadTask
		ctrl.Finish()
	}()
	queryFlow := newMockQueryFlow()
	metadata := metadb.NewMockMetadata(ctrl)
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	mockDatabase.EXPECT().Metadata().Return(metadata).AnyTimes()
	// case 1: normal case
	q, _ := sql.Parse("select f from cpu group by host")
	query := q.(*stmt.Query)

	exec := newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1}, query))
	exec1 := exec.(*storageExecutor)
	exec1.groupByTagKeyIDs = []tag.Meta{{ID: 1, Key: "host"}}
	exec1.tagValueIDs = make([]*roaring.Bitmap, len(exec1.groupByTagKeyIDs))
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()
	rs := flow.NewMockFilterResultSet(ctrl)
	rs.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	rs.EXPECT().SlotRange().Return(timeutil.SlotRange{}).AnyTimes()
	gCtx := series.NewMockGroupingContext(ctrl)
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(gCtx, nil)
	gCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(map[string][]uint16{"host": {1, 2, 3}})
	gCtx.EXPECT().GetGroupByTagValueIDs().Return([]*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)}).AnyTimes()
	tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	exec1.storageExecutePlan = &storageExecutePlan{groupByTags: []tag.Meta{{ID: 1, Key: "host"}}}
	exec1.executeGroupBy(shard, &timeSpanResultSet{}, roaring.BitmapOf(1, 2, 3))

	// case 2: get grouping context err
	gomock.InOrder(
		indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")),
	)
	exec1.executeGroupBy(shard, &timeSpanResultSet{}, roaring.BitmapOf(1, 2, 3))
	// case 3: get grouping context nil
	gomock.InOrder(
		indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, nil),
	)
	exec1.executeGroupBy(shard, &timeSpanResultSet{}, roaring.BitmapOf(1, 2, 3))

	// case 4: collect tag values err
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(gCtx, nil)
	gCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(map[string][]uint16{"host": {1, 2, 3}})
	tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	exec = newStorageExecutor(queryFlow, mockDatabase, newStorageExecuteContext([]int32{1}, query))
	exec1 = exec.(*storageExecutor)
	exec1.groupByTagKeyIDs = []tag.Meta{{ID: 1, Key: "host"}}
	exec1.tagValueIDs = make([]*roaring.Bitmap, len(exec1.groupByTagKeyIDs))
	exec1.storageExecutePlan = &storageExecutePlan{groupByTags: []tag.Meta{{ID: 1, Key: "host"}}}
	exec1.executeGroupBy(shard, &timeSpanResultSet{}, roaring.BitmapOf(1, 2, 3))

	// case 5: build group series err
	task := flow.NewMockQueryTask(ctrl)
	newBuildGroupTaskFunc = func(ctx *storageExecuteContext, shard tsdb.Shard, groupingCtx series.GroupingContext,
		highKey uint16, container roaring.Container, result *groupedSeriesResult) flow.QueryTask {
		return task
	}
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(gCtx, nil)
	task.EXPECT().Run().Return(fmt.Errorf("err"))
	exec1.executeGroupBy(shard, &timeSpanResultSet{}, roaring.BitmapOf(1, 2, 3))

	newBuildGroupTaskFunc = newBuildGroupTask
	// case 6: load data err
	newDataLoadTaskFunc = func(ctx *storageExecuteContext, shard tsdb.Shard, queryFlow flow.StorageQueryFlow,
		span *timeSpan,
		highKey uint16, seriesID roaring.Container,
	) flow.QueryTask {
		return task
	}
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(gCtx, nil)
	task.EXPECT().Run().Return(fmt.Errorf("err"))
	gCtx.EXPECT().BuildGroup(gomock.Any(), gomock.Any()).Return(map[string][]uint16{"host": {1, 2, 3}})
	exec1.executeGroupBy(shard, &timeSpanResultSet{spanMap: map[int64]*timeSpan{1: {}}, filterRSCount: 1}, roaring.BitmapOf(1, 2, 3))
}

func TestStorageExecutor_merge_groupBy_tagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	queryFlow.EXPECT().Load(gomock.Any()).AnyTimes()
	exec := newStorageExecutor(queryFlow, nil, newStorageExecuteContext([]int32{1}, &stmt.Query{}))
	exec1 := exec.(*storageExecutor)
	exec1.groupByTagKeyIDs = []tag.Meta{{ID: 1}, {ID: 2}, {ID: 3}}
	exec1.pendingForShard.Add(1)
	// case 1: has pending task return it
	exec1.mergeGroupByTagValueIDs(nil)
	// case 2: new tag values
	exec1.pendingForShard.Dec()
	exec1.mergeGroupByTagValueIDs([]*roaring.Bitmap{roaring.BitmapOf(1, 2, 3), nil, nil})
	// case 3: merge tag value
	exec1.mergeGroupByTagValueIDs([]*roaring.Bitmap{roaring.BitmapOf(4, 5, 6), roaring.BitmapOf(1, 2, 3), nil})
}

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

package storagequery

import (
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

type mockQueryFlow struct {
	err error
}

func (m *mockQueryFlow) Submit(_ flow.Stage, task concurrent.Task) {
	if m.err == nil {
		task()
	}
}

func (m *mockQueryFlow) ReduceTagValues(_ int, _ map[uint32]string) {
}

func (m *mockQueryFlow) Prepare() {
}

func (m *mockQueryFlow) Reduce(_ series.GroupedIterator) {
}

func (m *mockQueryFlow) Complete(err error) {
	m.err = err
}

func newMockQueryFlow() *mockQueryFlow {
	return &mockQueryFlow{}
}

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)

	mockDatabase := tsdb.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().GetOption().
		Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}).AnyTimes()
	mockDatabase.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.Interval(timeutil.OneSecond)}

	cases := []struct {
		name    string
		shards  []models.ShardID
		prepare func()
	}{
		{
			name:   "query shards is empty",
			shards: nil,
			prepare: func() {
				queryFlow.EXPECT().Complete(errNoShardID)
			},
		},
		{
			name:   "shards of engine is empty",
			shards: []models.ShardID{1, 2, 3},
			prepare: func() {
				mockDatabase.EXPECT().NumOfShards().Return(0)
				queryFlow.EXPECT().Complete(errNoShardInDatabase)
			},
		},
		{
			name:   "shard not found",
			shards: []models.ShardID{1, 2, 3},
			prepare: func() {
				mockDatabase.EXPECT().NumOfShards().Return(3)
				mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, false).MaxTimes(3)
				queryFlow.EXPECT().Complete(errShardNotFound)
			},
		},
		{
			name:   "shard not match",
			shards: []models.ShardID{1, 2, 3},
			prepare: func() {
				mockDatabase.EXPECT().NumOfShards().Return(3)
				mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, false)
				mockDatabase.EXPECT().GetShard(gomock.Any()).Return(nil, true).MaxTimes(2)
				queryFlow.EXPECT().Complete(errShardNumNotMatch)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			storageQuery := newStorageMetricQuery(queryFlow, newStorageExecuteContext(mockDatabase, tt.shards, query))
			if tt.prepare != nil {
				tt.prepare()
			}
			storageQuery.Execute()
		})
	}
}

func TestStorageExecute_Plan_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStorageExecutePlanFunc = newStorageExecutePlan
		newStoragePlanTaskFunc = newStoragePlanTask
		ctrl.Finish()
	}()

	queryFlow := flow.NewMockStorageQueryFlow(ctrl)

	mockDatabase := newMockDatabase(ctrl)
	newStorageExecutePlanFunc = func(ctx *executeContext) *storageExecutePlan {
		return &storageExecutePlan{}
	}
	task := flow.NewMockQueryTask(ctrl)
	newStoragePlanTaskFunc = func(_ *executeContext, _ *storageExecutePlan) flow.QueryTask {
		return task
	}

	// find metric name err
	q, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	query := q.(*stmt.Query)
	task.EXPECT().Run().Return(io.ErrClosedPipe)
	storageQuery := newStorageMetricQuery(queryFlow, newStorageExecuteContext(mockDatabase, []models.ShardID{1, 2, 3}, query))
	queryFlow.EXPECT().Complete(io.ErrClosedPipe)
	storageQuery.Execute()
}

func TestStorageExecutor_TagSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		newTagSearchFunc = newTagSearch
	}()
	tagSearch := NewMockTagSearch(ctrl)
	newTagSearchFunc = func(ctx *executeContext) TagSearch {
		return tagSearch
	}
	mockDatabase := newMockDatabase(ctrl)
	qFlow := flow.NewMockStorageQueryFlow(ctrl)
	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)

	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "tag search err",
			prepare: func() {
				tagSearch.EXPECT().Filter().Return(fmt.Errorf("err"))
				qFlow.EXPECT().Complete(fmt.Errorf("err"))
			},
		},
		{
			name: "tag search not result",
			prepare: func() {
				tagSearch.EXPECT().Filter().Return(nil)
				qFlow.EXPECT().Complete(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			exec := newStorageMetricQuery(qFlow, newStorageExecuteContext(mockDatabase, []models.ShardID{1, 2, 3}, query))
			if tt.prepare != nil {
				tt.prepare()
			}

			exec.Execute()
		})
	}
}

func TestStorageExecute_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagSearchFunc = newTagSearch
		newStoragePlanTaskFunc = newStoragePlanTask
		newTagFilterTaskFunc = newTagFilterTask
		ctrl.Finish()
	}()

	planTask := flow.NewMockQueryTask(ctrl)
	planTask.EXPECT().Run().Return(nil).AnyTimes()
	newStoragePlanTaskFunc = func(ctx *executeContext, plan *storageExecutePlan) flow.QueryTask {
		return planTask
	}
	tagFilterTask := flow.NewMockQueryTask(ctrl)
	tagFilterTask.EXPECT().Run().Return(nil).AnyTimes()
	newTagFilterTaskFunc = func(ctx *executeContext, tagSearch TagSearch) flow.QueryTask {
		return tagFilterTask
	}
	queryFlow := newMockQueryFlow()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	mockDatabase.EXPECT().GetShard(gomock.Any()).Return(shard, true).AnyTimes()
	mockDatabase.EXPECT().GetOption().
		Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}).AnyTimes()

	ctx := &executeContext{
		storageExecuteCtx: &flow.StorageExecuteContext{
			ShardIDs: []models.ShardID{1, 2, 3},
			TagFilterResult: map[string]*flow.TagFilterResult{
				"host": {TagValueIDs: roaring.BitmapOf(1, 2)},
			},
		},
		database: mockDatabase,
	}
	cases := []struct {
		name    string
		sql     string
		prepare func()
		assert  func()
	}{
		{
			name: "series search err",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'",
			prepare: func() {
				seriesSearchTask := flow.NewMockQueryTask(ctrl)
				seriesSearchTask.EXPECT().Run().Return(constants.ErrNotFound)
				seriesSearchTask.EXPECT().Run().Return(io.ErrClosedPipe)
				newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					return seriesSearchTask
				}
			},
			assert: func() {
				assert.Equal(t, io.ErrClosedPipe, queryFlow.err)
			},
		},
		{
			name: "series id not found",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'",
			prepare: func() {
				seriesSearchTask := flow.NewMockQueryTask(ctrl)
				seriesSearchTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					return seriesSearchTask
				}
			},
			assert: func() {
				assert.NoError(t, queryFlow.err)
			},
		},
		{
			name: "family data filter failure",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'",
			prepare: func() {
				seriesSearchTask := flow.NewMockQueryTask(ctrl)
				seriesSearchTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					shardExecuteContext.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2)
					return seriesSearchTask
				}
				familyTask := flow.NewMockQueryTask(ctrl)
				familyTask.EXPECT().Run().Return(constants.ErrNotFound).MaxTimes(2)
				familyTask.EXPECT().Run().Return(io.ErrClosedPipe)
				newFamilyFilterTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					return familyTask
				}
			},
			assert: func() {
				assert.Equal(t, io.ErrClosedPipe, queryFlow.err)
			},
		},
		{
			name: "family data filter, not data found",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'",
			prepare: func() {
				seriesSearchTask := flow.NewMockQueryTask(ctrl)
				seriesSearchTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					shardExecuteContext.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2)
					return seriesSearchTask
				}
				familyTask := flow.NewMockQueryTask(ctrl)
				familyTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newFamilyFilterTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					shardExecuteContext.SeriesIDsAfterFiltering = roaring.New()
					return familyTask
				}
			},
			assert: func() {
				assert.NoError(t, queryFlow.err)
			},
		},
		{
			name: "get group context failure",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00' group by node",
			prepare: func() {
				seriesSearchTask := flow.NewMockQueryTask(ctrl)
				seriesSearchTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					shardExecuteContext.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2)
					return seriesSearchTask
				}
				familyTask := flow.NewMockQueryTask(ctrl)
				familyTask.EXPECT().Run().Return(nil).MaxTimes(3)
				newFamilyFilterTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					shardExecuteContext.TimeSegmentContext.SeriesIDs = roaring.BitmapOf(1, 2)
					return familyTask
				}
				groupTask := flow.NewMockQueryTask(ctrl)
				groupTask.EXPECT().Run().Return(io.ErrClosedPipe)
				newGroupingContextFindTaskFunc = func(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
					return groupTask
				}
			},
			assert: func() {
				assert.Equal(t, io.ErrClosedPipe, queryFlow.err)
			},
		},
		{
			name: "build group failure",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00' group by node",
			prepare: func() {
				mockGroupData(ctrl)
				task := flow.NewMockQueryTask(ctrl)
				task.EXPECT().Run().Return(io.ErrClosedPipe)
				newBuildGroupTaskFunc = func(shard tsdb.Shard,
					loadCtx *flow.DataLoadContext) flow.QueryTask {
					return task
				}
			},
			assert: func() {
				assert.Equal(t, io.ErrClosedPipe, queryFlow.err)
			},
		},
		{
			name: "data scan failure",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00' group by node",
			prepare: func() {
				mockGroupData(ctrl)
				task := flow.NewMockQueryTask(ctrl)
				task.EXPECT().Run().Return(nil).AnyTimes()
				newBuildGroupTaskFunc = func(shard tsdb.Shard,
					loadCtx *flow.DataLoadContext) flow.QueryTask {
					loadCtx.ShardExecuteCtx = &flow.ShardExecuteContext{
						StorageExecuteCtx: &flow.StorageExecuteContext{
							Query: &stmt.Query{
								GroupBy: []string{"node"},
							},
						},
					}
					loadCtx.GroupingSeriesAgg = []*flow.GroupingSeriesAgg{nil}
					return task
				}
				scanTask := flow.NewMockQueryTask(ctrl)
				scanTask.EXPECT().Run().Return(io.ErrClosedPipe)
				newDataLoadTaskFunc = func(shard tsdb.Shard, queryFlow flow.StorageQueryFlow,
					dataLoadCtx *flow.DataLoadContext, segmentIdx int, segmentCtx *flow.TimeSegmentResultSet) flow.QueryTask {
					return scanTask
				}
			},
			assert: func() {
				assert.Equal(t, io.ErrClosedPipe, queryFlow.err)
			},
		},
		{
			name: "no group data",
			sql:  "select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00' group by node",
			prepare: func() {
				mockGroupData(ctrl)
				task := flow.NewMockQueryTask(ctrl)
				task.EXPECT().Run().Return(nil).AnyTimes()
				newBuildGroupTaskFunc = func(shard tsdb.Shard,
					loadCtx *flow.DataLoadContext) flow.QueryTask {
					loadCtx.ShardExecuteCtx = &flow.ShardExecuteContext{
						StorageExecuteCtx: &flow.StorageExecuteContext{
							Query: &stmt.Query{
								GroupBy: []string{"node"},
							},
						},
					}
					return task
				}
				newDataLoadTaskFunc = func(shard tsdb.Shard, queryFlow flow.StorageQueryFlow,
					dataLoadCtx *flow.DataLoadContext, segmentIdx int, segmentCtx *flow.TimeSegmentResultSet) flow.QueryTask {
					return task
				}
			},
			assert: func() {
				assert.Nil(t, queryFlow.err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newSeriesSearchFunc = newSeriesSearch
				newFamilyFilterTaskFunc = newFamilyFilterTask
				newGroupingContextFindTaskFunc = newGroupingContextFindTask
				newBuildGroupTaskFunc = newBuildGroupTask
				newDataLoadTaskFunc = newDataLoadTask
				queryFlow.err = nil
				ctx.shards = nil
			}()
			q, _ := sql.Parse(tt.sql)
			query := q.(*stmt.Query)
			ctx.storageExecuteCtx.Query = query
			exec := newStorageMetricQuery(queryFlow, ctx)

			if tt.prepare != nil {
				tt.prepare()
			}

			exec.Execute()

			if tt.assert != nil {
				tt.assert()
			}
		})
	}
}

func mockGroupData(ctrl *gomock.Controller) {
	seriesSearchTask := flow.NewMockQueryTask(ctrl)
	seriesSearchTask.EXPECT().Run().Return(nil).MaxTimes(3)
	newSeriesIDsSearchTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
		shardExecuteContext.SeriesIDsAfterFiltering = roaring.BitmapOf(1, 2)
		return seriesSearchTask
	}
	familyTask := flow.NewMockQueryTask(ctrl)
	familyTask.EXPECT().Run().Return(nil).MaxTimes(3)
	newFamilyFilterTaskFunc = func(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
		shardExecuteContext.TimeSegmentContext.SeriesIDs = roaring.BitmapOf(1, 2)
		shardExecuteContext.StorageExecuteCtx = &flow.StorageExecuteContext{
			DownSamplingSpecs: aggregation.AggregatorSpecs{aggregation.NewAggregatorSpec("f", field.SumField)},
		}
		shardExecuteContext.TimeSegmentContext.TimeSegments = map[int64]*flow.TimeSegmentResultSet{10: nil}
		return familyTask
	}
	groupTask := flow.NewMockQueryTask(ctrl)
	groupTask.EXPECT().Run().Return(nil).AnyTimes()
	newGroupingContextFindTaskFunc = func(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
		return groupTask
	}
}

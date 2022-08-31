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

package brokerquery

import (
	"context"
	"fmt"
	"io"
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

func Test_MetricQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newBrokerPlanFn = newBrokerPlan
		ctrl.Finish()
	}()

	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	stateMgr := broker.NewMockStateManager(ctrl)
	stateMgr.EXPECT().GetCurrentNode().Return(currentNode).AnyTimes()
	taskManager := NewMockTaskManager(ctrl)

	queryFactory := &queryFactory{
		stateMgr:    stateMgr,
		taskManager: taskManager,
	}
	brokerNodes := []models.StatelessNode{
		generateBrokerActiveNode("1.1.1.1", 8000),
		generateBrokerActiveNode("1.1.1.2", 8000),
		currentNode,
		generateBrokerActiveNode("1.1.1.4", 8000),
	}
	stateMgr.EXPECT().GetLiveNodes().Return(brokerNodes).AnyTimes()

	q, err := sql.Parse("select f from cpu")
	assert.NoError(t, err)
	opt := &option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}

	storageNodes := map[string][]models.ShardID{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 6, 9},
		"1.1.1.3:9000": {5, 7, 8},
		"1.1.1.4:9000": {10, 13, 15},
		"1.1.1.5:9000": {11, 12, 14},
	}
	cases := []struct {
		name    string
		prepare func() context.Context
		wantErr bool
	}{
		{
			name: "database not found",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").Return(models.Database{}, false)
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "storage nodes not exist",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").Return(nil, nil)
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "storage replica failure",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").Return(nil, fmt.Errorf("err"))
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "broker plan failure",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
				newBrokerPlanFn = func(_ *stmt.Query, _ models.Database,
					_ map[string][]models.ShardID, _ models.StatelessNode,
					_ []models.StatelessNode) *brokerPlan {
					return &brokerPlan{}
				}
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "submit task failure",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
				taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("err"))
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "timeout",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
				eventCh1 := make(chan *series.TimeSeriesEvent)
				taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(eventCh1, nil)
				ctx, cancel := context.WithCancel(context.Background())
				time.AfterFunc(time.Millisecond*200, cancel)
				return ctx
			},
			wantErr: true,
		},
		{
			name: "has error",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
					// has error
				eventCh2 := make(chan *series.TimeSeriesEvent)
				taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventCh2, nil)
				time.AfterFunc(time.Millisecond*200, func() {
					eventCh2 <- &series.TimeSeriesEvent{Err: io.ErrClosedPipe}
				})
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "close chan",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
					// has error
				eventCh2 := make(chan *series.TimeSeriesEvent)
				taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventCh2, nil)
				time.AfterFunc(time.Millisecond*200, func() { close(eventCh2) })
				return context.Background()
			},
			wantErr: true,
		},
		{
			name: "build result set failure",
			prepare: func() context.Context {
				stateMgr.EXPECT().GetDatabaseCfg("test_db").
					Return(models.Database{Option: opt}, true)
				stateMgr.EXPECT().GetQueryableReplicas("test_db").
					Return(storageNodes, nil)
				eventCh2 := make(chan *series.TimeSeriesEvent)
				taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventCh2, nil)
				time.AfterFunc(time.Millisecond*100, func() {
					eventCh2 <- &series.TimeSeriesEvent{}
				})
				query := q.(*stmt.Query)
				query.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{Expr: &stmt.FieldExpr{Name: "f1"}},
				}
				return context.Background()
			},
			wantErr: true,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newBrokerPlanFn = newBrokerPlan
			}()
			ctx := tt.prepare()
			qry := newMetricQuery(ctx,
				&models.StatelessNode{},
				"test_db",
				q.(*stmt.Query),
				queryFactory)
			_, err = qry.WaitResponse()
			if tt.wantErr != (err != nil) {
				t.Error(err)
			}
		})
	}
}

func Test_MetricQuery_makeResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newExpressionFn = aggregation.NewExpression
		ctrl.Finish()
	}()
	expression := aggregation.NewMockExpression(ctrl)
	newExpressionFn = func(_ timeutil.TimeRange, _ int64, _ []stmt.Expr) aggregation.Expression {
		return expression
	}
	expression.EXPECT().Eval(gomock.Any()).AnyTimes()

	var now, _ = timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	timeSeries := series.NewMockGroupedIterator(ctrl)
	timeSeries.EXPECT().Tags().Return("node2").AnyTimes()

	cases := []struct {
		name    string
		prepare func(stmtQuery *stmt.Query) series.GroupedIterator
		wantErr bool
	}{
		{
			name: "no order by",
			prepare: func(_ *stmt.Query) series.GroupedIterator {
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": collections.NewFloatArray(10)}).MaxTimes(2)
				return timeSeries
			},
		},
		{
			name: "group values not match",
			prepare: func(_ *stmt.Query) series.GroupedIterator {
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": collections.NewFloatArray(10)}).MaxTimes(2)
				timeSeries1 := series.NewMockGroupedIterator(ctrl)
				timeSeries1.EXPECT().Tags().Return("node1,node2").MaxTimes(2)
				return timeSeries1
			},
		},
		{
			name: "cannot parse order by function", prepare: func(query *stmt.Query) series.GroupedIterator {
				query.OrderByItems = []stmt.Expr{&stmt.OrderByExpr{}}
				return timeSeries
			},
			wantErr: true,
		},
		{
			name: "order by field",
			prepare: func(query *stmt.Query) series.GroupedIterator {
				query.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{Expr: &stmt.FieldExpr{Name: "f1"}},
				}
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": collections.NewFloatArray(10)}).MaxTimes(2)
				return timeSeries
			},
		},
		{
			name: "order by function",
			prepare: func(query *stmt.Query) series.GroupedIterator {
				query.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{
						Expr: &stmt.CallExpr{
							FuncType: function.Avg,
							Params:   []stmt.Expr{&stmt.FieldExpr{Name: "f1"}},
						},
					},
				}
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": collections.NewFloatArray(10)}).MaxTimes(2)
				return timeSeries
			},
		},
		{
			name: "no order by, not data",
			prepare: func(_ *stmt.Query) series.GroupedIterator {
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": nil}).MaxTimes(2)
				return timeSeries
			},
		},
		{
			name: "order by field, but not data",
			prepare: func(query *stmt.Query) series.GroupedIterator {
				query.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{Expr: &stmt.FieldExpr{Name: "f1"}},
				}
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": nil}).MaxTimes(2)
				return timeSeries
			},
		},
		{
			name: "order by field, has data",
			prepare: func(query *stmt.Query) series.GroupedIterator {
				query.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{Expr: &stmt.FieldExpr{Name: "f1"}},
				}
				values := collections.NewFloatArray(2)
				values.SetValue(0, math.NaN())
				values.SetValue(1, 1.0)
				expression.EXPECT().ResultSet().
					Return(map[string]*collections.FloatArray{"f1": values}).MaxTimes(2)
				return timeSeries
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			qry := &metricQuery{
				root: &models.StatelessNode{},
				stmtQuery: &stmt.Query{
					MetricName: "1",
					Interval:   timeutil.Interval(timeutil.OneMinute),
					TimeRange: timeutil.TimeRange{
						Start: now,
						End:   now + timeutil.OneHour*2,
					},
					GroupBy: []string{"node"},
					Limit:   10,
				},
			}
			ts := tt.prepare(qry.stmtQuery)
			_, err := qry.makeResultSet(&series.TimeSeriesEvent{
				AggregatorSpecs: map[string]*protoCommonV1.AggregatorSpec{
					"f1": {
						FieldName: "f1",
						FieldType: uint32(field.SumField),
					},
				},
				SeriesList: []series.GroupedIterator{ts, ts},
				Stats: &models.QueryStats{
					TotalCost:   100,
					ExpressCost: 200,
				},
			})
			if tt.wantErr != (err != nil) {
				t.Error(tt.name)
			}
		})
	}
}

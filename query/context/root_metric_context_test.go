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

package context

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	commonmodels "github.com/lindb/common/models"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

func TestRootMetricContext_WaitResponse(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		metricCtx := NewRootMetricContext(&RootMetricContextDeps{
			Ctx: ctx,
		})
		go func() {
			cancel()
		}()
		resp, err := metricCtx.WaitResponse()
		assert.Nil(t, resp)
		assert.Equal(t, constants.ErrTimeout, err)
	})
	t.Run("complete with result", func(t *testing.T) {
		metricCtx := NewRootMetricContext(&RootMetricContextDeps{
			Ctx:       context.TODO(),
			Statement: &stmt.Query{},
		})
		go func() {
			close(metricCtx.doneCh)
		}()
		resp, err := metricCtx.WaitResponse()
		assert.NotNil(t, resp)
		assert.NoError(t, err)
	})
}

func TestRootMetricDataContext_MakePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := models.Database{
		Option: &option.DatabaseOption{
			Intervals: option.Intervals{
				{Interval: timeutil.Interval(commontimeutil.OneSecond)},
				{Interval: timeutil.Interval(commontimeutil.OneMinute)},
			},
		},
	}
	stateMgr := broker.NewMockStateManager(ctrl)
	metricCtx := NewRootMetricContext(&RootMetricContextDeps{
		Ctx:     context.TODO(),
		Choose:  stateMgr,
		Request: &models.Request{},
		Statement: &stmt.Query{
			GroupBy: []string{"ip"},
		},
	})
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "choose failure",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "empty plan",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "database config not found",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(models.Database{}, false)
			},
			wantErr: true,
		},
		{
			name: "plan invalid",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(cfg, true)
			},
			wantErr: true,
		},
		{
			name: "make plan successfully",
			prepare: func() {
				stateMgr.EXPECT().Choose(gomock.Any(), gomock.Any()).Return([]*models.PhysicalPlan{{
					Database: "test",
					Targets:  []*models.Target{{}},
				}}, nil)
				stateMgr.EXPECT().GetDatabaseCfg(gomock.Any()).Return(cfg, true)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			err := metricCtx.MakePlan()
			if (err != nil) != tt.wantErr {
				t.Fatalf("run test case fail, case: %s", tt.name)
			}
		})
	}
}

func TestRootMetricDataContext_makeResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newExpressionFn = aggregation.NewExpression
		newResultLimiterFn = aggregation.NewResultLimiter
		ctrl.Finish()
	}()
	expr := aggregation.NewMockExpression(ctrl)
	orderBy := aggregation.NewMockOrderBy(ctrl)
	newExpressionFn = func(_ timeutil.TimeRange, _ int64, _ []stmt.Expr) aggregation.Expression {
		return expr
	}
	newResultLimiterFn = func(_ int) aggregation.OrderBy {
		return orderBy
	}
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)

	cases := []struct {
		name    string
		prepare func(ctx *RootMetricContext)
		assert  func(rs *commonmodels.ResultSet, err error)
	}{
		{
			name: "order by unknown field type",
			prepare: func(ctx *RootMetricContext) {
				ctx.Deps.Statement.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{
						Expr: &stmt.CallExpr{
							FuncType: function.Unknown,
							Params:   []stmt.Expr{&stmt.FieldExpr{}},
						},
					},
				}
			},
			assert: func(rs *commonmodels.ResultSet, err error) {
				assert.Nil(t, rs)
				assert.Error(t, err)
			},
		},
		{
			name: "build order by result set",
			prepare: func(ctx *RootMetricContext) {
				ctx.Deps.Statement.OrderByItems = []stmt.Expr{
					&stmt.OrderByExpr{
						Expr: &stmt.FieldExpr{
							Name: "f",
						},
					},
				}
			},
			assert: func(rs *commonmodels.ResultSet, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
		{
			name: "build result set",
			prepare: func(ctx *RootMetricContext) {
				ctx.Deps.Statement.GroupBy = []string{"a"}
				ctx.groupAgg = groupAgg
				groupIt := series.NewMockGroupedIterator(ctrl)
				groupAgg.EXPECT().ResultSet().Return(series.GroupedIterators{groupIt})
				expr.EXPECT().Eval(gomock.Any())
				groupIt.EXPECT().Tags().Return("tags")
				expr.EXPECT().ResultSet().Return(map[string]*collections.FloatArray{"f": collections.NewFloatArray(10)})
				orderBy.EXPECT().Push(gomock.Any())
				row := aggregation.NewMockRow(ctrl)
				row.EXPECT().ResultSet().Return("a,c", nil)                                        // group by not match
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"f": nil}) // field no value
				values := collections.NewFloatArray(10)
				values.SetValue(0, 1.1)
				values.SetValue(5, math.NaN())
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"f": values})
				orderBy.EXPECT().ResultSet().Return([]aggregation.Row{row, row, row})
			},
			assert: func(rs *commonmodels.ResultSet, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
		{
			name: "build all fields result set",
			prepare: func(ctx *RootMetricContext) {
				ctx.Deps.Statement.GroupBy = []string{"a"}
				ctx.Deps.Statement.AllFields = true
				ctx.groupAgg = groupAgg
				groupIt := series.NewMockGroupedIterator(ctrl)
				groupAgg.EXPECT().Fields().Return([]field.Name{"f"})
				groupAgg.EXPECT().ResultSet().Return(series.GroupedIterators{groupIt})
				expr.EXPECT().Eval(gomock.Any())
				groupIt.EXPECT().Tags().Return("tags")
				expr.EXPECT().ResultSet().Return(map[string]*collections.FloatArray{"f": collections.NewFloatArray(10)})
				orderBy.EXPECT().Push(gomock.Any())
				row := aggregation.NewMockRow(ctrl)
				row.EXPECT().ResultSet().Return("a,c", nil)                                        // group by not match
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"f": nil}) // field no value
				values := collections.NewFloatArray(10)
				values.SetValue(0, 1.1)
				values.SetValue(5, math.NaN())
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"f": values})
				orderBy.EXPECT().ResultSet().Return([]aggregation.Row{row, row, row})
			},
			assert: func(rs *commonmodels.ResultSet, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
		{
			name: "build all fields(histogram) result set",
			prepare: func(ctx *RootMetricContext) {
				ctx.Deps.Statement.GroupBy = []string{"a"}
				ctx.Deps.Statement.AllFields = true
				ctx.groupAgg = groupAgg
				groupIt := series.NewMockGroupedIterator(ctrl)
				groupAgg.EXPECT().Fields().Return([]field.Name{"__bucket_1"})
				groupAgg.EXPECT().ResultSet().Return(series.GroupedIterators{groupIt})
				expr.EXPECT().Eval(gomock.Any())
				groupIt.EXPECT().Tags().Return("tags")
				expr.EXPECT().ResultSet().Return(map[string]*collections.FloatArray{"__bucket_1": collections.NewFloatArray(10)})
				orderBy.EXPECT().Push(gomock.Any())
				row := aggregation.NewMockRow(ctrl)
				row.EXPECT().ResultSet().Return("a,c", nil)                                                 // group by not match
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"__bucket_1": nil}) // field no value
				values := collections.NewFloatArray(10)
				values.SetValue(0, 1.1)
				values.SetValue(5, math.NaN())
				row.EXPECT().ResultSet().Return("a", map[string]*collections.FloatArray{"__bucket_1": values})
				orderBy.EXPECT().ResultSet().Return([]aggregation.Row{row, row, row})
			},
			assert: func(rs *commonmodels.ResultSet, err error) {
				assert.NotNil(t, rs)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			metricCtx := NewRootMetricContext(&RootMetricContextDeps{
				Ctx:     context.TODO(),
				Request: &models.Request{},
				Statement: &stmt.Query{
					GroupBy: []string{"ip"},
				},
			})
			metricCtx.stats = &commonmodels.NodeStats{}
			metricCtx.aggregatorSpecs = map[string]*protoCommonV1.AggregatorSpec{
				"f": {
					FieldType: uint32(field.Sum),
				},
			}
			tt.prepare(metricCtx)
			rs, err := metricCtx.makeResultSet()
			tt.assert(rs, err)
		})
	}
}

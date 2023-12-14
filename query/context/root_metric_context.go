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
	"errors"
	"fmt"
	"github.com/lindb/lindb/sql"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"

	commonmodels "github.com/lindb/common/models"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/tracker"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

var (
	newExpressionFn    = aggregation.NewExpression
	newGroupingAgg     = aggregation.NewGroupingAggregator
	newResultLimiterFn = aggregation.NewResultLimiter
)

// RootMetricContextDeps represents root metric data search dependency.
type RootMetricContextDeps struct {
	Ctx          context.Context
	Request      *models.Request
	Database     string
	CurrentNode  models.StatelessNode
	Statement    *stmt.Query
	Choose       flow.NodeChoose
	TransportMgr rpc.TransportManager
}

// RootMetricContext represents root metric data search context.
type RootMetricContext struct {
	MetricContext

	Deps *RootMetricContextDeps
}

// NewRootMetricContext creates the root metric data search context.
func NewRootMetricContext(deps *RootMetricContextDeps) *RootMetricContext {
	return &RootMetricContext{
		MetricContext: newMetricContext(deps.Ctx, deps.TransportMgr),
		Deps:          deps,
	}
}

// MakePlan makes the metric data physical plan.
func (ctx *RootMetricContext) MakePlan() error {
	database := ctx.Deps.Database
	computeNodes := 1
	if ctx.Deps.Statement.HasGroupBy() {
		// max node num
		// TODO: need config?
		computeNodes = 5
	}
	physicalPlans, err := ctx.Deps.Choose.Choose(database, computeNodes)
	if err != nil {
		return err
	}
	if len(physicalPlans) == 0 {
		return constants.ErrTargetNodesNotFound
	}
	stateMgr, ok := ctx.Deps.Choose.(broker.StateManager)
	if ok {
		databaseCfg, ok := stateMgr.GetDatabaseCfg(database)
		if !ok {
			return constants.ErrDatabaseNotExist
		}
		calcTimeRangeAndInterval(ctx.Deps.Statement, databaseCfg)
	}
	payload, _ := ctx.Deps.Statement.MarshalJSON()
	for _, physicalPlan := range physicalPlans {
		//FIXME:
		physicalPlan.AddReceiver(ctx.Deps.CurrentNode.Indicator())
		if err := physicalPlan.Validate(); err != nil {
			return err
		}
		ctx.addRequests(
			&protoCommonV1.TaskRequest{
				RequestID:    ctx.Deps.Request.RequestID,
				RequestType:  protoCommonV1.RequestType_Data,
				PhysicalPlan: encoding.JSONMarshal(physicalPlan),
				Payload:      payload,
			}, physicalPlan)
	}
	return nil
}

// WaitResponse waits metric data search task completed, then returns the result set,
func (ctx *RootMetricContext) WaitResponse() (any, error) {
	err := ctx.waitResponse()
	if err != nil {
		return nil, err
	}

	return ctx.makeResultSet()
}

// makeResultSet makes final result set from time series event(GroupedIterators).
// TODO: can opt use stream, leaf node need return grouping if completed.
func (ctx *RootMetricContext) makeResultSet() (resultSet *commonmodels.ResultSet, err error) {
	makeResultStartTime := time.Now()
	orderBy, err := ctx.buildOrderBy()
	if err != nil {
		return nil, err
	}

	statement := ctx.Deps.Statement
	resultSet = new(commonmodels.ResultSet)
	// TODO: merge stats for cross idc query?
	groupByKeys := statement.GroupBy
	groupByKeysLength := len(groupByKeys)
	fieldsMap := make(map[string]struct{})
	timeRange := ctx.timeRange
	interval := ctx.interval
	if ctx.groupAgg != nil {
		groupIts := ctx.groupAgg.ResultSet()
		selectItems := ctx.getSelectItems()
		for _, it := range groupIts {
			// TODO: reuse expression??
			expression := newExpressionFn(
				timeRange,
				interval,
				selectItems,
			)
			// do expression eval
			expression.Eval(it)

			// result order by/limit
			orderBy.Push(aggregation.NewOrderByRow(it.Tags(), expression.ResultSet()))
		}

		rows := orderBy.ResultSet()
		for _, row := range rows {
			var tags map[string]string
			tagValues, fields := row.ResultSet()
			if groupByKeysLength > 0 {
				tagValues := tag.SplitTagValues(tagValues)
				if groupByKeysLength != len(tagValues) {
					// if tag values not match group by tag keys, ignore this time series
					continue
				}
				// build group by tags for final result
				tags = make(map[string]string)
				for idx, tagKey := range groupByKeys {
					tags[tagKey] = tagValues[idx]
				}
			}
			timeSeries := commonmodels.NewSeries(tags, tagValues)
			resultSet.AddSeries(timeSeries)

			having := ctx.Deps.Statement.Having
			notHavingSlots := make(map[int]struct{})
			slotValues := make(map[int]map[string]float64)

			if having != nil {
				for fieldName, values := range fields {
					if values == nil {
						continue
					}
					it := values.NewIterator()
					for it.HasNext() {
						slot, val := it.Next()
						if math.IsNaN(val) {
							continue
						}
						if v, ok := slotValues[slot]; ok {
							v[fieldName] = val
						} else {
							slotValues[slot] = map[string]float64{fieldName: val}
						}
					}
				}
				// calc and fill
				if len(slotValues) > 0 {
					calc := sql.NewCalc(having)
					for slot, fieldValue := range slotValues {
						result, err := calc.CalcExpr(fieldValue)
						if err != nil {
							return resultSet, err
						}
						if r, ok := result.(bool); !ok {
							return resultSet, fmt.Errorf("expected CalcExpr bool result got %v", reflect.TypeOf(result))
						} else if !r {
							notHavingSlots[slot] = struct{}{}
						}
					}
				}
			}

			for fieldName, values := range fields {
				if values == nil {
					continue
				}

				points := commonmodels.NewPoints()
				it := values.NewIterator()
				for it.HasNext() {
					slot, val := it.Next()
					if math.IsNaN(val) {
						// TODO: need check
						continue
					}
					if _, ok := notHavingSlots[slot]; ok {
						continue
					}
					points.AddPoint(timeutil.CalcTimestamp(timeRange.Start, slot, timeutil.Interval(interval)), val)
				}
				timeSeries.AddField(fieldName, points)
				fieldsMap[fieldName] = struct{}{}
			}
		}
	}

	sort.Slice(resultSet.Series, func(i, j int) bool {
		return resultSet.Series[i].TagValues < resultSet.Series[j].TagValues
	})

	resultSet.MetricName = statement.MetricName
	resultSet.GroupBy = statement.GroupBy
	for fName := range fieldsMap {
		resultSet.Fields = append(resultSet.Fields, fName)
	}
	resultSet.StartTime = timeRange.Start
	resultSet.EndTime = timeRange.End
	resultSet.Interval = interval

	if ctx.stats != nil {
		now := time.Now()
		ctx.stats.Node = ctx.Deps.CurrentNode.Indicator()
		ctx.stats.End = now.UnixNano()
		ctx.stats.TotalCost = now.Sub(ctx.startTime).Nanoseconds()

		ctx.stats.Stages = append(ctx.stats.Stages, &commonmodels.StageStats{
			Identifier: "Expression",
			Start:      makeResultStartTime.UnixNano(),
			End:        now.UnixNano(),
			Cost:       now.Sub(makeResultStartTime).Nanoseconds(),
			State:      tracker.CompleteState.String(),
			Async:      false,
		})
		resultSet.Stats = ctx.stats
	}
	return resultSet, nil
}

// buildOrderBy builds order by container.
func (ctx *RootMetricContext) buildOrderBy() (aggregation.OrderBy, error) {
	statement := ctx.Deps.Statement
	// build order by items if need do order by query
	orderByExprs := statement.OrderByItems
	if len(orderByExprs) == 0 {
		// use default limiter
		return newResultLimiterFn(statement.Limit), nil
	}
	var orderByItems []*aggregation.OrderByItem
	fields := ctx.aggregatorSpecs
	for _, orderBy := range orderByExprs {
		expr := orderBy.(*stmt.OrderByExpr)
		funcType := function.Unknown
		var fieldName string
		switch e := expr.Expr.(type) {
		case *stmt.FieldExpr:
			aggSpec, ok := fields[e.Name]
			if ok {
				funcType = field.Type(aggSpec.FieldType).GetOrderByFunc()
				fieldName = e.Name
			}
		case *stmt.CallExpr:
			funcType = e.FuncType
			fieldName = e.Params[0].Rewrite()
		}
		if funcType == function.Unknown {
			return nil, errors.New("cannot parse order by function")
		}
		orderByItems = append(orderByItems, &aggregation.OrderByItem{
			Expr:     expr,
			Name:     fieldName,
			FuncType: funcType,
			Desc:     expr.Desc,
		})
	}
	return aggregation.NewTopNOrderBy(orderByItems, statement.Limit), nil
}

// getSelectItems returns select field items.
func (ctx *RootMetricContext) getSelectItems() []stmt.Expr {
	statement := ctx.Deps.Statement
	selectItems := statement.SelectItems
	if statement.AllFields {
		// if select all fields, read field names from aggregator
		allAggFields := ctx.groupAgg.Fields()
		selectItems = []stmt.Expr{}
		isHistogram := false
		for _, fieldName := range allAggFields {
			if strings.HasPrefix(string(fieldName), "__bucket_") {
				// filter histogram raw field
				isHistogram = true
				continue
			}
			selectItems = append(selectItems, &stmt.SelectItem{Expr: &stmt.FieldExpr{Name: fieldName.String()}})
		}
		if isHistogram {
			// add histogram functions
			addQuantileFn := func(as string, num float64) {
				selectItems = append(selectItems, &stmt.SelectItem{
					Expr:  &stmt.CallExpr{FuncType: function.Quantile, Params: []stmt.Expr{&stmt.NumberLiteral{Val: num}}},
					Alias: as,
				})
			}
			addQuantileFn("p99", 0.99)
			addQuantileFn("p95", 0.95)
			addQuantileFn("p90", 0.90)
			addQuantileFn("mean", 0.50)
		}
	}
	return selectItems
}

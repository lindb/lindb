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

package sql

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// queryStmtParse represents query statement parser using visitor
type queryStmtParse struct {
	baseStmtParser
	explain bool

	selectItems []stmt.Expr
	fieldNames  map[string]struct{}

	startTime int64
	endTime   int64

	//orderByExpr stmt.Expr
	//desc        bool
	groupBy  []string
	interval int64
	fieldID  int
}

// newQueryStmtParse create a query statement parser
func newQueryStmtParse(explain bool) *queryStmtParse {
	return &queryStmtParse{
		explain:    explain,
		fieldNames: make(map[string]struct{}),
		fieldID:    1,
		baseStmtParser: baseStmtParser{
			exprStack: collections.NewStack(),
			namespace: constants.DefaultNamespace,
			limit:     20,
		},
	}
}

// build builds query statement based on parse result
func (q *queryStmtParse) build() (stmt.Statement, error) {
	if err := q.validation(); err != nil {
		return nil, err
	}

	query := &stmt.Query{}
	query.Explain = q.explain
	query.Namespace = q.namespace
	query.MetricName = q.metricName
	query.SelectItems = q.selectItems
	query.Condition = q.condition

	fieldNames := make([]string, len(q.fieldNames))
	idx := 0
	for fieldName := range q.fieldNames {
		fieldNames[idx] = fieldName
		idx++
	}
	sort.Slice(fieldNames, func(i, j int) bool {
		return fieldNames[i] < fieldNames[j]
	})
	query.FieldNames = fieldNames

	now := timeutil.Now()
	query.TimeRange = timeutil.TimeRange{Start: q.startTime, End: q.endTime}
	if query.TimeRange.Start <= 0 {
		query.TimeRange.Start = now - timeutil.OneHour
	}
	if query.TimeRange.End <= 0 {
		query.TimeRange.End = now
	}
	if query.TimeRange.End < query.TimeRange.Start {
		return nil, fmt.Errorf("start time cannot be larger than end time")
	}

	query.Interval = timeutil.Interval(q.interval)
	query.GroupBy = q.groupBy
	query.Limit = q.limit
	return query, nil
}

// validation tests data if invalid
func (q *queryStmtParse) validation() error {
	if q.err != nil {
		return q.err
	}
	if len(q.metricName) == 0 {
		return fmt.Errorf("metric name cannot be empty")
	}
	if len(q.selectItems) == 0 {
		return fmt.Errorf("select fields cannbe be empty")
	}
	return nil
}

// resetExprStack resets expr stack for next parse fragment
func (q *queryStmtParse) resetExprStack() {
	q.exprStack = collections.NewStack()
}

// visitGroupByKey visits when production groupBy key expression is entered
func (q *queryStmtParse) visitGroupByKey(ctx *grammar.GroupByKeyContext) {
	switch {
	case ctx.Ident() != nil:
		tagKey := strutil.GetStringValue(ctx.Ident().GetText())
		q.groupBy = append(q.groupBy, tagKey)
	case ctx.DurationLit() != nil:
		q.interval = q.parseDuration(ctx.DurationLit())
	}
}

// visitTimeRangeExpr visits when production timeRange expression is entered
func (q *queryStmtParse) visitTimeRangeExpr(ctx *grammar.TimeRangeExprContext) {
	timeExprCtxList := ctx.AllTimeExpr()
	for _, timeExpr := range timeExprCtxList {
		timeExprCtx, ok := timeExpr.(*grammar.TimeExprContext)
		if !ok {
			continue
		}
		var timestamp int64
		var err error
		switch {
		case timeExprCtx.Ident() != nil:
			timestamp, err = timeutil.ParseTimestamp(strutil.GetStringValue(timeExprCtx.Ident().GetText()))
		case timeExprCtx.NowExpr() != nil:
			timestamp = timeutil.Now()
			durationExpr, ok := timeExprCtx.NowExpr().(*grammar.NowExprContext)
			if ok {
				timestamp += q.parseDuration(durationExpr.DurationLit())
			}
		}
		if err != nil {
			q.err = err
			continue
		}
		binaryOp := timeExprCtx.BinaryOperator()
		if binaryOp == nil {
			continue
		}
		binaryOpCtx, ok := binaryOp.(*grammar.BinaryOperatorContext)
		if !ok {
			continue
		}
		if binaryOpCtx.T_GREATER() != nil || binaryOpCtx.T_GREATEREQUAL() != nil {
			q.startTime = timestamp
		}
		if binaryOpCtx.T_LESS() != nil || binaryOpCtx.T_LESSEQUAL() != nil {
			q.endTime = timestamp
		}
	}
}

// parseDuration parses time duration from duration string
func (q *queryStmtParse) parseDuration(ctx grammar.IDurationLitContext) int64 {
	if ctx == nil {
		return 0
	}
	durationCtx, ok := ctx.(*grammar.DurationLitContext)
	if !ok {
		return 0
	}

	duration, err := strconv.ParseInt(durationCtx.IntNumber().GetText(), 10, 64)
	if err != nil {
		q.err = err
		return 0
	}
	var result int64
	if durationCtx.IntervalItem() == nil {
		return result
	}
	unit, ok := durationCtx.IntervalItem().(*grammar.IntervalItemContext)
	if !ok {
		return result
	}
	switch {
	case nil != unit.T_SECOND():
		result = duration * timeutil.OneSecond
	case nil != unit.T_MINUTE():
		result = duration * timeutil.OneMinute
	case nil != unit.T_HOUR():
		result = duration * timeutil.OneHour
	case nil != unit.T_DAY():
		result = duration * timeutil.OneDay
	case nil != unit.T_WEEK():
		result = duration * timeutil.OneWeek
	case nil != unit.T_MONTH():
		result = duration * timeutil.OneMonth
	case nil != unit.T_YEAR():
		result = duration * timeutil.OneYear
	}
	return result
}

// visitFieldExpr visits when production field expression is entered
func (q *queryStmtParse) visitFieldExpr(ctx *grammar.FieldExprContext) {
	//var selectItem stmt.Expr
	switch {
	case ctx.ExprFunc() != nil:
		q.exprStack.Push(&stmt.CallExpr{})
	case ctx.T_OPEN_P() != nil:
		q.exprStack.Push(&stmt.ParenExpr{})
	case ctx.T_MUL() != nil:
		q.exprStack.Push(&stmt.BinaryExpr{Operator: stmt.MUL})
	case ctx.T_DIV() != nil:
		q.exprStack.Push(&stmt.BinaryExpr{Operator: stmt.DIV})
	case ctx.T_ADD() != nil:
		q.exprStack.Push(&stmt.BinaryExpr{Operator: stmt.ADD})
	case ctx.T_SUB() != nil:
		q.exprStack.Push(&stmt.BinaryExpr{Operator: stmt.SUB})
	}
}

// visitAlias visits when production alias expression is entered
func (q *queryStmtParse) visitAlias(ctx *grammar.AliasContext) {
	if len(q.selectItems) == 0 {
		return
	}
	selectItem, ok := (q.selectItems[0]).(*stmt.SelectItem)
	if ok {
		selectItem.Alias = strutil.GetStringValue(ctx.Ident().GetText())
	}
}

// visitFuncName visits when production function call expression is entered
func (q *queryStmtParse) visitFuncName(ctx *grammar.FuncNameContext) {
	if q.exprStack.Empty() {
		return
	}
	callExpr, ok := q.exprStack.Peek().(*stmt.CallExpr)
	if !ok {
		return
	}
	switch {
	case ctx.T_SUM() != nil:
		callExpr.FuncType = function.Sum
	case ctx.T_MIN() != nil:
		callExpr.FuncType = function.Min
	case ctx.T_MAX() != nil:
		callExpr.FuncType = function.Max
	case ctx.T_COUNT() != nil:
		callExpr.FuncType = function.Count
	case ctx.T_AVG() != nil:
		callExpr.FuncType = function.Avg
	case ctx.T_STDDEV() != nil:
		callExpr.FuncType = function.Stddev
	case ctx.T_HISTOGRAM() != nil:
		callExpr.FuncType = function.Histogram
	}
}

// completeFuncExpr completes a function call expression for select list
func (q *queryStmtParse) completeFuncExpr() {
	cur := q.exprStack.Pop()
	if cur != nil {
		expr, ok := cur.(stmt.Expr)
		if ok {
			q.setExprParam(expr)
		}
		if q.exprStack.Empty() {
			q.selectItems = append(q.selectItems, &stmt.SelectItem{Expr: expr})
		}
	}
}

// visitExprAtom visits when production atom expr expression is entered
func (q *queryStmtParse) visitExprAtom(ctx *grammar.ExprAtomContext) {
	switch {
	case ctx.Ident() != nil:
		val := strutil.GetStringValue(ctx.Ident().GetText())
		if q.exprStack.Empty() {
			q.selectItems = append(q.selectItems, &stmt.SelectItem{Expr: &stmt.FieldExpr{Name: val}})
		} else {
			q.setExprParam(&stmt.FieldExpr{Name: val})
		}
		q.fieldNames[val] = struct{}{}
	case ctx.DecNumber() != nil || ctx.IntNumber() != nil:
		valStr := ""
		switch {
		case ctx.DecNumber() != nil:
			valStr = ctx.DecNumber().GetText()
		case ctx.IntNumber() != nil:
			valStr = ctx.IntNumber().GetText()
		}

		val, _ := strconv.ParseFloat(valStr, 64)
		if !q.exprStack.Empty() {
			q.setExprParam(&stmt.NumberLiteral{Val: val})
		}
	default:
	}
}

// completeFieldExpr completes a field expr,
// only paren and binary expr need do set expr param,
// set func's param in complete func parse section.
func (q *queryStmtParse) completeFieldExpr(ctx *grammar.FieldExprContext) {
	switch {
	case ctx.T_OPEN_P() != nil:
	case ctx.T_MUL() != nil:
	case ctx.T_DIV() != nil:
	case ctx.T_ADD() != nil:
	case ctx.T_SUB() != nil:
	default:
		return
	}

	cur := q.exprStack.Pop()

	if cur != nil {
		expr, ok := cur.(stmt.Expr)
		if ok {
			q.setExprParam(expr)
		}
		if q.exprStack.Empty() {
			q.selectItems = append(q.selectItems, &stmt.SelectItem{Expr: expr})
		}
	}
}

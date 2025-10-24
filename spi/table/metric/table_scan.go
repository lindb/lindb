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

package metric

import (
	"context"
	"fmt"

	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
	metricstore "github.com/lindb/lindb/storage/metric"
)

type column struct {
	offset  int
	meta    field.Meta
	handles []*ColumnHandle

	rollups []rollupConfig
	aggs    []aggConfig
}

func (c *column) init() {
	rollupMap := make(map[tree.FuncName]struct{})
	index := 0
	for _, handle := range c.handles {
		c.aggs = append(c.aggs, aggConfig{target: index, aggType: getAggFunc(handle.Aggregation)})
		if _, ok := rollupMap[handle.Downsampling]; !ok {
			c.rollups = append(c.rollups, rollupConfig{aggType: getAggFunc(handle.Downsampling)})
			index++
		}
	}
}

type rollupConfig struct {
	aggType field.AggType
}

type aggConfig struct {
	target  int // ref to rollup result set
	aggType field.AggType
}

type TableScan struct {
	db        *metricstore.Database
	metricID  metric.ID       // table id
	schema    *metric.Schema  // metric(table) schema
	predicate tree.Expression // where clause

	grouping *Grouping

	// TODO: check if found all filter column values
	filterResult map[tree.NodeID]*flow.TagFilterResult

	fields        field.Metas
	columns       []*column
	columnMapping map[string]string
	maxOfRollups  int
	numOfAggs     int
	outputs       []types.ColumnMetadata

	isTimestampSelected bool
	timeRange           timeutil.TimeRange
	interval            timeutil.Interval
}

func (t *TableScan) isGrouping() bool {
	return t.grouping != nil && t.grouping.tags.Len() > 0
}

func (t *TableScan) createRollups() (rs rollups) {
	rs = make(rollups, t.maxOfRollups)
	for index := range rs {
		rs[index] = newRollup(t.timeRange.NumOfPoints(t.interval), t.interval.Int64())
	}
	return
}

type ColumnValuesLookupVisitor struct {
	evalCtx   expression.EvalContext
	tableScan *TableScan
}

func NewColumnValuesLookVisitor(ctx context.Context, tableScan *TableScan) *ColumnValuesLookupVisitor {
	return &ColumnValuesLookupVisitor{
		tableScan: tableScan,
		evalCtx:   expression.NewEvalContext(ctx),
	}
}

func (v *ColumnValuesLookupVisitor) Visit(ctx any, n tree.Node) any {
	var (
		column tree.Expression
		fn     func(columnName string) tree.Expr
	)
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		columnValue, _ := expression.EvalString(v.evalCtx, node.Right)
		column = node.Left
		// TODO: add other operator
		fn = func(columnName string) tree.Expr {
			return &tree.EqualsExpr{
				Name:  columnName,
				Value: columnValue,
			}
		}
	case *tree.InPredicate:
		var values []string
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
				columnValue, _ := expression.EvalString(v.evalCtx, item)
				return columnValue
			})
		}
		column = node.Value
		fn = func(columnName string) tree.Expr {
			return &tree.InExpr{
				Name:   columnName,
				Values: values,
			}
		}
	case *tree.LikePredicate:
		columnValue, _ := expression.EvalString(v.evalCtx, node.Pattern)
		column = node.Value
		fn = func(columnName string) tree.Expr {
			return &tree.LikeExpr{
				Name:  columnName,
				Value: columnValue,
			}
		}
	case *tree.RegexPredicate:
		regexp, _ := expression.EvalString(v.evalCtx, node.Pattern)
		column = node.Value
		fn = func(columnName string) tree.Expr {
			return &tree.RegexExpr{
				Name:   columnName,
				Regexp: regexp,
			}
		}
	case *tree.NotExpression:
		return node.Value.Accept(ctx, v)
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			term.Accept(ctx, v)
		}
		return nil
	case *tree.Cast:
		return node.Expression.Accept(ctx, v)
	default:
		panic(fmt.Sprintf("column values lookup error, not support node type: %T", n))
	}
	// visit predicate which finding tag value ids
	return v.visitPredicate(n, column, fn)
}

func (v *ColumnValuesLookupVisitor) visitPredicate(predicate tree.Node, column tree.Expression,
	buildExpr func(columnName string) tree.Expr,
) (r any) {
	columnSymbols := plan.ExtractSymbolsFromExpression(column)
	if len(columnSymbols) != 1 {
		panic(fmt.Sprintf("column values lookup error, column: %s, symbol size: %d",
			tree.FormatExpression(column), len(columnSymbols)))
	}
	columnName := getColumnName(columnSymbols[0].Name, v.tableScan.columnMapping)
	fmt.Printf("column values lookup, column: %s,name=%s\n", column, columnName)

	tagMeta, ok := v.tableScan.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	expr := buildExpr(columnName)
	fmt.Printf("filter expr==%v\n", expr)
	tagValueIDs, err = v.tableScan.db.MetaDB().FindTagValueDsByExpr(tagKeyID, expr)
	if err != nil {
		panic(err)
	}

	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnValueNotFound, columnName))
	}

	if v.tableScan.filterResult == nil {
		v.tableScan.filterResult = make(map[tree.NodeID]*flow.TagFilterResult)
	}

	v.tableScan.filterResult[predicate.GetID()] = &flow.TagFilterResult{
		TagKeyID:    tagKeyID,
		TagValueIDs: tagValueIDs,
	}
	return nil
}

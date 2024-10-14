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
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/tsdb"
)

type TableScan struct {
	ctx       context.Context
	db        tsdb.Database
	schema    *metric.Schema
	predicate tree.Expression
	grouping  *Grouping

	// TODO: check if found all filter column values
	filterResult map[tree.NodeID]*flow.TagFilterResult

	fields field.Metas
	output []types.ColumnMetadata

	timeRange       timeutil.TimeRange
	interval        timeutil.Interval
	storageInterval timeutil.Interval
	metricID        metric.ID
}

func (t *TableScan) hasGrouping() bool {
	return t.grouping != nil && t.grouping.tags.Len() > 0
}

func (t *TableScan) lookupColumnValues() {
	if t.predicate != nil {
		// lookup column if predicate not nil
		lookup := NewColumnValuesLookVisitor(t)
		_ = t.predicate.Accept(nil, lookup)
		// TODO: check filter result if empty????
	}
}

type ColumnValuesLookupVisitor struct {
	evalCtx   expression.EvalContext
	tableScan *TableScan
}

func NewColumnValuesLookVisitor(tableScan *TableScan) *ColumnValuesLookupVisitor {
	// timestamp, _ := expression.EvalTime(evalCtx, translations.Rewrite(timePredicate.Value))
	return &ColumnValuesLookupVisitor{
		tableScan: tableScan,
		evalCtx:   expression.NewEvalContext(tableScan.ctx),
	}
}

func (v *ColumnValuesLookupVisitor) Visit(context any, n tree.Node) any {
	var (
		column tree.Expression
		fn     func(columnName string) tree.Expr
	)
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		columnValue, _ := expression.EvalString(v.evalCtx, node.Right)
		column = node.Left
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
		return node.Value.Accept(context, v)
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			term.Accept(context, v)
		}
		return nil
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	default:
		panic(fmt.Sprintf("column values lookup error, not support node type: %T", n))
	}
	// visit predicate which finding tag value ids
	return v.visitPredicate(context, n, column, fn)
}

func (v *ColumnValuesLookupVisitor) visitPredicate(context any,
	predicate tree.Node, column tree.Expression,
	buildExpr func(columnName string) tree.Expr,
) (r any) {
	columnName, _ := expression.EvalString(v.evalCtx, column)

	tagMeta, ok := v.tableScan.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	tagValueIDs, err = v.tableScan.db.MetaDB().FindTagValueDsByExpr(tagKeyID, buildExpr(columnName))
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

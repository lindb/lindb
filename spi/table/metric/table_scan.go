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
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.InPredicate:
		return v.visitInPredicate(context, node)
	case *tree.NotExpression:
		return node.Value.Accept(context, v)
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	default:
		panic(fmt.Sprintf("column values lookup error, not support node type: %T", n))
	}
}

func (v *ColumnValuesLookupVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	// TODO: check error
	columnName, _ := expression.EvalString(v.evalCtx, node.Left)
	columnValue, _ := expression.EvalString(v.evalCtx, node.Right)

	tagMeta, ok := v.tableScan.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	// FIXME: impl other expr
	tagValueIDs, err = v.tableScan.db.MetaDB().FindTagValueDsByExpr(tagKeyID, &tree.EqualsExpr{
		Name:  columnName,
		Value: columnValue,
	})
	if err != nil {
		panic(err)
	}

	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		// TODO: panic if not found?
		return nil
	}

	if v.tableScan.filterResult == nil {
		v.tableScan.filterResult = make(map[tree.NodeID]*flow.TagFilterResult)
	}

	v.tableScan.filterResult[node.ID] = &flow.TagFilterResult{
		TagKeyID:    tagKeyID,
		TagValueIDs: tagValueIDs,
	}
	return nil
}

func (v *ColumnValuesLookupVisitor) visitInPredicate(context any, node *tree.InPredicate) (r any) {
	columnName, _ := expression.EvalString(v.evalCtx, node.Value)

	tagMeta, ok := v.tableScan.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	var values []string
	// TODO: check values
	if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
		values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
			columnValue, _ := expression.EvalString(v.evalCtx, item)
			return columnValue
		})
	}
	// FIXME: impl other expr
	tagValueIDs, err = v.tableScan.db.MetaDB().FindTagValueDsByExpr(tagKeyID, &tree.InExpr{
		Name:   columnName,
		Values: values,
	})
	if err != nil {
		panic(err)
	}

	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		// TODO: panic if not found?
		return nil
	}

	if v.tableScan.filterResult == nil {
		v.tableScan.filterResult = make(map[tree.NodeID]*flow.TagFilterResult)
	}

	v.tableScan.filterResult[node.ID] = &flow.TagFilterResult{
		TagKeyID:    tagKeyID,
		TagValueIDs: tagValueIDs,
	}
	return nil
}

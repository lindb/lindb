package metric

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/tsdb"
)

type TableScan struct {
	db        tsdb.Database
	schema    *metric.Schema
	predicate tree.Expression
	grouping  *Grouping

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
	tableScan *TableScan
}

func NewColumnValuesLookVisitor(tableScan *TableScan) *ColumnValuesLookupVisitor {
	return &ColumnValuesLookupVisitor{
		tableScan: tableScan,
	}
}

func (v *ColumnValuesLookupVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	case *tree.Identifier:
		return node.Value
	case *tree.StringLiteral:
		return node.Value
	}
	return nil
}

func (v *ColumnValuesLookupVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	column := node.Left.Accept(context, v)
	columnValue := node.Right.Accept(context, v)

	columnName, ok := column.(string)
	if !ok {
		panic(fmt.Sprintf("column name '%v' not support type '%T'", column, column))
	}

	tagMeta, ok := v.tableScan.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	switch val := columnValue.(type) {
	case string:
		// FIXME: impl other expr
		tagValueIDs, err = v.tableScan.db.MetaDB().FindTagValueDsByExpr(tagKeyID, &tree.EqualsExpr{
			Name:  columnName,
			Value: val,
		})
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("value of column '%v' not support type '%T'", columnName, val))
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

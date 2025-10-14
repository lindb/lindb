package log

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage/log"
)

type TableScan struct {
	db        *log.Database
	timeRange timeutil.TimeRange
	interval  timeutil.Interval

	nsID uint32

	predicate tree.Expression // where clause

	filterResult map[tree.NodeID]*Fields
}

type Fields struct {
	fieldID       uint32
	fieldValueIDs []uint32
}

type FieldValuesLookupVisitor struct {
	evalCtx   expression.EvalContext
	tableScan *TableScan
}

func NewFieldValuesLookupVisitor(ctx context.Context, tableScan *TableScan) *FieldValuesLookupVisitor {
	return &FieldValuesLookupVisitor{
		tableScan: tableScan,
		evalCtx:   expression.NewEvalContext(ctx),
	}
}

func (v *FieldValuesLookupVisitor) Visit(ctx any, n tree.Node) any {
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
	case *tree.NullPredicate:
		column = node.Value
		fn = func(columnName string) tree.Expr {
			return &tree.NullExpr{
				Name: columnName,
				Not:  node.Not,
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

func (v *FieldValuesLookupVisitor) visitPredicate(predicate tree.Node, column tree.Expression,
	buildExpr func(columnName string) tree.Expr,
) (r any) {
	var columnName string
	if str, ok := column.(*tree.Constant); ok {
		columnName = fmt.Sprintf("%v", str.Value)
	} else {
		columnSymbols := plan.ExtractSymbolsFromExpression(column)
		if len(columnSymbols) != 1 {
			panic(fmt.Sprintf("column values lookup error, column: %s, symbol size: %d",
				tree.FormatExpression(column), len(columnSymbols)))
		}
		columnName = columnSymbols[0].Name
	}
	fmt.Printf("column values lookup, column: %s,name=%s\n", column, columnName)

	fieldID, err := v.tableScan.db.IndexDatabase().GetFieldKeyID(v.tableScan.nsID, []byte(columnName))
	if err != nil {
		panic(fmt.Errorf("%w, column name: %s,%d", constants.ErrColumnNotFound, columnName, v.tableScan.nsID))
	}
	var fieldValueIDs []uint32
	if _, ok := predicate.(*tree.NullPredicate); !ok {
		expr := buildExpr(columnName)
		fmt.Printf("filter expr==%v\n", expr)
		fieldValueIDs, err = v.tableScan.db.IndexDatabase().FindFieldValueIDs(fieldID, expr)
		if err != nil {
			panic(err)
		}

		if len(fieldValueIDs) == 0 {
			panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnValueNotFound, columnName))
		}
	}

	if v.tableScan.filterResult == nil {
		v.tableScan.filterResult = make(map[tree.NodeID]*Fields)
	}

	v.tableScan.filterResult[predicate.GetID()] = &Fields{
		fieldID:       fieldID,
		fieldValueIDs: fieldValueIDs,
	}
	return nil
}

func getColumnName(name string, mapping map[string]string) string {
	if len(mapping) == 0 {
		return name
	}
	realName, ok := mapping[name]
	if !ok {
		return name
	}
	return realName
}

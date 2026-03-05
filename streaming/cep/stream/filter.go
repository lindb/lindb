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

package stream

import (
	"fmt"
	"regexp"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
)

type filter struct {
	expr   Expr
	sc     *sourceConnector
	schema *arrow.Schema

	evalContext expression.EvalContext
}

func newFilter(schema *arrow.Schema, expr tree.Expression, sc *sourceConnector) *filter {
	f := &filter{
		sc:          sc,
		schema:      schema,
		evalContext: expression.NewEvalContext(sc.ctx),
	}
	f.expr = f.rewrite(expr)
	return f
}

func (f *filter) eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	return f.expr.Eval(record)
}

func (f *filter) rewrite(n tree.Expression) Expr {
	switch node := n.(type) {
	case *tree.InPredicate:
		col := f.resolveColValue(node.Value)
		var values []string
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			values = lo.Map(inListExpression.Values, func(item tree.Expression, _ int) string {
				val, err := expression.EvalString(f.evalContext, item)
				if err != nil {
					panic(err)
				}
				return val
			})
		}
		return &InExpr{col: col, values: values}
	case *tree.ComparisonExpression:
		return &ComparisonExpr{
			left:  f.resolveColValue(node.Left),
			right: f.resolveColValue(node.Right),
		}
	case *tree.LikePredicate:
		pattern, err := expression.EvalString(f.evalContext, node.Pattern)
		if err != nil {
			panic(err)
		}
		return &RegexExpr{col: f.resolveColValue(node.Value), pattern: likeToRegexp(pattern)}
	case *tree.RegexPredicate:
		pattern, err := expression.EvalString(f.evalContext, node.Pattern)
		if err != nil {
			panic(err)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic(fmt.Errorf("invalid regex pattern %q: %w", pattern, err))
		}
		return &RegexExpr{col: f.resolveColValue(node.Value), pattern: re}
	case *tree.NullPredicate:
		return &NullExpr{col: f.resolveColValue(node.Value), not: node.Not}
	case *tree.NotExpression:
		return &NotExpr{expr: f.rewrite(node.Value)}
	case *tree.LogicalExpression:
		exprs := lo.Map(node.Terms, func(term tree.Expression, _ int) Expr {
			return f.rewrite(term)
		})
		return &LogicalExpr{op: node.Operator, exprs: exprs}
	default:
		panic(fmt.Errorf("where clause: not supported expression type %T", n))
	}
}

// resolveColValue resolves an expression to a colValue.
// Column references (Identifier, SubscriptExpression) map to runtime column reads;
// literals (StringLiteral etc.) become a constCol.
func (f *filter) resolveColValue(expr tree.Expression) colValue {
	fmt.Printf("resolving column value for expression: %T\n", expr)
	switch e := expr.(type) {
	case *tree.SubscriptExpression:
		key, err := expression.EvalString(f.evalContext, e.Key)
		if err != nil {
			panic(err)
		}
		return &subscriptCol{colIndex: f.getColIndex(e.Base), key: key}
	case *tree.SymbolReference:
		indexes := f.schema.FieldIndices(e.Name)
		if len(indexes) != 1 {
			panic(fmt.Sprintf("invalid column %s, found %d fields", e.Name, len(indexes)))
		}
		return &directCol{colIndex: indexes[0]}
	default:
		// literals and other scalar expressions evaluate to a constant
		val, err := expression.EvalString(f.evalContext, expr)
		if err != nil {
			panic(err)
		}
		// check if the string value matches a column name
		indexes := f.schema.FieldIndices(val)
		if len(indexes) == 1 {
			return &directCol{colIndex: indexes[0]}
		}
		return &constCol{val: val}
	}
}

func (f *filter) getColIndex(expr tree.Expression) int {
	var colName string
	var err error
	switch e := expr.(type) {
	case *tree.SymbolReference:
		colName = e.Name
	default:
		colName, err = expression.EvalString(f.evalContext, expr)
		if err != nil {
			panic(err)
		}
	}
	indexes := f.schema.FieldIndices(colName)
	if len(indexes) != 1 {
		panic(fmt.Sprintf("invalid column %s, found %d fields", colName, len(indexes)))
	}
	return indexes[0]
}

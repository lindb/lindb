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

package expression

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type RewriteContext struct {
	EvalContext  EvalContext
	SourceLayout []*plan.Symbol
}

func Rewrite(ctx *RewriteContext, node tree.Expression) Expression {
	return (&rewriter{ctx: ctx}).rewrite(node)
}

type rewriter struct {
	ctx *RewriteContext
}

func (r *rewriter) rewrite(node tree.Expression) Expression {
	switch expr := node.(type) {
	case *tree.FunctionCall:
		return r.rewriteCall(expr)
	case *tree.Identifier:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, Scalar)
	case *tree.StringLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, Scalar)
	case *tree.FloatLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, Scalar)
	case *tree.LongLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, Scalar)
	case *tree.Constant:
		return NewConstant(r.ctx.EvalContext, expr.Value, Scalar)
	case *tree.SymbolReference:
		// FIXME: add check,index not found
		_, index, _ := lo.FindIndexOf(r.ctx.SourceLayout, func(item *plan.Symbol) bool {
			return item.Name == expr.Name
		})
		return NewColumn(r.ctx.EvalContext, expr.Name, index, Array)
	case *tree.Cast:
		return NewCast(r.ctx.EvalContext, expr.Type, r.rewrite(expr.Expression))
	case *tree.SubscriptExpression:
		return NewMapAccess(r.ctx.EvalContext, r.rewrite(expr.Base), r.rewrite(expr.Key))
	case *tree.ComparisonExpression:
		return NewComparison(r.ctx.EvalContext, expr.Operator, r.rewrite(expr.Left), r.rewrite(expr.Right))
	case *tree.LogicalExpression:
		terms := make([]Expression, len(expr.Terms))
		for i, term := range expr.Terms {
			terms[i] = r.rewrite(term)
		}
		return NewLogical(r.ctx.EvalContext, expr.Operator, terms)
	case *tree.NotExpression:
		return NewNot(r.ctx.EvalContext, r.rewrite(expr.Value))
	case *tree.InPredicate:
		return r.rewriteInPredicate(expr)
	default:
		panic(fmt.Sprintf("expression rewrite unimplemented: %T", node))
	}
}

func (r *rewriter) rewriteInPredicate(node *tree.InPredicate) Expression {
	value := r.rewrite(node.Value)
	var candidates []Expression
	if inList, ok := node.ValueList.(*tree.InListExpression); ok {
		candidates = make([]Expression, len(inList.Values))
		for i, v := range inList.Values {
			candidates[i] = r.rewrite(v)
		}
	}
	return NewIn(r.ctx.EvalContext, value, candidates)
}

func (r *rewriter) rewriteCall(node *tree.FunctionCall) Expression {
	rt := Scalar
	args := lo.Map(node.Arguments,
		func(item tree.Expression, index int) Expression {
			arg := r.rewrite(item)
			if arg.ResultType() == Array {
				rt = Array
			}
			return arg
		},
	)
	return NewScalarFunc(r.ctx.EvalContext, node.Name, rt, args)
}

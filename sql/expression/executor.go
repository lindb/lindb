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

	"github.com/lindb/lindb/spi/types"
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
		return NewConstant(r.ctx.EvalContext, expr.Value, types.DTString)
	case *tree.StringLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, types.DTString)
	case *tree.FloatLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, types.DTFloat)
	case *tree.LongLiteral:
		// TODO: right?
		return NewConstant(r.ctx.EvalContext, expr.Value, types.DTInt)
	case *tree.Constant:
		return NewConstant(r.ctx.EvalContext, expr.Value, expr.Type)
	case *tree.SymbolReference:
		// FIXME: add check,index not found
		_, index, _ := lo.FindIndexOf(r.ctx.SourceLayout, func(item *plan.Symbol) bool {
			return item.Name == expr.Name
		})
		// fmt.Printf("expr rewrite %v,%v,%v,%v\n", r.ctx.SourceLayout, expr.Name, ok, index)
		return NewColumn(r.ctx.EvalContext, expr.Name, index, expr.DataType)
	case *tree.Cast:
		return NewCast(r.ctx.EvalContext, expr.Type, r.rewrite(expr.Expression))
	default:
		panic(fmt.Sprintf("expression rewrite unimplemented: %T", node))
	}
}

func (r *rewriter) rewriteCall(node *tree.FunctionCall) Expression {
	scalarFunc, err := NewScalarFunc(r.ctx.EvalContext, node.Name, node.RetType, lo.Map(node.Arguments,
		func(item tree.Expression, index int) Expression {
			return r.rewrite(item)
		},
	))
	if err != nil {
		panic(err)
	}
	return scalarFunc
}

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

package optimization

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type ExpressionRewrite struct {
	originSymbol, newSymbol string
}

func (r *ExpressionRewrite) Rewrite(originSymbol, newSymbol string, e tree.Expression) tree.Expression {
	// reset origin and new symbol
	r.originSymbol = originSymbol
	r.newSymbol = newSymbol
	return e.Accept(nil, r).(tree.Expression)
}

func (r *ExpressionRewrite) Visit(context any, e tree.Node) any {
	fmt.Printf("expression rewrite=%T\n", e)
	switch expr := e.(type) {
	case *tree.ComparisonExpression:
		return &tree.ComparisonExpression{
			BaseNode: tree.BaseNode{ID: expr.ID},
			Operator: expr.Operator,
			Left:     expr.Left.Accept(context, r).(tree.Expression),
			Right:    expr.Right.Accept(context, r).(tree.Expression),
		}
	case *tree.SymbolReference:
		if expr.Name == r.originSymbol {
			return &tree.SymbolReference{
				BaseNode: tree.BaseNode{ID: expr.ID},
				Name:     r.newSymbol,
			}
		}
		return expr
	default:
		return e
	}
}

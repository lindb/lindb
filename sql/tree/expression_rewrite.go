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

package tree

import (
	"fmt"
)

type ExpressionRewriter interface {
	RewriteExpression(context any, node Expression) Expression
}

type ExpressionTreeRewriter struct {
	rewriter ExpressionRewriter
	visitor  Visitor
}

func NewExpressionTreeRewriter(rewriter ExpressionRewriter) *ExpressionTreeRewriter {
	return &ExpressionTreeRewriter{
		rewriter: rewriter,
		visitor:  NewExpressionRewriteVisitor(rewriter),
	}
}

func (etr *ExpressionTreeRewriter) rewrite(context any, node Expression) Expression {
	return etr.visitor.Visit(context, node).(Expression)
}

func RewriteExpression(context any, rewriter ExpressionRewriter, node Expression) Expression {
	return NewExpressionTreeRewriter(rewriter).rewrite(context, node)
}

type ExpressionRewriteVisitor struct {
	rewriter ExpressionRewriter
}

func NewExpressionRewriteVisitor(rewriter ExpressionRewriter) *ExpressionRewriteVisitor {
	return &ExpressionRewriteVisitor{
		rewriter: rewriter,
	}
}

func (v *ExpressionRewriteVisitor) Visit(context any, n Node) any {
	// TODO: default????
	if expr, ok := n.(Expression); ok {
		result := v.rewriter.RewriteExpression(context, expr)
		if result != nil {
			return result
		}
	}
	panic(fmt.Sprintf("expression rewrite not support: %T", n))
}

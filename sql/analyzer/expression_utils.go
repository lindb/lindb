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

package analyzer

import (
	"github.com/lindb/lindb/sql/tree"
)

func asQualifiedName(expression tree.Expression) (name *tree.QualifiedName) {
	switch e := expression.(type) {
	case *tree.Identifier:
		name = tree.NewQualifiedName([]*tree.Identifier{e})
	case *tree.DereferenceExpression:
		name = e.ToQualifiedName()
		// TODO:????
	}
	return
}

func ExtractConjuncts(expression tree.Expression) (result []tree.Expression) {
	return ExtractPredicates(tree.LogicalAND, expression, result)
}

func ExtractTimePredicates(expression tree.Expression) (result []*tree.TimePredicate, newExpr tree.Expression) {
	if logicalExpression, ok := expression.(*tree.LogicalExpression); ok {
		var newTerms []tree.Expression
		for _, term := range logicalExpression.Terms {
			if timePredicate, ok := term.(*tree.TimePredicate); ok {
				result = append(result, timePredicate)
			} else {
				subResult, newTerm := ExtractTimePredicates(term)
				if len(subResult) > 0 && len(result) > 0 {
					panic("time predicate not support nested")
				}
				result = subResult
				if newTerm != nil {
					newTerms = append(newTerms, term)
				}
			}
		}
		if len(newTerms) == 0 {
			return result, nil
		}
		if len(newTerms) != len(logicalExpression.Terms) {
			logicalExpression.Terms = newTerms
			newExpr = logicalExpression
		} else {
			newExpr = logicalExpression
		}
	} else {
		newExpr = expression
	}
	return
}

func ExtractPredicates(operator tree.LogicalOperator,
	expression tree.Expression, result []tree.Expression,
) (r []tree.Expression) {
	if logicalExpression, ok := expression.(*tree.LogicalExpression); ok && logicalExpression.Operator == operator {
		for i := range logicalExpression.Terms {
			term := logicalExpression.Terms[i]
			result = ExtractPredicates(operator, term, result)
		}
	} else {
		result = append(result, expression)
	}
	return result
}

func ExtractAggregationFunctions(nodes []tree.Expression, handle func(node tree.Node)) {
	visitor := &tree.DefaultTraversalVisitor{
		PreProcess: handle,
	}
	for _, node := range nodes {
		visitor.Visit(nil, node)
	}
}

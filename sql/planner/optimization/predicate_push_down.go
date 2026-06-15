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
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PredicatePushDown struct {
	ctx     *context.PlannerContext
	rewrite *ExpressionRewrite
}

func NewPredicatePushDown() PlanOptimizer {
	return &PredicatePushDown{
		rewrite: &ExpressionRewrite{},
	}
}

// Optimize implements PlanOptimizer.
func (p *PredicatePushDown) Optimize(ctx *context.PlannerContext, node plan.PlanNode) plan.PlanNode {
	p.ctx = ctx
	return node.Accept(nil, p).(plan.PlanNode)
}

func (p *PredicatePushDown) Visit(context any, n plan.PlanNode) (r any) {
	switch node := n.(type) {
	case *plan.FilterNode:
		return p.visitFilter(context, node)
	case *plan.JoinNode:
		return p.visitJoin(context, node)
	case *plan.TableScanNode:
		return p.visitTableScan(context, node)
	}
	var sources []plan.PlanNode
	for _, source := range n.GetSources() {
		sources = append(sources, source.Accept(context, p).(plan.PlanNode))
	}
	return plan.ReplaceChildren(n, sources)
}

func (p *PredicatePushDown) visitTableScan(context any, node *plan.TableScanNode) any {
	expressions, ok := context.([]tree.Expression)
	if !ok {
		return node
	}
	return &plan.FilterNode{
		BaseNode:  plan.BaseNode{ID: p.ctx.PlanNodeIDAllocator.Next()},
		Source:    node,
		Predicate: p.createLogicalExpression(expressions),
	}
}

func (p *PredicatePushDown) visitJoin(context any, node *plan.JoinNode) any {
	expressions, ok := context.([]tree.Expression)
	if !ok {
		return node
	}

	var criteria []*plan.EqualJoinCriteria
	joinPredicates := make(map[string]string)
	var filterPredicates []tree.Expression
	// Classify each WHERE conjunct: equality comparisons between symbols from
	// different sides become equi-join criteria; everything else is a filter.
	for _, expr := range expressions {
		foundJoinPredicate := false
		// Only EQ comparisons can be promoted to hash-join keys; non-equality
		// comparisons (>, <, <>, etc.) must remain as post-join filters.
		if comparison, ok := expr.(*tree.ComparisonExpression); ok &&
			comparison.Operator == tree.ComparisonEQ {
			symbols := plan.ExtractSymbolsFromExpression(comparison)
			if len(symbols) == 2 {
				criteria = append(criteria, &plan.EqualJoinCriteria{Left: symbols[0], Right: symbols[1]})

				joinPredicates[symbols[0].Name] = symbols[1].Name
				joinPredicates[symbols[1].Name] = symbols[0].Name
				foundJoinPredicate = true
			}
		}

		if !foundJoinPredicate {
			filterPredicates = append(filterPredicates, expr)
		}
	}

	var leftPredicates []tree.Expression
	var rightPredicates []tree.Expression
	var filterExpressions []tree.Expression
	leftScope := node.Left.GetOutputSymbols()
	rightScope := node.Right.GetOutputSymbols()
	for _, expr := range filterPredicates {
		symbols := plan.ExtractSymbolsFromExpression(expr)
		if len(symbols) == 2 || len(symbols) == 0 {
			continue
		}
		symbolName := symbols[0].Name
		refName, ok := joinPredicates[symbolName]
		if ok {
			leftPredicates = append(leftPredicates, p.rewriteExpression(symbolName, refName, leftScope, expr))
			rightPredicates = append(rightPredicates, p.rewriteExpression(symbolName, refName, rightScope, expr))
		} else {
			filterExpressions = append(filterExpressions, expr)
		}
	}

	// Append WHERE-extracted equi-join criteria to the ON criteria already set by
	// the logical planner. Using assignment (=) would overwrite ON conditions.
	node.Criteria = append(node.Criteria, criteria...)

	node.Left = node.Left.Accept(leftPredicates, p).(plan.PlanNode)
	node.Right = node.Right.Accept(rightPredicates, p).(plan.PlanNode)

	if len(filterExpressions) > 0 {
		// if has filter expression after predicate push down, then create a filter node
		return &plan.FilterNode{
			BaseNode:  plan.BaseNode{ID: p.ctx.PlanNodeIDAllocator.Next()},
			Predicate: p.createLogicalExpression(filterExpressions),
			Source:    node,
		}
	}

	return node
}

func (p *PredicatePushDown) visitFilter(_ any, node *plan.FilterNode) any {
	// HAVING predicates sit directly above an AggregationNode and reference
	// aggregate output symbols — they must NOT be pushed below the aggregation.
	if _, ok := node.Source.(*plan.AggregationNode); ok {
		optimizedSource := node.Source.Accept(nil, p).(plan.PlanNode)
		return node.ReplaceChildren([]plan.PlanNode{optimizedSource})
	}
	// WHERE-style filter: push conjuncts down toward the TableScan.
	return node.Source.Accept(tree.ExtractConjuncts(node.Predicate), p)
}

func (p *PredicatePushDown) rewriteExpression(a, b string, scope []*plan.Symbol, expr tree.Expression) tree.Expression {
	_, ok := lo.Find(scope, func(item *plan.Symbol) bool {
		return item.Name == a
	})
	if ok {
		return expr
	}
	return p.rewrite.Rewrite(a, b, expr)
}

func (p *PredicatePushDown) createLogicalExpression(expressions []tree.Expression) tree.Expression {
	if len(expressions) == 1 {
		return expressions[0]
	}
	return &tree.LogicalExpression{
		BaseNode: tree.BaseNode{
			ID: p.ctx.AnalyzerContext.IDAllocator.Next(),
		},
		Operator: tree.LogicalAND,
		Terms:    expressions,
	}
}

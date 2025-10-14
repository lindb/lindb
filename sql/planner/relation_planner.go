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

package planner

import (
	"fmt"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/expression"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type RelationPlanner struct {
	outerContext     *TranslationMap
	context          *context.PlannerContext
	timePredicates   []*tree.TimePredicate
	groupingInterval *tree.IntervalLiteral
}

func NewRelationPlanner(ctx *context.PlannerContext, outerContext *TranslationMap,
	timePredicates []*tree.TimePredicate, groupingInterval *tree.IntervalLiteral,
) tree.Visitor {
	return &RelationPlanner{
		context:          ctx,
		outerContext:     outerContext,
		timePredicates:   timePredicates,
		groupingInterval: groupingInterval,
	}
}

func (p *RelationPlanner) Visit(context any, n tree.Node) (r any) {
	switch node := n.(type) {
	case *tree.Query:
		return p.visitQuery(context, node)
	case *tree.QuerySpecification:
		return p.visitQuerySpecification(context, node)
	case *tree.Join:
		return p.visitJoin(context, node)
	case *tree.AliasedRelation:
		return p.visitAliasedRelation(context, node)
	case *tree.Table:
		return p.visitTable(context, node)
	case *tree.Values:
		return p.visitValues(context, node)
	default:
		panic(fmt.Sprintf("relation analyzer unsupport node:%T", n))
	}
}

func (p *RelationPlanner) visitQuery(_ any, node *tree.Query) (r any) {
	return NewQueryPlanner(p.context, p.outerContext).planQuery(node)
}

func (p *RelationPlanner) visitQuerySpecification(_ any, node *tree.QuerySpecification) (r any) {
	return NewQueryPlanner(p.context, p.outerContext).planQuerySpecification(node)
}

func (p *RelationPlanner) visitJoin(context any, node *tree.Join) (r any) {
	leftPlan := node.Left.Accept(context, p).(*RelationPlan)
	rightPlan := node.Right.Accept(context, p).(*RelationPlan)
	criteria := node.Criteria
	if _, ok := criteria.(*tree.JoinUsing); ok {
		return p.planJoinUsing(node, leftPlan, rightPlan)
	}
	return p.planJoin(node, p.context.AnalyzerContext.Analysis.GetScope(node), leftPlan, rightPlan)
}

func (p *RelationPlanner) visitAliasedRelation(context any, node *tree.AliasedRelation) (r any) {
	subPlan := node.Relation.Accept(context, p).(*RelationPlan)
	root := subPlan.Root
	mappings := subPlan.FieldMappings
	return &RelationPlan{
		Root:          root,
		OutContext:    p.outerContext,
		Scope:         p.context.AnalyzerContext.Analysis.GetScope(node),
		FieldMappings: mappings,
	}
}

func (p *RelationPlanner) visitValues(_ any, node *tree.Values) (r any) {
	scope := p.context.AnalyzerContext.Analysis.GetScope(node)
	var outputSymbols []*planpkg.Symbol
	for _, f := range scope.RelationType.Fields {
		symbol := &planpkg.Symbol{
			Name:     f.Name,
			DataType: f.DataType,
			Hidden:   f.Hidden,
		}
		outputSymbols = append(outputSymbols, symbol)
	}
	return &RelationPlan{
		Root: &planpkg.ValuesNode{
			BaseNode: planpkg.BaseNode{
				ID: p.context.PlanNodeIDAllocator.Next(),
			},
			Rows:          node.Rows,
			RowCount:      node.Rows.NumRows(),
			OutputSymbols: outputSymbols,
		},
		OutContext:    p.outerContext,
		Scope:         scope,
		FieldMappings: outputSymbols,
	}
}

func (p *RelationPlanner) visitTable(_ any, node *tree.Table) (r any) {
	fmt.Printf("visit table, time predicates=%v\n", p.timePredicates)
	namedQuery := p.context.AnalyzerContext.Analysis.GetNamedQuery(node)
	scope := p.context.AnalyzerContext.Analysis.GetScope(node)
	var plan *RelationPlan
	if namedQuery != nil {
		// process named query ref
		subPlan := namedQuery.Accept(nil, p).(*RelationPlan)
		// FIXME:???
		coerced := coerce(subPlan, nil, nil, nil)
		plan = &RelationPlan{
			Root:          coerced.Node,
			Scope:         scope,
			FieldMappings: coerced.Fields,
		}
	} else {
		var outputSymbols []*planpkg.Symbol
		columnMapping := make(map[string]string)
		for _, f := range scope.RelationType.Fields {
			symbol := p.context.SymbolAllocator.NewSymbol(f.Name, f.DataType, f.Hidden)
			outputSymbols = append(outputSymbols, symbol)
			if symbol.Name != f.Name {
				columnMapping[symbol.Name] = f.Name
			}
		}

		fmt.Printf("table visit relation plan====%v\n", outputSymbols)
		tableHandle := p.context.AnalyzerContext.Analysis.GetTableHandle(node)
		tableMetadata := p.context.AnalyzerContext.Analysis.GetTableMetadata(tableHandle.String())
		root := planpkg.NewTableScanNode(p.context.PlanNodeIDAllocator.Next())
		root.Table = p.context.AnalyzerContext.Analysis.GetTableHandle(node)
		// time range(statement>context>default)
		// 1. default query time range(last hour)
		timeRange := timeutil.TimeRange{
			Start: time.Now().UnixMilli() - time.Hour.Milliseconds(),
			End:   time.Now().UnixMilli(),
		}
		fmt.Printf("default time range:%v\n", timeRange)
		// 2. time range from context
		currentParams := p.context.Context.Value(constants.ContextKeyParams)
		if currentParams != nil {
			if params, ok := currentParams.(*models.ExecuteParam); ok {
				if params.TimeRange.Start > 0 {
					timeRange.Start = params.TimeRange.Start
				}
				if params.TimeRange.End > 0 {
					timeRange.End = params.TimeRange.End
				}
			}
			fmt.Printf("params time range:%v\n", timeRange)
		}
		// 3. time range from statement condition
		if len(p.timePredicates) > 0 {
			translations := &TranslationMap{context: p.context}
			evalCtx := expression.NewEvalContext(p.context.Context)
			// if has time predicate, add time range filter for table handle
			for _, timePredicate := range p.timePredicates {
				timestamp, _ := expression.EvalTime(evalCtx, translations.Rewrite(timePredicate.Value))
				switch timePredicate.Operator {
				case tree.ComparisonGT:
					timeRange.Start = timestamp.UnixMilli()
				case tree.ComparisonLT:
					timeRange.End = timestamp.UnixMilli()
				}
			}
		}
		fmt.Printf("set time range:%v\n", timeRange)
		root.Table.SetTimeRange(timeRange)

		if p.groupingInterval != nil {
			root.Table.SetInterval(p.groupingInterval.Interval())
		}

		root.OutputSymbols = outputSymbols
		root.Partitions = tableMetadata.Partitions
		root.ColumnMapping = columnMapping

		plan = &RelationPlan{
			Root:          root,
			Scope:         scope,
			FieldMappings: outputSymbols,
			OutContext:    p.outerContext,
		}
	}
	return plan
}

func (p *RelationPlanner) planJoinUsing(node *tree.Join, left, right *RelationPlan) *RelationPlan {
	panic("need implement join using")
}

func (p *RelationPlanner) planJoin(node *tree.Join, scope *analyzer.Scope, left, right *RelationPlan) *RelationPlan {
	var outputSymbols []*planpkg.Symbol
	outputSymbols = append(outputSymbols, left.FieldMappings...)
	outputSymbols = append(outputSymbols, right.FieldMappings...)

	var joinCriteriaClauses []*planpkg.EqualJoinCriteria
	leftPlanBuilder := newPlanBuilder(p.context, left, nil)
	rightPlanBuilder := newPlanBuilder(p.context, right, nil)
	fmt.Printf("join type===%v\n", node.Type)
	if node.Type != tree.CROSS && node.Type != tree.IMPLICIT {
		criteria := p.context.AnalyzerContext.Analysis.GetJoinCriteria(node)
		expressions := tree.ExtractConjuncts(criteria)
		var leftComparisonExpressions []tree.Expression
		var rightComparisonExpressions []tree.Expression
		var joinConditionComparisonOperators []tree.ComparisonOperator
		for i := range expressions {
			conjunct := expressions[i]

			if comparisonExpression, ok := conjunct.(*tree.ComparisonExpression); ok {
				firstExpression := comparisonExpression.Left
				secondExpression := comparisonExpression.Right
				leftComparisonExpressions = append(leftComparisonExpressions, firstExpression)
				rightComparisonExpressions = append(rightComparisonExpressions, secondExpression)
				joinConditionComparisonOperators = append(joinConditionComparisonOperators, comparisonExpression.Operator)
			}
			// TODO: check not equal

			fmt.Println(conjunct)
		}

		// add projections for join criteria
		leftPlanBuilder = leftPlanBuilder.appendProjections(leftComparisonExpressions)
		rightPlanBuilder = rightPlanBuilder.appendProjections(rightComparisonExpressions)

		leftCoercions := coerceExpressions(leftPlanBuilder, leftComparisonExpressions,
			p.context.AnalyzerContext.Analysis,
			p.context.SymbolAllocator, p.context.PlanNodeIDAllocator)
		rightCoercions := coerceExpressions(rightPlanBuilder, rightComparisonExpressions,
			p.context.AnalyzerContext.Analysis,
			p.context.SymbolAllocator, p.context.PlanNodeIDAllocator)
		fmt.Printf("join......%v\n", leftCoercions)
		for i := range leftComparisonExpressions {
			if joinConditionComparisonOperators[i] == tree.ComparisonEQ {
				leftSymbol := leftCoercions.mappings[leftComparisonExpressions[i]]
				rightSymbol := rightCoercions.mappings[rightComparisonExpressions[i]]
				joinCriteriaClauses = append(joinCriteriaClauses, &planpkg.EqualJoinCriteria{
					Left:  leftSymbol,
					Right: rightSymbol,
				})
			}
		}
	}

	root := &planpkg.JoinNode{
		BaseNode: planpkg.BaseNode{
			ID: p.context.PlanNodeIDAllocator.Next(),
		},
		Type:               planpkg.JoinTypeConvert(node.Type),
		Left:               leftPlanBuilder.root,
		Right:              rightPlanBuilder.root,
		Criteria:           joinCriteriaClauses,
		LeftOutputSymbols:  leftPlanBuilder.root.GetOutputSymbols(),
		RightOutputSymbols: rightPlanBuilder.root.GetOutputSymbols(),
	}
	return &RelationPlan{
		Root:          root,
		Scope:         scope,
		FieldMappings: outputSymbols,
		OutContext:    p.outerContext,
	}
}

func coerce(plan *RelationPlan, types []types.Type,
	symbolAllocator *planpkg.SymbolAllocator, idAllocator *planpkg.PlanNodeIDAllocator,
) *NodeAndMappings {
	return nil
}

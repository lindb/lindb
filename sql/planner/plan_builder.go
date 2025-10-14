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

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PlanBuilder struct {
	root         plan.PlanNode
	translations *TranslationMap
}

func newPlanBuilder(context *context.PlannerContext,
	plan *RelationPlan, mappings map[tree.NodeID]*plan.Symbol, //nolint
) *PlanBuilder {
	// TODO: check mappings if nil(remove nolint)
	fmt.Printf("new plan builder=%v\n", plan.FieldMappings)
	return &PlanBuilder{
		root: plan.Root,
		translations: &TranslationMap{
			scope:        plan.Scope,
			context:      context,
			outerContext: plan.OutContext,
			fieldSymbols: plan.FieldMappings,
			astToSymbols: mappings,
		},
	}
}

func (pb *PlanBuilder) withNewRoot(root plan.PlanNode) *PlanBuilder {
	return &PlanBuilder{
		root:         root,
		translations: pb.translations,
	}
}

func (pb *PlanBuilder) appendProjections(expressions []tree.Expression) *PlanBuilder {
	var assignments plan.Assignments
	fmt.Printf("add root...%T\n", pb.root)
	assignments = assignments.Add(pb.root.GetOutputSymbols())
	symbolAllocator := pb.translations.context.SymbolAllocator
	idAllocator := pb.translations.context.PlanNodeIDAllocator

	mappings := make(map[tree.NodeID]*plan.Symbol)
	for i := range expressions {
		expression := expressions[i]
		fmt.Printf("check transs====%T,,,, %v=%v\n", expression, expression, pb.translations.CanTranslate(expression))
		if _, ok := mappings[expression.GetID()]; !ok && !pb.translations.CanTranslate(expression) {
			fmt.Println("kkkkkkkkkkkkk..............")
			symbol := symbolAllocator.FromExpression(expression, pb.translations.context.AnalyzerContext.Analysis.GetType(expression))
			expr := pb.translations.Rewrite(expression)
			assignments = append(assignments, &plan.Assignment{
				Symbol:     symbol,
				Expression: expr,
			})
			mappings[expression.GetID()] = symbol
			fmt.Println("kkkkkkkkkkkkk.............. done")
		}
	}
	fmt.Printf("project ass.......%v\n", string(encoding.JSONMarshal(assignments)))
	return &PlanBuilder{
		translations: pb.translations.withAdditionalMapping(mappings), // FIXME:
		root: &plan.ProjectionNode{
			BaseNode: plan.BaseNode{
				ID: idAllocator.Next(),
			},
			Source:      pb.root,
			Assignments: assignments,
		},
	}
}

func (pb *PlanBuilder) translate(node tree.Expression) *plan.Symbol {
	return plan.SymbolFrom(pb.translations.Rewrite(node))
}

func (pb *PlanBuilder) rewrite(node tree.Expression) tree.Expression {
	return pb.translations.Rewrite(node)
}

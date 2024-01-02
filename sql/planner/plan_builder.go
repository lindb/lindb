package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PlanBuilder struct {
	root         plan.PlanNode
	translations *TranslationMap
}

func newPlanBuilder(context *context.PlannerContext, plan *RelationPlan, mappings map[tree.NodeID]*plan.Symbol) *PlanBuilder {
	return &PlanBuilder{
		root: plan.Root,
		translations: &TranslationMap{
			scope:        plan.Scope,
			context:      context,
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
	assignments := make(plan.Assignments)
	fmt.Printf("add root...%T\n", pb.root)
	assignments.Add(pb.root.GetOutputSymbols())
	symbolAllocator := pb.translations.context.SymbolAllocator
	idAllocator := pb.translations.context.PlanNodeIDAllocator

	mappings := make(map[tree.NodeID]*plan.Symbol)
	for i := range expressions {
		expression := expressions[i]
		if _, ok := mappings[expression.GetID()]; !ok && !pb.translations.CanTranslate(expression) {
			fmt.Println("kkkkkkkkkkkkk..............")
			symbol := symbolAllocator.NewSymbol(expression, "")
			assignments[symbol] = pb.translations.Rewrite(expression)
			mappings[expression.GetID()] = symbol
			fmt.Println("kkkkkkkkkkkkk.............. done")
		}
	}
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

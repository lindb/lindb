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
			symbol := symbolAllocator.NewSymbol(expression, "", pb.translations.context.AnalyzerContext.Analysis.GetType(expression))
			expr := pb.translations.Rewrite(expression)
			assignments = append(assignments, &plan.Assignment{
				Symbol:     symbol,
				Expression: expr,
			})
			mappings[expression.GetID()] = symbol
			fmt.Println("kkkkkkkkkkkkk.............. done")
		}
	}
	fmt.Printf("proejct ass.......%v\n", assignments)
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

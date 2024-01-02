package context

import (
	"context"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PlannerContext struct {
	Context             context.Context
	AnalyzerContext     *analyzer.AnalyzerContext
	PlanNodeIDAllocator *plan.PlanNodeIDAllocator
	SymbolAllocator     *plan.SymbolAllocator
	Database            string
}

func NewPlannerContext(ctx context.Context, idallocator *tree.NodeIDAllocator, stmt tree.Statement) *PlannerContext {
	return &PlannerContext{
		Context:             ctx,
		AnalyzerContext:     analyzer.NewContext(stmt, idallocator),
		PlanNodeIDAllocator: plan.NewPlanNodeIDAllocator(),
		SymbolAllocator:     plan.NewSymbolAllocator(),
	}
}

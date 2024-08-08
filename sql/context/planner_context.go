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

func NewPlannerContext(ctx context.Context, database string, idAllocator *tree.NodeIDAllocator, stmt tree.Statement) *PlannerContext {
	analyzerContext := analyzer.NewAnalyzerContext(database, stmt, idAllocator)
	return &PlannerContext{
		Context:             ctx,
		Database:            database,
		AnalyzerContext:     analyzerContext,
		PlanNodeIDAllocator: plan.NewPlanNodeIDAllocator(),
		SymbolAllocator:     plan.NewSymbolAllocator(analyzerContext),
	}
}

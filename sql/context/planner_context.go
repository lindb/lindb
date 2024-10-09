package context

import (
	"context"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PlannerContext struct {
	Context             context.Context
	PlanNodeIDAllocator *plan.PlanNodeIDAllocator
	// create in analyze phase
	AnalyzerContext *analyzer.AnalyzerContext
	SymbolAllocator *plan.SymbolAllocator

	Database string
}

func NewPlannerContext(ctx context.Context, database string, idAllocator *tree.NodeIDAllocator, stmt tree.Statement) *PlannerContext {
	return &PlannerContext{
		Context:             ctx,
		Database:            database,
		PlanNodeIDAllocator: plan.NewPlanNodeIDAllocator(),
	}
}

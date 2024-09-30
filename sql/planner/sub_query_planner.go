package planner

import (
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/tree"
)

type SubQueryPlanner struct {
	context *context.PlannerContext
}

func NewSubQueryPlanner(context *context.PlannerContext) *SubQueryPlanner {
	return &SubQueryPlanner{
		context: context,
	}
}

func (p *SubQueryPlanner) handleSubQueries(builder *PlanBuilder, expression tree.Expression, subQueries *analyzer.SubQueryAnalysis) *PlanBuilder {
	return builder
}

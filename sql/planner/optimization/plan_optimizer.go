package optimization

import (
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PlanOptimizer interface {
	Optimize(ctx *context.PlannerContext, node plan.PlanNode) plan.PlanNode
}

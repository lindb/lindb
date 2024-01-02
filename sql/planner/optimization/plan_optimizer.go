package optimization

import "github.com/lindb/lindb/sql/planner/plan"

type PlanOptimizer interface {
	Optimize(node plan.PlanNode, idAllocator *plan.PlanNodeIDAllocator) plan.PlanNode
}

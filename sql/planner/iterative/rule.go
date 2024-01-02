package iterative

import "github.com/lindb/lindb/sql/planner/plan"

type Rule interface {
	Apply(context *Context, node plan.PlanNode) plan.PlanNode
}

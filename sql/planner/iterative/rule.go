package iterative

import (
	"github.com/lindb/lindb/sql/planner/plan"
)

type Rule interface {
	// GetPattern returns a pattern to which plan nodes this rule applies.
	// TODO: remove
	// GetPattern() *matching.Pattern
	Apply(context *Context, node plan.PlanNode) plan.PlanNode
}

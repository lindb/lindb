package iterative

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/plan"
)

type Rule interface {
	// GetPattern returns a pattern to which plan nodes this rule applies.
	// TODO: remove
	// GetPattern() *matching.Pattern
	Apply(context *Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode
}

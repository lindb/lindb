package rule

import (
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

// Base represents optimization base rule for common logic.
type Base[N plan.PlanNode] struct {
	// apply apply rule for specific node.
	apply func(context *iterative.Context, node N) plan.PlanNode
}

// Apply applies optimization rule for specific node.
func (rule *Base[N]) Apply(context *iterative.Context, node plan.PlanNode) plan.PlanNode {
	// check node if match target node type
	if targetNode, ok := node.(N); ok {
		return rule.apply(context, targetNode)
	}
	return nil
}

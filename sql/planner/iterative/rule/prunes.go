package rule

import (
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

// PruneOutputSourceColumns represents optimization rule for pruning output node's source columns.
type PruneOutputSourceColumns struct {
	Base[*plan.OutputNode]
}

// NewPruneOutputSourceColumns creates a PruneOutputSourceColumns instance.
func NewPruneOutputSourceColumns() iterative.Rule {
	rule := &PruneOutputSourceColumns{}
	rule.apply = func(context *iterative.Context, node *plan.OutputNode) plan.PlanNode {
		return restrictChildOutputs(context.IDAllocator, node, node.GetOutputSymbols())
	}
	return rule
}

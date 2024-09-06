package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PruneOutputSourceColumns struct {
	Base[*plan.OutputNode]
}

func NewPruneOutputSourceColumns() iterative.Rule {
	rule := &PruneOutputSourceColumns{}
	rule.apply = func(context *iterative.Context, captures *matching.Captures, node *plan.OutputNode) plan.PlanNode {
		return restrictChildOutputs(context.IDAllocator, node, node.GetOutputSymbols())
	}
	return rule
}

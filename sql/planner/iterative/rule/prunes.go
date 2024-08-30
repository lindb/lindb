package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PruneOutputSourceColumns struct{}

func NewPruneOutputSourceColumns() iterative.Rule {
	return &PruneOutputSourceColumns{}
}

func (rule *PruneOutputSourceColumns) GetPattern() *matching.Pattern {
	return output()
}

func (rule *PruneOutputSourceColumns) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if output, ok := node.(*plan.OutputNode); ok {
		return restrictChildOutputs(context.IDAllocator, output, node.GetOutputSymbols())
	}
	return nil
}

package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PruneOutputSourceColumnsRule struct{}

func NewPruneOutputSourceColumnsRule() iterative.Rule {
	return &PruneOutputSourceColumnsRule{}
}

func (rule *PruneOutputSourceColumnsRule) GetPattern() *matching.Pattern {
	return output()
}

func (rule *PruneOutputSourceColumnsRule) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if output, ok := node.(*plan.OutputNode); ok {
		return restrictChildOutputs(context.IDAllocator, output, node.GetOutputSymbols())
	}
	return nil
}

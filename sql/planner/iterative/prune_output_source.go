package iterative

import "github.com/lindb/lindb/sql/planner/plan"

type PruneOutputSourceColumnsRule struct{}

func NewPruneOutputSourceColumnsRule() Rule {
	return &PruneOutputSourceColumnsRule{}
}

func (rule *PruneOutputSourceColumnsRule) Apply(context *Context, node plan.PlanNode) plan.PlanNode {
	if output, ok := node.(*plan.OutputNode); ok {
		return restrictChildOutputs(context.idAllocator, output, node.GetOutputSymbols())
	}
	return nil
}

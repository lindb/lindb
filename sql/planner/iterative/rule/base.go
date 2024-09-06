package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type Base[N plan.PlanNode] struct {
	apply func(context *iterative.Context, captures *matching.Captures, node N) plan.PlanNode
}

func (rule *Base[N]) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if targetNode, ok := node.(N); ok {
		return rule.apply(context, captures, targetNode)
	}
	return nil
}

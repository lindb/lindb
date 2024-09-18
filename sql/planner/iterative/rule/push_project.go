package rule

import (
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PushProjectionIntoTableScan struct {
	Base[*plan.ProjectionNode]
}

func NewPushProjectionIntoTableScan() iterative.Rule {
	rule := &PushProjectionIntoTableScan{}
	rule.apply = func(context *iterative.Context, node *plan.ProjectionNode) plan.PlanNode {
		return nil
	}
	return rule
}

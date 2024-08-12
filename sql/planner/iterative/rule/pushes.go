package rule

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PushPartialAggregationThroughExchange struct{}

func NewPushPartialAggregationThroughExchange() *PushPartialAggregationThroughExchange {
	return &PushPartialAggregationThroughExchange{}
}

func (rule *PushPartialAggregationThroughExchange) Apply(context *iterative.Context, node plan.PlanNode) plan.PlanNode {
	if aggregationNode, ok := node.(*plan.AggregationNode); ok {
		fmt.Println(aggregationNode)
		return nil
	}
	return nil
}

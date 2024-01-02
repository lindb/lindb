package planner

import (
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type NodeAndMappings struct {
	Node   plan.PlanNode
	Fields []*plan.Symbol
}

type PlanAndMappings struct {
	subPlan  *PlanBuilder
	mappings map[tree.Expression]*plan.Symbol
}

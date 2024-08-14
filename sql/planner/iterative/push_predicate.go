package iterative

import "github.com/lindb/lindb/sql/planner/plan"

func PushFilterIntoTableScan(filter *plan.FilterNode, node *plan.TableScanNode) plan.PlanNode {
	tableScan := &plan.TableScanNode{
		BaseNode: plan.BaseNode{
			ID: node.GetNodeID(),
		},
		Table: node.Table,
	}
	return &plan.FilterNode{
		BaseNode: plan.BaseNode{
			ID: filter.GetNodeID(),
		},
		Source:    tableScan,
		Predicate: filter.Predicate,
	}
}

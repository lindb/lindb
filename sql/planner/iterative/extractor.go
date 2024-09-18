package iterative

import (
	"github.com/lindb/lindb/sql/planner/plan"
)

// ExtractTableScan extracts table scan from the sources of node.
func ExtractTableScan(context *Context, node plan.PlanNode) (tableScan *plan.TableScanNode) {
	visitor := &plan.DefaultTraversalVisitor{
		Process: func(n plan.PlanNode) {
			if node, ok := n.(*plan.TableScanNode); ok {
				tableScan = node
			}
		},
		Resolve: context.Lookup.Resolve,
	}

	// visit all plan node
	_ = visitor.Visit(nil, node)
	return
}

package execution

import (
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/iterative/rule"
	"github.com/lindb/lindb/sql/planner/optimization"
)

func planOptimizers() []optimization.PlanOptimizer {
	return []optimization.PlanOptimizer{
		// optimization.NewPruneColumns(),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
		// column pruning optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPruneAggregationSourceColumns(),
			rule.NewPruneFilterColumns(),
			rule.NewPruneOutputSourceColumns(),
			rule.NewPruneProjectionColumns(),
			rule.NewPruneTableScanColumns(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
		// push into table scan optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushProjectionIntoTableScan(),
			rule.NewPushAggregationIntoTableScan(),
		}),
		optimization.NewAddExchanges(),
		optimization.NewAddLocalExchanges(),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushPartialAggregationThroughExchange(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
	}
}

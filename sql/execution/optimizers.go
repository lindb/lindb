// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
			rule.NewPruneJoinColumns(),
			rule.NewPruneTableScanColumns(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
		// push into table scan optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushTimestampIntoTableScan(),
			rule.NewPushProjectionIntoTableScan(),
		}),
		optimization.NewPredicatePushDown(),
		optimization.NewAddExchanges(),
		optimization.NewAddLocalExchanges(),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushPartialAggregationThroughExchange(),
		}),
		// push into table scan optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushAggregationIntoTableScan(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
	}
}

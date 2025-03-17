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

package rule

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PushTimestampIntoTableScan struct {
	Base[*plan.OutputNode]
}

func NewPushTimestampIntoTableScan() iterative.Rule {
	rule := &PushTimestampIntoTableScan{}
	rule.apply = rule.pushTimestampIntoTableScan
	return rule
}

func (rule *PushTimestampIntoTableScan) pushTimestampIntoTableScan(context *iterative.Context, node *plan.OutputNode) plan.PlanNode {
	timestamp, ok := lo.Find(node.GetOutputSymbols(), func(item *plan.Symbol) bool {
		return item.DataType == types.DTTimestamp
	})
	if !ok {
		return nil
	}
	tableScan := iterative.ExtractTableScan(context, node)
	if tableScan == nil {
		return nil
	}
	if _, ok = lo.Find(tableScan.GetOutputSymbols(), func(item *plan.Symbol) bool {
		return item.DataType == types.DTTimestamp
	}); ok {
		return nil
	}
	tableScan.OutputSymbols = append(tableScan.OutputSymbols, timestamp)
	return nil
}

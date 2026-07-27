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
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

// RemoveHiddenColumns walks every node in the plan tree that owns an
// independent symbol list and strips hidden symbols from it.
// TableScanNode is deliberately skipped — the source connector relies on
// hidden symbols (e.g. the implicit timestamp column) being present there.
//
// Nodes that simply delegate GetOutputSymbols() to their child
// (FilterNode, InsertNode, …) need no special handling.
//
// The rule is anchored on OutputNode so it fires once per query root (or
// per fragment root in distributed plans).
type RemoveHiddenColumns struct {
	Base[*plan.OutputNode]
}

func NewRemoveHiddenColumns() iterative.Rule {
	r := &RemoveHiddenColumns{}
	r.apply = r.removeHidden
	return r
}

func (r *RemoveHiddenColumns) removeHidden(
	ctx *iterative.Context, node *plan.OutputNode,
) plan.PlanNode {
	visitor := &plan.DefaultTraversalVisitor{
		Process: pruneHiddenFromNode,
		Resolve: ctx.Lookup.Resolve,
	}
	_ = visitor.Visit(nil, node)
	// Return nil: all modifications are in-place (same pattern as PushTimestampIntoTableScan).
	return nil
}

// pruneHiddenFromNode removes hidden symbols from nodes that own an independent
// symbol list. Nodes that delegate to their child (FilterNode, InsertNode, …)
// are not listed here and need no special handling.
func pruneHiddenFromNode(n plan.PlanNode) {
	switch node := n.(type) {
	case *plan.TableScanNode:
		// Keep hidden symbols — source connector uses them.

	case *plan.OutputNode:
		// ColumnNames[i] corresponds to Outputs[i]; filter both in sync.
		var newOutputs []*plan.Symbol
		var newColumnNames []string
		for i, s := range node.Outputs {
			if s.Hidden {
				continue
			}
			newOutputs = append(newOutputs, s)
			if i < len(node.ColumnNames) {
				newColumnNames = append(newColumnNames, node.ColumnNames[i])
			}
		}
		node.Outputs = newOutputs
		node.ColumnNames = newColumnNames

	case *plan.ProjectionNode:
		var newAssignments plan.Assignments
		for _, a := range node.Assignments {
			if !a.Symbol.Hidden {
				newAssignments = append(newAssignments, a)
			}
		}
		node.Assignments = newAssignments

	case *plan.AggregationNode:
		node.Outputs = filterHidden(node.Outputs)

	case *plan.ExchangeNode:
		node.PartitioningScheme.OutputLayout = filterHidden(node.PartitioningScheme.OutputLayout)
		for i, inputs := range node.Inputs {
			node.Inputs[i] = filterHidden(inputs)
		}

	case *plan.RemoteSourceNode:
		node.OutputSymbols = filterHidden(node.OutputSymbols)
	}
}

// filterHidden returns a new slice with hidden symbols removed.
func filterHidden(symbols []*plan.Symbol) []*plan.Symbol {
	var result []*plan.Symbol
	for _, s := range symbols {
		if !s.Hidden {
			result = append(result, s)
		}
	}
	return result
}

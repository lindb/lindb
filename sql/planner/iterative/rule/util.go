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
	"fmt"

	lo "github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func restrictOutputs(idAllcator *plan.PlanNodeIDAllocator,
	node plan.PlanNode, permittedOutputs []*plan.Symbol,
) plan.PlanNode {
	outputs := node.GetOutputSymbols()
	restrictedOutputs := filter(outputs, permittedOutputs)

	if len(outputs) == len(restrictedOutputs) || len(restrictedOutputs) == 0 {
		fmt.Println("outputs same.....")
		return nil
	}
	fmt.Printf("restrictedOutputs, a=%v,b=%v,c=%v\n", outputs, restrictedOutputs, permittedOutputs)

	var assignments plan.Assignments
	assignments = assignments.Add(restrictedOutputs)

	return &plan.ProjectionNode{
		BaseNode: plan.BaseNode{
			ID: idAllcator.Next(),
		},
		Source:      node,
		Assignments: assignments,
	}
}

func restrictChildOutputs(idAllcator *plan.PlanNodeIDAllocator,
	node plan.PlanNode, permittedChildOutputs ...[]*plan.Symbol,
) plan.PlanNode {
	if len(node.GetSources()) != len(permittedChildOutputs) {
		panic(fmt.Sprintf("mismatched child (%d) and permitted outputs (%d) sizes",
			len(node.GetSources()), len(permittedChildOutputs)))
	}

	var newChildren []plan.PlanNode
	rewriteChildren := false

	for i, oldChild := range node.GetSources() {
		newChild := restrictOutputs(idAllcator, oldChild, permittedChildOutputs[i])
		if newChild != nil {
			rewriteChildren = true
			newChildren = append(newChildren, newChild)
		} else {
			newChildren = append(newChildren, oldChild)
		}
	}

	if !rewriteChildren {
		return nil
	}

	return node.ReplaceChildren(newChildren)
}

func pruneInputs(availableInputs []*plan.Symbol, expressions []tree.Expression) []*plan.Symbol {
	symbols := plan.ExtractSymbolsFromExpressions(expressions)
	prunedInputs := filter(availableInputs, symbols)
	if len(prunedInputs) == len(availableInputs) {
		return nil
	}
	return prunedInputs
}

func filter(source, predicate []*plan.Symbol) []*plan.Symbol {
	fmt.Printf("opt rule util.go filter source=%v,predicate=%v\n", source, predicate)
	return lo.Filter(source, func(item *plan.Symbol, index int) bool {
		return lo.ContainsBy(predicate, func(other *plan.Symbol) bool {
			return other.Name == item.Name
		})
	})
}

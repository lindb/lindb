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

package optimization

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/context"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
)

type PruneColumns struct{}

func NewPruneColumns() PlanOptimizer {
	return &PruneColumns{}
}

// Optimize implements PlanOptimizer
func (opt *PruneColumns) Optimize(ctx *context.PlannerContext, node planpkg.PlanNode) planpkg.PlanNode {
	outputs := node.GetOutputSymbols()
	result := node.Accept(outputs, &PruneColumnsVisitor{
		idAllocator: ctx.PlanNodeIDAllocator,
	})
	if r, ok := result.(planpkg.PlanNode); ok {
		return r
	}
	// FIXME: need remove
	return node
}

type PruneColumnsVisitor struct {
	idAllocator *planpkg.PlanNodeIDAllocator
}

func (v *PruneColumnsVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	permittedOutputs := context.([]*planpkg.Symbol)
	fmt.Printf("table scan permitted outputs,%T=%v\n", n, permittedOutputs)
	switch node := n.(type) {
	case *planpkg.ProjectionNode:
		restrictedOutputs := restrictOutputs(node, permittedOutputs)
		var assignments planpkg.Assignments
		node.Assignments = assignments.Add(restrictedOutputs)
	case *planpkg.TableScanNode:
		restrictedOutputs := restrictOutputs(node, permittedOutputs)
		node.OutputSymbols = restrictedOutputs
	}

	for _, child := range n.GetSources() {
		child.Accept(n.GetOutputSymbols(), v)
	}
	return n
}

func restrictOutputs(node planpkg.PlanNode, permittedOutputs []*planpkg.Symbol) (newOutputs []*planpkg.Symbol) {
	outputs := node.GetOutputSymbols()
	// restrict outputs based on permitted outputs
	for _, output := range permittedOutputs {
		o, ok := lo.Find(outputs, func(item *planpkg.Symbol) bool {
			return item.Name == output.Name
		})
		if ok {
			newOutputs = append(newOutputs, o)
		}
	}
	return
}

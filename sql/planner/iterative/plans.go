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

package iterative

import (
	lo "github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

type ResolvingVisitor struct {
	lookup func(groupRef *plan.GroupReference) plan.PlanNode
}

func NewResolvingVisitor(lookup func(groupRef *plan.GroupReference) plan.PlanNode) plan.Visitor {
	return ResolvingVisitor{
		lookup: lookup,
	}
}

func (v ResolvingVisitor) Visit(context any, n plan.PlanNode) (r any) {
	switch node := n.(type) {
	case *plan.GroupReference:
		pNode := v.lookup(node)
		return pNode.Accept(context, v)
	default:
		newChildren := lo.Map(node.GetSources(), func(child plan.PlanNode, index int) plan.PlanNode {
			return child.Accept(context, v).(plan.PlanNode)
		})
		return node.ReplaceChildren(newChildren)
	}
}

func resolveGroupReferences(node plan.PlanNode, lookup func(groupRef *plan.GroupReference) plan.PlanNode) plan.PlanNode {
	return node.Accept(nil, NewResolvingVisitor(lookup)).(plan.PlanNode)
}

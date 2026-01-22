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
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

type Memo struct {
	idAllocator *plan.PlanNodeIDAllocator
	groups      map[int]*plan.Group

	rootGroup   int
	nextGroupID int
}

func NewMemo(idAllocator *plan.PlanNodeIDAllocator, node plan.PlanNode) *Memo {
	memo := &Memo{
		idAllocator: idAllocator,
		nextGroupID: plan.RootGroupRef + 1,
		groups:      make(map[int]*plan.Group),
	}
	memo.rootGroup = memo.insertRecursive(node)
	// TODO: add ref????
	return memo
}

func (m *Memo) resolve(groupRef *plan.GroupReference) plan.PlanNode {
	return m.getNode(groupRef.GroupID)
}

func (m *Memo) replace(groupID int, node plan.PlanNode, _ string) plan.PlanNode {
	group := m.groups[groupID]
	// TODO:check old output???
	if groupRef, ok := node.(*plan.GroupReference); ok {
		node = m.getNode(groupRef.GroupID)
	} else {
		node = m.insertChildrenAndRewrite(node)
	}

	group.Membership = node

	return node
}

func (m *Memo) extract(node plan.PlanNode) plan.PlanNode {
	return resolveGroupReferences(node, m.resolve)
}

func (m *Memo) getNode(group int) plan.PlanNode {
	return m.groups[group].Membership
}

func (m *Memo) insertRecursive(node plan.PlanNode) int {
	if groupRef, ok := node.(*plan.GroupReference); ok {
		return groupRef.GroupID
	}
	group := m.genNextGroupID()
	rewritten := m.insertChildrenAndRewrite(node)
	m.groups[group] = plan.WithMember(rewritten)
	// TODO: inc ref

	return group
}

func (m *Memo) insertChildrenAndRewrite(node plan.PlanNode) plan.PlanNode {
	newChildren := lo.Map(node.GetSources(), func(child plan.PlanNode, index int) plan.PlanNode {
		return &plan.GroupReference{
			BaseNode: plan.BaseNode{
				ID: m.idAllocator.Next(),
			},
			GroupID: m.insertRecursive(child),
			Outputs: child.GetOutputSymbols(),
		}
	})
	return node.ReplaceChildren(newChildren)
}

func (m *Memo) genNextGroupID() int {
	m.nextGroupID++
	return m.nextGroupID
}

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

package plan

import "github.com/lindb/lindb/constants"

const RootGroupRef = 0

type Group struct {
	Membership PlanNode
}

func WithMember(node PlanNode) *Group {
	return &Group{
		Membership: node,
	}
}

type GroupReference struct {
	Outputs []*Symbol
	GroupID int

	BaseNode
}

func (n *GroupReference) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *GroupReference) GetOutputSymbols() []*Symbol {
	return n.Outputs
}

func (n *GroupReference) GetSources() []PlanNode {
	panic(constants.ErrNotSupportOperation)
}

func (n *GroupReference) ReplaceChildren(newChildren []PlanNode) PlanNode {
	panic(constants.ErrNotSupportOperation)
}

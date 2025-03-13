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

type OutputNode struct {
	Source      PlanNode  `json:"source"`
	ColumnNames []string  `json:"columnNames"`
	Outputs     []*Symbol `json:"outputs"`

	BaseNode
}

func (n *OutputNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *OutputNode) GetName() string {
	return "Output"
}

func (n *OutputNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *OutputNode) GetOutputSymbols() []*Symbol {
	return n.Outputs
}

func (n *OutputNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &OutputNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:      newChildren[0],
		ColumnNames: n.ColumnNames,
		Outputs:     n.Outputs,
	}
}

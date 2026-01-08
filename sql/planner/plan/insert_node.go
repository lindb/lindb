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

import "github.com/lindb/lindb/sql/tree"

type InsertNode struct {
	BaseNode

	Database string
	Table    *tree.Table

	Source PlanNode
}

func (n *InsertNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *InsertNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *InsertNode) GetOutputSymbols() []*Symbol {
	return n.Source.GetOutputSymbols()
}

func (n *InsertNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &InsertNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:   newChildren[0],
		Database: n.Database,
		Table:    n.Table,
	}
}

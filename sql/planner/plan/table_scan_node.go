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

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
)

type TableScanNode struct {
	Table         spi.TableHandle               `json:"table"`
	Partitions    map[models.InternalNode][]int `json:"-"`
	OutputSymbols []*Symbol                     `json:"outputs"`
	Assignments   []*spi.ColumnAssignment       `json:"assignments,omitempty"`
	ColumnMapping map[string]string             `json:"columnMapping,omitempty"`

	BaseNode
}

func NewTableScanNode(id PlanNodeID) *TableScanNode {
	return &TableScanNode{
		BaseNode: BaseNode{
			ID: id,
		},
	}
}

func (n *TableScanNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *TableScanNode) GetSources() []PlanNode {
	return nil
}

func (n *TableScanNode) GetOutputSymbols() []*Symbol {
	return n.OutputSymbols
}

func (n *TableScanNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return n
}

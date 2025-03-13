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

type ExchangeScope string

type ExchangeType string

var (
	Local  ExchangeScope = "Local"
	Remote ExchangeScope = "Remote"

	Gather      ExchangeType = "GATHER"
	Repartition ExchangeType = "REPARTITION"
)

type ExchangeNode struct {
	Type    ExchangeType  `json:"type"`
	Scope   ExchangeScope `json:"scope"`
	Sources []PlanNode    `json:"sources"`

	PartitioningScheme *PartitioningScheme `json:"PartitioningScheme"`
	// for each source, the list of inputs corresponding to each output
	Inputs [][]*Symbol `json:"inputs"`

	BaseNode
}

func (n *ExchangeNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *ExchangeNode) GetSources() []PlanNode {
	return n.Sources
}

func (n *ExchangeNode) GetOutputSymbols() []*Symbol {
	return n.PartitioningScheme.OutputLayout
}

func (n *ExchangeNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Type:               n.Type,
		Scope:              n.Scope,
		Sources:            newChildren,
		PartitioningScheme: n.PartitioningScheme,
		Inputs:             n.Inputs,
	}
}

func GatheringExchange(id PlanNodeID, scope ExchangeScope, child PlanNode) *ExchangeNode {
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Type:    Gather,
		Scope:   scope,
		Sources: []PlanNode{child},
		PartitioningScheme: &PartitioningScheme{
			// FIXME: Partitioning: &Partitioning{},
			OutputLayout: child.GetOutputSymbols(),
		},
		Inputs: [][]*Symbol{child.GetOutputSymbols()},
	}
}

func PartitionedExchange(id PlanNodeID, scope ExchangeScope, child PlanNode, partitioningScheme *PartitioningScheme) *ExchangeNode {
	// TODO: add check single
	return &ExchangeNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Type:               Repartition,
		Scope:              scope,
		Sources:            []PlanNode{child},
		PartitioningScheme: partitioningScheme,
		Inputs:             [][]*Symbol{partitioningScheme.OutputLayout},
	}
}

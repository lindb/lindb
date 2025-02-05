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
	"github.com/lindb/lindb/sql/tree"
)

var (
	PARTIAL      = &AggregationStep{InputRaw: true, OutputPartial: true}
	FINAL        = &AggregationStep{InputRaw: false, OutputPartial: false}
	INTERMEDIATE = &AggregationStep{InputRaw: false, OutputPartial: true}
	SINGLE       = &AggregationStep{InputRaw: true, OutputPartial: false}
)

type AggregationStep struct {
	InputRaw      bool
	OutputPartial bool
}

func (step *AggregationStep) String() string {
	switch {
	case step.InputRaw && step.OutputPartial:
		return "PARTIAL"
	case !step.InputRaw && !step.OutputPartial:
		return "FINAL"
	case !step.InputRaw && step.OutputPartial:
		return "INTERMEDIATE"
	case step.InputRaw && !step.OutputPartial:
		return "SINGLE"
	default:
		return "UNKNOWN"
	}
}

type GroupingSetDescriptor struct {
	GroupingKeys []*Symbol `json:"groupingKeys"`
}

type Aggregation struct {
	Function  tree.FuncName     `json:"function"`
	Arguments []tree.Expression `json:"arguments"`
	// TODO: add filter
}

type AggregationAssignment struct {
	Symbol        *Symbol
	ASTExpression tree.Expression
	Aggregation   *Aggregation
}

type AggregationNode struct {
	Source       PlanNode                 `json:"source"`
	GroupingSets *GroupingSetDescriptor   `json:"groupingSets"`
	Step         *AggregationStep         `json:"step"`
	Aggregations []*AggregationAssignment `json:"aggregations"`
	Outputs      []*Symbol                `json:"outputs"`
	BaseNode
}

func NewAggregationNode(id PlanNodeID, source PlanNode,
	aggregations []*AggregationAssignment, groupingSets *GroupingSetDescriptor,
	step *AggregationStep,
) *AggregationNode {
	aggregation := &AggregationNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Aggregations: aggregations,
		Source:       source,
		Step:         step,
		GroupingSets: groupingSets,
	}
	aggregation.Outputs = append(aggregation.Outputs, groupingSets.GroupingKeys...)
	for _, aa := range aggregations {
		aggregation.Outputs = append(aggregation.Outputs, aa.Symbol)
	}
	// TODO: add agg func
	return aggregation
}

func (n *AggregationNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *AggregationNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *AggregationNode) GetOutputSymbols() []*Symbol {
	return n.Outputs
}

func (n *AggregationNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return NewAggregationNode(n.GetNodeID(), newChildren[0], n.Aggregations, n.GroupingSets, n.Step)
}

func (n *AggregationNode) GetGroupingKeys() []*Symbol {
	return n.GroupingSets.GroupingKeys
}

func (n *AggregationNode) HasEmptyGroupingSet() bool {
	// TODO: add global check
	return len(n.GroupingSets.GroupingKeys) == 0
}

func (n *AggregationNode) IsSingleNodeExecutionPreference() bool {
	// 1. aggregations with only empty grouping sets like this:
	// select sum(order_count) from order
	// There is no need for distributed aggregation.
	// Single node FINAL aggregation will suffice, since all input have to be aggregated into one line output.
	// TODO: impl it
	return n.HasEmptyGroupingSet()
}

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

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type PushPartialAggregationThroughExchange struct {
	Base[*plan.AggregationNode]
}

func NewPushPartialAggregationThroughExchange() iterative.Rule {
	rule := &PushPartialAggregationThroughExchange{}
	rule.apply = func(context *iterative.Context, node *plan.AggregationNode) plan.PlanNode {
		exchangeNode, isExchange := context.Lookup.Resolve(node.Source).(*plan.ExchangeNode)
		if !isExchange {
			return nil
		}
		fmt.Printf("push partial through exchange=%v,%v\n", node.Step, exchangeNode.Type)
		// FIXME:add check(exchagne)
		if node.Step == plan.SINGLE &&
			exchangeNode.Type == plan.Repartition {
			return rule.split(context, node)
		}
		if exchangeNode.Type != plan.Gather && exchangeNode.Type != plan.Repartition {
			return nil
		}
		switch node.Step {
		case plan.SINGLE:
			return rule.split(context, node)
		case plan.PARTIAL:
			return rule.pushPartial(context, node, exchangeNode)
		}
		return nil
	}
	return rule
}

func (rule *PushPartialAggregationThroughExchange) pushPartial(context *iterative.Context,
	aggregation *plan.AggregationNode, exchange *plan.ExchangeNode,
) plan.PlanNode {
	var partials []plan.PlanNode
	var inputs [][]*plan.Symbol
	outputs := exchange.GetOutputSymbols()
	for i, source := range exchange.Sources {
		mapping := make(map[string]*plan.Symbol)

		// build symbol mapping
		for idx := range outputs {
			output := outputs[idx]
			input := exchange.Inputs[i][idx]
			if output.Name != input.Name {
				mapping[output.Name] = input
			}
		}

		symbolMapper := plan.NewSymbolMapper(mapping)
		mappedPartial := symbolMapper.MapAggregation(aggregation, source, context.PlannerContext.PlanNodeIDAllocator.Next())

		var assignments plan.Assignments
		for _, output := range aggregation.GetOutputSymbols() {
			input := symbolMapper.MapSymbol(output)
			assignments = assignments.Put(output, input.ToSymbolReference())
		}
		partials = append(partials, &plan.ProjectionNode{
			BaseNode: plan.BaseNode{
				ID: context.PlannerContext.PlanNodeIDAllocator.Next(),
			},
			Source:      mappedPartial,
			Assignments: assignments,
		})
		inputs = append(inputs, aggregation.GetOutputSymbols())
	}

	// TODO: check output

	partitioning := &plan.PartitioningScheme{
		Partitioning: exchange.PartitioningScheme.Partitioning,
		OutputLayout: aggregation.GetOutputSymbols(),
	}
	return &plan.ExchangeNode{
		BaseNode: plan.BaseNode{
			ID: context.PlannerContext.PlanNodeIDAllocator.Next(),
		},
		Type:               exchange.Type,
		Scope:              exchange.Scope,
		Sources:            partials,
		PartitioningScheme: partitioning,
		Inputs:             inputs,
	}
}

func (rule *PushPartialAggregationThroughExchange) split(context *iterative.Context,
	node *plan.AggregationNode,
) plan.PlanNode {
	// TODO: add agg fun
	partial := plan.NewAggregationNode(context.PlannerContext.PlanNodeIDAllocator.Next(),
		node.Source, node.Aggregations, node.GroupingSets, plan.PARTIAL)
	return plan.NewAggregationNode(node.GetNodeID(), partial, node.Aggregations, node.GroupingSets, plan.FINAL)
}

type PushAggregationIntoTableScan struct {
	Base[*plan.AggregationNode]
}

func NewPushAggregationIntoTableScan() iterative.Rule {
	rule := &PushAggregationIntoTableScan{}
	rule.apply = rule.pushAggregationIntoTableScan
	return rule
}

func (rule *PushAggregationIntoTableScan) pushAggregationIntoTableScan(context *iterative.Context,
	node *plan.AggregationNode,
) plan.PlanNode {
	if node.Step != plan.PARTIAL || len(node.Aggregations) == 0 {
		fmt.Printf("setp=%v,aggs=%v\n", node.Step, node.Aggregations)
		// if step is not partial or no aggregation, return nil
		return nil
	}
	// TODO: duplicate?
	var columnAggregations []spi.ColumnAggregation
	for _, agg := range node.Aggregations {
		for _, arg := range agg.Aggregation.Arguments {
			switch argument := arg.(type) {
			case *tree.SymbolReference:
				columnAggregations = append(columnAggregations,
					spi.ColumnAggregation{Column: argument.Name, AggFuncName: agg.Aggregation.Function})
			case *tree.Constant:
				columnAggregations = append(columnAggregations,
					// FIXME: add constant column??
					spi.ColumnAggregation{Column: string(agg.Aggregation.Function), AggFuncName: agg.Aggregation.Function})
			}
		}
	}
	if len(columnAggregations) == 0 {
		fmt.Println("no columnAggregations")
		return nil
	}
	tableScan := iterative.ExtractTableScan(context, node)
	if tableScan == nil {
		return nil
	}
	result := spi.ApplyAggregation(tableScan.Table,
		context.PlannerContext.AnalyzerContext.Analysis.GetTableMetadata(tableScan.Table.String()),
		columnAggregations,
	)
	fmt.Printf("agg columns aggignments: %v\n", result)
	if result == nil || len(result.ColumnAssignments) == 0 {
		return nil
	}
	// just replace column assignments of table scan
	tableScan.Assignments = result.ColumnAssignments
	// TODO: need check???
	tableScan.OutputSymbols = node.GetOutputSymbols()
	return node.Source
}

package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PushPartialAggregationThroughExchange struct {
	Base[*plan.AggregationNode]
}

func NewPushPartialAggregationThroughExchange() iterative.Rule {
	rule := &PushPartialAggregationThroughExchange{}
	rule.apply = func(context *iterative.Context, captures *matching.Captures, node *plan.AggregationNode) plan.PlanNode {
		exchangeNode, isExchange := context.Lookup.Resolve(node.Source).(*plan.ExchangeNode)
		if !isExchange {
			return nil
		}
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
		mappedPartial := symbolMapper.MapAggregation(aggregation, source, context.IDAllocator.Next())

		var assignments plan.Assignments
		for _, output := range aggregation.GetOutputSymbols() {
			input := symbolMapper.MapSymbol(output)
			assignments = assignments.Put(output, input.ToSymbolReference())
		}
		partials = append(partials, &plan.ProjectionNode{
			BaseNode: plan.BaseNode{
				ID: context.IDAllocator.Next(),
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
			ID: context.IDAllocator.Next(),
		},
		Type:               exchange.Type,
		Scope:              exchange.Scope,
		Sources:            partials,
		PartitioningScheme: partitioning,
		Inputs:             inputs,
	}
}

func (rule *PushPartialAggregationThroughExchange) split(context *iterative.Context, node *plan.AggregationNode) plan.PlanNode {
	// TODO: add agg fun
	partial := plan.NewAggregationNode(context.IDAllocator.Next(), node.Source, node.Aggregations, node.GroupingSets, plan.PARTIAL)
	return plan.NewAggregationNode(node.GetNodeID(), partial, node.Aggregations, node.GroupingSets, plan.FINAL)
}

package rule

import (
	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

var EXCHANGE_NODE = matching.NewCapture()

type PushPartialAggregationThroughExchange struct{}

func NewPushPartialAggregationThroughExchange() iterative.Rule {
	return &PushPartialAggregationThroughExchange{}
}

func (rule *PushPartialAggregationThroughExchange) GetPattern() *matching.Pattern {
	// FIXME: add order scheme filter
	// return aggregation()
	// return matching.With(source().Matching(matching.CapturedAs(EXCHANGE_NODE, exchange())), aggregation())
	exchangeNode := matching.CapturedAs(EXCHANGE_NODE, exchange())
	fn := func(source plan.PlanNode, lookup iterative.Lookup) plan.PlanNode {
		if len(source.GetSources()) == 1 {
			sourceNode := source.GetSources()[0]
			node := lookup.Resolve(sourceNode)
			return node
		}
		return nil
	}
	pattern := &matching.Pattern{
		Accept: func(context, val any, captures *matching.Captures) []*matching.Match {
			node := fn(val.(plan.PlanNode), context.(iterative.Lookup))
			result := exchangeNode.Match(context, node, captures)
			return result
		},
		Previous: aggregation(),
	}
	return pattern
}

func (rule *PushPartialAggregationThroughExchange) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if aggregationNode, ok := node.(*plan.AggregationNode); ok {
		exchangeNode, isExchange := context.Lookup.Resolve(aggregationNode.Source).(*plan.ExchangeNode)
		if !isExchange {
			return nil
		}
		// FIXME:add check(exchagne)
		if aggregationNode.Step == plan.SINGLE &&
			exchangeNode.Type == plan.Repartition {
			return rule.split(context, aggregationNode)
		}
		if exchangeNode.Type != plan.Gather && exchangeNode.Type != plan.Repartition {
			return nil
		}
		switch aggregationNode.Step {
		case plan.SINGLE:
			return rule.split(context, aggregationNode)
		case plan.PARTIAL:
			return rule.pushPartial(context, aggregationNode, exchangeNode)
		}
	}
	return nil
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

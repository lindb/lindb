package rule

import (
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

func (rule *PushPartialAggregationThroughExchange) split(context *iterative.Context, node *plan.AggregationNode) plan.PlanNode {
	// TODO: add agg fun
	partial := plan.NewAggregationNode(context.PlannerContext.PlanNodeIDAllocator.Next(), node.Source, node.Aggregations, node.GroupingSets, plan.PARTIAL)
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

func (rule *PushAggregationIntoTableScan) pushAggregationIntoTableScan(context *iterative.Context, node *plan.AggregationNode) plan.PlanNode {
	if node.Step != plan.SINGLE || len(node.Aggregations) == 0 {
		// if step is single or no aggregation, return nil
		return nil
	}
	// TODO: duplicate
	var columnAggregations []spi.ColumnAggregation
	var assigments plan.Assignments
	for _, agg := range node.Aggregations {
		for _, arg := range agg.Aggregation.Arguments {
			if symbol, ok := arg.(*tree.SymbolReference); ok {
				columnAggregations = append(columnAggregations, spi.ColumnAggregation{Column: symbol.Name, AggFuncName: agg.Aggregation.Function})
				assigments = assigments.Put(plan.SymbolFrom(symbol), agg.ASTExpression)
			}
		}
	}
	if len(columnAggregations) == 0 {
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
	if result == nil || len(result.ColumnAssignments) == 0 {
		return nil
	}
	// just replace column assignments of table scan
	tableScan.Assignments = result.ColumnAssignments
	return nil
	//
	// project := &plan.ProjectionNode{
	// 	BaseNode: plan.BaseNode{
	// 		ID: context.PlannerContext.PlanNodeIDAllocator.Next(),
	// 	},
	// 	Source: &plan.TableScanNode{
	// 		BaseNode: plan.BaseNode{
	// 			ID: context.PlannerContext.PlanNodeIDAllocator.Next(),
	// 		},
	// 		Table:         tableScan.Table,
	// 		OutputSymbols: tableScan.OutputSymbols,
	// 		Partitions:    tableScan.Partitions,
	// 		Assignments:   result.ColumnAssignments,
	// 	},
	// 	Assignments: assigments,
	// }
	// return plan.NewAggregationNode(context.PlannerContext.PlanNodeIDAllocator.Next(), project, node.Aggregations, node.GroupingSets, node.Step)
}

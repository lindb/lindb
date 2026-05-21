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
	"github.com/apache/arrow-go/v18/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

var log = logger.GetLogger("Planner", "PushAggregation")

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
		// FIXME: add check (exchange)
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
		// if step is not partial or no aggregation, return nil
		return nil
	}
	columnAggregations := buildColumnAggregations(node.Aggregations)
	log.Info("PushAggregation: built column aggregations",
		logger.Int("count", len(columnAggregations)))
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
		log.Info("PushAggregation: no assignments, skip push")
		return nil
	}
	log.Info("PushAggregation: pushed into table scan",
		logger.Int("assignments", len(result.ColumnAssignments)))
	// Replace physical column assignments of the table scan.
	tableScan.Assignments = result.ColumnAssignments
	// Set output symbols to match the physical assignment column names.
	// For histogram aggregations the aggregation node produces function-result symbols
	// (e.g. "histogram_quantile") which don't correspond to any physical storage field;
	// the source connector's findFieldMeta would fail to resolve them.
	// Using the physical field names directly lets findFieldMeta succeed for every
	// output column, keeping numOfAggs == len(outputs) and preventing a reducer panic.
	tableScan.OutputSymbols = buildTableScanOutputSymbols(node, result.ColumnAssignments)
	return node.Source
}

// buildTableScanOutputSymbols builds the output symbol list for the table scan node
// after aggregation push-down.
//
// For histogram aggregations the AggregationNode's output symbols carry function-result
// names (e.g. "histogram_quantile_0") that don't map to any physical storage field.
// We replace them with one symbol per physical column assignment (bucket fields, stat fields)
// so that source connector's findFieldMeta can resolve every output column by name.
//
// All grouping keys (tags + timestamp) are kept as-is. Timestamp inclusion causes
// isTimestampSelected=true in buildTableScan, which makes the reducer emit TimeSeries
// columns for range queries; HashAggregationOperator then computes the histogram function
// per time slot and emits a TimeSeries result.
func buildTableScanOutputSymbols(node *plan.AggregationNode, assignments []*spi.ColumnAssignment) []*plan.Symbol {
	// Check whether any aggregation is a histogram function.
	hasHisto := false
	for _, aggAssign := range node.Aggregations {
		if tree.IsHistogramFunc(aggAssign.Aggregation.Function) {
			hasHisto = true
			break
		}
	}
	if !hasHisto {
		// No histogram aggregation: keep the existing output symbols unchanged.
		// ColumnMapping in the table scan already translates agg result names to field names.
		return node.GetOutputSymbols()
	}

	// Histogram case: replace aggregation result symbols with physical field symbols.
	// The broker aggregation layer computes the final histogram function from raw values.
	var outputs []*plan.Symbol
	// Keep ALL grouping keys (tags AND timestamp).
	// Timestamp is needed so that isTimestampSelected=true is detected in buildTableScan,
	// which causes the reducer to output TimeSeries columns (per-interval values).
	// HashAggregationOperator will compute the quantile per time slot from those TimeSeries
	// and emit a TimeSeries result; it explicitly sets timestamp output to null because the
	// time information is embedded in the TimeSeries struct itself.
	if node.GroupingSets != nil {
		outputs = append(outputs, node.GroupingSets.GroupingKeys...)
	}
	// One symbol per physical assignment column (bucket fields, stat fields, plain fields).
	for _, a := range assignments {
		outputs = append(outputs, &plan.Symbol{
			Name:     a.Column,
			DataType: larray.NewAggregationType(larray.Sum),
		})
	}
	return outputs
}

// buildColumnAggregations converts a list of AggregationAssignments into ColumnAggregations
// by classifying each argument:
//   - *tree.SymbolReference (non-numeric type) → the metric column name (first wins)
//   - *tree.SymbolReference (Float64/Int64 type) → numeric literal saved in Arguments
//   - *tree.Constant         → treated as column name fallback (function name used)
//   - *tree.FloatLiteral     → literal argument saved in Arguments (e.g. phi for histogram_quantile)
//   - *tree.LongLiteral      → same as FloatLiteral
//
// Entries where no column name is found are skipped.
func buildColumnAggregations(aggregations []*plan.AggregationAssignment) []spi.ColumnAggregation {
	// TODO: duplicate?
	var result []spi.ColumnAggregation
	for _, agg := range aggregations {
		var colName string
		for _, arg := range agg.Aggregation.Arguments {
			switch argument := arg.(type) {
			case *tree.SymbolReference:
				// SymbolReferences with numeric type (Float64/Int64) originated from literal
				// arguments such as the phi value in histogram_quantile(0.99, col).
				// They should be treated as literal parameters, not as column names.
				if isNumericSymbolRef(argument) {
					continue
				} else if colName == "" {
					// Non-numeric SymbolReference: this is the metric column name.
					colName = argument.Name
				}
			case *tree.Constant:
				// Constants are not column references; skip them.
			}
		}
		if colName != "" {
			result = append(result, spi.ColumnAggregation{
				Column:      colName,
				AggFuncName: agg.Aggregation.Function,
			})
		}
	}
	return result
}

// isNumericSymbolRef reports whether a SymbolReference has a numeric data type
// (Float64 or Int64), meaning it originated from a literal argument rather than a
// column reference.  Used to distinguish phi (0.99) from the histogram column name
// in histogram_quantile(0.99, sent_duration) after planAggregation converts arguments.
func isNumericSymbolRef(ref *tree.SymbolReference) bool {
	if ref.DataType == nil {
		return false
	}
	return arrow.TypeEqual(ref.DataType, arrow.PrimitiveTypes.Float64) ||
		arrow.TypeEqual(ref.DataType, arrow.PrimitiveTypes.Int64)
}


package rule

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PruneTableScanColumns struct{}

func NewPruneTableScanColumns() iterative.Rule {
	return &PruneTableScanColumns{}
}

func (rule *PruneTableScanColumns) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if parent, ok := node.(*plan.ProjectionNode); ok {
		table, isTable := context.Lookup.Resolve(parent.Source).(*plan.TableScanNode)
		fmt.Printf("prune table columns 1.....%T,%v\n", context.Lookup.Resolve(parent.Source), isTable)
		if isTable {
			prunedOutputs := pruneInputs(table.GetOutputSymbols(), parent.Assignments.GetExpressions())
			if len(prunedOutputs) == 0 {
				return nil
			}
			newTable := rule.pruneColumns(table, prunedOutputs)
			fmt.Printf("prune table columns 2.....%v,%v\n", newTable, parent.Assignments)
			if newTable != nil {
				fmt.Printf("prune table columns 3.....%v\n", newTable.GetOutputSymbols())
				return parent.ReplaceChildren([]plan.PlanNode{newTable})
			}
			return nil
		}
	}

	return nil
}

func (rule *PruneTableScanColumns) pruneColumns(node *plan.TableScanNode, referencedOutputs []*plan.Symbol) plan.PlanNode {
	newOutputs := lo.Filter(node.GetOutputSymbols(), func(output *plan.Symbol, index int) bool {
		return lo.ContainsBy(referencedOutputs, func(ref *plan.Symbol) bool {
			return ref.Name == output.Name
		})
	})
	fmt.Printf("column 4....%v=%v=%v\n", node.GetOutputSymbols(), referencedOutputs, newOutputs)

	if len(newOutputs) == len(node.GetOutputSymbols()) {
		return nil
	}

	// // TODO: refact?
	// newTable := plan.NewTableScanNode(node.GetNodeID())
	// newTable.Table = node.Table
	// newTable.Partitions = node.Partitions
	// newTable.OutputSymbols = newOutputs
	//
	node.OutputSymbols = newOutputs
	return node
}

type PruneProjectionColumns struct{}

func NewPruneProjectionColumns() iterative.Rule {
	return &PruneProjectionColumns{}
}

func (rule *PruneProjectionColumns) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if parent, ok := node.(*plan.ProjectionNode); ok {
		childProjection, isTable := context.Lookup.Resolve(parent.Source).(*plan.ProjectionNode)
		if isTable {
			prunedOutputs := pruneInputs(childProjection.GetOutputSymbols(), parent.Assignments.GetExpressions())
			return &plan.ProjectionNode{
				BaseNode: plan.BaseNode{
					ID: childProjection.GetNodeID(),
				},
				Source: childProjection.Source,
				Assignments: lo.Filter(childProjection.Assignments, func(assignment *plan.Assignment, index int) bool {
					return lo.ContainsBy(prunedOutputs, func(item *plan.Symbol) bool {
						return assignment.Symbol.Name == item.Name
					})
				}),
			}
		}
	}

	return nil
}

type PruneAggregationSourceColumns struct{}

func NewPruneAggregationSourceColumns() iterative.Rule {
	return &PruneAggregationSourceColumns{}
}

func (rule *PruneAggregationSourceColumns) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if aggregationNode, ok := node.(*plan.AggregationNode); ok {
		var requiredInputs []*plan.Symbol
		requiredInputs = append(requiredInputs, aggregationNode.GetGroupingKeys()...)
		for _, agg := range aggregationNode.Aggregations {
			requiredInputs = append(requiredInputs, planner.ExtractSymbolsFromAggreation(agg.Aggregation)...)
		}
		// TODO: remove
		requiredInputs = lo.UniqBy(requiredInputs, func(item *plan.Symbol) string {
			return item.Name
		})
		fmt.Printf("pure agg sources....%v\n", requiredInputs)
		return restrictChildOutputs(
			context.IDAllocator,
			aggregationNode,
			requiredInputs,
		)
	}
	return nil
}

package rule

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ProjectionOffPushDown[N plan.PlanNode] struct {
	pushDownProjectOff func(context *iterative.Context, targetNode N, referencedOutputs []*plan.Symbol) plan.PlanNode
}

func (rule *ProjectionOffPushDown[N]) Apply(context *iterative.Context, node plan.PlanNode) plan.PlanNode {
	if parent, ok := node.(*plan.ProjectionNode); ok {
		target, isMatch := context.Lookup.Resolve(parent.Source).(N)
		if isMatch {
			prunedOutputs := pruneInputs(target.GetOutputSymbols(), parent.Assignments.GetExpressions())
			if len(prunedOutputs) == 0 {
				return nil
			}
			child := rule.pushDownProjectOff(context, target, prunedOutputs)
			if child != nil {
				return parent.ReplaceChildren([]plan.PlanNode{child})
			}
		}
	}

	return nil
}

type PruneTableScanColumns struct {
	ProjectionOffPushDown[*plan.TableScanNode]
}

func NewPruneTableScanColumns() iterative.Rule {
	rule := &PruneTableScanColumns{}
	rule.pushDownProjectOff = func(context *iterative.Context, table *plan.TableScanNode, referencedOutputs []*plan.Symbol) plan.PlanNode {
		return rule.pruneColumns(table, referencedOutputs)
	}
	return rule
}

func (rule *PruneTableScanColumns) pruneColumns(node *plan.TableScanNode, referencedOutputs []*plan.Symbol) plan.PlanNode {
	newOutputs := lo.Filter(node.GetOutputSymbols(), func(output *plan.Symbol, index int) bool {
		return lo.ContainsBy(referencedOutputs, func(ref *plan.Symbol) bool {
			return ref.Name == output.Name
		})
	})
	if len(newOutputs) == len(node.GetOutputSymbols()) {
		return nil
	}

	// TODO: create new table node?
	node.OutputSymbols = newOutputs
	return node
}

type PruneProjectionColumns struct {
	ProjectionOffPushDown[*plan.ProjectionNode]
}

func NewPruneProjectionColumns() iterative.Rule {
	rule := &PruneProjectionColumns{}
	rule.pushDownProjectOff = func(context *iterative.Context, childProjection *plan.ProjectionNode, referencedOutputs []*plan.Symbol) plan.PlanNode {
		return &plan.ProjectionNode{
			BaseNode: plan.BaseNode{
				ID: childProjection.GetNodeID(),
			},
			Source: childProjection.Source,
			Assignments: lo.Filter(childProjection.Assignments, func(assignment *plan.Assignment, index int) bool {
				return lo.ContainsBy(referencedOutputs, func(item *plan.Symbol) bool {
					return assignment.Symbol.Name == item.Name
				})
			}),
		}
	}
	return rule
}

type PruneAggregationSourceColumns struct {
	Base[*plan.AggregationNode]
}

func NewPruneAggregationSourceColumns() iterative.Rule {
	rule := &PruneAggregationSourceColumns{}
	rule.apply = func(context *iterative.Context, node *plan.AggregationNode) plan.PlanNode {
		var requiredInputs []*plan.Symbol
		requiredInputs = append(requiredInputs, node.GetGroupingKeys()...)
		for _, agg := range node.Aggregations {
			requiredInputs = append(requiredInputs, planner.ExtractSymbolsFromAggreation(agg.Aggregation)...)
		}
		return restrictChildOutputs(
			context.IDAllocator,
			node,
			requiredInputs,
		)
	}
	return rule
}

type PruneFilterColumns struct {
	ProjectionOffPushDown[*plan.FilterNode]
}

func NewPruneFilterColumns() iterative.Rule {
	rule := &PruneFilterColumns{}
	rule.pushDownProjectOff = func(context *iterative.Context, filter *plan.FilterNode, referencedOutputs []*plan.Symbol) plan.PlanNode {
		// TODO: add filter columns
		return restrictChildOutputs(context.IDAllocator, filter, referencedOutputs)
	}
	return rule
}

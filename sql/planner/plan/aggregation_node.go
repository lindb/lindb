package plan

import (
	"github.com/lindb/lindb/spi/function"
	"github.com/lindb/lindb/sql/tree"
)

type AggregationStep struct {
	InputRaw      bool
	OutputPartial bool
}

type GroupingSetDescriptor struct {
	GroupingKeys []*Symbol `json:"groupingKeys"`
}

type Aggregation struct {
	ResolvedFunction *function.ResolvedFunction `json:"resolvedFunction"`
	Arguments        []tree.Expression          `json:"arguments"`
}

type AggregationAssignment struct {
	Symbol        *Symbol
	ASTExpression tree.Expression
	Aggregation   *Aggregation
}

type AggregationNode struct {
	Source       PlanNode                 `json:"source"`
	Aggregations []*AggregationAssignment `json:"aggregations"`
	GroupingSets *GroupingSetDescriptor   `json:"groupingSets"`
	Outputs      []*Symbol                `json:"outputs"`

	BaseNode
}

func NewAggregationNode(id PlanNodeID, source PlanNode, aggregations []*AggregationAssignment, groupingSets *GroupingSetDescriptor) *AggregationNode {
	aggregation := &AggregationNode{
		BaseNode: BaseNode{
			ID: id,
		},
		Aggregations: aggregations,
		Source:       source,
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
	return NewAggregationNode(n.GetNodeID(), newChildren[0], n.Aggregations, n.GroupingSets)
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
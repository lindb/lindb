package plan

import "github.com/lindb/lindb/sql/tree"

type FilterNode struct {
	Source    PlanNode        `json:"source"`
	Predicate tree.Expression `json:"predicate"`

	BaseNode
}

func (n *FilterNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *FilterNode) GetName() string {
	return "Filter"
}

func (n *FilterNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *FilterNode) GetOutputSymbols() []*Symbol {
	return n.Source.GetOutputSymbols()
}

func (n *FilterNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &FilterNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:    newChildren[0],
		Predicate: n.Predicate,
	}
}

package plan

import "github.com/lindb/lindb/sql/tree"

type InsertNode struct {
	BaseNode

	Source PlanNode `json:"source"`

	Table *tree.Table
}

func (n *InsertNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *InsertNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *InsertNode) GetOutputSymbols() []*Symbol {
	return n.Source.GetOutputSymbols()
}

func (n *InsertNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &InsertNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source: newChildren[0],
		Table:  n.Table,
	}
}

package plan

import "github.com/lindb/lindb/sql/tree"

type Assignments map[*Symbol]tree.Expression

func (a Assignments) Add(symbols []*Symbol) {
	for _, symbol := range symbols {
		a[symbol] = symbol.ToSymbolReference()
	}
}

func (a Assignments) GetOutputs() (outputs []*Symbol) {
	for k := range a {
		outputs = append(outputs, k)
	}
	return
}

type ProjectionNode struct {
	Source      PlanNode    `json:"source"`
	Assignments Assignments `json:"-"` // FIXME:

	BaseNode
}

func (n *ProjectionNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *ProjectionNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *ProjectionNode) GetOutputSymbols() []*Symbol {
	return n.Assignments.GetOutputs()
}

func (n *ProjectionNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &ProjectionNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:      newChildren[0],
		Assignments: n.Assignments,
	}
}

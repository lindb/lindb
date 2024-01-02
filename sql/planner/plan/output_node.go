package plan

type OutputNode struct {
	Source      PlanNode  `json:"source"`
	ColumnNames []string  `json:"columnNames"`
	Outputs     []*Symbol `json:"outputs"`

	BaseNode
}

func (n *OutputNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *OutputNode) GetName() string {
	return "Output"
}

func (n *OutputNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *OutputNode) GetOutputSymbols() []*Symbol {
	return n.Outputs
}

func (n *OutputNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &OutputNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:      newChildren[0],
		ColumnNames: n.ColumnNames,
		Outputs:     n.Outputs,
	}
}

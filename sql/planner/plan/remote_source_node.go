package plan

type RemoteSourceNode struct {
	SourceFragmentIDs []FragmentID `json:"sourceFragmentIDs,omitempty"`
	OutputSymbols     []*Symbol    `json:"outputs"`

	BaseNode
}

func (n *RemoteSourceNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *RemoteSourceNode) GetSources() []PlanNode {
	return nil
}

func (n *RemoteSourceNode) GetOutputSymbols() []*Symbol {
	return n.OutputSymbols
}

func (n *RemoteSourceNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return n
}

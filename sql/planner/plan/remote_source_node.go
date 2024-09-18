package plan

type RemoteSourceNode struct {
	SourceFragmentIDs []FragmentID `json:"sourceFragmentIDs,omitempty"`

	BaseNode
}

func (n *RemoteSourceNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *RemoteSourceNode) GetSources() []PlanNode {
	return nil
}

func (n *RemoteSourceNode) GetOutputSymbols() []*Symbol {
	return nil
}

func (n *RemoteSourceNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return n
}

package plan

type RemoteSourceNode struct {
	BaseNode

	SourceFragmentIDs []FragmentID `json:"sourceFragmentIDs,omitempty"`
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

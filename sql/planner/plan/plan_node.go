package plan

type FragmentID int

const RootFragmentID = FragmentID(0)

type PlanNodeID int64

type PlanNode interface {
	GetNodeID() PlanNodeID
	GetSources() []PlanNode
	GetOutputSymbols() []*Symbol

	ReplaceChildren(newChildren []PlanNode) PlanNode

	Accept(context any, visitor Visitor) (r any)
}

type BaseNode struct {
	ID PlanNodeID `json:"id"`
}

func (n *BaseNode) GetNodeID() PlanNodeID {
	return n.ID
}

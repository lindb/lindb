package plan

import "github.com/lindb/lindb/constants"

const RootGroupRef = 0

type Group struct {
	Membership PlanNode
}

func WithMember(node PlanNode) *Group {
	return &Group{
		Membership: node,
	}
}

type GroupReference struct {
	Outputs []*Symbol
	GroupID int

	BaseNode
}

func (n *GroupReference) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *GroupReference) GetOutputSymbols() []*Symbol {
	return n.Outputs
}

func (n *GroupReference) GetSources() []PlanNode {
	panic(constants.ErrNotSupportOperation)
}

func (n *GroupReference) ReplaceChildren(newChildren []PlanNode) PlanNode {
	panic(constants.ErrNotSupportOperation)
}

package tree

type NodeID int

type Node interface {
	Accept(context any, visitor Visitor) any

	GetID() NodeID
}

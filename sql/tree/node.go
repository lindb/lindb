package tree

type NodeID int64

type Node interface {
	Accept(context any, visitor Visitor) any

	GetID() NodeID
	SetID(id NodeID)
}

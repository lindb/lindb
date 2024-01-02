package tree

type NodeIDAllocator struct {
	next NodeID
}

func NewNodeIDAllocator() *NodeIDAllocator {
	return &NodeIDAllocator{}
}

func (a *NodeIDAllocator) Next() NodeID {
	a.next++
	return a.next
}

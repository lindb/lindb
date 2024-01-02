package tree

type Limit struct {
	BaseNode
	RowCount Expression
}

// Accept implements Node
func (n *Limit) Accept(context any, visitor Visitor) any {
	panic("unimplemented")
}

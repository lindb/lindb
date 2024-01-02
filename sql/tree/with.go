package tree

type With struct {
	BaseNode
	Queries []*WithQuery
}

// Accept implements Node
func (n *With) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type WithQuery struct {
	BaseNode
	Name  *Identifier
	Query *Query
}

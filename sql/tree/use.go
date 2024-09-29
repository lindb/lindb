package tree

// Use represents the statement that select search database.
type Use struct {
	BaseNode
	Database *Identifier
}

func (n *Use) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

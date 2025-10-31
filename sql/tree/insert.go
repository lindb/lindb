package tree

type Insert struct {
	BaseNode

	Table *Table
	Query *Query
}

func (n *Insert) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

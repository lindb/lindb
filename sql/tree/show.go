package tree

type ShowColumns struct {
	BaseNode

	Table *Table
}

func (n *ShowColumns) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

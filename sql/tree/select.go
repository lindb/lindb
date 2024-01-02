package tree

type Select struct {
	SelectItems []SelectItem
}

type SelectItem interface {
	Node
}

type AllColumns struct {
	BaseNode
	Target Expression
}

func (n *AllColumns) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type SingleColumn struct {
	BaseNode
	Expression Expression
	Aliase     *Identifier
}

func (n *SingleColumn) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

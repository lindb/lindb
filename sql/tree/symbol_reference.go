package tree

type SymbolReference struct {
	BaseNode

	Name string
}

func (n *SymbolReference) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

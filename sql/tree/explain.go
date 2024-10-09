package tree

type Explain struct {
	BaseNode
	Statement Statement
}

func (n *Explain) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

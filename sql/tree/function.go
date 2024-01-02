package tree

type FunctionCall struct {
	BaseNode
	Name      QualifiedName
	Arguments []Expression
}

// Accept implements Expression
func (n *FunctionCall) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

package tree

type FunctionCall struct {
	BaseNode
	RefField  *Field
	Name      FunctionName
	Arguments []Expression
}

// Accept implements Expression
func (n *FunctionCall) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

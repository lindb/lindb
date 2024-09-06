package tree

type Call struct {
	BaseNode

	Function FunctionName // FIXME: add func
	Args     []Expression
}

func (n *Call) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

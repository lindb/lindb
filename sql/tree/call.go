package tree

import "github.com/lindb/lindb/spi/types"

type Call struct {
	BaseNode

	Function FunctionName
	Args     []Expression
	RetType  types.DataType
}

func (n *Call) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

package tree

import "github.com/lindb/lindb/spi/types"

type Call struct {
	BaseNode

	Function FunctionName // FIXME: add func
	RetType  types.DataType
	Args     []Expression
}

func (n *Call) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

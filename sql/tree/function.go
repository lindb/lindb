package tree

import "github.com/lindb/lindb/spi/types"

type FunctionCall struct {
	BaseNode
	RefField  *Field
	Name      FuncName
	Arguments []Expression
	RetType   types.DataType
}

// Accept implements Expression
func (n *FunctionCall) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

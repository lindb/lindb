package tree

import "github.com/lindb/lindb/spi/types"

type SymbolReference struct {
	BaseNode

	Name     string
	DataType types.DataType
}

func (n *SymbolReference) Accept(context any, vistor Visitor) any {
	return vistor.Visit(context, n)
}

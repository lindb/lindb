package plan

import (
	"github.com/lindb/lindb/spi/types"
)

type ValuesNode struct {
	Rows          *types.Page
	OutputSymbols []*Symbol
	BaseNode
	RowCount int
}

func (n *ValuesNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *ValuesNode) GetSources() []PlanNode {
	return nil
}

func (n *ValuesNode) GetOutputSymbols() []*Symbol {
	return n.OutputSymbols
}

func (n *ValuesNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return n
}

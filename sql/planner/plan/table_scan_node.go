package plan

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
)

type TableScanNode struct {
	Table         spi.TableHandle               `json:"table"`
	Partitions    map[models.InternalNode][]int `json:"-"`
	Database      string                        `json:"database"`
	OutputSymbols []*Symbol                     `json:"outputSymbols"`

	BaseNode
}

func NewTableScanNode(id PlanNodeID) *TableScanNode {
	return &TableScanNode{
		BaseNode: BaseNode{
			ID: id,
		},
	}
}

func (n *TableScanNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *TableScanNode) GetSources() []PlanNode {
	return nil
}

func (n *TableScanNode) GetOutputSymbols() []*Symbol {
	return n.OutputSymbols
}

func (n *TableScanNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return n
}

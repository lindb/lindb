package plan

import "github.com/lindb/lindb/sql/tree"

type PartitioningScheme struct {
	Partitioning *Partitioning `json:"partitioning"`
	OutputLayout []*Symbol     `json:"outputLayout"`
}

type Partitioning struct {
	Handle    *PartitioningHandle `json:"handle"`
	Arguments []*ArgumentBinding  `json:"arguments"`
}

func (p *Partitioning) Translate(translator func(symbol *Symbol) *Symbol) *Partitioning {
	// FIXME: imple translate
	return nil
}

type PartitioningHandle struct{}

func (h *PartitioningHandle) IsSingleNode() bool {
	return false
}

type ArgumentBinding struct {
	Expression tree.Expression `json:"expression"`
}

func (arg *ArgumentBinding) Translate() *ArgumentBinding {
	// FIXME: imple arg binding translate
	return nil
}

type PartitioningProps struct {
	PartitioningColumns []*Symbol
}

func singlePartition() *PartitioningProps {
	return &PartitioningProps{}
}

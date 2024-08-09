package plan

import "github.com/lindb/lindb/sql/tree"

type Partitioning struct {
	Handle    *PartitioningHandle `json:"handle"`
	Arguments []*ArgumentBinding  `json:"arguments"`
}

type PartitioningHandle struct{}

func (h *PartitioningHandle) IsSingleNode() bool {
	return false
}

type ArgumentBinding struct {
	Expression tree.Expression `json:"expression"`
}

func (arg *ArgumentBinding) Translate() *ArgumentBinding {
	return nil
}

type PartitioningProps struct {
	PartitioningColumns []*Symbol
}

func singlePartition() *PartitioningProps {
	return &PartitioningProps{}
}

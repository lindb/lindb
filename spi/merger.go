package spi

import "github.com/lindb/lindb/spi/types"

type Merger interface {
	AddSplit(split *BinarySplit)
	GetOutputPage() *types.Page
}

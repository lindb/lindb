package spi

type Merger interface {
	AddSplit(split *BinarySplit)
	GetOutputPage() *Page
}

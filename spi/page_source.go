package spi

import (
	"github.com/lindb/lindb/sql/tree"
)

type PageSource interface {
	AddSplit(split Split)
	GetNextPage() *Page
}

type SplitSource interface {
	Prepare()
	HasSplit() bool
	GetNextSplit() Split
}

type SplitSourceProvider interface {
	CreateSplitSources(database string, table TableHandle, partitions []int, columns []ColumnMetadata, filter tree.Expression) (splits []SplitSource)
}
type PageSourceProvider interface {
	CreatePageSource(table TableHandle) PageSource
}

type PageSourceManager struct{}

type SplitSourceFactory struct{}

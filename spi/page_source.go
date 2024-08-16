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
	HasNext() bool
	Next() Split
}

type SplitSourceProvider interface {
	CreateSplitSources(table TableHandle, partitions []int, outputColumns []ColumnMetadata, predicate tree.Expression) (splits []SplitSource)
}

type PageSourceProvider interface {
	CreatePageSource(table TableHandle) PageSource
}

type PageSourceManager struct{}

type SplitSourceFactory struct{}

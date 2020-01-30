package indexdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series"
)

type groupingContext struct {
	gCtx series.GroupingContext
}

func (g *groupingContext) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	// need add read lock
	return g.gCtx.BuildGroup(highKey, container)
}

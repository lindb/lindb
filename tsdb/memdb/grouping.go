package memdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series"
)

type groupingContext struct {
	ms *metricStore

	gCtx series.GroupingContext
}

func (g *groupingContext) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	// need add read lock
	g.ms.mux.RLock()
	defer g.ms.mux.RUnlock()
	return g.gCtx.BuildGroup(highKey, container)
}

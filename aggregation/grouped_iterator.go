package aggregation

import "github.com/lindb/lindb/series"

//FIXME stone1100 need refactor
type groupedIterator struct {
	tags       map[string]string
	aggregates map[uint16]FieldAggregator
	fieldIDs   []uint16

	idx int
}

func newGroupedIterator(tags map[string]string, aggregates map[uint16]FieldAggregator) series.GroupedIterator {
	fieldIDs := make([]uint16, len(aggregates))
	idx := 0
	for fieldID := range aggregates {
		fieldIDs[idx] = fieldID
		idx++
	}
	return &groupedIterator{tags: tags, aggregates: aggregates, fieldIDs: fieldIDs}
}

func (g *groupedIterator) Tags() map[string]string {
	return g.tags
}

func (g *groupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldIDs) {
		return false
	}
	g.idx++
	return true
}

func (g *groupedIterator) Next() series.FieldIterator {
	return g.aggregates[g.fieldIDs[g.idx-1]].Iterator()
}

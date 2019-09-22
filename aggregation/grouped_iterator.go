package aggregation

import "github.com/lindb/lindb/series"

type groupedIterator struct {
	tags       map[string]string
	aggregates FieldAggregates
	len        int
	idx        int
}

func newGroupedIterator(tags map[string]string, aggregates FieldAggregates) series.GroupedIterator {
	return &groupedIterator{tags: tags, aggregates: aggregates, len: len(aggregates)}
}

func (g *groupedIterator) Tags() map[string]string {
	return g.tags
}

func (g *groupedIterator) HasNext() bool {
	if g.idx >= g.len {
		return false
	}
	g.idx++
	return true
}

func (g *groupedIterator) Next() series.Iterator {
	return g.aggregates[g.idx-1].ResultSet()
}

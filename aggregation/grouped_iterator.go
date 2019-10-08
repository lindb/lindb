package aggregation

import "github.com/lindb/lindb/series"

//////////////////////////////////////////////////////
// binaryGroupedIterator implements GroupedIterator
//////////////////////////////////////////////////////
type groupedIterator struct {
	tags       map[string]string
	aggregates FieldAggregates
	len        int
	idx        int
}

// newGroupedIterator creates a grouped iterator for field aggregates
func newGroupedIterator(tags map[string]string, aggregates FieldAggregates) series.GroupedIterator {
	return &groupedIterator{tags: tags, aggregates: aggregates, len: len(aggregates)}
}

// Tags returns the tags of series
func (g *groupedIterator) Tags() map[string]string {
	return g.tags
}

// HasNext returns if the iteration has more field's iterator
func (g *groupedIterator) HasNext() bool {
	if g.idx >= g.len {
		return false
	}
	g.idx++
	return true
}

// Next returns the result set of aggregator
func (g *groupedIterator) Next() series.Iterator {
	return g.aggregates[g.idx-1].ResultSet()
}

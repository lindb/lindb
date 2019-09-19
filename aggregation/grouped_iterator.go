package aggregation

import "github.com/lindb/lindb/series"

//FIXME stone1100 need refactor
type groupedIterator struct {
	tags       map[string]string
	aggregates map[string]FieldAggregator
	fieldNames []string

	idx int
}

func newGroupedIterator(tags map[string]string, aggregates map[string]FieldAggregator) series.GroupedIterator {
	fieldNames := make([]string, len(aggregates))
	idx := 0
	for fieldName := range aggregates {
		fieldNames[idx] = fieldName
		idx++
	}
	return &groupedIterator{tags: tags, aggregates: aggregates, fieldNames: fieldNames}
}

func (g *groupedIterator) Tags() map[string]string {
	return g.tags
}

func (g *groupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldNames) {
		return false
	}
	g.idx++
	return true
}

func (g *groupedIterator) Next() series.FieldIterator {
	return g.aggregates[g.fieldNames[g.idx-1]].Iterator()
}

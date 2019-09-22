package aggregation

import "github.com/lindb/lindb/series"

// seriesIterator implements series.Iterator
type seriesIterator struct {
	fieldName  string
	aggregates []FieldAggregator
	idx        int
	len        int
}

// newSeriesIterator creates the time series iterator
func newSeriesIterator(agg SeriesAggregator) series.Iterator {
	it := &seriesIterator{fieldName: agg.FieldName(), aggregates: agg.Aggregates()}
	it.len = len(it.aggregates)
	return it
}

// FieldName return field name
func (s *seriesIterator) FieldName() string {
	return s.fieldName
}

// HasNext returns if the iteration has more field's iterator
func (s *seriesIterator) HasNext() bool {
	for s.idx < s.len {
		if s.aggregates[s.idx] != nil {
			s.idx++
			return true
		}
		s.idx++
	}
	return false
}

// Next returns the field's iterator and segment start time
func (s *seriesIterator) Next() (startTime int64, fieldIt series.FieldIterator) {
	agg := s.aggregates[s.idx-1]
	return agg.ResultSet()
}

package aggregation

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// seriesIterator implements series.Iterator
type seriesIterator struct {
	fieldName   string
	fieldType   field.Type
	aggregators []FieldAggregator
	idx         int
	len         int
}

// newSeriesIterator creates the time series iterator
func newSeriesIterator(agg SeriesAggregator) series.Iterator {
	it := &seriesIterator{fieldName: agg.FieldName(), fieldType: agg.GetFieldType(), aggregators: agg.Aggregators()}
	it.len = len(it.aggregators)
	return it
}

// FieldName returns field name
func (s *seriesIterator) FieldName() string {
	return s.fieldName
}

// FieldType returns field type
func (s *seriesIterator) FieldType() field.Type {
	return s.fieldType
}

// HasNext returns if the iteration has more field's iterator
func (s *seriesIterator) HasNext() bool {
	for s.idx < s.len {
		if s.aggregators[s.idx] != nil {
			s.idx++
			return true
		}
		s.idx++
	}
	return false
}

// Next returns the field's iterator and segment start time
func (s *seriesIterator) Next() (startTime int64, fieldIt series.FieldIterator) {
	agg := s.aggregators[s.idx-1]
	return agg.ResultSet()
}

func (s *seriesIterator) MarshalBinary() ([]byte, error) {
	return series.MarshalIterator(s)
}

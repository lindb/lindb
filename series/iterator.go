package series

import (
	enc "encoding"

	"github.com/lindb/lindb/models"

	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./iterator.go -destination=./iterator_mock.go -package=series

// TimeSeriesEvent represents time series event for query
type TimeSeriesEvent struct {
	SeriesList []GroupedIterator

	Stats *models.QueryStats
	Err   error
}

// GroupedIterator represents a iterator for the grouped time series data
type GroupedIterator interface {
	// HasNext returns if the iteration has more field's iterator
	HasNext() bool
	// Next returns the field's iterator
	Next() Iterator
	// Tags returns group tags, tags is tag values concat string
	Tags() string
}

// Iterator represents an iterator for the time series data
type Iterator interface {
	// FieldName returns the field name
	FieldName() string
	// FieldType returns the field type
	FieldType() field.Type
	// HasNext returns if the iteration has more field's iterator
	HasNext() bool
	// Next returns the field's iterator
	Next() (startTime int64, fieldIt FieldIterator)
	// MarshalBinary marshals the data
	enc.BinaryMarshaler
}

// FieldIterator represents a field's data iterator, support multi field for one series
type FieldIterator interface {
	// HasNext returns if the iteration has more fields
	HasNext() bool
	// Next returns the primitive field iterator
	// because there are some primitive fields if field type is complex
	Next() PrimitiveIterator
	// MarshalBinary marshals the data
	enc.BinaryMarshaler
}

// PrimitiveIterator represents an iterator over a primitive field, iterator points data of primitive field
type PrimitiveIterator interface {
	// FieldID returns the primitive field id
	FieldID() field.PrimitiveID
	// AggType returns the primitive field's agg type
	AggType() field.AggType
	// HasNext returns if the iteration has more data points
	HasNext() bool
	// Next returns the data point in the iteration
	Next() (timeSlot int, value float64)
}

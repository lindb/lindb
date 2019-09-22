package series

import (
	"io"

	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./iterator.go -destination=./iterator_mock.go -package=series

// TimeSeriesEvent represents time series event for query
type TimeSeriesEvent struct {
	SeriesList []GroupedIterator

	Err error
}

// VersionIterator represents a multi-version iterator
type VersionIterator interface {
	// Version returns the version no.
	Version() Version
	// HasNext returns if the iteration has more time-series's iterator
	HasNext() bool
	// Next returns the time-series's iterator
	Next() Iterator
	// Close closes the underlying resource
	io.Closer
}

// GroupedIterator represents a iterator for the grouped time series data
type GroupedIterator interface {
	// HasNext returns if the iteration has more field's iterator
	HasNext() bool
	// Next returns the field's iterator
	Next() Iterator
	// Tags returns group tags
	Tags() map[string]string
}

// Iterator represents an iterator for the time series data
type Iterator interface {
	// FieldName return field name
	FieldName() string
	// HasNext returns if the iteration has more field's iterator
	HasNext() bool
	// Next returns the field's iterator
	Next() (startTime int64, fieldIt FieldIterator)
}

// FieldIterator represents a field's data iterator, support multi field for one series
type FieldIterator interface {
	// HasNext returns if the iteration has more fields
	HasNext() bool
	// Next returns the primitive field iterator
	// because there are some primitive fields if field type is complex
	Next() PrimitiveIterator
	// Bytes returns the binary data for field iterator
	Bytes() ([]byte, error)
}

// PrimitiveIterator represents an iterator over a primitive field, iterator points data of primitive field
type PrimitiveIterator interface {
	// FieldID returns the primitive field id
	FieldID() uint16
	// AggType returns the primitive field's agg type
	AggType() field.AggType
	// HasNext returns if the iteration has more data points
	HasNext() bool
	// Next returns the data point in the iteration
	Next() (timeSlot int, value float64)
}

package field

//go:generate mockgen -source ./iterator.go -destination=./iterator_mock.go -package=field

// TimeSeries represents an iterator for the time series data
type TimeSeries interface {
	// HasNext returns if the iteration has more field's iterator
	HasNext() bool
	// Next returns the field's iterator
	Next() Iterator
}

// GroupedTimeSeries represents a iterator for the grouped time series data
type GroupedTimeSeries interface {
	TimeSeries
	// Tags returns group tags
	Tags() map[string]string
}

// MultiTimeSeries represents a multi-version iterator for the time series data
type MultiTimeSeries interface {
	TimeSeries
	// Version returns the version no.
	Version() uint32
	// ID returns the time series id under current metric
	ID() uint32
}

// Iterator represents a field's data iterator, support multi field for one series
type Iterator interface {
	// ID returns the field's id
	ID() uint16
	// Name return the field's name
	Name() string
	// FieldType returns the field's type
	FieldType() Type
	// HasNext returns if the iteration has more fields
	HasNext() bool
	// Next returns the primitive field iterator
	// because there are some primitive fields if field type is complex
	Next() PrimitiveIterator
}

// PrimitiveIterator represents an iterator over a primitive field, iterator points data of primitive field
type PrimitiveIterator interface {
	// ID returns the primitive field id
	ID() uint16
	// HasNext returns if the iteration has more data points
	HasNext() bool
	// Next returns the data point in the iteration
	Next() (timeSlot int, value float64)
}

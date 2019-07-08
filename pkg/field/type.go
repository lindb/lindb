package field

// ValueType represents primitive field's value type
type ValueType int

// Defines all value type of primitive field
const (
	Integer ValueType = iota + 1
	Float
)

// AggType represents primitive field's aggregator type
type AggType int

// Defines all aggregator types for primitive field
const (
	Sum AggType = iota + 1
	Min
	Max
)

// Type represents field type for LinDB support
type Type int16

// Defines all field types for LinDB support(user write)
const (
	SumField Type = iota + 1
	MinField
	MaxField
	HistogramField
)

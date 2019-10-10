package field

import "github.com/lindb/lindb/aggregation/function"

// ValueType represents primitive field's value type
type ValueType uint8

// Defines all value type of primitive field
const (
	Integer ValueType = iota + 1
	Float
)

// AggType represents primitive field's aggregator type
type AggType uint8

// Defines all aggregator types for primitive field
const (
	Sum AggType = iota + 1
	Count
	Min
	Max
)

var schemas = map[Type]schema{}

func init() {
	schemas[SumField] = newSumSchema()
	schemas[MinField] = newMinSchema()
	schemas[SummaryField] = newSummarySchema()
}

// Type represents field type for LinDB support
type Type uint8

// Defines all field types for LinDB support(user write)
const (
	SumField Type = iota + 1
	MinField
	MaxField
	SummaryField
	HistogramField

	Unknown
)

// String returns the field type's string value
func (t Type) String() string {
	switch t {
	case SumField:
		return "sum"
	case MinField:
		return "min"
	case MaxField:
		return "max"
	case SummaryField:
		return "summary"
	case HistogramField:
		return "histogram"
	default:
		return "unknown"
	}
}

// GetPrimitiveFields returns the primitive fields for aggregator
func (t Type) GetPrimitiveFields(funcType function.FuncType) map[uint16]AggType {
	schema := schemas[t]
	if schema == nil {
		return nil
	}
	return schema.getPrimitiveFields(funcType)
}

// GetDefaultPrimitiveFields returns the default primitive fields for aggregator
func (t Type) GetDefaultPrimitiveFields() map[uint16]AggType {
	schema := schemas[t]
	if schema == nil {
		return nil
	}
	return schema.getDefaultPrimitiveFields()
}

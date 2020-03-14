package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

// AggType represents primitive field's aggregator type
type AggType uint8
type PrimitiveID uint8
type ID uint8

// Field key represents field id[1byte] + primitive field id[1byte]
type Key uint16

// Defines all aggregator types for primitive field
const (
	Sum AggType = iota + 1
	Count
	Min
	Max
	Replace
)

var schemas = map[Type]Schema{}

func init() {
	schemas[SumField] = newSumSchema()
	schemas[MinField] = newMinSchema()
	schemas[MaxField] = newMaxSchema()
	schemas[GaugeField] = newGaugeSchema()
	schemas[SummaryField] = newSummarySchema()
	schemas[IncreaseField] = newIncreaseSchema()
}

// Type represents field type for LinDB support
type Type uint8

// Defines all field types for LinDB support(user write)
const (
	SumField Type = iota + 1
	MinField
	MaxField
	GaugeField
	IncreaseField
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
	case GaugeField:
		return "gauge"
	case IncreaseField:
		return "increase"
	case SummaryField:
		return "summary"
	case HistogramField:
		return "histogram"
	default:
		return "unknown"
	}
}

func (t Type) DownSamplingFunc() function.FuncType {
	switch t {
	case SumField:
		return function.Sum
	case MinField:
		return function.Min
	case MaxField:
		return function.Max
	case GaugeField:
		return function.Replace
	case IncreaseField:
		return function.Sum
	case SummaryField:
		return function.Count
	case HistogramField:
		return function.Histogram
	default:
		return function.Unknown
	}
}

func (t Type) IsFuncSupported(funcType function.FuncType) bool {
	switch t {
	case SumField:
		switch funcType {
		case function.Sum, function.Min, function.Max:
			return true
		default:
			return false
		}
	case MinField:
		switch funcType {
		case function.Min:
			return true
		default:
			return false
		}
	case MaxField:
		switch funcType {
		case function.Max:
			return true
		default:
			return false
		}
	case GaugeField:
		switch funcType {
		case function.Sum, function.Min, function.Max, function.Replace:
			return true
		default:
			return false
		}
	case SummaryField:
		return true
	case HistogramField:
		return true
	default:
		return false
	}
}

// GetPrimitiveFields returns the primitive fields for aggregator
func (t Type) GetPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	schema := schemas[t]
	if schema == nil {
		return nil
	}
	return schema.getPrimitiveFields(funcType)
}

func (t Type) GetSchema() Schema {
	return schemas[t]
}

// GetDefaultPrimitiveFields returns the default primitive fields for aggregator
func (t Type) GetDefaultPrimitiveFields() PrimitiveFields {
	schema := schemas[t]
	if schema == nil {
		return nil
	}
	return schema.getDefaultPrimitiveFields()
}

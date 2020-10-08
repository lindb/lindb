package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

// AggType represents field's aggregator type.
type AggType uint8

// ID represents field id.
type ID uint16

// Name represents field name.
type Name string

// Defines all aggregator types for field
const (
	Sum AggType = iota + 1
	Count
	Min
	Max
	Replace
)

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

// GetAggFunc returns the aggregate function
func (t Type) GetAggFunc() AggFunc {
	switch t {
	case SumField:
		return sumAggregator
	case MinField:
		return minAggregator
	case MaxField:
		return maxAggregator
	default:
		return nil
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

// GetFuncFieldParams returns the fields for aggregator's function params.
func (t Type) GetFuncFieldParams(funcType function.FuncType) []AggType {
	switch t {
	case SumField:
		return getFieldParamsForSumField(funcType)
	case MinField:
		return getFieldParamsForMinField(funcType)
	}
	return nil
}

func getFieldParamsForSumField(funcType function.FuncType) []AggType {
	switch funcType {
	case function.Max:
		return []AggType{Max}
	default:
		return []AggType{Sum}
	}
}

func getFieldParamsForMinField(funcType function.FuncType) []AggType {
	switch funcType {
	case function.Max:
		return []AggType{Max}
	default:
		return []AggType{Min}
	}
}

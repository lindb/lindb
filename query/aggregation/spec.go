package aggregation

import (
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/function"
)

// dummy value to keep field name unique
const dummy bool = false

type AggregatorStreamSpec struct {
}
type AggregatorSpec struct {
	fieldID   uint16
	fieldName string
	fieldType field.Type
	functions map[function.Type]bool
}

func NewAggregatorSpec(fieldID uint16, fieldName string, fieldType field.Type) *AggregatorSpec {
	return &AggregatorSpec{
		fieldID:   fieldID,
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.Type]bool),
	}
}

func (a *AggregatorSpec) AddFunctionType(funcType function.Type) {
	_, exist := a.functions[funcType]
	if !exist {
		a.functions[funcType] = dummy
	}
}

func DownSamplingFunc(fieldType field.Type) function.Type {
	switch fieldType {
	case field.SumField:
		return function.Sum
	case field.MinField:
		return function.Min
	case field.MaxField:
		return function.Max
	case field.HistogramField:
		return function.Histogram
	default:
		return function.Unknown
	}
}

func IsSupportFunc(fieldType field.Type, funcType function.Type) bool {
	switch fieldType {
	case field.SumField:
		switch funcType {
		case function.Sum, function.Min, function.Max:
			return true
		default:
			return false
		}
	case field.MinField:
		switch funcType {
		case function.Min:
			return true
		default:
			return false
		}
	case field.MaxField:
		switch funcType {
		case function.Max:
			return true
		default:
			return false
		}
	case field.HistogramField:
		return true
	default:
		return false
	}
}

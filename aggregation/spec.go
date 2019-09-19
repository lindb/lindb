package aggregation

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

// dummy value to keep field name unique
const dummy bool = false

type AggregatorStreamSpec struct{}

type AggregatorSpec struct {
	fieldName string
	fieldType field.Type
	functions map[function.FuncType]bool
}

func NewAggregatorSpec(fieldName string, fieldType field.Type) *AggregatorSpec {
	return &AggregatorSpec{
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.FuncType]bool),
	}
}

func (a *AggregatorSpec) FieldName() string {
	return a.fieldName
}

func (a *AggregatorSpec) AddFunctionType(funcType function.FuncType) {
	_, exist := a.functions[funcType]
	if !exist {
		a.functions[funcType] = dummy
	}
}

func DownSamplingFunc(fieldType field.Type) function.FuncType {
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

func IsSupportFunc(fieldType field.Type, funcType function.FuncType) bool {
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

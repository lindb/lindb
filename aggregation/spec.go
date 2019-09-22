package aggregation

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

type AggregatorSpecs []AggregatorSpec

type AggregatorSpec interface {
	FieldName() string
	FieldType() field.Type
	AddFunctionType(funcType function.FuncType)
	Functions() map[function.FuncType]function.FuncType
}

type aggregatorSpec struct {
	fieldName string
	fieldType field.Type
	functions map[function.FuncType]function.FuncType
}

func NewAggregatorSpec(fieldName string, fieldType field.Type) AggregatorSpec {
	return &aggregatorSpec{
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.FuncType]function.FuncType),
	}
}

func (a *aggregatorSpec) FieldType() field.Type {
	return a.fieldType
}

func (a *aggregatorSpec) FieldName() string {
	return a.fieldName
}

func (a *aggregatorSpec) AddFunctionType(funcType function.FuncType) {
	_, exist := a.functions[funcType]
	if !exist {
		a.functions[funcType] = funcType
	}
}

func (a *aggregatorSpec) Functions() map[function.FuncType]function.FuncType {
	return a.functions
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

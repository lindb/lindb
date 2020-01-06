package aggregation

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

type AggregatorSpecs []AggregatorSpec

type AggregatorSpec interface {
	FieldName() string
	GetFieldType() field.Type
	SetFieldType(fieldType field.Type)
	AddFunctionType(funcType function.FuncType)
	Functions() map[function.FuncType]function.FuncType
}

type aggregatorSpec struct {
	fieldName string
	fieldType field.Type
	functions map[function.FuncType]function.FuncType
}

func NewAggregatorSpec(fieldName string) AggregatorSpec {
	return &aggregatorSpec{
		fieldName: fieldName,
		functions: make(map[function.FuncType]function.FuncType),
	}
}

func NewDownSamplingSpec(fieldName string, fieldType field.Type) AggregatorSpec {
	return &aggregatorSpec{
		fieldName: fieldName,
		fieldType: fieldType,
		functions: make(map[function.FuncType]function.FuncType),
	}
}

func (a *aggregatorSpec) GetFieldType() field.Type {
	return a.fieldType
}

func (a *aggregatorSpec) SetFieldType(fieldType field.Type) {
	a.fieldType = fieldType
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

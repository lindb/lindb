package function

import "github.com/lindb/lindb/spi/value"

type ResolvedFunction struct{}

func (f *ResolvedFunction) GetSignature() *BoundSignature {
	return &BoundSignature{
		returnType: &value.RowType{},
	}
}

type FunctionResolver struct {
}

func NewFunctionResolver() *FunctionResolver {
	return &FunctionResolver{}
}

func (r *FunctionResolver) ResolveOperator(operatorType OperatorType, argumentTypes []value.Type) *ResolvedFunction {
	return &ResolvedFunction{}
}

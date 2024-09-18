package function

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ResolvedFunction struct {
	Signature *BoundSignature `json:"signature"`
}

type FunctionResolver struct{}

func NewFunctionResolver() *FunctionResolver {
	return &FunctionResolver{}
}

func (r *FunctionResolver) ResolveOperator(operatorType types.OperatorType, argumentTypes []types.Type) *ResolvedFunction {
	return &ResolvedFunction{
		// FIXME: function name/types
		Signature: NewBoundSignature(operatorType.Operator, types.DTFloat, []types.DataType{types.DTFloat}),
	}
}

func (r *FunctionResolver) ResolveFunction(name *tree.QualifiedName) *ResolvedFunction {
	return &ResolvedFunction{
		// FIXME: function name/types
		Signature: NewBoundSignature(name.Suffix, types.DTFloat, []types.DataType{types.DTFloat}),
	}
}

package function

import "github.com/lindb/lindb/spi/value"

type BoundSignature struct {
	returnType    value.Type
	argumentTypes []value.Type
}

func NewBoundSignature(returnType value.Type, argumentTypes []value.Type) *BoundSignature {
	return &BoundSignature{
		returnType:    returnType,
		argumentTypes: argumentTypes,
	}
}

func (s *BoundSignature) GetReturnType() value.Type {
	return s.returnType
}

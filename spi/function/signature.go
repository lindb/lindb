package function

import "github.com/lindb/lindb/spi/types"

type BoundSignature struct {
	Name          string           `json:"name"`
	ArgumentTypes []types.DataType `json:"argumentTypes"`
	ReturnType    types.DataType   `json:"returnType"`
}

func NewBoundSignature(name string, returnType types.DataType, argumentTypes []types.DataType) *BoundSignature {
	return &BoundSignature{
		Name:          name,
		ReturnType:    returnType,
		ArgumentTypes: argumentTypes,
	}
}

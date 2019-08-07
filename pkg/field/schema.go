package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

type schema interface {
	getPrimitiveFields(funcType function.FuncType) map[uint16]AggType
}

type sumSchema struct {
	primitiveFieldID uint16
}

func newSumSchema() schema {
	return &sumSchema{
		primitiveFieldID: uint16(1),
	}
}

func (s *sumSchema) getPrimitiveFields(funcType function.FuncType) map[uint16]AggType {
	switch funcType {
	case function.Sum:
		return map[uint16]AggType{s.primitiveFieldID: Sum}
	default:
		return nil
	}
}

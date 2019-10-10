package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

type schema interface {
	getPrimitiveFields(funcType function.FuncType) map[uint16]AggType
	getDefaultPrimitiveFields() map[uint16]AggType
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

func (s *sumSchema) getDefaultPrimitiveFields() map[uint16]AggType {
	return map[uint16]AggType{s.primitiveFieldID: Sum}
}

type minSchema struct {
	primitiveFieldID uint16
}

func newMinSchema() schema {
	return &minSchema{
		primitiveFieldID: uint16(1),
	}
}

func (s *minSchema) getPrimitiveFields(funcType function.FuncType) map[uint16]AggType {
	switch funcType {
	case function.Min:
		return map[uint16]AggType{s.primitiveFieldID: Min}
	default:
		return nil
	}
}

func (s *minSchema) getDefaultPrimitiveFields() map[uint16]AggType {
	return map[uint16]AggType{s.primitiveFieldID: Min}
}

type summarySchema struct {
	sumFieldID, countFieldID, minFieldID, maxFieldID uint16
}

func newSummarySchema() schema {
	return &summarySchema{
		sumFieldID:   uint16(1),
		countFieldID: uint16(2),
		maxFieldID:   uint16(3),
		minFieldID:   uint16(4),
	}
}

func (s *summarySchema) getPrimitiveFields(funcType function.FuncType) map[uint16]AggType {
	switch funcType {
	case function.Sum:
		return map[uint16]AggType{s.sumFieldID: Sum}
	case function.Min:
		return map[uint16]AggType{s.minFieldID: Min}
	case function.Max:
		return map[uint16]AggType{s.maxFieldID: Max}
	case function.Count:
		return map[uint16]AggType{s.countFieldID: Sum}
	case function.Avg:
		return map[uint16]AggType{s.sumFieldID: Sum, s.countFieldID: Sum}
	default:
		return nil
	}
}

func (s *summarySchema) getDefaultPrimitiveFields() map[uint16]AggType {
	return map[uint16]AggType{s.countFieldID: Sum}
}

package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

const SimpleFieldPFieldID = uint16(1)

// Schema represents the field schema internal definition
type Schema interface {
	// GetAggFunc gets agg func type by primitive field id
	GetAggFunc(pFieldID uint16) AggFunc
	// GetAllPrimitiveFields returns all primitive field ids, sort by field id
	GetAllPrimitiveFields() []uint16
	// getPrimitiveFields gets need extract primitive fields
	getPrimitiveFields(funcType function.FuncType) PrimitiveFields
	// getDefaultPrimitiveFields gets the default extract primitive fields
	getDefaultPrimitiveFields() PrimitiveFields
}

type sumSchema struct {
	primitiveFieldID uint16
	fieldIDs         []uint16
}

func newSumSchema() Schema {
	return &sumSchema{
		primitiveFieldID: SimpleFieldPFieldID,
		fieldIDs:         []uint16{SimpleFieldPFieldID},
	}
}

func (s *sumSchema) GetAggFunc(pFieldID uint16) AggFunc {
	return sumAggregator
}

func (s *sumSchema) GetAllPrimitiveFields() []uint16 {
	return s.fieldIDs
}

func (s *sumSchema) getPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	switch funcType {
	case function.Sum:
		return s.getDefaultPrimitiveFields()
	default:
		return nil
	}
}

func (s *sumSchema) getDefaultPrimitiveFields() PrimitiveFields {
	return PrimitiveFields{
		{FieldID: s.primitiveFieldID, AggType: Sum},
	}
}

type minSchema struct {
	primitiveFieldID uint16
	fieldIDs         []uint16
}

func newMinSchema() Schema {
	return &minSchema{
		primitiveFieldID: SimpleFieldPFieldID,
		fieldIDs:         []uint16{SimpleFieldPFieldID},
	}
}

func (s *minSchema) GetAggFunc(pFieldID uint16) AggFunc {
	return minAggregator
}

func (s *minSchema) GetAllPrimitiveFields() []uint16 {
	return s.fieldIDs
}

func (s *minSchema) getPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	switch funcType {
	case function.Min:
		return s.getDefaultPrimitiveFields()
	default:
		return nil
	}
}

func (s *minSchema) getDefaultPrimitiveFields() PrimitiveFields {
	return PrimitiveFields{
		{FieldID: s.primitiveFieldID, AggType: Min},
	}
}

type maxSchema struct {
	primitiveFieldID uint16
	fieldIDs         []uint16
}

func newMaxSchema() Schema {
	return &maxSchema{
		primitiveFieldID: SimpleFieldPFieldID,
		fieldIDs:         []uint16{SimpleFieldPFieldID},
	}
}

func (s *maxSchema) GetAggFunc(pFieldID uint16) AggFunc {
	return maxAggregator
}

func (s *maxSchema) GetAllPrimitiveFields() []uint16 {
	return s.fieldIDs
}

func (s *maxSchema) getPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	switch funcType {
	case function.Max:
		return s.getDefaultPrimitiveFields()
	default:
		return nil
	}
}

func (s *maxSchema) getDefaultPrimitiveFields() PrimitiveFields {
	return PrimitiveFields{
		{FieldID: s.primitiveFieldID, AggType: Max},
	}
}

type gaugeSchema struct {
	primitiveFieldID uint16
	fieldIDs         []uint16
}

func newGaugeSchema() Schema {
	return &gaugeSchema{
		primitiveFieldID: SimpleFieldPFieldID,
		fieldIDs:         []uint16{SimpleFieldPFieldID},
	}
}

func (s *gaugeSchema) GetAggFunc(pFieldID uint16) AggFunc {
	return replaceAggregator
}

func (s *gaugeSchema) GetAllPrimitiveFields() []uint16 {
	return s.fieldIDs
}

func (s *gaugeSchema) getPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	switch funcType {
	case function.Replace:
		return s.getDefaultPrimitiveFields()
	default:
		return nil
	}
}

func (s *gaugeSchema) getDefaultPrimitiveFields() PrimitiveFields {
	return PrimitiveFields{
		{FieldID: s.primitiveFieldID, AggType: Replace},
	}
}

type summarySchema struct {
	sumFieldID, countFieldID, minFieldID, maxFieldID uint16
	fieldIDs                                         []uint16
}

func newSummarySchema() Schema {
	return &summarySchema{
		sumFieldID:   uint16(1),
		countFieldID: uint16(2),
		maxFieldID:   uint16(3),
		minFieldID:   uint16(4),
		fieldIDs:     []uint16{1, 2, 3, 4},
	}
}
func (s *summarySchema) GetAggFunc(pFieldID uint16) AggFunc {
	switch pFieldID {
	case uint16(1), uint16(2):
		return sumAggregator
	case uint16(3):
		return maxAggregator
	case uint16(4):
		return minAggregator
	default:
		return replaceAggregator
	}
}

func (s *summarySchema) GetAllPrimitiveFields() []uint16 {
	return s.fieldIDs
}

func (s *summarySchema) getPrimitiveFields(funcType function.FuncType) PrimitiveFields {
	switch funcType {
	case function.Sum:
		return PrimitiveFields{
			{FieldID: s.sumFieldID, AggType: Sum},
		}
	case function.Min:
		return PrimitiveFields{
			{FieldID: s.minFieldID, AggType: Min},
		}
	case function.Max:
		return PrimitiveFields{
			{FieldID: s.maxFieldID, AggType: Max},
		}
	case function.Count:
		return PrimitiveFields{
			{FieldID: s.countFieldID, AggType: Sum},
		}
	case function.Avg:
		return PrimitiveFields{
			{FieldID: s.sumFieldID, AggType: Sum},
			{FieldID: s.countFieldID, AggType: Sum},
		}
	default:
		return nil
	}
}

func (s *summarySchema) getDefaultPrimitiveFields() PrimitiveFields {
	return PrimitiveFields{
		{FieldID: s.countFieldID, AggType: Sum},
	}
}

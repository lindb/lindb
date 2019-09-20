package field

import (
	"github.com/lindb/lindb/aggregation/function"
)

// ValueType represents primitive field's value type
type ValueType uint8

// Defines all value type of primitive field
const (
	Integer ValueType = iota + 1
	Float
)

// AggType represents primitive field's aggregator type
type AggType uint8

// Defines all aggregator types for primitive field
const (
	Sum AggType = iota + 1
	Count
	Min
	Max
)

// Type represents field type for LinDB support
type Type uint8

// Defines all field types for LinDB support(user write)
const (
	SumField Type = iota + 1
	MinField
	MaxField
	HistogramField

	Unknown
)

var schemas = map[Type]schema{}

func init() {
	schemas[SumField] = newSumSchema()
}

// GetPrimitiveFields returns the primitive fields for down sampling
func GetPrimitiveFields(fieldType Type, funcType function.FuncType) map[uint16]AggType {
	schema := schemas[fieldType]
	if schema == nil {
		return nil
	}
	return schema.getPrimitiveFields(funcType)
}

func GetPrimitiveFieldsValue() {

}

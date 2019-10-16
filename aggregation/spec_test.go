package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

func TestAggregatorSpec_FieldName(t *testing.T) {
	agg := NewAggregatorSpec("f1", field.SumField)
	assert.Equal(t, "f1", agg.FieldName())
	assert.Equal(t, field.SumField, agg.FieldType())
}

func TestAggregatorSpec_AddFunctionType(t *testing.T) {
	agg := NewAggregatorSpec("f1", field.SumField)
	agg.AddFunctionType(function.Sum)
	agg.AddFunctionType(function.Sum)
	assert.Equal(t, 1, len(agg.Functions()))
}

package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

func TestAggregatorSpec_FieldName(t *testing.T) {
	agg := NewDownSamplingSpec("f1", field.SumField)
	assert.Equal(t, field.Name("f1"), agg.FieldName())
	assert.Equal(t, field.SumField, agg.GetFieldType())
}

func TestAggregatorSpec_AddFunctionType(t *testing.T) {
	agg := NewDownSamplingSpec("f1", field.SumField)
	agg.AddFunctionType(function.Sum)
	agg.AddFunctionType(function.Sum)
	assert.Equal(t, 1, len(agg.Functions()))

	agg = NewAggregatorSpec("f1")
	assert.Equal(t, 0, int(agg.GetFieldType()))
	agg.SetFieldType(field.SumField)
	assert.Equal(t, field.SumField, agg.GetFieldType())
}

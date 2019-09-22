package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/series/field"
)

func TestDownSamplingFunc(t *testing.T) {
	assert.Equal(t, function.Sum, DownSamplingFunc(field.SumField))
	assert.Equal(t, function.Min, DownSamplingFunc(field.MinField))
	assert.Equal(t, function.Max, DownSamplingFunc(field.MaxField))
	assert.Equal(t, function.Histogram, DownSamplingFunc(field.HistogramField))
	assert.Equal(t, function.Unknown, DownSamplingFunc(field.Unknown))
}

func TestIsSupportFunc(t *testing.T) {
	assert.True(t, IsSupportFunc(field.SumField, function.Sum))
	assert.True(t, IsSupportFunc(field.SumField, function.Min))
	assert.True(t, IsSupportFunc(field.SumField, function.Max))
	assert.False(t, IsSupportFunc(field.SumField, function.Histogram))

	assert.True(t, IsSupportFunc(field.MaxField, function.Max))
	assert.False(t, IsSupportFunc(field.MaxField, function.Histogram))

	assert.True(t, IsSupportFunc(field.MinField, function.Min))
	assert.False(t, IsSupportFunc(field.MinField, function.Histogram))

	assert.True(t, IsSupportFunc(field.HistogramField, function.Min))
	assert.True(t, IsSupportFunc(field.HistogramField, function.Sum))
	assert.True(t, IsSupportFunc(field.HistogramField, function.Max))
	assert.True(t, IsSupportFunc(field.HistogramField, function.Histogram))

	assert.False(t, IsSupportFunc(field.Unknown, function.Histogram))
}

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

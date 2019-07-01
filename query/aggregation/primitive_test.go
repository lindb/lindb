package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/field"
)

func TestPrimitiveSumIntegerAgg(t *testing.T) {
	agg := NewPrimitiveAggregator(field.Integer, field.GetAggFunc("sum"), 10)
	//ignore wrong idx
	agg.Aggregate(-1, 10)
	agg.Aggregate(10, 10)

	agg.Aggregate(1, int64(10))
	agg.Aggregate(1, int64(10))

	agg.Aggregate(2, 100)
	agg.Aggregate(2, int64(10))

	vals, ok := agg.Values().([]int64)
	assert.True(t, ok)
	assert.Equal(t, 10, len(vals))
	assert.Equal(t, int64(20), vals[1])
	assert.Equal(t, int64(110), vals[2])
}

func TestPrimitiveSumFloatAgg(t *testing.T) {
	agg := NewPrimitiveAggregator(field.Float, field.GetAggFunc("sum"), 10)
	//ignore wrong idx
	agg.Aggregate(-1, 10)
	agg.Aggregate(10, 10)
	//ignore wrong value
	agg.Aggregate(2, 100)
	agg.Aggregate(2, int64(10))

	agg.Aggregate(1, 10.0)
	agg.Aggregate(1, 30.0)

	vals, ok := agg.Values().([]float64)
	assert.True(t, ok)
	assert.Equal(t, 10, len(vals))
	assert.Equal(t, 40.0, vals[1])
}

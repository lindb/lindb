package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAggFunc(t *testing.T) {
	assert.NotNil(t, GetAggFunc(Sum))
	assert.NotNil(t, GetAggFunc(Min))
	assert.NotNil(t, GetAggFunc(Max))
	assert.Nil(t, GetAggFunc(1000))
}

func TestSumAgg(t *testing.T) {
	agg := GetAggFunc(Sum)
	assert.Equal(t, int64(100), agg.AggregateInt(1, 99))
	assert.Equal(t, float64(100.0), agg.AggregateFloat(1, 99.0))
}

func TestMinAgg(t *testing.T) {
	agg := GetAggFunc(Min)
	assert.Equal(t, int64(1), agg.AggregateInt(1, 99))
	assert.Equal(t, int64(1), agg.AggregateInt(99, 1))
	assert.Equal(t, float64(1.0), agg.AggregateFloat(1, 99.0))
	assert.Equal(t, float64(1.0), agg.AggregateFloat(99.0, 1))
}

func TestMaxAgg(t *testing.T) {
	agg := GetAggFunc(Max)
	assert.Equal(t, int64(99), agg.AggregateInt(1, 99))
	assert.Equal(t, int64(99), agg.AggregateInt(99, 1))
	assert.Equal(t, float64(99.0), agg.AggregateFloat(1, 99.0))
	assert.Equal(t, float64(99.0), agg.AggregateFloat(99.0, 1))
}

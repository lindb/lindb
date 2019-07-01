package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAggFunc(t *testing.T) {
	assert.NotNil(t, GetAggFunc("sum"))
	assert.NotNil(t, GetAggFunc("max"))
	assert.NotNil(t, GetAggFunc("min"))
	assert.Nil(t, GetAggFunc("sum111"))
}

func TestSumAgg(t *testing.T) {
	agg := GetAggFunc("sum")
	assert.Equal(t, int64(100), agg.AggregateInt(1, 99))
	assert.Equal(t, float64(100.0), agg.AggregateFloat(1, 99.0))
}

func TestMinAgg(t *testing.T) {
	agg := GetAggFunc("min")
	assert.Equal(t, int64(1), agg.AggregateInt(1, 99))
	assert.Equal(t, int64(1), agg.AggregateInt(99, 1))
	assert.Equal(t, float64(1.0), agg.AggregateFloat(1, 99.0))
	assert.Equal(t, float64(1.0), agg.AggregateFloat(99.0, 1))
}

func TestMaxAgg(t *testing.T) {
	agg := GetAggFunc("max")
	assert.Equal(t, int64(99), agg.AggregateInt(1, 99))
	assert.Equal(t, int64(99), agg.AggregateInt(99, 1))
	assert.Equal(t, float64(99.0), agg.AggregateFloat(1, 99.0))
	assert.Equal(t, float64(99.0), agg.AggregateFloat(99.0, 1))
}

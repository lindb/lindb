package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAggFunc(t *testing.T) {
	assert.NotNil(t, Sum.AggFunc())
	assert.NotNil(t, Min.AggFunc())
	assert.NotNil(t, Max.AggFunc())
	assert.NotNil(t, Count.AggFunc())
	assert.NotNil(t, Replace.AggFunc())
	assert.Nil(t, AggType(99).AggFunc())
}

func TestSumAgg(t *testing.T) {
	agg := Sum.AggFunc()
	assert.Equal(t, Sum, agg.AggType())
	assert.Equal(t, 100.0, agg.Aggregate(1, 99.0))
}

func TestMinAgg(t *testing.T) {
	agg := Min.AggFunc()
	assert.Equal(t, Min, agg.AggType())
	assert.Equal(t, 1.0, agg.Aggregate(1, 99.0))
	assert.Equal(t, 1.0, agg.Aggregate(99.0, 1))
}

func TestMaxAgg(t *testing.T) {
	agg := Max.AggFunc()
	assert.Equal(t, Max, agg.AggType())
	assert.Equal(t, 99.0, agg.Aggregate(1, 99.0))
	assert.Equal(t, 99.0, agg.Aggregate(99.0, 1))
}

func TestReplaceAgg(t *testing.T) {
	agg := Replace.AggFunc()
	assert.Equal(t, Replace, agg.AggType())
	assert.Equal(t, 99.0, agg.Aggregate(1, 99.0))
}

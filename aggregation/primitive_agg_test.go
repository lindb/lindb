package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
)

func TestPrimitiveSumFloatAgg(t *testing.T) {
	agg := NewPrimitiveAggregator(1, 10, 5, field.Sum.AggFunc())
	agg.Aggregate(1, 10.0)
	agg.Aggregate(1, 30.0)
	agg.Aggregate(-1, 30.0)
	agg.Aggregate(10, 30.0)

	expect := map[int]float64{11: 40.0}
	it := agg.Iterator()
	assert.Equal(t, uint16(1), agg.FieldID())
	AssertPrimitiveIt(t, it, expect)

	agg.reset()
	it = agg.Iterator()
	assert.False(t, it.HasNext())
	agg.Aggregate(2, 20.0)
	expect = map[int]float64{12: 20.0}
	it = agg.Iterator()
	assert.Equal(t, uint16(1), agg.FieldID())
	AssertPrimitiveIt(t, it, expect)
}

package aggregation

import (
	"testing"

	"github.com/lindb/lindb/pkg/field"
)

func TestPrimitiveSumFloatAgg(t *testing.T) {
	agg := newPrimitiveAggregator(uint16(1), 5, field.GetAggFunc(field.Sum))
	agg.Aggregate(1, 10.0)
	agg.Aggregate(1, 30.0)
	agg.Aggregate(-1, 30.0)
	agg.Aggregate(10, 30.0)

	expect := map[int]float64{1: 40.0}
	it := agg.Iterator()
	AssertPrimitiveIt(t, it, expect)
}

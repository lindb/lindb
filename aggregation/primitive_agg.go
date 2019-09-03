package aggregation

import (
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// PrimitiveAggregator represents a primitive field aggregator under spec field aggregator.
// NOTICE: not-safe for goroutine concurrently
type PrimitiveAggregator interface {
	// Aggregate aggregates value with time slot(index)
	Aggregate(idx int, value float64)
	// Iterator returns an iterator for aggregator results
	Iterator() series.PrimitiveIterator
}

// primitiveAggregator implements primitive aggregator interface, using array for storing aggregation result
type primitiveAggregator struct {
	id         uint16
	values     collections.FloatArray
	pointCount int
	aggFunc    field.AggFunc
}

// newPrimitiveAggregator creates primitive aggregator
func newPrimitiveAggregator(id uint16, pointCount int, aggFunc field.AggFunc) PrimitiveAggregator {
	return &primitiveAggregator{
		id:         id,
		pointCount: pointCount,
		aggFunc:    aggFunc,
	}
}

// Iterator returns an iterator for aggregator results
func (agg *primitiveAggregator) Iterator() series.PrimitiveIterator {
	return newPrimitiveIterator(agg.id, agg.values)
}

// Aggregate aggregates value with time slot(index)
func (agg *primitiveAggregator) Aggregate(idx int, value float64) {
	if idx < 0 || idx >= agg.pointCount {
		return
	}
	if agg.values == nil {
		agg.values = collections.NewFloatArray(agg.pointCount)
	}

	if agg.values.HasValue(idx) {
		agg.values.SetValue(idx, agg.aggFunc.AggregateFloat(agg.values.GetValue(idx), value))
	} else {
		agg.values.SetValue(idx, value)
	}
}

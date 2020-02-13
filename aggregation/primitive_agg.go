package aggregation

import (
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./primitive_agg.go -destination=./primitive_agg_mock.go -package=aggregation

// PrimitiveAggregator represents a primitive field aggregator under spec field aggregator.
// NOTICE: not-safe for goroutine concurrently
type PrimitiveAggregator interface {
	// FieldID returns the primitive field id
	FieldID() field.PrimitiveID
	// Aggregate aggregates value with time slot(index)
	// true: aggregate completed
	Aggregate(idx int, value float64) (completed bool)
	// Iterator returns an iterator for aggregator results
	Iterator() series.PrimitiveIterator

	reset()
}

// primitiveAggregator implements primitive aggregator interface, using array for storing aggregation result
type primitiveAggregator struct {
	id         field.PrimitiveID
	start      int
	values     collections.FloatArray
	pointCount int
	aggFunc    field.AggFunc
}

// newPrimitiveAggregator creates primitive aggregator
func NewPrimitiveAggregator(fieldID field.PrimitiveID, start int, pointCount int, aggFunc field.AggFunc) PrimitiveAggregator {
	return &primitiveAggregator{
		id:         fieldID,
		start:      start,
		pointCount: pointCount,
		aggFunc:    aggFunc,
	}
}

// FieldID returns the primitive field id
func (agg *primitiveAggregator) FieldID() field.PrimitiveID {
	return agg.id
}

// Iterator returns an iterator for aggregator results
func (agg *primitiveAggregator) Iterator() series.PrimitiveIterator {
	return newPrimitiveIterator(agg.id, agg.start, agg.aggFunc.AggType(), agg.values)
}

func (agg *primitiveAggregator) reset() {
	if agg.values != nil {
		agg.values.Reset()
		//agg.values = collections.NewFloatArray(agg.pointCount)
	}
}

// Aggregate aggregates value with time slot(index)
func (agg *primitiveAggregator) Aggregate(idx int, value float64) (completed bool) {
	if idx < 0 {
		return
	}
	if idx >= agg.pointCount {
		return true
	}
	if agg.values == nil {
		agg.values = collections.NewFloatArray(agg.pointCount)
	}

	if agg.values.HasValue(idx) {
		agg.values.SetValue(idx, agg.aggFunc.Aggregate(agg.values.GetValue(idx), value))
	} else {
		agg.values.SetValue(idx, value)
	}
	return
}

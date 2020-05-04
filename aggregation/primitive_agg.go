package aggregation

import (
	"github.com/lindb/lindb/aggregation/selector"
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
	Aggregate(slot int, value float64) (completed bool)
	// Iterator returns an iterator for aggregator results
	Iterator() series.PrimitiveIterator

	// reset resets the aggregate context for reuse
	reset()
}

// primitiveAggregator implements primitive aggregator interface, using array for storing aggregation result
type primitiveAggregator struct {
	id       field.PrimitiveID
	selector selector.SlotSelector
	aggFunc  field.AggFunc
	values   collections.FloatArray
}

// newPrimitiveAggregator creates primitive aggregator
func NewPrimitiveAggregator(fieldID field.PrimitiveID, selector selector.SlotSelector, aggFunc field.AggFunc) PrimitiveAggregator {
	return &primitiveAggregator{
		id:       fieldID,
		selector: selector,
		aggFunc:  aggFunc,
	}
}

// FieldID returns the primitive field id
func (agg *primitiveAggregator) FieldID() field.PrimitiveID {
	return agg.id
}

// Iterator returns an iterator for aggregator results
func (agg *primitiveAggregator) Iterator() series.PrimitiveIterator {
	start, _ := agg.selector.Range()
	return newPrimitiveIterator(agg.id, start, agg.aggFunc.AggType(), agg.values)
}

// reset resets the aggregate context for reuse
func (agg *primitiveAggregator) reset() {
	if agg.values != nil {
		agg.values.Reset()
	}
}

// Aggregate aggregates value with time slot(index), if returns true, aggregate completed
func (agg *primitiveAggregator) Aggregate(slot int, value float64) bool {
	// 1. calc index by time slot via selector
	idx, completed := agg.selector.IndexOf(slot)
	if completed {
		return true
	}

	if idx < 0 {
		// if slot index < 0, returns it
		return false
	}

	if agg.values == nil {
		agg.values = collections.NewFloatArray(agg.selector.PointCount())
	}

	if agg.values.HasValue(idx) {
		// aggregate value with function
		agg.values.SetValue(idx, agg.aggFunc.Aggregate(agg.values.GetValue(idx), value))
	} else {
		// set value
		agg.values.SetValue(idx, value)
	}
	return false
}

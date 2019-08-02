package aggregation

import (
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query/selector"
)

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id
type FieldAggregator interface {
	// TimeRange returns the time range of current aggregator
	TimeRange() timeutil.TimeRange
	// Aggregate aggregates the field series into current aggregator
	Aggregate(it field.Iterator)
	// Iterator returns an iterator for aggregator result
	Iterator() field.Iterator
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	baseTime   int64
	timeRange  timeutil.TimeRange
	interval   int64
	aggregates map[uint16]PrimitiveAggregator
	pointCount int

	aggSpec  *AggregatorSpec
	selector selector.SlotSelector
}

// NewFieldAggregator creates a field aggregator
func NewFieldAggregator(baseTime, interval, start, end int64, intervalRatio int, aggSpec *AggregatorSpec) FieldAggregator {
	agg := &fieldAggregator{
		baseTime:   baseTime,
		interval:   interval,
		pointCount: timeutil.CalPointCount(baseTime+interval*start, baseTime+interval*end, interval),
		aggSpec:    aggSpec,
		aggregates: make(map[uint16]PrimitiveAggregator),
	}

	agg.timeRange = timeutil.TimeRange{Start: baseTime + interval*start, End: baseTime + interval*int64(agg.pointCount)}
	agg.selector = selector.NewIndexSlotSelector(int(start), intervalRatio)

	for funcType := range aggSpec.functions {
		primitiveFields := field.GetPrimitiveFields(aggSpec.fieldType, funcType)
		for id, aggType := range primitiveFields {
			agg.aggregates[id] = newPrimitiveAggregator(id, agg.pointCount, field.GetAggFunc(aggType))
		}
	}

	return agg
}

// TimeRange returns the time range of current aggregator
func (a *fieldAggregator) TimeRange() timeutil.TimeRange {
	return a.timeRange
}

// Iterator returns an iterator for aggregator result
func (a *fieldAggregator) Iterator() field.Iterator {
	its := make([]field.PrimitiveIterator, len(a.aggregates))
	idx := 0
	for _, it := range a.aggregates {
		its[idx] = it.Iterator()
		idx++
	}
	return newFieldIterator(a.aggSpec.fieldID, a.aggSpec.fieldType, its)
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it field.Iterator) {
	slotSelector := a.selector

	for it.HasNext() {
		primitiveIt := it.Next()
		if primitiveIt == nil {
			continue
		}
		primitiveFieldID := primitiveIt.ID()
		aggregator, ok := a.aggregates[primitiveFieldID]
		if !ok {
			continue
		}

		for primitiveIt.HasNext() {
			timeSlot, value := primitiveIt.Next()
			idx := slotSelector.IndexOf(timeSlot)
			if idx < 0 {
				continue
			}
			if idx > a.pointCount {
				break
			}
			aggregator.Aggregate(idx, value)
		}
	}
}

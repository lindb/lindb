package aggregation

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query/selector"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id
type FieldAggregator interface {
	// TimeRange returns the time range of current aggregator
	TimeRange() timeutil.TimeRange
	// Aggregate aggregates the field series into current aggregator
	Aggregate(it series.FieldIterator)
	// Iterator returns an iterator for aggregator result
	Iterator() series.FieldIterator
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	familyStartTime int64
	startSlot       int
	timeRange       timeutil.TimeRange
	interval        int64
	aggregates      map[uint16]PrimitiveAggregator
	pointCount      int

	aggSpec  *AggregatorSpec
	selector selector.SlotSelector
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on base time and interval.
// e.g. family start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(familyStartTime, interval, startIdx, endIdx int64, intervalRatio int, aggSpec *AggregatorSpec) FieldAggregator {
	agg := &fieldAggregator{
		familyStartTime: familyStartTime,
		interval:        interval,
		startSlot:       int(startIdx),
		pointCount:      timeutil.CalPointCount(familyStartTime+interval*startIdx, familyStartTime+interval*endIdx, interval),
		aggSpec:         aggSpec,
		aggregates:      make(map[uint16]PrimitiveAggregator),
	}

	agg.timeRange = timeutil.TimeRange{Start: familyStartTime + interval*startIdx, End: familyStartTime + interval*int64(agg.pointCount)}
	agg.selector = selector.NewIndexSlotSelector(int(startIdx), intervalRatio)

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
func (a *fieldAggregator) Iterator() series.FieldIterator {
	its := make([]series.PrimitiveIterator, len(a.aggregates))
	idx := 0
	for _, it := range a.aggregates {
		its[idx] = it.Iterator()
		idx++
	}
	return newFieldIterator(a.aggSpec.fieldID, a.aggSpec.fieldName, a.aggSpec.fieldType, a.familyStartTime, a.startSlot, its)
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it series.FieldIterator) {
	slotSelector := a.selector

	for it.HasNext() {
		primitiveIt := it.Next()
		if primitiveIt == nil {
			continue
		}
		primitiveFieldID := primitiveIt.FieldID()
		//FIXME stone1100 multi-aggs
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

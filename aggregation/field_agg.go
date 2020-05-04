package aggregation

import (
	"sort"

	"github.com/lindb/lindb/aggregation/selector"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field_agg.go -destination=./field_agg_mock.go -package=aggregation

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id.
type FieldAggregator interface {
	// Aggregate aggregates the field series into current aggregator
	Aggregate(it series.FieldIterator)
	// GetAllAggregators returns all primitive aggregates
	GetAllAggregators() []PrimitiveAggregator
	// ResultSet returns the result set of field aggregator
	ResultSet() (startTime int64, it series.FieldIterator)

	// reset resets the aggregate context for reusing
	reset()
}

// aggKey represents aggregate key for supporting multi aggregate function with a same primitive field id
type aggKey struct {
	primitiveID field.PrimitiveID
	aggType     field.AggType
}

// downSamplingFieldAggregator represents
type downSamplingFieldAggregator struct {
	segmentStartTime int64
	start            int
	aggregators      []PrimitiveAggregator
}

// NewDownSamplingFieldAggregator creates a field aggregator for down sampling,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewDownSamplingFieldAggregator(
	segmentStartTime int64,
	selector selector.SlotSelector,
	aggSpec AggregatorSpec,
) FieldAggregator {
	start, _ := selector.Range()
	agg := &downSamplingFieldAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
	}
	aggregatorMap := make(map[aggKey]PrimitiveAggregator)
	// if down sampling spec need init all aggregator
	for funcType := range aggSpec.Functions() {
		primitiveFields := aggSpec.GetFieldType().GetPrimitiveFields(funcType)
		for _, pField := range primitiveFields {
			key := aggKey{
				primitiveID: pField.FieldID,
				aggType:     pField.AggType,
			}
			aggregatorMap[key] = NewPrimitiveAggregator(pField.FieldID, selector, pField.AggType.AggFunc())
		}
	}
	length := len(aggregatorMap)
	agg.aggregators = make([]PrimitiveAggregator, length)
	idx := 0
	for _, pAgg := range aggregatorMap {
		agg.aggregators[idx] = pAgg
		idx++
	}
	// sort field ids
	sort.Slice(agg.aggregators, func(i, j int) bool {
		return agg.aggregators[i].FieldID() < agg.aggregators[j].FieldID()
	})
	return agg
}

// Aggregate aggregates the field series into current aggregator
func (agg *downSamplingFieldAggregator) Aggregate(it series.FieldIterator) {
	// do nothing for down sampling
}

// GetAllAggregators returns all primitive aggregates
func (agg *downSamplingFieldAggregator) GetAllAggregators() []PrimitiveAggregator {
	return agg.aggregators
}

// ResultSet returns the result set of field aggregator
func (agg *downSamplingFieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	its := make([]series.PrimitiveIterator, len(agg.aggregators))
	idx := 0
	for _, it := range agg.aggregators {
		its[idx] = it.Iterator()
		idx++
	}
	return agg.segmentStartTime, newFieldIterator(agg.start, its)
}

// reset resets the aggregate context for reusing
func (agg *downSamplingFieldAggregator) reset() {
	for _, aggregator := range agg.aggregators {
		aggregator.reset()
	}
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	segmentStartTime int64
	start            int

	aggregateMap map[aggKey]PrimitiveAggregator

	selector selector.SlotSelector
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(segmentStartTime int64, selector selector.SlotSelector) FieldAggregator {
	start, _ := selector.Range()
	agg := &fieldAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
		selector:         selector,
		aggregateMap:     make(map[aggKey]PrimitiveAggregator),
	}

	return agg
}

// ResultSet returns the result set of field aggregator
func (a *fieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	its := make([]series.PrimitiveIterator, len(a.aggregateMap))
	idx := 0
	for _, agg := range a.aggregateMap {
		its[idx] = agg.Iterator()
		idx++
	}
	return a.segmentStartTime, newFieldIterator(a.start, its)
}

// GetAllAggregators returns all primitive aggregates
func (a *fieldAggregator) GetAllAggregators() []PrimitiveAggregator {
	result := make([]PrimitiveAggregator, len(a.aggregateMap))
	idx := 0
	for _, agg := range a.aggregateMap {
		result[idx] = agg
		idx++
	}
	return result
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it series.FieldIterator) {
	for it.HasNext() {
		primitiveIt := it.Next()
		if primitiveIt == nil {
			continue
		}
		primitiveFieldID := primitiveIt.FieldID()
		aggregator := a.getAggregator(primitiveFieldID, primitiveIt.AggType())
		for primitiveIt.HasNext() {
			timeSlot, value := primitiveIt.Next()
			aggregator.Aggregate(timeSlot, value)
		}
	}
}

// getAggregator returns the primitive aggregator by primitive id and function type
func (a *fieldAggregator) getAggregator(primitiveFieldID field.PrimitiveID, aggType field.AggType) PrimitiveAggregator {
	key := aggKey{
		primitiveID: primitiveFieldID,
		aggType:     aggType,
	}
	agg, ok := a.aggregateMap[key]
	if ok {
		return agg
	}
	agg = NewPrimitiveAggregator(primitiveFieldID, a.selector, aggType.AggFunc())
	a.aggregateMap[key] = agg
	return agg
}

// reset resets the aggregate context for reusing
func (a *fieldAggregator) reset() {
	for _, aggregator := range a.aggregateMap {
		aggregator.reset()
	}
}

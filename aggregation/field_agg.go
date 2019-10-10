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
	// GetAllAggregates returns all primitive aggregates
	GetAllAggregates() []PrimitiveAggregator
	// ResultSet returns the result set of field aggregator
	ResultSet() (startTime int64, it series.FieldIterator)
	// reset resets the context for reusing
	reset()
}

type aggKey struct {
	primitiveID uint16
	aggType     field.AggType
}

type downSamplingFieldAggregator struct {
	segmentStartTime int64
	start            int
	aggregates       []PrimitiveAggregator
}

// NewDownSamplingFieldAggregator creates a field aggregator for down sampling,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewDownSamplingFieldAggregator(segmentStartTime int64, selector selector.SlotSelector,
	aggSpec AggregatorSpec) FieldAggregator {
	start, _ := selector.Range()
	agg := &downSamplingFieldAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
	}
	aggregatorMap := make(map[aggKey]PrimitiveAggregator)
	// if down sampling spec need init all aggregator
	for funcType := range aggSpec.Functions() {
		primitiveFields := aggSpec.FieldType().GetPrimitiveFields(funcType)
		for id, aggType := range primitiveFields {
			key := aggKey{
				primitiveID: id,
				aggType:     aggType,
			}
			aggregatorMap[key] = NewPrimitiveAggregator(id, start, selector.PointCount(), field.GetAggFunc(aggType))
		}
	}
	length := len(aggregatorMap)
	agg.aggregates = make([]PrimitiveAggregator, length)
	idx := 0
	for _, pAgg := range aggregatorMap {
		agg.aggregates[idx] = pAgg
		idx++
	}
	// sort field ids
	sort.Slice(agg.aggregates, func(i, j int) bool {
		return agg.aggregates[i].FieldID() < agg.aggregates[j].FieldID()
	})
	return agg
}

func (agg *downSamplingFieldAggregator) Aggregate(it series.FieldIterator) {
	// do nothing for down sampling
}

func (agg *downSamplingFieldAggregator) GetAllAggregates() []PrimitiveAggregator {
	return agg.aggregates
}

func (agg *downSamplingFieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	its := make([]series.PrimitiveIterator, len(agg.aggregates))
	idx := 0
	for _, it := range agg.aggregates {
		its[idx] = it.Iterator()
		idx++
	}
	return agg.segmentStartTime, newFieldIterator(agg.start, its)
}

func (agg *downSamplingFieldAggregator) reset() {
	for _, aggregator := range agg.aggregates {
		aggregator.reset()
	}
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	segmentStartTime int64
	start            int

	aggregateMap map[aggKey]PrimitiveAggregator

	aggSpec  AggregatorSpec
	selector selector.SlotSelector
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(segmentStartTime int64, selector selector.SlotSelector, aggSpec AggregatorSpec) FieldAggregator {
	start, _ := selector.Range()
	agg := &fieldAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
		selector:         selector,
		aggSpec:          aggSpec,
		aggregateMap:     make(map[aggKey]PrimitiveAggregator),
	}

	return agg
}

func (a *fieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	its := make([]series.PrimitiveIterator, len(a.aggregateMap))
	idx := 0
	for _, agg := range a.aggregateMap {
		its[idx] = agg.Iterator()
		idx++
	}
	return a.segmentStartTime, newFieldIterator(a.start, its)
}

func (a *fieldAggregator) GetAllAggregates() []PrimitiveAggregator {
	result := make([]PrimitiveAggregator, len(a.aggregateMap))
	idx := 0
	for _, agg := range a.aggregateMap {
		result[idx] = agg
		idx++
	}
	return result
}

func (a *fieldAggregator) reset() {
	for _, aggregator := range a.aggregateMap {
		aggregator.reset()
	}
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
			idx, completed := a.selector.IndexOf(timeSlot)
			if completed {
				break
			}
			if idx < 0 {
				continue
			}
			aggregator.Aggregate(idx, value)
		}
	}
}

func (a *fieldAggregator) getAggregator(primitiveFieldID uint16, aggType field.AggType) PrimitiveAggregator {
	key := aggKey{
		primitiveID: primitiveFieldID,
		aggType:     aggType,
	}
	agg, ok := a.aggregateMap[key]
	if ok {
		return agg
	}
	start, _ := a.selector.Range()
	agg = NewPrimitiveAggregator(primitiveFieldID, start, a.selector.PointCount(), field.GetAggFunc(aggType))
	a.aggregateMap[key] = agg
	return agg
}

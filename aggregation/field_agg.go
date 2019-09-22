package aggregation

import (
	"sort"

	"github.com/lindb/lindb/aggregation/selector"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field_agg.go -destination=./field_agg_mock.go -package=aggregation

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id
type FieldAggregator interface {
	// Aggregate aggregates the field series into current aggregator
	Aggregate(it series.FieldIterator)

	GetAllAggregates() []PrimitiveAggregator
	ResultSet() (startTime int64, it series.FieldIterator)
	reset()
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	isDownSampling   bool
	segmentStartTime int64
	start            int
	aggregates       []PrimitiveAggregator

	aggregateMap map[aggKey]PrimitiveAggregator

	aggSpec  AggregatorSpec
	selector selector.SlotSelector
}

type aggKey struct {
	primitiveID uint16
	aggType     field.AggType
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(segmentStartTime int64, selector selector.SlotSelector, isDownSampling bool,
	aggSpec AggregatorSpec) FieldAggregator {
	start, _ := selector.Range()
	agg := &fieldAggregator{
		segmentStartTime: segmentStartTime,
		isDownSampling:   isDownSampling,
		start:            start,
		selector:         selector,
		aggSpec:          aggSpec,
	}

	agg.aggregateMap = make(map[aggKey]PrimitiveAggregator)
	// if down sampling spec need init all aggregator
	if isDownSampling {
		for funcType := range aggSpec.Functions() {
			primitiveFields := field.GetPrimitiveFields(aggSpec.FieldType(), funcType)
			for id, aggType := range primitiveFields {
				key := aggKey{
					primitiveID: id,
					aggType:     aggType,
				}
				agg.aggregateMap[key] = NewPrimitiveAggregator(id, agg.selector.PointCount(), field.GetAggFunc(aggType))
			}
		}
		length := len(agg.aggregateMap)
		agg.aggregates = make([]PrimitiveAggregator, length)
		idx := 0
		for _, pAgg := range agg.aggregateMap {
			agg.aggregates[idx] = pAgg
			idx++
		}
		// sort field ids
		sort.Slice(agg.aggregates, func(i, j int) bool {
			return agg.aggregates[i].FieldID() < agg.aggregates[j].FieldID()
		})
	}
	return agg
}

func (a *fieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	if a.isDownSampling {
		its := make([]series.PrimitiveIterator, len(a.aggregates))
		idx := 0
		for _, it := range a.aggregates {
			its[idx] = it.Iterator()
			idx++
		}
		return a.segmentStartTime, newFieldIterator(a.start, its)
	}
	its := make([]series.PrimitiveIterator, len(a.aggregateMap))
	idx := 0
	for _, it := range a.aggregateMap {
		its[idx] = it.Iterator()
		idx++
	}
	return a.segmentStartTime, newFieldIterator(a.start, its)
}
func (a *fieldAggregator) GetAllAggregates() []PrimitiveAggregator {
	return a.aggregates
}
func (a *fieldAggregator) reset() {
	for _, aggregator := range a.aggregates {
		aggregator.reset()
	}
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it series.FieldIterator) {
	if a.isDownSampling {
		return
	}
	for it.HasNext() {
		primitiveIt := it.Next()
		if primitiveIt == nil {
			continue
		}
		primitiveFieldID := primitiveIt.FieldID()
		aggregator, ok := a.getAggregator(primitiveFieldID, primitiveIt.AggType())
		if !ok {
			continue
		}
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

func (a *fieldAggregator) getAggregator(primitiveFieldID uint16, aggType field.AggType) (agg PrimitiveAggregator, ok bool) {
	key := aggKey{
		primitiveID: primitiveFieldID,
		aggType:     aggType,
	}
	agg, ok = a.aggregateMap[key]
	if ok {
		return
	}
	if a.isDownSampling {
		return
	}
	ok = true
	agg = NewPrimitiveAggregator(primitiveFieldID, a.selector.PointCount(), field.GetAggFunc(aggType))
	a.aggregateMap[key] = agg
	return
}

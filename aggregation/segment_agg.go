package aggregation

import (
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

// SegmentAggregator represents a fragment aggregator based on timestamp
// NOTICE: no safe for goroutine
type SegmentAggregator interface {
	// Aggregate aggregates the field series into fragment aggregator based on fragment start time
	Aggregate(it series.FieldIterator)
	// Iterator returns an iterator for aggregator result
	Iterator(tags map[string]string) series.GroupedIterator
}

type baseAggregator struct {
	segmentStartTime   int64
	startSlot, endSlot int64
	intervalVal        int64
	intervalRatio      int

	aggSpecs   map[string]*AggregatorSpec // aggregator spec (down sampling/aggregator)
	aggregates map[string]FieldAggregator
}

// Aggregate aggregates the field series into fragment aggregator based on fragment start time
func (ba *baseAggregator) Aggregate(it series.FieldIterator) {
	if it == nil {
		return
	}
	fieldName := it.FieldMeta().Name
	var agg FieldAggregator
	ok := false
	agg, ok = ba.aggregates[fieldName]
	if !ok {
		agg = NewFieldAggregator(ba.segmentStartTime, ba.intervalVal, ba.startSlot, ba.endSlot,
			ba.intervalRatio, ba.aggSpecs[fieldName])
		ba.aggregates[fieldName] = agg
	}
	agg.Aggregate(it)
}

// Iterator returns an iterator for aggregator result
func (ba *baseAggregator) Iterator(tags map[string]string) series.GroupedIterator {
	return newGroupedIterator(tags, ba.aggregates)
}

type segmentAggregator struct {
	baseAggregator
}

// NewSegmentAggregator creates the segment aggregator based on family start time for one series,
// aggSpecs is down sampling aggregator specs.
func NewSegmentAggregator(queryInterval, storageInterval int64, queryTimeRange *timeutil.TimeRange, familyStartTime int64,
	aggSpecs map[string]*AggregatorSpec) SegmentAggregator {
	// 1. calc interval, default use storageInterval's interval if user not input
	intervalVal := storageInterval
	intervalRatio := 1
	if queryInterval > 0 {
		intervalRatio = timeutil.CalIntervalRatio(queryInterval, intervalVal)
		intervalVal *= int64(intervalRatio)
	}
	// 2. get calculator by interval
	intervalType := interval.CalcIntervalType(intervalVal)
	calc := interval.GetCalculator(intervalType)
	storageTimeRange := &timeutil.TimeRange{
		Start: familyStartTime,
		End:   calc.CalcFamilyEndTime(familyStartTime),
	}
	// 3. calc final query time range
	finalQueryTimeRange := storageTimeRange.Intersect(queryTimeRange)
	// 4. calc start/end index based on storage's start time and interval
	startIdx := int64(calc.CalcSlot(finalQueryTimeRange.Start, familyStartTime, intervalVal))
	endIdx := int64(calc.CalcSlot(finalQueryTimeRange.End, familyStartTime, intervalVal))

	return &segmentAggregator{
		baseAggregator: baseAggregator{
			segmentStartTime: familyStartTime,
			startSlot:        startIdx,
			endSlot:          endIdx,
			intervalRatio:    intervalRatio,
			intervalVal:      intervalVal,
			aggregates:       make(map[string]FieldAggregator),
			aggSpecs:         aggSpecs,
		},
	}
}

type seriesSegmentAggregator struct {
	baseAggregator
}

// NewSeriesSegmentAggregator creates the segment aggregator based on query start time for one series,
// aggSpecs is aggregator specs.
func NewSeriesSegmentAggregator(queryInterval int64, queryTimeRange *timeutil.TimeRange,
	aggSpecs map[string]*AggregatorSpec) SegmentAggregator {
	startIdx := 0
	endIdx := timeutil.CalPointCount(queryTimeRange.Start, queryTimeRange.End, queryInterval)
	return &seriesSegmentAggregator{
		baseAggregator: baseAggregator{
			segmentStartTime: queryTimeRange.Start,
			intervalVal:      queryInterval,
			startSlot:        int64(startIdx),
			endSlot:          int64(endIdx),
			intervalRatio:    1,
			aggSpecs:         aggSpecs,
			aggregates:       make(map[string]FieldAggregator),
		},
	}
}

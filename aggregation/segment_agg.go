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

type segmentAggregator struct {
	familyStartTime    int64
	startSlot, endSlot int64
	intervalVal        int64
	intervalRatio      int

	aggSpecs   map[uint16]*AggregatorSpec
	aggregates map[uint16]FieldAggregator
}

// NewSegmentAggregator creates the segment aggregator based on family start time
func NewSegmentAggregator(queryInterval, storageInterval int64, queryTimeRange *timeutil.TimeRange, familyStartTime int64,
	appSpecs map[uint16]*AggregatorSpec) SegmentAggregator {
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
		familyStartTime: familyStartTime,
		startSlot:       startIdx,
		endSlot:         endIdx,
		intervalRatio:   intervalRatio,
		intervalVal:     intervalVal,
		aggregates:      make(map[uint16]FieldAggregator),
		aggSpecs:        appSpecs,
	}
}

// Aggregate aggregates the field series into fragment aggregator based on fragment start time
func (fa *segmentAggregator) Aggregate(it series.FieldIterator) {
	if it == nil {
		return
	}
	fieldID := it.FieldMeta().ID
	var agg FieldAggregator
	ok := false
	agg, ok = fa.aggregates[fieldID]
	if !ok {
		agg = NewFieldAggregator(fa.familyStartTime, fa.intervalVal, fa.startSlot, fa.endSlot,
			fa.intervalRatio, fa.aggSpecs[fieldID])
		fa.aggregates[fieldID] = agg
	}
	agg.Aggregate(it)
}

// Iterator returns an iterator for aggregator result
func (fa *segmentAggregator) Iterator(tags map[string]string) series.GroupedIterator {
	return newGroupedIterator(tags, fa.aggregates)
}

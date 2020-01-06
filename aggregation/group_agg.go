package aggregation

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./group_agg.go -destination=./group_agg_mock.go -package=aggregation

// GroupingAggregator represents an aggregator which merges time series and does grouping if need
type GroupingAggregator interface {
	// Aggregate aggregates the time series data
	Aggregate(it series.GroupedIterator)
	// ResultSet returns the result set of aggregator
	ResultSet() []series.GroupedIterator
}

type groupingAggregator struct {
	aggSpecs   AggregatorSpecs
	interval   timeutil.Interval
	timeRange  timeutil.TimeRange
	aggregates map[string]FieldAggregates // tag values => field aggregates
}

// NewGroupingAggregator creates a grouping aggregator
func NewGroupingAggregator(
	interval timeutil.Interval,
	timeRange timeutil.TimeRange,
	aggSpecs AggregatorSpecs,
) GroupingAggregator {
	return &groupingAggregator{
		aggSpecs:   aggSpecs,
		interval:   interval,
		timeRange:  timeRange,
		aggregates: make(map[string]FieldAggregates),
	}
}

// Aggregate aggregates the time series data
func (ga *groupingAggregator) Aggregate(it series.GroupedIterator) {
	tags := it.Tags()
	seriesAgg := ga.getAggregator(tags)
	var sAgg SeriesAggregator
	for it.HasNext() {
		seriesIt := it.Next()
		fieldName := seriesIt.FieldName()
		fieldType := seriesIt.FieldType()
		// 1. find field aggregator
		sAgg = nil
		for _, aggregator := range seriesAgg {
			if aggregator.FieldName() == fieldName {
				sAgg = aggregator
				break
			}
		}
		if sAgg == nil {
			continue
		}
		// set field type for aggregate
		sAgg.SetFieldType(fieldType)
		// 2. merge the field series data
		for seriesIt.HasNext() {
			startTime, fieldIt := seriesIt.Next()
			if fieldIt == nil {
				continue
			}
			fAgg, ok := sAgg.GetAggregator(startTime)
			if ok {
				fAgg.Aggregate(fieldIt)
			}
		}
	}
}

// ResultSet returns the result set of aggregator
func (ga *groupingAggregator) ResultSet() []series.GroupedIterator {
	length := len(ga.aggregates)
	if length == 0 {
		return nil
	}
	seriesList := make([]series.GroupedIterator, length)
	idx := 0
	for tags, aggregator := range ga.aggregates {
		seriesList[idx] = aggregator.ResultSet(tags)
		idx++
	}
	return seriesList
}

// getAggregator returns the time series aggregator by time series's tags
func (ga *groupingAggregator) getAggregator(tags string) (agg FieldAggregates) {
	// 2. get series aggregator
	agg, ok := ga.aggregates[tags]
	if !ok {
		agg = NewFieldAggregates(ga.interval, 1, ga.timeRange, false, ga.aggSpecs)
		ga.aggregates[tags] = agg
	}
	return
}

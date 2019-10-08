package aggregation

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source=./group_agg.go -destination=./group_agg_mock.go -package=aggregation

// GroupingAggregator represents an aggregator which merges time series and does grouping if need
type GroupingAggregator interface {
	// Aggregate aggregates the time series data
	Aggregate(it series.GroupedIterator)
	// ResultSet returns the result set of aggregator
	ResultSet() []series.GroupedIterator
}

// timeSeriesAggregator represents the aggregator of a time series
type timeSeriesAggregator struct {
	tags       map[string]string // tags of time series
	aggregator FieldAggregates   // fields aggregator
}

type groupingAggregator struct {
	aggSpecs   AggregatorSpecs
	interval   int64
	timeRange  *timeutil.TimeRange
	aggregates map[string]*timeSeriesAggregator
}

// NewGroupingAggregator creates a grouping aggregator
func NewGroupingAggregator(interval int64, timeRange *timeutil.TimeRange,
	aggSpecs AggregatorSpecs) GroupingAggregator {
	return &groupingAggregator{
		aggSpecs:   aggSpecs,
		interval:   interval,
		timeRange:  timeRange,
		aggregates: make(map[string]*timeSeriesAggregator),
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
		// 1. find field aggregator
		sAgg = nil
		for _, aggregator := range seriesAgg.aggregator {
			if aggregator.FieldName() == fieldName {
				sAgg = aggregator
				break
			}
		}
		if sAgg == nil {
			continue
		}
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
	for _, result := range ga.aggregates {
		seriesList[idx] = result.aggregator.ResultSet(result.tags)
		idx++
	}
	return seriesList
}

// getAggregator returns the time series aggregator by time series's tags
func (ga *groupingAggregator) getAggregator(tags map[string]string) (agg *timeSeriesAggregator) {
	// 1. prepare series tags
	tagsStr := constants.EmptyGroupTagsStr
	if len(tags) > 0 {
		tagsStr = tag.Concat(tags)
	}
	// 2. get series aggregator
	agg, ok := ga.aggregates[tagsStr]
	if !ok {
		agg = &timeSeriesAggregator{
			tags:       tags,
			aggregator: NewFieldAggregates(ga.interval, 1, ga.timeRange, false, ga.aggSpecs),
		}
		ga.aggregates[tagsStr] = agg
	}
	return
}

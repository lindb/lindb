package aggregation

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source=./group_agg.go -destination=./group_agg_mock.go -package=aggregation

type GroupByAggregator interface {
	Aggregate(it series.GroupedIterator)
	Merge(tags map[string]string, agg FieldAggregates)
	ResultSet() []series.GroupedIterator
}

type SeriesAgg struct {
	tags       map[string]string
	aggregator FieldAggregates
}

type groupByAggregator struct {
	aggSpecs       AggregatorSpecs
	isDownSampling bool
	interval       int64
	timeRange      *timeutil.TimeRange
	aggregates     map[string]*SeriesAgg
}

func NewGroupByAggregator(interval int64, timeRange *timeutil.TimeRange,
	isDownSampling bool, aggSpecs AggregatorSpecs) GroupByAggregator {
	return &groupByAggregator{
		aggSpecs:       aggSpecs,
		isDownSampling: isDownSampling,
		interval:       interval,
		timeRange:      timeRange,
		aggregates:     make(map[string]*SeriesAgg),
	}
}

func (ga *groupByAggregator) Merge(tags map[string]string, agg FieldAggregates) {
	seriesAgg := ga.getAggregator(tags)
	for idx, aggregator := range seriesAgg.aggregator {
		for _, segAgg := range agg[idx].Aggregates() {
			if segAgg == nil {
				continue
			}
			startTime, fieldIt := segAgg.ResultSet()
			sAgg, ok := aggregator.GetAggregator(startTime)
			if ok {
				sAgg.Aggregate(fieldIt)
			}
		}
	}
}

func (ga *groupByAggregator) Aggregate(it series.GroupedIterator) {
	tags := it.Tags()
	seriesAgg := ga.getAggregator(tags)
	for _, aggregator := range seriesAgg.aggregator {
		for it.HasNext() {
			seriesIt := it.Next()
			for seriesIt.HasNext() {
				startTime, fieldIt := seriesIt.Next()
				if fieldIt == nil {
					continue
				}
				sAgg, ok := aggregator.GetAggregator(startTime)
				if ok {
					sAgg.Aggregate(fieldIt)
				}
			}
		}
	}
}
func (ga *groupByAggregator) ResultSet() []series.GroupedIterator {
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

func (ga *groupByAggregator) getAggregator(tags map[string]string) (agg *SeriesAgg) {
	// 1. prepare series tags
	tagsStr := constants.EmptyGroupTagsStr
	if len(tags) > 0 {
		tagsStr = tag.Concat(tags)
	}
	// 2. get series aggregator
	agg, ok := ga.aggregates[tagsStr]
	if !ok {
		agg = &SeriesAgg{
			tags:       tags,
			aggregator: NewFieldAggregates(ga.interval, 1, ga.timeRange, ga.isDownSampling, ga.aggSpecs),
		}
		ga.aggregates[tagsStr] = agg
	}
	return
}

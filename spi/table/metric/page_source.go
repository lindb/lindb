package metric

import (
	"fmt"

	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type MetricPageSourceProvider struct{}

func NewMetricPageSourceProvider() *MetricPageSourceProvider {
	return &MetricPageSourceProvider{}
}

func (p *MetricPageSourceProvider) CreatePageSource(table spi.TableHandle, outputs []spi.ColumnMetadata) spi.PageSource {
	return &MetricPageSource{
		table:   table.(*TableHandle),
		outputs: outputs,
		decoder: encoding.GetTSDDecoder(),
	}
}

type MetricPageSource struct {
	table *TableHandle
	split *ScanSplit

	decoder *encoding.TSDDecoder

	outputs []spi.ColumnMetadata
}

func (mps *MetricPageSource) AddSplit(split spi.Split) {
	if metricScanSplit, ok := split.(*ScanSplit); ok {
		mps.split = metricScanSplit
	}
}

func (mps *MetricPageSource) GetNextPage() *spi.Page {
	if mps.split == nil {
		return nil
	}

	defer func() {
		mps.split = nil
	}()

	dataLoadCtx := &flow.DataLoadContext{
		Fields:                mps.split.tableScan.fields,
		LowSeriesIDsContainer: mps.split.LowSeriesIDsContainer,
		SeriesIDHighKey:       mps.split.HighSeriesID,

		TimeRange:            mps.table.TimeRange,
		IntervalRatio:        mps.table.IntervalRatio,
		Interval:             mps.table.Interval,
		IsMultiField:         mps.split.tableScan.fields.Len() > 1,
		IsGrouping:           mps.split.tableScan.hasGrouping(),
		PendingDataLoadTasks: atomic.NewInt32(0),
		Decoder:              mps.decoder,
	}
	dataLoadCtx.DownSamplingSpecs = make(aggregation.AggregatorSpecs, len(dataLoadCtx.Fields))
	dataLoadCtx.AggregatorSpecs = make(aggregation.AggregatorSpecs, len(dataLoadCtx.Fields))
	for i := range dataLoadCtx.Fields {
		a := aggregation.NewAggregatorSpec(dataLoadCtx.Fields[i].Name, field.SumField)
		a.AddFunctionType(function.Sum)
		dataLoadCtx.DownSamplingSpecs[i] = a
		b := aggregation.NewAggregatorSpec(dataLoadCtx.Fields[i].Name, field.SumField)
		b.AddFunctionType(function.Sum)
		dataLoadCtx.AggregatorSpecs[i] = b
	}

	dataLoadCtx.Prepare()

	var loaders []flow.DataLoader
	for i := range mps.split.ResultSet {
		rs := mps.split.ResultSet[i]
		// check series ids if match
		loader := rs.Load(dataLoadCtx)
		if loader != nil {
			loaders = append(loaders, loader)
		}
	}
	if len(loaders) == 0 {
		return nil
	}
	if mps.split.tableScan.hasGrouping() {
		// set collect grouping tag value ids func
		dataLoadCtx.Grouping = mps.split.tableScan.grouping.CollectTagValueIDs
		mps.split.groupingContext.BuildGroup(dataLoadCtx)
	} else {
		dataLoadCtx.PrepareAggregatorWithoutGrouping()
	}

	// for each low series ids
	for _, loader := range loaders {
		var familyTime int64
		// load field series data by series ids
		dataLoadCtx.DownSampling = func(slotRange timeutil.SlotRange, lowSeriesIdx uint16, fieldIdx int, getter encoding.TSDValueGetter) {
			seriesAggregator := dataLoadCtx.GetSeriesAggregator(lowSeriesIdx, fieldIdx)

			agg := seriesAggregator.GetAggregator(familyTime)
			for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
				value, ok := getter.GetValue(movingSourceSlot)
				if !ok {
					// no data, goto next loop
					continue
				}
				agg.AggregateBySlot(int(movingSourceSlot), value)
			}
		}

		// loads the metric data by given series id from load result.
		// if found data need to do down sampling aggregate.
		loader.Load(dataLoadCtx)
	}
	// FIXME: need do agg
	// down sampling
	// reduce aggreator
	fmt.Println("metric source page done")
	reduceAgg := aggregation.NewGroupingAggregator(mps.table.Interval,
		mps.table.IntervalRatio, mps.table.TimeRange, dataLoadCtx.AggregatorSpecs)
	// TODO:
	if dataLoadCtx.IsMultiField {
		fmt.Println(dataLoadCtx.WithoutGroupingSeriesAgg)
		reduceAgg.Aggregate(dataLoadCtx.WithoutGroupingSeriesAgg.Aggregators.ResultSet(""))
	} else {
		if mps.split.tableScan.hasGrouping() {
			for _, groupAgg := range dataLoadCtx.GroupingSeriesAgg {
				reduceAgg.Aggregate(aggregation.FieldAggregates{groupAgg.Aggregator}.ResultSet(groupAgg.Key))
				// TODO: reset
				groupAgg.Aggregator.Reset()
			}
		} else {
			reduceAgg.Aggregate(aggregation.FieldAggregates{dataLoadCtx.WithoutGroupingSeriesAgg.Aggregator}.ResultSet(""))
		}
	}
	// TODO: remove it?
	mps.split.tableScan.grouping.CollectTagValues()

	rs := reduceAgg.ResultSet()
	return mps.buildOutputPage(rs)
}

func (mps *MetricPageSource) buildOutputPage(groupedSeriesList series.GroupedIterators) *spi.Page {
	page := spi.NewPage()
	var (
		fields          []*spi.Column
		grouping        []*spi.Column
		groupingIndexes []int
	)
	for idx, output := range mps.outputs {
		column := spi.NewColumn()
		page.AppendColumn(output, column)
		if lo.ContainsBy(mps.split.tableScan.fields, func(item field.Meta) bool {
			return item.Name.String() == output.Name
		}) {
			fields = append(fields, column)
		} else {
			grouping = append(grouping, column)
			groupingIndexes = append(groupingIndexes, idx)
		}
	}
	// set grouping columns' index
	page.SetGrouping(groupingIndexes)

	hasGrouping := mps.split.tableScan.hasGrouping()
	for _, groupedSeriesItr := range groupedSeriesList {
		for groupedSeriesItr.HasNext() {
			if hasGrouping {
				tagValueIDs := groupedSeriesItr.Tags()
				tags := mps.split.tableScan.grouping.GetTagValues(tagValueIDs)
				for idx, tag := range tags {
					grouping[idx].AppendString(tag)
				}
			}

			seriesItr := groupedSeriesItr.Next()
			fieldIdx := 0
			for seriesItr.HasNext() {
				_, fieldIt := seriesItr.Next()
				for fieldIt.HasNext() {
					pField := fieldIt.Next()

					timeSeries := types.NewTimeSeries(mps.table.TimeRange, mps.table.Interval)

					for pField.HasNext() {
						timestamp, value := pField.Next()
						timeSeries.Put(timestamp, value)
					}

					fields[fieldIdx].AppendTimeSeries(timeSeries)
					fieldIdx++
				}
			}
		}
	}

	return page
}

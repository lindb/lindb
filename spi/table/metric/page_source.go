package metric

import (
	"context"
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

type PageSourceProvider struct{}

func NewPageSourceProvider() spi.PageSourceProvider {
	return &PageSourceProvider{}
}

func (p *PageSourceProvider) CreatePageSource(ctx context.Context, table spi.TableHandle, outputs []types.ColumnMetadata, assignments []*spi.ColumnAssignment) spi.PageSource {
	return &PageSource{
		table:       table.(*TableHandle),
		assignments: assignments,
		outputs:     outputs,
		decoder:     encoding.GetTSDDecoder(),
	}
}

type PageSource struct {
	table       *TableHandle
	assignments []*spi.ColumnAssignment

	split *ScanSplit

	decoder *encoding.TSDDecoder

	outputs []types.ColumnMetadata
}

func (mps *PageSource) AddSplit(split spi.Split) {
	if metricScanSplit, ok := split.(*ScanSplit); ok {
		mps.split = metricScanSplit
	}
}

func (mps *PageSource) GetNextPage() *types.Page {
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
	for i, field := range dataLoadCtx.Fields {
		// build down sampling aggregator spec based on field type
		columnHandles := lo.Filter(mps.assignments, func(item *spi.ColumnAssignment, index int) bool {
			return item.Column == field.Name.String()
		})

		downSamplingSpec := aggregation.NewAggregatorSpec(field.Name, field.Type)
		dataLoadCtx.DownSamplingSpecs[i] = downSamplingSpec

		aggregatorSpec := aggregation.NewAggregatorSpec(field.Name, field.Type)
		dataLoadCtx.AggregatorSpecs[i] = aggregatorSpec

		for _, columnHandle := range columnHandles {
			if handle, ok := columnHandle.Handler.(*ColumnHandle); ok {
				// TODO: react it
				downSamplingSpec.AddFunctionType(function.FuncTypeValueOf(string(handle.Downsampling)))
				aggregatorSpec.AddFunctionType(function.FuncTypeValueOf(string(handle.Aggregation)))
			}
		}
	}

	dataLoadCtx.Prepare()

	if mps.split.tableScan.fields.Len() == 0 {
		fmt.Printf("grouping...%v\n", mps.split.tableScan.grouping.tags)
		dataLoadCtx.Grouping = mps.split.tableScan.grouping.CollectTagValueIDs
		result := make(map[string]struct{})
		dataLoadCtx.CollectGrouping = func(tagValueIDs string, seriesIdxFromQuery uint16) {
			_, ok := result[tagValueIDs]
			if !ok {
				result[tagValueIDs] = struct{}{}
			}
		}
		mps.split.groupingContext.BuildGroup(dataLoadCtx)
		mps.split.tableScan.grouping.CollectTagValues()

		page := types.NewPage()
		var (
			grouping        []*types.Column
			groupingIndexes []int
		)
		for idx, output := range mps.outputs {
			column := types.NewColumn()
			page.AppendColumn(output, column)
			grouping = append(grouping, column)
			groupingIndexes = append(groupingIndexes, idx)
		}
		// set grouping index of the columns
		page.SetGrouping(groupingIndexes)
		for tagValueIDs := range result {
			tags := mps.split.tableScan.grouping.GetTagValues(tagValueIDs)
			for idx, tag := range tags {
				grouping[idx].AppendString(tag)
			}
		}
		return page
	}

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
		result := make(map[string]uint16)
		dataLoadCtx.CollectGrouping = func(tagValueIDs string, seriesIdxFromQuery uint16) {
			// last tag key
			aggIdx, ok := result[tagValueIDs]
			if !ok {
				groupingSeriesAggIdx := dataLoadCtx.NewSeriesAggregator(tagValueIDs)
				aggIdx = groupingSeriesAggIdx
				result[tagValueIDs] = aggIdx
			}
			dataLoadCtx.GroupingSeriesAggRefs[seriesIdxFromQuery] = aggIdx
		}
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

func (mps *PageSource) buildOutputPage(groupedSeriesList series.GroupedIterators) *types.Page {
	page := types.NewPage()
	var (
		fields          []*types.Column
		grouping        []*types.Column
		groupingIndexes []int
	)
	for idx, output := range mps.outputs {
		column := types.NewColumn()
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

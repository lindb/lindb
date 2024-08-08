package metric

import (
	"fmt"
	"reflect"

	commontimeutil "github.com/lindb/common/pkg/timeutil"
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

func (p *MetricPageSourceProvider) CreatePageSource(table spi.TableHandle) spi.PageSource {
	return &MetricPageSource{
		table:   table.(*MetricTableHandle),
		decoder: encoding.GetTSDDecoder(),
	}
}

type MetricPageSource struct {
	table *MetricTableHandle
	split *MetricScanSplit

	decoder *encoding.TSDDecoder
}

func (mps *MetricPageSource) AddSplit(split spi.Split) {
	if metricScanSplit, ok := split.(*MetricScanSplit); ok {
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
		Fields:                mps.split.Fields,
		LowSeriesIDsContainer: mps.split.LowSeriesIDsContainer,
		SeriesIDHighKey:       mps.split.HighSeriesID,
		IntervalRatio:         1,
		Interval:              timeutil.Interval(10 * commontimeutil.OneSecond),
		IsMultiField:          len(mps.split.Fields) > 1,
		// IsGrouping:            shardExecuteCtx.StorageExecuteCtx.Query.HasGroupBy(),
		PendingDataLoadTasks: atomic.NewInt32(0),
		TimeRange:            mps.table.TimeRange,
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

	dataLoadCtx.Grouping()
	dataLoadCtx.PrepareAggregatorWithoutGrouping()

	var loaders []flow.DataLoader
	for i := range mps.split.ResultSet {
		rs := mps.split.ResultSet[i]
		loader := rs.Load(dataLoadCtx)
		if loader != nil {
			loaders = append(loaders, loader)
			fmt.Printf("ident11=%s,loader=%v\n", rs.Identifier(), reflect.TypeOf(loader))
		}
	}
	if len(loaders) == 0 {
		return nil
	}

	min := dataLoadCtx.LowSeriesIDsContainer.Minimum()
	dataLoadCtx.MinSeriesID = min
	dataLoadCtx.MaxSeriesID = dataLoadCtx.LowSeriesIDsContainer.Maximum()
	lengthOfSeriesIDs := int(dataLoadCtx.MaxSeriesID-dataLoadCtx.MinSeriesID) + 1
	dataLoadCtx.LowSeriesIDs = make([]uint16, lengthOfSeriesIDs)
	if dataLoadCtx.IsGrouping {
		dataLoadCtx.GroupingSeriesAggRefs = make([]uint16, lengthOfSeriesIDs)
	}
	it := dataLoadCtx.LowSeriesIDsContainer.PeekableIterator()
	for it.HasNext() {
		lowSeriesID := it.Next()
		seriesIdx := lowSeriesID - min
		dataLoadCtx.LowSeriesIDs[seriesIdx] = lowSeriesID
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
		reduceAgg.Aggregate(dataLoadCtx.WithoutGroupingSeriesAgg.Aggregators.ResultSet(""))
	} else {
		reduceAgg.Aggregate(aggregation.FieldAggregates{dataLoadCtx.WithoutGroupingSeriesAgg.Aggregator}.ResultSet(""))
	}
	rs := reduceAgg.ResultSet()

	return mps.buildOutputPage(rs)
}

func (mps *MetricPageSource) buildOutputPage(groupedSeriesList series.GroupedIterators) *spi.Page {
	page := spi.NewPage()
	// TODO: refact
	for _, groupedSeriesItr := range groupedSeriesList {
		for groupedSeriesItr.HasNext() {
			column := spi.NewColumn()
			// TODO: parse tags
			seriesItr := groupedSeriesItr.Next()
			for seriesItr.HasNext() {
				_, fieldIt := seriesItr.Next()
				for fieldIt.HasNext() {
					pField := fieldIt.Next()

					timeSeries := types.NewTimeSeries(mps.table.TimeRange, mps.table.Interval)

					for pField.HasNext() {
						timestamp, value := pField.Next()
						timeSeries.Put(timestamp, value)
					}

					column.AppendTimeSeries(timeSeries)
				}
			}

			page.AppendColumn(
				spi.NewColumnInfo(string(groupedSeriesItr.Next().FieldName()), types.DataTypeSum), // TODO: set type
				column)
			fmt.Println("add colll")
		}
	}

	fmt.Println("return page.....")
	return page
}

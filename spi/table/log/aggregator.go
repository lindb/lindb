package log

import (
	"fmt"

	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
	logstore "github.com/lindb/lindb/storage/log"
)

type Aggregator interface {
	Initialize()
	Aggregate(output chan<- *types.Page)
}

type aggregatorByTime struct {
	source    *sourceConnector
	tableScan *TableScan
	page      *types.Page

	statsColumn *types.Column
	timeColumn  *types.Column
}

func newAggregatorByTime(source *sourceConnector, tableScan *TableScan) Aggregator {
	return &aggregatorByField{
		source:    source,
		tableScan: tableScan,
		fields:    source.fields,
		fieldKeys: source.fieldKeys,
	}
}

func (agg *aggregatorByTime) Initialize() {
	agg.page = types.NewPage()
	agg.timeColumn = types.NewColumn()
	agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimestamp, Name: "timestamp", Hidden: true},
		agg.timeColumn)
	agg.statsColumn = types.NewColumn()
	agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimeSeries, Name: "count"}, agg.statsColumn)
}

func (agg *aggregatorByTime) Aggregate(output chan<- *types.Page) {
	timeseries := types.NewTimeSeries(agg.tableScan.timeRange, timeutil.Interval(60_000))
	agg.source.findLogs(agg.tableScan, func(segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		segment.FindLogIDsByTimeRange(agg.tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
			logIDsFromStore.And(logIDs)
			pos := int((timestamp - agg.tableScan.timeRange.Start) / 60_000)
			fmt.Printf("--------------logID===>%v=%v,%v\n", pos, logIDsFromStore.GetCardinality(), logIDs.GetCardinality())

			timeseries.Put(pos, timeseries.Get(pos)+float64(logIDsFromStore.GetCardinality()))
		})
		return true
	})
	agg.statsColumn.AppendTimeSeries(timeseries)

	output <- agg.page
}

type aggregatorByField struct {
	source    *sourceConnector
	tableScan *TableScan
	page      *types.Page

	statsColumn *types.Column
	fieldColumn *types.Column
	fields      []string
	fieldKeys   []uint32
}

func newAggregatorByField(source *sourceConnector, tableScan *TableScan) Aggregator {
	return &aggregatorByField{
		source:    source,
		tableScan: tableScan,
		fields:    source.fields,
		fieldKeys: source.fieldKeys,
	}
}

func (agg *aggregatorByField) Initialize() {
	agg.page = types.NewPage()
	agg.statsColumn = types.NewColumn()
	agg.fieldColumn = types.NewColumn()
	_, index, ok := lo.FindIndexOf(agg.source.outputColumns, func(item types.ColumnMetadata) bool {
		return item.DataType == types.DTDynamic
	})
	if !ok {
		agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimeSeries, Name: "count"}, agg.statsColumn)
		return
	}
	if index > 0 {
		agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimeSeries, Name: "count"}, agg.statsColumn)
		agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTDynamic, Name: agg.fields[0]}, agg.fieldColumn)
	} else {
		agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTDynamic, Name: agg.fields[0]}, agg.fieldColumn)
		agg.page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimeSeries, Name: "count"}, agg.statsColumn)
	}
}

func (agg *aggregatorByField) Aggregate(output chan<- *types.Page) {
	stats := uint64(0)
	rows := 0
	var grouping map[string][]uint32
	hasGrouping := len(agg.fieldKeys) == 1
	if hasGrouping {
		grouping = make(map[string][]uint32)
		agg.tableScan.db.IndexDatabase().ScanField(agg.fieldKeys[0], nil, func(key []byte, value uint32) bool {
			grouping[string(key)] = []uint32{value, 0}
			rows++
			// TODO: set limit??
			return rows <= 100
		})
	}
	agg.source.findLogs(agg.tableScan, func(segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		if hasGrouping {
			for _, v := range grouping {
				fieldLogIDs := segment.FindLogIDsByFields([]uint32{v[0]})
				fieldLogIDs.And(logIDs)
				v[1] += uint32(fieldLogIDs.GetCardinality())
			}
		} else {
			stats += logIDs.GetCardinality()
		}
		return true
	})
	if hasGrouping {
		for k, v := range grouping {
			agg.fieldColumn.AppendString(k)
			agg.statsColumn.AppendTimeSeries(types.NewTimeSeriesWithSingleValue(float64(v[1])))
		}
	} else {
		agg.statsColumn.AppendTimeSeries(types.NewTimeSeriesWithSingleValue(float64(stats)))
	}

	output <- agg.page
}

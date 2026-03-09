// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package log

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/arrow/pkg/arrow/builder"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	logstore "github.com/lindb/lindb/storage/log"
)

type Aggregator interface {
	Initialize()
	Aggregate(output chan<- arrow.RecordBatch)
	Close()
}

type aggregatorByTime struct {
	source    *sourceConnector
	tableScan *TableScan

	rb *builder.RecordBuilder
	// statsBuilder *larray.TimeSeriesStructBuilder
}

func newAggregatorByTime(source *sourceConnector, tableScan *TableScan) Aggregator {
	return &aggregatorByTime{
		source:    source,
		tableScan: tableScan,
	}
}

func (agg *aggregatorByTime) Initialize() {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: constants.TimestampColumnName, Type: arrow.FixedWidthTypes.Timestamp_ns, Metadata: arrow.MetadataFrom(map[string]string{"hidden": "true"})},
		{Name: "count", Type: larrow.ExtensionTypes.TimeSeries},
	}, nil)
	agg.rb = builder.NewRecordBuilder(memory.NewGoAllocator(), schema)
	// exBuilder := agg.rb.Fields()[1].(*array.ExtensionBuilder)
	// agg.statsBuilder = larray.NewTimeSeriesStructBuilder(exBuilder)
}

func (agg *aggregatorByTime) Aggregate(output chan<- arrow.RecordBatch) {
	numOfPoints := agg.tableScan.timeRange.NumOfPoints(timeutil.Interval(60_000))
	values := make([]float64, numOfPoints)

	agg.source.findLogs(agg.tableScan, func(segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		segment.FindLogIDsByTimeRange(agg.tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
			logIDsFromStore.And(logIDs)
			pos := int((timestamp - agg.tableScan.timeRange.Start) / 60_000)
			values[pos] += float64(logIDsFromStore.GetCardinality())
		})
		return true
	})

	// agg.statsBuilder.Append(agg.tableScan.timeRange.Start, agg.tableScan.timeRange.End, 60_000, values)

	output <- agg.rb.NewRecord()
}

func (agg *aggregatorByTime) Close() {
	agg.rb.Release()
}

type aggregatorByField struct {
	source    *sourceConnector
	tableScan *TableScan

	rb           *builder.RecordBuilder
	statsBuilder *larray.AggregationBuilder
	fieldBuilder *array.StringBuilder
	fields       []string
	fieldKeys    []uint32
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
	hasGrouping := len(agg.fieldKeys) == 1
	var fields []arrow.Field
	if hasGrouping {
		fields = []arrow.Field{
			{Name: agg.fields[0], Type: arrow.BinaryTypes.String},
			{Name: "count", Type: larrow.ExtensionTypes.Sum},
		}
	} else {
		fields = []arrow.Field{
			{Name: "count", Type: larrow.ExtensionTypes.Sum},
		}
	}
	agg.rb = builder.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema(fields, nil))
	rbFields := agg.rb.Fields()
	if hasGrouping {
		agg.fieldBuilder = rbFields[0].(*array.StringBuilder)
		agg.statsBuilder = larray.NewAggregationBuilder(rbFields[1].(*array.ExtensionBuilder))
	} else {
		agg.statsBuilder = larray.NewAggregationBuilder(rbFields[0].(*array.ExtensionBuilder))
	}
}

func (agg *aggregatorByField) Aggregate(output chan<- arrow.RecordBatch) {
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
			agg.fieldBuilder.Append(k)
			agg.statsBuilder.Append(float64(v[1]))
		}
	} else {
		agg.statsBuilder.Append(float64(stats))
	}

	output <- agg.rb.NewRecord()
}

func (agg *aggregatorByField) Close() {
	agg.rb.Release()
}

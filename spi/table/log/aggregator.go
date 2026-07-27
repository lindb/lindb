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
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/arrow/pkg/arrow/builder"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/timeutil"
	logstore "github.com/lindb/lindb/storage/log"
)

// Aggregator drives a single aggregation pass over a log scan.
// Initialize sets up the output RecordBatch schema and builders.
// Aggregate performs the scan and sends exactly one RecordBatch to output.
// Close releases the underlying Arrow memory.
type Aggregator interface {
	Initialize()
	Aggregate(output chan<- arrow.RecordBatch)
	Close()
}

// fixedInterval is the bucket width used by the time-based aggregator (60 s in ms).
const fixedInterval = int64(60_000)

// aggregatorByTime counts log entries per fixed-width time bucket.
// Output schema: [count (TimeSeries)]
// One row is emitted; the TimeSeries value encodes the per-bucket counts over
// the full query time range.  The timestamp is not emitted as a standalone
// column — time information is already embedded in the TimeSeries struct
// (start, end, interval, values), following the same pattern as the metric connector.
type aggregatorByTime struct {
	source    *sourceConnector
	tableScan *TableScan

	rb           *builder.RecordBuilder
	statsBuilder *larray.TimeSeriesBuilder
}

func newAggregatorByTime(source *sourceConnector, tableScan *TableScan) Aggregator {
	return &aggregatorByTime{
		source:    source,
		tableScan: tableScan,
	}
}

func (agg *aggregatorByTime) Initialize() {
	// Output schema contains only the count column encoded as a TimeSeries.
	// The timestamp is not a separate column — time metadata (start, end, interval)
	// lives inside the TimeSeries struct, matching the metric connector pattern.
	// HashAggregationOperator fills any timestamp grouping-key symbol with nulls.
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "count", Type: larrow.ExtensionTypes.TimeSeries},
	}, nil)
	agg.rb = builder.NewRecordBuilder(memory.NewGoAllocator(), schema)
	// Fields()[0] is the ExtensionBuilder for the TimeSeries column.
	exBuilder := agg.rb.Fields()[0].(*array.ExtensionBuilder)
	agg.statsBuilder = larray.NewTimeSeriesBuilder(exBuilder)
	fmt.Println("aggregatorByTime.Initialize() done")
}

func (agg *aggregatorByTime) Aggregate(output chan<- arrow.RecordBatch) {
	fmt.Printf("aggregatorByTime.Aggregate() called: timeRange=[%d,%d]\n",
		agg.tableScan.timeRange.Start, agg.tableScan.timeRange.End)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("aggregatorByTime.Aggregate() PANIC: %v\n", r)
		}
	}()

	numOfPoints := agg.tableScan.timeRange.NumOfPoints(timeutil.Interval(fixedInterval))
	values := make([]float64, numOfPoints)

	// Accumulate matching log counts into the correct time bucket.
	agg.source.findLogs(agg.tableScan, func(_ *Partition, segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		segment.FindLogIDsByTimeRange(agg.tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
			// Keep only the log IDs that satisfy the WHERE predicate.
			logIDsFromStore.And(logIDs)
			pos := int((timestamp - agg.tableScan.timeRange.Start) / fixedInterval)
			fmt.Printf("  agg bucket: ts=%d pos=%d/%d count=%d\n",
				timestamp, pos, numOfPoints, logIDsFromStore.GetCardinality())
			// Guard bounds: timestamps from walkTimeRange may exceed the query end by one step.
			if pos >= 0 && pos < numOfPoints {
				values[pos] += float64(logIDsFromStore.GetCardinality())
			}
		})
		return true
	})

	// Emit the bucketed counts as a single TimeSeries row.
	// Time metadata (start, end, interval) is embedded in the struct; no separate
	// timestamp column is needed — HashAggregationOperator fills any timestamp
	// grouping-key symbol with nulls, matching the metric connector pattern.
	agg.statsBuilder.Append(agg.tableScan.timeRange.Start, agg.tableScan.timeRange.End, fixedInterval, values)

	results := agg.rb.NewRecord()
	fmt.Println("aggregatorByTime results:", results, values)
	output <- results
}

func (agg *aggregatorByTime) Close() {
	agg.rb.Release()
}

// aggregatorByField counts log entries grouped by a single field value (tag),
// or returns a single total count when no grouping key is requested.
//
// Output schema (with grouping):    [<fieldName> (String), count (Sum)]
// Output schema (without grouping): [count (Sum)]
//
// Each distinct field value becomes one output row.
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
	// hasGrouping is true when exactly one grouping field key was requested.
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
	hasGrouping := len(agg.fieldKeys) == 1

	// grouping maps each field string value → [fieldValueID, accumulatedCount].
	// The slice has exactly 2 elements: index 0 = FK id, index 1 = running count.
	var grouping map[string][]uint32
	if hasGrouping {
		grouping = make(map[string][]uint32)
		// Pre-populate the grouping map with all known values for the field.
		agg.tableScan.db.IndexDatabase().ScanField(agg.fieldKeys[0], nil, func(key []byte, value uint32) bool {
			grouping[string(key)] = []uint32{value, 0}
			rows++
			// TODO: set a configurable limit to avoid unbounded result size.
			return rows <= 100
		})
	}

	agg.source.findLogs(agg.tableScan, func(_ *Partition, segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		if hasGrouping {
			for _, v := range grouping {
				// v[0] = FK id for this field value; v[1] = accumulated count.
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
			// Skip groups with zero count: ScanField enumerates all ever-indexed values
			// across the whole database; values outside the queried time range will have
			// count=0 after the bitmap intersection and must not appear in GROUP BY results.
			if v[1] == 0 {
				continue
			}
			agg.fieldBuilder.Append(k)
			agg.statsBuilder.Append(float64(v[1]))
		}
	} else {
		agg.statsBuilder.Append(float64(stats))
	}

	results := agg.rb.NewRecord()
	fmt.Println("aggregatorByField results:", results, "hasGrouping:", len(agg.fieldKeys) == 1, "grouping:", grouping, "stats:", stats)
	output <- results
}

func (agg *aggregatorByField) Close() {
	agg.rb.Release()
}

// fieldEntry holds the per-group state used by aggregatorByFieldAndTime.
// A pointer to this struct is stored in the grouping map to avoid map-value copy
// issues when mutating the values slice during accumulation.
type fieldEntry struct {
	fvID   uint32    // field value ID used for bitmap lookup in the segment
	values []float64 // per-time-bucket counts indexed by bucket position
}

// aggregatorByFieldAndTime counts log entries grouped by a single field value (tag)
// AND bucketed into fixed-width time intervals.
//
// Output schema: [<fieldName> (String), count (TimeSeries)]
//
// Each distinct field value produces one output row; the count column is a TimeSeries
// encoding per-bucket counts over the full query time range, matching the metric connector
// pattern used by aggregatorByTime.
type aggregatorByFieldAndTime struct {
	source    *sourceConnector
	tableScan *TableScan

	rb           *builder.RecordBuilder
	statsBuilder *larray.TimeSeriesBuilder
	fieldBuilder *array.StringBuilder
	fields       []string
	fieldKeys    []uint32
}

func newAggregatorByFieldAndTime(source *sourceConnector, tableScan *TableScan) Aggregator {
	return &aggregatorByFieldAndTime{
		source:    source,
		tableScan: tableScan,
		fields:    source.fields,
		fieldKeys: source.fieldKeys,
	}
}

func (agg *aggregatorByFieldAndTime) Initialize() {
	// Output schema: [fieldName (String), count (TimeSeries)].
	// The timestamp is not a separate column — time metadata lives inside the TimeSeries struct.
	schema := arrow.NewSchema([]arrow.Field{
		{Name: agg.fields[0], Type: arrow.BinaryTypes.String},
		{Name: "count", Type: larrow.ExtensionTypes.TimeSeries},
	}, nil)
	agg.rb = builder.NewRecordBuilder(memory.NewGoAllocator(), schema)
	rbFields := agg.rb.Fields()
	agg.fieldBuilder = rbFields[0].(*array.StringBuilder)
	exBuilder := rbFields[1].(*array.ExtensionBuilder)
	agg.statsBuilder = larray.NewTimeSeriesBuilder(exBuilder)
}

func (agg *aggregatorByFieldAndTime) Aggregate(output chan<- arrow.RecordBatch) {
	numOfPoints := agg.tableScan.timeRange.NumOfPoints(timeutil.Interval(fixedInterval))

	// grouping maps each field string value → *fieldEntry{fvID, per-bucket counts}.
	// Pointer values ensure in-place mutation of the values slice is visible across iterations.
	grouping := make(map[string]*fieldEntry)
	rows := 0
	agg.tableScan.db.IndexDatabase().ScanField(agg.fieldKeys[0], nil, func(key []byte, value uint32) bool {
		grouping[string(key)] = &fieldEntry{
			fvID:   value,
			values: make([]float64, numOfPoints),
		}
		rows++
		// TODO: set a configurable limit to avoid unbounded result size.
		return rows <= 100
	})

	agg.source.findLogs(agg.tableScan, func(_ *Partition, segment *logstore.Segment, logIDs *roaring.Bitmap) bool {
		for _, entry := range grouping {
			// Find log IDs in this segment that match the field value.
			fieldLogIDs := segment.FindLogIDsByFields([]uint32{entry.fvID})
			fieldLogIDs.And(logIDs) // keep only IDs that also satisfy the WHERE predicate
			if fieldLogIDs.IsEmpty() {
				continue
			}
			// Distribute the matched log IDs into their respective time buckets.
			segment.FindLogIDsByTimeRange(agg.tableScan.timeRange, func(timestamp int64, logIDsFromStore *roaring.Bitmap) {
				logIDsFromStore.And(fieldLogIDs)
				pos := int((timestamp - agg.tableScan.timeRange.Start) / fixedInterval)
				// Guard bounds: timestamps from walkTimeRange may exceed the query end by one step.
				if pos >= 0 && pos < numOfPoints {
					entry.values[pos] += float64(logIDsFromStore.GetCardinality())
				}
			})
		}
		return true
	})

	for k, entry := range grouping {
		// Skip groups whose every bucket is zero: ScanField enumerates all ever-indexed values
		// across the whole database; values outside the queried time range will be all-zero
		// after the bitmap intersection and must not appear in GROUP BY results.
		allZero := true
		for _, v := range entry.values {
			if v > 0 {
				allZero = false
				break
			}
		}
		if allZero {
			continue
		}
		agg.fieldBuilder.Append(k)
		agg.statsBuilder.Append(agg.tableScan.timeRange.Start, agg.tableScan.timeRange.End, fixedInterval, entry.values)
	}

	results := agg.rb.NewRecord()
	fmt.Printf("aggregatorByFieldAndTime results: rows=%d\n", results.NumRows())
	output <- results
}

func (agg *aggregatorByFieldAndTime) Close() {
	agg.rb.Release()
}

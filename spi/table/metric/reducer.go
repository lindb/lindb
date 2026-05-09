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

package metric

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	larrowModel "github.com/lindb/arrow/pkg/model"
	"github.com/lindb/common/models"
	"github.com/lindb/roaring"
)

type reducer struct {
	ctx       *ExecutionContext
	tableScan *TableScan
	inbound   <-chan any // *DataSplit or []*roaring.Bitmap
	outbound  chan<- arrow.RecordBatch

	result map[*GroupingKey][]Result // tags => series data of fields(aggregators)
}

func NewReducer(ctx *ExecutionContext, tableScan *TableScan, inbound <-chan any, outbound chan<- arrow.RecordBatch) *reducer {
	return &reducer{
		ctx:       ctx,
		tableScan: tableScan,
		inbound:   inbound,
		outbound:  outbound,
		result:    make(map[*GroupingKey][]Result),
	}
}

func (r *reducer) Run() {
	for {
		select {
		case <-r.ctx.ctx.Done():
			err := r.ctx.ctx.Err()
			if err != nil {
				panic(err)
			}
			goto END
		case event := <-r.inbound:
			if event == nil {
				goto END
			}
			switch e := event.(type) {
			case *DataSplit:
				r.process(e)
			case []*roaring.Bitmap:
				r.tableScan.grouping.CollectTagValueIDs(e)
			}
		}
	}

END:
	if r.tableScan.isGrouping() {
		r.tableScan.grouping.CollectTagValues()
	}

	r.outbound <- r.buildOutputPage()
}

func (r *reducer) process(split *DataSplit) {
	if r.tableScan.fields.Len() == 0 {
		// find series data(columns)
		r.findSeries(split)
		return
	}

	// merge the data of time series
	split.groupingAgg.ForEach(func(tags *GroupingKey, rs []Result) {
		if _, ok := r.result[tags]; !ok {
			r.result[tags] = rs
		} else {
			panic("need impl merge")
		}
	})
}

func (r *reducer) findSeries(split *DataSplit) {
	tagScanners := split.groupingContext.BuildGroup(split.seriesIDHighKey, split.lowSeriesIDs)
	tagsScanner := NewTagsScanner(tagScanners)
	it := split.lowSeriesIDs.PeekableIterator()
	for it.HasNext() {
		// loop each low series ids
		lowSeriesID := it.Next()
		// build grouping keys
		key := tagsScanner.FindTagValues(lowSeriesID)
		if _, ok := r.result[key]; !ok {
			r.result[key.Clone()] = nil
		}
	}
	tableScan := split.partition.tableScan

	tableScan.grouping.CollectTagValueIDs(tagsScanner.GetTagValueIDs())
}

func (r *reducer) buildOutputPage() arrow.RecordBatch {
	// For metrics, timestamp is a hidden dimension — time info is already embedded
	// in each TimeSeries struct (start/end/interval/values). When timestamp is selected,
	// isTimestampSelected controls the output format of float64 fields (TimeSeries vs
	// scalar AggregationType); no separate timestamp column is emitted.
	var adjustedOutputs []arrow.Field
	for _, f := range r.tableScan.outputs {
		if f.Type.ID() == arrow.TIMESTAMP {
			// Timestamp is implicit in TimeSeries structs; skip as a standalone column.
			continue
		}
		if _, ok := f.Type.(*larray.AggregationType); ok && r.tableScan.isTimestampSelected {
			// Remap AggregationType → TimeSeries when time axis is requested.
			adjustedOutputs = append(adjustedOutputs, arrow.Field{Name: f.Name, Type: larrow.ExtensionTypes.TimeSeries, Nullable: f.Nullable})
		} else {
			adjustedOutputs = append(adjustedOutputs, f)
		}
	}
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema(adjustedOutputs, nil))
	defer rb.Release()

	// classify builders into grouping (string) vs field (timeseries/exemplar/aggregation)
	type fieldBuilder struct {
		idx        int
		isExemplar bool
		isAgg      bool
		tsBuilder  *larray.TimeSeriesBuilder
		exBuilder  *larray.ExemplarBuilder
		aggBuilder *larray.AggregationBuilder
	}
	var (
		groupingBuilders []*array.StringBuilder
		fieldBuilders    []fieldBuilder
	)
	rbIdx := 0 // tracks position in rb (excludes skipped timestamp columns)
	for _, output := range r.tableScan.outputs {
		if output.Type.ID() == arrow.TIMESTAMP {
			// Timestamp is implicit in TimeSeries structs; no standalone column needed.
			continue
		}
		b := rb.Field(rbIdx)
		if arrow.TypeEqual(output.Type, arrow.BinaryTypes.String) {
			groupingBuilders = append(groupingBuilders, b.(*array.StringBuilder))
		} else if arrow.TypeEqual(output.Type, larrow.ExtensionTypes.Exemplar) {
			exB := larray.NewExemplarBuilder(b.(*array.ExtensionBuilder))
			fieldBuilders = append(fieldBuilders, fieldBuilder{idx: rbIdx, isExemplar: true, exBuilder: exB})
		} else {
			// Float64 field: output format depends on whether timestamp is selected.
			// With timestamp → TimeSeries struct{start, end, interval, values}.
			// Without timestamp → scalar AggregationType (single aggregated float64).
			if r.tableScan.isTimestampSelected {
				tsB := larray.NewTimeSeriesBuilder(b.(*array.ExtensionBuilder))
				fieldBuilders = append(fieldBuilders, fieldBuilder{idx: rbIdx, tsBuilder: tsB})
			} else {
				aggB := larray.NewAggregationBuilder(b.(*array.ExtensionBuilder))
				fieldBuilders = append(fieldBuilders, fieldBuilder{idx: rbIdx, isAgg: true, aggBuilder: aggB})
			}
		}
		rbIdx++
	}

	hasGrouping := r.tableScan.isGrouping()
	for tags, seriesData := range r.result {
		if hasGrouping {
			tagValues := r.tableScan.grouping.GetTagValues(*tags)
			for i, tag := range tagValues {
				groupingBuilders[i].Append(tag)
			}
		}
		for fieldIdx, stream := range seriesData {
			fb := fieldBuilders[fieldIdx]
			if stream == nil {
				rb.Field(fb.idx).AppendNull()
				continue
			}
			switch dst := stream.(type) {
			case *result[float64]:
				values := dst.array.Values()
				if fb.isAgg {
					// Aggregation type: write single scalar aggregated value.
					if len(values) == 0 {
						fb.aggBuilder.AppendNull()
					} else {
						fb.aggBuilder.Append(values[0])
					}
				} else {
					// TimeSeries: write start/end/interval and the full values slice.
					if len(values) == 0 {
						fb.tsBuilder.AppendNull()
					} else {
						fb.tsBuilder.Append(r.tableScan.timeRange.Start, r.tableScan.timeRange.End,
							r.tableScan.interval.Int64(), values)
					}
				}
			case *result[*models.Exemplar]:
				for _, ex := range dst.array.Values() {
					if ex == nil {
						fb.exBuilder.AppendNull()
					} else {
						fb.exBuilder.Append(&larrowModel.Exemplar{
							TraceID:  []byte(ex.TraceID),
							SpanID:   []byte(ex.SpanID),
							Duration: ex.Duration,
						})
					}
				}
			}
		}
	}

	return rb.NewRecordBatch()
}

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

package aggregation

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
)

// tsData holds the raw values extracted from a single TimeSeries column cell.
type tsData struct {
	start    int64
	end      int64
	interval int64
	values   []float64
}

// readFloat64 reads the float64 value at rowIdx from col.
// Handles AggregationType (larray.Aggregation) and plain Float64 columns.
func readFloat64(col arrow.Array, rowIdx int) float64 {
	if col == nil || col.IsNull(rowIdx) {
		return 0
	}
	switch a := col.(type) {
	case *larray.Aggregation:
		return a.Value(rowIdx)
	case *array.Float64:
		return a.Value(rowIdx)
	}
	return 0
}

// readTimeSeries extracts the time metadata and values from a TimeSeries column cell.
func readTimeSeries(col arrow.Array, rowIdx int) *tsData {
	ts, ok := col.(*larray.TimeSeries)
	if !ok || ts.IsNull(rowIdx) {
		return nil
	}
	structArr := ts.Storage().(*array.Struct)
	start := structArr.Field(0).(*array.Int64).Value(rowIdx)
	end := structArr.Field(1).(*array.Int64).Value(rowIdx)
	interval := structArr.Field(2).(*array.Int64).Value(rowIdx)
	listArr := structArr.Field(3).(*array.List)
	offsets := listArr.Offsets()
	from, to := int(offsets[rowIdx]), int(offsets[rowIdx+1])
	floats := listArr.ListValues().(*array.Float64)
	values := make([]float64, to-from)
	for i := range values {
		values[i] = floats.Value(from + i)
	}
	return &tsData{start: start, end: end, interval: interval, values: values}
}

// readTimeSeriesSlot reads one time slot's value from a TimeSeries column cell.
// Returns 0 when the column is a scalar (Aggregation type) — used in range queries
// where physical fields may include non-TimeSeries stat columns.
func readTimeSeriesSlot(col arrow.Array, rowIdx, slot int) float64 {
	if col == nil || col.IsNull(rowIdx) {
		return 0
	}
	switch a := col.(type) {
	case *larray.TimeSeries:
		td := readTimeSeries(a, rowIdx)
		if td == nil || slot >= len(td.values) {
			return 0
		}
		return td.values[slot]
	case *larray.Aggregation:
		// Stat columns (sum/count) that didn't get promoted to TimeSeries: use scalar.
		return a.Value(rowIdx)
	case *array.Float64:
		return a.Value(rowIdx)
	}
	return 0
}

// appendFloat64 writes a float64 value to b, which may be a Float64Builder or
// an ExtensionBuilder (for AggregationType output columns).
func appendFloat64(b array.Builder, val float64) {
	switch eb := b.(type) {
	case *array.Float64Builder:
		eb.Append(val)
	case *array.ExtensionBuilder:
		// Guard: only call NewAggregationBuilder when the underlying builder is Float64.
		// If the ExtensionType wraps a non-Float64 builder (e.g. TimeSeries wraps Struct),
		// fall back to AppendNull rather than panicking.
		if _, ok := eb.Builder.(*array.Float64Builder); ok {
			larray.NewAggregationBuilder(eb).Append(val)
		} else {
			eb.AppendNull()
		}
	default:
		b.AppendNull()
	}
}

// appendColumnValue copies the value at rowIdx from src into dst.
// Handles String, Float64, Aggregation, and TimeSeries column types.
func appendColumnValue(dst array.Builder, src arrow.Array, rowIdx int) {
	if src.IsNull(rowIdx) {
		dst.AppendNull()
		return
	}
	switch a := src.(type) {
	case *array.String:
		dst.(*array.StringBuilder).Append(a.Value(rowIdx))
	case *array.Float64:
		appendFloat64(dst, a.Value(rowIdx))
	case *larray.Aggregation:
		appendFloat64(dst, a.Value(rowIdx))
	case *larray.TimeSeries:
		eb, ok := dst.(*array.ExtensionBuilder)
		if !ok {
			dst.AppendNull()
			return
		}
		tsB := larray.NewTimeSeriesBuilder(eb)
		td := readTimeSeries(a, rowIdx)
		if td == nil {
			tsB.AppendNull()
		} else {
			tsB.Append(td.start, td.end, td.interval, td.values)
		}
	default:
		dst.AppendNull()
	}
}

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

package main

import (
	"strings"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildStringRecord builds a simple RecordBatch with two string columns for regression testing.
func buildStringRecord(t *testing.T) arrow.RecordBatch {
	t.Helper()
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "version", Type: arrow.BinaryTypes.String},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	rb.Field(0).(*array.StringBuilder).Append("lindb")
	rb.Field(1).(*array.StringBuilder).Append("v1.0")

	rec := rb.NewRecordBatch()
	return rec
}

// buildTimeSeriesRecord builds a RecordBatch with one string tag column and one TimeSeries column.
func buildTimeSeriesRecord(t *testing.T) arrow.RecordBatch {
	t.Helper()
	// start = 1000 ms, interval = 500 ms, 3 data points → timestamps 1000, 1500, 2000
	tsType := larrow.ExtensionTypes.TimeSeries
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "host", Type: arrow.BinaryTypes.String},
		{Name: "cpu", Type: tsType},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	rb.Field(0).(*array.StringBuilder).Append("server-1")
	tsB := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB.Append(1000, 2000, 500, []float64{0.1, 0.2, 0.3})

	rec := rb.NewRecordBatch()
	return rec
}

// buildMultiTimeSeriesRecord builds a RecordBatch with two TimeSeries columns and one tag column.
func buildMultiTimeSeriesRecord(t *testing.T) arrow.RecordBatch {
	t.Helper()
	tsType := larrow.ExtensionTypes.TimeSeries
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "host", Type: arrow.BinaryTypes.String},
		{Name: "cpu", Type: tsType},
		{Name: "mem", Type: tsType},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	rb.Field(0).(*array.StringBuilder).Append("server-1")
	cpuB := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	cpuB.Append(1000, 2000, 500, []float64{10.0, 20.0, 30.0})
	memB := larray.NewTimeSeriesBuilder(rb.Field(2).(*array.ExtensionBuilder))
	memB.Append(1000, 2000, 500, []float64{60.0, 70.0, 80.0})

	rec := rb.NewRecordBatch()
	return rec
}

// buildNullTimeSeriesRecord builds a RecordBatch where the TimeSeries column is null.
func buildNullTimeSeriesRecord(t *testing.T) arrow.RecordBatch {
	t.Helper()
	tsType := larrow.ExtensionTypes.TimeSeries
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "host", Type: arrow.BinaryTypes.String},
		{Name: "cpu", Type: tsType},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	rb.Field(0).(*array.StringBuilder).Append("server-1")
	tsB := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB.AppendNull()

	rec := rb.NewRecordBatch()
	return rec
}

// buildMultiRowTimeSeriesRecord builds a two-row RecordBatch with one TimeSeries column.
func buildMultiRowTimeSeriesRecord(t *testing.T) arrow.RecordBatch {
	t.Helper()
	tsType := larrow.ExtensionTypes.TimeSeries
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "host", Type: arrow.BinaryTypes.String},
		{Name: "cpu", Type: tsType},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	// Row 0: server-1 with 2 data points
	rb.Field(0).(*array.StringBuilder).Append("server-1")
	tsB0 := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB0.Append(1000, 1500, 500, []float64{0.5, 0.6})

	// Row 1: server-2 with 2 data points
	rb.Field(0).(*array.StringBuilder).Append("server-2")
	tsB1 := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB1.Append(2000, 2500, 500, []float64{0.7, 0.8})

	rec := rb.NewRecordBatch()
	return rec
}

func Test_toTable_noTimeSeries(t *testing.T) {
	rec := buildStringRecord(t)
	defer rec.Release()

	out := toTable(rec)
	// Should contain column headers and the single row.
	assert.Contains(t, out, "name")
	assert.Contains(t, out, "version")
	assert.Contains(t, out, "lindb")
	assert.Contains(t, out, "v1.0")
	// Must NOT contain a timestamp column.
	assert.NotContains(t, out, "timestamp")
}

func Test_toTable_singleTimeSeries(t *testing.T) {
	rec := buildTimeSeriesRecord(t)
	defer rec.Release()

	out := toTable(rec)
	// Must have a timestamp column as the first header.
	assert.Contains(t, out, "timestamp")
	// The tag column must be present.
	assert.Contains(t, out, "host")
	assert.Contains(t, out, "server-1")
	// The TimeSeries column header must appear.
	assert.Contains(t, out, "cpu")
	// Values should be present.
	assert.Contains(t, out, "0.1")
	assert.Contains(t, out, "0.2")
	assert.Contains(t, out, "0.3")
	// Three data points → 3 output rows (server-1 repeated).
	require.Equal(t, 3, strings.Count(out, "server-1"), "expected one row per data point")
}

func Test_toTable_multiTimeSeries_singleTimestampColumn(t *testing.T) {
	rec := buildMultiTimeSeriesRecord(t)
	defer rec.Release()

	out := toTable(rec)
	// Only one "timestamp" header, not two.
	require.Equal(t, 1, strings.Count(out, "timestamp"), "expected exactly one timestamp column")
	// Both metric columns must appear.
	assert.Contains(t, out, "cpu")
	assert.Contains(t, out, "mem")
	// All data point values must appear.
	assert.Contains(t, out, "10")
	assert.Contains(t, out, "60")
	// 3 data points → "server-1" appears 3 times.
	require.Equal(t, 3, strings.Count(out, "server-1"), "expected one row per data point")
}

func Test_toTable_nullTimeSeries(t *testing.T) {
	rec := buildNullTimeSeriesRecord(t)
	defer rec.Release()

	out := toTable(rec)
	// Null TimeSeries still produces the timestamp header.
	assert.Contains(t, out, "timestamp")
	// The TS column shows "null" but the tag column still shows its real value.
	assert.Contains(t, out, "null")
	assert.Contains(t, out, "server-1", "non-TS column should retain its actual value even when TS is null")
}

func Test_toTable_multiRow_timeSeries(t *testing.T) {
	rec := buildMultiRowTimeSeriesRecord(t)
	defer rec.Release()

	out := toTable(rec)
	// Both host values must appear.
	assert.Contains(t, out, "server-1")
	assert.Contains(t, out, "server-2")
	// 2 rows × 2 data points each = 4 expanded rows total.
	require.Equal(t, 2, strings.Count(out, "server-1"), "server-1 should appear once per data point")
	require.Equal(t, 2, strings.Count(out, "server-2"), "server-2 should appear once per data point")
	// All values should appear.
	assert.Contains(t, out, "0.5")
	assert.Contains(t, out, "0.8")
}

func Test_toTable_timeSeries_sortedByTimestamp(t *testing.T) {
	// Build a two-row RecordBatch where server-2 starts earlier than server-1.
	// After sorting by timestamp, server-2's rows must appear before server-1's rows.
	tsType := larrow.ExtensionTypes.TimeSeries
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "host", Type: arrow.BinaryTypes.String},
		{Name: "cpu", Type: tsType},
	}, nil)
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()

	// Row 0: server-1 starts at a later time (t=2000 ms)
	rb.Field(0).(*array.StringBuilder).Append("server-1")
	tsB0 := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB0.Append(2000, 2000, 0, []float64{9.9})

	// Row 1: server-2 starts at an earlier time (t=1000 ms)
	rb.Field(0).(*array.StringBuilder).Append("server-2")
	tsB1 := larray.NewTimeSeriesBuilder(rb.Field(1).(*array.ExtensionBuilder))
	tsB1.Append(1000, 1000, 0, []float64{1.1})

	rec := rb.NewRecordBatch()
	defer rec.Release()

	out := toTable(rec)
	// server-2 (t=1000) must appear before server-1 (t=2000) in the output.
	pos1 := strings.Index(out, "server-1")
	pos2 := strings.Index(out, "server-2")
	require.True(t, pos2 < pos1, "server-2 (earlier timestamp) should appear before server-1")
}

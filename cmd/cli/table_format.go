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
	"fmt"
	"sort"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/mattn/go-runewidth"

	"github.com/lindb/lindb/pkg/terminal"
)

// toTable renders an Arrow RecordBatch as a formatted table string.
// Column headers come from the schema field names; values are formatted
// according to each field's Arrow data type.
// If the schema contains TimeSeries columns, the output is expanded so
// that each data point becomes its own row, with a leading "timestamp" column.
func toTable(rs arrow.RecordBatch) string {
	// Detect TimeSeries columns; collect their indices.
	var tsColIndices []int
	for i, f := range rs.Schema().Fields() {
		if arrow.TypeEqual(f.Type, larrow.ExtensionTypes.TimeSeries) {
			tsColIndices = append(tsColIndices, i)
		}
	}
	if len(tsColIndices) > 0 {
		return toTableExpandTimeSeries(rs, tsColIndices)
	}

	schema := rs.Schema()
	numCols := int(schema.NumFields())
	numRows := int(rs.NumRows())

	// Build header row from schema field names.
	header := make(table.Row, numCols)
	maxWidths := make([]int, numCols)
	for i, f := range schema.Fields() {
		header[i] = f.Name
		maxWidths[i] = len(f.Name)
	}

	// Extract all rows.
	rows := make([]table.Row, numRows)
	for rowIdx := range numRows {
		row := make(table.Row, numCols)
		for colIdx := range numCols {
			col := rs.Column(colIdx)
			val := recordCellValue(col, rowIdx)
			row[colIdx] = val
			w := runewidth.StringWidth(fmt.Sprintf("%v", val))
			if w > maxWidths[colIdx] {
				maxWidths[colIdx] = w
			}
		}
		rows[rowIdx] = row
	}

	// Render with go-pretty.
	w := table.NewWriter()
	w.SetStyle(terminal.TableSylte())
	w.AppendHeader(header)
	w.AppendRows(rows)

	// Compute per-column max widths for terminal-aware truncation.
	termWidth := terminal.GetTerminalWidth()
	colConfigs := make([]table.ColumnConfig, numCols)
	remainingWidth := termWidth - numCols
	colWidth := remainingWidth / numCols
	for i := range colConfigs {
		w := colWidth
		if maxWidths[i] < w {
			w = maxWidths[i]
		}
		colConfigs[i] = table.ColumnConfig{
			WidthMax:         w,
			WidthMaxEnforcer: text.WrapSoft,
		}
	}
	w.SetColumnConfigs(colConfigs)
	return w.Render()
}

// toTableExpandTimeSeries renders a RecordBatch that contains TimeSeries columns
// by expanding each data point into its own output row.
// A single "timestamp" column is prepended (shared by all TimeSeries columns).
// Non-TimeSeries columns repeat their value for every expanded row.
func toTableExpandTimeSeries(rs arrow.RecordBatch, tsColIndices []int) string {
	schema := rs.Schema()
	numCols := int(schema.NumFields())
	numRows := int(rs.NumRows())

	// Build header: "timestamp" first, then original column names.
	outCols := numCols + 1
	header := make(table.Row, outCols)
	header[0] = "timestamp"
	maxWidths := make([]int, outCols)
	maxWidths[0] = len("timestamp")
	for i, f := range schema.Fields() {
		header[i+1] = f.Name
		maxWidths[i+1] = len(f.Name)
	}

	// Build a set for O(1) lookup of TimeSeries column indices.
	tsSet := make(map[int]bool, len(tsColIndices))
	for _, idx := range tsColIndices {
		tsSet[idx] = true
	}

	var rows []table.Row

	for rowIdx := range numRows {
		// Determine the number of data points from the first TimeSeries column.
		firstTSCol, ok := rs.Column(tsColIndices[0]).(*larray.TimeSeries)
		if !ok || firstTSCol.IsNull(rowIdx) {
			// Null (or unexpected type) TimeSeries: emit one row.
			// Non-TimeSeries columns still show their actual values.
			row := make(table.Row, outCols)
			row[0] = "null"
			for colIdx := range numCols {
				col := rs.Column(colIdx)
				if tsSet[colIdx] {
					row[colIdx+1] = "null"
				} else {
					row[colIdx+1] = recordCellValue(col, rowIdx)
				}
			}
			updateMaxWidths(maxWidths, row)
			rows = append(rows, row)
			continue
		}

		start := firstTSCol.Start(rowIdx)
		interval := firstTSCol.Interval(rowIdx)
		values := firstTSCol.Values(rowIdx)
		numPoints := len(values)
		// A series with no data points is skipped — no rows are emitted for it.
		if numPoints == 0 {
			continue
		}

		// Expand: one output row per data point.
		for pos := range numPoints {
			row := make(table.Row, outCols)

			// timestamp column: start + interval * pos (milliseconds)
			ts := start + interval*int64(pos)
			row[0] = timeutil.FormatTimestamp(ts, timeutil.DataTimeFormat2)

			// Fill remaining columns.
			for colIdx := range numCols {
				col := rs.Column(colIdx)
				if tsSet[colIdx] {
					// TimeSeries column: extract the value at this position.
					tsCol, ok := col.(*larray.TimeSeries)
					if !ok || tsCol.IsNull(rowIdx) {
						row[colIdx+1] = "null"
					} else {
						vals := tsCol.Values(rowIdx)
						// Guard against this column having fewer points than the first column.
						if pos < len(vals) {
							row[colIdx+1] = vals[pos]
						} else {
							row[colIdx+1] = "null"
						}
					}
				} else {
					// Non-TimeSeries column: repeat the scalar value for every row.
					row[colIdx+1] = recordCellValue(col, rowIdx)
				}
			}
			updateMaxWidths(maxWidths, row)
			rows = append(rows, row)
		}
	}

	// Sort all expanded rows by timestamp (column 0) ascending.
	// "2006-01-02 15:04:05" format is lexicographically ordered, so string comparison is safe.
	// Rows with a null timestamp are sorted to the end.
	sort.SliceStable(rows, func(i, j int) bool {
		ti, iok := rows[i][0].(string)
		tj, jok := rows[j][0].(string)
		if !iok {
			return false // null sorts last
		}
		if !jok {
			return true
		}
		return ti < tj
	})

	// Render with go-pretty.
	w := table.NewWriter()
	w.SetStyle(terminal.TableSylte())
	w.AppendHeader(header)
	w.AppendRows(rows)

	// Compute per-column max widths for terminal-aware truncation.
	termWidth := terminal.GetTerminalWidth()
	colConfigs := make([]table.ColumnConfig, outCols)
	remainingWidth := termWidth - outCols
	colWidth := remainingWidth / outCols
	for i := range colConfigs {
		cw := colWidth
		if maxWidths[i] < cw {
			cw = maxWidths[i]
		}
		colConfigs[i] = table.ColumnConfig{
			WidthMax:         cw,
			WidthMaxEnforcer: text.WrapSoft,
		}
	}
	w.SetColumnConfigs(colConfigs)
	return w.Render()
}

// updateMaxWidths updates the per-column max display widths based on a newly built row.
func updateMaxWidths(maxWidths []int, row table.Row) {
	for i, val := range row {
		w := runewidth.StringWidth(fmt.Sprintf("%v", val))
		if w > maxWidths[i] {
			maxWidths[i] = w
		}
	}
}

// recordCellValue extracts a display value from an Arrow array at the given row.
func recordCellValue(col arrow.Array, row int) any {
	if col.IsNull(row) {
		return "null"
	}
	switch c := col.(type) {
	case *array.String:
		return c.Value(row)
	case *array.Int64:
		return c.Value(row)
	case *array.Int32:
		return c.Value(row)
	case *array.Float64:
		return c.Value(row)
	case *array.Boolean:
		return c.Value(row)
	case *array.Timestamp:
		return c.Value(row).ToTime(arrow.Nanosecond).Format("2006-01-02 15:04:05.000 -07:00 MST")
	case *array.Duration:
		return time.Duration(c.Value(row)).String()
	default:
		return fmt.Sprintf("%v", col.ValueStr(row))
	}
}

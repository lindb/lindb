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
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-runewidth"

	"github.com/lindb/lindb/pkg/terminal"
)

// toTable renders an Arrow RecordBatch as a formatted table string.
// Column headers come from the schema field names; values are formatted
// according to each field's Arrow data type.
func toTable(rs arrow.RecordBatch) string {
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

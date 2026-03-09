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

package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	larrow "github.com/lindb/arrow/pkg/arrow"
	commonmodels "github.com/lindb/common/models"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/mapstructure"

	"github.com/lindb/lindb/pkg/terminal"
	"github.com/lindb/lindb/spi/types"
)

type ResultSet struct {
	Schema *arrow.Schema `json:"-"`
	Rows   [][]any       `json:"rows,omitempty"`

	Error string `json:"-"`
}

func (rs *ResultSet) MarshalJSON() ([]byte, error) {
	type alias struct {
		Schema []byte  `json:"schema,omitempty"`
		Rows   [][]any `json:"rows,omitempty"`
	}
	a := alias{Rows: rs.Rows}
	if rs.Schema != nil {
		b, err := larrow.MarshalSchema(rs.Schema)
		if err != nil {
			return nil, err
		}
		a.Schema = b
	}
	return json.Marshal(a)
}

func (rs *ResultSet) UnmarshalJSON(data []byte) error {
	fmt.Println(string(data))
	type alias struct {
		Schema []byte  `json:"schema,omitempty"`
		Rows   [][]any `json:"rows,omitempty"`
	}
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	rs.Rows = a.Rows
	if len(a.Schema) > 0 {
		schema, err := larrow.UnmarshalSchema(a.Schema)
		if err != nil {
			return err
		}
		rs.Schema = schema
	}
	return nil
}

func NewResultSet() *ResultSet {
	return &ResultSet{}
}

// ToTable returns stateless node list as table if it has value, else return empty string.
func (rs *ResultSet) ToTable() (tableStr string) {
	if rs.Schema == nil {
		return
	}
	fmt.Println(rs.Schema)
	writer := commonmodels.NewTableFormatter()
	writer.SetStyle(terminal.TableSylte())
	var headers table.Row
	var columnTypes []arrow.DataType
	var (
		hasTimeSeries bool
		timeSeriesIdx int
		dataPoints    int
		rows          []table.Row
	)
	var maxWidths []int
	columns := rs.Schema.Fields()
	for i, col := range columns {
		if !hasTimeSeries && arrow.TypeEqual(col.Type, larrow.ExtensionTypes.TimeSeries) {
			timeSeriesIdx = i
			hasTimeSeries = true
			timeSeries := &types.TimeSeries{}
			_ = mapstructure.Decode(rs.Rows[0][i], timeSeries)
			dataPoints = len(timeSeries.Values)
			headers = append(headers, "timestamp") // add timestamp column
			columnTypes = append(columnTypes, arrow.FixedWidthTypes.Timestamp_ns)
			maxWidths = append(maxWidths, len("timestamp"))
		}
		headers = append(headers, col.Name)
		columnTypes = append(columnTypes, col.Type)
		maxWidths = append(maxWidths, len(col.Name))
	}
	writer.AppendHeader(headers)

	for _, row := range rs.Rows {
		if hasTimeSeries {
			// has time series, build row based on data points
			for pos := range dataPoints {
				cols := make(table.Row, len(columns)+1) // add timestamp column
				colIdx := 0
				for i, col := range row {
					if arrow.TypeEqual(columns[i].Type, larrow.ExtensionTypes.TimeSeries) {
						timeSeries := &types.TimeSeries{}
						_ = mapstructure.Decode(col, timeSeries)
						if timeSeriesIdx == i {
							cols[colIdx] = timeutil.FormatTimestamp(timeSeries.TimeRange.Start+timeSeries.Interval*int64(pos),
								timeutil.DataTimeFormat2)
							maxWidths[colIdx] = stringWidth(maxWidths[colIdx], cols[colIdx])
							colIdx++
							cols[colIdx] = timeSeries.Values[pos]
							maxWidths[colIdx] = stringWidth(maxWidths[colIdx], cols[colIdx])
						} else {
							cols[colIdx] = timeSeries.Values[pos]
							maxWidths[colIdx] = stringWidth(maxWidths[colIdx], cols[colIdx])
						}
					} else {
						appendColumn(cols, columnTypes[colIdx], col, colIdx)
						maxWidths[colIdx] = stringWidth(maxWidths[colIdx], cols[colIdx])
					}
					colIdx++
				}
				rows = append(rows, cols)
			}
		} else {
			cols := make(table.Row, len(columns))
			for colIdx, col := range row {
				appendColumn(cols, columnTypes[colIdx], col, colIdx)
				maxWidths[colIdx] = stringWidth(maxWidths[colIdx], cols[colIdx])
			}
			rows = append(rows, cols)
		}
	}
	writer.AppendRows(rows)
	writer.SetColumnConfigs(columnStyles(maxWidths))
	return writer.Render()
}

// appendColumn appends column value to row.
func appendColumn(row table.Row, colType arrow.DataType, col any, index int) {
	if col == nil {
		row[index] = "null"
		return
	}
	fmt.Printf("appendColumn,%v,%T,%v\n", col, col, colType)
	// FIXME: check type???
	switch colType.ID() {
	case arrow.BinaryTypes.String.ID():
		row[index] = strings.ReplaceAll(col.(string), "\t", "  ") // replace tab with space
	case arrow.FixedWidthTypes.Duration_ns.ID():
		row[index] = time.Duration(col.(float64))
	case arrow.PrimitiveTypes.Int64.ID(), arrow.PrimitiveTypes.Float64.ID(), arrow.PrimitiveTypes.Int32.ID():
		row[index] = fmt.Sprintf("%v", col)
	case arrow.FixedWidthTypes.Timestamp_ns.ID():
		fmt.Printf("dataPoints,%v,%T\n", col, col)
		switch val := col.(type) {
		case string:
			row[index] = val
		case int64:
			row[index] = timeutil.FormatTimestamp(val/1000_1000, timeutil.DataTimeFormat2)
		case float64:
			row[index] = timeutil.FormatTimestamp(int64(val)/1000_000, timeutil.DataTimeFormat2)
		}
	// case &arrow.MapType{}.ID():
	// 	// fixme: need fix
	// 	m := col.(map[string]any)
	// 	var sb strings.Builder
	// 	for key, value := range m {
	// 		fmt.Fprintf(&sb, "%s:%v\n", key, value)
	// 	}
	// 	row[index] = sb.String()
	case larrow.ExtensionTypes.Exemplar.ID():
		exemplars := col.([]any)
		var values []string
		for _, exemplarData := range exemplars {
			if exemplarData != nil {
				exemplar := &commonmodels.Exemplar{}
				_ = mapstructure.Decode(exemplarData, exemplar)
				values = append(values, fmt.Sprintf("%s:%s@%v", exemplar.TraceID, exemplar.SpanID, exemplar.Duration))
			}
		}
		row[index] = fmt.Sprintf("[%s]", strings.Join(values, ", "))
	}
}

func stringWidth(width int, v any) int {
	return max(width, runewidth.StringWidth(fmt.Sprintf("%v", v)))
}

func columnStyles(maxWidths []int) []table.ColumnConfig {
	// get terminal width
	terminalWidth := terminal.GetTerminalWidth()
	// calculate the width of each column
	numCols := len(maxWidths)
	colWidths := make([]int, numCols)
	remainingWidth := terminalWidth - numCols // subtract the width of separators end

	//  initialize all column widths
	for i := range colWidths {
		colWidths[i] = remainingWidth / numCols
	}
	// dynamically adjust column widths
	for i := range colWidths {
		maxWidth := maxWidths[i]
		if maxWidth < colWidths[i] {
			remainingWidth += colWidths[i] - maxWidth
			colWidths[i] = maxWidth
		}
	}

	//  allocate remaining width
	for i := range colWidths {
		if remainingWidth <= 0 {
			break
		}
		colWidths[i] += remainingWidth / numCols
		remainingWidth -= remainingWidth / numCols
	}

	//  set the maximum width of the column
	columnConfigs := make([]table.ColumnConfig, numCols)
	for i := range columnConfigs {
		columnConfigs[i] = table.ColumnConfig{
			WidthMax:         colWidths[i],
			WidthMaxEnforcer: text.WrapSoft,
		}
	}
	return columnConfigs
}

package model

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	commonmodels "github.com/lindb/common/models"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/mapstructure"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/terminal"
	"github.com/lindb/lindb/spi/types"
)

type Schema struct {
	Columns   []types.ColumnMetadata `json:"columns,omitempty"`
	Partition []models.Partition     `json:"partitions,omitempty"`
}

type ResultSet struct {
	Schema *Schema `json:"schema,omitempty"`
	Rows   [][]any `json:"rows,omitempty"`
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		Schema: &Schema{},
	}
}

// ToTable returns stateless node list as table if it has value, else return empty string.
func (rs *ResultSet) ToTable() (tableStr string) {
	writer := commonmodels.NewTableFormatter()
	writer.SetStyle(terminal.TableSylte())
	var headers table.Row
	var columnTypes []types.DataType
	var (
		hasTimeSeries bool
		timeSeriesIdx int
		dataPoints    int
		rows          []table.Row
	)
	var maxWidths []int
	for i, col := range rs.Schema.Columns {
		if !hasTimeSeries && col.DataType == types.DTTimeSeries {
			timeSeriesIdx = i
			hasTimeSeries = true
			timeSeries := &types.TimeSeries{}
			_ = mapstructure.Decode(rs.Rows[0][i], timeSeries)
			dataPoints = len(timeSeries.Values)
			headers = append(headers, "timestamp") // add timestamp column
			columnTypes = append(columnTypes, types.DTTimestamp)
			maxWidths = append(maxWidths, len("timestamp"))
		}
		headers = append(headers, col.Name)
		columnTypes = append(columnTypes, col.DataType)
		maxWidths = append(maxWidths, len(col.Name))
	}
	writer.AppendHeader(headers)

	for _, row := range rs.Rows {
		if hasTimeSeries {
			// has time series, build row based on data points
			for pos := 0; pos < dataPoints; pos++ {
				cols := make(table.Row, len(rs.Schema.Columns)+1) // add timestamp column
				colIdx := 0
				for i, col := range row {
					if rs.Schema.Columns[colIdx].DataType == types.DTTimeSeries {
						timeSeries := &types.TimeSeries{}
						_ = mapstructure.Decode(row[colIdx], timeSeries)
						if timeSeriesIdx == i {
							cols[colIdx] = timeSeries.TimeRange.Start + timeSeries.Interval*int64(pos)
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
			cols := make(table.Row, len(rs.Schema.Columns))
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
func appendColumn(row table.Row, colType types.DataType, col any, index int) {
	switch colType {
	case types.DTString:
		row[index] = col.(string)
	case types.DTDuration:
		row[index] = time.Duration(col.(float64))
	case types.DTFloat, types.DTInt:
		row[index] = fmt.Sprintf("%v", col)
	case types.DTTimestamp:
		switch val := col.(type) {
		case string:
			row[index] = val
		case int64:
			row[index] = timeutil.FormatTimestamp(val, timeutil.DataTimeFormat2)
		}
	}
}

func stringWidth(width int, v any) int {
	return max(width, runewidth.StringWidth(fmt.Sprintf("%v", v)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		columnConfigs[i] = table.ColumnConfig{Number: i + 1, WidthMax: colWidths[i], WidthMaxEnforcer: text.WrapText}
	}
	return columnConfigs
}

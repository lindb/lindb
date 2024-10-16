package model

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	commonmodels "github.com/lindb/common/models"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/mitchellh/mapstructure"

	"github.com/lindb/lindb/models"
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
	var headers table.Row
	var columnTypes []types.DataType
	var (
		hasTimeSeries bool
		timeSeriesIdx int
		dataPoints    int
	)
	for i, col := range rs.Schema.Columns {
		if !hasTimeSeries && col.DataType == types.DTTimeSeries {
			timeSeriesIdx = i
			hasTimeSeries = true
			timeSeries := &types.TimeSeries{}
			_ = mapstructure.Decode(rs.Rows[0][i], timeSeries)
			dataPoints = len(timeSeries.Values)
			headers = append(headers, "timestamp") // add timestamp column
			columnTypes = append(columnTypes, types.DTTimestamp)
		}
		headers = append(headers, col.Name)
		columnTypes = append(columnTypes, col.DataType)
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
							colIdx++
							cols[colIdx] = timeSeries.Values[pos]
						} else {
							cols[colIdx] = timeSeries.Values[pos]
						}
					} else {
						appendColumn(cols, columnTypes[colIdx], col, colIdx)
					}
					colIdx++
				}
				writer.AppendRow(cols)
			}
		} else {
			cols := make(table.Row, len(rs.Schema.Columns))
			for i, col := range row {
				appendColumn(cols, columnTypes[i], col, i)
			}
			writer.AppendRow(cols)
		}
	}
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

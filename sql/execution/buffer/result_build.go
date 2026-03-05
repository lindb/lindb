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

package buffer

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/execution/model"
)

type ResultSetBuild struct {
	inbound   chan arrow.RecordBatch
	completed chan struct{}
	resultSet *model.ResultSet
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		inbound:   make(chan arrow.RecordBatch),
		completed: make(chan struct{}),
		resultSet: model.NewResultSet(),
	}
}

func (rsb *ResultSetBuild) AddRecord(record arrow.RecordBatch) {
	if record != nil && record.NumRows() > 0 {
		rsb.inbound <- record
	}
}

func (rsb *ResultSetBuild) Process() {
	defer func() {
		close(rsb.completed)
	}()
	// TODO: need close when timeout
	isTimestampSelected := false
	hasTimeSeries := false
	for record := range rsb.inbound {
		// TODO: how to handle error, maybe add error channel
		// if page.Error != "" {
		// 	rsb.resultSet.Error = page.Error
		// 	break
		// }
		if len(rsb.resultSet.Schema.Columns) == 0 {
			fields := record.Schema().Fields()
			lo.ForEach(fields, func(item arrow.Field, index int) {
				if arrow.TypeEqual(item.Type, arrow.FixedWidthTypes.Timestamp_ms) {
					isTimestampSelected = true
				}
				if arrow.TypeEqual(item.Type, larrow.ExtensionTypes.TimeSeries) {
					hasTimeSeries = true
				}
			})
			fmt.Printf("isTimestampSelected:%v, hasTimeSeries:%v\n", isTimestampSelected, hasTimeSeries)

			// lo.ForEach(record.Layout, func(item types.ColumnMetadata, index int) {
			// 	column := types.ColumnMetadata{
			// 		Name:     item.Name,
			// 		DataType: item.DataType,
			// 		Ref:      index,
			// 	}
			// 	if !hasTimeSeries {
			// 		rsb.resultSet.Schema.Columns = append(rsb.resultSet.Schema.Columns, column)
			// 		return
			// 	}
			//
			// 	if item.DataType == types.DTTimeSeries && !isTimestampSelected {
			// 		column.DataType = types.DTFloat
			// 	}
			// 	if item.DataType != types.DTTimestamp {
			// 		// ignore timestamp if select item list has time series
			// 		rsb.resultSet.Schema.Columns = append(rsb.resultSet.Schema.Columns, column)
			// 	}
			// })
		}
		// it := record.Iterator()
		// for row := it.Begin(); row != it.End(); row = it.Next() {
		// 	columns := make([]any, len(rsb.resultSet.Schema.Columns))
		// 	for i, c := range rsb.resultSet.Schema.Columns {
		// 		meta := record.Layout[c.Ref]
		// 		// TODO: add more type
		// 		switch meta.DataType {
		// 		case types.DTString, types.DTDynamic:
		// 			columns[i] = row.GetString(i)
		// 		case types.DTJSON:
		// 			columns[i] = row.GetJSON(i)
		// 		case types.DTInt:
		// 			columns[i] = row.GetInt(i)
		// 		case types.DTFloat:
		// 			columns[i] = row.GetFloat(i)
		// 		case types.DTTimeSeries:
		// 			timeSeries := row.GetTimeSeries(i)
		// 			if timeSeries == nil {
		// 				columns[i] = nil
		// 				continue
		// 			}
		// 			if isTimestampSelected || timeSeries.NumOfPoints > 1 {
		// 				columns[i] = timeSeries
		// 			} else {
		// 				columns[i] = timeSeries.GetValue()
		// 			}
		// 		case types.DTTimestamp:
		// 			// FIXME: maybe timestamp is nil
		// 			columns[i] = row.GetTimestamp(i).UnixMilli()
		// 		case types.DTDuration:
		// 			columns[i] = row.GetDuration(i)
		// 		case types.DTExemplar, types.DTMap:
		// 			// FIXME: exemplar not support now
		// 			columns[i] = row.Get(i)
		// 		default:
		// 			panic(fmt.Sprintf("build resultset error, column:%v, unknown data type:%v", meta.Name, meta.DataType))
		// 		}
		// 	}
		// 	rsb.resultSet.Rows = append(rsb.resultSet.Rows, columns)
		// }
	}
}

func (rsb *ResultSetBuild) Complete() {
	close(rsb.inbound)
}

func (rsb *ResultSetBuild) ResultSet() *model.ResultSet {
	// waiting process result page completed
	<-rsb.completed
	return rsb.resultSet
}

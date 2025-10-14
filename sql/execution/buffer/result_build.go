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

	"github.com/lindb/common/pkg/encoding"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/model"
)

type ResultSetBuild struct {
	inbound   chan *types.Page
	completed chan struct{}
	resultSet *model.ResultSet
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		inbound:   make(chan *types.Page),
		completed: make(chan struct{}),
		resultSet: model.NewResultSet(),
	}
}

func (rsb *ResultSetBuild) AddPage(page *types.Page) {
	if page != nil {
		rsb.inbound <- page
	}
}

func (rsb *ResultSetBuild) Process() {
	defer func() {
		close(rsb.completed)
	}()
	// TODO: need close when timeout
	isTimestampSelected := false
	hasTimeSeries := false
	for page := range rsb.inbound {
		fmt.Printf("page.....=.....%v\n", string(encoding.JSONMarshal(page))) // TODO: remove page)
		if page.Error != "" {
			rsb.resultSet.Error = page.Error
			break
		}
		if len(rsb.resultSet.Schema.Columns) == 0 {
			lo.ForEach(page.Layout, func(item types.ColumnMetadata, index int) {
				if item.DataType == types.DTTimestamp {
					isTimestampSelected = true
				}
				if item.DataType == types.DTTimeSeries {
					hasTimeSeries = true
				}
			})

			lo.ForEach(page.Layout, func(item types.ColumnMetadata, index int) {
				column := types.ColumnMetadata{
					Name:     item.Name,
					DataType: item.DataType,
					Ref:      index,
				}
				if !hasTimeSeries {
					rsb.resultSet.Schema.Columns = append(rsb.resultSet.Schema.Columns, column)
					return
				}

				if item.DataType == types.DTTimeSeries && !isTimestampSelected {
					column.DataType = types.DTFloat
				}
				if item.DataType != types.DTTimestamp {
					// ignore timestamp if select item list has time series
					rsb.resultSet.Schema.Columns = append(rsb.resultSet.Schema.Columns, column)
				}
			})
		}
		it := page.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			columns := make([]any, len(rsb.resultSet.Schema.Columns))
			for i, c := range rsb.resultSet.Schema.Columns {
				meta := page.Layout[c.Ref]
				// TODO: add more type
				switch meta.DataType {
				case types.DTString:
					columns[i] = row.GetString(i)
				case types.DTJSON:
					columns[i] = row.GetJSON(i)
				case types.DTInt:
					columns[i] = row.GetInt(i)
				case types.DTFloat:
					columns[i] = row.GetFloat(i)
				case types.DTTimeSeries:
					timeSeries := row.GetTimeSeries(i)
					if isTimestampSelected || timeSeries.NumOfPoints > 1 {
						columns[i] = timeSeries
					} else {
						columns[i] = timeSeries.GetValue()
					}
				case types.DTTimestamp:
					// FIXME: maybe timestamp is nil
					columns[i] = row.GetTimestamp(i).UnixMilli()
				case types.DTDuration:
					columns[i] = row.GetDuration(i)
				default:
					panic(fmt.Sprintf("build resultset error, column:%v, unknown data type:%v", meta.Name, meta.DataType))
				}
			}
			rsb.resultSet.Rows = append(rsb.resultSet.Rows, columns)
		}
		fmt.Println("merge result page")
	}
}

func (rsb *ResultSetBuild) Complete() {
	fmt.Println("ResultSetBuild close result page")
	close(rsb.inbound)
}

func (rsb *ResultSetBuild) ResultSet() *model.ResultSet {
	// waiting process result page completed
	<-rsb.completed
	fmt.Println("result.....")
	return rsb.resultSet
}

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

package types

import (
	"github.com/apache/arrow-go/v18/arrow"
)

func MergeRecords(pages []arrow.RecordBatch) arrow.RecordBatch {
	switch len(pages) {
	case 0:
		return nil
	case 1:
		return pages[0]
	default:
		panic("not implemented merge")
		// mergedPage := NewPage()
		// allColumns := make([]*Column, len(pages[0].Columns))
		// grouping := make(map[string]int) // grouping column values => row number
		// var mergeColumns []*mergeColumn
		// var keys []any
		// groupingColumns := pages[0].Grouping
		// if len(groupingColumns) > 0 {
		// 	keys = make([]any, len(groupingColumns))
		// }
		// for idx, c := range pages[0].Layout {
		// 	column := NewColumn()
		// 	allColumns[idx] = column
		// 	mergedPage.AppendColumn(c, column)
		// 	if !slices.Contains(groupingColumns, idx) {
		// 		// columns need merge value
		// 		mergeColumns = append(mergeColumns, &mergeColumn{meta: c, target: column, index: idx})
		// 	}
		// }
		// for _, page := range pages {
		// 	if page.Error != "" {
		// 		// NOTE: if has error, return it
		// 		mergedPage.Error = page.Error
		// 		break
		// 	}
		// 	it := page.Iterator()
		// 	rowNum := 0
		// 	for row := it.Begin(); row != it.End(); row = it.Next() {
		// 		for idx, columnIndex := range groupingColumns {
		// 			keys[idx] = row.Get(columnIndex)
		// 		}
		// 		keysStr := strutil.SliceToTypedString(keys) // merge keys
		// 		mergedRow, ok := grouping[keysStr]
		// 		if !ok {
		// 			grouping[keysStr] = rowNum
		// 			for idx, column := range allColumns {
		// 				column.Append(row.Get(idx))
		// 			}
		// 		} else {
		// 			for _, column := range mergeColumns {
		// 				// merge values of columns
		// 				column.Merge(mergedRow, row)
		// 			}
		// 		}
		//
		// 		rowNum++
		// 	}
		// }
		// return mergedPage
	}
}

type mergeColumn struct {
	meta   ColumnMetadata
	target *Column
	index  int
}

func (mc *mergeColumn) Merge(mergedRow int, row Row) {
	if mc.meta.DataType == DTTimeSeries {
		// NOTE: merge time series
		oldVal := mc.target.GetTimeSeries(mergedRow)
		newVal := row.GetTimeSeries(mc.index)

		switch {
		case oldVal == nil && newVal == nil:
		case oldVal == nil:
			mc.target.Reset(mergedRow, newVal)
		default:
			if len(oldVal.Values) != len(newVal.Values) {
				panic("merge time series values length not equal")
			} else {
				for i, v := range oldVal.Values {
					oldVal.Values[i] = mc.meta.AggType.Aggregate(v, newVal.Values[i])
				}
			}
		}
	}
	// TODO: merge value for other data type
}

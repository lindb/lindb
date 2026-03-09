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
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	commonmodels "github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
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
	for record := range rsb.inbound {
		fmt.Println(record)
		rsb.processRecord(record)
		record.Release()
	}
}

func (rsb *ResultSetBuild) processRecord(record arrow.RecordBatch) {
	fields := record.Schema().Fields()

	// build schema on first record
	if rsb.resultSet.Schema == nil {
		var visibleFields []arrow.Field
		for _, f := range fields {
			if f.Metadata.FindKey("hidden") >= 0 {
				continue
			}
			visibleFields = append(visibleFields, f)
		}
		rsb.resultSet.Schema = arrow.NewSchema(visibleFields, nil)
	}

	numRows := int(record.NumRows())
	for row := range numRows {
		cols := make([]any, 0, rsb.resultSet.Schema.NumFields())
		for colIdx, f := range fields {
			if f.Metadata.FindKey("hidden") >= 0 {
				continue
			}
			col := record.Column(colIdx)
			if col.IsNull(row) {
				cols = append(cols, nil)
				continue
			}
			cols = append(cols, extractValue(f.Type, col, row))
		}
		rsb.resultSet.Rows = append(rsb.resultSet.Rows, cols)
	}
}

func extractValue(dt arrow.DataType, col arrow.Array, row int) any {
	fmt.Printf("extract value from column type: %s,%v\n", dt, col)
	switch c := col.(type) {
	case *array.String:
		return c.Value(row)
	case *array.Int64:
		return c.Value(row)
	case *array.Int32:
		return c.Value(row)
	case *array.Float64:
		return c.Value(row)
	case *array.Timestamp:
		return int64(c.Value(row))
	case *larray.TimeSeries:
		structArr := c.Storage().(*array.Struct)
		start := structArr.Field(0).(*array.Int64).Value(row)
		end := structArr.Field(1).(*array.Int64).Value(row)
		interval := structArr.Field(2).(*array.Int64).Value(row)
		listArr := structArr.Field(3).(*array.List)
		offsets := listArr.Offsets()
		from, to := int(offsets[row]), int(offsets[row+1])
		floats := listArr.ListValues().(*array.Float64)
		values := make([]float64, to-from)
		for i := range values {
			values[i] = floats.Value(from + i)
		}
		return types.NewTimeSeriesWithValues(
			timeutil.TimeRange{Start: start, End: end},
			timeutil.Interval(interval),
			values,
		)
	case *larray.Exemplar:
		ex := c.Value(row)
		if ex == nil {
			return nil
		}
		return &commonmodels.Exemplar{
			TraceID:  string(ex.TraceID),
			SpanID:   string(ex.SpanID),
			Duration: ex.Duration,
		}
	case *larray.Aggregation:
		return c.Value(row)
	default:
		_ = dt
		return nil
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

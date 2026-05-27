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

package expression

import (
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/spi/scalar"
)

type addSubDateFunc struct {
	ctx  EvalContext
	args []Expression
}

func newAddSubDateFunc(ctx EvalContext, args []Expression) Func {
	return &addSubDateFunc{
		ctx:  ctx,
		args: args,
	}
}

func (n *addSubDateFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	// Scalar date argument (e.g. DATE_ADD('2024-01-01', INTERVAL 1 DAY)):
	// evaluate once and broadcast the result to every row.
	if n.args[0].ResultType() == Scalar {
		result, err := n.EvalScalar()
		if err != nil {
			return nil, err
		}
		return broadcastTimeScalar(scalar.ToTime(result), record.NumRows()), nil
	}

	// Array date argument: apply the duration to each timestamp in the column.
	input, err := n.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	durationScalar, err := n.args[1].EvalScalar()
	if err != nil {
		return nil, err
	}
	duration := scalar.ToDuration(durationScalar)

	timestampArray, ok := input.(*larray.Generic[arrow.Timestamp])
	if !ok {
		return nil, fmt.Errorf("date_add: unexpected input type %T, expected Timestamp array", input)
	}
	dataType := timestampArray.DataType().(*arrow.TimestampType)
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, dataType)
	builder.Reserve(timestampArray.Len())
	for i := range timestampArray.Len() {
		if timestampArray.IsNull(i) {
			builder.AppendNull()
		} else {
			ts := timestampArray.Value(i).ToTime(dataType.Unit).Add(duration)
			builder.UnsafeAppend(arrow.Timestamp(ts.UnixMilli()))
		}
	}
	return builder.NewArray(), nil
}

func (n *addSubDateFunc) EvalScalar() (scalar.Scalar, error) {
	tsScalar, err := n.args[0].EvalScalar()
	if err != nil {
		return nil, err
	}
	tsStr := scalar.ToString(tsScalar)
	durationScalar, err := n.args[1].EvalScalar()
	if err != nil {
		return nil, err
	}
	duration := scalar.ToDuration(durationScalar)

	format := timeutil.DataTimeFormat2
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return nil, err
	}
	return scalar.NewTimeScalar(time.UnixMilli(timestamp).Add(duration)), nil
}

type nowFunc struct {
	ctx  EvalContext
	args []Expression

	builder *array.TimestampBuilder
}

func newNowFunc(ctx EvalContext, args []Expression) Func {
	return &nowFunc{
		ctx:  ctx,
		args: args,
	}
}

func (n *nowFunc) EvalScalar() (scalar.Scalar, error) {
	return scalar.NewTimeScalar(n.ctx.CurrentTime()), nil
}

func (n *nowFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	numRows := int(record.NumRows())
	if n.builder == nil {
		// NOW() return type is Timestamp_ns per defaultFuncReturnTypes.
		n.builder = array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_ns.(*arrow.TimestampType))
	}
	builder := n.builder
	// NOW() is a scalar — every row in the batch gets the same timestamp.
	ts := arrow.Timestamp(n.ctx.CurrentTime().UnixNano())
	builder.Reserve(numRows)
	for range numRows {
		builder.UnsafeAppend(ts)
	}
	return builder.NewArray(), nil
}

type strToDateFunc struct {
	ctx  EvalContext
	args []Expression
}

func newStrToDateFunc(ctx EvalContext, args []Expression) Func {
	return &strToDateFunc{
		ctx:  ctx,
		args: args,
	}
}

func (n *strToDateFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	// Scalar string argument: parse once and broadcast to all rows.
	if n.args[0].ResultType() == Scalar {
		result, err := n.EvalScalar()
		if err != nil {
			return nil, err
		}
		return broadcastTimeScalar(scalar.ToTime(result), record.NumRows()), nil
	}

	// Array string argument: parse each string value in the column.
	input, err := n.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	formatScalar, err := n.args[1].EvalScalar()
	if err != nil {
		return nil, err
	}
	format := scalar.ToString(formatScalar)
	switch format {
	case "YYYYMMDD HH:mm:ss":
		format = timeutil.DataTimeFormat1
	case "YYYY-MM-DD HH:mm:ss":
		format = timeutil.DataTimeFormat2
	case "YYYY/MM/DD HH:mm:ss":
		format = timeutil.DataTimeFormat3
	case "YYYYMMDDHHmmss":
		format = timeutil.DataTimeFormat4
	}

	strArray, ok := input.(*array.String)
	if !ok {
		return nil, fmt.Errorf("str_to_date: unexpected input type %T, expected String array", input)
	}
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_ns.(*arrow.TimestampType))
	builder.Reserve(strArray.Len())
	for i := range strArray.Len() {
		if strArray.IsNull(i) {
			builder.AppendNull()
		} else {
			ts, err := timeutil.ParseTimestamp(strArray.Value(i), format)
			if err != nil {
				return nil, fmt.Errorf("str_to_date: failed to parse %q: %w", strArray.Value(i), err)
			}
			builder.UnsafeAppend(arrow.Timestamp(time.UnixMilli(ts).UnixNano()))
		}
	}
	return builder.NewArray(), nil
}

func (n *strToDateFunc) EvalScalar() (scalar.Scalar, error) {
	tsScalar, err := n.args[0].EvalScalar()
	if err != nil {
		return nil, err
	}
	tsStr := scalar.ToString(tsScalar)
	formatScalar, err := n.args[1].EvalScalar()
	if err != nil {
		return nil, err
	}
	format := scalar.ToString(formatScalar)
	switch format {
	case "YYYYMMDD HH:mm:ss":
		format = timeutil.DataTimeFormat1
	case "YYYY-MM-DD HH:mm:ss":
		format = timeutil.DataTimeFormat2
	case "YYYY/MM/DD HH:mm:ss":
		format = timeutil.DataTimeFormat3
	case "YYYYMMDDHHmmss":
		format = timeutil.DataTimeFormat4
	}
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return nil, err
	}
	return scalar.NewTimeScalar(time.UnixMilli(timestamp)), nil
}

type timeTruncFunc struct {
	ctx  EvalContext
	args []Expression

	duration time.Duration

	builder *array.TimestampBuilder
}

func newTimeTruncFunc(ctx EvalContext, args []Expression) Func {
	duration, err := args[1].EvalScalar()
	if err != nil {
		panic(fmt.Sprintf("failed to evaluate duration argument in time_trunc function: %v", err))
	}
	return &timeTruncFunc{
		ctx:      ctx,
		args:     args,
		duration: scalar.ToDuration(duration),
	}
}

func (n *timeTruncFunc) EvalScalar() (scalar.Scalar, error) {
	panic("time_trunc does not support scalar evaluation")
}

func (n *timeTruncFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	input, err := n.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release() // release the input array after processing

	timestampArray, ok := input.(*larray.Generic[arrow.Timestamp])
	if !ok {
		return nil, fmt.Errorf("unexpected input type for time_trunc: %T", input)
	}

	dataType := timestampArray.DataType().(*arrow.TimestampType)
	if n.builder == nil {
		n.builder = array.NewTimestampBuilder(memory.DefaultAllocator, dataType)
	}
	builder := n.builder

	builder.Reserve(timestampArray.Len())
	for i := 0; i < timestampArray.Len(); i++ {
		if timestampArray.IsNull(i) {
			builder.AppendNull()
		} else {
			timestamp := timestampArray.Value(i).ToTime(dataType.Unit)
			truncated := timestamp.Truncate(n.duration)
			// return milliseconds
			builder.Append(arrow.Timestamp(truncated.UnixMilli()))
		}
	}

	return builder.NewArray(), nil
}

// broadcastTimeScalar creates a Timestamp_ns array where every row has the same value t.
// Used by scalar time functions (DATE_ADD, STR_TO_DATE) when their date argument is a constant.
func broadcastTimeScalar(t time.Time, numRows int64) arrow.Array {
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_ns.(*arrow.TimestampType))
	ts := arrow.Timestamp(t.UnixNano())
	builder.Reserve(int(numRows))
	for range numRows {
		builder.UnsafeAppend(ts)
	}
	return builder.NewArray()
}

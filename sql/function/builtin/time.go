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

package builtin

import (
	"errors"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// ── DATE_ADD ──────────────────────────────────────────────────────────────────

// addSubDateInstance holds a per-query DATE_ADD/DATE_SUB state.
// If the date argument is a constant, it is cached as dateScalar at construction time.
type addSubDateInstance struct {
	args       []function.Expr
	dateScalar scalar.Scalar  // non-nil when first arg is a constant date string
	duration   time.Duration  // cached when second arg is a constant interval
	hasDur     bool
}

var AddSubDateFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	f := &addSubDateInstance{args: args}
	f.dateScalar, _ = args[0].EvalScalar()
	if ds, err := args[1].EvalScalar(); err == nil {
		f.duration = scalar.ToDuration(ds)
		f.hasDur = true
	}
	return f
}

func (f *addSubDateInstance) EvalScalar() (scalar.Scalar, error) {
	if f.dateScalar == nil || !f.hasDur {
		return nil, errors.New("date_add: not a constant expression")
	}
	timestamp, err := timeutil.ParseTimestamp(scalar.ToString(f.dateScalar), timeutil.DataTimeFormat2)
	if err != nil {
		return nil, err
	}
	return scalar.NewTimeScalar(time.UnixMilli(timestamp).Add(f.duration)), nil
}

func (f *addSubDateInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	if f.dateScalar != nil {
		result, err := f.EvalScalar()
		if err != nil {
			return nil, err
		}
		return broadcastTimeScalar(scalar.ToTime(result), record.NumRows()), nil
	}

	input, err := f.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	var duration time.Duration
	if f.hasDur {
		duration = f.duration
	} else {
		ds, err := f.args[1].EvalScalar()
		if err != nil {
			return nil, err
		}
		duration = scalar.ToDuration(ds)
	}

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

// ── NOW ───────────────────────────────────────────────────────────────────────

type nowInstance struct {
	ctx function.EvalContext
}

var NowFactory function.VectorFuncFactory = func(ctx function.EvalContext, _ []function.Expr) function.VectorFunc {
	return &nowInstance{ctx: ctx}
}

func (f *nowInstance) EvalScalar() (scalar.Scalar, error) {
	return scalar.NewTimeScalar(f.ctx.CurrentTime()), nil
}

func (f *nowInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	numRows := int(record.NumRows())
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_ns.(*arrow.TimestampType))
	ts := arrow.Timestamp(f.ctx.CurrentTime().UnixNano())
	builder.Reserve(numRows)
	for range numRows {
		builder.UnsafeAppend(ts)
	}
	return builder.NewArray(), nil
}

// ── STR_TO_DATE ───────────────────────────────────────────────────────────────

type strToDateInstance struct {
	args       []function.Expr
	dateScalar scalar.Scalar // non-nil when first arg is a constant
	format     string        // cached when second arg is a constant format string
}

var StrToDateFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	f := &strToDateInstance{args: args}
	f.dateScalar, _ = args[0].EvalScalar()
	if fs, err := args[1].EvalScalar(); err == nil {
		f.format = normaliseTimeFormat(scalar.ToString(fs))
	}
	return f
}

func (f *strToDateInstance) EvalScalar() (scalar.Scalar, error) {
	if f.dateScalar == nil || f.format == "" {
		return nil, errors.New("str_to_date: not a constant expression")
	}
	timestamp, err := timeutil.ParseTimestamp(scalar.ToString(f.dateScalar), f.format)
	if err != nil {
		return nil, err
	}
	return scalar.NewTimeScalar(time.UnixMilli(timestamp)), nil
}

func (f *strToDateInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	if f.dateScalar != nil {
		result, err := f.EvalScalar()
		if err != nil {
			return nil, err
		}
		return broadcastTimeScalar(scalar.ToTime(result), record.NumRows()), nil
	}

	input, err := f.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	format := f.format
	if format == "" {
		if fs, err := f.args[1].EvalScalar(); err == nil {
			format = normaliseTimeFormat(scalar.ToString(fs))
		}
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

// ── TIME_TRUNC ────────────────────────────────────────────────────────────────

type timeTruncInstance struct {
	args     []function.Expr
	duration time.Duration // cached when second arg is a constant
	hasDur   bool
}

var TimeTruncFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	f := &timeTruncInstance{args: args}
	if ds, err := args[1].EvalScalar(); err == nil {
		f.duration = scalar.ToDuration(ds)
		f.hasDur = true
	}
	return f
}

func (f *timeTruncInstance) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("time_trunc does not support scalar evaluation")
}

func (f *timeTruncInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	var duration time.Duration
	if f.hasDur {
		duration = f.duration
	} else {
		ds, err := f.args[1].EvalScalar()
		if err != nil {
			return nil, fmt.Errorf("time_trunc: failed to evaluate duration argument: %w", err)
		}
		duration = scalar.ToDuration(ds)
	}

	input, err := f.args[0].Eval(record)
	if err != nil {
		return nil, err
	}
	defer input.Release()

	timestampArray, ok := input.(*larray.Generic[arrow.Timestamp])
	if !ok {
		return nil, fmt.Errorf("time_trunc: unexpected input type %T", input)
	}
	dataType := timestampArray.DataType().(*arrow.TimestampType)
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, dataType)
	builder.Reserve(timestampArray.Len())
	for i := range timestampArray.Len() {
		if timestampArray.IsNull(i) {
			builder.AppendNull()
		} else {
			truncated := timestampArray.Value(i).ToTime(dataType.Unit).Truncate(duration)
			builder.Append(arrow.Timestamp(truncated.UnixMilli()))
		}
	}
	return builder.NewArray(), nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

// normaliseTimeFormat converts user-friendly format strings to timeutil constants.
func normaliseTimeFormat(format string) string {
	switch format {
	case "YYYYMMDD HH:mm:ss":
		return timeutil.DataTimeFormat1
	case "YYYY-MM-DD HH:mm:ss":
		return timeutil.DataTimeFormat2
	case "YYYY/MM/DD HH:mm:ss":
		return timeutil.DataTimeFormat3
	case "YYYYMMDDHHmmss":
		return timeutil.DataTimeFormat4
	}
	return format
}

// broadcastTimeScalar creates a Timestamp_ns array where every row has the same value t.
func broadcastTimeScalar(t time.Time, numRows int64) arrow.Array {
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, arrow.FixedWidthTypes.Timestamp_ns.(*arrow.TimestampType))
	ts := arrow.Timestamp(t.UnixNano())
	builder.Reserve(int(numRows))
	for range numRows {
		builder.UnsafeAppend(ts)
	}
	return builder.NewArray()
}

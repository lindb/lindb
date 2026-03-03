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
	panic("add_sub_date is not supported in vectorized execution")
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
	panic("now is not supported in vectorized execution")
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
	panic("str_to_date is not supported in vectorized execution")
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

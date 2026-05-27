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

	"github.com/lindb/lindb/spi/scalar"
)

type Constant struct {
	value any
	rt    ResultType
}

func NewConstant(ctx EvalContext, value any, rt ResultType) Expression {
	return &Constant{
		value: value,
		rt:    rt,
	}
}

func (c *Constant) EvalScalar() (scalar.Scalar, error) {
	switch val := c.value.(type) {
	case string:
		return scalar.NewStringScalar(val), nil
	case int64:
		return scalar.NewInt64Scalar(val), nil
	case float64:
		return scalar.NewFloat64Scalar(val), nil
	case time.Duration:
		return scalar.NewDurationScalar(val), nil
	default:
		panic(fmt.Sprintf("unsupported data type for constant: %T", val))
	}
}

func (c *Constant) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	numRows := int(record.NumRows())
	// Broadcast the constant value to every row in the batch.
	switch val := c.value.(type) {
	case int64:
		b := array.NewInt64Builder(memory.DefaultAllocator)
		defer b.Release()
		b.Reserve(numRows)
		for range numRows {
			b.UnsafeAppend(val)
		}
		return b.NewArray(), nil
	case float64:
		b := array.NewFloat64Builder(memory.DefaultAllocator)
		defer b.Release()
		b.Reserve(numRows)
		for range numRows {
			b.UnsafeAppend(val)
		}
		return b.NewArray(), nil
	case string:
		b := array.NewStringBuilder(memory.DefaultAllocator)
		defer b.Release()
		b.Reserve(numRows)
		for range numRows {
			b.Append(val)
		}
		return b.NewArray(), nil
	default:
		return nil, fmt.Errorf("constant: unsupported type for vectorized evaluation: %T", c.value)
	}
}

// func (c *Constant) EvalExemplar(_ types.Row) (val *models.Exemplar, isNull bool, err error) {
// 	return
// }
//
// func (c *Constant) EvalTime(_ types.Row) (val time.Time, isNull bool, err error) {
// 	switch v := c.value.(type) {
// 	case string:
// 		timestamp, err := timeutil.ParseTimestamp(v, timeutil.DataTimeFormat2)
// 		if err != nil {
// 			return time.Time{}, true, err
// 		}
// 		return time.UnixMilli(timestamp), false, nil
// 	default:
// 		return
// 	}
// }

func (c *Constant) ResultType() ResultType {
	return c.rt
}

// String returns the constant in string format.
func (c *Constant) String() string {
	return fmt.Sprintf("%v", c.value)
}

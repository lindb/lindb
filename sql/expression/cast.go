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

	"github.com/lindb/lindb/spi/types"
)

type Cast struct {
	function Func
	arg      Expression
	retType  types.DataType
}

func NewCast(retType types.DataType, arg Expression) Expression {
	return &Cast{
		retType: retType,
		arg:     arg,
		function: &castFunc{
			baseFunc: baseFunc{args: []Expression{arg}},
		},
	}
}

// EvalString implements Expression.
func (c *Cast) EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error) {
	panic("unimplemented")
}

func (c *Cast) EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error) {
	fmt.Printf("cast eval int=%v\n", c.retType)
	return c.function.EvalInt(ctx, row)
}

func (c *Cast) EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error) {
	return c.function.EvalFloat(ctx, row)
}

func (c *Cast) EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return c.function.EvalTimeSeries(ctx, row)
}

func (c *Cast) EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error) {
	return
}

func (c *Cast) EvalTime(ctx EvalContext, _ types.Row) (val time.Time, isNull bool, err error) {
	return
}

// GetType implements Expression.
func (c *Cast) GetType() types.DataType {
	return c.retType
}

func (c *Cast) String() string {
	return fmt.Sprintf("CAST(%s as %s)", c.arg.String(), c.retType)
}

type castFunc struct {
	baseFunc
}

func (f *castFunc) EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(ctx, row)
	fmt.Println("cast int..........")
	return lv, false, nil
}

// EvalFloat implements Func.
func (f *castFunc) EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("cast float..........")
	return
}

// EvalTimeSeries evaluates the expression, cast result to types.TimeSeries type.
func (f *castFunc) EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Printf("cast time series..........,type =%T,%s,%s\n", f.args[0], f.args[0].GetType(), f.args[0].String())
	switch f.args[0].GetType() {
	case types.DTInt:
		val, isNull, err := f.args[0].EvalInt(ctx, row)
		if err != nil {
			return nil, false, err
		}
		if isNull {
			return nil, true, nil
		}
		return types.NewTimeSeriesWithSingleValue(float64(val)), false, nil
	case types.DTFloat:
		val, isNull, err := f.args[0].EvalFloat(ctx, row)
		if err != nil {
			return nil, false, err
		}
		if isNull {
			return nil, true, nil
		}
		return types.NewTimeSeriesWithSingleValue(val), false, nil
	case types.DTTimeSeries:
		return f.args[0].EvalTimeSeries(ctx, row)
	}
	return
}

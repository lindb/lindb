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
	"errors"
	"fmt"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
)

type arithmeticPlusFunc struct {
	baseFunc
}

func newArithmeticPlusFunc(ctx EvalContext, args []Expression) Func {
	return &arithmeticPlusFunc{
		baseFunc: baseFunc{ctx: ctx, args: args},
	}
}

func (f *arithmeticPlusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("plus int.....")
	return lv + rv, false, nil
}

func (f *arithmeticPlusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("plus float.....")
	return
}

func (f *arithmeticPlusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("plus time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv + rv
	})
}

func (f *arithmeticPlusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalTime(row)
	rv, _, _ := f.args[1].EvalDuration(row)
	val = lv.Add(rv)
	return
}

func evalTimeSeries(row types.Row, args []Expression,
	math func(lv, rv float64) float64,
) (val *types.TimeSeries, isNull bool, err error) {
	l, lIsNull, err := args[0].EvalTimeSeries(row)
	if err != nil {
		return nil, false, err
	}
	r, rIsNull, err := args[1].EvalTimeSeries(row)
	if err != nil {
		return nil, false, err
	}
	if lIsNull {
		return r, rIsNull, nil
	}
	if rIsNull {
		return l, lIsNull, nil
	}
	// check num. of points whether match
	if !l.IsSingleValue() && !r.IsSingleValue() && l.Size() != r.Size() {
		return nil, true, errors.New("num. of points not match")
	}
	var result *types.TimeSeries
	if !l.IsSingleValue() {
		result = types.NewTimeSeries(l.TimeRange, timeutil.Interval(l.Interval))
	} else if !r.IsSingleValue() {
		result = types.NewTimeSeries(r.TimeRange, timeutil.Interval(r.Interval))
	} else {
		result = types.NewTimeSeriesWithSingleValue(0)
	}
	for i := range result.Size() {
		result.Put(i, math(l.Get(i), r.Get(i)))
	}
	return result, false, nil
}

type arithmeticMinusFunc struct {
	baseFunc
}

func newArithmeticMinusFunc(ctx EvalContext, args []Expression) Func {
	return &arithmeticMinusFunc{
		baseFunc: baseFunc{ctx: ctx, args: args},
	}
}

func (f *arithmeticMinusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("minus int.....")
	return lv - rv, false, nil
}

func (f *arithmeticMinusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	fmt.Println("minus float.....")
	return
}

func (f *arithmeticMinusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("minus time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv - rv
	})
}

func (f *arithmeticMinusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalTime(row)
	rv, _, _ := f.args[1].EvalDuration(row)
	val = lv.Add(-rv)
	return
}

type arithmeticMulFunc struct {
	baseFunc
}

func newArithmeticMulFunc(ctx EvalContext, args []Expression) Func {
	return &arithmeticMulFunc{
		baseFunc: baseFunc{ctx: ctx, args: args},
	}
}

func (f *arithmeticMulFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("mul int.....")
	return lv * rv, false, nil
}

func (f *arithmeticMulFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("mul time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		return lv * rv
	})
}

type arithmeticDivFunc struct{ baseFunc }

func newArithmeticDivFunc(ctx EvalContext, args []Expression) Func {
	return &arithmeticDivFunc{
		baseFunc: baseFunc{ctx: ctx, args: args},
	}
}

func (f *arithmeticDivFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)
	fmt.Println("div int.....")
	return lv / rv, false, nil
}

func (f *arithmeticDivFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("div time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		if rv == 0 {
			return 0
		}
		return lv / rv
	})
}

type arithmeticModFunc struct {
	baseFunc
}

func newArithmeticModFunc(ctx EvalContext, args []Expression) Func {
	return &arithmeticModFunc{
		baseFunc: baseFunc{ctx: ctx, args: args},
	}
}

func (f *arithmeticModFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	rv, _, _ := f.args[1].EvalInt(row)

	fmt.Println("mod int.....")
	return lv % rv, false, nil
}

func (f *arithmeticModFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Println("mod time series.....")
	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
		if rv == 0 {
			return 0
		}
		return float64(int64(lv) % int64(rv))
	})
}

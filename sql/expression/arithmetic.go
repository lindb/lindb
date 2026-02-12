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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/lindb/lindb/spi/scalar"
)

type arithmeticFunc struct {
	ctx  EvalContext
	args []Expression

	op string
}

func newArithmeticFunc(ctx EvalContext, op string, args []Expression) Func {
	return &arithmeticFunc{
		ctx:  ctx,
		args: args,
		op:   op,
	}
}

func (f *arithmeticFunc) EvalScalar() (scalar.Scalar, error) {
	return nil, nil
}

func (f *arithmeticFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	left := f.args[0]
	right := f.args[1]

	switch {
	case left.ResultType() == Scalar:
		return evalScalarArray(record, left, right, f.op)
	case right.ResultType() == Scalar:
		return evalArrayScalar(record, left, right, f.op)
	case left.ResultType() == Array && right.ResultType() == Array:

	default:
		return nil, fmt.Errorf("unsupported result type: %v, %v", left.ResultType(), right.ResultType())
	}
	return nil, nil
}

func evalArrayScalar(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	leftArray, err := left.Eval(record)
	if err != nil {
		return nil, err
	}
	defer leftArray.Release()

	rightScalar, err := right.EvalScalar()
	if err != nil {
		return nil, err
	}
	switch left := leftArray.(type) {
	case *array.Int64:
		right := scalar.ToInt64(rightScalar)

		result := array.NewInt64Builder(memory.DefaultAllocator)
		defer result.Release()

		result.Reserve(left.Len())

		// do function logic
		for i := 0; i < left.Len(); i++ {
			if left.IsNull(i) {
				result.AppendNull()
			} else {
				result.Append(left.Value(i) + right)
			}
		}

		return result.NewArray(), nil
	default:
		panic(fmt.Sprintf("unsupported array type: %T", leftArray))
	}
}

func evalScalarArray(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	return nil, nil
}

func evalArrayArray(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	return nil, nil
}

// type arithmeticPlusFunc struct {
// 	ctx  EvalContext
// 	args []Expression
// }
//
// func newArithmeticPlusFunc(ctx EvalContext, args []Expression) Func {
// 	return &arithmeticPlusFunc{
// 		ctx:  ctx,
// 		args: args,
// 	}
// }
//
// func (f *arithmeticPlusFunc) EvalScalar() (scalar.Scalar, error) {
// 	left, err := f.args[0].EvalScalar()
// 	if err != nil {
// 		return nil, err
// 	}
// 	right, err := f.args[1].EvalScalar()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
//
// func (f *arithmeticPlusFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
// 	return nil, nil
// }
//
// func (f *arithmeticPlusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalInt(row)
// 	rv, _, _ := f.args[1].EvalInt(row)
// 	fmt.Println("plus int.....")
// 	return lv + rv, false, nil
// }
//
// func (f *arithmeticPlusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
// 	fmt.Println("plus float.....")
// 	return
// }
//
// func (f *arithmeticPlusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
// 	fmt.Println("plus time series.....")
// 	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
// 		return lv + rv
// 	})
// }
//
// func (f *arithmeticPlusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalTime(row)
// 	rv, _, _ := f.args[1].EvalDuration(row)
// 	val = lv.Add(rv)
// 	return
// }
//
// func evalTimeSeries(row types.Row, args []Expression,
// 	math func(lv, rv float64) float64,
// ) (val *types.TimeSeries, isNull bool, err error) {
// 	l, lIsNull, err := args[0].EvalTimeSeries(row)
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	r, rIsNull, err := args[1].EvalTimeSeries(row)
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	if lIsNull {
// 		return r, rIsNull, nil
// 	}
// 	if rIsNull {
// 		return l, lIsNull, nil
// 	}
// 	// check num. of points whether match
// 	if !l.IsSingleValue() && !r.IsSingleValue() && l.Size() != r.Size() {
// 		return nil, true, errors.New("num. of points not match")
// 	}
// 	var result *types.TimeSeries
// 	if !l.IsSingleValue() {
// 		result = types.NewTimeSeries(l.TimeRange, timeutil.Interval(l.Interval))
// 	} else if !r.IsSingleValue() {
// 		result = types.NewTimeSeries(r.TimeRange, timeutil.Interval(r.Interval))
// 	} else {
// 		result = types.NewTimeSeriesWithSingleValue(0)
// 	}
// 	for i := range result.Size() {
// 		result.Put(i, math(l.Get(i), r.Get(i)))
// 	}
// 	return result, false, nil
// }
//
// type arithmeticMinusFunc struct {
// 	baseFunc
// }
//
// func newArithmeticMinusFunc(ctx EvalContext, args []Expression) Func {
// 	return &arithmeticMinusFunc{
// 		baseFunc: baseFunc{ctx: ctx, args: args},
// 	}
// }
//
// func (f *arithmeticMinusFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalInt(row)
// 	rv, _, _ := f.args[1].EvalInt(row)
// 	fmt.Println("minus int.....")
// 	return lv - rv, false, nil
// }
//
// func (f *arithmeticMinusFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
// 	fmt.Println("minus float.....")
// 	return
// }
//
// func (f *arithmeticMinusFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
// 	fmt.Println("minus time series.....")
// 	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
// 		return lv - rv
// 	})
// }
//
// func (f *arithmeticMinusFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalTime(row)
// 	rv, _, _ := f.args[1].EvalDuration(row)
// 	val = lv.Add(-rv)
// 	return
// }
//
// type arithmeticMulFunc struct {
// 	baseFunc
// }
//
// func newArithmeticMulFunc(ctx EvalContext, args []Expression) Func {
// 	return &arithmeticMulFunc{
// 		baseFunc: baseFunc{ctx: ctx, args: args},
// 	}
// }
//
// func (f *arithmeticMulFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalInt(row)
// 	rv, _, _ := f.args[1].EvalInt(row)
// 	fmt.Println("mul int.....")
// 	return lv * rv, false, nil
// }
//
// func (f *arithmeticMulFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
// 	fmt.Println("mul time series.....")
// 	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
// 		return lv * rv
// 	})
// }
//
// type arithmeticDivFunc struct{ baseFunc }
//
// func newArithmeticDivFunc(ctx EvalContext, args []Expression) Func {
// 	return &arithmeticDivFunc{
// 		baseFunc: baseFunc{ctx: ctx, args: args},
// 	}
// }
//
// func (f *arithmeticDivFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalInt(row)
// 	rv, _, _ := f.args[1].EvalInt(row)
// 	fmt.Println("div int.....")
// 	return lv / rv, false, nil
// }
//
// func (f *arithmeticDivFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
// 	fmt.Println("div time series.....")
// 	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
// 		if rv == 0 {
// 			return 0
// 		}
// 		return lv / rv
// 	})
// }
//
// type arithmeticModFunc struct {
// 	baseFunc
// }
//
// func newArithmeticModFunc(ctx EvalContext, args []Expression) Func {
// 	return &arithmeticModFunc{
// 		baseFunc: baseFunc{ctx: ctx, args: args},
// 	}
// }
//
// func (f *arithmeticModFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
// 	lv, _, _ := f.args[0].EvalInt(row)
// 	rv, _, _ := f.args[1].EvalInt(row)
//
// 	fmt.Println("mod int.....")
// 	return lv % rv, false, nil
// }
//
// func (f *arithmeticModFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
// 	fmt.Println("mod time series.....")
// 	return evalTimeSeries(row, f.args, func(lv, rv float64) float64 {
// 		if rv == 0 {
// 			return 0
// 		}
// 		return float64(int64(lv) % int64(rv))
// 	})
// }

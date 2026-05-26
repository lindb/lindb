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
	larray "github.com/lindb/arrow/pkg/arrow/array"

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
	// Scalar OP Scalar must be checked first; otherwise the Scalar OP Array branch fires
	// and calls right.Eval(record) on a Constant, which panics.
	case left.ResultType() == Scalar && right.ResultType() == Scalar:
		return evalScalarScalar(record, left, right, f.op)
	case left.ResultType() == Scalar:
		return evalScalarArray(record, left, right, f.op)
	case right.ResultType() == Scalar:
		return evalArrayScalar(record, left, right, f.op)
	case left.ResultType() == Array && right.ResultType() == Array:
		return evalArrayArray(record, left, right, f.op)
	default:
		return nil, fmt.Errorf("unsupported result type: %v, %v", left.ResultType(), right.ResultType())
	}
}

// applyOp applies the arithmetic operator to two float64 values.
func applyOp(lv, rv float64, op string) float64 {
	switch op {
	case "plus":
		return lv + rv
	case "minus":
		return lv - rv
	case "mul":
		return lv * rv
	case "div":
		if rv == 0 {
			return 0
		}
		return lv / rv
	case "mod":
		if rv == 0 {
			return 0
		}
		return float64(int64(lv) % int64(rv))
	default:
		panic(fmt.Sprintf("unsupported arithmetic op: %s", op))
	}
}

// scalarToFloat64 converts a Scalar (Int64 or Float64) to float64.
// Handles both integer literals (e.g. "* 100") and float literals (e.g. "* 1.5").
func scalarToFloat64(s scalar.Scalar) float64 {
	switch v := s.(type) {
	case *scalar.Int64:
		return float64(v.Value)
	case *scalar.Float64:
		return v.Value
	default:
		panic(fmt.Sprintf("arithmetic scalar operand must be numeric, got: %T", s))
	}
}

// evalScalarScalar handles scalar OP scalar (e.g. "38/48").
// Computes the result once and broadcasts it as a constant float64 array.
func evalScalarScalar(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	leftScalar, err := left.EvalScalar()
	if err != nil {
		return nil, err
	}
	rightScalar, err := right.EvalScalar()
	if err != nil {
		return nil, err
	}
	val := applyOp(scalarToFloat64(leftScalar), scalarToFloat64(rightScalar), op)
	// Broadcast the constant result to match the number of rows in the record.
	numRows := int(record.NumRows())
	result := array.NewFloat64Builder(memory.DefaultAllocator)
	defer result.Release()
	result.Reserve(numRows)
	for range numRows {
		result.Append(val)
	}
	return result.NewArray(), nil
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

	rv := scalarToFloat64(rightScalar)

	switch leftTyped := leftArray.(type) {
	case *array.Int64:
		result := array.NewFloat64Builder(memory.DefaultAllocator)
		defer result.Release()
		result.Reserve(leftTyped.Len())
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) {
				result.AppendNull()
			} else {
				result.Append(applyOp(float64(leftTyped.Value(i)), rv, op))
			}
		}
		return result.NewArray(), nil

	case *larray.Aggregation:
		// Direct Aggregation (not wrapped by FilterableRecord).
		extType := leftTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(leftTyped.Value(i), rv, op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.Generic[float64]:
		// Aggregation wrapped by FilterableRecord.ToGenericWithMask — the storage still
		// carries the original AggregationType so we can reconstruct the correct kind.
		extType, ok := leftTyped.Storage().DataType().(arrow.ExtensionType)
		if !ok {
			// Fallback: plain float64 array result.
			result := array.NewFloat64Builder(memory.DefaultAllocator)
			defer result.Release()
			result.Reserve(leftTyped.Len())
			for i := range leftTyped.Len() {
				if leftTyped.IsNull(i) {
					result.AppendNull()
				} else {
					result.Append(applyOp(leftTyped.Value(i), rv, op))
				}
			}
			return result.NewArray(), nil
		}
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(leftTyped.Value(i), rv, op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.TimeSeries:
		// TimeSeries column (range query with timestamp): apply the scalar factor to
		// every data point in each series while preserving start/end/interval.
		extType := leftTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		tsBuilder := larray.NewTimeSeriesBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) {
				tsBuilder.AppendNull()
			} else {
				vals := leftTyped.Values(i)
				for j := range vals {
					vals[j] = applyOp(vals[j], rv, op)
				}
				tsBuilder.Append(leftTyped.Start(i), leftTyped.End(i), leftTyped.Interval(i), vals)
			}
		}
		return extBuilder.NewExtensionArray(), nil

	default:
		panic(fmt.Sprintf("unsupported array type in array-scalar arithmetic: %T", leftArray))
	}
}

// evalScalarArray handles scalar OP array (scalar on the left, e.g. "100 - last(a)").
// For commutative ops (plus/mul) the scalar side is swapped; for non-commutative ops
// (minus/div/mod) the scalar is kept as lv so order is preserved.
func evalScalarArray(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	leftScalar, err := left.EvalScalar()
	if err != nil {
		return nil, err
	}
	rightArray, err := right.Eval(record)
	if err != nil {
		return nil, err
	}
	defer rightArray.Release()

	lv := scalarToFloat64(leftScalar)

	switch rightTyped := rightArray.(type) {
	case *array.Int64:
		result := array.NewFloat64Builder(memory.DefaultAllocator)
		defer result.Release()
		result.Reserve(rightTyped.Len())
		for i := range rightTyped.Len() {
			if rightTyped.IsNull(i) {
				result.AppendNull()
			} else {
				result.Append(applyOp(lv, float64(rightTyped.Value(i)), op))
			}
		}
		return result.NewArray(), nil

	case *larray.Aggregation:
		extType := rightTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range rightTyped.Len() {
			if rightTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(lv, rightTyped.Value(i), op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.Generic[float64]:
		extType, ok := rightTyped.Storage().DataType().(arrow.ExtensionType)
		if !ok {
			result := array.NewFloat64Builder(memory.DefaultAllocator)
			defer result.Release()
			result.Reserve(rightTyped.Len())
			for i := range rightTyped.Len() {
				if rightTyped.IsNull(i) {
					result.AppendNull()
				} else {
					result.Append(applyOp(lv, rightTyped.Value(i), op))
				}
			}
			return result.NewArray(), nil
		}
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range rightTyped.Len() {
			if rightTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(lv, rightTyped.Value(i), op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.TimeSeries:
		extType := rightTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		tsBuilder := larray.NewTimeSeriesBuilder(extBuilder)
		for i := range rightTyped.Len() {
			if rightTyped.IsNull(i) {
				tsBuilder.AppendNull()
			} else {
				vals := rightTyped.Values(i)
				for j := range vals {
					vals[j] = applyOp(lv, vals[j], op)
				}
				tsBuilder.Append(rightTyped.Start(i), rightTyped.End(i), rightTyped.Interval(i), vals)
			}
		}
		return extBuilder.NewExtensionArray(), nil

	default:
		panic(fmt.Sprintf("unsupported array type in scalar-array arithmetic: %T", rightArray))
	}
}

// evalArrayArray handles array OP array (e.g. "last(a) - last(b)").
// Both arrays must have the same length and compatible types.
func evalArrayArray(record arrow.RecordBatch, left, right Expression, op string) (arrow.Array, error) {
	leftArray, err := left.Eval(record)
	if err != nil {
		return nil, err
	}
	defer leftArray.Release()

	rightArray, err := right.Eval(record)
	if err != nil {
		return nil, err
	}
	defer rightArray.Release()

	switch leftTyped := leftArray.(type) {
	case *array.Int64:
		rightTyped, ok := rightArray.(*array.Int64)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		result := array.NewFloat64Builder(memory.DefaultAllocator)
		defer result.Release()
		result.Reserve(leftTyped.Len())
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) || rightTyped.IsNull(i) {
				result.AppendNull()
			} else {
				result.Append(applyOp(float64(leftTyped.Value(i)), float64(rightTyped.Value(i)), op))
			}
		}
		return result.NewArray(), nil

	case *larray.Aggregation:
		// Both sides are direct Aggregation arrays (same extension kind).
		rightTyped, ok := rightArray.(*larray.Aggregation)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		extType := leftTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) || rightTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(leftTyped.Value(i), rightTyped.Value(i), op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.Generic[float64]:
		// Both sides wrapped by FilterableRecord (ToGenericWithMask).
		rightTyped, ok := rightArray.(*larray.Generic[float64])
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		extType, hasExt := leftTyped.Storage().DataType().(arrow.ExtensionType)
		if !hasExt {
			result := array.NewFloat64Builder(memory.DefaultAllocator)
			defer result.Release()
			result.Reserve(leftTyped.Len())
			for i := range leftTyped.Len() {
				if leftTyped.IsNull(i) || rightTyped.IsNull(i) {
					result.AppendNull()
				} else {
					result.Append(applyOp(leftTyped.Value(i), rightTyped.Value(i), op))
				}
			}
			return result.NewArray(), nil
		}
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		aggBuilder := larray.NewAggregationBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) || rightTyped.IsNull(i) {
				aggBuilder.AppendNull()
			} else {
				aggBuilder.Append(applyOp(leftTyped.Value(i), rightTyped.Value(i), op))
			}
		}
		return extBuilder.NewExtensionArray(), nil

	case *larray.TimeSeries:
		// Both sides are TimeSeries (range query result): apply op point-by-point.
		rightTyped, ok := rightArray.(*larray.TimeSeries)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		extType := leftTyped.DataType().(arrow.ExtensionType)
		extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
		defer extBuilder.Release()
		tsBuilder := larray.NewTimeSeriesBuilder(extBuilder)
		for i := range leftTyped.Len() {
			if leftTyped.IsNull(i) || rightTyped.IsNull(i) {
				tsBuilder.AppendNull()
			} else {
				lVals := leftTyped.Values(i)
				rVals := rightTyped.Values(i)
				n := min(len(lVals), len(rVals))
				result := make([]float64, n)
				for j := range n {
					result[j] = applyOp(lVals[j], rVals[j], op)
				}
				tsBuilder.Append(leftTyped.Start(i), leftTyped.End(i), leftTyped.Interval(i), result)
			}
		}
		return extBuilder.NewExtensionArray(), nil

	default:
		panic(fmt.Sprintf("unsupported array type in array-array arithmetic: %T", leftArray))
	}
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

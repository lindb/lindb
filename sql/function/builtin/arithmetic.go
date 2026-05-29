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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// validArithmeticOps is the set of operator strings accepted by arithmeticInstance.
var validArithmeticOps = map[string]bool{
	"plus": true, "minus": true, "mul": true, "div": true, "mod": true,
}

// arithmeticInstance is a per-query VectorFunc for binary arithmetic operators.
//
// Constant arguments are detected at construction time (via EvalScalar probe) and
// cached in lsv / rsv.  This means constant-OP-array and array-OP-constant paths
// never re-evaluate the constant argument on subsequent record batches.
type arithmeticInstance struct {
	op   string
	args []function.Expr // baked in at construction; column refs are evaluated per-batch

	lsv scalar.Scalar // non-nil when left arg is a constant
	rsv scalar.Scalar // non-nil when right arg is a constant
}

// newArithmeticFactory returns a VectorFuncFactory for the given arithmetic operator.
func newArithmeticFactory(op string) function.VectorFuncFactory {
	if !validArithmeticOps[op] {
		panic(fmt.Sprintf("unknown arithmetic op: %q", op))
	}
	return func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
		if len(args) != 2 {
			panic(fmt.Sprintf("arithmetic op requires exactly 2 args, got %d", len(args)))
		}
		f := &arithmeticInstance{op: op, args: args}
		// Probe once at construction: cache constant args to avoid per-record re-evaluation.
		f.lsv, _ = args[0].EvalScalar()
		f.rsv, _ = args[1].EvalScalar()
		return f
	}
}

func (f *arithmeticInstance) EvalScalar() (scalar.Scalar, error) {
	if f.lsv == nil || f.rsv == nil {
		return nil, errors.New("arithmetic: not a constant expression")
	}
	val := applyOp(scalarToFloat64(f.lsv), scalarToFloat64(f.rsv), f.op)
	return scalar.NewFloat64Scalar(val), nil
}

func (f *arithmeticInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	switch {
	case f.lsv != nil && f.rsv != nil:
		// Scalar OP Scalar: broadcast constant result (cached lsv/rsv, no per-batch work).
		val := applyOp(scalarToFloat64(f.lsv), scalarToFloat64(f.rsv), f.op)
		return broadcastFloat64(val, int(record.NumRows())), nil

	case f.lsv != nil:
		// Scalar OP Array: right is a column, left value is cached.
		rightArray, err := f.args[1].Eval(record)
		if err != nil {
			return nil, err
		}
		defer rightArray.Release()
		return evalScalarOpArray(scalarToFloat64(f.lsv), rightArray, f.op)

	case f.rsv != nil:
		// Array OP Scalar: left is a column, right value is cached.
		leftArray, err := f.args[0].Eval(record)
		if err != nil {
			return nil, err
		}
		defer leftArray.Release()
		return evalArrayOpScalar(leftArray, scalarToFloat64(f.rsv), f.op)

	default:
		// Array OP Array: both are columns.
		leftArray, err := f.args[0].Eval(record)
		if err != nil {
			return nil, err
		}
		defer leftArray.Release()
		rightArray, err := f.args[1].Eval(record)
		if err != nil {
			return nil, err
		}
		defer rightArray.Release()
		return evalArrayOpArray(leftArray, rightArray, f.op)
	}
}

// ── Pure computation helpers ───────────────────────────────────────────────────

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

// scalarToFloat64 converts a numeric Scalar to float64.
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

// broadcastFloat64 creates a float64 array containing val repeated numRows times.
func broadcastFloat64(val float64, numRows int) arrow.Array {
	b := array.NewFloat64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(numRows)
	for range numRows {
		b.Append(val)
	}
	return b.NewArray()
}

// buildAggExtArray builds an Aggregation extension array applying compute(i) per element.
func buildAggExtArray(
	extType arrow.ExtensionType, n int,
	isNull func(int) bool,
	compute func(int) float64,
) (arrow.Array, error) {
	extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
	defer extBuilder.Release()
	b := larray.NewAggregationBuilder(extBuilder)
	for i := range n {
		if isNull(i) {
			b.AppendNull()
		} else {
			b.Append(compute(i))
		}
	}
	return extBuilder.NewExtensionArray(), nil
}

// buildTSExtArray builds a TimeSeries extension array row by row.
func buildTSExtArray(
	extType arrow.ExtensionType, n int,
	isNull func(int) bool,
	appendRow func(b *larray.TimeSeriesBuilder, i int) error,
) (arrow.Array, error) {
	extBuilder := array.NewExtensionBuilder(memory.DefaultAllocator, extType)
	defer extBuilder.Release()
	b := larray.NewTimeSeriesBuilder(extBuilder)
	for i := range n {
		if isNull(i) {
			b.AppendNull()
		} else {
			if err := appendRow(b, i); err != nil {
				return nil, err
			}
		}
	}
	return extBuilder.NewExtensionArray(), nil
}

// ── Array evaluation helpers ───────────────────────────────────────────────────

func evalArrayOpScalar(leftArray arrow.Array, rv float64, op string) (arrow.Array, error) {
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
		return buildAggExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			leftTyped.IsNull,
			func(i int) float64 { return applyOp(leftTyped.Value(i), rv, op) },
		)
	case *larray.Generic[float64]:
		extType, ok := leftTyped.Storage().DataType().(arrow.ExtensionType)
		if !ok {
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
		return buildAggExtArray(
			extType, leftTyped.Len(),
			leftTyped.IsNull,
			func(i int) float64 { return applyOp(leftTyped.Value(i), rv, op) },
		)
	case *larray.TimeSeries:
		return buildTSExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			leftTyped.IsNull,
			func(b *larray.TimeSeriesBuilder, i int) error {
				vals := leftTyped.Values(i)
				for j := range vals {
					vals[j] = applyOp(vals[j], rv, op)
				}
				b.Append(leftTyped.Start(i), leftTyped.End(i), leftTyped.Interval(i), vals)
				return nil
			},
		)
	default:
		panic(fmt.Sprintf("unsupported array type in array-scalar arithmetic: %T", leftArray))
	}
}

func evalScalarOpArray(lv float64, rightArray arrow.Array, op string) (arrow.Array, error) {
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
		return buildAggExtArray(
			rightTyped.DataType().(arrow.ExtensionType), rightTyped.Len(),
			rightTyped.IsNull,
			func(i int) float64 { return applyOp(lv, rightTyped.Value(i), op) },
		)
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
		return buildAggExtArray(
			extType, rightTyped.Len(),
			rightTyped.IsNull,
			func(i int) float64 { return applyOp(lv, rightTyped.Value(i), op) },
		)
	case *larray.TimeSeries:
		return buildTSExtArray(
			rightTyped.DataType().(arrow.ExtensionType), rightTyped.Len(),
			rightTyped.IsNull,
			func(b *larray.TimeSeriesBuilder, i int) error {
				vals := rightTyped.Values(i)
				for j := range vals {
					vals[j] = applyOp(lv, vals[j], op)
				}
				b.Append(rightTyped.Start(i), rightTyped.End(i), rightTyped.Interval(i), vals)
				return nil
			},
		)
	default:
		panic(fmt.Sprintf("unsupported array type in scalar-array arithmetic: %T", rightArray))
	}
}

func evalArrayOpArray(leftArray, rightArray arrow.Array, op string) (arrow.Array, error) {
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
		rightTyped, ok := rightArray.(*larray.Aggregation)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		return buildAggExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			func(i int) bool { return leftTyped.IsNull(i) || rightTyped.IsNull(i) },
			func(i int) float64 { return applyOp(leftTyped.Value(i), rightTyped.Value(i), op) },
		)
	case *larray.Generic[float64]:
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
		return buildAggExtArray(
			extType, leftTyped.Len(),
			func(i int) bool { return leftTyped.IsNull(i) || rightTyped.IsNull(i) },
			func(i int) float64 { return applyOp(leftTyped.Value(i), rightTyped.Value(i), op) },
		)
	case *larray.TimeSeries:
		rightTyped, ok := rightArray.(*larray.TimeSeries)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		return buildTSExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			func(i int) bool { return leftTyped.IsNull(i) || rightTyped.IsNull(i) },
			func(b *larray.TimeSeriesBuilder, i int) error {
				lVals := leftTyped.Values(i)
				rVals := rightTyped.Values(i)
				if len(lVals) != len(rVals) {
					return fmt.Errorf("time series length mismatch at row %d: left=%d right=%d",
						i, len(lVals), len(rVals))
				}
				for j := range lVals {
					lVals[j] = applyOp(lVals[j], rVals[j], op)
				}
				b.Append(leftTyped.Start(i), leftTyped.End(i), leftTyped.Interval(i), lVals)
				return nil
			},
		)
	default:
		panic(fmt.Sprintf("unsupported array type in array-array arithmetic: %T", leftArray))
	}
}

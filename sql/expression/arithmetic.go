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

// validArithmeticOps is the set of operator strings accepted by arithmeticFunc.
// Validated at construction so the execution path never needs to handle unknowns.
var validArithmeticOps = map[string]bool{
	"plus": true, "minus": true, "mul": true, "div": true, "mod": true,
}

type arithmeticFunc struct {
	ctx  EvalContext
	args []Expression

	op string
}

func newArithmeticFunc(ctx EvalContext, op string, args []Expression) Func {
	// Fail fast at construction rather than at query execution time.
	if !validArithmeticOps[op] {
		panic(fmt.Sprintf("unknown arithmetic op: %q", op))
	}
	if len(args) != 2 {
		panic(fmt.Sprintf("arithmetic op requires exactly 2 args, got %d", len(args)))
	}
	return &arithmeticFunc{
		ctx:  ctx,
		args: args,
		op:   op,
	}
}

// EvalScalar evaluates the arithmetic expression when both operands are scalars.
// This enables nested scalar expressions (e.g. "(3+4) * last(cpu)" where the
// outer mul calls EvalScalar on the inner plus before branching to evalScalarArray).
func (f *arithmeticFunc) EvalScalar() (scalar.Scalar, error) {
	lv, err := f.args[0].EvalScalar()
	if err != nil {
		return nil, err
	}
	rv, err := f.args[1].EvalScalar()
	if err != nil {
		return nil, err
	}
	result := applyOp(scalarToFloat64(lv), scalarToFloat64(rv), f.op)
	return scalar.NewFloat64Scalar(result), nil
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
// The op string is validated at construction (newArithmeticFunc); the default
// branch is a safety net that should never be reached in practice.
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

// buildAggExtArray builds an Aggregation extension array by applying compute(i) for each element.
// extType is the AggregationType to preserve on the output (e.g. Sum, Last).
// Note: AggregationBuilder is a thin wrapper over Float64Builder; it holds no
// independent reference count, so only extBuilder needs a Release call.
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
// appendRow receives the builder and the current index; it should call b.Append(...)
// and may return an error (e.g. mismatched series lengths) to abort construction.
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
		return buildAggExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			leftTyped.IsNull,
			func(i int) float64 { return applyOp(leftTyped.Value(i), rv, op) },
		)

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
		return buildAggExtArray(
			extType, leftTyped.Len(),
			leftTyped.IsNull,
			func(i int) float64 { return applyOp(leftTyped.Value(i), rv, op) },
		)

	case *larray.TimeSeries:
		// TimeSeries column (range query): apply the scalar factor to every data point
		// while preserving start/end/interval.
		// Values() returns a fresh copy of each series' float64 slice — safe to mutate.
		return buildTSExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			leftTyped.IsNull,
			func(b *larray.TimeSeriesBuilder, i int) error {
				vals := leftTyped.Values(i) // copy; safe to mutate
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
		// Values() returns a copy; safe to mutate.
		return buildTSExtArray(
			rightTyped.DataType().(arrow.ExtensionType), rightTyped.Len(),
			rightTyped.IsNull,
			func(b *larray.TimeSeriesBuilder, i int) error {
				vals := rightTyped.Values(i) // copy; safe to mutate
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

// evalArrayArray handles array OP array (e.g. "last(a) - last(b)").
// Both arrays must have the same length and compatible element types.
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
		// The result extType is taken from the left operand. For operations between
		// different aggregation kinds (e.g. Sum - Last) the semantic is undefined;
		// we default to the left side's kind as a reasonable convention.
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
		// extType is taken from the left side (see Aggregation note above).
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
		// Both sides are TimeSeries (range query result): apply op point-by-point.
		// Mismatched lengths indicate misaligned query windows/steps and are treated as
		// an error rather than silently truncating data.
		rightTyped, ok := rightArray.(*larray.TimeSeries)
		if !ok {
			return nil, fmt.Errorf("array-array arithmetic type mismatch: %T OP %T", leftArray, rightArray)
		}
		return buildTSExtArray(
			leftTyped.DataType().(arrow.ExtensionType), leftTyped.Len(),
			func(i int) bool { return leftTyped.IsNull(i) || rightTyped.IsNull(i) },
			func(b *larray.TimeSeriesBuilder, i int) error {
				lVals := leftTyped.Values(i)  // copy; safe to mutate
				rVals := rightTyped.Values(i) // copy
				if len(lVals) != len(rVals) {
					return fmt.Errorf(
						"time series length mismatch at row %d: left=%d right=%d",
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

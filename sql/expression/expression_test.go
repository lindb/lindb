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
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

// ─────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────

// int64Batch builds a single-column int64 RecordBatch with the given values.
func int64Batch(vals []int64) arrow.RecordBatch {
	schema := arrow.NewSchema([]arrow.Field{{Name: "a", Type: arrow.PrimitiveTypes.Int64}}, nil)
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return array.NewRecordBatch(schema, []arrow.Array{b.NewArray()}, int64(len(vals)))
}

// fixedBoolExpr is a test-only Expression that always returns a predetermined
// boolean array, used to drive Logical and Not tests without Arrow plumbing.
type fixedBoolExpr struct {
	vals []bool
}

func (f *fixedBoolExpr) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("not supported")
}

func (f *fixedBoolExpr) Eval(_ arrow.RecordBatch) (arrow.Array, error) {
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	for _, v := range f.vals {
		b.Append(v)
	}
	return b.NewArray(), nil
}

func (f *fixedBoolExpr) ResultType() ResultType { return Array }
func (f *fixedBoolExpr) String() string          { return "fixed" }

// boolResults extracts the boolean values from a *array.Boolean.
func boolResults(arr arrow.Array) []bool {
	ba := arr.(*array.Boolean)
	out := make([]bool, ba.Len())
	for i := range ba.Len() {
		out[i] = ba.Value(i)
	}
	return out
}

// ─────────────────────────────────────────────────────────────
// Constant
// ─────────────────────────────────────────────────────────────

func TestConstant_EvalScalar(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"int64", int64(42)},
		{"float64", float64(3.14)},
		{"string", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConstant(nil, tt.value, Scalar)
			s, err := c.EvalScalar()
			require.NoError(t, err)
			require.NotNil(t, s)
		})
	}
}

func TestConstant_Eval_BroadcastsToAllRows(t *testing.T) {
	record := int64Batch([]int64{1, 2, 3}) // 3-row batch

	t.Run("int64", func(t *testing.T) {
		c := NewConstant(nil, int64(7), Scalar)
		arr, err := c.Eval(record)
		require.NoError(t, err)
		got := arr.(*array.Int64)
		require.Equal(t, 3, got.Len())
		for i := range 3 {
			assert.Equal(t, int64(7), got.Value(i))
		}
	})

	t.Run("float64", func(t *testing.T) {
		c := NewConstant(nil, float64(2.5), Scalar)
		arr, err := c.Eval(record)
		require.NoError(t, err)
		got := arr.(*array.Float64)
		require.Equal(t, 3, got.Len())
		for i := range 3 {
			assert.Equal(t, float64(2.5), got.Value(i))
		}
	})

	t.Run("string", func(t *testing.T) {
		c := NewConstant(nil, "foo", Scalar)
		arr, err := c.Eval(record)
		require.NoError(t, err)
		got := arr.(*array.String)
		require.Equal(t, 3, got.Len())
		for i := range 3 {
			assert.Equal(t, "foo", got.Value(i))
		}
	})
}

func TestConstant_ResultType(t *testing.T) {
	assert.Equal(t, Scalar, NewConstant(nil, int64(1), Scalar).ResultType())
	assert.Equal(t, Array, NewConstant(nil, int64(1), Array).ResultType())
}

func TestConstant_String(t *testing.T) {
	assert.Equal(t, "42", NewConstant(nil, int64(42), Scalar).String())
	assert.Equal(t, "hello", NewConstant(nil, "hello", Scalar).String())
}

// ─────────────────────────────────────────────────────────────
// Column
// ─────────────────────────────────────────────────────────────

func TestColumn_Eval_ExtractsColumnByIndex(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "x", Type: arrow.PrimitiveTypes.Int64},
		{Name: "y", Type: arrow.PrimitiveTypes.Float64},
	}, nil)
	ib := array.NewInt64Builder(memory.DefaultAllocator)
	defer ib.Release()
	ib.AppendValues([]int64{10, 20}, nil)
	fb := array.NewFloat64Builder(memory.DefaultAllocator)
	defer fb.Release()
	fb.AppendValues([]float64{1.1, 2.2}, nil)
	record := array.NewRecordBatch(schema, []arrow.Array{ib.NewArray(), fb.NewArray()}, 2)

	// Extract the second column (index 1 = float64).
	col := NewColumn(nil, "y", 1, Array)
	arr, err := col.Eval(record)
	require.NoError(t, err)
	got := arr.(*array.Float64)
	require.Equal(t, 2, got.Len())
	assert.InDelta(t, 1.1, got.Value(0), 1e-9)
	assert.InDelta(t, 2.2, got.Value(1), 1e-9)
}

func TestColumn_EvalScalar_ReturnsError(t *testing.T) {
	col := NewColumn(nil, "x", 0, Array)
	_, err := col.EvalScalar()
	assert.Error(t, err)
}

func TestColumn_String(t *testing.T) {
	assert.Equal(t, "myCol", NewColumn(nil, "myCol", 0, Array).String())
}

// ─────────────────────────────────────────────────────────────
// Comparison
// ─────────────────────────────────────────────────────────────

// twoColInt64Batch creates a batch with two int64 columns for comparison tests.
func twoColInt64Batch(left, right []int64) arrow.RecordBatch {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "l", Type: arrow.PrimitiveTypes.Int64},
		{Name: "r", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	lb := array.NewInt64Builder(memory.DefaultAllocator)
	defer lb.Release()
	lb.AppendValues(left, nil)
	rb := array.NewInt64Builder(memory.DefaultAllocator)
	defer rb.Release()
	rb.AppendValues(right, nil)
	return array.NewRecordBatch(schema, []arrow.Array{lb.NewArray(), rb.NewArray()}, int64(len(left)))
}

// evalIntCmp is a shorthand helper for two-int64-column comparison tests.
func evalIntCmp(t *testing.T, left, right []int64, op tree.ComparisonOperator) []bool {
	t.Helper()
	record := twoColInt64Batch(left, right)
	cmp := NewComparison(nil, op, NewColumn(nil, "l", 0, Array), NewColumn(nil, "r", 1, Array))
	arr, err := cmp.Eval(record)
	require.NoError(t, err)
	return boolResults(arr)
}

func TestComparison_Int64_AllOperators(t *testing.T) {
	// rows: (1,2), (2,2), (3,2), (3,3), (3,4)
	lefts := []int64{1, 2, 3, 3, 3}
	rights := []int64{2, 2, 2, 3, 4}

	assert.Equal(t, []bool{false, true, false, true, false}, evalIntCmp(t, lefts, rights, tree.ComparisonEQ))
	assert.Equal(t, []bool{true, false, true, false, true}, evalIntCmp(t, lefts, rights, tree.ComparisonNEQ))
	assert.Equal(t, []bool{false, false, true, false, false}, evalIntCmp(t, lefts, rights, tree.ComparisonGT))
	assert.Equal(t, []bool{false, true, true, true, false}, evalIntCmp(t, lefts, rights, tree.ComparisonGTE))
	assert.Equal(t, []bool{true, false, false, false, true}, evalIntCmp(t, lefts, rights, tree.ComparisonLT))
	assert.Equal(t, []bool{true, true, false, true, true}, evalIntCmp(t, lefts, rights, tree.ComparisonLTE))
}

func TestComparison_Float64_AllOperators(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "l", Type: arrow.PrimitiveTypes.Float64},
		{Name: "r", Type: arrow.PrimitiveTypes.Float64},
	}, nil)
	lb := array.NewFloat64Builder(memory.DefaultAllocator)
	defer lb.Release()
	lb.AppendValues([]float64{1.0, 2.0, 3.0}, nil)
	rb := array.NewFloat64Builder(memory.DefaultAllocator)
	defer rb.Release()
	rb.AppendValues([]float64{2.0, 2.0, 2.0}, nil)
	record := array.NewRecordBatch(schema, []arrow.Array{lb.NewArray(), rb.NewArray()}, 3)

	cmp := NewComparison(nil, tree.ComparisonGT,
		NewColumn(nil, "l", 0, Array), NewColumn(nil, "r", 1, Array))
	arr, err := cmp.Eval(record)
	require.NoError(t, err)
	assert.Equal(t, []bool{false, false, true}, boolResults(arr))
}

// TestComparison_Float64Left_Int64Right_Widening is the regression test for Bug 2:
// HAVING count(*) > 2000 produces float64 left and int64 right; must compare via widening.
func TestComparison_Float64Left_Int64Right_Widening(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "cnt", Type: arrow.PrimitiveTypes.Float64},
		{Name: "threshold", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	fb := array.NewFloat64Builder(memory.DefaultAllocator)
	defer fb.Release()
	fb.AppendValues([]float64{1500, 2500, 2000}, nil)
	ib := array.NewInt64Builder(memory.DefaultAllocator)
	defer ib.Release()
	ib.AppendValues([]int64{2000, 2000, 2000}, nil)
	record := array.NewRecordBatch(schema, []arrow.Array{fb.NewArray(), ib.NewArray()}, 3)

	cmp := NewComparison(nil, tree.ComparisonGT,
		NewColumn(nil, "cnt", 0, Array), NewColumn(nil, "threshold", 1, Array))
	arr, err := cmp.Eval(record)
	require.NoError(t, err)
	// 1500>2000=false, 2500>2000=true, 2000>2000=false
	assert.Equal(t, []bool{false, true, false}, boolResults(arr))
}

// TestComparison_Int64Left_Float64Right_Widening exercises the symmetric case: int64 left
// compared against float64 right (also widened to float64 for comparison).
func TestComparison_Int64Left_Float64Right_Widening(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "a", Type: arrow.PrimitiveTypes.Int64},
		{Name: "b", Type: arrow.PrimitiveTypes.Float64},
	}, nil)
	ib := array.NewInt64Builder(memory.DefaultAllocator)
	defer ib.Release()
	ib.AppendValues([]int64{1, 3, 2}, nil)
	fb := array.NewFloat64Builder(memory.DefaultAllocator)
	defer fb.Release()
	fb.AppendValues([]float64{2.0, 2.0, 2.0}, nil)
	record := array.NewRecordBatch(schema, []arrow.Array{ib.NewArray(), fb.NewArray()}, 3)

	cmp := NewComparison(nil, tree.ComparisonGT,
		NewColumn(nil, "a", 0, Array), NewColumn(nil, "b", 1, Array))
	arr, err := cmp.Eval(record)
	require.NoError(t, err)
	// 1>2.0=false, 3>2.0=true, 2>2.0=false
	assert.Equal(t, []bool{false, true, false}, boolResults(arr))
}

func TestComparison_String_EQandNEQ(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "l", Type: arrow.BinaryTypes.String},
		{Name: "r", Type: arrow.BinaryTypes.String},
	}, nil)
	lb := array.NewStringBuilder(memory.DefaultAllocator)
	defer lb.Release()
	lb.AppendValues([]string{"a", "b", "b"}, nil)
	rb := array.NewStringBuilder(memory.DefaultAllocator)
	defer rb.Release()
	rb.AppendValues([]string{"b", "b", "a"}, nil)
	record := array.NewRecordBatch(schema, []arrow.Array{lb.NewArray(), rb.NewArray()}, 3)

	left := NewColumn(nil, "l", 0, Array)
	right := NewColumn(nil, "r", 1, Array)

	for _, tc := range []struct {
		op   tree.ComparisonOperator
		want []bool
	}{
		{tree.ComparisonEQ, []bool{false, true, false}},
		{tree.ComparisonNEQ, []bool{true, false, true}},
		{tree.ComparisonLT, []bool{true, false, false}},
		{tree.ComparisonGT, []bool{false, false, true}},
	} {
		cmp := NewComparison(nil, tc.op, left, right)
		arr, err := cmp.Eval(record)
		require.NoError(t, err, "op=%s", tc.op)
		assert.Equal(t, tc.want, boolResults(arr), "op=%s", tc.op)
	}
}

func TestComparison_EvalScalar_ReturnsError(t *testing.T) {
	cmp := NewComparison(nil, tree.ComparisonEQ,
		NewConstant(nil, int64(1), Scalar), NewConstant(nil, int64(1), Scalar))
	_, err := cmp.EvalScalar()
	assert.Error(t, err)
}

func TestComparison_UnsupportedLeftType_ReturnsError(t *testing.T) {
	// Boolean left operand — not supported by compareXxx helpers.
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "b", Type: arrow.FixedWidthTypes.Boolean},
		{Name: "n", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	bb := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer bb.Release()
	bb.Append(true)
	ib := array.NewInt64Builder(memory.DefaultAllocator)
	defer ib.Release()
	ib.Append(1)
	record := array.NewRecordBatch(schema, []arrow.Array{bb.NewArray(), ib.NewArray()}, 1)

	cmp := NewComparison(nil, tree.ComparisonEQ,
		NewColumn(nil, "b", 0, Array), NewColumn(nil, "n", 1, Array))
	_, err := cmp.Eval(record)
	assert.Error(t, err)
}

func TestComparison_String_ContainsOperator(t *testing.T) {
	cmp := NewComparison(nil, tree.ComparisonEQ,
		NewConstant(nil, int64(1), Scalar), NewConstant(nil, int64(2), Scalar))
	assert.Contains(t, cmp.String(), "=")
}

// ─────────────────────────────────────────────────────────────
// Logical (AND / OR) and Not
// ─────────────────────────────────────────────────────────────

func TestLogical_AND(t *testing.T) {
	record := int64Batch([]int64{1, 2, 3, 4})
	// rows: T&&T=T, T&&F=F, F&&T=F, F&&F=F
	logical := NewLogical(nil, tree.LogicalAND, []Expression{
		&fixedBoolExpr{vals: []bool{true, true, false, false}},
		&fixedBoolExpr{vals: []bool{true, false, true, false}},
	})
	arr, err := logical.Eval(record)
	require.NoError(t, err)
	assert.Equal(t, []bool{true, false, false, false}, boolResults(arr))
}

func TestLogical_OR(t *testing.T) {
	record := int64Batch([]int64{1, 2, 3, 4})
	// rows: T||T=T, T||F=T, F||T=T, F||F=F
	logical := NewLogical(nil, tree.LogicalOR, []Expression{
		&fixedBoolExpr{vals: []bool{true, true, false, false}},
		&fixedBoolExpr{vals: []bool{true, false, true, false}},
	})
	arr, err := logical.Eval(record)
	require.NoError(t, err)
	assert.Equal(t, []bool{true, true, true, false}, boolResults(arr))
}

func TestLogical_EvalScalar_ReturnsError(t *testing.T) {
	logical := NewLogical(nil, tree.LogicalAND, nil)
	_, err := logical.EvalScalar()
	assert.Error(t, err)
}

func TestLogical_String(t *testing.T) {
	logical := NewLogical(nil, tree.LogicalAND, []Expression{
		&fixedBoolExpr{}, &fixedBoolExpr{},
	})
	assert.Contains(t, logical.String(), "AND")
}

func TestNot_Negates(t *testing.T) {
	record := int64Batch([]int64{1, 2, 3})
	not := NewNot(nil, &fixedBoolExpr{vals: []bool{true, false, true}})
	arr, err := not.Eval(record)
	require.NoError(t, err)
	assert.Equal(t, []bool{false, true, false}, boolResults(arr))
}

func TestNot_EvalScalar_ReturnsError(t *testing.T) {
	not := NewNot(nil, &fixedBoolExpr{})
	_, err := not.EvalScalar()
	assert.Error(t, err)
}

func TestNot_String(t *testing.T) {
	not := NewNot(nil, &fixedBoolExpr{})
	assert.Contains(t, not.String(), "NOT")
}

// ─────────────────────────────────────────────────────────────
// Cast (transparent pass-through)
// ─────────────────────────────────────────────────────────────

// TestCast_TransparentArrayPassThrough verifies that Cast.Eval() delegates
// directly to its wrapped expression without changing the array type.
func TestCast_TransparentArrayPassThrough(t *testing.T) {
	record := int64Batch([]int64{1, 2, 3})
	// Wrap an int64 constant; Cast declares float64 target but must not convert.
	inner := NewConstant(nil, int64(99), Array)
	cast := NewCast(nil, arrow.PrimitiveTypes.Float64, inner)

	arr, err := cast.Eval(record)
	require.NoError(t, err)
	// Array type must remain int64 — Cast is a deliberate pass-through.
	assert.IsType(t, &array.Int64{}, arr, "Cast must not convert array type")
}

// TestCast_TransparentScalarPassThrough verifies Cast.EvalScalar() delegates
// to the wrapped expression unchanged.
func TestCast_TransparentScalarPassThrough(t *testing.T) {
	inner := NewConstant(nil, int64(7), Scalar)
	cast := NewCast(nil, arrow.PrimitiveTypes.Float64, inner)
	s, err := cast.EvalScalar()
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestCast_String(t *testing.T) {
	inner := NewConstant(nil, int64(1), Scalar)
	cast := NewCast(nil, arrow.PrimitiveTypes.Float64, inner)
	assert.Contains(t, cast.String(), "CAST")
}

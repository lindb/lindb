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

package operator

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

// ─────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────

// staticBoolExpr is a test-only expression.Expression that always returns
// a fixed boolean array, used to inject pre-evaluated predicates into FilterOperator.
type staticBoolExpr struct {
	vals []bool
}

func (s *staticBoolExpr) EvalScalar() (scalar.Scalar, error) {
	return nil, errors.New("not supported")
}

func (s *staticBoolExpr) Eval(_ arrow.RecordBatch) (arrow.Array, error) {
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	for _, v := range s.vals {
		b.Append(v)
	}
	return b.NewArray(), nil
}

func (s *staticBoolExpr) ResultType() expression.ResultType { return expression.Array }
func (s *staticBoolExpr) String() string                    { return "static" }

// errorBoolExpr is a test-only expression that always returns an error from Eval.
type errorBoolExpr struct{}

func (e *errorBoolExpr) EvalScalar() (scalar.Scalar, error) { return nil, errors.New("not supported") }
func (e *errorBoolExpr) Eval(_ arrow.RecordBatch) (arrow.Array, error) {
	return nil, errors.New("predicate eval error")
}
func (e *errorBoolExpr) ResultType() expression.ResultType { return expression.Array }
func (e *errorBoolExpr) String() string                    { return "error" }

// notBoolExpr returns an int64 array instead of boolean, used to test the type-check path.
type notBoolExpr struct{}

func (n *notBoolExpr) EvalScalar() (scalar.Scalar, error) { return nil, errors.New("not supported") }
func (n *notBoolExpr) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	// Return an int64 column, which is not *array.Boolean.
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	for range record.NumRows() {
		b.Append(0)
	}
	return b.NewArray(), nil
}
func (n *notBoolExpr) ResultType() expression.ResultType { return expression.Array }
func (n *notBoolExpr) String() string                    { return "notBool" }

// makeFilterOp creates a FilterOperator with the compiled predicate pre-set,
// bypassing the AST-rewrite step (prepare()).  This lets us test the filter
// method in isolation without needing a real child operator or planner context.
func makeFilterOp(pred expression.Expression) *FilterOperator {
	return &FilterOperator{
		ctx:     context.Background(),
		pred:    pred,
		inbound: NewQueue(make(chan arrow.RecordBatch, 1024)),
	}
}

// makeRecord builds a RecordBatch with the provided columns and schema.
func makeRecord(schema *arrow.Schema, cols []arrow.Array, numRows int64) arrow.RecordBatch {
	return array.NewRecordBatch(schema, cols, numRows)
}

// ─────────────────────────────────────────────────────────────
// Type-agnostic value extractors
//
// FilterOperator wraps its output in larrow.NewFilterableRecord, which calls
// larray.ToGenericWithMask for every column, converting *array.Int64 →
// *array.Generic[int64] (and so on for other types).  Direct concrete-type
// assertions on RecordBatch.Column() therefore panic at runtime.
// These helpers use structural interface assertions instead.
// ─────────────────────────────────────────────────────────────

// getInt64Val reads an int64 value at position i from any array type that
// exposes a Value(int) int64 method, including both *array.Int64 and
// *array.Generic[int64] returned by FilterableRecord.
func getInt64Val(arr arrow.Array, i int) int64 {
	type valuer interface{ Value(int) int64 }
	v, ok := arr.(valuer)
	if !ok {
		panic(fmt.Sprintf("getInt64Val: unexpected type %T", arr))
	}
	return v.Value(i)
}

// getFloat64Val reads a float64 value at position i; accepts *array.Float64
// and *array.Generic[float64].
func getFloat64Val(arr arrow.Array, i int) float64 {
	type valuer interface{ Value(int) float64 }
	v, ok := arr.(valuer)
	if !ok {
		panic(fmt.Sprintf("getFloat64Val: unexpected type %T", arr))
	}
	return v.Value(i)
}

// getStringVal reads a string value at position i; accepts *array.String and
// *array.Generic[string].
func getStringVal(arr arrow.Array, i int) string {
	type valuer interface{ Value(int) string }
	v, ok := arr.(valuer)
	if !ok {
		panic(fmt.Sprintf("getStringVal: unexpected type %T", arr))
	}
	return v.Value(i)
}

// int64Col builds an *array.Int64 from the given values.
func int64Col(vals []int64) arrow.Array {
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

// float64Col builds an *array.Float64 from the given values.
func float64Col(vals []float64) arrow.Array {
	b := array.NewFloat64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

// stringCol builds an *array.String from the given values.
func stringCol(vals []string) arrow.Array {
	b := array.NewStringBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

// boolCol builds an *array.Boolean from the given values.
func boolCol(vals []bool) arrow.Array {
	b := array.NewBooleanBuilder(memory.DefaultAllocator)
	defer b.Release()
	for _, v := range vals {
		b.Append(v)
	}
	return b.NewArray()
}

// ─────────────────────────────────────────────────────────────
// selectRows
// ─────────────────────────────────────────────────────────────

func TestSelectRows_Int64_SelectsCorrectRows(t *testing.T) {
	col := int64Col([]int64{10, 20, 30, 40, 50})
	result, err := SelectRows(col, []int{0, 2, 4})
	require.NoError(t, err)
	got := result.(*array.Int64)
	require.Equal(t, 3, got.Len())
	assert.Equal(t, int64(10), got.Value(0))
	assert.Equal(t, int64(30), got.Value(1))
	assert.Equal(t, int64(50), got.Value(2))
}

func TestSelectRows_Float64_SelectsCorrectRows(t *testing.T) {
	col := float64Col([]float64{1.1, 2.2, 3.3})
	result, err := SelectRows(col, []int{1})
	require.NoError(t, err)
	got := result.(*array.Float64)
	require.Equal(t, 1, got.Len())
	assert.InDelta(t, 2.2, got.Value(0), 1e-9)
}

func TestSelectRows_String_SelectsCorrectRows(t *testing.T) {
	col := stringCol([]string{"a", "b", "c", "d"})
	result, err := SelectRows(col, []int{0, 3})
	require.NoError(t, err)
	got := result.(*array.String)
	require.Equal(t, 2, got.Len())
	assert.Equal(t, "a", got.Value(0))
	assert.Equal(t, "d", got.Value(1))
}

func TestSelectRows_Boolean_SelectsCorrectRows(t *testing.T) {
	col := boolCol([]bool{true, false, true, false})
	result, err := SelectRows(col, []int{0, 2})
	require.NoError(t, err)
	got := result.(*array.Boolean)
	require.Equal(t, 2, got.Len())
	assert.True(t, got.Value(0))
	assert.True(t, got.Value(1))
}

func TestSelectRows_Int64_PreservesNulls(t *testing.T) {
	// Build an int64 column with a null at index 1.
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.Append(10)
	b.AppendNull()
	b.Append(30)
	col := b.NewArray()

	result, err := SelectRows(col, []int{0, 1, 2})
	require.NoError(t, err)
	got := result.(*array.Int64)
	require.Equal(t, 3, got.Len())
	assert.False(t, got.IsNull(0))
	assert.True(t, got.IsNull(1), "null must be preserved")
	assert.False(t, got.IsNull(2))
}

func TestSelectRows_EmptyRowList_ReturnsEmptyArray(t *testing.T) {
	col := int64Col([]int64{1, 2, 3})
	result, err := SelectRows(col, []int{})
	require.NoError(t, err)
	assert.Equal(t, 0, result.Len())
}

func TestSelectRows_UnsupportedType_ReturnsError(t *testing.T) {
	// Use a raw column type not handled by selectRows (e.g. *array.Binary would
	// not be the usual path; here we use the unsupported path via a cast of bool
	// to test the default branch). Actually, use a list array to trigger default.
	bldr := array.NewListBuilder(memory.DefaultAllocator, arrow.PrimitiveTypes.Int64)
	defer bldr.Release()
	bldr.Append(true)
	col := bldr.NewArray()
	_, err := SelectRows(col, []int{0})
	assert.Error(t, err)
}

// ─────────────────────────────────────────────────────────────
// buildEmptyColumn
// ─────────────────────────────────────────────────────────────

func TestBuildEmptyColumn_ReturnsZeroLengthForAllSupportedTypes(t *testing.T) {
	types := []struct {
		name string
		col  arrow.Array
	}{
		{"int64", int64Col([]int64{1, 2, 3})},
		{"float64", float64Col([]float64{1.1, 2.2})},
		{"string", stringCol([]string{"a", "b"})},
		{"bool", boolCol([]bool{true, false})},
	}
	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			empty := BuildEmptyColumn(tt.col)
			assert.Equal(t, 0, empty.Len(), "buildEmptyColumn must produce a zero-length array")
		})
	}
}

// ─────────────────────────────────────────────────────────────
// FilterOperator.filter
// ─────────────────────────────────────────────────────────────

func TestFilterOperator_Filter_AllRowsPass(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "v", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	record := makeRecord(schema, []arrow.Array{int64Col([]int64{10, 20, 30})}, 3)

	// Predicate passes every row.
	op := makeFilterOp(&staticBoolExpr{vals: []bool{true, true, true}})
	result, err := op.filter(record)
	require.NoError(t, err)
	assert.Equal(t, int64(3), result.NumRows())
	assert.Equal(t, int64(10), getInt64Val(result.Column(0), 0))
	assert.Equal(t, int64(30), getInt64Val(result.Column(0), 2))
}

func TestFilterOperator_Filter_SomeRowsPass(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "v", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	record := makeRecord(schema, []arrow.Array{int64Col([]int64{10, 20, 30})}, 3)

	// Only rows 0 and 2 pass the predicate.
	op := makeFilterOp(&staticBoolExpr{vals: []bool{true, false, true}})
	result, err := op.filter(record)
	require.NoError(t, err)
	require.Equal(t, int64(2), result.NumRows())
	assert.Equal(t, int64(10), getInt64Val(result.Column(0), 0))
	assert.Equal(t, int64(30), getInt64Val(result.Column(0), 1))
}

func TestFilterOperator_Filter_NoRowsPass_ReturnsEmptyBatch(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "v", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	record := makeRecord(schema, []arrow.Array{int64Col([]int64{10, 20, 30})}, 3)

	// No rows pass.
	op := makeFilterOp(&staticBoolExpr{vals: []bool{false, false, false}})
	result, err := op.filter(record)
	require.NoError(t, err)
	// Empty batch is returned — downstream operators discard it via NumRows check.
	assert.Equal(t, int64(0), result.NumRows())
}

func TestFilterOperator_Filter_PredicateEvalError_ReturnError(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "v", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	record := makeRecord(schema, []arrow.Array{int64Col([]int64{1})}, 1)

	op := makeFilterOp(&errorBoolExpr{})
	_, err := op.filter(record)
	assert.Error(t, err)
}

func TestFilterOperator_Filter_NonBoolPredicate_ReturnsError(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "v", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	record := makeRecord(schema, []arrow.Array{int64Col([]int64{1})}, 1)

	// Predicate returns int64 instead of boolean.
	op := makeFilterOp(&notBoolExpr{})
	_, err := op.filter(record)
	assert.Error(t, err)
}

func TestFilterOperator_Filter_MultipleColumns_AllPreserved(t *testing.T) {
	// Verify that every column is correctly filtered, not just the first.
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "score", Type: arrow.PrimitiveTypes.Float64},
		{Name: "label", Type: arrow.BinaryTypes.String},
	}, nil)
	record := makeRecord(schema, []arrow.Array{
		int64Col([]int64{1, 2, 3}),
		float64Col([]float64{1.0, 2.0, 3.0}),
		stringCol([]string{"a", "b", "c"}),
	}, 3)

	// Keep only row 1 (index 1).
	op := makeFilterOp(&staticBoolExpr{vals: []bool{false, true, false}})
	result, err := op.filter(record)
	require.NoError(t, err)
	require.Equal(t, int64(1), result.NumRows())
	assert.Equal(t, int64(2), getInt64Val(result.Column(0), 0))
	assert.InDelta(t, 2.0, getFloat64Val(result.Column(1), 0), 1e-9)
	assert.Equal(t, "b", getStringVal(result.Column(2), 0))
}

// ─────────────────────────────────────────────────────────────
// FilterOperator interface methods
// ─────────────────────────────────────────────────────────────

func TestFilterOperator_String(t *testing.T) {
	op := makeFilterOp(&staticBoolExpr{})
	assert.Equal(t, "FilterOperator", op.String())
}

func TestFilterOperator_GetLayout_ProxiesToChild(t *testing.T) {
	// GetLayout must delegate to the child operator's layout.
	expected := []*plan.Symbol{{Name: "x"}}
	child := &mockChildOp{layout: expected}
	op := &FilterOperator{
		ctx:     context.Background(),
		child:   child,
		inbound: NewQueue(make(chan arrow.RecordBatch, 1)),
	}
	assert.Equal(t, expected, op.GetLayout())
}

// mockChildOp is a minimal Operator stub used only to provide a GetLayout answer.
type mockChildOp struct {
	layout []*plan.Symbol
}

func (m *mockChildOp) Run(_ context.Context, _ chan<- arrow.RecordBatch) {}
func (m *mockChildOp) GetLayout() []*plan.Symbol                         { return m.layout }
func (m *mockChildOp) Children() []Operator                              { return nil }
func (m *mockChildOp) GetInbounds() []chan arrow.RecordBatch             { return nil }
func (m *mockChildOp) String() string                                    { return "mock" }

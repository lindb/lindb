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
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// ── test helpers ──────────────────────────────────────────────────────────────

// makeTestRecord creates an Arrow RecordBatch with numRows float64 zero rows.
func makeTestRecord(numRows int) arrow.RecordBatch {
	schema := arrow.NewSchema([]arrow.Field{{Name: "x", Type: arrow.PrimitiveTypes.Float64}}, nil)
	b := array.NewFloat64Builder(memory.DefaultAllocator)
	defer b.Release()
	for range numRows {
		b.Append(0)
	}
	col := b.NewArray()
	defer col.Release()
	return array.NewRecordBatch(schema, []arrow.Array{col}, int64(numRows))
}

// constExpr is a minimal function.Expr that wraps a compile-time constant scalar.
type constExpr struct{ val scalar.Scalar }

func (c *constExpr) Eval(_ arrow.RecordBatch) (arrow.Array, error) {
	return nil, nil // never called in these tests
}
func (c *constExpr) EvalScalar() (scalar.Scalar, error) { return c.val, nil }

func int64Arg(seed int64) function.Expr  { return &constExpr{val: &scalar.Int64{Value: seed}} }
func float64Arg(v float64) function.Expr { return &constExpr{val: &scalar.Float64{Value: v}} }

// newRand creates a randInstance via RandFactory.
func newRand(args ...function.Expr) function.VectorFunc {
	return RandFactory(nil, args)
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestRand_EvalScalar_NoSeed_InRange(t *testing.T) {
	fn := newRand()
	for range 100 {
		v, err := fn.EvalScalar()
		require.NoError(t, err)
		f := scalar.ToFloat64(v)
		assert.True(t, f >= 0 && f < 1, "RAND() must be in [0,1), got %v", f)
	}
}

func TestRand_EvalScalar_WithSeed_Deterministic(t *testing.T) {
	// Same seed → same first value (PRNG reset per-instance).
	v1, err := newRand(int64Arg(42)).EvalScalar()
	require.NoError(t, err)
	v2, err := newRand(int64Arg(42)).EvalScalar()
	require.NoError(t, err)
	assert.Equal(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2), "RAND(42) must be deterministic")
}

func TestRand_EvalScalar_DifferentSeeds_DifferentValues(t *testing.T) {
	v1, _ := newRand(int64Arg(1)).EvalScalar()
	v2, _ := newRand(int64Arg(2)).EvalScalar()
	assert.NotEqual(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2), "RAND(1) != RAND(2)")
}

func TestRand_EvalScalar_SequenceAdvances(t *testing.T) {
	// Consecutive EvalScalar calls on the same instance must advance the PRNG.
	fn := newRand(int64Arg(42))
	v1, _ := fn.EvalScalar()
	v2, _ := fn.EvalScalar()
	assert.NotEqual(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2),
		"consecutive EvalScalar calls must advance the PRNG")
}

func TestRand_Eval_ReturnsCorrectRowCount(t *testing.T) {
	rec := makeTestRecord(5)
	defer rec.Release()

	arr, err := newRand().Eval(rec)
	require.NoError(t, err)
	defer arr.Release()
	assert.Equal(t, 5, arr.Len())
}

func TestRand_Eval_AllValuesInRange(t *testing.T) {
	rec := makeTestRecord(50)
	defer rec.Release()

	arr, err := newRand().Eval(rec)
	require.NoError(t, err)
	defer arr.Release()

	fa, ok := arr.(*array.Float64)
	require.True(t, ok, "result must be *array.Float64")
	for i := range 50 {
		v := fa.Value(i)
		assert.True(t, v >= 0 && v < 1, "row %d: RAND() must be in [0,1), got %v", i, v)
	}
}

func TestRand_Eval_SameSeed_SameSequence(t *testing.T) {
	rec := makeTestRecord(10)
	defer rec.Release()

	arr1, err := newRand(int64Arg(99)).Eval(rec)
	require.NoError(t, err)
	defer arr1.Release()

	arr2, err := newRand(int64Arg(99)).Eval(rec)
	require.NoError(t, err)
	defer arr2.Release()

	fa1 := arr1.(*array.Float64)
	fa2 := arr2.(*array.Float64)
	for i := range 10 {
		assert.Equal(t, fa1.Value(i), fa2.Value(i),
			"RAND(99) row %d must match across same-seed instances", i)
	}
}

func TestRand_Eval_MultipleCallsPreserveState(t *testing.T) {
	// Calling Eval twice on the same instance advances the PRNG (sequence continues).
	rec := makeTestRecord(3)
	defer rec.Release()

	fn := newRand(int64Arg(7))

	arr1, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr1.Release()

	arr2, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr2.Release()

	assert.Equal(t, 3, arr2.Len())
	fa2 := arr2.(*array.Float64)
	for i := range 3 {
		v := fa2.Value(i)
		assert.True(t, v >= 0 && v < 1, "row %d must be in [0,1)", i)
	}
}

func TestRand_Eval_ZeroRows(t *testing.T) {
	rec := makeTestRecord(0)
	defer rec.Release()

	arr, err := newRand().Eval(rec)
	require.NoError(t, err)
	defer arr.Release()
	assert.Equal(t, 0, arr.Len())
}

func TestRand_FloatSeed_NoPanic(t *testing.T) {
	fn := newRand(float64Arg(3.14))
	v, err := fn.EvalScalar()
	require.NoError(t, err)
	f := scalar.ToFloat64(v)
	assert.True(t, f >= 0 && f < 1, "RAND(3.14) must be in [0,1), got %v", f)
}

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
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/spi/scalar"
)

// makeRecord creates an Arrow RecordBatch with numRows rows of float64 zeros for testing.
func makeRecord(numRows int) arrow.RecordBatch {
	schema := arrow.NewSchema([]arrow.Field{{Name: "x", Type: arrow.PrimitiveTypes.Float64}}, nil)
	b := array.NewFloat64Builder(memory.DefaultAllocator)
	defer b.Release()
	for i := 0; i < numRows; i++ {
		b.Append(0)
	}
	col := b.NewArray()
	defer col.Release()
	return array.NewRecordBatch(schema, []arrow.Array{col}, int64(numRows))
}

// seedArg wraps a constant int64 expression used as the seed argument.
func seedArg(seed int64) Expression {
	return NewConstant(nil, seed, Scalar)
}

func TestRandFunc_EvalScalar_NoSeed_InRange(t *testing.T) {
	fn := newRandFunc(nil, nil)
	for range 100 {
		v, err := fn.EvalScalar()
		require.NoError(t, err)
		f := scalar.ToFloat64(v)
		assert.True(t, f >= 0 && f < 1, "rand() must be in [0, 1), got %v", f)
	}
}

func TestRandFunc_EvalScalar_WithSeed_Deterministic(t *testing.T) {
	// Two instances with the same seed must produce the same value.
	fn1 := newRandFunc(nil, []Expression{seedArg(42)})
	fn2 := newRandFunc(nil, []Expression{seedArg(42)})

	v1, err := fn1.EvalScalar()
	require.NoError(t, err)
	v2, err := fn2.EvalScalar()
	require.NoError(t, err)

	assert.Equal(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2),
		"RAND(42) must return the same value for the same seed")
}

func TestRandFunc_EvalScalar_DifferentSeeds_DifferentValues(t *testing.T) {
	fn1 := newRandFunc(nil, []Expression{seedArg(1)})
	fn2 := newRandFunc(nil, []Expression{seedArg(2)})

	v1, err := fn1.EvalScalar()
	require.NoError(t, err)
	v2, err := fn2.EvalScalar()
	require.NoError(t, err)

	assert.NotEqual(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2),
		"RAND(1) and RAND(2) should produce different values")
}

func TestRandFunc_Eval_ReturnsCorrectRowCount(t *testing.T) {
	const numRows = 5
	rec := makeRecord(numRows)
	defer rec.Release()

	fn := newRandFunc(nil, nil)
	arr, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr.Release()

	assert.Equal(t, numRows, arr.Len(), "Eval must return one value per row")
}

func TestRandFunc_Eval_AllValuesInRange(t *testing.T) {
	const numRows = 50
	rec := makeRecord(numRows)
	defer rec.Release()

	fn := newRandFunc(nil, nil)
	arr, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr.Release()

	floatArr, ok := arr.(*array.Float64)
	require.True(t, ok, "Eval result must be *array.Float64")

	for i := range numRows {
		v := floatArr.Value(i)
		assert.True(t, v >= 0 && v < 1, "row %d: rand() must be in [0, 1), got %v", i, v)
	}
}

func TestRandFunc_Eval_WithSeed_SameSequence(t *testing.T) {
	const numRows = 10
	rec := makeRecord(numRows)
	defer rec.Release()

	fn1 := newRandFunc(nil, []Expression{seedArg(99)})
	fn2 := newRandFunc(nil, []Expression{seedArg(99)})

	arr1, err := fn1.Eval(rec)
	require.NoError(t, err)
	defer arr1.Release()

	arr2, err := fn2.Eval(rec)
	require.NoError(t, err)
	defer arr2.Release()

	fa1, ok1 := arr1.(*array.Float64)
	fa2, ok2 := arr2.(*array.Float64)
	require.True(t, ok1 && ok2)

	for i := range numRows {
		assert.Equal(t, fa1.Value(i), fa2.Value(i),
			"RAND(99) row %d must be identical across same-seed instances", i)
	}
}

func TestRandFunc_Eval_ReusableBuilder(t *testing.T) {
	// Calling Eval multiple times on the same instance must not accumulate state.
	const numRows = 3
	rec := makeRecord(numRows)
	defer rec.Release()

	fn := newRandFunc(nil, []Expression{seedArg(7)})

	arr1, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr1.Release()

	arr2, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr2.Release()

	assert.Equal(t, numRows, arr2.Len(), "second Eval call must still return correct row count")

	// The second call should produce a different sequence (PRNG state advances).
	fa2, ok := arr2.(*array.Float64)
	require.True(t, ok)
	for i := range numRows {
		v := fa2.Value(i)
		assert.True(t, v >= 0 && v < 1, "row %d: rand() must be in [0, 1), got %v", i, v)
	}
}

func TestRandFunc_Eval_ZeroRows(t *testing.T) {
	rec := makeRecord(0)
	defer rec.Release()

	fn := newRandFunc(nil, nil)
	arr, err := fn.Eval(rec)
	require.NoError(t, err)
	defer arr.Release()

	assert.Equal(t, 0, arr.Len(), "Eval on zero-row batch must return empty array")
}

func TestRandFunc_EvalScalar_SeededSequenceAdvances(t *testing.T) {
	// Repeated calls on the same seeded instance must produce distinct values (PRNG advances).
	fn := newRandFunc(nil, []Expression{seedArg(42)})
	v1, err := fn.EvalScalar()
	require.NoError(t, err)
	v2, err := fn.EvalScalar()
	require.NoError(t, err)
	assert.NotEqual(t, scalar.ToFloat64(v1), scalar.ToFloat64(v2),
		"consecutive EvalScalar calls must advance the PRNG and return different values")
}

func TestRandFunc_NewFunc_FloatSeed(t *testing.T) {
	// MySQL truncates float seed to integer; LinDB must not panic on float seed.
	floatSeedArg := NewConstant(nil, float64(3.14), Scalar)
	fn := newRandFunc(nil, []Expression{floatSeedArg})

	v, err := fn.EvalScalar()
	require.NoError(t, err)
	f := scalar.ToFloat64(v)
	assert.True(t, f >= 0 && f < 1, "rand(3.14) must be in [0, 1), got %v", f)
}

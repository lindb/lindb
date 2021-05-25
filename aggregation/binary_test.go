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

package aggregation

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/sql/stmt"
)

func TestBinary_eval(t *testing.T) {
	assert.Equal(t, float64(10), eval(stmt.ADD, 4, 6))
	assert.Equal(t, float64(-2), eval(stmt.SUB, 4, 6))
	assert.Equal(t, float64(24), eval(stmt.MUL, 4, 6))
	assert.Equal(t, 0.5, eval(stmt.DIV, 4, 8))
	assert.Equal(t, float64(0), eval(stmt.DIV, 4, 0))

	// wrong binary operator
	assert.Equal(t, float64(0), eval(stmt.OR, 4, 8))
}

func TestBinary_Eval_Single(t *testing.T) {
	left := collections.NewFloatArray(10)
	left.SetValue(0, 1.1)
	left.SetValue(5, 5.5)
	left.SetValue(8, 9.9)

	right := collections.NewFloatArray(10)
	right.SetSingle(true)
	for i := 0; i < 10; i++ {
		right.SetValue(i, 10)
	}

	result := binaryEval(stmt.MUL, left, right)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 11.0, result.GetValue(0))
	assert.Equal(t, 55.0, result.GetValue(5))
	assert.Equal(t, 99.0, result.GetValue(8))

	result = binaryEval(stmt.MUL, right, left)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 11.0, result.GetValue(0))
	assert.Equal(t, 55.0, result.GetValue(5))
	assert.Equal(t, 99.0, result.GetValue(8))
}

func TestBinaryEval(t *testing.T) {
	assert.Nil(t, binaryEval(stmt.DIV, collections.NewFloatArray(10), collections.NewFloatArray(10)))

	fa := collections.NewFloatArray(10)
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa.SetValue(8, 9.9)
	result := binaryEval(stmt.ADD, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, fa, result)
	result = binaryEval(stmt.ADD, collections.NewFloatArray(10), nil)
	assert.Nil(t, result)
	result = binaryEval(stmt.ADD, nil, collections.NewFloatArray(10))
	assert.Nil(t, result)

	result = binaryEval(stmt.SUB, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, -1.1, result.GetValue(0))
	assert.Equal(t, -5.5, result.GetValue(5))
	assert.Equal(t, -9.9, result.GetValue(8))
	result = binaryEval(stmt.MUL, collections.NewFloatArray(10), fa)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.MUL, fa, collections.NewFloatArray(10))
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.SUB, fa, collections.NewFloatArray(10))
	assert.Equal(t, fa, result)
	result = binaryEval(stmt.ADD, fa, collections.NewFloatArray(10))
	assert.Equal(t, fa, result)

	fa = collections.NewFloatArray(10)
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa2 := collections.NewFloatArray(10)
	fa2.SetValue(0, 1.1)
	fa2.SetValue(8, 9.9)
	result = binaryEval(stmt.ADD, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 2.2, result.GetValue(0))
	assert.Equal(t, 5.5, result.GetValue(5))
	assert.Equal(t, 9.9, result.GetValue(8))
	result = binaryEval(stmt.SUB, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 0.0, result.GetValue(0))
	assert.Equal(t, 5.5, result.GetValue(5))
	assert.Equal(t, -9.9, result.GetValue(8))
	result = binaryEval(stmt.MUL, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 1.21, math.Floor(result.GetValue(0)*100)/100)
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
	result = binaryEval(stmt.DIV, fa, fa2)
	assert.Equal(t, 3, result.Size())
	assert.Equal(t, 1.0, result.GetValue(0))
	assert.Equal(t, 0.0, result.GetValue(5))
	assert.Equal(t, 0.0, result.GetValue(8))
}

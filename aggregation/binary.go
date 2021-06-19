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
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/sql/stmt"
)

// binaryEval evaluates two float array and returns float array
// 1. capacity not equals, return nil
// 2. merge two array based on binary operator, return other float array
// NOTICE: make sure both left and right array are not nil and same capacity
func binaryEval(binaryOp stmt.BinaryOP, left, right collections.FloatArray) collections.FloatArray {
	if left == nil || right == nil {
		return nil
	}
	if left.IsEmpty() && right.IsEmpty() {
		return nil
	}

	capacity := left.Capacity()
	result := collections.NewFloatArray(capacity)

	for i := 0; i < capacity; i++ {
		leftHasValue := left.HasValue(i)
		rightHasValue := right.HasValue(i)
		switch {
		case !leftHasValue && right.IsSingle():
		case left.IsSingle() && !rightHasValue:
		case leftHasValue || rightHasValue:
			result.SetValue(i, eval(binaryOp, left.GetValue(i), right.GetValue(i)))
		}
	}

	return result
}

// eval evaluates two values and returns another value
func eval(binaryOp stmt.BinaryOP, left, right float64) float64 {
	switch binaryOp {
	case stmt.ADD:
		return left + right
	case stmt.SUB:
		return left - right
	case stmt.MUL:
		return left * right
	case stmt.DIV:
		if right == 0 {
			return 0
		}
		return left / right
	default:
		return 0
	}
}

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

package tree

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/assert"
)

func TestHistogramFunctions_IsFuncSupported(t *testing.T) {
	cases := []FuncName{
		HistogramQuantile,
		HistogramAvg,
		HistogramSum,
		HistogramCount,
	}
	for _, fn := range cases {
		assert.True(t, IsFuncSupported(fn), "expected %s to be supported", fn)
	}
}

func TestHistogramFunctions_IsAggFunc(t *testing.T) {
	cases := []FuncName{
		HistogramQuantile,
		HistogramAvg,
		HistogramSum,
		HistogramCount,
	}
	for _, fn := range cases {
		assert.True(t, IsAggFunc(fn), "expected %s to be an aggregation function", fn)
	}
}

func TestHistogramFunctions_Constants(t *testing.T) {
	assert.Equal(t, FuncName("histogram_quantile"), HistogramQuantile)
	assert.Equal(t, FuncName("histogram_avg"), HistogramAvg)
	assert.Equal(t, FuncName("histogram_sum"), HistogramSum)
	assert.Equal(t, FuncName("histogram_count"), HistogramCount)
}

func TestRandFunction_IsFuncSupported(t *testing.T) {
	assert.True(t, IsFuncSupported(Rand), "rand must be registered as a supported function")
}

func TestRandFunction_IsNotAggFunc(t *testing.T) {
	assert.False(t, IsAggFunc(Rand), "rand must not be treated as an aggregation function")
}

func TestRandFunction_ReturnType(t *testing.T) {
	retType := GetDefaultFuncReturnType(Rand)
	assert.Equal(t, arrow.PrimitiveTypes.Float64, retType, "rand must return Float64")
}

func TestRandFunction_Constant(t *testing.T) {
	assert.Equal(t, FuncName("rand"), Rand)
}

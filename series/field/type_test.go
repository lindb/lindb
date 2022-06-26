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

package field

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
)

func TestDownSamplingFunc(t *testing.T) {
	assert.Equal(t, function.Sum, SumField.DownSamplingFunc())
	assert.Equal(t, function.Sum, HistogramField.DownSamplingFunc())
	assert.Equal(t, function.Min, MinField.DownSamplingFunc())
	assert.Equal(t, function.Max, MaxField.DownSamplingFunc())
	assert.Equal(t, function.Last, LastField.DownSamplingFunc())
	assert.Equal(t, function.First, FirstField.DownSamplingFunc())
	assert.Equal(t, function.Unknown, Unknown.DownSamplingFunc())
}

func TestType_String(t *testing.T) {
	assert.Equal(t, "sum", SumField.String())
	assert.Equal(t, "max", MaxField.String())
	assert.Equal(t, "min", MinField.String())
	assert.Equal(t, "last", LastField.String())
	assert.Equal(t, "first", FirstField.String())
	assert.Equal(t, "histogram", HistogramField.String())
	assert.Equal(t, "unknown", Unknown.String())
	assert.Equal(t, "name", Name("name").String())
}

func TestIsSupportFunc(t *testing.T) {
	assert.True(t, HistogramField.IsFuncSupported(function.Sum))
	assert.False(t, HistogramField.IsFuncSupported(function.Last))

	assert.True(t, SumField.IsFuncSupported(function.Sum))
	assert.True(t, SumField.IsFuncSupported(function.Min))
	assert.True(t, SumField.IsFuncSupported(function.Max))
	assert.False(t, SumField.IsFuncSupported(function.Quantile))

	assert.True(t, MaxField.IsFuncSupported(function.Max))
	assert.False(t, MaxField.IsFuncSupported(function.Quantile))

	assert.True(t, LastField.IsFuncSupported(function.Last))
	assert.False(t, LastField.IsFuncSupported(function.Quantile))

	assert.True(t, FirstField.IsFuncSupported(function.First))
	assert.False(t, FirstField.IsFuncSupported(function.Quantile))

	assert.True(t, MinField.IsFuncSupported(function.Min))
	assert.False(t, MinField.IsFuncSupported(function.Quantile))

	assert.False(t, Unknown.IsFuncSupported(function.Quantile))
}

func TestAggType_Aggregate(t *testing.T) {
	assert.Equal(t, 100.0, SumField.AggType().Aggregate(1, 99.0))

	assert.Equal(t, 1.0, MinField.AggType().Aggregate(1, 99.0))
	assert.Equal(t, 1.0, MinField.AggType().Aggregate(99.0, 1))

	assert.Equal(t, 99.0, MaxField.AggType().Aggregate(1, 99.0))
	assert.Equal(t, 99.0, MaxField.AggType().Aggregate(99.0, 1))

	assert.Equal(t, 99.0, LastField.AggType().Aggregate(1, 99.0))

	assert.Equal(t, 1.0, FirstField.AggType().Aggregate(1, 99.0))

	assert.Panics(t, func() {
		AggType(22).Aggregate(1, 2)
	})
}

func TestPanicAgg(t *testing.T) {
	assert.Panics(t, func() {
		Type(99).AggType().Aggregate(1, 99.0)
	})
}

func TestType_GetFuncFieldParams(t *testing.T) {
	assert.Empty(t, Type(99).GetFuncFieldParams(function.Min))
	assert.Equal(t, []AggType{Sum}, HistogramField.GetFuncFieldParams(function.Min))

	assert.Equal(t, []AggType{Max}, MaxField.GetFuncFieldParams(function.Max))
	assert.Equal(t, []AggType{Min}, MaxField.GetFuncFieldParams(function.Min))

	assert.Equal(t, []AggType{Max}, MinField.GetFuncFieldParams(function.Max))
	assert.Equal(t, []AggType{Min}, MinField.GetFuncFieldParams(function.Min))

	assert.Equal(t, []AggType{Sum}, SumField.GetFuncFieldParams(function.Sum))
	assert.Equal(t, []AggType{Max}, SumField.GetFuncFieldParams(function.Max))
	assert.Equal(t, []AggType{Min}, SumField.GetFuncFieldParams(function.Min))

	assert.Equal(t, []AggType{Sum}, LastField.GetFuncFieldParams(function.Sum))
	assert.Equal(t, []AggType{Max}, LastField.GetFuncFieldParams(function.Max))
	assert.Equal(t, []AggType{Min}, LastField.GetFuncFieldParams(function.Min))
	assert.Equal(t, []AggType{Last}, LastField.GetFuncFieldParams(function.Last))

	assert.Equal(t, []AggType{Sum}, FirstField.GetFuncFieldParams(function.Sum))
	assert.Equal(t, []AggType{Max}, FirstField.GetFuncFieldParams(function.Max))
	assert.Equal(t, []AggType{Min}, FirstField.GetFuncFieldParams(function.Min))
	assert.Equal(t, []AggType{First}, FirstField.GetFuncFieldParams(function.First))
}

func TestType_GetDefaultFuncFieldParams(t *testing.T) {
	assert.Empty(t, Type(99).GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{Sum}, HistogramField.GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{Sum}, SumField.GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{Max}, MaxField.GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{Min}, MinField.GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{Last}, LastField.GetDefaultFuncFieldParams())
	assert.Equal(t, []AggType{First}, FirstField.GetDefaultFuncFieldParams())
}

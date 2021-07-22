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

package function

import (
	"math"
	"testing"

	"github.com/lindb/lindb/pkg/collections"

	"github.com/stretchr/testify/assert"
)

func makeFloatArray(data []float64) []*collections.FloatArray {
	array := collections.NewFloatArray(len(data))
	for idx := range data {
		array.SetValue(idx, data[idx])
	}
	return []*collections.FloatArray{array}
}

func getDataFloatArray(array *collections.FloatArray) []float64 {
	var data []float64
	itr := array.NewIterator()
	for itr.HasNext() {
		_, v := itr.Next()
		data = append(data, v)
	}
	return data
}

func Test_Buckets(t *testing.T) {
	var bkt0 = buckets{}
	bkt0.EnsureCountFieldCumulative()

	var bkt = buckets{
		{upperBound: 0, count: 1},
		{upperBound: 2, count: 2},
		{upperBound: 4, count: 19}}
	bkt.EnsureCountFieldCumulative()
	assert.Equal(t, buckets{
		{upperBound: 0, count: 1},
		{upperBound: 2, count: 3},
		{upperBound: 4, count: 22}}, bkt)
}

func Test_QuantileCall(t *testing.T) {
	fields := map[float64][]*collections.FloatArray{
		1:               makeFloatArray([]float64{1, 2, 3, 4}),
		2:               makeFloatArray([]float64{2, 2, 3, 4}),
		4:               makeFloatArray([]float64{5, 2, 3, 4}),
		8:               makeFloatArray([]float64{4, 2, 3, 4}),
		20:              makeFloatArray([]float64{3, 2, 3, 4}),
		50:              makeFloatArray([]float64{1, 2, 3, 4}),
		math.Inf(1) + 1: makeFloatArray([]float64{1, 2, 3, 4}),
	}

	// invalid fields
	_, err := QuantileCall(-1, fields)
	assert.NotNil(t, err)
	_, err = QuantileCall(1.01, fields)
	assert.NotNil(t, err)

	// quantile(0) returns 0
	array, _ := QuantileCall(0, fields)
	assert.Equal(t, []float64{0, 0, 0, 0}, getDataFloatArray(array))

	// quantile(1) returns upperBound before inf
	array, _ = QuantileCall(1, fields)
	assert.Equal(t, []float64{50, 50, 50, 50}, getDataFloatArray(array))

	array, _ = QuantileCall(0.99, fields)
	assert.Equal(t, []float64{50, 50, 50, 50}, getDataFloatArray(array))

	array, _ = QuantileCall(0.9, fields)
	assert.InDeltaSlice(t, []float64{29, 50, 50, 50}, getDataFloatArray(array), 0.001)

	array, _ = QuantileCall(0.8, fields)
	assert.InDeltaSlice(t, []float64{14.4, 38, 38, 38}, getDataFloatArray(array), 0.001)

	array, _ = QuantileCall(0.5, fields)
	assert.InDeltaSlice(t, []float64{4.5, 6, 6, 6}, getDataFloatArray(array), 0.001)

}

func Test_QuantileCallBadCases(t *testing.T) {
	// starts from upperBound < 0
	fields := map[float64][]*collections.FloatArray{
		-0.2:            makeFloatArray([]float64{1, 2, 3, 4}),
		2:               makeFloatArray([]float64{2, 2, 3, 4}),
		4:               makeFloatArray([]float64{5, 2, 3, 4}),
		8:               makeFloatArray([]float64{4, 2, 3, 4}),
		20:              makeFloatArray([]float64{3, 2, 3, 4}),
		50:              makeFloatArray([]float64{1, 2, 3, 4}),
		math.Inf(1) + 1: makeFloatArray([]float64{1, 2, 3, 4}),
	}
	array, _ := QuantileCall(0, fields)
	assert.InDeltaSlice(t, []float64{-0.2, -0.2, -0.2, -0.2}, getDataFloatArray(array), 0.001)

	// fields buckets less than 2
	fields = map[float64][]*collections.FloatArray{
		-0.2: makeFloatArray([]float64{1, 2, 3, 4}),
	}
	_, err := QuantileCall(0, fields)
	assert.Error(t, err)

	// last upper bound not Inf
	fields = map[float64][]*collections.FloatArray{
		2: makeFloatArray([]float64{2, 2, 3, 4}),
		4: makeFloatArray([]float64{5, 2, 3, 4}),
	}
	_, err = QuantileCall(0, fields)
	assert.Error(t, err)

	// all data zero
	fields = map[float64][]*collections.FloatArray{
		1:               makeFloatArray([]float64{0, 0, 0, 0}),
		2:               makeFloatArray([]float64{0, 0, 0, 0}),
		4:               makeFloatArray([]float64{0, 0, 0, 0}),
		8:               makeFloatArray([]float64{0, 0, 0, 0}),
		20:              makeFloatArray([]float64{0, 0, 0, 0}),
		50:              makeFloatArray([]float64{0, 0, 0, 0}),
		math.Inf(1) + 1: makeFloatArray([]float64{0, 0, 0, 0}),
	}
	array, err = QuantileCall(0, fields)
	assert.Nil(t, err)
	assert.InDeltaSlice(t, []float64{0, 0, 0, 0}, getDataFloatArray(array), 0.001)

	// float64 array length not ok
	fields = map[float64][]*collections.FloatArray{
		2:               makeFloatArray([]float64{2, 2, 3, 4}),
		4:               append(makeFloatArray([]float64{5, 2, 3, 4}), makeFloatArray([]float64{5, 2, 3, 4})...),
		math.Inf(1) + 1: makeFloatArray([]float64{1, 2, 3, 4}),
	}
	_, err = QuantileCall(0.9, fields)
	assert.Error(t, err)
	// data length not match
	fields = map[float64][]*collections.FloatArray{
		1:               makeFloatArray([]float64{0, 0, 0, 0, 1, 2, 3}),
		2:               makeFloatArray([]float64{0, 0, 0, 0}),
		4:               makeFloatArray([]float64{0, 0, 0, 0}),
		8:               makeFloatArray([]float64{0, 0, 0, 0}),
		20:              makeFloatArray([]float64{0, 0, 0, 0, 1}),
		50:              makeFloatArray([]float64{0, 0, 0, 0, 2}),
		math.Inf(1) + 1: makeFloatArray([]float64{0, 0, 0, 0}),
	}
	_, err = QuantileCall(0.9, fields)
	assert.Error(t, err)
}

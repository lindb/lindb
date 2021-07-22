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

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatArray(t *testing.T) {
	fa := NewFloatArray(10)
	assert.Equal(t, 10, fa.Capacity())
	assert.Equal(t, 0, fa.Size())
	assert.True(t, fa.IsEmpty())
	assert.False(t, fa.IsSingle())
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa.SetValue(8, 9.9)
	fa.SetValue(-1, 1.1)
	fa.SetValue(10, 11.1)
	fa.SetValue(11, 11.1)
	assert.False(t, fa.IsEmpty())
	assert.True(t, fa.HasValue(0))
	assert.True(t, fa.HasValue(5))
	assert.False(t, fa.HasValue(-1))
	assert.False(t, fa.HasValue(10))
	assert.False(t, fa.HasValue(11))

	assert.Equal(t, float64(0), fa.GetValue(-1))
	assert.Equal(t, 1.1, fa.GetValue(0))
	assert.Equal(t, 5.5, fa.GetValue(5))
	assert.Equal(t, 9.9, fa.GetValue(8))
	assert.Equal(t, float64(0), fa.GetValue(10))
	assert.Equal(t, float64(0), fa.GetValue(11))

	assert.Equal(t, 3, fa.Size())

	for i := 0; i < 3; i++ {
		it := fa.NewIterator()
		assert.True(t, it.HasNext())
		idx, value := it.Next()
		assert.Equal(t, 0, idx)
		assert.Equal(t, 1.1, value)
		assert.True(t, it.HasNext())
		idx, value = it.Next()
		assert.Equal(t, 5, idx)
		assert.Equal(t, 5.5, value)
		assert.True(t, it.HasNext())
		idx, value = it.Next()
		assert.Equal(t, 8, idx)
		assert.Equal(t, 9.9, value)
		assert.False(t, it.HasNext())
		idx, value = it.Next()
		assert.Equal(t, -1, idx)
		assert.Equal(t, float64(0), value)
	}

	// reset
	fa.SetValue(8, 10.10)
	assert.Equal(t, 10.10, fa.GetValue(8))
	assert.Equal(t, 3, fa.Size())
}

func TestFloatArray_Single(t *testing.T) {
	fa := NewFloatArray(10)
	assert.False(t, fa.IsSingle())
	for i := 0; i < 10; i++ {
		fa.SetValue(i, 10)
	}
	assert.Equal(t, 10, fa.Size())
	fa.SetSingle(true)
	assert.True(t, fa.IsSingle())
	fa.Reset()
	assert.False(t, fa.IsSingle())
	assert.True(t, fa.IsEmpty())
}

func BenchmarkFloatArray_HasValue(b *testing.B) {
	for x := 0; x < b.N; x++ {
		pos := 1000
		_ = pos / blockSize
		_ = pos - pos/blockSize*blockSize
	}
}

func BenchmarkFloatArray_HasValue_2(b *testing.B) {
	for x := 0; x < b.N; x++ {
		pos := 1000
		_ = pos / blockSize
		_ = pos % blockSize
	}
}

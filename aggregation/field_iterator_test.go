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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestFieldIterator(t *testing.T) {
	it := newFieldIterator(20, []field.AggType{field.Sum}, []*collections.FloatArray{generateFloatArray(nil)})
	assert.True(t, it.HasNext())
	assert.NotNil(t, it.Next())
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	it = newFieldIterator(20, []field.AggType{field.Min}, []*collections.FloatArray{generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0})})

	expect := map[int]float64{20: 0, 21: 10, 22: 10.0, 23: 100.4, 24: 50.0}
	AssertFieldIt(t, it, expect)
	assert.False(t, it.HasNext())
	assert.Nil(t, it.Next())

	// marshal has data, reset idx
	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// test empty data
	it = newFieldIterator(20, nil, nil)
	assert.False(t, it.HasNext())
	assert.Nil(t, it.Next())

	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestFieldIterator_MarshalBinary(t *testing.T) {
	defer func() {
		toBytesFn = toBytes
	}()
	pData := generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0})
	it := newFieldIterator(10, []field.AggType{field.Sum}, []*collections.FloatArray{pData})
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	fIt := series.NewFieldIterator(data)
	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	AssertFieldIt(t, fIt, expect)
	assert.False(t, fIt.HasNext())

	floatArray := collections.NewFloatArray(4)
	floatArray.SetValue(3, float64(3))
	it = newFieldIterator(5, []field.AggType{field.Sum}, []*collections.FloatArray{floatArray})
	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	fIt = series.NewFieldIterator(data)
	expect = map[int]float64{8: 3.0}
	AssertFieldIt(t, fIt, expect)
	assert.False(t, fIt.HasNext())

	it = newFieldIterator(10, []field.AggType{field.Sum, field.Sum}, []*collections.FloatArray{pData, pData})
	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	toBytesFn = func(e *encoding.TSDEncoder) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	it = newFieldIterator(10, []field.AggType{field.Sum, field.Sum}, []*collections.FloatArray{pData, pData})
	data, err = it.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestPrimitiveIterator(t *testing.T) {
	it := newPrimitiveIterator(1, field.Sum, nil)
	assert.False(t, it.HasNext())
	slot, v := it.Next()
	assert.Equal(t, -1, slot)
	assert.Equal(t, float64(0), v)

	floatArray := collections.NewFloatArray(4)
	floatArray.SetValue(3, float64(3))
	it = newPrimitiveIterator(1, field.Sum, floatArray)
	slot, v = it.Next()
	assert.Equal(t, -1, slot)
	assert.Equal(t, float64(0), v)
}

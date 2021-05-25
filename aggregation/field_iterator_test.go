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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

var encodeFunc = encoding.NewTSDEncoder

func TestFieldIterator(t *testing.T) {
	it := newFieldIterator(20, field.Sum, generateFloatArray(nil))
	assert.False(t, it.HasNext())
	slot, value := it.Next()
	assert.Equal(t, -1, slot)
	assert.Equal(t, 0.0, value)
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)

	it = newFieldIterator(20, field.Min, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	expect := map[int]float64{20: 0, 21: 10, 22: 10.0, 23: 100.4, 24: 50.0}
	AssertFieldIt(t, it, expect)
	assert.False(t, it.HasNext())
	slot, value = it.Next()
	assert.Equal(t, -1, slot)
	assert.Equal(t, 0.0, value)

	// marshal empty, because field iterator already read
	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestFieldIterator_MarshalBinary(t *testing.T) {
	it := newFieldIterator(10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	reader := stream.NewReader(data)
	aggType := field.AggType(reader.ReadByte()) // read field agg type
	assert.Equal(t, field.Sum, aggType)
	length := reader.ReadVarint32()
	data1 := reader.ReadBytes(int(length))

	fIt := series.NewFieldIterator(aggType, encoding.NewTSDDecoder(data1))
	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	AssertFieldIt(t, fIt, expect)
	assert.False(t, fIt.HasNext())
}

func TestFieldIterator_MarshalBinary_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()
	encoder := encoding.NewMockTSDEncoder(ctrl)
	encoding.TSDEncodeFunc = func(startTime uint16) encoding.TSDEncoder {
		return encoder
	}
	floatArray := collections.NewFloatArray(5)
	floatArray.SetValue(2, 10.0)
	encoder.EXPECT().AppendTime(gomock.Any()).AnyTimes()
	encoder.EXPECT().AppendValue(gomock.Any()).AnyTimes()
	encoder.EXPECT().Bytes().Return(nil, fmt.Errorf("err"))
	it := newFieldIterator(10, field.Sum, floatArray)
	data, err := it.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)
}

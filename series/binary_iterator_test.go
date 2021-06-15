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

package series

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

func TestBinaryGroupedIterator(t *testing.T) {
	writer := stream.NewBufferWriter(nil)
	d := buildFieldIterator()
	writer.PutByte(byte(field.SumField))
	writer.PutVarint64(10)
	writer.PutVarint32(int32(len(d)))
	writer.PutBytes(d)
	data, err := writer.Bytes()
	assert.NoError(t, err)
	result := make(map[field.Name]field.Name)
	it := NewGroupedIterator("1.1.1.1", map[field.Name][]byte{
		"f1": data,
		"f2": data,
	})
	assert.Equal(t, "1.1.1.1", it.Tags())
	assert.True(t, it.HasNext())
	fIt := it.Next()
	assert.Equal(t, field.SumField, fIt.FieldType())
	assert.True(t, fIt.HasNext())
	result[fIt.FieldName()] = fIt.FieldName()
	startTime, fIt1 := fIt.Next()
	assert.Equal(t, int64(10), startTime)
	assertFieldIterator(t, fIt1)

	assert.True(t, it.HasNext())
	fIt = it.Next()
	assert.Equal(t, field.SumField, fIt.FieldType())
	assert.True(t, fIt.HasNext())
	result[fIt.FieldName()] = fIt.FieldName()
	startTime, fIt1 = fIt.Next()
	assert.Equal(t, int64(10), startTime)
	assertFieldIterator(t, fIt1)

	assert.False(t, it.HasNext())

	assert.Equal(t, 2, len(result))
}

func TestBinaryIterator(t *testing.T) {
	writer := stream.NewBufferWriter(nil)
	writer.PutByte(byte(field.SumField))
	d := buildFieldIterator()
	writer.PutVarint64(10)
	writer.PutVarint32(int32(len(d)))
	writer.PutBytes(d)
	d = buildFieldIterator()
	writer.PutVarint64(11)
	writer.PutVarint32(int32(len(d)))
	writer.PutBytes(d)
	data, err := writer.Bytes()
	assert.NoError(t, err)
	it := NewIterator("f1", data)
	assert.Equal(t, field.Name("f1"), it.FieldName())
	assert.Equal(t, field.SumField, it.fieldType)
	assert.True(t, it.HasNext())
	startTime, fIt := it.Next()
	assert.Equal(t, int64(10), startTime)
	assertFieldIterator(t, fIt)
	assert.True(t, it.HasNext())
	startTime, fIt = it.Next()
	assert.Equal(t, int64(11), startTime)
	assertFieldIterator(t, fIt)
	assert.False(t, it.HasNext())

	// test marshal binary
	it = NewIterator("f1", data)
	data2, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, data, data2)

	writer = stream.NewBufferWriter(nil)
	writer.PutByte(byte(field.SumField))
	writer.PutVarint64(10)
	writer.PutVarint32(int32(0))
	data, _ = writer.Bytes()
	it = NewIterator("f1", data)
	assert.True(t, it.HasNext())
	startTime, fIt = it.Next()
	assert.Equal(t, int64(10), startTime)
	assert.Nil(t, fIt)
	assert.False(t, it.HasNext())
}

func TestBinaryFieldIterator(t *testing.T) {
	d := buildFieldIterator()
	it := NewFieldIterator(d)
	assertFieldIterator(t, it)

	_, err := it.MarshalBinary()
	assert.Error(t, err)
}

func assertFieldIterator(t *testing.T, it FieldIterator) {
	assert.True(t, it.HasNext())
	pIt := it.Next()
	assert.Equal(t, field.Sum, pIt.AggType())
	assert.True(t, pIt.HasNext())
	s, v := pIt.Next()
	assert.Equal(t, 12, s)
	assert.Equal(t, 10.0, v)
	assert.False(t, it.HasNext())
}

func buildFieldIterator() []byte {
	writer := stream.NewBufferWriter(nil)
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.Bytes()
	writer.PutByte(byte(field.Sum))
	writer.PutVarint32(int32(len(data)))
	writer.PutBytes(data)
	d, _ := writer.Bytes()
	return d
}

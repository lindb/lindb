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

package encoding

import (
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedOffsetEncoder_IsEmpty(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	assert.True(t, encoder.IsEmpty())
	encoder.Add(10)
	assert.False(t, encoder.IsEmpty())
}

func TestFixedOffsetDecoder_Get(t *testing.T) {
	// test all empty value
	encoder := NewFixedOffsetEncoder()
	for i := 1; i < 10000; i++ {
		encoder.Add(i * -1)
	}
	data := encoder.MarshalBinary()
	decoder := NewFixedOffsetDecoder(data)
	assert.Equal(t, 1, decoder.width)
	for i := 1; i < 10000; i++ {
		_, ok := decoder.Get(i)
		assert.False(t, ok)
	}

	// test with empty offset
	encoder.Reset()
	encoder.Add(10)
	encoder.Add(-10)
	encoder.Add(20)
	encoder.Add(-10)
	data = encoder.MarshalBinary()
	decoder = NewFixedOffsetDecoder(data)
	assert.Equal(t, 1, decoder.width)
	v, ok := decoder.Get(0)
	assert.True(t, ok)
	assert.Equal(t, 10, v)
	_, ok = decoder.Get(1)
	assert.False(t, ok)
	v, ok = decoder.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 20, v)
	_, ok = decoder.Get(3)
	assert.False(t, ok)
}

func TestFixedOffsetDecoder_Codec(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	data := encoder.MarshalBinary()
	assert.Len(t, data, 0)
	assert.Equal(t, 0, encoder.Size())

	encoder.FromValues([]int{
		0,
		1,
		1 << 8,
		(1 << 8) + 1,
		1 << 16,
		(1 << 16) + 1,
		1 << 24,
		1<<24 + 1,
	})
	assert.Equal(t, 8, encoder.Size())
	data = encoder.MarshalBinary()
	assert.True(t, len(data) > 1)
	decoder := NewFixedOffsetDecoder(data)
	assert.Equal(t, 4, decoder.width)

	var ok bool
	var value int
	value, _ = decoder.Get(7)
	assert.Equal(t, (1<<24)+1, value)
	value, _ = decoder.Get(5)
	assert.Equal(t, (1<<16)+1, value)
	value, _ = decoder.Get(3)
	assert.Equal(t, (1<<8)+1, value)
	value, _ = decoder.Get(2)
	assert.Equal(t, 1<<8, value)
	value, _ = decoder.Get(0)
	assert.Equal(t, 0, value)
	_, ok = decoder.Get(8)
	assert.False(t, ok)
	_, ok = decoder.Get(-8)
	assert.False(t, ok)
	assert.Equal(t, 8, decoder.Size())
	assert.Equal(t, 4, decoder.ValueWidth())

	// test empty
	decoder = NewFixedOffsetDecoder([]byte{0})
	assert.Equal(t, 0, decoder.width)
	assert.Zero(t, decoder.Size())
	assert.Equal(t, lengthOfHeader, decoder.Header())
}

type mockWriter struct{}

func (w *mockWriter) Write(p []byte) (int, error) {
	return 0, io.ErrShortBuffer
}

func TestFixedOffsetEncoder_WriteTo(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	encoder.FromValues([]int{1, 2, 3})
	assert.NotNil(t, encoder.WriteTo(&mockWriter{}))
}

func TestFixedOffsetEncoder_Reset(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	encoder.Add(0) //0
	encoder.Add(1) //1
	data := encoder.MarshalBinary()
	assert.True(t, len(data) > 1)
	// reset
	encoder.Reset()

	encoder.Add(1 << 24)       //0
	encoder.Add((1 << 24) + 1) //1
	data = encoder.MarshalBinary()
	assert.True(t, len(data) > 1)

	decoder := NewFixedOffsetDecoder(data)
	assert.Equal(t, 4, decoder.width)

	value, _ := decoder.Get(1)
	assert.Equal(t, (1<<24)+1, value)
	value, _ = decoder.Get(0)
	assert.Equal(t, 1<<24, value)
}

func TestFixedOffset_codec_int(t *testing.T) {
	assertGet := func(buf []byte, index int, value int, ok bool) {
		decoder := NewFixedOffsetDecoder(buf)
		v, exist := decoder.Get(index)
		assert.Equal(t, value, v)
		assert.Equal(t, ok, exist)
	}

	assertGet(nil, 0, 0, false)
	// width:1
	assertGet([]byte{1}, 0, 0, false)
	assertGet([]byte{1, 0xff}, -1, 0, false)
	assertGet([]byte{1, 0xff}, 0, 254, true)
	assertGet([]byte{1, 0xff}, 1, 0, false)
	// width: 4
	assertGet([]byte{2, 0xff}, 0, 0, false)
	// width: 4
	assertGet([]byte{4, 0xff}, 0, 0, false)
	assertGet([]byte{4, 0xff, 0xff, 0xff, 0xff}, 0, 0xfffffffe, true)

	// data corruption
	assertGet([]byte{5, 0xff, 0xff, 0xff, 0xff}, 0, 0, false)

}

func TestByteSlice2Uint32(t *testing.T) {
	assert.Equal(t, uint32(1), ByteSlice2Uint32([]byte{1}))
	assert.Equal(t, uint32(0xffffffff), ByteSlice2Uint32([]byte{0xff, 0xff, 0xff, 0xff}))
	assert.Equal(t, uint32(0xffffffff), ByteSlice2Uint32([]byte{0xff, 0xff, 0xff, 0xff, 0xff}))
	assert.Equal(t, uint32(0), ByteSlice2Uint32(nil))
}

func Test_GetAdd_Consistency(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var expects []int
	for i := 0; i < 100000; i++ {
		expects = append(expects, rand.Intn(100000000))
	}
	encoder := NewFixedOffsetEncoder()
	encoder.FromValues(expects)
	decoder := NewFixedOffsetDecoder(encoder.MarshalBinary())

	for i := 0; i < 100000; i++ {
		value, ok := decoder.Get(i)
		assert.Equal(t, expects[i], value)
		assert.True(t, ok)
	}
}

func BenchmarkFixedOffsetDecoder_Get(b *testing.B) {
	encoder := NewFixedOffsetEncoder()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++ {
		x := rand.Intn(10000000)
		encoder.Add(x)
	}
	data := encoder.MarshalBinary()
	decoder := NewFixedOffsetDecoder(data)
	b.ResetTimer()
	for round := 0; round < b.N; round++ {
		for i := 0; i < 100000; i++ {
			_, _ = decoder.Get(i)
		}
	}
}

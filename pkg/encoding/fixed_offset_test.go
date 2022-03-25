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
	"fmt"
	"io"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedOffsetEncoder_IsEmpty(t *testing.T) {
	encoder := NewFixedOffsetEncoder(true)
	assert.True(t, encoder.IsEmpty())
	encoder.Add(10)
	assert.False(t, encoder.IsEmpty())
}

func TestFixedOffsetDecoder_NotIncreasing(t *testing.T) {
	encoder := NewFixedOffsetEncoder(false)
	assert.Panics(t, func() {
		encoder.Add(-1)
	})
	encoder.Add(1)
	encoder.Add(3)
	assert.NotPanics(t, func() {
		encoder.Add(2)
	})
}

func TestFixedOffsetDecoder_Get(t *testing.T) {
	// test with negative value
	encoder := NewFixedOffsetEncoder(true)
	assert.Panics(t, func() {
		for i := 1; i < 10000; i++ {
			encoder.Add(i * -1)
		}
	})
	data := encoder.MarshalBinary()
	assert.Len(t, data, 0)

	// test with non-increasing value
	encoder.Add(1)
	encoder.Add(3)
	assert.Panics(t, func() {
		encoder.Add(2)
	})

	decoder := NewFixedOffsetDecoder()
	_, err := decoder.Unmarshal(nil)
	assert.Error(t, err)
	assert.Equal(t, 0, decoder.Size())

	// corrupted block
	encoder.Reset()
	encoder.Add(0)
	encoder.Add(5)
	encoder.Add(15)
	encoder.Add(30)

	data = encoder.MarshalBinary()
	fmt.Println(data)
	decoder = NewFixedOffsetDecoder()
	_, err = decoder.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, decoder.width)

	block := make([]byte, 28)
	_, err = decoder.Unmarshal(data)
	assert.NoError(t, err)
	partition, err := decoder.GetBlock(0, block)
	assert.NoError(t, err)
	assert.Len(t, partition, 5)

	partition, err = decoder.GetBlock(1, block)
	assert.NoError(t, err)
	assert.Len(t, partition, 10)

	_, err = decoder.GetBlock(2, block)
	assert.Error(t, err)
}

func TestFixedOffsetDecoder_Codec(t *testing.T) {
	encoder := NewFixedOffsetEncoder(true)
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
	decoder := NewFixedOffsetDecoder()
	_, err := decoder.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, 4, decoder.width)

	assert.Equal(t, 8, decoder.Size())
	assert.Equal(t, 4, decoder.ValueWidth())

	// test empty
	decoder = NewFixedOffsetDecoder()
	_, _ = decoder.Unmarshal([]byte{0})
	assert.Equal(t, 0, decoder.width)
	assert.Zero(t, decoder.Size())
}

type mockWriter struct {
	count   int
	errorOn int
}

func (w *mockWriter) Write(p []byte) (int, error) {
	w.count++
	if w.count >= w.errorOn {
		return 0, io.ErrShortBuffer
	}
	return len(p), nil
}

func TestFixedOffsetEncoder_WriteTo(t *testing.T) {
	encoder := NewFixedOffsetEncoder(true)
	encoder.FromValues([]int{1, 2, 3})
	assert.NotNil(t, encoder.Write(&mockWriter{errorOn: 1}))
	assert.NotNil(t, encoder.Write(&mockWriter{errorOn: 2}))
	assert.NotNil(t, encoder.Write(&mockWriter{errorOn: 3}))
}

func TestFixedOffsetEncoder_Reset(t *testing.T) {
	encoder := NewFixedOffsetEncoder(true)
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

	decoder := NewFixedOffsetDecoder()
	_, err := decoder.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, 2, decoder.Size())
	assert.Equal(t, 4, decoder.width)
}

func TestFixedOffset_codec_int(t *testing.T) {
	assertGet := func(buf []byte, index int, value int, ok bool) {
		decoder := NewFixedOffsetDecoder()
		_, _ = decoder.Unmarshal(buf)
		v, exist := decoder.Get(index)
		assert.Equal(t, value, v)
		assert.Equal(t, ok, exist)
	}

	assertGet(nil, 0, 0, false)
	// width:1
	assertGet([]byte{1}, 0, 0, false)
	assertGet([]byte{1, 1, 0xff}, -1, 0, false)
	assertGet([]byte{1, 1, 0xff}, 0, 255, true)
	assertGet([]byte{1, 1, 0xff}, 1, 0, false)
	// width: 4
	assertGet([]byte{2, 1, 0xff}, 0, 0, false)
	// width: 4
	assertGet([]byte{4, 1, 0xff}, 0, 0, false)
	assertGet([]byte{4, 1, 0xff, 0xff, 0xff, 0xff}, 0, 0xffffffff, true)

	// data corruption
	assertGet([]byte{5, 1, 0xff, 0xff, 0xff, 0xff}, 0, 0, false)
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
	sort.Ints(expects)

	encoder := NewFixedOffsetEncoder(true)
	encoder.FromValues(expects)
	decoder := NewFixedOffsetDecoder()
	_, _ = decoder.Unmarshal(encoder.MarshalBinary())

	for i := 0; i < 100000; i++ {
		value, ok := decoder.Get(i)
		assert.Equal(t, expects[i], value)
		assert.True(t, ok)
	}
}

func Test_Decoder_Corrupted_data(t *testing.T) {
	data := []byte{4, 1, 0xff, 0xff, 0xff, 0xff}
	decoder := NewFixedOffsetDecoder()
	_, _ = decoder.Unmarshal(data)

	block, err := decoder.GetBlock(1, nil)
	assert.Error(t, err)
	assert.Len(t, block, 0)

	block, err = decoder.GetBlock(0, nil)
	assert.Error(t, err)
	assert.Len(t, block, 0)
}

func TestGetFixedOffsetDecoder(t *testing.T) {
	decoder := GetFixedOffsetDecoder()
	assert.NotNil(t, decoder)
	ReleaseFixedOffsetDecoder(decoder)
}

func BenchmarkFixedOffsetDecoder_Get(b *testing.B) {
	encoder := NewFixedOffsetEncoder(true)
	var expects = make([]int, 100000)
	for i := 0; i < 100000; i++ {
		rand.Seed(time.Now().UnixNano())
		expects[i] = rand.Intn(10000000)
	}
	sort.Ints(expects)
	encoder.FromValues(expects)

	data := encoder.MarshalBinary()
	decoder := NewFixedOffsetDecoder()
	_, _ = decoder.Unmarshal(data)
	b.ResetTimer()

	for round := 0; round < b.N; round++ {
		for i := 0; i < 100000; i++ {
			_, _ = decoder.Get(i)
		}
	}
}

func Benchmark_ByteSlice2Uint32(b *testing.B) {
	var slice = []byte{1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		ByteSlice2Uint32(slice)
	}
}

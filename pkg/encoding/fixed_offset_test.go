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

func TestFixedOffsetDecoder_Codec(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	data := encoder.MarshalBinary()
	assert.Len(t, data, 0)
	assert.Equal(t, 0, encoder.Size())

	encoder.FromValues([]uint32{
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
	var value uint32
	value, _ = decoder.Get(7)
	assert.Equal(t, uint32((1<<24)+1), value)
	value, _ = decoder.Get(5)
	assert.Equal(t, uint32((1<<16)+1), value)
	value, _ = decoder.Get(3)
	assert.Equal(t, uint32((1<<8)+1), value)
	value, _ = decoder.Get(2)
	assert.Equal(t, uint32(1<<8), value)
	value, _ = decoder.Get(0)
	assert.Equal(t, uint32(0), value)
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
}

type mockWriter struct{}

func (w *mockWriter) Write(p []byte) (int, error) {
	return 0, io.ErrShortBuffer
}

func TestFixedOffsetEncoder_WriteTo(t *testing.T) {
	encoder := NewFixedOffsetEncoder()
	encoder.FromValues([]uint32{1, 2, 3})
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
	assert.Equal(t, uint32((1<<24)+1), value)
	value, _ = decoder.Get(0)
	assert.Equal(t, uint32(1<<24), value)
}

func TestFixedOffset_codec_int32(t *testing.T) {
	assert.Equal(t, 1, Uint32MinWidth(0))
	assert.Equal(t, 1, Uint32MinWidth(1))
	assert.Equal(t, 2, Uint32MinWidth(1<<8))
	assert.Equal(t, 2, Uint32MinWidth((1<<8)+1))
	assert.Equal(t, 3, Uint32MinWidth(1<<16))
	assert.Equal(t, 3, Uint32MinWidth((1<<16)+1))
	assert.Equal(t, 4, Uint32MinWidth(1<<24))
	assert.Equal(t, 4, Uint32MinWidth((1<<24)+1))

	assertGet := func(buf []byte, index int, value uint32, ok bool) {
		decoder := NewFixedOffsetDecoder(buf)
		v, exist := decoder.Get(index)
		assert.Equal(t, value, v)
		assert.Equal(t, ok, exist)
	}

	assertGet(nil, 0, 0, false)
	// width:1
	assertGet([]byte{1}, 0, 0, false)
	assertGet([]byte{1, 0xff}, -1, 0, false)
	assertGet([]byte{1, 0xff}, 0, 255, true)
	assertGet([]byte{1, 0xff}, 1, 0, false)
	// width: 4
	assertGet([]byte{2, 0xff}, 0, 0, false)
	// width: 4
	assertGet([]byte{4, 0xff}, 0, 0, false)
	assertGet([]byte{4, 0xff, 0xff, 0xff, 0xff}, 0, 0xffffffff, true)

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
	var expects []uint32
	for i := 0; i < 100000; i++ {
		expects = append(expects, rand.Uint32())
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
		x := rand.Uint32()
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

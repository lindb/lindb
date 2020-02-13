package encoding

import (
	"testing"

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
	encoder.Add(0)             //0
	encoder.Add(1)             //1
	encoder.Add(1 << 8)        //2
	encoder.Add((1 << 8) + 1)  //3
	encoder.Add(1 << 16)       //4
	encoder.Add((1 << 16) + 1) //5
	encoder.Add(1 << 24)       //6
	encoder.Add((1 << 24) + 1) //7
	data = encoder.MarshalBinary()
	assert.True(t, len(data) > 1)
	decoder := NewFixedOffsetDecoder(data)
	assert.Equal(t, 4, decoder.valueLength)

	assert.Equal(t, (1<<24)+1, decoder.Get(7))
	assert.Equal(t, (1<<16)+1, decoder.Get(5))
	assert.Equal(t, (1<<8)+1, decoder.Get(3))
	assert.Equal(t, 1<<8, decoder.Get(2))
	assert.Equal(t, 0, decoder.Get(0))
	assert.Equal(t, -1, decoder.Get(8))
	assert.Equal(t, -1, decoder.Get(-8))
	assert.Equal(t, 8, decoder.Size())
	assert.Equal(t, 4, decoder.ValueWidth())

	// test empty
	decoder = NewFixedOffsetDecoder([]byte{0})
	assert.Equal(t, 0, decoder.valueLength)
	assert.Zero(t, decoder.Size())
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
	assert.Equal(t, 4, decoder.valueLength)

	assert.Equal(t, (1<<24)+1, decoder.Get(1))
	assert.Equal(t, 1<<24, decoder.Get(0))
}

func TestFixedOffset_codec_int32(t *testing.T) {
	assert.Equal(t, 1, getValueLen(0))
	assert.Equal(t, 1, getValueLen(1))
	assert.Equal(t, 2, getValueLen(1<<8))
	assert.Equal(t, 2, getValueLen((1<<8)+1))
	assert.Equal(t, 3, getValueLen(1<<16))
	assert.Equal(t, 3, getValueLen((1<<16)+1))
	assert.Equal(t, 4, getValueLen(1<<24))
	assert.Equal(t, 4, getValueLen((1<<24)+1))

	buf := make([]byte, 1)
	putInt32(buf, 0, 1)
	assert.Equal(t, 0, getInt(buf, 0, 1))

	buf = make([]byte, 1)
	putInt32(buf, 1, 1)
	assert.Equal(t, 1, getInt(buf, 0, 1))

	buf = make([]byte, 2)
	putInt32(buf, 1<<8, 2)
	assert.Equal(t, 1<<8, getInt(buf, 0, 2))

	buf = make([]byte, 2)
	putInt32(buf, 1+(1<<8), 2)
	assert.Equal(t, 1+(1<<8), getInt(buf, 0, 2))

	buf = make([]byte, 3)
	putInt32(buf, 1<<16, 3)
	assert.Equal(t, 1<<16, getInt(buf, 0, 3))

	buf = make([]byte, 3)
	putInt32(buf, 1+(1<<16), 3)
	assert.Equal(t, 1+(1<<16), getInt(buf, 0, 3))

	buf = make([]byte, 4)
	putInt32(buf, 1<<24, 4)
	assert.Equal(t, 1<<24, getInt(buf, 0, 4))

	buf = make([]byte, 4)
	putInt32(buf, 1+(1<<24), 4)
	assert.Equal(t, 1+(1<<24), getInt(buf, 0, 4))
}

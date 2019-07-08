package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/bit"
)

func TestCodec(t *testing.T) {
	encoder := NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(100))
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(50))
	assert.Nil(t, encoder.Error())

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)
	assert.Equal(t, 10, decoder.StartTime())
	assert.Equal(t, 13, decoder.EndTime())
	assert.Equal(t, 4, decoder.count)

	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(10), decoder.Value())
	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(100), decoder.Value())
	assert.True(t, decoder.Next())
	assert.False(t, decoder.HasValue())
	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(50), decoder.Value())

	assert.False(t, decoder.Next())
}

func TestHasValueWithSlot(t *testing.T) {
	encoder := NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(100))
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(50))
	assert.Nil(t, encoder.Error())

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)

	assert.True(t, decoder.HasValueWithSlot(0))
	assert.Equal(t, uint64(10), decoder.Value())
	assert.True(t, decoder.HasValueWithSlot(1))
	assert.Equal(t, uint64(100), decoder.Value())
	assert.False(t, decoder.HasValueWithSlot(2))
	assert.True(t, decoder.HasValueWithSlot(3))
	assert.Equal(t, uint64(50), decoder.Value())

	// out of range
	assert.False(t, decoder.HasValueWithSlot(-2))
	assert.False(t, decoder.HasValueWithSlot(100))
}

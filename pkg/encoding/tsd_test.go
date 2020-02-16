package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
)

func TestCodec(t *testing.T) {
	encoder := NewTSDEncoder(10)
	data, _ := encoder.Bytes()
	assert.Len(t, data, 4)
	encoder.Reset()

	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(100))
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(50))

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)
	assert.Equal(t, uint16(10), decoder.StartTime())
	assert.Equal(t, uint16(13), decoder.EndTime())
	startTime, endTime := DecodeTSDTime(data)
	assert.Equal(t, uint16(10), startTime)
	assert.Equal(t, uint16(13), endTime)

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
	data, err = encoder.BytesWithoutTime()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder.ResetWithTimeRange(data, 10, 13)
	assert.Equal(t, uint16(10), decoder.StartTime())
	assert.Equal(t, uint16(13), decoder.EndTime())

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

	encoder.Reset()
	data, _ = encoder.Bytes()
	assert.Len(t, data, 4)
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

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)

	assert.True(t, decoder.HasValueWithSlot(10))
	assert.Equal(t, uint64(10), decoder.Value())
	assert.True(t, decoder.HasValueWithSlot(11))
	assert.Equal(t, uint64(100), decoder.Value())
	assert.False(t, decoder.HasValueWithSlot(12))
	assert.True(t, decoder.HasValueWithSlot(13))
	assert.Equal(t, uint64(50), decoder.Value())
	// out of range
	assert.False(t, decoder.HasValueWithSlot(9))
	assert.False(t, decoder.HasValueWithSlot(100))

	decoder.Reset(data)
	result := map[uint16]uint64{
		10: uint64(10),
		11: uint64(100),
		13: uint64(50),
	}
	c := 0
	total := 0
	for decoder.Next() {
		if decoder.HasValue() {
			assert.Equal(t, result[decoder.Slot()], decoder.Value())
			c++
		}
		total++
	}
	assert.Equal(t, 3, c)
	assert.Equal(t, 4, total)
}

func Test_Empty_TSDDecoder(t *testing.T) {
	decoder := NewTSDDecoder(nil)
	assert.Nil(t, decoder.Error())
}

package encoding

import (
	"fmt"
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
	assert.Nil(t, encoder.Error())

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)
	assert.Equal(t, 10, decoder.StartTime())
	assert.Equal(t, 13, decoder.EndTime())
	assert.Equal(t, 4, decoder.count)
	startTime, endTime := DecodeTSDTime(data)
	assert.Equal(t, 10, startTime)
	assert.Equal(t, 13, endTime)

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

	decoder = NewTSDDecoder(data)
	result := map[int]uint64{
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

func Test_Empty_TSDEncoderDecoder(t *testing.T) {
	encoder := NewTSDEncoder(1)
	encoder.AppendTime(bit.One)
	encoder.err = fmt.Errorf("error")
	encoder.AppendTime(bit.One)
	encoder.AppendValue(2)
	assert.NotNil(t, encoder.Error())

	decoder := NewTSDDecoder(nil)
	assert.Nil(t, decoder.Error())
}

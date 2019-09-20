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

func TestPrimitiveIterator(t *testing.T) {
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(100.0))
	data, _ := encoder.Bytes()
	it := NewPrimitiveIterator(10, field.Sum, data)
	assert.Equal(t, uint16(10), it.FieldID())
	assert.True(t, it.HasNext())
	slot, val := it.Next()
	assert.Equal(t, 10, slot)
	assert.Equal(t, 10.0, val)
	assert.True(t, it.HasNext())
	slot, val = it.Next()
	assert.Equal(t, 11, slot)
	assert.Equal(t, 100.0, val)
	assert.False(t, it.HasNext())
}

func TestBinaryFieldIterator(t *testing.T) {
	writer := stream.NewBufferWriter(nil)

	writer.PutVarint64(int64(10))
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.Bytes()
	writer.PutUInt16(uint16(10))
	writer.PutByte(byte(field.Sum))
	writer.PutVarint32(int32(len(data)))
	writer.PutBytes(data)
	d, _ := writer.Bytes()

	it := NewFieldIterator("f1", d)
	assert.Equal(t, int64(10), it.SegmentStartTime())
	assert.Equal(t, "f1", it.FieldMeta().Name)
	assert.True(t, it.HasNext())
	pIt := it.Next()
	assert.Equal(t, uint16(10), pIt.FieldID())
	assert.Equal(t, field.Sum, pIt.AggType())
	assert.True(t, pIt.HasNext())
	s, v := pIt.Next()
	assert.Equal(t, 12, s)
	assert.Equal(t, 10.0, v)
	assert.False(t, pIt.HasNext())
	assert.False(t, it.HasNext())

	_, err := it.Bytes()
	assert.Nil(t, err)
}

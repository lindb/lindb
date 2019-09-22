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
	writer.PutVarint64(10)
	writer.PutVarint32(int32(len(d)))
	writer.PutBytes(d)
	data, err := writer.Bytes()
	assert.NoError(t, err)
	result := make(map[string]string)
	it := NewGroupedIterator(nil, map[string][]byte{
		"f1": data,
		"f2": data,
	})
	assert.Nil(t, it.Tags())
	assert.True(t, it.HasNext())
	fIt := it.Next()
	assert.True(t, fIt.HasNext())
	result[fIt.FieldName()] = fIt.FieldName()
	startTime, fIt1 := fIt.Next()
	assert.Equal(t, int64(10), startTime)
	assertFieldIterator(t, fIt1)

	assert.True(t, it.HasNext())
	fIt = it.Next()
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
	assert.Equal(t, "f1", it.FieldName())
	assert.True(t, it.HasNext())
	startTime, fIt := it.Next()
	assert.Equal(t, int64(10), startTime)
	assertFieldIterator(t, fIt)
	assert.True(t, it.HasNext())
	startTime, fIt = it.Next()
	assert.Equal(t, int64(11), startTime)
	assertFieldIterator(t, fIt)
	assert.False(t, it.HasNext())

	writer = stream.NewBufferWriter(nil)
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

func TestPrimitiveIterator(t *testing.T) {
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(100.0))
	data, _ := encoder.Bytes()
	decoder := encoding.GetTSDDecoder()
	decoder.Reset(data)
	it := NewPrimitiveIterator(10, field.Sum, decoder)
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
	d := buildFieldIterator()
	it := NewFieldIterator(d)
	assertFieldIterator(t, it)

	_, err := it.Bytes()
	assert.Error(t, err)
}

func assertFieldIterator(t *testing.T, it FieldIterator) {
	for i := 0; i < 2; i++ {
		assert.True(t, it.HasNext())
		pIt := it.Next()
		assert.Equal(t, uint16(10), pIt.FieldID())
		assert.Equal(t, field.Sum, pIt.AggType())
		assert.True(t, pIt.HasNext())
		s, v := pIt.Next()
		assert.Equal(t, 12, s)
		assert.Equal(t, 10.0, v)
	}
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
	for i := 0; i < 2; i++ {
		writer.PutUInt16(uint16(10))
		writer.PutByte(byte(field.Sum))
		writer.PutVarint32(int32(len(data)))
		writer.PutBytes(data)
	}
	d, _ := writer.Bytes()
	return d
}

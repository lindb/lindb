package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutInt(t *testing.T) {
	writer := BinaryWriter()
	writer.PutUint32(uint32(123))
	writer.PutUint64(uint64(456))

	buf, err := writer.Bytes()
	if err != nil {
		t.Error(err)
	}

	reader := BinaryReader(buf)
	assert.Equal(t, reader.ReadUint32(), uint32(123))
	assert.Equal(t, reader.ReadUint64(), uint64(456))
}

func Test_Write(t *testing.T) {
	writer := BinaryBufWriter(nil)
	writer.PutByte(1)
	assert.Equal(t, 1, writer.Len())

	writer.PutLenBytes([]byte{2, 3, 4})
	assert.Equal(t, 5, writer.Len())

	writer.PutVarint32(2)
	assert.Equal(t, 6, writer.Len())

	writer.PutVarint64(2)
	assert.Equal(t, 7, writer.Len())

	writer.PutUvarint32(11)
	assert.Equal(t, 8, writer.Len())

	writer.err = fmt.Errorf("error")
	writer.PutUvarint64(12)
	assert.Equal(t, 8, writer.Len())
	writer.err = nil
	writer.PutUvarint64(12)
	assert.Equal(t, 9, writer.Len())

	writer.PutInt32(1)
	assert.Equal(t, 13, writer.Len())

	writer.PutInt64(1)
	assert.Equal(t, 21, writer.Len())

	writer.err = fmt.Errorf("error")
	_, err := writer.Bytes()
	assert.NotNil(t, err)
}

func Test_Read(t *testing.T) {
	var buf []byte
	writer := BinaryBufWriter(buf)

	newReader := func() *Binary {
		data, _ := writer.Bytes()
		writer.buf.Reset()
		return BinaryReader(data)
	}

	writer.PutVarint32(1)
	assert.Equal(t, int32(1), newReader().ReadVarint32())

	writer.PutVarint64(2)
	assert.Equal(t, int64(2), newReader().ReadVarint64())

	writer.PutUvarint32(3)
	writer.PutUvarint64(4)
	reader := newReader()
	assert.Equal(t, uint32(3), reader.ReadUvarint32())
	assert.Equal(t, uint64(4), reader.ReadUvarint64())

	writer.PutInt32(5)
	writer.PutInt64(6)
	writer.PutByte(7)
	reader = newReader()
	assert.Equal(t, int32(5), reader.ReadInt32())
	assert.Equal(t, int64(6), reader.ReadInt64())
	assert.Equal(t, byte(7), reader.ReadByte())
	assert.Equal(t, byte(0), reader.ReadByte())
	assert.True(t, reader.Empty())
	assert.Nil(t, reader.Error())

	// test error read uvariant64
	reader = BinaryReader(nil)
	reader.ReadUvarint64()
	assert.NotNil(t, reader.Error())
}

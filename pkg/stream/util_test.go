package stream

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutUint64(t *testing.T) {
	buf := make([]byte, 16)
	PutUint64(buf, 0, 20)
	PutUint64(buf, 8, 40)
	assert.Equal(t, uint64(20), ReadUint64(buf, 0))
	assert.Equal(t, uint64(40), ReadUint64(buf, 8))
}

func Test_ReadUint32(t *testing.T) {
	buf := make([]byte, 8)
	PutUint32(buf, 0, 20)
	PutUint32(buf, 4, 40)
	assert.Equal(t, uint32(20), ReadUint32(buf, 0))
	assert.Equal(t, uint32(40), ReadUint32(buf, 4))
}

func Test_ReadUint16(t *testing.T) {
	var buf bytes.Buffer
	writer2 := NewBufferWriter(&buf)
	writer2.PutUInt16(2)
	writer2.PutUInt16(999)
	data, err := writer2.Bytes()
	assert.NoError(t, err)
	assert.Len(t, data, 4)
	assert.Equal(t, uint16(999), ReadUint16(data, 2))
	assert.Equal(t, uint16(2), ReadUint16(data, 0))
}

func Test_ReadUvarint(t *testing.T) {
	var buf bytes.Buffer
	writer2 := NewBufferWriter(&buf)
	writer2.PutUvarint64(999)
	writer2.PutUvarint64(889)
	data, err := writer2.Bytes()
	assert.NoError(t, err)
	v, l, err := ReadUvarint(data, 0)
	assert.NoError(t, err)
	assert.True(t, l > 0)
	assert.Equal(t, uint64(999), v)
	v, l2, err := ReadUvarint(data, l)
	assert.NoError(t, err)
	assert.True(t, l2 > 0)
	assert.Equal(t, uint64(889), v)
	assert.Equal(t, len(data), l+l2)

	d := make([]byte, 10)
	for i := 0; i < 10; i++ {
		d[i] = 0xa0
	}
	d[9] = 0x60
	_, l, err = ReadUvarint(d, 0)
	assert.Error(t, err)
	assert.Equal(t, 10, l)

	d = make([]byte, 20)
	for i := 0; i < 20; i++ {
		d[i] = 0xa0
	}
	d[19] = 0x60
	_, l, err = ReadUvarint(d, 0)
	assert.Error(t, err)
	assert.Equal(t, 20, l)
}

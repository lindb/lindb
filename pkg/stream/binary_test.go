package stream

import (
	"testing"

	"github.com/magiconair/properties/assert"
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

func TestByteBufReader_ReadBytes(t *testing.T) {

	keys := []string{
		"apple",
		"orange",
		"apple pie",
		"lemon meringue",
		"lemon",
	}

	writer := BinaryWriter()
	writer.PutUint32(uint32(len(keys)))

	for _, k := range keys {
		writer.PutLenBytes([]byte(k))
	}

	by, err := writer.Bytes()
	assert.Equal(t, nil, err)

	assert.Equal(t, writer.Len(), len(by))

	reader := NewBufReader(by)
	count := reader.ReadUint32()
	assert.Equal(t, uint32(len(keys)), count)

	for i := 0; i < int(count); i++ {
		_, k := reader.ReadLenBytes()
		assert.Equal(t, keys[i], string(k))
	}

}

package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteBufReader_ReadBytes(t *testing.T) {

	keys := []string{
		"apple",
		"orange",
		"apple pie",
		"lemon meringue",
		"lemon",
	}

	writer := BinaryWriter()
	writer.PutUInt32(uint32(len(keys)))

	for _, k := range keys {
		writer.PutKey([]byte(k))
	}

	by, err := writer.Bytes()
	assert.Equal(t, nil, err)

	assert.Equal(t, writer.Len(), len(by))

	reader := NewBufReader(by)
	count := reader.ReadUInt32()
	assert.Equal(t, uint32(len(keys)), count)

	for i := 0; i < int(count); i++ {
		_, k := reader.ReadKey()
		assert.Equal(t, keys[i], string(k))
	}

}

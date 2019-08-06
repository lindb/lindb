package stream

import (
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

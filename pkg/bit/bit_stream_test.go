package bit

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOkWriter struct{}

func (w *mockOkWriter) Write(p []byte) (n int, err error) {
	return 1, nil
}

type mockErrWriter struct{}

func (w *mockErrWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error")
}

var (
	okWriter  = NewWriter(&mockOkWriter{})
	badWriter = NewWriter(&mockErrWriter{})
)

func Test_Writer_WriteBit(t *testing.T) {
	for range [10]struct{}{} {
		assert.Nil(t, okWriter.WriteBit(Zero))
		assert.Nil(t, okWriter.WriteBit(One))
	}

	for range [7]struct{}{} {
		assert.Nil(t, badWriter.WriteBit(Zero))
	}
	assert.NotNil(t, badWriter.WriteBit(One))
}

func Test_Writer_WriteBytes(t *testing.T) {
	assert.Nil(t, okWriter.WriteBits(math.MaxUint64, 63))

	assert.NotNil(t, badWriter.WriteBits(math.MaxUint64, 63))
}

func Test_Writer_Flush(t *testing.T) {
	okWriter.count = 8
	assert.Nil(t, okWriter.Flush())

	badWriter.count = 1
	assert.NotNil(t, badWriter.Flush())
	assert.NotNil(t, badWriter.Flush())
}

func Test_Reader(t *testing.T) {
	var buffer bytes.Buffer
	reader := NewReader(&buffer)

	_, err := reader.ReadBit()
	assert.NotNil(t, err)
	_, err = reader.ReadByte()
	assert.NotNil(t, err)
	_, err = reader.ReadBits(10)
	assert.NotNil(t, err)
	_, err = reader.ReadBits(1)
	assert.NotNil(t, err)

	buffer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	reader = NewReader(&buffer)
	reader.ReadBits(10)
}

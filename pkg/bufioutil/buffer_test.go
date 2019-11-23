package bufioutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	buffer := NewBuffer([]byte{1, 2, 3})
	b, err := buffer.GetByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), b)

	buffer.SetIdx(100)
	_, err = buffer.GetByte()
	assert.Equal(t, errOutOfRange, err)

	// reset
	buffer.SetBuf([]byte{1, 2, 3})
	buffer.SetIdx(1)
	b, err = buffer.GetByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(2), b)
}

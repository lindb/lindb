package bufpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BufferPool(t *testing.T) {
	buf1 := GetBuffer()
	assert.Equal(t, 0, buf1.Len())

	buf1.WriteByte(byte(66))
	assert.Equal(t, 1, buf1.Len())
	assert.Equal(t, 64, buf1.Cap())

	PutBuffer(buf1)

	buf2 := GetBuffer()
	assert.Equal(t, 0, buf2.Len())
}

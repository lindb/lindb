package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BitArray(t *testing.T) {
	ba, err := NewBitArray(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", ba.String())
	assert.False(t, ba.GetBit(0))

	ba.SetBit(uint16(0))
	ba.SetBit(uint16(2))
	assert.Len(t, ba.Bytes(), 1)

	assert.Equal(t, "10100000", ba.String())
	ba.SetBit(uint16(4))
	assert.Equal(t, "10101000", ba.String())

	ba.SetBit(uint16(8))
	assert.Equal(t, "1010100010000000", ba.String())
	ba.SetBit(uint16(9))
	ba.SetBit(uint16(9))
	assert.Equal(t, "1010100011000000", ba.String())
	ba.SetBit(uint16(16))
	assert.Equal(t, "101010001100000010000000", ba.String())

	assert.True(t, ba.GetBit(0))
	assert.False(t, ba.GetBit(1))
	assert.True(t, ba.GetBit(8))
	assert.True(t, ba.GetBit(9))
	assert.False(t, ba.GetBit(23))
	assert.False(t, ba.GetBit(24))
	assert.False(t, ba.GetBit(800))

	ba.Reset()
	assert.False(t, ba.GetBit(0))

	buf := make([]byte, 65537)
	ba2, err := NewBitArray(buf)
	assert.Nil(t, ba2)
	assert.NotNil(t, err)

	buf = make([]byte, 2)
	ba2, err = NewBitArray(buf)
	assert.Nil(t, err)
	assert.NotNil(t, ba2)

	ba3, err := NewBitArray([]byte{})
	assert.Nil(t, err)
	assert.NotNil(t, ba3)
	assert.False(t, ba3.GetBit(23))

}

package metrictbl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bitArray(t *testing.T) {
	ba := bitArray{}
	assert.Equal(t, "", ba.String())
	assert.False(t, ba.getBit(0))

	ba.setBit(uint16(0))
	ba.setBit(uint16(2))
	assert.Len(t, ba.payload, 1)

	assert.Equal(t, "10100000", ba.String())
	ba.setBit(uint16(4))
	assert.Equal(t, "10101000", ba.String())

	ba.setBit(uint16(8))
	assert.Equal(t, "1010100010000000", ba.String())
	ba.setBit(uint16(9))
	ba.setBit(uint16(9))
	assert.Equal(t, "1010100011000000", ba.String())
	ba.setBit(uint16(16))
	assert.Equal(t, "101010001100000010000000", ba.String())

	assert.True(t, ba.getBit(0))
	assert.False(t, ba.getBit(1))
	assert.True(t, ba.getBit(8))
	assert.True(t, ba.getBit(9))
	assert.False(t, ba.getBit(23))
	assert.False(t, ba.getBit(24))
	assert.False(t, ba.getBit(800))

	ba.reset()
	assert.False(t, ba.getBit(0))

	buf := make([]byte, 65537)
	ba2, err := newBitArray(buf)
	assert.Nil(t, ba2)
	assert.NotNil(t, err)

	buf = make([]byte, 2)
	ba2, err = newBitArray(buf)
	assert.Nil(t, err)
	assert.NotNil(t, ba2)

	ba3, err := newBitArray([]byte{})
	assert.Nil(t, err)
	assert.NotNil(t, ba3)
	assert.False(t, ba3.getBit(23))

}

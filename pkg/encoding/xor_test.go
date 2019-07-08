package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	e := NewXOREncoder()
	e.Write(uint64(76))
	e.Write(uint64(50))
	e.Write(uint64(50))
	e.Write(uint64(999999999))
	e.Write(uint64(100))

	data, err := e.Bytes()
	assert.Nil(t, err)

	d := NewXORDecoder(data)
	exceptIntValue(d, t, uint64(76))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(999999999))
	exceptIntValue(d, t, uint64(100))
}

func exceptIntValue(d *XORDecoder, t *testing.T, except uint64) {
	assert.True(t, d.Next())
	assert.Equal(t, except, d.Value())
}

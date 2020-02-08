package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZigZag(t *testing.T) {
	v := ZigZagEncode(1)
	assert.Equal(t, int64(1), ZigZagDecode(v))

	v = ZigZagEncode(-99999)
	assert.Equal(t, int64(-99999), ZigZagDecode(v))
}

func TestHighLowBits(t *testing.T) {
	v := uint32(67043434)
	high := HighBits(v)
	low := LowBits(v)
	assert.Equal(t, v, ValueWithHighLowBits(uint32(high)<<16, low))
}

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

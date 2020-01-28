package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntBlockAlloc(t *testing.T) {
	bs := newBlockStore(0)
	assert.NotNil(t, bs)

	// int block
	b1 := bs.allocBlock()
	assert.NotNil(t, b1)
	bs.freeBlock(b1)
	b2 := bs.allocBlock()
	assert.NotNil(t, b2)
	b3 := bs.allocBlock()
	assert.True(t, b1 != b3)
}

func TestTimeWindowRange(t *testing.T) {
	bs := newBlockStore(30)

	b1 := bs.allocBlock()
	assert.True(t, b1.isEmpty())
	b1.setValue(10, float64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, float64(100), b1.getValue(10))
	assert.True(t, b1.memsize() > 0)
	assert.Equal(t, uint16(10), b1.getSize())
}

func TestFloatReset(t *testing.T) {
	bs := newBlockStore(30)

	b1 := bs.allocBlock()
	b1.setValue(11, float64(100))
	assert.True(t, b1.hasValue(11))
	assert.Equal(t, float64(100), b1.getValue(11))
	assert.Equal(t, uint16(11), b1.getSize())
	b1.setValue(20, float64(100))
	assert.Equal(t, uint16(20), b1.getSize())
	b1.reset()
	assert.False(t, b1.hasValue(11))
	assert.Equal(t, uint16(0), b1.getSize())
}

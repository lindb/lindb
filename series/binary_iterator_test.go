package series

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestPrimitiveIterator(t *testing.T) {
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(100.0))
	data, _ := encoder.Bytes()
	it := NewPrimitiveIterator(10, data)
	assert.Equal(t, uint16(10), it.FieldID())
	assert.True(t, it.HasNext())
	slot, val := it.Next()
	assert.Equal(t, 10, slot)
	assert.Equal(t, 10.0, val)
	assert.True(t, it.HasNext())
	slot, val = it.Next()
	assert.Equal(t, 11, slot)
	assert.Equal(t, 100.0, val)
	assert.False(t, it.HasNext())
}

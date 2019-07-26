package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatArray(t *testing.T) {
	fa := NewFloatArray(10)
	assert.True(t, fa.IsEmpty())
	fa.SetValue(0, 1.1)
	fa.SetValue(5, 5.5)
	fa.SetValue(8, 9.9)
	fa.SetValue(-1, 1.1)
	fa.SetValue(10, 11.1)
	fa.SetValue(11, 11.1)
	assert.False(t, fa.IsEmpty())
	assert.True(t, fa.HasValue(0))
	assert.True(t, fa.HasValue(5))
	assert.False(t, fa.HasValue(-1))
	assert.False(t, fa.HasValue(10))
	assert.False(t, fa.HasValue(11))

	assert.Equal(t, float64(0), fa.GetValue(-1))
	assert.Equal(t, 1.1, fa.GetValue(0))
	assert.Equal(t, 5.5, fa.GetValue(5))
	assert.Equal(t, 9.9, fa.GetValue(8))
	assert.Equal(t, float64(0), fa.GetValue(10))
	assert.Equal(t, float64(0), fa.GetValue(11))

	assert.Equal(t, 3, fa.Size())

	it := fa.Iterator()
	assert.True(t, it.HasNext())
	idx, value := it.Next()
	assert.Equal(t, 0, idx)
	assert.Equal(t, 1.1, value)
	assert.True(t, it.HasNext())
	idx, value = it.Next()
	assert.Equal(t, 5, idx)
	assert.Equal(t, 5.5, value)
	assert.True(t, it.HasNext())
	idx, value = it.Next()
	assert.Equal(t, 8, idx)
	assert.Equal(t, 9.9, value)
	assert.False(t, it.HasNext())
	idx, value = it.Next()
	assert.Equal(t, -1, idx)
	assert.Equal(t, float64(0), value)

	// reset
	fa.SetValue(8, 10.10)
	assert.Equal(t, 10.10, fa.GetValue(8))
	assert.Equal(t, 3, fa.Size())
}

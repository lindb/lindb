package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/field"
)

func TestFieldIterator(t *testing.T) {
	primitiveIt := newPrimitiveIterator(uint16(10), generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))
	primitiveIt1 := newPrimitiveIterator(uint16(10), generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	it := newFieldIterator(uint16(111), field.SumField, []field.PrimitiveIterator{primitiveIt, primitiveIt1})

	except := map[int]float64{0: 0, 1: 10, 2: 10.0, 3: 100.4, 4: 50.0}
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), except)
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), except)

	assert.False(t, it.HasNext())
	assert.Nil(t, it.Next())
	assert.Equal(t, uint16(111), it.ID())
	assert.Equal(t, field.SumField, it.FieldType())
}

func TestPrimitiveIterator(t *testing.T) {
	it := newPrimitiveIterator(uint16(10), generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	except := map[int]float64{0: 0, 1: 10, 2: 10.0, 3: 100.4, 4: 50.0}
	AssertPrimitiveIt(t, it, except)

	assert.False(t, it.HasNext())
	timeSlot, value := it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
	assert.Equal(t, uint16(10), it.ID())

	it = newPrimitiveIterator(uint16(10), nil)
	assert.False(t, it.HasNext())
	timeSlot, value = it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
}

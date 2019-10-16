package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func TestFieldIterator(t *testing.T) {
	it := newFieldIterator(10, []series.PrimitiveIterator{})
	assert.False(t, it.HasNext())
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)

	primitiveIt := newPrimitiveIterator(uint16(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))
	primitiveIt1 := newPrimitiveIterator(uint16(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	it = newFieldIterator(10, []series.PrimitiveIterator{primitiveIt, primitiveIt1})

	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), expect)
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), expect)

	assert.False(t, it.HasNext())
	assert.Nil(t, it.Next())

	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
}

func TestFieldIterator_MarshalBinary(t *testing.T) {
	floatArray := collections.NewFloatArray(5)
	floatArray.SetValue(2, 10.0)
	primitiveIt := newPrimitiveIterator(uint16(10), 10, field.Sum, floatArray)
	primitiveIt1 := newPrimitiveIterator(uint16(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	it := newFieldIterator(10, []series.PrimitiveIterator{primitiveIt, primitiveIt1})
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
}

func TestPrimitiveIterator(t *testing.T) {
	it := newPrimitiveIterator(uint16(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	AssertPrimitiveIt(t, it, expect)

	assert.False(t, it.HasNext())
	timeSlot, value := it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
	assert.Equal(t, uint16(10), it.FieldID())

	it = newPrimitiveIterator(uint16(10), 10, field.Sum, nil)
	assert.False(t, it.HasNext())
	timeSlot, value = it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
}

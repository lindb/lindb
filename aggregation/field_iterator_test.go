package aggregation

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

var encodeFunc = encoding.NewTSDEncoder

func TestFieldIterator(t *testing.T) {
	it := newFieldIterator(10, []series.PrimitiveIterator{})
	assert.False(t, it.HasNext())
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)

	primitiveIt := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))
	primitiveIt1 := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	it = newFieldIterator(10, []series.PrimitiveIterator{primitiveIt, primitiveIt1})

	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), expect)
	assert.True(t, it.HasNext())
	AssertPrimitiveIt(t, it.Next(), expect)

	assert.False(t, it.HasNext())
	assert.Nil(t, it.Next())

	// marshal empty primitive iterator, because primitive iterator already read
	data, err = it.MarshalBinary()
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestFieldIterator_MarshalBinary(t *testing.T) {
	floatArray := collections.NewFloatArray(5)
	floatArray.SetValue(2, 10.0)
	primitiveIt := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, floatArray)
	primitiveIt1 := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	it := newFieldIterator(10, []series.PrimitiveIterator{primitiveIt, primitiveIt1})
	data, err := it.MarshalBinary()
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	reader := stream.NewReader(data)
	pFieldID := reader.ReadByte() // read primitive field id
	assert.Equal(t, field.PrimitiveID(10), field.PrimitiveID(pFieldID))
	aggType := field.AggType(reader.ReadByte())
	assert.Equal(t, field.Sum, aggType)
	length := reader.ReadVarint32()
	data1 := reader.ReadBytes(int(length))

	pIt := series.NewPrimitiveIterator(field.PrimitiveID(pFieldID), aggType, encoding.NewTSDDecoder(data1))
	assert.True(t, pIt.HasNext())
	i, value := pIt.Next()
	assert.Equal(t, 12, i)
	assert.Equal(t, 10.0, value)
}

func TestFieldIterator_MarshalBinary_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()
	encoder := encoding.NewMockTSDEncoder(ctrl)
	encoding.TSDEncodeFunc = func(startTime uint16) encoding.TSDEncoder {
		return encoder
	}
	floatArray := collections.NewFloatArray(5)
	floatArray.SetValue(2, 10.0)
	primitiveIt := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, floatArray)
	encoder.EXPECT().AppendTime(gomock.Any()).AnyTimes()
	encoder.EXPECT().AppendValue(gomock.Any()).AnyTimes()
	encoder.EXPECT().Bytes().Return(nil, fmt.Errorf("err"))
	it := newFieldIterator(10, []series.PrimitiveIterator{primitiveIt})
	data, err := it.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestPrimitiveIterator(t *testing.T) {
	it := newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, generateFloatArray([]float64{0, 10, 10.0, 100.4, 50.0}))

	expect := map[int]float64{10: 0, 11: 10, 12: 10.0, 13: 100.4, 14: 50.0}
	AssertPrimitiveIt(t, it, expect)

	assert.False(t, it.HasNext())
	timeSlot, value := it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
	assert.Equal(t, field.PrimitiveID(10), it.FieldID())

	it = newPrimitiveIterator(field.PrimitiveID(10), 10, field.Sum, nil)
	assert.False(t, it.HasNext())
	timeSlot, value = it.Next()
	assert.Equal(t, -1, timeSlot)
	assert.Equal(t, float64(0), value)
}

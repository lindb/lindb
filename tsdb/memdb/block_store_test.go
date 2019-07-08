package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/field"
)

func TestBlockAlloc(t *testing.T) {
	bs := newBlockStore(10)

	// int block
	b1 := bs.allocIntBlock()
	assert.NotNil(t, b1)
	bs.freeIntBlock(b1)
	b2 := bs.allocIntBlock()
	assert.NotNil(t, b2)
	b3 := bs.allocIntBlock()
	assert.True(t, b1 != b3)

	// float block
	bf := bs.allocFloatBlock()
	assert.NotNil(t, bf)
	bs.freeFloatBlock(bf)
	bf2 := bs.allocFloatBlock()
	assert.NotNil(t, bf2)
	bf3 := bs.allocFloatBlock()
	assert.True(t, bf != bf3)
}
func TestTimeWindowRange(t *testing.T) {
	bs := newBlockStore(30)

	// int block
	b1 := bs.allocIntBlock()
	b1.setStartTime(10)
	b1.setValue(10)
	b1.updateValue(10, int64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, int64(100), b1.getValue(10))
	assert.Equal(t, 10, b1.getStartTime())
	assert.Equal(t, 20, b1.getEndTime())
	b1.setStartTime(40)
	b1.setValue(0)
	b1.updateValue(0, int64(100))
	assert.False(t, b1.hasValue(10))
	assert.Equal(t, 40, b1.getStartTime())
	assert.Equal(t, 40, b1.getEndTime())

	// float block
	b2 := bs.allocFloatBlock()
	b2.setStartTime(10)
	b2.setValue(10)
	b2.updateValue(10, 10.0)
	assert.True(t, b2.hasValue(10))
	assert.Equal(t, 10.0, b2.getValue(10))
	assert.Equal(t, 10, b2.getStartTime())
	b2.setStartTime(40)
	b2.setValue(0)
	b2.updateValue(0, 10.90)
	assert.False(t, b2.hasValue(10))
	assert.Equal(t, 40, b2.getStartTime())
	assert.Equal(t, 40, b2.getEndTime())
}

func TestReset(t *testing.T) {
	bs := newBlockStore(30)

	// int block
	b1 := bs.allocIntBlock()
	b1.setValue(10)
	b1.updateValue(10, int64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, int64(100), b1.getValue(10))
	b1.reset()
	assert.False(t, b1.hasValue(10))

	// float block
	b2 := bs.allocFloatBlock()
	b2.setValue(10)
	b2.updateValue(10, 10.0)
	assert.True(t, b2.hasValue(10))
	assert.Equal(t, 10.0, b2.getValue(10))
	b2.reset()
	assert.False(t, b2.hasValue(10))
}

func TestCompactIntBlock(t *testing.T) {
	bs := newBlockStore(30)

	// int block
	b1 := bs.allocIntBlock()
	b1.setStartTime(10)
	b1.setValue(10)
	b1.updateValue(10, int64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, int64(100), b1.getValue(10))
	assert.Equal(t, 10, b1.getStartTime())
	assert.Equal(t, 20, b1.getEndTime())

	// test compact [10,20] and no compress => [10,20]
	b1.compact(field.GetAggFunc(field.Sum))

	tsd := encoding.NewTSDDecoder(b1.compress)
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(100), encoding.ZigZagDecode(tsd.Value()))

	b1.setStartTime(10)
	b1.setValue(10)
	b1.updateValue(10, int64(100))

	// test compact [10,20] and compress[10,20] => [10,20]
	b1.compact(field.GetAggFunc(field.Sum))

	tsd = encoding.NewTSDDecoder(b1.compress)
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(200), encoding.ZigZagDecode(tsd.Value()))

	b1.setStartTime(10)
	b1.setValue(0)
	b1.updateValue(0, int64(50))
	b1.setValue(11)
	b1.updateValue(11, int64(100))

	// test compact [10,21] and compress[10,20] => [10,21]
	b1.compact(field.GetAggFunc(field.Sum))

	tsd = encoding.NewTSDDecoder(b1.compress)
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 21, tsd.EndTime())
	assert.True(t, tsd.HasValueWithSlot(0))
	assert.Equal(t, int64(50), encoding.ZigZagDecode(tsd.Value()))
	for i := 1; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(200), encoding.ZigZagDecode(tsd.Value()))
	assert.True(t, tsd.HasValueWithSlot(11))
	assert.Equal(t, int64(100), encoding.ZigZagDecode(tsd.Value()))

	b1.setStartTime(40)
	b1.setValue(11)
	b1.updateValue(11, int64(90))

	// test compact [40,51] and compress[10,21] => [10,51]
	b1.compact(field.GetAggFunc(field.Sum))

	tsd = encoding.NewTSDDecoder(b1.compress)
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 51, tsd.EndTime())
	assert.True(t, tsd.HasValueWithSlot(0))
	assert.Equal(t, int64(50), encoding.ZigZagDecode(tsd.Value()))
	for i := 1; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(200), encoding.ZigZagDecode(tsd.Value()))
	assert.True(t, tsd.HasValueWithSlot(11))
	assert.Equal(t, int64(100), encoding.ZigZagDecode(tsd.Value()))
	for i := 12; i < 41; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(41))
	assert.Equal(t, int64(90), encoding.ZigZagDecode(tsd.Value()))
}

package memdb

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

func TestBlockAlloc(t *testing.T) {
	bs := newBlockStore(-1)
	assert.NotNil(t, bs)
	bs = newBlockStore(10)

	// int block
	b1 := bs.allocIntBlock()
	assert.NotNil(t, b1)
	bs.freeBlock(b1)
	b2 := bs.allocIntBlock()
	assert.NotNil(t, b2)
	b3 := bs.allocIntBlock()
	assert.True(t, b1 != b3)

	// float block
	bf := bs.allocFloatBlock()
	assert.NotNil(t, bf)
	bs.freeBlock(bf)
	bf2 := bs.allocFloatBlock()
	assert.NotNil(t, bf2)
	bf3 := bs.allocFloatBlock()
	assert.True(t, bf != bf3)
}

func TestTimeWindowRange(t *testing.T) {
	bs := newBlockStore(30)

	// int block
	b1 := bs.allocIntBlock()
	assert.True(t, b1.isEmpty())
	b1.setStartTime(10)
	b1.setIntValue(10, int64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, int64(100), b1.getIntValue(10))
	assert.Equal(t, 10, b1.getStartTime())
	assert.Equal(t, 20, b1.getEndTime())
	b1.setStartTime(40)
	b1.setIntValue(0, int64(100))
	assert.False(t, b1.hasValue(10))
	assert.Equal(t, 40, b1.getStartTime())
	assert.Equal(t, 40, b1.getEndTime())
	assert.False(t, b1.isEmpty())

	// float block
	b2 := bs.allocFloatBlock()
	b2.setStartTime(10)
	b2.setFloatValue(10, 10.0)
	assert.True(t, b2.hasValue(10))
	assert.Equal(t, 10.0, b2.getFloatValue(10))
	assert.Equal(t, 10, b2.getStartTime())
	b2.setStartTime(40)
	b2.setFloatValue(0, 10.90)
	assert.False(t, b2.hasValue(10))
	assert.Equal(t, 40, b2.getStartTime())
	assert.Equal(t, 40, b2.getEndTime())
}

func TestReset(t *testing.T) {
	bs := newBlockStore(30)

	// int block
	b1 := bs.allocIntBlock()
	b1.setIntValue(11, int64(100))
	assert.True(t, b1.hasValue(11))
	assert.Equal(t, int64(100), b1.getIntValue(11))
	b1.reset()
	assert.False(t, b1.hasValue(11))

	// float block
	b2 := bs.allocFloatBlock()
	b2.setFloatValue(11, 10.0)
	assert.True(t, b2.hasValue(11))
	assert.Equal(t, 10.0, b2.getFloatValue(11))
	b2.reset()
	assert.False(t, b2.hasValue(11))
}

func TestCompactIntBlock(t *testing.T) {
	bs := newBlockStore(30)

	assert.Nil(t, bs.allocBlock(field.ValueType(uint8(9))))
	// int block
	b1 := bs.allocBlock(field.Integer)
	start, end, err := b1.compact(field.GetAggFunc(field.Sum), true)
	assert.Nil(t, err)
	assert.Equal(t, 0, start)
	assert.Equal(t, 0, end)

	b1.setStartTime(10)
	b1.setIntValue(10, int64(100))
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, int64(100), b1.getIntValue(10))
	assert.Equal(t, 10, b1.getStartTime())
	assert.Equal(t, 20, b1.getEndTime())

	// test compact [10,20] and no compress => [10,20]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	assert.Nil(t, err)
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)

	tsd := encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(100), encoding.ZigZagDecode(tsd.Value()))

	b1.setStartTime(10)
	b1.setIntValue(10, int64(100))

	// test compact [10,20] and compress[10,20] => [10,20]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, int64(200), encoding.ZigZagDecode(tsd.Value()))

	b1.setStartTime(10)
	b1.setIntValue(0, int64(50))
	b1.setIntValue(11, int64(100))

	// test compact [10,21] and compress[10,20] => [10,21]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 21, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
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
	b1.setIntValue(11, int64(90))

	// test compact [40,51] and compress[10,21] => [10,51]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 51, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
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

func TestCompactFloatBlock(t *testing.T) {
	bs := newBlockStore(30)

	// float block
	b1 := bs.allocBlock(field.Float)
	start, end, err := b1.compact(field.GetAggFunc(field.Sum), true)
	assert.Nil(t, err)
	assert.Equal(t, 0, start)
	assert.Equal(t, 0, end)

	b1.setStartTime(10)
	b1.setFloatValue(10, 100.05)
	assert.True(t, b1.hasValue(10))
	assert.Equal(t, 100.05, b1.getFloatValue(10))
	assert.Equal(t, 10, b1.getStartTime())
	assert.Equal(t, 20, b1.getEndTime())

	// test compact [10,20] and no compress => [10,20]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)
	data1 := b1.bytes()
	assert.Equal(t, data1, b1.bytes())

	tsd := encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, 100.05, math.Float64frombits(tsd.Value()))

	b1.setStartTime(10)
	b1.setFloatValue(10, 100.05)

	// test compact [10,20] and compress[10,20] => [10,20]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 20, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 20, tsd.EndTime())
	for i := 0; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, 200.1, math.Float64frombits(tsd.Value()))

	b1.setStartTime(10)
	b1.setFloatValue(0, 50.0)
	b1.setFloatValue(11, 100.0)

	// test compact [10,21] and compress[10,20] => [10,21]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 21, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 21, tsd.EndTime())
	assert.True(t, tsd.HasValueWithSlot(0))
	assert.Equal(t, 50.0, math.Float64frombits(tsd.Value()))
	for i := 1; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, 200.1, math.Float64frombits(tsd.Value()))
	assert.True(t, tsd.HasValueWithSlot(11))
	assert.Equal(t, 100.0, math.Float64frombits(tsd.Value()))

	b1.setStartTime(40)
	b1.setFloatValue(11, 90.0)

	// test compact [40,51] and compress[10,21] => [10,51]
	start, end, err = b1.compact(field.GetAggFunc(field.Sum), true)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, start)
	assert.Equal(t, 51, end)

	tsd = encoding.NewTSDDecoder(b1.bytes())
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 51, tsd.EndTime())
	assert.True(t, tsd.HasValueWithSlot(0))
	assert.Equal(t, 50.0, math.Float64frombits(tsd.Value()))
	for i := 1; i < 10; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(10))
	assert.Equal(t, 200.1, math.Float64frombits(tsd.Value()))
	assert.True(t, tsd.HasValueWithSlot(11))
	assert.Equal(t, 100.0, math.Float64frombits(tsd.Value()))
	for i := 12; i < 41; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(41))
	assert.Equal(t, 90.0, math.Float64frombits(tsd.Value()))
}

func TestContainer_Get_Set(t *testing.T) {
	c := &container{}
	c.setFloatValue(10, 10.0)
	assert.Equal(t, 0.0, c.getFloatValue(10))
	c.setIntValue(10, 10)
	assert.Equal(t, int64(0), c.getIntValue(10))
}

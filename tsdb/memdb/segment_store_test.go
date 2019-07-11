package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/field"
)

func TestSimpleSegmentStore(t *testing.T) {
	aggFunc := field.GetAggFunc(field.Sum)
	store := newSimpleFieldStore(aggFunc)
	assert.NotNil(t, store)
	ss, ok := store.(*simpleFieldStore)
	assert.True(t, ok)

	compress, startSlot, endSlot, err := store.bytes()
	assert.Nil(t, compress)
	assert.NotNil(t, err)
	assert.Equal(t, 0, startSlot)
	assert.Equal(t, 0, endSlot)

	bs := newBlockStore(30)
	ss.writeInt(bs, 10, int64(100))
	// memory auto rollup
	ss.writeInt(bs, 11, int64(110))
	// memory auto rollup
	ss.writeInt(bs, 10, int64(100))
	// compact because slot out of current time window
	ss.writeInt(bs, 40, int64(20))
	// compact before time window
	ss.writeInt(bs, 10, int64(100))
	// compact because slot out of current time window
	ss.writeInt(bs, 41, int64(50))

	compress, startSlot, endSlot, err = store.bytes()
	assert.Nil(t, err)
	assert.Equal(t, 10, startSlot)
	assert.Equal(t, 41, endSlot)

	tsd := encoding.NewTSDDecoder(compress)
	assert.Equal(t, 10, tsd.StartTime())
	assert.Equal(t, 41, tsd.EndTime())
	assert.True(t, tsd.HasValueWithSlot(0))
	assert.Equal(t, int64(300), encoding.ZigZagDecode(tsd.Value()))
	assert.True(t, tsd.HasValueWithSlot(1))
	assert.Equal(t, int64(110), encoding.ZigZagDecode(tsd.Value()))
	for i := 1; i < 30; i++ {
		assert.False(t, tsd.HasValueWithSlot(i))
	}
	assert.True(t, tsd.HasValueWithSlot(30))
	assert.Equal(t, int64(20), encoding.ZigZagDecode(tsd.Value()))

	assert.True(t, tsd.HasValueWithSlot(31))
	assert.Equal(t, int64(50), encoding.ZigZagDecode(tsd.Value()))
}

func BenchmarkSimpleSegmentStore(b *testing.B) {
	aggFunc := field.GetAggFunc(field.Sum)
	store := newSimpleFieldStore(aggFunc)
	ss, _ := store.(*simpleFieldStore)

	bs := newBlockStore(30)
	ss.writeInt(bs, 10, int64(100))
	// memory auto rollup
	ss.writeInt(bs, 11, int64(110))
	// memory auto rollup
	ss.writeInt(bs, 10, int64(100))
	// compact because slot out of current time window
	ss.writeInt(bs, 40, int64(20))
	// compact before time window
	ss.writeInt(bs, 10, int64(100))
	// compact because slot out of current time window
	ss.writeInt(bs, 41, int64(50))

	store.bytes()
}

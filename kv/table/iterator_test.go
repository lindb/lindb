package table

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMergedIterator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it1 := NewMockIterator(ctrl)
	gomock.InOrder(
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(10)),
		it1.EXPECT().Value().Return([]byte("value10")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(100)),
		it1.EXPECT().Value().Return([]byte("value100")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(1000)),
		it1.EXPECT().Value().Return([]byte("value1000")),
		it1.EXPECT().HasNext().Return(false),
	)

	it2 := NewMockIterator(ctrl)
	gomock.InOrder(
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(1)),
		it2.EXPECT().Value().Return([]byte("value1")),
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(100)),
		it2.EXPECT().Value().Return([]byte("value100")),
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(2000)),
		it2.EXPECT().Value().Return([]byte("value2000")),
		it2.EXPECT().HasNext().Return(false),
	)

	keys := []uint32{1, 10, 100, 100, 1000, 2000}
	expects := map[uint32][]byte{
		uint32(1):    []byte("value1"),
		uint32(10):   []byte("value10"),
		uint32(100):  []byte("value100"),
		uint32(1000): []byte("value1000"),
		uint32(2000): []byte("value2000"),
	}
	// same length
	mergedIt := NewMergedIterator([]Iterator{it1, it2})
	i := 0
	for mergedIt.HasNext() {
		assert.Equal(t, keys[i], mergedIt.Key())
		assert.Equal(t, expects[keys[i]], mergedIt.Value())
		i++
	}
	assert.Equal(t, len(keys), i)

	gomock.InOrder(
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(10)),
		it1.EXPECT().Value().Return([]byte("value10")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(100)),
		it1.EXPECT().Value().Return([]byte("value100")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(1000)),
		it1.EXPECT().Value().Return([]byte("value1000")),
		it1.EXPECT().HasNext().Return(false),
	)
	// only one iterator has value
	it2.EXPECT().HasNext().Return(false)
	mergedIt = NewMergedIterator([]Iterator{it1, it2})
	keys = []uint32{10, 100, 1000}
	i = 0
	for mergedIt.HasNext() {
		assert.Equal(t, keys[i], mergedIt.Key())
		assert.Equal(t, expects[keys[i]], mergedIt.Value())
		i++
	}
	assert.Equal(t, len(keys), i)

	it1.EXPECT().HasNext().Return(false)
	it2.EXPECT().HasNext().Return(false)
	// both its are empty
	mergedIt = NewMergedIterator([]Iterator{it1, it2})
	assert.False(t, mergedIt.HasNext())

	gomock.InOrder(
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(10)),
		it1.EXPECT().Value().Return([]byte("value10")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(100)),
		it1.EXPECT().Value().Return([]byte("value100")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(1000)),
		it1.EXPECT().Value().Return([]byte("value1000")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(2000)),
		it1.EXPECT().Value().Return([]byte("value2000")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(3000)),
		it1.EXPECT().Value().Return([]byte("value3000")),
		it1.EXPECT().HasNext().Return(true),
		it1.EXPECT().Key().Return(uint32(5000)),
		it1.EXPECT().Value().Return([]byte("value5000")),
		it1.EXPECT().HasNext().Return(false),
	)

	gomock.InOrder(
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(1)),
		it2.EXPECT().Value().Return([]byte("value1")),
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(100)),
		it2.EXPECT().Value().Return([]byte("value100")),
		it2.EXPECT().HasNext().Return(true),
		it2.EXPECT().Key().Return(uint32(2000)),
		it2.EXPECT().Value().Return([]byte("value2000")),
		it2.EXPECT().HasNext().Return(false),
	)
	expects[uint32(3000)] = []byte("value3000")
	expects[uint32(5000)] = []byte("value5000")

	keys = []uint32{1, 10, 100, 100, 1000, 2000, 2000, 3000, 5000}
	// it's length is diff
	mergedIt = NewMergedIterator([]Iterator{it1, it2})
	i = 0
	for mergedIt.HasNext() {
		assert.Equal(t, keys[i], mergedIt.Key())
		assert.Equal(t, expects[keys[i]], mergedIt.Value())
		i++
	}
	assert.Equal(t, len(keys), i)
}

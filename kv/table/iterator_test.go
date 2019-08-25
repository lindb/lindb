package table

import (
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMergedIterator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it1 := generateIterator(ctrl, map[uint32][]byte{
		10:   []byte("value10"),
		100:  []byte("value100"),
		1000: []byte("value1000"),
	})
	it2 := generateIterator(ctrl, map[uint32][]byte{
		1:    []byte("value1"),
		100:  []byte("value100"),
		2000: []byte("value2000"),
	})

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

	it1 = generateIterator(ctrl, map[uint32][]byte{
		10:   []byte("value10"),
		100:  []byte("value100"),
		1000: []byte("value1000"),
	})
	it2 = NewMockIterator(ctrl)
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

	it1 = generateIterator(ctrl, map[uint32][]byte{
		10:   []byte("value10"),
		100:  []byte("value100"),
		1000: []byte("value1000"),
		2000: []byte("value2000"),
		3000: []byte("value3000"),
		5000: []byte("value5000"),
	})
	it2 = generateIterator(ctrl, map[uint32][]byte{
		1:    []byte("value1"),
		100:  []byte("value100"),
		2000: []byte("value2000"),
	})

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

func TestMergedIterator_complex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it1 := generateIterator(ctrl, map[uint32][]byte{
		1:  []byte("value1"),
		3:  []byte("value3"),
		10: []byte("value10"),
	})
	it2 := generateIterator(ctrl, map[uint32][]byte{
		10: []byte("value10"),
		30: []byte("value30"),
		40: []byte("value40"),
	})
	it3 := generateIterator(ctrl, map[uint32][]byte{
		1:  []byte("value1"),
		10: []byte("value10"),
	})
	it4 := generateIterator(ctrl, map[uint32][]byte{
		10:  []byte("value10"),
		30:  []byte("value30"),
		100: []byte("value100"),
	})
	expects := map[uint32][]byte{
		uint32(1):   []byte("value1"),
		uint32(3):   []byte("value3"),
		uint32(10):  []byte("value10"),
		uint32(30):  []byte("value30"),
		uint32(40):  []byte("value40"),
		uint32(100): []byte("value100"),
	}
	keys := []uint32{1, 1, 3, 10, 10, 10, 10, 30, 30, 40, 100}
	mergedIt := NewMergedIterator([]Iterator{it1, it2, it3, it4})
	i := 0
	for mergedIt.HasNext() {
		assert.Equal(t, keys[i], mergedIt.Key())
		assert.Equal(t, expects[keys[i]], mergedIt.Value())
		i++
	}
	assert.Equal(t, len(keys), i)
}

func generateIterator(ctrl *gomock.Controller, values map[uint32][]byte) *MockIterator {
	it1 := NewMockIterator(ctrl)
	var keys []uint32
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var calls []*gomock.Call
	for _, key := range keys {
		calls = append(calls,
			it1.EXPECT().HasNext().Return(true),
			it1.EXPECT().Key().Return(key),
			it1.EXPECT().Value().Return(values[key]))
	}
	calls = append(calls, it1.EXPECT().HasNext().Return(false))

	gomock.InOrder(calls...)
	return it1
}

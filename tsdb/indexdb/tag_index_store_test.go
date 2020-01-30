package indexdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// hack test
func _assertTagIndexStoreData(t *testing.T, tagKeyIDs []uint32, m *tagIndexStore) {
	for _, tagKeyID := range tagKeyIDs {
		found, highIdx, lowIdx := m.tagKeyIDs.ContainsAndRank(tagKeyID)
		assert.True(t, found)
		assert.NotNil(t, m.indexes[highIdx][lowIdx-1])
	}
}

func TestTagIndexStore_put(t *testing.T) {
	m := newTagIndexStore()
	m.put(1, newTagIndex())
	m.put(8, newTagIndex())
	m.put(3, newTagIndex())
	m.put(5, newTagIndex())
	m.put(6, newTagIndex())
	m.put(7, newTagIndex())
	m.put(4, newTagIndex())
	m.put(2, newTagIndex())
	// test insert new high
	m.put(2000000, newTagIndex())
	// test insert new high
	m.put(200000, newTagIndex())

	_assertTagIndexStoreData(t, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000}, m)
}

func TestTagIndexStore_gut(t *testing.T) {
	m := newTagIndexStore()
	store, ok := m.get(uint32(10))
	assert.Nil(t, store)
	assert.False(t, ok)
	m.put(1, newTagIndex())
	m.put(8, newTagIndex())
	_, ok = m.get(1)
	assert.True(t, ok)
	_, ok = m.get(2)
	assert.False(t, ok)
	_, ok = m.get(0)
	assert.False(t, ok)
	_, ok = m.get(9)
	assert.False(t, ok)
}

func TestTagIndexStore_getAllTagKeyIDs(t *testing.T) {
	m := newTagIndexStore()
	m.put(1, newTagIndex())
	m.put(8, newTagIndex())
	assert.Equal(t, m.tagKeyIDs, m.getAllTagKeyIDs())
}

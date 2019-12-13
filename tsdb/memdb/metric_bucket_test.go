package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// hack test
func _assertSortedOrderBucket(t *testing.T, m *metricBucket) {
	for highIndex, mStores := range m.stores {
		for lowIndex, mStore := range mStores {
			metricID := mStore.GetMetricID()
			found, highIdx, lowIdx := m.metricsIDs.ContainsAndRank(metricID)
			assert.True(t, found)
			assert.Equal(t, highIdx, highIndex)
			assert.Equal(t, lowIndex, lowIdx-1)
		}
	}
}
func _newTestMStore(metricID uint32) mStoreINTF {
	mStore := newMetricStore(metricID)
	return mStore
}

func Test_metricBucket_put(t *testing.T) {
	m := newMetricBucket()
	m.put(1, _newTestMStore(1))
	m.put(8, _newTestMStore(8))
	m.put(3, _newTestMStore(3))
	m.put(5, _newTestMStore(5))
	m.put(6, _newTestMStore(6))
	m.put(7, _newTestMStore(7))
	m.put(4, _newTestMStore(4))
	m.put(2, _newTestMStore(2))

	m.put(200000, _newTestMStore(200000))

	_assertSortedOrderBucket(t, m)
}

func Test_metricBucket_gut(t *testing.T) {
	m := newMetricBucket()
	store, ok := m.get(uint32(10))
	assert.Nil(t, store)
	assert.False(t, ok)
	m.put(1, _newTestMStore(1))
	m.put(8, _newTestMStore(8))
	_, ok = m.get(1)
	assert.True(t, ok)
	_, ok = m.get(2)
	assert.False(t, ok)
	_, ok = m.get(0)
	assert.False(t, ok)
	_, ok = m.get(9)
	assert.False(t, ok)
}

func Test_metricBucket_delete(t *testing.T) {
	m := newMetricBucket()
	m.put(1, _newTestMStore(1))
	m.put(8, _newTestMStore(8))
	m.put(3, _newTestMStore(3))
	m.put(5, _newTestMStore(5))
	m.put(6, _newTestMStore(6))
	m.put(7, _newTestMStore(7))
	m.put(4, _newTestMStore(4))
	m.put(2, _newTestMStore(2))

	_assertSortedOrderBucket(t, m)

	m.delete(0)
	m.delete(2)
	_assertSortedOrderBucket(t, m)

	m.delete(1)
	m.delete(10)
	_assertSortedOrderBucket(t, m)

	m.delete(8)
	_assertSortedOrderBucket(t, m)
	assert.Equal(t, m.size(), int(m.metricsIDs.GetCardinality()))

	for i := 0; i < 10; i++ {
		m.delete(uint32(i))
	}
	assert.Len(t, m.stores, 0)
}

func Test_metricBucket_getAllMetricIDs(t *testing.T) {
	m := newMetricBucket()
	m.put(1, _newTestMStore(1))
	m.put(8, _newTestMStore(8))
	assert.Equal(t, m.metricsIDs, m.getAllMetricIDs())
}

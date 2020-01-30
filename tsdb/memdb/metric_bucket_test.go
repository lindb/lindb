package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// hack test
func _assertSortedOrderBucket(t *testing.T, metricIDs []uint32, m *metricBucket) {
	for _, metricID := range metricIDs {
		found, highIdx, lowIdx := m.metricsIDs.ContainsAndRank(metricID)
		assert.True(t, found)
		assert.NotNil(t, m.stores[highIdx][lowIdx-1])
	}
}
func _newTestMStore() mStoreINTF {
	mStore := newMetricStore()
	return mStore
}

func Test_metricBucket_put(t *testing.T) {
	m := newMetricBucket()
	m.put(1, _newTestMStore())
	m.put(8, _newTestMStore())
	m.put(3, _newTestMStore())
	m.put(5, _newTestMStore())
	m.put(6, _newTestMStore())
	m.put(7, _newTestMStore())
	m.put(4, _newTestMStore())
	m.put(2, _newTestMStore())

	m.put(200000, _newTestMStore())

	_assertSortedOrderBucket(t, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000}, m)
}

func Test_metricBucket_gut(t *testing.T) {
	m := newMetricBucket()
	store, ok := m.get(uint32(10))
	assert.Nil(t, store)
	assert.False(t, ok)
	m.put(1, _newTestMStore())
	m.put(8, _newTestMStore())
	_, ok = m.get(1)
	assert.True(t, ok)
	_, ok = m.get(2)
	assert.False(t, ok)
	_, ok = m.get(0)
	assert.False(t, ok)
	_, ok = m.get(9)
	assert.False(t, ok)
}

func Test_metricBucket_getAllMetricIDs(t *testing.T) {
	m := newMetricBucket()
	m.put(1, _newTestMStore())
	m.put(8, _newTestMStore())
	assert.Equal(t, m.metricsIDs, m.getAllMetricIDs())
}

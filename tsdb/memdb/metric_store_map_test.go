package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
)

// hack test
func _assertSortedOrder(t *testing.T, m *metricMap) {
	for highIndex, tStores := range m.stores {
		for lowIndex, tStore := range tStores {
			seriesID := tStore.(*timeSeriesStore).lastWroteTime.Load()
			found, highIdx, lowIdx := m.seriesIDs.ContainsAndRank(seriesID)
			assert.True(t, found)
			assert.Equal(t, highIdx, highIndex)
			assert.Equal(t, lowIndex, lowIdx-1)
		}
	}
}

func _newTestTStore(seriesID uint32) tStoreINTF {
	tStore := newTimeSeriesStore()
	tStore.(*timeSeriesStore).lastWroteTime.Store(seriesID)
	return tStore
}

func Test_metricMap_put(t *testing.T) {
	m := newMetricMap()
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))
	m.put(3, _newTestTStore(3))
	m.put(5, _newTestTStore(5))
	m.put(6, _newTestTStore(6))
	m.put(7, _newTestTStore(7))
	m.put(4, _newTestTStore(4))
	m.put(2, _newTestTStore(2))

	assert.Equal(t, uint64(8), m.getAllSeriesIDs().GetCardinality())

	_assertSortedOrder(t, m)
}

func Test_metricMap_get(t *testing.T) {
	m := newMetricMap()
	store, ok := m.get(uint32(10))
	assert.Nil(t, store)
	assert.False(t, ok)
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))
	_, ok = m.get(1)
	assert.True(t, ok)
	_, ok = m.get(2)
	assert.False(t, ok)
	_, ok = m.get(0)
	assert.False(t, ok)
	_, ok = m.get(9)
	assert.False(t, ok)

	s := m.getAtIndex(0, 0)
	assert.NotNil(t, s)
}

func Test_metricMap_loadData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := newMetricMap()
	for i := 100; i < 4199; i++ {
		m.put(uint32(i), _newTestTStore(uint32(i)))
	}
	queryFlow := flow.NewMockStorageQueryFlow(ctrl)
	gomock.InOrder(
		queryFlow.EXPECT().GetAggregator().Return(nil),
		queryFlow.EXPECT().Reduce("1.1.1.1", gomock.Any()),
	)
	m.loadData(queryFlow, nil, 0, map[string][]uint16{"1.1.1.1": {1, 2, 3, 4}})

	// high key not exist
	m.loadData(queryFlow, nil, 1, map[string][]uint16{"1.1.1.1": {1, 2, 3, 4}})

	gomock.InOrder(
		queryFlow.EXPECT().GetAggregator().Return(nil),
		queryFlow.EXPECT().Reduce("1.1.1.1", gomock.Any()),
	)
	m.loadData(queryFlow, nil, 0, map[string][]uint16{"1.1.1.1": {100, 101}})
}

func TestMetricStore_Filter(t *testing.T) {
	m := newMetricMap()
	for i := 100; i < 4199; i++ {
		m.put(uint32(i), _newTestTStore(uint32(i)))
	}
	assert.False(t, m.filter(roaring.BitmapOf(1, 3)))
	assert.True(t, m.filter(roaring.BitmapOf(100, 40000)))
}

func Benchmark_get(b *testing.B) {
	m := newMetricMap()
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))

	for i := 0; i < b.N; i++ {
		_, _ = m.get(8)
	}
}

func Benchmark_iterate_bitmap(b *testing.B) {
	r := roaring.New()
	r.AddRange(1, 100000)
	s := make([]uint32, r.GetCardinality())

	for i := 0; i < b.N; i++ {
		itr := r.ManyIterator()
		length := itr.NextMany(s)
		for x := 0; x < length; x++ {
			_ = s[x]
		}
	}
}

func Benchmark_copy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int64, 10000)
		for y := 0; y < 10000; y++ {
			s[y] = 0
		}
	}
}

package memdb

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
)

// hack test
func _assertSortedOrder(t *testing.T, m *metricMap) {
	for idx, tStore := range m.stores {
		seriesID := tStore.(*timeSeriesStore).lastWroteTime.Load()
		assert.True(t, m.seriesIDs.Contains(seriesID))
		expectedIdx := m.seriesIDs.Rank(seriesID) - 1
		assert.Equal(t, idx, int(expectedIdx))
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
}

func Test_metricMap_iterator(t *testing.T) {
	m := newMetricMap()
	it := m.iterator()
	assert.False(t, it.hasNext())
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))
	it = m.iterator()
	assert.True(t, it.hasNext())
	seriesID, tStore := it.next()
	assert.Equal(t, uint32(1), seriesID)
	assert.NotNil(t, tStore)
	assert.True(t, it.hasNext())
	seriesID, tStore = it.next()
	assert.Equal(t, uint32(8), seriesID)
	assert.NotNil(t, tStore)
	assert.False(t, it.hasNext())
}

func Test_metricMap_scan(t *testing.T) {
	m := newMetricMap()
	for i := 100; i < 4199; i++ {
		m.put(uint32(i), _newTestTStore(uint32(i)))
	}
	foundSeriesIDs := roaring.BitmapOf()
	multiVer1 := series.NewMultiVerSeriesIDSet()
	multiVer1.Add(series.Version(12), roaring.BitmapOf())
	worker := &mockScanWorker{}
	mCtx := &series.ScanContext{
		SeriesIDSet: multiVer1,
		HasGroupBy:  true,
		Worker:      worker,
	}
	// not match
	m.scan(series.Version(12), mCtx)
	assert.Equal(t, 0, len(worker.events))

	// find all series ids
	multiVer1.Add(series.Version(13), m.seriesIDs.Clone())
	m.scan(series.Version(13), mCtx)
	assert.Equal(t, 2, len(worker.events))
	foundSeriesIDs.Or(worker.events[0].SeriesIDs())
	foundSeriesIDs.Or(worker.events[1].SeriesIDs())
	assert.True(t, foundSeriesIDs.Equals(m.seriesIDs))

	// find some series ids
	seriesIDs := m.seriesIDs.Clone()
	seriesIDs.Remove(uint32(300))
	multiVer1.Add(series.Version(14), seriesIDs)
	worker.events = nil
	m.scan(series.Version(14), mCtx)
	assert.Equal(t, 2, len(worker.events))
	foundSeriesIDs.Clear()
	foundSeriesIDs.Or(worker.events[0].SeriesIDs())
	foundSeriesIDs.Or(worker.events[1].SeriesIDs())
	seriesIDs = m.seriesIDs.Clone()
	seriesIDs.Remove(uint32(300))
	assert.True(t, foundSeriesIDs.Equals(seriesIDs))

	seriesIDs = roaring.New()
	seriesIDs.AddRange(uint64(30), uint64(10240))
	multiVer1.Add(series.Version(15), seriesIDs)
	worker.events = nil
	m.scan(series.Version(15), mCtx)
	assert.Equal(t, 2, len(worker.events))
	foundSeriesIDs.Clear()
	foundSeriesIDs.Or(worker.events[0].SeriesIDs())
	foundSeriesIDs.Or(worker.events[1].SeriesIDs())
	seriesIDs = roaring.New()
	seriesIDs.AddRange(uint64(100), uint64(4199))
	assert.True(t, foundSeriesIDs.Equals(seriesIDs))
}

func Benchmark_get(b *testing.B) {
	m := newMetricMap()
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))

	for i := 0; i < b.N; i++ {
		_, _ = m.get(8)
	}
}

func Test_metricMap_delete(t *testing.T) {
	m := newMetricMap()
	m.put(1, _newTestTStore(1))
	m.put(8, _newTestTStore(8))
	m.put(3, _newTestTStore(3))
	m.put(5, _newTestTStore(5))
	m.put(6, _newTestTStore(6))
	m.put(7, _newTestTStore(7))
	m.put(4, _newTestTStore(4))
	m.put(2, _newTestTStore(2))

	_assertSortedOrder(t, m)

	m.delete(0)
	m.delete(2)
	_assertSortedOrder(t, m)

	m.delete(1)
	m.delete(10)
	_assertSortedOrder(t, m)

	m.delete(8)
	_assertSortedOrder(t, m)
	assert.Equal(t, m.size(), int(m.seriesIDs.GetCardinality()))
}

func Test_metricMap_deleteMany(t *testing.T) {
	m := newMetricMap()
	for i := uint32(1); i <= 100000; i++ {
		m.put(i, _newTestTStore(i))
	}
	var seriesIDs []uint32
	for i := uint32(1); i < 5000; i += 2 {
		seriesIDs = append(seriesIDs, i)
	}
	assert.Len(t, m.deleteMany(seriesIDs...), 2500)
	assert.Len(t, m.stores, 100000-2500)
	assert.Equal(t, 100000-2500, int(m.seriesIDs.GetCardinality()))
	_assertSortedOrder(t, m)

	m.deleteMany()
	assert.Equal(t, 100000-2500, int(m.seriesIDs.GetCardinality()))

	m.deleteMany(0, 100001, 100002, 100003)
	assert.Len(t, m.stores, 100000-2500)
	assert.Equal(t, 100000-2500, int(m.seriesIDs.GetCardinality()))
	_assertSortedOrder(t, m)
}

func Benchmark_deleteMany(b *testing.B) {
	var seriesIDs []uint32
	for i := uint32(1); i < 5000; i += 2 {
		seriesIDs = append(seriesIDs, i)
	}

	for x := 0; x < b.N; x++ {
		b.StopTimer()
		m := newMetricMap()
		for i := uint32(1); i <= 100000; i++ {
			m.put(i, _newTestTStore(i))
		}
		b.StartTimer()
		m.deleteMany(seriesIDs...)
	}
}

func Benchmark_delete(b *testing.B) {
	for x := 0; x < b.N; x++ {
		b.StopTimer()
		m := newMetricMap()
		for i := uint32(1); i <= 100000; i++ {
			m.put(i, _newTestTStore(i))
		}
		b.StartTimer()
		for i := uint32(1); i < 5000; i += 2 {
			m.delete(i)
		}
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

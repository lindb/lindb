package memdb

import (
	"github.com/RoaringBitmap/roaring"
)

// metricMap represents a map structure for storing time series data.
// keys => bitmap, values => slice
type metricMap struct {
	seriesIDs *roaring.Bitmap
	stores    []tStoreINTF
}

// newMetricMap creates a metric map
func newMetricMap() *metricMap {
	return &metricMap{
		seriesIDs: roaring.New(),
	}
}

// get returns time series store by series id, if exist returns it, else returns nil, false
func (m *metricMap) get(seriesID uint32) (tStoreINTF, bool) {
	if len(m.stores) == 0 {
		return nil, false
	}
	idx := m.seriesIDs.Rank(seriesID)
	if idx == 0 {
		// not found
		return nil, false
	}
	tStore := m.stores[idx-1]
	if tStore.GetSeriesID() == seriesID {
		return tStore, true
	}
	return nil, false
}

// put puts the time series store
func (m *metricMap) put(seriesID uint32, tStore tStoreINTF) {
	if len(m.stores) == 0 {
		m.seriesIDs.Add(seriesID)
		m.stores = append(m.stores, tStore)
		return
	}
	idx := m.seriesIDs.Rank(seriesID)
	// add new value
	m.seriesIDs.Add(seriesID)
	m.stores = append(m.stores, nil)
	copy(m.stores[idx+1:], m.stores[idx:len(m.stores)-1])
	m.stores[idx] = tStore
}

// delete deletes the time series store by series id
func (m *metricMap) delete(seriesID uint32) {
	idx := m.seriesIDs.Rank(seriesID)
	if idx == 0 {
		// not found
		return
	}

	tStore := m.stores[idx-1]
	if tStore.GetSeriesID() == seriesID {
		m.seriesIDs.Remove(seriesID)
		copy(m.stores[idx-1:], m.stores[idx:])
		m.stores = m.stores[:len(m.stores)-1]
	}
}

// size returns the size of map
func (m *metricMap) size() int {
	return len(m.stores)
}

// iterator returns an iterator for iterating the map data
func (m *metricMap) iterator() *mStoreIterator {
	return &mStoreIterator{
		it:     m.seriesIDs.Iterator(),
		stores: m.stores,
	}
}

// mStoreIterator represents an iterator over the metric map
type mStoreIterator struct {
	it     roaring.IntIterable
	stores []tStoreINTF

	idx int
}

// hasNext returns if the iteration has more time series store
func (it *mStoreIterator) hasNext() bool {
	return it.it.HasNext()
}

// next returns the series id and store
func (it *mStoreIterator) next() (seriesID uint32, store tStoreINTF) {
	seriesID = it.it.Next()
	store = it.stores[it.idx]
	it.idx++
	return
}

package memdb

import (
	"sort"

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
	idx := int(m.seriesIDs.Rank(seriesID))
	if !m.seriesIDs.Contains(seriesID) {
		return nil, false
	}
	return m.stores[idx-1], true
}

// put puts the time series store
func (m *metricMap) put(seriesID uint32, tStore tStoreINTF) {
	if len(m.stores) == 0 {
		m.seriesIDs.Add(seriesID)
		m.stores = append(m.stores, tStore)
		return
	}
	idx := int(m.seriesIDs.Rank(seriesID))
	m.seriesIDs.Add(seriesID)
	// append operation
	if idx == len(m.stores) {
		m.stores = append(m.stores, tStore)
		return
	}
	// insert operation
	m.stores = append(m.stores, nil)
	copy(m.stores[idx+1:], m.stores[idx:len(m.stores)-1])
	m.stores[idx] = tStore
}

// delete deletes the time series store by series id
func (m *metricMap) delete(seriesID uint32) tStoreINTF {
	if !m.seriesIDs.Contains(seriesID) {
		return nil
	}
	idx := m.seriesIDs.Rank(seriesID)
	tStore := m.stores[idx-1]
	m.seriesIDs.Remove(seriesID)
	copy(m.stores[idx-1:], m.stores[idx:])
	m.stores[len(m.stores)-1] = nil
	m.stores = m.stores[:len(m.stores)-1]
	return tStore
}

// deleteMany deletes the time series store by multi seriesIDs
func (m *metricMap) deleteMany(
	seriesIDs ...uint32,
) (
	removedTStores []tStoreINTF,
) {
	if len(seriesIDs) == 0 {
		return nil
	}
	sort.Slice(seriesIDs, func(i, j int) bool {
		return seriesIDs[i] < seriesIDs[j]
	})
	removedTStores = make([]tStoreINTF, len(seriesIDs))[:0]
	var (
		nextRemoveIndex = 0
		manyItrLen      = 0
		buffer          = make([]uint32, 4096)
		exhausted       = false
		n               = 0
		counter         = 0
	)
	keep := func(seriesID uint32) bool {
		if exhausted {
			return true
		}
		if int(seriesID) > int(seriesIDs[nextRemoveIndex]) {
			nextRemoveIndex++
			if nextRemoveIndex >= len(seriesIDs) {
				exhausted = true
				return true
			}
		}
		if int(seriesID) == int(seriesIDs[nextRemoveIndex]) {
			return false
		}
		return true
	}
	manyItr := m.seriesIDs.ManyIterator()
	for {
		manyItrLen = manyItr.NextMany(buffer)
		if manyItrLen == 0 {
			break
		}
		for idx := 0; idx < manyItrLen; idx++ {
			thisSeriesID := buffer[idx]
			if keep(thisSeriesID) {
				m.stores[n] = m.stores[counter]
				n++
			} else {
				removedTStores = append(removedTStores, m.stores[counter])
			}
			counter++
		}
	}

	for idx := n; idx < counter; idx++ {
		m.stores[idx] = nil
	}
	for _, seriesID := range seriesIDs {
		m.seriesIDs.Remove(seriesID)
	}
	m.stores = m.stores[:n]
	return removedTStores
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

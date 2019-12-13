package memdb

import (
	"sync"

	"github.com/lindb/roaring"
)

// metricBucket represents the metric store bucket
type metricBucket struct {
	metricsIDs *roaring.Bitmap
	stores     [][]mStoreINTF
	rwLock     sync.RWMutex // read-write lock of hash2MStore
}

// newMetricBucket creates new metric bucket for storing metric store
func newMetricBucket() *metricBucket {
	return &metricBucket{
		metricsIDs: roaring.New(),
	}
}

// get returns metric store by metric id, if exist returns it, else returns nil, false
func (m *metricBucket) get(metricID uint32) (mStoreINTF, bool) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()

	if len(m.stores) == 0 {
		return nil, false
	}
	found, highIdx, lowIdx := m.metricsIDs.ContainsAndRank(metricID)
	if !found {
		return nil, false
	}
	return m.stores[highIdx][lowIdx-1], true
}

// put puts the metric store by metric id
func (m *metricBucket) put(metricID uint32, mStore mStoreINTF) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	if len(m.stores) == 0 {
		m.metricsIDs.Add(metricID)
		m.stores = append(m.stores, []mStoreINTF{mStore})
		return
	}

	found, highIdx, lowIdx := m.metricsIDs.ContainsAndRank(metricID)
	if !found {
		m.metricsIDs.Add(metricID)
		if highIdx < 0 {
			// high container not exist, append operation
			m.stores = append(m.stores, []mStoreINTF{mStore})
		} else {
			// high container exist
			stores := m.stores[highIdx]
			// insert operation
			stores = append(stores, nil)
			copy(stores[lowIdx+1:], stores[lowIdx:len(stores)-1])
			stores[lowIdx] = mStore
			m.stores[highIdx] = stores
		}
	}
}

// delete deletes the metric store by metric id
func (m *metricBucket) delete(metricID uint32) mStoreINTF {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	found, highIdx, lowIdx := m.metricsIDs.ContainsAndRank(metricID)
	if !found {
		return nil
	}
	// get high container
	stores := m.stores[highIdx]
	// get mStore
	mStore := stores[lowIdx-1]
	// remove metric id
	m.metricsIDs.Remove(metricID)

	if len(stores) > 1 {
		// remove mStore from high container
		copy(stores[lowIdx-1:], stores[lowIdx:])
		stores[len(stores)-1] = nil
		// reset high container
		m.stores[highIdx] = stores[:len(stores)-1]
	} else {
		// remove high container
		copy(m.stores[highIdx:], m.stores[highIdx+1:])
		m.stores = m.stores[:len(m.stores)-1]
	}
	return mStore
}

// size returns the size of metric store
func (m *metricBucket) size() int {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()

	return int(m.metricsIDs.GetCardinality())
}

// getAllMetricIDs gets all metric ids
func (m *metricBucket) getAllMetricIDs() *roaring.Bitmap {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()

	return m.metricsIDs.Clone()
}

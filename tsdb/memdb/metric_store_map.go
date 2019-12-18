package memdb

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
)

// metricMap represents a map structure for storing time series data.
// keys => bitmap, values => slice
type metricMap struct {
	seriesIDs *roaring.Bitmap
	stores    [][]tStoreINTF
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
	found, highIdx, lowIdx := m.seriesIDs.ContainsAndRank(seriesID)
	if !found {
		return nil, false
	}
	return m.stores[highIdx][lowIdx-1], true
}

// getAtIndex gets by given high/low index, make sure index exist
func (m *metricMap) getAtIndex(highIdx, lowIdx int) tStoreINTF {
	return m.stores[highIdx][lowIdx]
}

// put puts the time series store
func (m *metricMap) put(seriesID uint32, tStore tStoreINTF) {
	if len(m.stores) == 0 {
		m.seriesIDs.Add(seriesID)
		m.stores = append(m.stores, []tStoreINTF{tStore})
		return
	}

	found, highIdx, lowIdx := m.seriesIDs.ContainsAndRank(seriesID)
	if !found {
		m.seriesIDs.Add(seriesID)
		if highIdx < 0 {
			// high container not exist, append operation
			m.stores = append(m.stores, []tStoreINTF{tStore})
		} else {
			// high container exist
			stores := m.stores[highIdx]
			// insert operation
			stores = append(stores, nil)
			copy(stores[lowIdx+1:], stores[lowIdx:len(stores)-1])
			stores[lowIdx] = tStore
			m.stores[highIdx] = stores
		}
	}
}

// delete deletes the time series store by series id
func (m *metricMap) delete(seriesID uint32) tStoreINTF {
	found, highIdx, lowIdx := m.seriesIDs.ContainsAndRank(seriesID)
	if !found {
		return nil
	}
	// get high container
	stores := m.stores[highIdx]
	// get tStore
	tStore := stores[lowIdx-1]
	// remove series id
	m.seriesIDs.Remove(seriesID)

	if len(stores) > 1 {
		// remove tStore from high container
		copy(stores[lowIdx-1:], stores[lowIdx:])
		stores[len(stores)-1] = nil
		// reset high container
		m.stores[highIdx] = stores[:len(stores)-1]
	} else {
		// remove high container
		copy(m.stores[highIdx:], m.stores[highIdx+1:])
		m.stores = m.stores[:len(m.stores)-1]
	}
	return tStore
}

// deleteMany deletes the time series store by multi seriesIDs
func (m *metricMap) deleteMany(
	seriesIDs ...uint32,
) (
	removedTStores []tStoreINTF,
) {
	if len(seriesIDs) == 0 {
		return
	}
	sort.Slice(seriesIDs, func(i, j int) bool {
		return seriesIDs[i] < seriesIDs[j]
	})
	for _, seriesID := range seriesIDs {
		tStore := m.delete(seriesID)
		if tStore != nil {
			removedTStores = append(removedTStores, tStore)
		}
	}
	return
}

// size returns the size of map
func (m *metricMap) size() int {
	return int(m.seriesIDs.GetCardinality())
}

// iterator returns an iterator for iterating the map data
func (m *metricMap) iterator() *mStoreIterator {
	return newStoreIterator(m)
}

// getAllSeriesIDs gets all series ids
func (m *metricMap) getAllSeriesIDs() *roaring.Bitmap {
	return m.seriesIDs.Clone()
}

func (m *metricMap) loadData(flow flow.StorageQueryFlow, fieldIDs []uint16,
	highKey uint16, groupedSeries map[string][]uint16,
) {
	//FIXME need add lock?????

	// 1. get high container index by the high key of series ID
	highContainerIdx := m.seriesIDs.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series ID not exist) return it
		return
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := m.seriesIDs.GetContainerAtIndex(highContainerIdx)

	memScanCtx := &memScanContext{
		fieldIDs:   fieldIDs,
		tsd:        encoding.GetTSDDecoder(),
		fieldCount: len(fieldIDs),
	}
	for groupByTags, lowSeriesIDs := range groupedSeries {
		aggregator := flow.GetAggregator()
		memScanCtx.aggregators = aggregator
		for _, lowSeriesID := range lowSeriesIDs {
			// check low series id if exist
			if !lowContainer.Contains(lowSeriesID) {
				continue
			}
			// get the index of low series id in container
			idx := lowContainer.Rank(lowSeriesID)
			// scan the data and aggregate the values
			store := m.stores[highContainerIdx][idx-1]
			store.scan(memScanCtx)
		}
		flow.Reduce(groupByTags, aggregator)
	}
	encoding.ReleaseTSDDecoder(memScanCtx.tsd)
}

// filter filters if seriesIDSs exist in data storage
func (m *metricMap) filter(seriesIDs *roaring.Bitmap) bool {
	// after and operator, query bitmap is sub of store bitmap
	matchSeriesIDs := roaring.FastAnd(seriesIDs, m.seriesIDs)
	return !matchSeriesIDs.IsEmpty()
}

// mStoreIterator represents an iterator over the metric map
type mStoreIterator struct {
	mStore      *metricMap
	highKeys    []uint16
	highKeysLen int

	highKey                     uint32
	lowIt                       roaring.PeekableShortIterator
	highIdx, lowIdx, curHighIdx int
}

func newStoreIterator(mStore *metricMap) *mStoreIterator {
	highKeys := mStore.seriesIDs.GetHighKeys()
	return &mStoreIterator{
		mStore:      mStore,
		highKeys:    highKeys,
		highKeysLen: len(highKeys),
	}
}

// hasNext returns if the iteration has more time series store
func (it *mStoreIterator) hasNext() bool {
	if it.highKeysLen == 0 {
		return false
	}

	notFound := it.lowIt == nil || !it.lowIt.HasNext()
	if notFound {
		if it.highIdx == it.highKeysLen {
			return false
		}
		it.curHighIdx = it.highIdx
		it.highKey = uint32(it.highKeys[it.curHighIdx]) << 16
		it.lowIt = it.mStore.seriesIDs.GetContainerAtIndex(it.curHighIdx).PeekableIterator()

		// for next loop
		it.highIdx++
		it.lowIdx = 0 // reset low index
	}
	return it.lowIt.HasNext()
}

// next returns the series id and store
func (it *mStoreIterator) next() (seriesID uint32, store tStoreINTF) {
	seriesID = uint32(it.lowIt.Next()) | it.highKey
	store = it.mStore.stores[it.curHighIdx][it.lowIdx]
	it.lowIdx++
	return
}

package memdb

import (
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

// size returns the size of map
func (m *metricMap) size() int {
	return int(m.seriesIDs.GetCardinality())
}

// getAllSeriesIDs gets all series ids
func (m *metricMap) getAllSeriesIDs() *roaring.Bitmap {
	return m.seriesIDs
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

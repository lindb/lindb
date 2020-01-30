package memdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./metric_store_index.go -destination=./metric_store_index_mock.go -package memdb

//FIXME
//const emptyTagIndexSize = 24 + // tagKVEntrySet slice
//	24 + // metric-map slice
//	4 + // idCounter
//	8 // version

// tagIndexINTF abstracts the index of tStores, not thread-safe
type tagIndexINTF interface {
	// GetOrCreateTStore constructs the index and return a tStore,
	// error of too may tag keys may be return
	GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int)
	// FlushVersionDataTo flush metric to the tableFlusher
	FlushVersionDataTo(
		flusher metricsdata.Flusher,
		flushCtx flushContext,
	)
	// Version returns a version(uptime in milliseconds) of the index
	Version() series.Version
	// filter filters if series ids exist in data storage
	filter(seriesIDs *roaring.Bitmap) bool
	// loadData loads time series data based grouped series ids
	loadData(flow flow.StorageQueryFlow, fieldIDs []uint16, highKey uint16, groupedSeries map[string][]uint16)
}

// tagIndex implements tagIndexINTF,
// it is a composition of both inverted and forward index,
// not thread-safe
type tagIndex struct {
	seriesID2TStore *metricMap
	// version is the uptime in milliseconds
	version series.Version
}

// newTagIndex returns a new tagIndexINTF with version.
func newTagIndex() tagIndexINTF {
	return &tagIndex{
		seriesID2TStore: newMetricMap(),
		version:         series.NewVersion(),
	}
}

// GetTStoreBySeriesID returns a tStoreINTF from series-id.
func (index *tagIndex) GetTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool) {
	return index.seriesID2TStore.get(seriesID)
}

// GetOrCreateTStore get or creates the tStore from string tags,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *tagIndex) GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int) {
	// double check
	tStore, ok := index.GetTStoreBySeriesID(seriesID)
	if !ok {
		tStore = newTimeSeriesStore()
		index.seriesID2TStore.put(seriesID, tStore)
		createdSize += tStore.MemSize()
	}
	return tStore, createdSize
}

// AllTStores returns the map of seriesID and tStores
func (index *tagIndex) AllTStores() *metricMap {
	return index.seriesID2TStore
}

// FlushVersionDataTo flushes metric-block of mStore to the writer.
func (index *tagIndex) FlushVersionDataTo(
	tableFlusher metricsdata.Flusher,
	flushCtx flushContext,
) {
	seriesIDs := index.seriesID2TStore.getAllSeriesIDs()
	if seriesIDs.IsEmpty() {
		// if no series data, returns it
		return
	}

	// get the all high key of series ids, flush data by roaring.Bitmap's container
	highKeys := seriesIDs.GetHighKeys()

	for highIdx := range highKeys {
		lowContainer := seriesIDs.GetContainerAtIndex(highIdx)
		it := lowContainer.PeekableIterator()
		lowIdx := 0
		for it.HasNext() {
			_ = it.Next() //skip to next value
			tStore := index.seriesID2TStore.getAtIndex(highIdx, lowIdx)
			tStore.FlushSeriesTo(tableFlusher, flushCtx)
			lowIdx++
		}
		tableFlusher.FlushSeriesBucket()
	}
	tableFlusher.FlushVersion(index.Version(), seriesIDs)
}

// Version returns a version(uptime) of the index
func (index *tagIndex) Version() series.Version {
	return index.version
}

// filter filters if seriesIDSs exist in data storage
func (index *tagIndex) filter(seriesIDs *roaring.Bitmap) bool {
	return index.seriesID2TStore.filter(seriesIDs)
}

// loadData loads the data points data from storage based on high series key and grouped low series ids
func (index *tagIndex) loadData(flow flow.StorageQueryFlow, fieldIDs []uint16,
	highKey uint16, groupedSeries map[string][]uint16,
) {
	index.seriesID2TStore.loadData(flow, fieldIDs, highKey, groupedSeries)
}

// staticNopTagIndex is the static nop-tagIndex,
// it is used as a placeholder of immutable atomic.Value
var staticNopTagIndex = newNopTagIndex()

func newNopTagIndex() tagIndexINTF {
	ti := newTagIndex().(*tagIndex)
	ti.version = 0
	return ti
}

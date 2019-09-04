package memdb

import (
	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/series"
)

//////////////////////////////////////////////////////
// Scanner methods
//////////////////////////////////////////////////////

// retain locks the metricStore, then release a unlock function
func (ms *metricStore) retain() func() {
	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	return func() {
		ms.mutex4Immutable.RUnlock()
		ms.mutex4Mutable.RUnlock()
	}
}

// scan returns a iterator for scanning data
func (ms *metricStore) scan(sCtx series.ScanContext) series.VersionIterator {
	return newMStoreIterator(ms, sCtx)
}

//////////////////////////////////////////////////////
// tStoreIterator implements series.Iterator
//////////////////////////////////////////////////////
type tStoreIterator struct {
	tagIndex    tagIndexINTF
	releaseFunc func()              // release lock handler for tStore
	intItr      roaring.IntIterable // bitmap iterator
	fStoreItr   *fStoreIterator     // reusable
	seriesID    uint32              // current ts id
}

func newTStoreIterator(metas fieldsMetas, sCtx series.ScanContext) *tStoreIterator {
	return &tStoreIterator{fStoreItr: newFStoreIterator(metas, sCtx)}
}

// reset resets the multiTimeSeries to a different tStore
func (tsi *tStoreIterator) reset(tagIndex tagIndexINTF, itr roaring.IntIterable) {
	tsi.tagIndex = tagIndex
	tsi.intItr = itr
}

func (tsi *tStoreIterator) SeriesID() uint32 { return tsi.seriesID }
func (tsi *tStoreIterator) HasNext() bool {
	if tsi.intItr == nil {
		return false
	}
Loop:
	{
		if tsi.releaseFunc != nil {
			// release lock of tStore
			tsi.releaseFunc()
			tsi.releaseFunc = nil
		}
		if !tsi.intItr.HasNext() {
			return false
		}
		tsi.seriesID = tsi.intItr.Next()
		tStore, ok := tsi.tagIndex.getTStoreBySeriesID(tsi.seriesID)
		if !ok {
			tsi.seriesID = 0
			goto Loop
		}
		// hold lock of tStore
		tsi.releaseFunc = tStore.retain()
		tsi.fStoreItr.reset(tStore)
		return true
	}
}

func (tsi *tStoreIterator) Next() series.FieldIterator {
	return tsi.fStoreItr
}

//////////////////////////////////////////////////////
// mStoreIterator implements series.VersionIterator
//////////////////////////////////////////////////////
type mStoreIterator struct {
	releaseFunc func() // release lock handler for mStore
	mStore      *metricStore
	sCtx        series.ScanContext
	tagIndexes  []tagIndexINTF
	version     uint32
	tStoreItr   *tStoreIterator
}

func newMStoreIterator(mStore *metricStore, sCtx series.ScanContext) *mStoreIterator {
	msi := &mStoreIterator{
		mStore: mStore, sCtx: sCtx,
		tStoreItr: newTStoreIterator(mStore.fieldsMetas, sCtx),
	}
	msi.releaseFunc = mStore.retain()

	// collect all tagIndexes whose version matches the idSet
	collectOnVersionMatch := func(idx tagIndexINTF) {
		if _, ok := msi.sCtx.SeriesIDSet.Versions()[idx.getVersion()]; ok {
			msi.tagIndexes = append(msi.tagIndexes, idx)
		}
	}
	for _, idx := range msi.mStore.immutable {
		collectOnVersionMatch(idx)
	}
	collectOnVersionMatch(msi.mStore.mutable)
	return msi
}

func (msi *mStoreIterator) Close() error {
	if msi.releaseFunc != nil {
		msi.releaseFunc()
		msi.releaseFunc = nil
	}
	return nil
}

func (msi *mStoreIterator) Version() uint32 {
	return msi.version
}

func (msi *mStoreIterator) HasNext() bool {
Loop:
	{ // version exhaustion
		if len(msi.tagIndexes) == 0 {
			return false
		}
		thisTagIndex := msi.tagIndexes[0]
		msi.tagIndexes = msi.tagIndexes[1:]
		msi.version = thisTagIndex.getVersion()
		intItr := msi.sCtx.SeriesIDSet.Versions()[msi.version].Iterator()
		msi.tStoreItr.reset(thisTagIndex, intItr)

		startTime, endTime := thisTagIndex.getTimeRange()
		if !msi.sCtx.TimeRange.Overlap(&timeutil.TimeRange{
			Start: int64(startTime) * 1000, End: int64(endTime) * 1000}) {

			msi.tStoreItr.reset(nil, nil)
			msi.version = 0
			goto Loop
		}
		return true
	}
}

func (msi *mStoreIterator) Next() series.Iterator {
	return msi.tStoreItr
}

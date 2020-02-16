package memdb

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock.go -package memdb

// for testing
var (
	flushFunc = flush
)

const emptyMStoreSize = 8 + // mutable
	8 + // atomic.Value
	4 + // uint32
	4 + // uint32
	4 // int32

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// Filter filters the data based on fieldIDs/seriesIDs/familyIDs,
	// if finds data then returns the FilterResultSet, else returns constants.ErrNotFound
	Filter(fieldIDs []field.ID,
		seriesIDs *roaring.Bitmap, familyIDs map[familyID]int64,
	) ([]flow.FilterResultSet, error)
	// SetTimestamp sets the current write timestamp
	SetTimestamp(familyID familyID, slot uint16)
	// AddField adds field meta into metric level
	AddField(fieldID field.ID, fieldType field.Type)
	// GetOrCreateTStore constructs the index and return a tStore
	GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int)
	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(tableFlusher metricsdata.Flusher, flushCtx flushContext) (err error)
}

type familyIDSlotRangeEntry struct {
	id        familyID
	slotRange familySlotRange
}

// familyIDSlotRangeEntries implements sort.Interface
// sorted in ascending order
type familyIDSlotRangeEntries []familyIDSlotRangeEntry

func (e familyIDSlotRangeEntries) Len() int           { return len(e) }
func (e familyIDSlotRangeEntries) Less(i, j int) bool { return e[i].id < e[j].id }
func (e familyIDSlotRangeEntries) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e familyIDSlotRangeEntries) GetRange(id familyID) (familySlotRange, bool) {
	idx := sort.Search(e.Len(), func(i int) bool {
		return e[i].id >= id
	})
	if idx >= e.Len() || e[idx].id != id {
		return familySlotRange{}, false
	}
	return e[idx].slotRange, true
}

func (e familyIDSlotRangeEntries) SetRange(id familyID, slotRange familySlotRange) familyIDSlotRangeEntries {
	idx := sort.Search(e.Len(), func(i int) bool {
		return e[i].id >= id
	})
	if idx >= e.Len() || e[idx].id != id {
		newE := append(e, familyIDSlotRangeEntry{id: id, slotRange: slotRange})
		sort.Sort(newE)
		return newE
	}
	e[idx].slotRange = slotRange
	return e
}

// metricStore represents metric level storage, stores all series data, and fields/family times metadata
type metricStore struct {
	MetricStore

	families familyIDSlotRangeEntries // time slot range
	fields   field.Metas              // field metadata
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	var ms metricStore
	ms.keys = roaring.New() // init keys
	return &ms
}

// SetTimestamp sets the current write timestamp
func (ms *metricStore) SetTimestamp(familyID familyID, slot uint16) {
	slotRange, ok := ms.families.GetRange(familyID)
	if ok {
		slotRange.setSlot(slot)
	} else {
		slotRange = newFamilySlotRange(slot, slot)
	}
	ms.families = ms.families.SetRange(familyID, slotRange)
}

// AddField adds field meta into metric level
func (ms *metricStore) AddField(fieldID field.ID, fieldType field.Type) {
	_, ok := ms.fields.GetFromID(fieldID)
	if !ok {
		ms.fields = ms.fields.Insert(field.Meta{
			ID:   fieldID,
			Type: fieldType,
		})
		// sort by field id
		sort.Slice(ms.fields, func(i, j int) bool { return ms.fields[i].ID < ms.fields[j].ID })
	}
}

// GetOrCreateTStore constructs the index and return a tStore
func (ms *metricStore) GetOrCreateTStore(seriesID uint32) (tStore tStoreINTF, createdSize int) {
	tStore, ok := ms.Get(seriesID)
	if !ok {
		tStore = newTimeSeriesStore()
		ms.Put(seriesID, tStore)
		createdSize += emptyTimeSeriesStoreSize
	}
	return tStore, createdSize
}

// FlushMetricsTo Writes metric-data to the table.
func (ms *metricStore) FlushMetricsDataTo(flusher metricsdata.Flusher, flushCtx flushContext) (err error) {
	// family time not exist, return
	slotRange, ok := ms.families.GetRange(flushCtx.familyID)
	if !ok {
		return
	}
	// field not exist, return
	fieldLen := len(ms.fields)
	if fieldLen == 0 {
		return
	}
	// flush field meta info
	flusher.FlushFieldMetas(ms.fields)
	// set current family's slot range
	flushCtx.start, flushCtx.end = slotRange.getSlotRange()
	if err := ms.WalkEntry(func(key uint32, value tStoreINTF) error {
		return flushFunc(flusher, flushCtx, key, value)
	}); err != nil {
		return err
	}
	return flusher.FlushMetric(flushCtx.metricID, slotRange.start, slotRange.end)
}

// flush flushes series data
func flush(flusher metricsdata.Flusher, flushCtx flushContext, key uint32, value tStoreINTF) error {
	value.FlushSeriesTo(flusher, flushCtx)
	flusher.FlushSeries(key)
	return nil
}

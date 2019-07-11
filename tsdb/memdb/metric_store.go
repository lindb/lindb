package memdb

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eleme/lindb/pkg/hashers"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

// versionedTSMap holds a mapping relation of tags and TsStore.
// a version is assigned since newed,
type versionedTSMap struct {
	tsMap       map[uint32]*timeSeriesStore // map-key: FNV32a(sortedTag)
	familyTimes map[int64]struct{}          // all segments
	version     int64                       // uptime in nanoseconds
}

// allTSStores returns a tsStore list in order of tsID.
func (vm *versionedTSMap) allTSStores() (buf *[]*timeSeriesStore, release func()) {
	buf = tsStoresListPool.get(len(vm.tsMap))
	var count = 0
	for _, tsStore := range vm.tsMap {
		(*buf)[count] = tsStore
		count++
	}
	return buf, func() {
		tsStoresListPool.put(buf)
	}
}

// unionFamilyTimesTo update familyTimes to the input map.
func (vm *versionedTSMap) unionFamilyTimesTo(segments map[int64]struct{}) {
	for familyTime := range vm.familyTimes {
		segments[familyTime] = struct{}{}
	}
}

// flushMetricBlockTo flushes metric-block of mStore to the writer.
func (vm *versionedTSMap) flushMetricBlocksTo(
	writer metrictbl.TableWriter, metricID uint32, familyTime int64, generator index.IDGenerator) error {

	// this familyTime doesn't exist
	if _, ok := vm.familyTimes[familyTime]; !ok {
		return nil
	}
	all, release := vm.allTSStores()
	defer release()

	for _, tsStore := range *all {
		tsStore.flushTSEntryTo(writer, metricID, familyTime, generator)
	}
	return writer.WriteMetricBlock(metricID)
}

// newVersionedTSMap returns a new versionedTSMap.
func newVersionedTSMap() *versionedTSMap {
	return &versionedTSMap{
		tsMap:       make(map[uint32]*timeSeriesStore),
		familyTimes: make(map[int64]struct{}),
		version:     time.Now().UnixNano()}
}

// metricStore is composed of the immutable part and mutable part of tsMap.
// evictor scans the lruList to check which of them should be purged from the mutable part.
// flusher flushes both the immutable and mutable tsMap to disk, after flushing, the immutable part will be removed.
type metricStore struct {
	immutable    []*versionedTSMap // immutable TSMap list that has not been flushed to disk
	sl4immutable lockers.SpinLock  // spin-lock of immutable versionedTSMap list
	mutable      *versionedTSMap   // current mutable TSMap in use
	mu4Mutable   sync.RWMutex      // sync.RWMutex for mutable tsMap
	name         string            // metric name
	maxTagsLimit uint32            // maximum number of combinations of tags
	metricID     uint32            // default 0, unset
}

// newMetricStore returns a new metricStore from name.
func newMetricStore(name string) *metricStore {
	ms := metricStore{
		name:         name,
		mutable:      newVersionedTSMap(),
		maxTagsLimit: defaultMaxTagsLimit}
	return &ms
}

// setMaxTagsLimit removes race condition.
func (ms *metricStore) setMaxTagsLimit(limit uint32) {
	atomic.StoreUint32(&ms.maxTagsLimit, limit)
}

// mustGetMetricID returns metricID, if unset, generate a new one.
func (ms *metricStore) mustGetMetricID(generator index.IDGenerator) uint32 {
	metricID := atomic.LoadUint32(&ms.metricID)
	if metricID > 0 {
		return metricID
	}
	atomic.CompareAndSwapUint32(&ms.metricID, 0, generator.GenMetricID(ms.name))
	return atomic.LoadUint32(&ms.metricID)
}

// getMaxTagsLimit removes race condition.
func (ms *metricStore) getMaxTagsLimit() uint32 {
	return atomic.LoadUint32(&ms.maxTagsLimit)
}

// addFamilyTime marked this familyTime.
func (ms *metricStore) addFamilyTime(familyTime int64) {
	ms.mu4Mutable.RLock()
	_, ok := ms.mutable.familyTimes[familyTime]
	ms.mu4Mutable.RUnlock()
	if ok {
		return
	}
	ms.mu4Mutable.Lock()
	ms.mutable.familyTimes[familyTime] = struct{}{}
	ms.mu4Mutable.Unlock()
}

// getTSStore returns timeSeriesStore, return false when not exist.
func (ms *metricStore) getTSStore(tagsHash uint32) (tsStore *timeSeriesStore, ok bool) {
	ms.mu4Mutable.RLock()
	tsStore, ok = ms.mutable.tsMap[tagsHash]
	ms.mu4Mutable.RUnlock()
	return
}

// getOrCreateTSStore returns timeSeriesStore by sortedTags.
func (ms *metricStore) getOrCreateTSStore(sortedTags string) *timeSeriesStore {
	tagsHash := hashers.Fnv32a(sortedTags)

	tsStore, ok := ms.getTSStore(tagsHash)
	if !ok {
		ms.mu4Mutable.Lock()
		tsStore, ok = ms.mutable.tsMap[tagsHash]
		if !ok {
			tsStore = newTimeSeriesStore(sortedTags)
			ms.mutable.tsMap[tagsHash] = tsStore
		}
		ms.mu4Mutable.Unlock()
	}
	return tsStore
}

// getTagsCount return the map's length.
func (ms *metricStore) getTagsCount() int {
	ms.mu4Mutable.RLock()
	length := len(ms.mutable.tsMap)
	ms.mu4Mutable.RUnlock()
	return length
}

// isFull detects if timeSeriesMap exceeds the tagsID limitation.
func (ms *metricStore) isFull() bool {
	return uint32(ms.getTagsCount()) >= ms.getMaxTagsLimit()
}

// isEmpty detects if timeSeriesMap is empty or not.
func (ms *metricStore) isEmpty() bool {
	return ms.getTagsCount() == 0
}

// evict scans all metric-stores and removes which are not in use for a while.
func (ms *metricStore) evict() {
	var evictList []uint32
	ms.mu4Mutable.RLock()
	for tagsHash, tStore := range ms.mutable.tsMap {
		if tStore.shouldBeEvicted() {
			evictList = append(evictList, tagsHash)
		}
	}
	ms.mu4Mutable.RUnlock()

	ms.mu4Mutable.Lock()
	for _, tagsHash := range evictList {
		tsStore, ok := ms.mutable.tsMap[tagsHash]
		if !ok {
			continue
		}
		if tsStore.shouldBeEvicted() {
			delete(ms.mutable.tsMap, tagsHash)
		}
	}
	ms.mu4Mutable.Unlock()
}

// unionFamilyTimesTo updates familyTimes of mutable and immutable to the input map.
func (ms *metricStore) unionFamilyTimesTo(families map[int64]struct{}) {
	ms.mu4Mutable.RLock()
	ms.mutable.unionFamilyTimesTo(families)
	ms.mu4Mutable.RUnlock()

	ms.sl4immutable.Lock()
	for _, vm := range ms.immutable {
		vm.unionFamilyTimesTo(families)
	}
	ms.sl4immutable.Unlock()
}

// assignNewVersion moves the mutable TSMap to immutable list, then creates a new mutable map.
func (ms *metricStore) assignNewVersion() error {
	ms.mu4Mutable.Lock()
	if ms.mutable.version+minIntervalForResetMetricStore*int64(time.Millisecond) > time.Now().UnixNano() {
		ms.mu4Mutable.Unlock()
		return fmt.Errorf("reset version of metric: %s too frequently", ms.name)
	}
	oldMutable := ms.mutable
	ms.mutable = newVersionedTSMap()
	ms.mu4Mutable.Unlock()

	ms.sl4immutable.Lock()
	ms.immutable = append(ms.immutable, oldMutable)
	ms.sl4immutable.Unlock()

	return nil
}

// updateTSIDAndFieldID calls mustGet to generate tsID and fieldID.
func (ms *metricStore) updateTSIDAndFieldID(metricID uint32, generator index.IDGenerator) {
	ms.mu4Mutable.RLock()
	tsStores, release := ms.mutable.allTSStores()
	ms.mu4Mutable.RUnlock()
	defer release()

	for _, tsStore := range *tsStores {
		tsStore.mustGetTSID(metricID, generator)
		tsStore.generateFieldsID(metricID, generator)
	}
}

// flushMetricBlocksTo writes metric-blocks to the writer.
func (ms *metricStore) flushMetricBlocksTo(
	writer metrictbl.TableWriter, familyTime int64, generator index.IDGenerator) error {

	var err error
	metricID := ms.mustGetMetricID(generator)

	// pick immutable
	ms.sl4immutable.Lock()
	// build immutable metric-blocks
	for _, vm := range ms.immutable {
		if vm == nil {
			continue
		}
		err = vm.flushMetricBlocksTo(writer, metricID, familyTime, generator)
		delete(vm.familyTimes, familyTime)
		if err != nil {
			ms.sl4immutable.Unlock()
			return err
		}
	}
	ms.sl4immutable.Unlock()

	ms.mu4Mutable.RLock()
	err = ms.mutable.flushMetricBlocksTo(writer, metricID, familyTime, generator)
	ms.mu4Mutable.RUnlock()

	ms.mu4Mutable.Lock()
	delete(ms.mutable.familyTimes, familyTime)
	ms.mu4Mutable.Unlock()

	if err != nil {
		return err
	}
	return nil
}

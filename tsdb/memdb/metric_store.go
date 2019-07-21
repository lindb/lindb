package memdb

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

// versionedTSMap holds a mapping relation of tags and TsStore.
// a version is assigned since newed,
type versionedTSMap struct {
	tsMap       map[string]*timeSeriesStore // map-key: sortedTag
	familyTimes map[int64]struct{}          // all segments
	version     int64                       // uptime in nanoseconds
}

// allTSStores returns tsStore list and related tags list of tsID.
func (vm *versionedTSMap) allTSStores() (tagsList *[]string, tsStoreList *[]*timeSeriesStore, release func()) {
	length := len(vm.tsMap)
	tsStoreList = tsStoresListPool.get(length)
	tagsList = stringListPool.get(length)

	var count = 0
	for tags, tsStore := range vm.tsMap {
		(*tagsList)[count] = tags
		(*tsStoreList)[count] = tsStore
		count++
	}
	return tagsList, tsStoreList, func() {
		stringListPool.put(tagsList)
		tsStoresListPool.put(tsStoreList)
	}
}

// unionFamilyTimesTo update familyTimes to the input map.
func (vm *versionedTSMap) unionFamilyTimesTo(segments map[int64]struct{}) {
	for familyTime := range vm.familyTimes {
		segments[familyTime] = struct{}{}
	}
}

// updateTSIDAndFieldID calls mustGet to generate tsID and fieldID.
func (vm *versionedTSMap) updateTSIDAndFieldID(generator index.IDGenerator, metricID uint32) {
	tagsList, tsStoreList, release := vm.allTSStores()
	defer release()

	for idx, tsStore := range *tsStoreList {
		tags := (*tagsList)[idx]
		tsStore.mustGetTSID(generator, metricID, tags, vm.version)
		tsStore.generateFieldsID(metricID, generator)
	}
}

// flushMetricBlockTo flushes metric-block of mStore to the writer.
func (vm *versionedTSMap) flushMetricBlocksTo(writer metrictbl.TableWriter, familyTime int64,
	generator index.IDGenerator, metricID uint32) error {

	// this familyTime doesn't exist
	if _, ok := vm.familyTimes[familyTime]; !ok {
		return nil
	}
	tagsList, tsStoreList, release := vm.allTSStores()
	defer release()

	for idx, tsStore := range *tsStoreList {
		tags := (*tagsList)[idx]
		tsStore.flushTSEntryTo(writer, familyTime, generator, metricID, tags, vm.version)
	}
	return writer.WriteMetricBlock(metricID)
}

// newVersionedTSMap returns a new versionedTSMap.
func newVersionedTSMap() *versionedTSMap {
	return &versionedTSMap{
		tsMap:       make(map[string]*timeSeriesStore),
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
	maxTagsLimit uint32            // maximum number of combinations of tags
	metricID     uint32            // default 0, unset
}

// newMetricStore returns a new metricStore.
func newMetricStore() *metricStore {
	ms := metricStore{
		mutable:      newVersionedTSMap(),
		maxTagsLimit: defaultMaxTagsLimit}
	return &ms
}

// setMaxTagsLimit removes race condition.
func (ms *metricStore) setMaxTagsLimit(limit uint32) {
	atomic.StoreUint32(&ms.maxTagsLimit, limit)
}

// mustGetMetricID returns metricID, if unset, generate a new one.
func (ms *metricStore) mustGetMetricID(generator index.IDGenerator, metricName string) uint32 {
	metricID := atomic.LoadUint32(&ms.metricID)
	if metricID > 0 {
		return metricID
	}
	atomic.CompareAndSwapUint32(&ms.metricID, 0, generator.GenMetricID(metricName))
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
func (ms *metricStore) getTSStore(tags string) (tsStore *timeSeriesStore, ok bool) {
	ms.mu4Mutable.RLock()
	tsStore, ok = ms.mutable.tsMap[tags]
	ms.mu4Mutable.RUnlock()
	return
}

// getOrCreateTSStore returns timeSeriesStore by sortedTags.
func (ms *metricStore) getOrCreateTSStore(tags string) *timeSeriesStore {
	tsStore, ok := ms.getTSStore(tags)
	if !ok {
		ms.mu4Mutable.Lock()
		tsStore, ok = ms.mutable.tsMap[tags]
		if !ok {
			tsStore = newTimeSeriesStore()
			ms.mutable.tsMap[tags] = tsStore
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

	var evictList []string
	ms.mu4Mutable.RLock()
	for tags, tStore := range ms.mutable.tsMap {
		if tStore.shouldBeEvicted() {
			evictList = append(evictList, tags)
		}
	}
	ms.mu4Mutable.RUnlock()

	ms.mu4Mutable.Lock()
	for _, tags := range evictList {
		tsStore, ok := ms.mutable.tsMap[tags]
		if !ok {
			continue
		}
		if tsStore.shouldBeEvicted() {
			delete(ms.mutable.tsMap, tags)
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
		return fmt.Errorf("reset version too frequently")
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
func (ms *metricStore) updateTSIDAndFieldID(generator index.IDGenerator, metricID uint32) {
	ms.mu4Mutable.RLock()
	ms.mutable.updateTSIDAndFieldID(generator, metricID)
	ms.mu4Mutable.RUnlock()
}

// flushMetricBlocksTo writes metric-blocks to the writer.
func (ms *metricStore) flushMetricBlocksTo(writer metrictbl.TableWriter, familyTime int64,
	generator index.IDGenerator, metricName string) error {

	var err error
	metricID := ms.mustGetMetricID(generator, metricName)

	// pick immutable
	ms.sl4immutable.Lock()
	// build immutable metric-blocks
	for _, vm := range ms.immutable {
		if vm == nil {
			continue
		}
		err = vm.flushMetricBlocksTo(writer, familyTime, generator, metricID)
		delete(vm.familyTimes, familyTime)
		if err != nil {
			ms.sl4immutable.Unlock()
			return err
		}
	}
	ms.sl4immutable.Unlock()

	ms.mu4Mutable.RLock()
	err = ms.mutable.flushMetricBlocksTo(writer, familyTime, generator, metricID)
	ms.mu4Mutable.RUnlock()

	ms.mu4Mutable.Lock()
	delete(ms.mutable.familyTimes, familyTime)
	ms.mu4Mutable.Unlock()

	if err != nil {
		return err
	}
	return nil
}

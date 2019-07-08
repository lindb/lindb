package memdb

import (
	"container/list"
	"regexp"
	"sort"
	"sync"
	"sync/atomic"
)

// metricStore holds a mapping relation of tag and TsStore.
type metricStore struct {
	// map-key: tagsID
	tsMap map[string]*timeSeriesStore
	// sync.RWMutex for tsMap.
	mu4Map sync.RWMutex
	// timeSeriesList records the last accessed timestamps of TimeSeriesStore
	// evictor will scans the list to check which of them should be purged from the map.
	lruList *list.List
	// Sync.Mutex for timeSeriesList
	mu4List sync.Mutex
	name    string
	// maximum number of combinations of tags
	maxTagsLimit uint32
	// tsSeq           uint32
}

// newMetricStore returns a new metricStore from name.
func newMetricStore(name string) *metricStore {
	ms := metricStore{
		name:         name,
		tsMap:        make(map[string]*timeSeriesStore),
		maxTagsLimit: defaultMaxTagsLimit,
		lruList:      list.New()}
	return &ms
}

// setMaxTagsLimit removes race condition.
func (ms *metricStore) setMaxTagsLimit(limit uint32) {
	atomic.StoreUint32(&ms.maxTagsLimit, limit)
}

// getMaxTagsLimit removes race condition.
func (ms *metricStore) getMaxTagsLimit() uint32 {
	return atomic.LoadUint32(&ms.maxTagsLimit)
}

// regexSearchTags search tags which matches the pattern.
func (ms *metricStore) regexSearchTags(tagsIDPattern string) []string {
	if tagsIDPattern == "" {
		return nil
	}
	validPattern, err := regexp.Compile(tagsIDPattern)
	if err != nil {
		return nil
	}
	var matched []string
	ms.mu4Map.RLock()
	for tagsID := range ms.tsMap {
		if validPattern.MatchString(tagsID) {
			matched = append(matched, tagsID)
		}
	}
	ms.mu4Map.RUnlock()
	sort.Slice(matched, func(i, j int) bool {
		return matched[i] < matched[j]
	})
	return matched
}

// getTimeSeries returns timeSeriesStore by tagsID.
func (ms *metricStore) getTimeSeries(tagsID string) *timeSeriesStore {
	ms.mu4Map.RLock()
	tsStore, exist := ms.tsMap[tagsID]
	ms.mu4Map.RUnlock()

	if !exist {
		ms.mu4Map.Lock()
		tsStore, exist = ms.tsMap[tagsID]
		if !exist {
			tsStore = newTimeSeriesStore(tagsID)
			ms.tsMap[tagsID] = tsStore
		}
		ms.mu4Map.Unlock()
	}

	ms.mu4List.Lock()
	if exist {
		ms.lruList.MoveToBack(tsStore.element)
	} else {
		tsStore.element = ms.lruList.PushBack(tsStore)
	}
	ms.mu4List.Unlock()

	return tsStore
}

// getTagsCount return the map's length.
func (ms *metricStore) getTagsCount() int {
	ms.mu4Map.RLock()
	length := len(ms.tsMap)
	ms.mu4Map.RUnlock()
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
	ms.mu4List.Lock()
	defer ms.mu4List.Unlock()

	var next *list.Element
	for e := ms.lruList.Front(); e != nil; e = next {
		next = e.Next()
		tsStore := e.Value.(*timeSeriesStore)
		if tsStore.shouldBeEvicted() {
			// remove from map
			ms.mu4Map.Lock()
			delete(ms.tsMap, tsStore.tagsID)
			ms.mu4Map.Unlock()
			// remove from list
			ms.lruList.Remove(e)
		} else {
			// elements behind this are still in use.
			break
		}
	}
}

package memdb

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/series"
	"github.com/lindb/lindb/tsdb/tblstore"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock_test.go -package memdb

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// getMetricID returns the metricID
	getMetricID() uint32
	// suggestTagKeys returns tagKeys by prefix-search
	suggestTagKeys(tagKeyPrefix string, limit int) []string
	// suggestTagValues returns tagValues by prefix-search
	suggestTagValues(tagKey, tagValuePrefix string, limit int) []string
	// getTagValues get tagValues from the specified version and tagKeys
	getTagValues(tagKeys []string, version uint32) (tagValues [][]string, err error)
	// write writes the metric
	write(metric *pb.Metric, writeCtx writeContext) error
	// setMaxTagsLimit sets the max tags-limit
	setMaxTagsLimit(limit uint32)
	// isEmpty detects whether if tags number is empty or not.
	isEmpty() bool
	// evict scans all tsStore and removes which are not in use for a while.
	evict()
	// getTagsInUse return the in-use tStores count.
	getTagsInUse() int
	// getTagsUsed return count of all used tStores.
	getTagsUsed() int
	// flushMetricsTo flushes metric-block of mStore to the writer.
	flushMetricsTo(tableFlusher tblstore.MetricsDataFlusher, flushCtx flushContext) error
	// flushForwardIndexTo flushes metric-block of mStore to the writer.
	flushForwardIndexTo(tableFlusher tblstore.ForwardIndexFlusher) error
	// flushInvertedIndexTo flushes series-index of mStore to the writer
	flushInvertedIndexTo(tableFlusher tblstore.InvertedIndexFlusher, idGenerator diskdb.IDGenerator) error
	// resetVersion moves the current running mutable index to immutable list,
	// then creates a new mutable map.
	resetVersion() error
	// findSeriesIDsByExpr finds series ids by tag filter expr
	findSeriesIDsByExpr(expr stmt.TagFilter) (*series.MultiVerSeriesIDSet, error)
	// getSeriesIDsForTag get series ids by tagKey
	getSeriesIDsForTag(tagKey string) (*series.MultiVerSeriesIDSet, error)
	mStoreFieldIDGetter
	// scan returns a iterator for scanning data
	scan(sCtx series.ScanContext) series.VersionIterator
}

// mStoreFieldIDGetter gets fieldID from fieldsMeta, and calls the id-generator when not exist
type mStoreFieldIDGetter interface {
	getFieldIDOrGenerate(fieldName string, fieldType field.Type, generator diskdb.IDGenerator) (uint16, error)
}

// metricStore is composed of the immutable part and mutable part of indexes.
// evictor scans the index to check which of them should be purged from the mutable part.
// flusher flushes both the immutable and mutable index to disk,
// after flushing, the immutable part will be removed.
type metricStore struct {
	immutable       []tagIndexINTF // immutable index that has not been flushed to disk
	mutex4Immutable sync.RWMutex   // read-write lock  for immutable index
	mutable         tagIndexINTF   // current mutable index in use
	mutex4Mutable   sync.RWMutex   // read-write lock  for mutable index
	maxTagsLimit    uint32         // maximum number of combinations of tags
	metricID        uint32         // persistent on the disk
	mutex4Fields    sync.RWMutex   // read-write lock for fieldsMeta
	fieldsMetas     fieldsMetas    // mapping of fieldNames and fieldIDs
}

// fieldMeta contains the meta info of field
type fieldMeta struct {
	fieldName string
	fieldID   uint16
	fieldType field.Type
}

// fieldsMetas implements sort.Interface
type fieldsMetas []fieldMeta

func (fm fieldsMetas) Len() int           { return len(fm) }
func (fm fieldsMetas) Less(i, j int) bool { return fm[i].fieldName < fm[j].fieldName }
func (fm fieldsMetas) Swap(i, j int)      { fm[i], fm[j] = fm[j], fm[i] }

// getField get the fieldMeta from fieldName, return false when not exist
func (fm fieldsMetas) getFieldMeta(fieldName string) (fieldMeta, bool) {
	idx := sort.Search(len(fm), func(i int) bool { return fm[i].fieldName >= fieldName })
	if idx >= len(fm) || fm[idx].fieldName != fieldName {
		return fieldMeta{}, false
	}
	return fm[idx], true
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore(metricID uint32) mStoreINTF {
	ms := metricStore{
		metricID:     metricID,
		mutable:      newTagIndex(),
		maxTagsLimit: constants.DefaultMStoreMaxTagsCount}
	return &ms
}

// getFieldIDOrGenerate gets fieldID from fieldsMeta, and calls the id-generator when not exist
func (ms *metricStore) getFieldIDOrGenerate(fieldName string, fieldType field.Type,
	generator diskdb.IDGenerator) (uint16, error) {

	ms.mutex4Fields.RLock()
	fm, ok := ms.fieldsMetas.getFieldMeta(fieldName)
	ms.mutex4Fields.RUnlock()
	// exist, check fieldType
	if ok {
		if fm.fieldType == fieldType {
			return fm.fieldID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// not exist
	ms.mutex4Fields.Lock()
	defer ms.mutex4Fields.Unlock()
	fm, ok = ms.fieldsMetas.getFieldMeta(fieldName)
	// double check
	if ok {
		return fm.fieldID, nil
	}
	// forbid creating new fStore when full
	if len(ms.fieldsMetas) >= constants.TStoreMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}
	// generate and check fieldType
	newFieldID, err := generator.GenFieldID(ms.metricID, fieldName, fieldType)
	if err != nil { // fieldType not matches to the existed
		return 0, err
	}
	ms.fieldsMetas = append(ms.fieldsMetas, fieldMeta{
		fieldName: fieldName, fieldID: newFieldID, fieldType: fieldType})
	sort.Sort(ms.fieldsMetas)
	return newFieldID, nil

}

// getMetricID returns the metricID
func (ms *metricStore) getMetricID() uint32 {
	return atomic.LoadUint32(&ms.metricID)
}

// suggestTagKeys returns tagKeys by prefix-search
func (ms *metricStore) suggestTagKeys(tagKeyPrefix string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	var tagKeys []string
	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	defer ms.mutex4Immutable.RUnlock()
	defer ms.mutex4Mutable.RUnlock()

	prefixSearchTagKey := func(tagIndex tagIndexINTF) {
		for _, entrySet := range tagIndex.getTagKVEntrySets() {
			if len(tagKeys) >= limit {
				return
			}
			if strings.HasPrefix(entrySet.key, tagKeyPrefix) {
				tagKeys = append(tagKeys, entrySet.key)
			}
		}
	}
	for _, indexINTF := range ms.immutable {
		prefixSearchTagKey(indexINTF)
	}
	prefixSearchTagKey(ms.mutable)
	return tagKeys
}

// suggestTagValues returns tagValues by prefix-search
func (ms *metricStore) suggestTagValues(tagKey, tagValuePrefix string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	var tagValues []string
	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	defer ms.mutex4Immutable.RUnlock()
	defer ms.mutex4Mutable.RUnlock()

	prefixSearchTagValue := func(tagIndex tagIndexINTF) {
		for _, entrySet := range tagIndex.getTagKVEntrySets() {
			if len(tagValues) >= limit {
				return
			}
			for tagValue := range entrySet.values {
				if strings.HasPrefix(tagValue, tagValuePrefix) {
					tagValues = append(tagValues, tagValue)
				}
			}
		}
	}
	for _, indexINTF := range ms.immutable {
		prefixSearchTagValue(indexINTF)
	}
	prefixSearchTagValue(ms.mutable)
	return tagValues
}

// getTagValues get tagValues from the specified version and tagKeys
func (ms *metricStore) getTagValues(tagKeys []string, version uint32) (tagValues [][]string, err error) {
	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	defer ms.mutex4Immutable.RUnlock()
	defer ms.mutex4Mutable.RUnlock()

	var found tagIndexINTF
	for _, indexINTF := range ms.immutable {
		if indexINTF.getVersion() == version {
			found = indexINTF
		}
	}
	if ms.mutable.getVersion() == version {
		found = ms.mutable
	}
	if found == nil {
		return nil, series.ErrNotFound
	}
	for _, tagKey := range tagKeys {
		entrySet, ok := found.getTagKVEntrySet(tagKey)
		if !ok {
			tagValues = append(tagValues, nil)
			continue
		}
		var tagValueList []string
		for tagValue := range entrySet.values {
			tagValueList = append(tagValueList, tagValue)
		}
		tagValues = append(tagValues, tagValueList)
	}
	return tagValues, nil
}

// write writes the metric to the tStore
func (ms *metricStore) write(metric *pb.Metric, writeCtx writeContext) error {
	if ms.isFull() {
		return series.ErrTooManyTags
	}
	var err error
	tStore, ok := ms.getTStore(metric.Tags)
	if !ok {
		ms.mutex4Mutable.Lock()
		tStore, err = ms.mutable.getOrCreateTStore(metric.Tags)
		if err != nil {
			ms.mutex4Mutable.Unlock()
			return err
		}
		ms.mutex4Mutable.Unlock()
	}
	err = tStore.write(metric, writeCtx)
	if err == nil {
		ms.mutex4Mutable.RLock()
		ms.mutable.updateTime(uint32(writeCtx.PointTime() / 1000))
		ms.mutex4Mutable.RUnlock()
	}
	return err
}

// setMaxTagsLimit sets the max tags-limit of the metricStore
func (ms *metricStore) setMaxTagsLimit(limit uint32) {
	atomic.StoreUint32(&ms.maxTagsLimit, limit)
}

// getMaxTagsLimit return the max tags limit without race condition.
func (ms *metricStore) getMaxTagsLimit() uint32 {
	return atomic.LoadUint32(&ms.maxTagsLimit)
}

// getTStore returns timeSeriesStore, return false when not exist.
func (ms *metricStore) getTStore(tags map[string]string) (tStore tStoreINTF, ok bool) {
	ms.mutex4Mutable.RLock()
	tStore, ok = ms.mutable.getTStore(tags)
	ms.mutex4Mutable.RUnlock()
	return
}

// getTagsInUse return the tStores count.
func (ms *metricStore) getTagsInUse() int {
	ms.mutex4Mutable.RLock()
	size := ms.mutable.tagsInUse()
	ms.mutex4Mutable.RUnlock()
	return size
}

// getTagsUsed return count of all used tStores.
func (ms *metricStore) getTagsUsed() int {
	ms.mutex4Mutable.RLock()
	size := ms.mutable.tagsUsed()
	ms.mutex4Mutable.RUnlock()
	return size
}

// isFull detects if timeSeriesMap exceeds the tagsID limitation.
func (ms *metricStore) isFull() bool {
	return uint32(ms.getTagsUsed()) >= ms.getMaxTagsLimit()
}

// isEmpty detects if tStores were all evicted or not.
func (ms *metricStore) isEmpty() bool {
	return ms.getTagsInUse() == 0
}

// evict scans all tsStore and removes which are not in use for a while.
func (ms *metricStore) evict() {
	var (
		evictList            []uint32
		doubleCheckEvictList []uint32
	)
	// first check
	ms.mutex4Mutable.RLock()
	for seriesID, tStore := range ms.mutable.allTStores() {
		if tStore.isNoData() && tStore.isExpired() {
			evictList = append(evictList, seriesID)
		}
	}
	ms.mutex4Mutable.RUnlock()
	// double check
	ms.mutex4Mutable.Lock()
	for _, seriesID := range evictList {
		tStore, ok := ms.mutable.getTStoreBySeriesID(seriesID)
		if !ok {
			continue
		}
		if tStore.isNoData() && tStore.isExpired() {
			doubleCheckEvictList = append(doubleCheckEvictList, seriesID)
		}
	}
	ms.mutable.removeTStores(doubleCheckEvictList...)
	ms.mutex4Mutable.Unlock()
}

// resetVersion moves the mutable index to immutable list, then creates a new active index.
func (ms *metricStore) resetVersion() error {
	ms.mutex4Mutable.Lock()
	if ms.mutable.getVersion()+minIntervalForResetMetricStore > uint32(timeutil.Now()/1000) {
		ms.mutex4Mutable.Unlock()
		return fmt.Errorf("reset version too frequently")
	}
	oldMutable := ms.mutable
	ms.mutable = newTagIndex()
	ms.mutex4Mutable.Unlock()

	ms.mutex4Immutable.Lock()
	ms.immutable = append(ms.immutable, oldMutable)
	ms.mutex4Immutable.Unlock()
	return nil
}

// flushMetricsTo writes metric-blocks to the writer.
func (ms *metricStore) flushMetricsTo(flusher tblstore.MetricsDataFlusher, flushCtx flushContext) error {
	// flush field meta info
	ms.mutex4Fields.RLock()
	for _, fm := range ms.fieldsMetas {
		flusher.FlushFieldMeta(fm.fieldID, fm.fieldType)
	}
	ms.mutex4Fields.RUnlock()

	var err error
	// pick immutable
	ms.mutex4Immutable.Lock()
	// build immutable metric-blocks
	for _, tagIdx := range ms.immutable {
		if err = tagIdx.flushMetricTo(flusher, flushCtx); err != nil {
			ms.mutex4Immutable.Unlock()
			return err
		}
	}
	// reset the immutable part
	ms.immutable = nil
	ms.mutex4Immutable.Unlock()

	// reset the mutable part
	ms.mutex4Mutable.RLock()
	err = ms.mutable.flushMetricTo(flusher, flushCtx)
	ms.mutex4Mutable.RUnlock()
	return err
}

// flushForwardIndexTo flushes metric-block of mStore to the writer.
func (ms *metricStore) flushForwardIndexTo(flusher tblstore.ForwardIndexFlusher) error {
	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	defer ms.mutex4Mutable.RUnlock()
	defer ms.mutex4Immutable.RUnlock()

	flushIndexINTF := func(indexINTF tagIndexINTF) {
		for _, entrySet := range indexINTF.getTagKVEntrySets() {
			for tagValue, bitmap := range entrySet.values {
				flusher.FlushTagValue(tagValue, bitmap)
			}
			flusher.FlushTagKey(entrySet.key)
		}
		startTime, endTime := indexINTF.getTimeRange()
		flusher.FlushVersion(indexINTF.getVersion(), startTime, endTime)
	}
	// real flush process
	for _, indexINTF := range ms.immutable {
		flushIndexINTF(indexINTF)
	}
	flushIndexINTF(ms.mutable)
	return flusher.FlushMetricID(ms.metricID)
}

// flushInvertedIndexTo flushes the inverted-index of mStore to the writer
func (ms *metricStore) flushInvertedIndexTo(flusher tblstore.InvertedIndexFlusher,
	idGenerator diskdb.IDGenerator) error {

	ms.mutex4Immutable.RLock()
	ms.mutex4Mutable.RLock()
	defer ms.mutex4Mutable.RUnlock()
	defer ms.mutex4Immutable.RUnlock()

	// build relation of tagKey -> {tagValue1...}
	tagKeyValues := make(map[string]map[string]struct{})
	for _, indexINTF := range ms.immutable {
		for _, entrySet := range indexINTF.getTagKVEntrySets() {
			tagValues := make(map[string]struct{})
			for tagValue := range entrySet.values {
				tagValues[tagValue] = struct{}{}
			}
			tagKeyValues[entrySet.key] = tagValues
		}
	}
	for _, entrySet := range ms.mutable.getTagKVEntrySets() {
		tagValues, ok := tagKeyValues[entrySet.key]
		if !ok {
			tagValues = make(map[string]struct{})
		}
		for tagValue := range entrySet.values {
			tagValues[tagValue] = struct{}{}
		}
		tagKeyValues[entrySet.key] = tagValues
	}

	// flush data process
	for tagKey, tagValues := range tagKeyValues {
		for tagValue := range tagValues {
			for _, indexINTF := range ms.immutable {
				entrySet, ok := indexINTF.getTagKVEntrySet(tagKey)
				if !ok {
					continue
				}
				if bitmap, ok := entrySet.values[tagValue]; ok {
					startTime, endTime := indexINTF.getTimeRange()
					flusher.FlushVersion(indexINTF.getVersion(), startTime, endTime, bitmap)
				}
			}
			entrySet, ok := ms.mutable.getTagKVEntrySet(tagKey)
			if ok {
				if bitmap, ok := entrySet.values[tagValue]; ok {
					startTime, endTime := ms.mutable.getTimeRange()
					flusher.FlushVersion(ms.mutable.getVersion(), startTime, endTime, bitmap)
				}
			}
			flusher.FlushTagValue(tagValue)
		}
		if err := flusher.FlushTagID(idGenerator.GenTagID(ms.metricID, tagKey)); err != nil {
			return err
		}
	}
	return nil
}

// findSeriesIDsByExpr finds series ids by tag filter expr
func (ms *metricStore) findSeriesIDsByExpr(expr stmt.TagFilter) (*series.MultiVerSeriesIDSet, error) {

	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()

	ms.mutex4Immutable.RLock()
	for _, tagIdx := range ms.immutable {
		if bitMap := tagIdx.findSeriesIDsByExpr(expr); bitMap != nil {
			multiVerSeriesIDSet.Add(tagIdx.getVersion(), bitMap)
		}
	}
	ms.mutex4Immutable.RUnlock()

	ms.mutex4Mutable.RLock()
	if bitMap := ms.mutable.findSeriesIDsByExpr(expr); bitMap != nil {
		multiVerSeriesIDSet.Add(ms.mutable.getVersion(), bitMap)
	}
	ms.mutex4Mutable.RUnlock()

	return multiVerSeriesIDSet, nil
}

// getSeriesIDsForTag get series ids by tagKey
func (ms *metricStore) getSeriesIDsForTag(tagKey string) (*series.MultiVerSeriesIDSet, error) {

	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()

	ms.mutex4Immutable.RLock()
	for _, tagIdx := range ms.immutable {
		if bitMap := tagIdx.getSeriesIDsForTag(tagKey); bitMap != nil {
			multiVerSeriesIDSet.Add(tagIdx.getVersion(), bitMap)
		}
	}
	ms.mutex4Immutable.RUnlock()

	ms.mutex4Mutable.RLock()
	if bitMap := ms.mutable.getSeriesIDsForTag(tagKey); bitMap != nil {
		multiVerSeriesIDSet.Add(ms.mutable.getVersion(), bitMap)
	}
	ms.mutex4Mutable.RUnlock()

	return multiVerSeriesIDSet, nil
}

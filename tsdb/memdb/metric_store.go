package memdb

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/lindb/lindb/constants"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/RoaringBitmap/roaring"
	"go.uber.org/atomic"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock_test.go -package memdb

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// GetMetricID returns the metricID
	GetMetricID() uint32

	// SuggestTagKeys returns tagKeys by prefix-search
	SuggestTagKeys(tagKeyPrefix string, limit int) []string

	// SuggestTagValues returns tagValues by prefix-search
	SuggestTagValues(tagKey, tagValuePrefix string, limit int) []string

	// GetTagValues get tagValues from the specified version and tagKeys
	GetTagValues(
		tagKeys []string,
		version series.Version,
		seriesID *roaring.Bitmap,
	) (
		seriesID2TagValues map[uint32][]string,
		err error)

	// Write Writes the metric
	Write(metric *pb.Metric, WriteCtx writeContext) error

	// SetMaxTagsLimit sets the max tags-limit
	SetMaxTagsLimit(limit uint32)

	// IsEmpty detects whether if tags number is empty or not.
	IsEmpty() bool

	// Evict scans all tsStore and removes which are not in use for a while.
	Evict()

	// GetTagsInUse return the in-use tStores count.
	GetTagsInUse() int

	// GetTagsUsed return count of all used tStores.
	GetTagsUsed() int

	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(
		tableFlusher tblstore.MetricsDataFlusher,
		flushCtx flushContext,
	) error

	// FlushForwardIndexTo flushes metric-block of mStore to the Writer.
	FlushForwardIndexTo(tableFlusher tblstore.ForwardIndexFlusher) error

	// FlushInvertedIndexTo flushes series-index of mStore to the Writer
	FlushInvertedIndexTo(
		tableFlusher tblstore.InvertedIndexFlusher,
		idGenerator diskdb.IDGenerator,
	) error

	// ResetVersion moves the current running mutable index to immutable list,
	// then creates a new mutable map.
	ResetVersion() error

	// FindSeriesIDsByExpr finds series ids by tag filter expr
	FindSeriesIDsByExpr(expr stmt.TagFilter) (*series.MultiVerSeriesIDSet, error)

	// GetSeriesIDsForTag get series ids by tagKey
	GetSeriesIDsForTag(tagKey string) (*series.MultiVerSeriesIDSet, error)

	mStoreFieldIDGetter

	series.Scanner
}

type mStoreFieldIDGetter interface {
	// GetFieldIDOrGenerate gets fieldID from fieldsMeta
	// and calls the id-generator when it's not exist
	GetFieldIDOrGenerate(
		fieldName string,
		fieldType field.Type,
		generator diskdb.IDGenerator,
	) (
		fieldID uint16, err error)
}

// metricStore is composed of the immutable part and mutable part of indexes.
// evictor scans the index to check which of them should be purged from the mutable part.
// flusher flushes both the immutable and mutable index to disk,
// after flushing, the immutable part will be removed.
type metricStore struct {
	immutable    atomic.Value  // lock free immutable index that has not been flushed to disk
	mutable      tagIndexINTF  // active mutable index in use
	mux          sync.RWMutex  // read-Write lock for mutable index and fieldMetas
	fieldsMetas  atomic.Value  // read only, storing (*[]fieldMeta), hold mux before storing new value
	maxTagsLimit atomic.Uint32 // maximum number of combinations of tags
	metricID     uint32        // persistent on the disk
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

// copy returns a new copy of fieldsMetas
func (fm fieldsMetas) copy() (clone fieldsMetas) {
	clone = make([]fieldMeta, fm.Len())
	for idx, fm := range fm {
		clone[idx] = fm
	}
	return clone
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore(metricID uint32) mStoreINTF {
	ms := metricStore{
		metricID:     metricID,
		mutable:      newTagIndex(),
		maxTagsLimit: *atomic.NewUint32(constants.DefaultMStoreMaxTagsCount)}
	var fm fieldsMetas
	ms.fieldsMetas.Store(&fm)
	return &ms
}

// getFieldIDOrGenerate gets fieldID from fieldsMeta, and calls the id-generator when not exist
func (ms *metricStore) GetFieldIDOrGenerate(
	fieldName string,
	fieldType field.Type,
	generator diskdb.IDGenerator,
) (
	fieldID uint16,
	err error,
) {
	fmList := ms.fieldsMetas.Load().(*fieldsMetas)
	fm, ok := fmList.getFieldMeta(fieldName)
	// exist, check fieldType
	if ok {
		if fm.fieldType == fieldType {
			return fm.fieldID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// forbid creating new fStore when full
	if fmList.Len() >= constants.TStoreMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}
	// not exist, create a new one
	ms.mux.Lock()
	defer ms.mux.Unlock()

	fmList = ms.fieldsMetas.Load().(*fieldsMetas)
	fm, ok = fmList.getFieldMeta(fieldName)
	// double check
	if ok {
		return fm.fieldID, nil
	}
	// generate and check fieldType
	newFieldID, err := generator.GenFieldID(ms.metricID, fieldName, fieldType)
	if err != nil { // fieldType not matches to the existed
		return 0, err
	}
	clone := fmList.copy()
	clone = append(clone, fieldMeta{
		fieldName: fieldName,
		fieldID:   newFieldID,
		fieldType: fieldType})
	sort.Sort(clone)
	// store the new clone
	ms.fieldsMetas.Store(&clone)
	return newFieldID, nil

}

// GetMetricID returns the metricID
func (ms *metricStore) GetMetricID() uint32 {
	return ms.metricID
}

// SuggestTagKeys returns tagKeys by prefix-search
func (ms *metricStore) SuggestTagKeys(
	tagKeyPrefix string,
	limit int,
) (
	tagKeysList []string,
) {
	if limit <= 0 {
		return nil
	}
	var tagKeysMap = make(map[string]struct{})
	prefixSearchTagKey := func(tagIndex tagIndexINTF) {
		for _, entrySet := range tagIndex.GetTagKVEntrySets() {
			if len(tagKeysMap) >= limit {
				return
			}
			if strings.HasPrefix(entrySet.key, tagKeyPrefix) {
				tagKeysMap[entrySet.key] = struct{}{}
			}
		}
	}
	ms.mux.RLock()
	immutable := ms.atomicGetImmutable()
	prefixSearchTagKey(ms.mutable)
	ms.mux.RUnlock()
	if immutable != nil {
		prefixSearchTagKey(immutable)
	}

	for tagKey := range tagKeysMap {
		tagKeysList = append(tagKeysList, tagKey)
	}
	return tagKeysList
}

// SuggestTagValues returns tagValues by prefix-search
func (ms *metricStore) SuggestTagValues(
	tagKey,
	tagValuePrefix string,
	limit int,
) (
	tagValuesList []string,
) {
	if limit <= 0 {
		return nil
	}
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	var tagValuesMap = make(map[string]struct{})
	prefixSearchTagValue := func(tagIndex tagIndexINTF) {
		for _, entrySet := range tagIndex.GetTagKVEntrySets() {
			if len(tagValuesMap) >= limit {
				return
			}
			for tagValue := range entrySet.values {
				if strings.HasPrefix(tagValue, tagValuePrefix) {
					tagValuesMap[tagValue] = struct{}{}
				}
			}
		}
	}
	ms.mux.RLock()
	immutable := ms.atomicGetImmutable()
	prefixSearchTagValue(ms.mutable)
	ms.mux.RUnlock()
	if immutable != nil {
		prefixSearchTagValue(immutable)
	}

	for tagValue := range tagValuesMap {
		tagValuesList = append(tagValuesList, tagValue)
	}
	return tagValuesList
}

// GetTagValues get tagValues from the specified version and tagKeys
func (ms *metricStore) GetTagValues(
	tagKeys []string,
	version series.Version,
	seriesID *roaring.Bitmap,
) (
	seriesID2TagValues map[uint32][]string,
	err error,
) {
	seriesID2TagValues = make(map[uint32][]string)
	var found tagIndexINTF

	ms.mux.RLock()
	// release the lock when immutable matches to the version
	immutable := ms.atomicGetImmutable()
	if immutable != nil && immutable.Version() == version {
		found = immutable
		ms.mux.RUnlock()
	} else {
		defer ms.mux.RUnlock()
	}
	if ms.mutable.Version() == version {
		found = ms.mutable
	}
	if found == nil {
		return nil, series.ErrNotFound
	}
	// validate tagKeys
	for _, tagKey := range tagKeys {
		_, ok := found.GetTagKVEntrySet(tagKey)
		if !ok {
			return nil, fmt.Errorf("tagKey: %s not exist", tagKey)
		}
	}
	itr := seriesID.Iterator()
	for itr.HasNext() {
		seriesID := itr.Next()
		var tagValues []string
		for _, tagKey := range tagKeys {
			entrySet, ok := found.GetTagKVEntrySet(tagKey)
			if !ok {
				tagValues = append(tagValues, "")
				continue
			}
			var found bool
			for tagValue, bitmap := range entrySet.values {
				if bitmap.Contains(seriesID) {
					found = true
					tagValues = append(tagValues, tagValue)
					break
				}
			}
			if !found {
				tagValues = append(tagValues, "")
			}
		}
		seriesID2TagValues[seriesID] = tagValues
	}
	return seriesID2TagValues, nil
}

// Write Writes the metric to the tStore
func (ms *metricStore) Write(
	metric *pb.Metric,
	writeCtx writeContext,
) error {
	if ms.isFull() {
		return series.ErrTooManyTags
	}
	var err error
	ms.mux.RLock()
	tStore, ok := ms.mutable.GetTStore(metric.Tags)
	ms.mux.RUnlock()
	if !ok {
		ms.mux.Lock()
		tStore, err = ms.mutable.GetOrCreateTStore(metric.Tags)
		if err != nil {
			ms.mux.Unlock()
			return err
		}
		ms.mux.Unlock()
	}
	err = tStore.Write(metric, writeCtx)
	if err == nil {
		ms.mux.RLock()
		ms.mutable.UpdateIndexTimeRange(writeCtx.PointTime())
		ms.mux.RUnlock()
	}
	return err
}

// SetMaxTagsLimit sets the max tags-limit of the metricStore
func (ms *metricStore) SetMaxTagsLimit(limit uint32) {
	ms.maxTagsLimit.Store(limit)
}

// getMaxTagsLimit return the max tags limit without race condition.
func (ms *metricStore) getMaxTagsLimit() uint32 {
	return ms.maxTagsLimit.Load()
}

// GetTagsInUse return the tStores count.
func (ms *metricStore) GetTagsInUse() int {
	ms.mux.RLock()
	size := ms.mutable.TagsInUse()
	ms.mux.RUnlock()
	return size
}

// GetTagsUsed return count of all used tStores.
func (ms *metricStore) GetTagsUsed() int {
	ms.mux.RLock()
	size := ms.mutable.TagsUsed()
	ms.mux.RUnlock()
	return size
}

// isFull detects if timeSeriesMap exceeds the tagsID limitation.
func (ms *metricStore) isFull() bool {
	return uint32(ms.GetTagsUsed()) >= ms.getMaxTagsLimit()
}

// IsEmpty detects if tStores were all Evicted or not.
func (ms *metricStore) IsEmpty() bool {
	return ms.GetTagsInUse() == 0 && ms.atomicGetImmutable() == nil
}

func (ms *metricStore) atomicGetImmutable() tagIndexINTF {
	immutable, ok := ms.immutable.Load().(tagIndexINTF)
	// version zero is the placeholder tagIndexINTF stored in atomic.Value
	if ok && immutable.Version() != 0 {
		return immutable
	}
	return nil
}

// Evict scans all tsStore and removes which are not in use for a while.
func (ms *metricStore) Evict() {
	var (
		evictList            []uint32
		doubleCheckEvictList []uint32
	)
	// first check
	ms.mux.RLock()
	for seriesID, tStore := range ms.mutable.AllTStores() {
		if tStore.IsExpired() && tStore.IsNoData() {
			evictList = append(evictList, seriesID)
		}
	}
	ms.mux.RUnlock()
	// double check
	ms.mux.Lock()
	for _, seriesID := range evictList {
		tStore, ok := ms.mutable.GetTStoreBySeriesID(seriesID)
		if !ok {
			continue
		}
		if tStore.IsExpired() && tStore.IsNoData() {
			doubleCheckEvictList = append(doubleCheckEvictList, seriesID)
		}
	}
	ms.mutable.RemoveTStores(doubleCheckEvictList...)
	ms.mux.Unlock()
}

// ResetVersion marks the mutable index's status to immutable, then creates a new active index.
func (ms *metricStore) ResetVersion() error {
	immutable := ms.atomicGetImmutable()
	if immutable != nil {
		return series.ErrResetVersionUnavailable
	}

	ms.mux.Lock()
	defer ms.mux.Unlock()
	// double check
	immutable = ms.atomicGetImmutable()
	if immutable != nil {
		return series.ErrResetVersionUnavailable
	}
	ms.immutable.Store(ms.mutable)
	ms.mutable = newTagIndex()
	return nil
}

// FlushMetricsTo Writes metric-data to the table.
// immutable tagIndex will be removed after call,
// index shall be flushed before flushing data.
func (ms *metricStore) FlushMetricsDataTo(
	flusher tblstore.MetricsDataFlusher,
	flushCtx flushContext,
) error {
	// flush field meta info
	fmList := ms.fieldsMetas.Load().(*fieldsMetas)
	for _, fm := range *fmList {
		flusher.FlushFieldMeta(fm.fieldID, fm.fieldType)
	}
	// reset the mutable part
	ms.mux.RLock()
	ms.mutable.FlushVersionDataTo(flusher, flushCtx)
	immutable := ms.atomicGetImmutable()
	// remove the immutable, put the nopTagIndex into it
	ms.immutable.Store(staticNopTagIndex)
	ms.mux.RUnlock()

	if immutable != nil {
		immutable.FlushVersionDataTo(flusher, flushCtx)
	}
	return flusher.FlushMetric(flushCtx.metricID)
}

// FlushForwardIndexTo flushes metric-block of mStore to the Writer.
func (ms *metricStore) FlushForwardIndexTo(
	flusher tblstore.ForwardIndexFlusher,
) error {
	flushForwardIndex := func(tagIndex tagIndexINTF) {
		for _, entrySet := range tagIndex.GetTagKVEntrySets() {
			for tagValue, bitmap := range entrySet.values {
				flusher.FlushTagValue(tagValue, bitmap)
			}
			flusher.FlushTagKey(entrySet.key)
		}
		flusher.FlushVersion(tagIndex.Version(), tagIndex.IndexTimeRange())
	}

	ms.mux.RLock()
	immutable := ms.atomicGetImmutable()
	flushForwardIndex(ms.mutable)
	ms.mux.RUnlock()

	if immutable != nil {
		flushForwardIndex(immutable)
	}
	return flusher.FlushMetricID(ms.metricID)
}

// FlushInvertedIndexTo flushes the inverted-index of mStore to the Writer
func (ms *metricStore) FlushInvertedIndexTo(
	flusher tblstore.InvertedIndexFlusher,
	idGenerator diskdb.IDGenerator,
) error {
	// build relation of tagKey -> {tagValue1...}
	tagKeyValues := make(map[string]map[string]struct{})

	ms.mux.RLock()
	defer ms.mux.RUnlock()
	immutable := ms.atomicGetImmutable()
	if immutable != nil {
		for _, entrySet := range immutable.GetTagKVEntrySets() {
			tagValues := make(map[string]struct{})
			for tagValue := range entrySet.values {
				tagValues[tagValue] = struct{}{}
			}
			tagKeyValues[entrySet.key] = tagValues
		}
	}
	for _, entrySet := range ms.mutable.GetTagKVEntrySets() {
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
	flushInvertedIndex := func(tagIndex tagIndexINTF, tagKey, tagValue string) {
		entrySet, ok := tagIndex.GetTagKVEntrySet(tagKey)
		if !ok {
			return
		}
		if bitmap, ok := entrySet.values[tagValue]; ok {
			flusher.FlushVersion(tagIndex.Version(), tagIndex.IndexTimeRange(), bitmap)
		}
	}
	for tagKey, tagValues := range tagKeyValues {
		for tagValue := range tagValues {
			if immutable != nil {
				flushInvertedIndex(immutable, tagKey, tagValue)
			}
			flushInvertedIndex(ms.mutable, tagKey, tagValue)
			flusher.FlushTagValue(tagValue)
		}
		if err := flusher.FlushTagKeyID(idGenerator.GenTagKeyID(ms.metricID, tagKey)); err != nil {
			return err
		}
	}
	return nil
}

// FindSeriesIDsByExpr finds series ids by tag filter expr
func (ms *metricStore) FindSeriesIDsByExpr(
	expr stmt.TagFilter,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()

	findSeriesIDsByExpr := func(tagIdx tagIndexINTF) {
		if bitMap := tagIdx.FindSeriesIDsByExpr(expr); bitMap != nil {
			multiVerSeriesIDSet.Add(tagIdx.Version(), bitMap)
		}
	}
	ms.mux.RLock()
	findSeriesIDsByExpr(ms.mutable)
	immutable := ms.atomicGetImmutable()
	ms.mux.RUnlock()
	if immutable != nil {
		findSeriesIDsByExpr(immutable)
	}
	return multiVerSeriesIDSet, nil
}

// GetSeriesIDsForTag get series ids by tagKey
func (ms *metricStore) GetSeriesIDsForTag(
	tagKey string,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()
	getSeriesIDsForTag := func(tagIdx tagIndexINTF) {
		if bitMap := tagIdx.GetSeriesIDsForTag(tagKey); bitMap != nil {
			multiVerSeriesIDSet.Add(ms.mutable.Version(), bitMap)
		}
	}

	ms.mux.RLock()
	getSeriesIDsForTag(ms.mutable)
	immutable := ms.atomicGetImmutable()
	ms.mux.RUnlock()

	if immutable != nil {
		getSeriesIDsForTag(immutable)
	}
	return multiVerSeriesIDSet, nil
}

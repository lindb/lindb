package memdb

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"

	"go.uber.org/atomic"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock_test.go -package memdb

const emptyMStoreSize = 8 + // immutable
	8 + // mutable
	24 + // rwmutex
	8 + // atomic.Value
	4 + // uint32
	4 + // uint32
	4 // int32

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// SuggestTagKeys returns tagKeys by prefix-search
	SuggestTagKeys(tagKeyPrefix string, limit int) []string

	// SuggestTagValues returns tagValues by prefix-search
	SuggestTagValues(tagKey, tagValuePrefix string, limit int) []string

	// SetMaxTagsLimit sets the max tags-limit
	SetMaxTagsLimit(limit uint32)

	// IsEmpty detects whether if tags number is empty or not.
	IsEmpty() bool

	// GetTagsInUse return the in-use tStores count.
	GetTagsInUse() int

	// GetTagsUsed return count of all used tStores.
	GetTagsUsed() int

	// FlushInvertedIndexTo flushes series-index of mStore to the Writer
	FlushInvertedIndexTo(
		metricID uint32,
		tableFlusher invertedindex.Flusher,
		idGenerator metadb.IDGenerator,
	) error

	// FindSeriesIDsByExpr finds series ids by tag filter expr
	FindSeriesIDsByExpr(expr stmt.TagFilter) (*series.MultiVerSeriesIDSet, error)

	// GetSeriesIDsForTag get series ids by tagKey
	GetSeriesIDsForTag(tagKey string) (*series.MultiVerSeriesIDSet, error)
	// GetSeriesIDsForMetric get series ids by tagKey with min tag values
	GetSeriesIDsForMetric(timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)

	// GetGroupingContext returns the context of group by from the specified version and tagKeys
	GetGroupingContext(tagKeys []string, version series.Version) (series.GroupingContext, error)

	mStoreFieldIDGetter

	// flow.DataFilter filters the data based on condition
	flow.DataFilter

	// MemSize returns the memory-size of this metric-store
	MemSize() int

	///////////////////////////////////
	// Methods below will change the memory size
	///////////////////////////////////
	// Write Writes the metric
	Write(
		metric *pb.Metric,
		writeCtx writeContext,
	) (
		writtenSize int,
		err error)

	// Evict scans all tsStore and removes which are not in use for a while.
	Evict() (evictedSize int)

	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(
		tableFlusher metricsdata.Flusher,
		flushCtx flushContext,
	) (
		flushedSize int,
		err error)

	// ResetVersion moves the current running mutable index to immutable list,
	// then creates a new mutable map.
	ResetVersion() (createdSize int, err error)
}

type mStoreFieldIDGetter interface {
	// GetFieldIDOrGenerate gets fieldID from fieldsMeta
	// and calls the id-generator when it's not exist
	GetFieldIDOrGenerate(
		metricID uint32,
		fieldName string,
		fieldType field.Type,
		generator metadb.IDGenerator,
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
	fieldsMetas  atomic.Value  // read only, storing (field.Metas), hold mux before storing new value
	maxTagsLimit atomic.Uint32 // maximum number of combinations of tags
	size         atomic.Int32  // memory-size
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	mutable := newTagIndex()
	ms := metricStore{
		mutable:      mutable,
		maxTagsLimit: *atomic.NewUint32(constants.DefaultMStoreMaxTagsCount),
		size:         *atomic.NewInt32(int32(mutable.MemSize()))}
	var fm field.Metas
	ms.fieldsMetas.Store(fm)
	return &ms
}

// getFieldIDOrGenerate gets fieldID from fieldsMeta, and calls the id-generator when not exist
func (ms *metricStore) GetFieldIDOrGenerate(
	metricID uint32,
	fieldName string,
	fieldType field.Type,
	generator metadb.IDGenerator,
) (
	fieldID uint16,
	err error,
) {
	fmList := ms.fieldsMetas.Load().(field.Metas)
	fm, ok := fmList.GetFromName(fieldName)
	// exist, check fieldType
	if ok {
		if fm.Type == fieldType {
			return fm.ID, nil
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

	fmList = ms.fieldsMetas.Load().(field.Metas)
	fm, ok = fmList.GetFromName(fieldName)
	// double check
	if ok {
		return fm.ID, nil
	}
	// generate and check fieldType
	newFieldID, err := generator.GenFieldID(metricID, fieldName, fieldType)
	if err != nil { // fieldType not matches to the existed
		return 0, err
	}
	x2 := fmList.Clone()
	x2 = x2.Insert(field.Meta{
		Name: fieldName,
		ID:   newFieldID,
		Type: fieldType})
	// store the new clone
	ms.fieldsMetas.Store(x2)
	return newFieldID, nil

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

func (ms *metricStore) GetGroupingContext(tagKeys []string,
	version series.Version,
) (series.GroupingContext, error) {
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

	tagKeysLen := len(tagKeys)
	gCtx := query.NewGroupContext(tagKeysLen)
	// validate tagKeys
	for idx, tagKey := range tagKeys {
		tagKVEntry, ok := found.GetTagKVEntrySet(tagKey)
		if !ok {
			return nil, fmt.Errorf("tagKey: %s not exist", tagKey)
		}
		tagValuesEntrySet := query.NewTagValuesEntrySet()
		gCtx.SetTagValuesEntrySet(idx, tagValuesEntrySet)
		tagValuesEntrySet.SetTagValues(tagKVEntry.values)
	}
	return &groupingContext{
		ms:   ms,
		gCtx: gCtx,
	}, nil
}

// Write Writes the metric to the tStore
func (ms *metricStore) Write(
	metric *pb.Metric,
	writeCtx writeContext,
) (
	writtenSize int,
	err error,
) {
	if ms.isFull() {
		return 0, series.ErrTooManyTags
	}
	var createdSize int
	ms.mux.RLock()
	tStore, ok := ms.mutable.GetTStore(metric.Tags)
	ms.mux.RUnlock()
	if !ok {
		ms.mux.Lock()
		tStore, createdSize, err = ms.mutable.GetOrCreateTStore(metric.Tags, writeCtx)
		if err != nil {
			ms.mux.Unlock()
			return 0, err
		}
		ms.mux.Unlock()
		ms.size.Add(int32(createdSize))
	}

	writtenSize, err = tStore.Write(metric, writeCtx)
	if err == nil {
		ms.mux.RLock()
		ms.mutable.UpdateIndexTimeRange(writeCtx.PointTime())
		ms.mux.RUnlock()
	}
	ms.size.Add(int32(writtenSize))
	return writtenSize + createdSize, err
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
	count := ms.mutable.TagsInUse()
	ms.mux.RUnlock()
	return count
}

// GetTagsUsed return count of all used tStores.
func (ms *metricStore) GetTagsUsed() int {
	ms.mux.RLock()
	count := ms.mutable.TagsUsed()
	ms.mux.RUnlock()
	return count
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
func (ms *metricStore) Evict() (evictedSize int) {
	var (
		evictList            []uint32
		doubleCheckEvictList []uint32
	)
	// first check
	ms.mux.RLock()
	metricMap := ms.mutable.AllTStores()
	it := metricMap.iterator()
	for it.hasNext() {
		seriesID, tStore := it.next()
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
	removedTStores := ms.mutable.RemoveTStores(doubleCheckEvictList...)
	ms.mux.Unlock()

	for _, tStore := range removedTStores {
		evictedSize += tStore.MemSize()
	}
	ms.size.Sub(int32(evictedSize))
	return evictedSize
}

// ResetVersion marks the mutable index's status to immutable, then creates a new active index.
func (ms *metricStore) ResetVersion() (createdSize int, err error) {
	immutable := ms.atomicGetImmutable()
	if immutable != nil {
		return 0, series.ErrResetVersionUnavailable
	}

	ms.mux.Lock()
	defer ms.mux.Unlock()
	// double check
	immutable = ms.atomicGetImmutable()
	if immutable != nil {
		return 0, series.ErrResetVersionUnavailable
	}
	ms.immutable.Store(ms.mutable)
	ms.mutable = newTagIndex()
	createdSize = ms.mutable.MemSize()
	ms.size.Store(int32(createdSize))
	return createdSize, nil
}

// FlushMetricsTo Writes metric-data to the table.
// immutable tagIndex will be removed after call,
// index shall be flushed before flushing data.
func (ms *metricStore) FlushMetricsDataTo(
	flusher metricsdata.Flusher,
	flushCtx flushContext,
) (
	flushedSize int,
	err error,
) {
	// flush field meta info
	fmList := ms.fieldsMetas.Load().(field.Metas)
	flusher.FlushFieldMetas(fmList)

	// reset the mutable part
	ms.mux.RLock()
	flushedSize = ms.mutable.FlushVersionDataTo(flusher, flushCtx)
	immutable := ms.atomicGetImmutable()
	// remove the immutable, put the nopTagIndex into it
	ms.immutable.Store(staticNopTagIndex)
	ms.mux.RUnlock()

	if immutable != nil {
		flushedSize += immutable.FlushVersionDataTo(flusher, flushCtx)
	}
	ms.size.Sub(int32(flushedSize))
	return flushedSize, flusher.FlushMetric(flushCtx.metricID)
}

// FlushInvertedIndexTo flushes the inverted-index of mStore to the Writer
func (ms *metricStore) FlushInvertedIndexTo(
	metricID uint32,
	flusher invertedindex.Flusher,
	idGenerator metadb.IDGenerator,
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
		if err := flusher.FlushTagKeyID(idGenerator.GenTagKeyID(metricID, tagKey)); err != nil {
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

// GetSeriesIDsForMetric get series ids by tagKey with min tag values
func (ms *metricStore) GetSeriesIDsForMetric(timeRange timeutil.TimeRange) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()
	getSeriesIDsForMetric := func(tagIdx tagIndexINTF) {
		if !timeRange.Contains(tagIdx.Version().Int64()) {
			return
		}
		if bitMap := tagIdx.GetSeriesIDsForMetric(); bitMap != nil && !bitMap.IsEmpty() {
			multiVerSeriesIDSet.Add(ms.mutable.Version(), bitMap)
		}
	}

	ms.mux.RLock()
	getSeriesIDsForMetric(ms.mutable)
	immutable := ms.atomicGetImmutable()
	ms.mux.RUnlock()

	if immutable != nil {
		getSeriesIDsForMetric(immutable)
	}
	return multiVerSeriesIDSet, nil
}

func (ms *metricStore) MemSize() int {
	size := emptyMStoreSize + int(ms.size.Load())
	immutable := ms.atomicGetImmutable()
	if immutable != nil {
		size += immutable.MemSize()
	}
	return size
}

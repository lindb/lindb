package memdb

import (
	"regexp"
	"sort"
	"strings"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"

	"github.com/RoaringBitmap/roaring"
	"github.com/cespare/xxhash"
	"go.uber.org/atomic"
)

//go:generate mockgen -source ./metric_store_index.go -destination=./metric_store_index_mock_test.go -package memdb

const emptyTagIndexSize = 24 + // tagKVEntrySet slice
	24 + // metric-map slice
	4 + // idCounter
	8 + // version
	4 + // earliestTimeDelta
	4 // latestTimeDelta

// tagIndexINTF abstracts the index of tStores, not thread-safe
type tagIndexINTF interface {
	// UpdateIndexTimeRange updates the start and endTime by CAS
	UpdateIndexTimeRange(pointTime int64)

	// IndexTimeRange returns the time range of index
	IndexTimeRange() timeutil.TimeRange

	// GetTagKVEntrySet returns the kv-entrySet by tagKey
	GetTagKVEntrySet(tagKey string) (*tagKVEntrySet, bool)

	// GetTagKVEntrySets returns the kv-entrySets for flushing index data.
	GetTagKVEntrySets() []*tagKVEntrySet

	// GetTStore get tStore from map tags
	GetTStore(tags map[string]string) (tStoreINTF, bool)

	// GetTStoreBySeriesID get tStore from seriesID
	GetTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool)

	// GetOrCreateTStore constructs the index and return a tStore,
	// error of too may tag keys may be return
	GetOrCreateTStore(
		tags map[string]string,
		writeCtx writeContext,
	) (
		tStore tStoreINTF,
		createdSize int,
		err error)

	// RemoveTStores removes tStores from a list of seriesID
	RemoveTStores(seriesIDs ...uint32) (removedTStores []tStoreINTF)

	// TagsUsed returns the count of all used tags, it is used for restricting write.
	TagsUsed() int

	// TagsInUse returns how many tags are still in use, it is used for evicting
	TagsInUse() int

	// AllTStores returns the map of seriesID and tStores
	AllTStores() *metricMap

	// FlushVersionDataTo flush metric to the tableFlusher
	FlushVersionDataTo(
		flusher metricsdata.Flusher,
		flushCtx flushContext,
	) (flushedSize int)

	// Version returns a version(uptime in milliseconds) of the index
	Version() series.Version

	// FindSeriesIDsByExpr finds series ids by tag filter expr
	FindSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap

	// GetSeriesIDsForTag get series ids by tagKey
	GetSeriesIDsForTag(tagKey string) *roaring.Bitmap

	// MemSize returns the memory size in bytes
	MemSize() int

	// scan scans metric store data based on scanner context
	scan(sCtx *series.ScanContext)
}

// tagKVEntrySet is a inverted mapping relation of tag-value and seriesID group.
type tagKVEntrySet struct {
	key    string
	values map[string]*roaring.Bitmap
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagKVEntrySet(tagKey string) *tagKVEntrySet {
	return &tagKVEntrySet{
		key:    tagKey,
		values: make(map[string]*roaring.Bitmap)}
}

// tagIndex implements tagIndexINTF,
// it is a composition of both inverted and forward index,
// not thread-safe
type tagIndex struct {
	// invertedIndex part for storing a mapping from tag-keys to the tsStore list,
	// the purpose of this index is to allow fast filtering and querying
	tagKVEntrySet   []*tagKVEntrySet
	seriesID2TStore *metricMap
	// forwardIndex for storing a mapping from tag-hash to the seriesID,
	// purpose of this index is used for fast writing
	hash2SeriesID map[uint64]uint32
	idCounter     atomic.Uint32
	// version is the uptime in milliseconds
	version series.Version
	// index time-range
	earliestTimeDelta atomic.Int32 // earliestTime = versionTime + earliestTimeDelta
	latestTimeDelta   atomic.Int32 // latestTime = versionTime + latestTimeDelta
}

// newTagIndex returns a new tagIndexINTF with version.
func newTagIndex() tagIndexINTF {
	return &tagIndex{
		seriesID2TStore:   newMetricMap(),
		hash2SeriesID:     make(map[uint64]uint32),
		version:           series.NewVersion(),
		idCounter:         *atomic.NewUint32(0), // first value is 1
		earliestTimeDelta: *atomic.NewInt32(0),
		latestTimeDelta:   *atomic.NewInt32(0)}
}

// UpdateIndexTimeRange updates the start and endTime by CAS
// lock-free
func (index *tagIndex) UpdateIndexTimeRange(pointTime int64) {
	timeDelta := int32((pointTime - index.version.Int64()) / 1000)
	for {
		oldStartTimeDelta := index.earliestTimeDelta.Load()
		if oldStartTimeDelta <= timeDelta {
			break
		}
		if index.earliestTimeDelta.CAS(oldStartTimeDelta, timeDelta) {
			break
		}
	}
	for {
		oldEndTimeDelta := index.latestTimeDelta.Load()
		if oldEndTimeDelta >= timeDelta {
			break
		}
		if index.latestTimeDelta.CAS(oldEndTimeDelta, timeDelta) {
			break
		}
	}
}

// lock-free
func (index *tagIndex) IndexTimeRange() timeutil.TimeRange {
	startTimeDelta, endTimeDelta := index.earliestTimeDelta.Load(), index.latestTimeDelta.Load()
	return timeutil.TimeRange{
		Start: index.version.Int64() + int64(startTimeDelta)*1000,
		End:   index.version.Int64() + int64(endTimeDelta)*1000}
}

// GetTagKVEntrySets returns the kv-entrySet for flushing index data.
func (index *tagIndex) GetTagKVEntrySets() []*tagKVEntrySet {
	return index.tagKVEntrySet
}

// insertNewTStore binds a new tStore to the inverted index to the seriesID.
func (index *tagIndex) insertNewTStore(
	tags map[string]string,
	newSeriesID uint32,
	tStore tStoreINTF,
	writeCtx writeContext,
) error {
	// insert to bitmap
	if tags == nil {
		tags = make(map[string]string)
	}
	if len(tags) == 0 {
		tags[""] = ""
	}
	for tagKey, tagValue := range tags {
		entrySet, created, err := index.getOrInsertTagKeyEntry(tagKey)
		if err != nil {
			return err
		}
		if created {
			// create the tagKeyID synchronously
			_ = writeCtx.generator.GenTagKeyID(writeCtx.metricID, tagKey)
		}
		// create the tagKeyID
		bitMap, ok := entrySet.values[tagValue]
		if !ok {
			bitMap = roaring.NewBitmap()
		}
		bitMap.Add(newSeriesID)
		entrySet.values[tagValue] = bitMap
	}
	// insert to the id mapping
	index.seriesID2TStore.put(newSeriesID, tStore)
	return nil
}

// GetTagKVEntrySet search the tagKeyEntry by binary-search
func (index *tagIndex) GetTagKVEntrySet(tagKey string) (*tagKVEntrySet, bool) {
	offset := sort.Search(len(index.tagKVEntrySet), func(i int) bool { return index.tagKVEntrySet[i].key >= tagKey })
	// not present in the slice
	if offset >= len(index.tagKVEntrySet) || index.tagKVEntrySet[offset].key != tagKey {
		return nil, false
	}
	return index.tagKVEntrySet[offset], true
}

// getOrInsertTagKeyEntry get or insert a new entrySet, return error when tag keys exceeds the limit.
func (index *tagIndex) getOrInsertTagKeyEntry(
	tagKey string,
) (
	entrySet *tagKVEntrySet,
	created bool,
	err error,
) {
	length := len(index.tagKVEntrySet)
	offset := sort.Search(length, func(i int) bool { return index.tagKVEntrySet[i].key >= tagKey })
	// present in the slice
	if offset < len(index.tagKVEntrySet) && index.tagKVEntrySet[offset].key == tagKey {
		return index.tagKVEntrySet[offset], false, nil
	}
	if length >= constants.MStoreMaxTagKeysCount {
		return nil, false, series.ErrTooManyTagKeys
	}
	// not present
	newEntry := newTagKVEntrySet(tagKey)
	index.tagKVEntrySet = append(index.tagKVEntrySet, newEntry)
	// insert, and sort
	if offset < length {
		sort.Slice(index.tagKVEntrySet, func(i, j int) bool {
			return index.tagKVEntrySet[i].key < index.tagKVEntrySet[j].key
		})
	}
	return newEntry, true, nil
}

// GetTStore returns a tStoreINTF from map tags.
func (index *tagIndex) GetTStore(tags map[string]string) (tStoreINTF, bool) {
	hash := xxhash.Sum64String(tag.Concat(tags))
	seriesID, ok := index.hash2SeriesID[hash]
	if ok {
		return index.seriesID2TStore.get(seriesID)
	}
	return nil, false
}

// GetTStoreBySeriesID returns a tStoreINTF from series-id.
func (index *tagIndex) GetTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool) {
	return index.seriesID2TStore.get(seriesID)
}

// GetOrCreateTStore get or creates the tStore from string tags,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *tagIndex) GetOrCreateTStore(
	tags map[string]string,
	writeCtx writeContext,
) (
	tStore tStoreINTF,
	createdSize int,
	err error,
) {
	hash := xxhash.Sum64String(tag.Concat(tags))
	seriesID, ok := index.hash2SeriesID[hash]
	// hash is already existed before
	if ok {
		tStore, ok := index.seriesID2TStore.get(seriesID)
		// has been evicted before, reuse the old seriesID
		if !ok {
			tStore = newTimeSeriesStore()
			index.seriesID2TStore.put(seriesID, tStore)
			createdSize += tStore.MemSize()
		}
		return tStore, createdSize, nil
	}
	// seriesID is not allocated before, assign a new one.
	incrSeriesID := index.idCounter.Inc()
	newTStore := newTimeSeriesStore()
	// bind relation of tag kv pairs to the tStore
	err = index.insertNewTStore(tags, incrSeriesID, newTStore, writeCtx)
	if err != nil {
		index.idCounter.Dec()
		return nil, createdSize, err
	}
	createdSize += newTStore.MemSize()
	// bind relation of hash and seriesID to the forward index
	index.hash2SeriesID[hash] = incrSeriesID
	return newTStore, createdSize, nil
}

// RemoveTStores removes the tStores from seriesIDs
// RemoveTStores does not remove the mapping relation of hash and seriesID and keep the seriesID in bitmap
func (index *tagIndex) RemoveTStores(
	seriesIDs ...uint32,
) (
	removedTStores []tStoreINTF,
) {
	if len(seriesIDs) == 0 {
		return nil
	}
	return index.seriesID2TStore.deleteMany(seriesIDs...)
}

// TagsUsed returns the count of all used tStores
func (index *tagIndex) TagsUsed() int {
	return len(index.hash2SeriesID)
}

// TagsInUse returns how many tags are still in use, it is used for evicting
func (index *tagIndex) TagsInUse() int {
	return index.seriesID2TStore.size()
}

// AllTStores returns the map of seriesID and tStores
func (index *tagIndex) AllTStores() *metricMap {
	return index.seriesID2TStore
}

// FlushVersionDataTo flushes metric-block of mStore to the writer.
func (index *tagIndex) FlushVersionDataTo(
	tableFlusher metricsdata.Flusher,
	flushCtx flushContext,
) (
	flushedSize int,
) {
	it := index.seriesID2TStore.iterator()
	for it.hasNext() {
		seriesID, tStore := it.next()
		flushedSize += tStore.FlushSeriesTo(tableFlusher, flushCtx, seriesID)
	}
	if flushedSize > 0 {
		tableFlusher.FlushVersion(index.Version())
	}
	return flushedSize
}

// Version returns a version(uptime) of the index
func (index *tagIndex) Version() series.Version {
	return index.version
}

// FindSeriesIDsByExpr finds series ids by tag filter expr
func (index *tagIndex) FindSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap {
	entrySet, ok := index.GetTagKVEntrySet(expr.TagKey())
	if !ok {
		return nil
	}
	switch expression := expr.(type) {
	case *stmt.EqualsExpr:
		return index.findSeriesIDsByEqual(entrySet, expression)
	case *stmt.InExpr:
		return index.findSeriesIDsByIn(entrySet, expression)
	case *stmt.LikeExpr:
		return index.findSeriesIDsByLike(entrySet, expression)
	case *stmt.RegexExpr:
		return index.findSeriesIDsByRegex(entrySet, expression)
	}
	return nil
}

func (index *tagIndex) findSeriesIDsByEqual(entrySet *tagKVEntrySet, expr *stmt.EqualsExpr) *roaring.Bitmap {
	bitmap, ok := entrySet.values[expr.Value]
	if !ok {
		return nil
	}
	return bitmap.Clone()
}

func (index *tagIndex) findSeriesIDsByIn(entrySet *tagKVEntrySet, expr *stmt.InExpr) *roaring.Bitmap {
	union := roaring.New()
	for _, value := range expr.Values {
		bitmap, ok := entrySet.values[value]
		if !ok {
			continue
		}
		union.Or(bitmap)
	}
	return union
}

func (index *tagIndex) findSeriesIDsByLike(entrySet *tagKVEntrySet, expr *stmt.LikeExpr) *roaring.Bitmap {
	union := roaring.New()

	likeTo := expr.Value
	switch expr.Value {
	case "":
		return union
	case "*":
		likeTo = ""
	}
	for value, bitmap := range entrySet.values {
		if strings.Contains(value, likeTo) {
			union.Or(bitmap)
		}
	}
	return union
}

func (index *tagIndex) findSeriesIDsByRegex(entrySet *tagKVEntrySet, expr *stmt.RegexExpr) *roaring.Bitmap {
	pattern, err := regexp.Compile(expr.Regexp)
	if err != nil {
		return nil
	}
	// the regex pattern is regarded as a prefix string + pattern
	literalPrefix, _ := pattern.LiteralPrefix()
	union := roaring.New()
	for value, bitmap := range entrySet.values {
		if !strings.HasPrefix(value, literalPrefix) {
			continue
		}
		if pattern.MatchString(value) {
			union.Or(bitmap)
		}
	}
	return union
}

func (index *tagIndex) MemSize() int {
	// tagKVEntrySet, map is not calculated
	size := emptyTagIndexSize
	for _, tStore := range index.seriesID2TStore.stores {
		size += tStore.MemSize()
	}
	return size
}

// GetSeriesIDsForTag get series ids by tagKey
func (index *tagIndex) GetSeriesIDsForTag(tagKey string) *roaring.Bitmap {
	entrySet, ok := index.GetTagKVEntrySet(tagKey)
	if !ok {
		return nil
	}
	union := roaring.New()
	for _, bitMap := range entrySet.values {
		union.Or(bitMap)
	}
	return union
}

// scan scans metric store data based on scanner context
func (index *tagIndex) scan(sCtx *series.ScanContext) {
	index.seriesID2TStore.scan(index.version, sCtx)
}

// staticNopTagIndex is the static nop-tagIndex,
// it is used as a placeholder of immutable atomic.Value
var staticNopTagIndex = newNopTagIndex()

func newNopTagIndex() tagIndexINTF {
	ti := newTagIndex().(*tagIndex)
	ti.version = 0
	return ti
}

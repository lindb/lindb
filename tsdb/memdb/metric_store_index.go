package memdb

import (
	"regexp"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/series"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/RoaringBitmap/roaring"
	"github.com/segmentio/fasthash/fnv1a"
)

//go:generate mockgen -source ./metric_store_index.go -destination=./metric_store_index_mock_test.go -package memdb

// tagIndexINTF abstracts the index of tStores, not thread-safe
type tagIndexINTF interface {
	// updateTime updates the start and endTime by CAS
	updateTime(pointTime uint32)
	// return timeRange in seconds
	getTimeRange() (startTime, endTime uint32)
	// getTagKVEntrySet returns the kv-entrySet by tagKey
	getTagKVEntrySet(tagKey string) (*tagKVEntrySet, bool)
	// getTagKVEntrySets returns the kv-entrySets for flushing index data.
	getTagKVEntrySets() []tagKVEntrySet
	// getTStore get tStore from map tags
	getTStore(tags map[string]string) (tStoreINTF, bool)
	// getTStoreBySeriesID get tStore from seriesID
	getTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool)
	// getOrCreateTStore constructs the index and return a tStore,
	// error of too may tag keys may be return
	getOrCreateTStore(tags map[string]string) (tStoreINTF, error)
	// removeTStores removes tStores from a list of seriesID
	removeTStores(seriesIDs ...uint32)
	// tagsUsed returns the count of all used tags, it is used for restricting write.
	tagsUsed() int
	// tagsInUse returns how many tags are still in use, it is used for evicting
	tagsInUse() int
	// allTStores returns the map of seriesID and tStores
	allTStores() map[uint32]tStoreINTF
	// flushMetricTo flush metric to the tableFlusher
	flushMetricTo(flusher tblstore.MetricsDataFlusher, flushCtx flushContext) error
	// getVersion returns a version(uptime in seconds) of the index
	getVersion() uint32
	// findSeriesIDsByExpr finds series ids by tag filter expr
	findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap
	// getSeriesIDsForTag get series ids by tagKey
	getSeriesIDsForTag(tagKey string) *roaring.Bitmap
}

// tagKVEntrySet is a inverted mapping relation of tag-value and seriesID group.
type tagKVEntrySet struct {
	key    string
	values map[string]*roaring.Bitmap
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagKVEntrySet(tagKey string) tagKVEntrySet {
	return tagKVEntrySet{
		key:    tagKey,
		values: make(map[string]*roaring.Bitmap)}
}

// tagIndex implements tagIndexINTF,
// it is a composition of both inverted and forward index,
// not thread-safe
type tagIndex struct {
	// invertedIndex part for storing a mapping from tag-keys to the tsStore list,
	// the purpose of this index is to allow fast filtering and querying
	tagKVEntrySet   []tagKVEntrySet
	seriesID2TStore map[uint32]tStoreINTF
	// forwardIndex for storing a mapping from tag-hash to the seriesID,
	// purpose of this index is used for fast writing
	hash2SeriesID map[uint64]uint32
	idCounter     uint32
	// version is the uptime in seconds
	version   uint32
	startTime uint32 // startTime, endTime of all data written
	endTime   uint32
}

// newTagIndex returns a new tagIndexINTF with version.
func newTagIndex() tagIndexINTF {
	return &tagIndex{
		seriesID2TStore: make(map[uint32]tStoreINTF),
		hash2SeriesID:   make(map[uint64]uint32),
		version:         uint32(timeutil.Now() / 1000),
		startTime:       uint32(timeutil.Now() / 1000),
		endTime:         uint32(timeutil.Now() / 1000)}
}

// updateTime updates the start and endTime by CAS
func (index *tagIndex) updateTime(pointTime uint32) {
	for {
		startTime := atomic.LoadUint32(&index.startTime)
		if startTime <= pointTime {
			break
		}
		if atomic.CompareAndSwapUint32(&index.startTime, startTime, pointTime) {
			break
		}
	}
	for {
		endTime := atomic.LoadUint32(&index.endTime)
		if endTime >= pointTime {
			break
		}
		if atomic.CompareAndSwapUint32(&index.endTime, endTime, pointTime) {
			break
		}
	}
}

// getTimeRange return timeRange in seconds
func (index *tagIndex) getTimeRange() (startTime, endTime uint32) {
	return atomic.LoadUint32(&index.startTime), atomic.LoadUint32(&index.endTime)
}

// getTagKVEntrySets returns the kv-entrySet for flushing index data.
func (index *tagIndex) getTagKVEntrySets() []tagKVEntrySet {
	return index.tagKVEntrySet
}

// insertNewTStore binds a new tStore to the inverted index to the seriesID.
func (index *tagIndex) insertNewTStore(tags map[string]string, newSeriesID uint32, tStore tStoreINTF) error {
	// insert to bitmap
	if tags == nil {
		tags = make(map[string]string)
	}
	if len(tags) == 0 {
		tags[""] = ""
	}
	for tagKey, tagValue := range tags {
		entrySet, err := index.getOrInsertTagKeyEntry(tagKey)
		if err != nil {
			return err
		}
		bitMap, ok := entrySet.values[tagValue]
		if !ok {
			bitMap = roaring.NewBitmap()
		}
		bitMap.Add(newSeriesID)
		entrySet.values[tagValue] = bitMap
	}
	// insert to the id mapping
	index.seriesID2TStore[newSeriesID] = tStore
	return nil
}

// getTagKVEntrySet search the tagKeyEntry by binary-search
func (index *tagIndex) getTagKVEntrySet(tagKey string) (*tagKVEntrySet, bool) {
	offset := sort.Search(len(index.tagKVEntrySet), func(i int) bool { return index.tagKVEntrySet[i].key >= tagKey })
	// not present in the slice
	if offset >= len(index.tagKVEntrySet) || index.tagKVEntrySet[offset].key != tagKey {
		return nil, false
	}
	return &index.tagKVEntrySet[offset], true
}

// getOrInsertTagKeyEntry get or insert a new entrySet, return error when tag keys exceeds the limit.
func (index *tagIndex) getOrInsertTagKeyEntry(tagKey string) (*tagKVEntrySet, error) {
	length := len(index.tagKVEntrySet)
	offset := sort.Search(length, func(i int) bool { return index.tagKVEntrySet[i].key >= tagKey })
	// present in the slice
	if offset < len(index.tagKVEntrySet) && index.tagKVEntrySet[offset].key == tagKey {
		return &index.tagKVEntrySet[offset], nil
	}
	if length >= constants.MStoreMaxTagKeysCount {
		return nil, series.ErrTooManyTagKeys
	}
	// not present
	newEntry := newTagKVEntrySet(tagKey)
	index.tagKVEntrySet = append(index.tagKVEntrySet, newEntry)
	// insert, not append at the tail
	if offset < length {
		sort.Slice(index.tagKVEntrySet, func(i, j int) bool {
			return index.tagKVEntrySet[i].key < index.tagKVEntrySet[j].key
		})
	}
	return &newEntry, nil
}

// getTStore returns a tStoreINTF from map tags.
func (index *tagIndex) getTStore(tags map[string]string) (tStoreINTF, bool) {
	hash := fnv1a.HashString64(models.TagsAsString(tags))
	seriesID, ok := index.hash2SeriesID[hash]
	if ok {
		return index.seriesID2TStore[seriesID], true
	}
	return nil, false
}

// getTStoreBySeriesID returns a tStoreINTF from series-id.
func (index *tagIndex) getTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool) {
	tStore, ok := index.seriesID2TStore[seriesID]
	return tStore, ok
}

// getOrCreateTStore get or creates the tStore from string tags,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *tagIndex) getOrCreateTStore(tags map[string]string) (tStoreINTF, error) {
	tagsStr := models.TagsAsString(tags)
	hash := fnv1a.HashString64(tagsStr)
	seriesID, ok := index.hash2SeriesID[hash]
	// hash is already existed before
	if ok {
		tStore, ok := index.seriesID2TStore[seriesID]
		// has been evicted before, reuse the old seriesID
		if !ok {
			tStore = newTimeSeriesStore(hash)
			index.seriesID2TStore[seriesID] = tStore
		}
		return tStore, nil
	}
	// seriesID is not allocated before, assign a new one.
	incrSeriesID := atomic.AddUint32(&index.idCounter, 1)
	newTStore := newTimeSeriesStore(hash)
	// bind relation of tag kv pairs to the tStore
	err := index.insertNewTStore(tags, incrSeriesID, newTStore)
	if err != nil {
		index.idCounter--
		return nil, err
	}
	// bind relation of hash and seriesID to the forward index
	index.hash2SeriesID[hash] = incrSeriesID
	return newTStore, nil
}

// removeTStores removes the tStores from seriesIDs
// removeTStores does not remove the mapping relation of hash and seriesID and keep the seriesID in bitmap
func (index *tagIndex) removeTStores(seriesIDs ...uint32) {
	if len(seriesIDs) == 0 {
		return
	}
	// remove from seriesID2TStore
	for _, id := range seriesIDs {
		delete(index.seriesID2TStore, id)
	}
}

// tagsUsed returns the count of all used tStores
func (index *tagIndex) tagsUsed() int {
	return len(index.hash2SeriesID)
}

// tagsInUse returns how many tags are still in use, it is used for evicting
func (index *tagIndex) tagsInUse() int {
	return len(index.seriesID2TStore)
}

// allTStores returns the map of seriesID and tStores
func (index *tagIndex) allTStores() map[uint32]tStoreINTF {
	return index.seriesID2TStore
}

// flushMetricTo flushes metric-block of mStore to the writer.
func (index *tagIndex) flushMetricTo(tableFlusher tblstore.MetricsDataFlusher, flushCtx flushContext) error {
	flushed := false
	for seriesID, tStore := range index.seriesID2TStore {
		tStoreDataFlushed := tStore.flushSeriesTo(tableFlusher, flushCtx, seriesID)
		flushed = flushed || tStoreDataFlushed
	}
	if !flushed {
		return nil
	}
	return tableFlusher.FlushMetric(flushCtx.metricID)
}

// getVersion returns a version(uptime) of the index
func (index *tagIndex) getVersion() uint32 {
	return index.version
}

// findSeriesIDsByExpr finds series ids by tag filter expr
func (index *tagIndex) findSeriesIDsByExpr(expr stmt.TagFilter) *roaring.Bitmap {
	entrySet, ok := index.getTagKVEntrySet(expr.TagKey())
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
	for value, bitmap := range entrySet.values {
		if strings.Contains(value, expr.Value) {
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

// getSeriesIDsForTag get series ids by tagKey
func (index *tagIndex) getSeriesIDsForTag(tagKey string) *roaring.Bitmap {
	entrySet, ok := index.getTagKVEntrySet(tagKey)
	if !ok {
		return nil
	}
	union := roaring.New()
	for _, bitMap := range entrySet.values {
		union.Or(bitMap)
	}
	return union
}

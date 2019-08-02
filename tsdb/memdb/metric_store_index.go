package memdb

import (
	"regexp"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metrictbl"

	"github.com/RoaringBitmap/roaring"
	"github.com/segmentio/fasthash/fnv1a"
)

//go:generate mockgen -source ./metric_store_index.go -destination=./metric_store_index_mock_test.go -package memdb

// tagIndexINTF abstracts the index of tStores, not thread-safe
type tagIndexINTF interface {
	// getTStore get tStore from string tags
	getTStore(tags string) (tStoreINTF, bool)
	// getTStoreBySeriesID get tStore from seriesID
	getTStoreBySeriesID(seriesID uint32) (tStoreINTF, bool)
	// getOrCreateTStore constructs the index and return a tStore,
	// error of too may tag keys may be return
	getOrCreateTStore(tags string) (tStoreINTF, error)
	// removeTStores removes tStores from a list of seriesID
	removeTStores(seriesIDs ...uint32)
	// len returns the count of tStores
	len() int
	// allTStores returns the map of seriesID and tStores
	allTStores() map[uint32]tStoreINTF
	// flushMetricTo flush metric to the tableFlusher
	flushMetricTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) error
	// getVersion returns a version(uptime) of the index
	getVersion() int64
	// findSeriesIDsByExpr finds series ids by tag filter expr and timeRange
	findSeriesIDsByExpr(expr stmt.TagFilter, timeRange timeutil.TimeRange) *roaring.Bitmap
	// getSeriesIDsForTag get series ids by tagKey and timeRange
	getSeriesIDsForTag(tagKey string, timeRange timeutil.TimeRange) *roaring.Bitmap
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
	// the purpose of this index is to allow fast writing
	hash2SeriesID map[uint64]uint32
	idCounter     uint32
	// version is the uptime in milliseconds
	version int64
}

// newTagIndex returns a new tagIndexINTF with version.
func newTagIndex() tagIndexINTF {
	return &tagIndex{
		seriesID2TStore: make(map[uint32]tStoreINTF),
		hash2SeriesID:   make(map[uint64]uint32),
		version:         timeutil.Now()}
}

// insertNewTStore binds a new tStore to the inverted index to the seriesID.
func (index *tagIndex) insertNewTStore(tag string, newSeriesID uint32, tStore tStoreINTF) error {
	// insert to bitmap
	tagPairs := models.NewTags(tag)
	if len(tagPairs) == 0 {
		tagPairs = []models.Tag{{Key: "", Value: ""}}
	}
	for _, tagPair := range tagPairs {
		entrySet, err := index.getOrInsertTagKeyEntry(tagPair.Key)
		if err != nil {
			return err
		}
		bitMap, ok := entrySet.values[tagPair.Value]
		if !ok {
			bitMap = roaring.NewBitmap()
		}
		bitMap.Add(newSeriesID)
		entrySet.values[tagPair.Value] = bitMap
	}
	// insert to the id mapping
	index.seriesID2TStore[newSeriesID] = tStore
	return nil
}

// getTagKeyEntry search the tagKeyEntry by binary-search
func (index *tagIndex) getTagKeyEntry(tagKey string) (*tagKVEntrySet, bool) {
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
	if length >= defaultMaxTagKeys {
		return nil, models.ErrTooManyTagKeys
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

// getTStore returns a tStoreINTF from string tags.
func (index *tagIndex) getTStore(tags string) (tStoreINTF, bool) {
	hash := fnv1a.HashString64(tags)
	theTagID, ok := index.hash2SeriesID[hash]
	if ok {
		return index.seriesID2TStore[theTagID], true
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
func (index *tagIndex) getOrCreateTStore(tags string) (tStoreINTF, error) {
	hash := fnv1a.HashString64(tags)
	seriesID, ok := index.hash2SeriesID[hash]
	if ok {
		return index.seriesID2TStore[seriesID], nil
	}
	// seriesID is not allocated before, assign a new one.
	incrSeriesID := atomic.AddUint32(&index.idCounter, 1)
	newTStore := newTimeSeriesStore(incrSeriesID, fnv1a.HashString64(tags))
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

// removeTStores removes the seriesID from both the forward and inverted index.
func (index *tagIndex) removeTStores(seriesIDs ...uint32) {
	if len(seriesIDs) == 0 {
		return
	}
	var tagHashes []uint64
	// remove from bitmap
	for _, entrySet := range index.tagKVEntrySet {
		for _, bitMap := range entrySet.values {
			for _, id := range seriesIDs {
				bitMap.Remove(id)
			}
		}
	}
	// remove from seriesID2TStore
	for _, id := range seriesIDs {
		tStore, ok := index.seriesID2TStore[id]
		if ok {
			delete(index.seriesID2TStore, id)
			tagHashes = append(tagHashes, tStore.getHash())
		}
	}
	// remove from forward index
	for _, tagHash := range tagHashes {
		delete(index.hash2SeriesID, tagHash)
	}
}

// len returns the count of tStores
func (index *tagIndex) len() int {
	return len(index.hash2SeriesID)
}

// allTStores returns the map of seriesID and tStores
func (index *tagIndex) allTStores() map[uint32]tStoreINTF {
	return index.seriesID2TStore
}

// flushMetricTo flushes metric-block of mStore to the writer.
func (index *tagIndex) flushMetricTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) error {
	flushed := false
	for _, tStore := range index.seriesID2TStore {
		tStoreDataFlushed := tStore.flushSeriesTo(tableFlusher, flushCtx)
		flushed = flushed || tStoreDataFlushed
	}
	if !flushed {
		return nil
	}
	return tableFlusher.FlushMetric(flushCtx.metricID)
}

// getVersion returns a version(uptime) of the index
func (index *tagIndex) getVersion() int64 {
	return index.version
}

// findSeriesIDsByExpr finds series ids by tag filter expr and timeRange
func (index *tagIndex) findSeriesIDsByExpr(expr stmt.TagFilter, timeRange timeutil.TimeRange) *roaring.Bitmap {
	entrySet, ok := index.getTagKeyEntry(expr.TagKey())
	if !ok {
		return nil
	}
	switch expression := expr.(type) {
	case *stmt.EqualsExpr:
		return index.findSeriesIDsByEqual(entrySet, expression, timeRange)
	case *stmt.InExpr:
		return index.findSeriesIDsByIn(entrySet, expression, timeRange)
	case *stmt.LikeExpr:
		return index.findSeriesIDsByLike(entrySet, expression, timeRange)
	case *stmt.RegexExpr:
		return index.findSeriesIDsByRegex(entrySet, expression, timeRange)
	}
	return nil
}
func (index *tagIndex) findSeriesIDsByEqual(entrySet *tagKVEntrySet, expr *stmt.EqualsExpr,
	timeRange timeutil.TimeRange) *roaring.Bitmap {

	bitmap, ok := entrySet.values[expr.Value]
	if !ok {
		return nil
	}
	union := roaring.New()
	index.unionBitMap(union, bitmap, timeRange)
	return union
}

func (index *tagIndex) findSeriesIDsByIn(entrySet *tagKVEntrySet, expr *stmt.InExpr,
	timeRange timeutil.TimeRange) *roaring.Bitmap {

	union := roaring.New()
	for _, value := range expr.Values {
		bitmap, ok := entrySet.values[value]
		if !ok {
			continue
		}
		index.unionBitMap(union, bitmap, timeRange)
	}
	return union
}

func (index *tagIndex) findSeriesIDsByLike(entrySet *tagKVEntrySet, expr *stmt.LikeExpr,
	timeRange timeutil.TimeRange) *roaring.Bitmap {

	union := roaring.New()
	for value, bitmap := range entrySet.values {
		if strings.Contains(value, expr.Value) {
			index.unionBitMap(union, bitmap, timeRange)
		}
	}
	return union
}
func (index *tagIndex) findSeriesIDsByRegex(entrySet *tagKVEntrySet, expr *stmt.RegexExpr,
	timeRange timeutil.TimeRange) *roaring.Bitmap {

	pattern, err := regexp.Compile(expr.Regexp)
	if err != nil {
		return nil
	}
	union := roaring.New()
	for value, bitmap := range entrySet.values {
		if pattern.MatchString(value) {
			index.unionBitMap(union, bitmap, timeRange)
		}
	}
	return union
}

// getSeriesIDsForTag get series ids by tagKey and timeRange
func (index *tagIndex) getSeriesIDsForTag(tagKey string, timeRange timeutil.TimeRange) *roaring.Bitmap {
	entrySet, ok := index.getTagKeyEntry(tagKey)
	if !ok {
		return nil
	}

	union := roaring.New()
	for _, bitMap := range entrySet.values {
		index.unionBitMap(union, bitMap, timeRange)
	}
	return union
}

// unionBitMap computes the union between two bitmaps and stores the result in the first bitmap.
func (index *tagIndex) unionBitMap(union *roaring.Bitmap, x2 *roaring.Bitmap, timeRange timeutil.TimeRange) {
	iterator := x2.Iterator()
	for iterator.HasNext() {
		seriesID := iterator.Next()
		tStore, ok := index.getTStoreBySeriesID(seriesID)
		if !ok {
			continue
		}
		tRange, ok := tStore.timeRange()
		if !ok {
			continue
		}
		if timeRange.Overlap(&tRange) {
			union.Add(seriesID)
		}
	}
}

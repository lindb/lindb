package indexdb

import (
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

// for testing
var (
	newFlusherFunc = invertedindex.NewFlusher
	newReaderFunc  = invertedindex.NewReader
)

// InvertedIndex represents the tag's inverted index (tag values => series id list)
type InvertedIndex interface {
	// GetSeriesIDsByTagValueIDs returns series ids by tag value ids for spec metric's tag key
	GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(tagKeyIDs []uint32) (series.GroupingContext, error)
	// buildInvertIndex builds the inverted index for tag value => series ids,
	// the tags is considered as a empty key-value pair while tags is nil.
	buildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32)
	// Flush flushes the inverted-index of tag value id=>series ids under tag key
	Flush() error
}

type invertedIndex struct {
	family   kv.Family // store tag value inverted index
	metadata metadb.Metadata

	mutable   *TagIndexStore
	immutable *TagIndexStore

	rwMutex sync.RWMutex
}

func newInvertedIndex(metadata metadb.Metadata, family kv.Family) InvertedIndex {
	return &invertedIndex{
		family:   family,
		metadata: metadata,
		mutable:  NewTagIndexStore(),
	}
}

// FindSeriesIDsByExpr finds series ids by tag filter expr
func (index *invertedIndex) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		seriesIDs := tagIndex.getSeriesIDsByTagValueIDs(tagValueIDs)
		if seriesIDs != nil {
			result.Or(seriesIDs)
		}
	})

	// read data from kv store
	if err := index.loadSeriesIDsInKV(tagKeyID, func(reader invertedindex.Reader) error {
		seriesIDs, err := reader.FindSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
		if err != nil {
			return err
		}
		result.Or(seriesIDs)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

// GetSeriesIDsForTag get series ids by tagKeyId
func (index *invertedIndex) GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		result.Or(tagIndex.getAllSeriesIDs())
	})

	// read data from kv store
	if err := index.loadSeriesIDsInKV(tagKeyID, func(reader invertedindex.Reader) error {
		seriesIDs, err := reader.GetSeriesIDsForTagKeyID(tagKeyID)
		if err != nil {
			return err
		}
		result.Or(seriesIDs)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (index *invertedIndex) GetGroupingContext(tagKeyIDs []uint32) (series.GroupingContext, error) {
	tagKeysLen := len(tagKeyIDs)
	gCtx := query.NewGroupContext(tagKeysLen)
	// validate tagKeys
	for idx, tagKeyID := range tagKeyIDs {
		_, ok := index.mutable.Get(tagKeyID)
		if !ok {
			return nil, constants.ErrNotFound
		}
		tagValuesEntrySet := query.NewTagValuesEntrySet()
		gCtx.SetTagValuesEntrySet(idx, tagValuesEntrySet)
		//FIXME stone1100
		//tagValuesEntrySet.SetTagValues(tagIndex.getValues())
	}
	return &groupingContext{
		gCtx: gCtx,
	}, nil
}

// buildInvertIndex builds the inverted index for tag value => series ids,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *invertedIndex) buildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32) {
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()

	metadataDB := index.metadata.MetadataDatabase()
	tagMetadata := index.metadata.TagMetadata()
	for tagKey, tagValue := range tags {
		tagKeyID, err := metadataDB.GenTagKeyID(namespace, metricName, tagKey)
		if err != nil {
			//FIXME stone1100 add metric???
			indexLogger.Error("gen tag key id fail, ignore index build for this tag key",
				logger.String("key", tagKey), logger.Error(err))
			continue
		}
		tagIndex, ok := index.mutable.Get(tagKeyID)
		if !ok {
			tagIndex = newTagIndex()
			index.mutable.Put(tagKeyID, tagIndex)
		}
		tagValueID, err := tagMetadata.GenTagValueID(tagKeyID, tagValue)
		if err != nil {
			//FIXME stone1100 add metric???
			indexLogger.Error("gen tag value id fail, ignore index build for this tag key",
				logger.String("key", tagKey), logger.String("value", tagValue), logger.Error(err))
			continue
		}
		tagIndex.buildInvertedIndex(tagValueID, seriesID)
	}
}

// Flush flushes the inverted-index of tag value id=>series ids under tag key
func (index *invertedIndex) Flush() error {
	if !index.checkFlush() {
		return nil
	}

	// flush immutable data into kv store
	flusher := index.family.NewFlusher()
	indexFlusher := newFlusherFunc(flusher)
	if err := index.immutable.WalkEntry(func(key uint32, value TagIndex) error {
		if err := value.flush(indexFlusher); err != nil {
			return err
		}
		if err := indexFlusher.FlushTagKeyID(key); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	// commit kv stone meta
	if err := indexFlusher.Commit(); err != nil {
		return err
	}
	// finally clear immutable
	index.rwMutex.Lock()
	index.immutable = nil
	index.rwMutex.Unlock()
	return nil
}

// checkFlush checks if need do flush job, if need, do switch mutable/immutable
func (index *invertedIndex) checkFlush() bool {
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()

	if index.mutable.Size() == 0 && index.immutable == nil {
		// no new data or immutable is not nil
		return false
	}
	if index.mutable.Size() > 0 && index.immutable == nil {
		// reset mutable, if flush fail immutable is not nil
		index.immutable = index.mutable
		index.mutable = NewTagIndexStore()
	}
	return true
}

// loadTagValueIDsInKV loads series ids in kv store
func (index *invertedIndex) loadSeriesIDsInKV(tagKeyID uint32, fn func(reader invertedindex.Reader) error) error {
	// try get tag key id from kv store
	snapshot := index.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return err
	}
	var reader invertedindex.Reader
	if len(readers) > 0 {
		// found tag data in kv store, try load series ids data
		reader = newReaderFunc(readers)
		if err := fn(reader); err != nil {
			return err
		}
	}
	return nil
}

// loadSeriesIDsInMem loads series ids from mutable/immutable store
func (index *invertedIndex) loadSeriesIDsInMem(tagKeyID uint32, fn func(tagIndex TagIndex)) {
	// define get tag series ids func
	getSeriesIDsIDs := func(tagIndexStore *TagIndexStore) {
		tag, ok := tagIndexStore.Get(tagKeyID)
		if ok {
			fn(tag)
		}
	}

	// read data with read lock
	index.rwMutex.RLock()
	defer index.rwMutex.RUnlock()

	getSeriesIDsIDs(index.mutable)
	if index.immutable != nil {
		getSeriesIDsIDs(index.immutable)
	}
}

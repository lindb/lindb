package indexdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

// InvertedIndex represents the tag's inverted index (tag values => series id list)
type InvertedIndex interface {
	series.TagValueSuggester
	// FindSeriesIDsByExpr finds series ids by tag filter expr for tag key id
	FindSeriesIDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(tagKeyIDs []uint32) (series.GroupingContext, error)
	buildInvertIndex(metricID uint32, tags map[string]string, seriesID uint32)
}

type invertedIndex struct {
	store     *tagIndexStore
	generator metadb.IDGenerator
}

func newInvertedIndex(generator metadb.IDGenerator) InvertedIndex {
	return &invertedIndex{
		generator: generator,
		store:     newTagIndexStore(),
	}
}

// FindSeriesIDsByExpr finds series ids by tag filter expr
func (index *invertedIndex) FindSeriesIDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	tagIndex, ok := index.store.get(tagKeyID)
	if !ok {
		return nil, constants.ErrNotFound
	}
	return tagIndex.findSeriesIDsByExpr(expr), nil
}

// GetSeriesIDsForTag get series ids by tagKeyId
func (index *invertedIndex) GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	tagIndex, ok := index.store.get(tagKeyID)
	if !ok {
		return nil, constants.ErrNotFound
	}
	return tagIndex.getAllSeriesIDs(), nil
}

func (index *invertedIndex) GetGroupingContext(tagKeyIDs []uint32) (series.GroupingContext, error) {
	tagKeysLen := len(tagKeyIDs)
	gCtx := query.NewGroupContext(tagKeysLen)
	// validate tagKeys
	for idx, tagKeyID := range tagKeyIDs {
		tagIndex, ok := index.store.get(tagKeyID)
		if !ok {
			return nil, constants.ErrNotFound
		}
		tagValuesEntrySet := query.NewTagValuesEntrySet()
		gCtx.SetTagValuesEntrySet(idx, tagValuesEntrySet)
		tagValuesEntrySet.SetTagValues(tagIndex.getValues())
	}
	return &groupingContext{
		gCtx: gCtx,
	}, nil
}

func (index *invertedIndex) SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string {
	tagIndex, ok := index.store.get(tagKeyID)
	if !ok {
		return nil
	}
	return tagIndex.suggestTagValues(tagValuePrefix, limit)
}

// buildInvertIndex builds the inverted index for tag value => series ids,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *invertedIndex) buildInvertIndex(metricID uint32, tags map[string]string, seriesID uint32) {
	for tagKey, tagValue := range tags {
		tagKeyID := index.generator.GenTagKeyID(metricID, tagKey)
		tagIndex, ok := index.store.get(tagKeyID)
		if !ok {
			tagIndex = newTagIndex()
			index.store.put(tagKeyID, tagIndex)
		}
		tagIndex.buildInvertedIndex(tagValue, seriesID)
	}
}

// FlushInvertedIndexTo flushes the inverted-index of mStore to the Writer
func (index *invertedIndex) FlushInvertedIndexTo(flusher invertedindex.TagFlusher) error {
	//seriesIDBitmap := index.store.tagKeyIDs
	//for idx, highKey := range seriesIDBitmap.GetHighKeys() {
	//	container := seriesIDBitmap.GetContainer(highKey)
	//	tagIndexes := index.store.indexes[idx]
	//	it := container.PeekableIterator()
	//	i := 0
	//	for it.HasNext() {
	//		lowKeyID := it.Next()
	//		tagIndex := tagIndexes[i]
	//		tagValues := tagIndex.getValues()
	//		for tagValue, seriesIDs := range tagValues {
	//			flusher.FlushTagValue(tagValue, seriesIDs)
	//		}
	//		tagKeyID := uint32(lowKeyID) | uint32(highKey)
	//		if err := flusher.FlushTagKeyID(tagKeyID); err != nil {
	//			return err
	//		}
	//	}
	//}
	return nil
}

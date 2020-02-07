package metadb

import (
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/metricsmeta"
)

//go:generate mockgen -source ./tag_metadata.go -destination=./tag_metadata_mock.go -package metadb

// for testing
var (
	newTagReaderFunc  = metricsmeta.NewTagReader
	newTagFlusherFunc = metricsmeta.NewTagFlusher
)

// TagMetadata represents the tag metadata, stores all tag values under spec tag key
type TagMetadata interface {
	// GenTagValueID generates the tag value id for spec tag key
	GenTagValueID(tagKeyID uint32, tagValue string) (uint32, error)

	// SuggestTagValues returns suggestions from given tag key id and prefix of tag value
	SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string

	// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	FindTagValueDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error)
	// GetTagValueIDsForTag get tag value ids for spec metric's tag key,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	GetTagValueIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)

	// Flush flushes the memory tag metadata into kv store
	Flush() error
}

// tagMetadata implements TagMetadata interface
type tagMetadata struct {
	family    kv.Family // store tag key/value data using common kv store
	mutable   *TagStore // mutable store current writeable memory store
	immutable *TagStore // immutable need to flush into kv store

	rwMutex sync.RWMutex
}

// NewTagMetadata creates a tag metadata
func NewTagMetadata(family kv.Family) TagMetadata {
	m := &tagMetadata{
		family:  family,
		mutable: NewTagStore(),
	}
	return m
}

// GenTagValueID generates the tag value id for spec tag key
func (m *tagMetadata) GenTagValueID(tagKeyID uint32, tagValue string) (tagValueID uint32, err error) {
	// get tag value id from memory with read lock
	m.rwMutex.RLock()
	tagValueID, ok := m.getTagValueIDInMem(tagKeyID, tagValue)
	if ok {
		m.rwMutex.RUnlock()
		return tagValueID, nil
	}
	m.rwMutex.RUnlock()

	// try load tag value id from kv store
	snapshot := m.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return
	}
	var reader metricsmeta.TagReader
	if len(readers) > 0 {
		// found tag data in kv store, try load tag value data
		reader = newTagReaderFunc(readers)
		tagValueID, err = reader.GetTagValueID(tagKeyID, tagValue)
		if err == nil {
			// got tag value id from kv store
			return tagValueID, nil
		}
		if err != constants.ErrNotFound {
			// if load tag value id err, return it
			return
		}
	}

	// tag value not exist, need assign new tag value id for tag value with write lock
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	// double check, memory if exist tag value
	tagValueID, ok = m.getTagValueIDInMem(tagKeyID, tagValue)
	if ok {
		return tagValueID, nil
	}

	// assign new tag value id
	tag, ok := m.mutable.Get(tagKeyID)
	if !ok {
		if reader != nil {
			// if tag data exist in kv store, need load tag value id auto sequence
			seq, err := reader.GetTagValueSeq(tagKeyID)
			if err != nil {
				return 0, err
			}
			tag = newTagEntry(seq)
		} else {
			// for new tag, auto sequence start with 0
			tag = newTagEntry(0)
		}
		// cache tag entry
		m.mutable.Put(tagKeyID, tag)
	}

	// assign new id
	tagValueID = tag.genTagValueID()
	tag.addTagValue(tagValue, tagValueID)
	return tagValueID, nil
}

// SuggestTagValues returns suggestions from given tag key id and prefix of tag value
func (m *tagMetadata) SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string {
	//FIXME stone1100
	panic("implement me")
}

// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (m *tagMetadata) FindTagValueDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	result := roaring.New()
	result.Or(m.findTagValueIDsByExprTagInMem(tagKeyID, expr))

	err := m.loadTagValueIDsInKV(tagKeyID, func(reader metricsmeta.TagReader) error {
		tagValueIDs, err := reader.FindValueIDsByExprForTagKeyID(tagKeyID, expr)
		if err != nil {
			return err
		}
		result.Or(tagValueIDs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTagValueIDsForTag get tag value ids for spec metric's tag key,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (m *tagMetadata) GetTagValueIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	result := roaring.New()
	result.Or(m.getTagValueIDsForTagInMem(tagKeyID))

	err := m.loadTagValueIDsInKV(tagKeyID, func(reader metricsmeta.TagReader) error {
		tagValueIDs, err := reader.GetTagValueIDsForTagKeyID(tagKeyID)
		if err != nil {
			return err
		}
		result.Or(tagValueIDs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Flush flushes the memory tag metadata into kv store
func (m *tagMetadata) Flush() error {
	m.rwMutex.Lock()
	if m.mutable.Size() == 0 && m.immutable == nil {
		// no new data or immutable is not nil
		m.rwMutex.Unlock()
		return nil
	}
	if m.mutable.Size() > 0 && m.immutable == nil {
		// reset mutable, if flush fail immutable is not nil
		m.immutable = m.mutable
		m.mutable = NewTagStore()
	}
	m.rwMutex.Unlock()

	// flush immutable data into kv store
	fluster := m.family.NewFlusher()
	tagFluster := newTagFlusherFunc(fluster)
	keys := m.immutable.Keys()
	values := m.immutable.Values()
	highKeys := keys.GetHighKeys()
	for highIdx, highKey := range highKeys {
		hk := uint32(highKey) << 16
		lowValues := values[highIdx]
		lowContainer := keys.GetContainerAtIndex(highIdx)
		it := lowContainer.PeekableIterator()
		idx := 0
		for it.HasNext() {
			lowKey := it.Next()
			tag := lowValues[idx]
			idx++
			tagValues := tag.getTagValues()
			for tagValue, tagValueID := range tagValues {
				tagFluster.FlushTagValue(tagValue, tagValueID)
			}
			if err := tagFluster.FlushTagKeyID(uint32(lowKey&0xFFFF)|hk, tag.getTagValueIDSeq()); err != nil {
				return err
			}
		}
	}
	if err := tagFluster.Commit(); err != nil {
		return err
	}
	// finally clear immutable
	m.rwMutex.Lock()
	m.immutable = nil
	m.rwMutex.Unlock()
	return nil
}

// findTagValueIDsByExprTagInMem finds tag value ids by expr from mutable/immutable store
func (m *tagMetadata) findTagValueIDsByExprTagInMem(tagKeyID uint32, expr stmt.TagFilter) *roaring.Bitmap {
	result := roaring.New()

	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		ids := tagEntry.findSeriesIDsByExpr(expr)
		if ids != nil {
			result.Or(ids)
		}
	})
	return result
}

// getTagValueIDsForTagInMem gets tag value ids from mutable/immutable store
func (m *tagMetadata) getTagValueIDsForTagInMem(tagKeyID uint32) *roaring.Bitmap {
	result := roaring.New()

	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		ids := tagEntry.getTagValueIDs()
		if ids != nil {
			result.Or(ids)
		}
	})
	return result
}

func (m *tagMetadata) loadTagValueIDsInKV(tagKeyID uint32, fn func(reader metricsmeta.TagReader) error) error {
	// try load tag value id from kv store
	snapshot := m.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return err
	}
	var reader metricsmeta.TagReader
	if len(readers) > 0 {
		// found tag data in kv store, try load tag value data
		reader = newTagReaderFunc(readers)
		if err := fn(reader); err != nil {
			return err
		}
	}
	return nil
}

// loadTagValueIDsInMem loads tag value ids from mutable/immutable store
func (m *tagMetadata) loadTagValueIDsInMem(tagKeyID uint32, fn func(tagEntry TagEntry)) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	// define get tag value ids func
	getTagValueIDs := func(tagStore *TagStore) {
		tag, ok := tagStore.Get(tagKeyID)
		if ok {
			fn(tag)
		}
	}

	getTagValueIDs(m.mutable)
	if m.immutable != nil {
		getTagValueIDs(m.immutable)
	}
}

// getTagValueIDInMem gets tag value id from mutable/immutable store
func (m *tagMetadata) getTagValueIDInMem(tagKeyID uint32, tagValue string) (tagValueID uint32, ok bool) {
	tagValueID, ok = getTagValueID(m.mutable, tagKeyID, tagValue)
	if ok {
		return
	}
	if m.immutable != nil {
		tagValueID, ok = getTagValueID(m.immutable, tagKeyID, tagValue)
		if ok {
			return
		}
	}
	return
}

// getTagValueID gets tag value id from tag store based on tag key id and tag value
func getTagValueID(tags *TagStore, tagKeyID uint32, tagValue string) (tagValueID uint32, ok bool) {
	tag, ok := tags.Get(tagKeyID)
	if ok {
		tagValueID, ok = tag.getTagValueID(tagValue)
		if ok {
			return tagValueID, true
		}
	}
	return
}

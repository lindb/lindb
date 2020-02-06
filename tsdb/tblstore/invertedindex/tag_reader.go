package invertedindex

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/sql/stmt"
)

var invertedIndexReaderLogger = logger.GetLogger("tsdb", "InvertedIndexReader")

//go:generate mockgen -source ./tag_reader.go -destination=./tag_reader_mock.go -package invertedindex

// for testing
var (
	newTagKVEntrySetFunc = newTagKVEntrySet
)

const (
	tagFooterSize = 4 + // tag value id sequence
		4 + // tag value ids position
		4 // crc32 checksum
)

// TagReader reads tag value data from tag-index-table
type TagReader interface {
	// GetTagValueSeq returns the auto sequence of tag value under the tag key, if not exist return constants.ErrNotFound
	GetTagValueSeq(tagKeyID uint32) (tagValueSeq uint32, err error)

	// GetTagValueID returns the tag value id for spec metric's tag key id, if not exist return constants.ErrNotFound
	GetTagValueID(tagKeyID uint32, tagValue string) (tagValueID uint32, err error)

	// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
	GetTagValueIDsForTagKeyID(tagKeyID uint32) (tagValueIDs *roaring.Bitmap, err error)

	// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
	FindValueIDsByExprForTagKeyID(tagKeyID uint32, expr stmt.TagFilter) (tagValueIDs *roaring.Bitmap, err error)

	// SuggestTagValues finds tag values by prefix search
	SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string

	// WalkTagValues walks each tag value and bitmap via fn.
	// If fn returns false, the iteration is stopped.
	// The values are the raw byte slices and not the converted types.
	WalkTagValues(
		tagID uint32,
		tagValuePrefix string,
		fn func(tagValue []byte, tagValueID uint32) bool,
	) error
}

// tagReader implements TagReader
type tagReader struct {
	readers []table.Reader
}

// NewReader returns a new TagReader
func NewTagReader(readers []table.Reader) TagReader {
	return &tagReader{readers: readers}
}

// GetTagValueSeq returns the auto sequence of tag value under the tag key, if not exist return constants.ErrNotFound
// kv store returns the table.readers in order,
// so the max sequence will be stored in the first table.reader that is tag key store.
func (r *tagReader) GetTagValueSeq(tagKeyID uint32) (tagValueSeq uint32, err error) {
	for _, reader := range r.readers {
		value, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		entrySet, err := newTagKVEntrySetFunc(value)
		if err != nil {
			return 0, err
		}
		return entrySet.TagValueSeq(), nil
	}
	return 0, constants.ErrNotFound
}

// GetTagValueID returns the tag value id for spec metric's tag key id, if not exist return constants.ErrNotFound
func (r *tagReader) GetTagValueID(tagID uint32, tagValue string) (tagValueID uint32, err error) {
	for _, reader := range r.readers {
		value, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		entrySet, err := newTagKVEntrySetFunc(value)
		if err != nil {
			return 0, err
		}
		q, err := entrySet.TrieTree()
		if err != nil {
			return 0, err
		}
		offsets := q.FindOffsetsByEqual(tagValue)
		if len(offsets) == 0 {
			continue
		}
		if len(offsets) > 1 {
			return 0, fmt.Errorf("found too many offsets for tag value")
		}
		return entrySet.GetTagValueID(offsets[0]), nil
	}
	return 0, constants.ErrNotFound
}

// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
func (r *tagReader) FindValueIDsByExprForTagKeyID(tagID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	entrySets := r.filterEntrySets(tagID)
	if len(entrySets) == 0 {
		return nil, constants.ErrNotFound
	}
	tagValueIDs := roaring.New()
	for _, entrySet := range entrySets {
		var offsets []int
		q, err := entrySet.TrieTree()
		if err != nil {
			invertedIndexReaderLogger.Error("failed reading trie-tree block", logger.Error(err))
			continue
		}
		switch expression := expr.(type) {
		case *stmt.EqualsExpr:
			offsets = q.FindOffsetsByEqual(expression.Value)
		case *stmt.InExpr:
			offsets = q.FindOffsetsByIn(expression.Values)
		case *stmt.LikeExpr:
			offsets = q.FindOffsetsByLike(expression.Value)
		case *stmt.RegexExpr:
			offsets = q.FindOffsetsByRegex(expression.Regexp)
		default:
			return nil, constants.ErrNotFound
		}
		if len(offsets) == 0 {
			continue
		}
		for _, offset := range offsets {
			tagValueIDs.Add(entrySet.GetTagValueID(offset))
		}
	}
	if tagValueIDs.IsEmpty() {
		return nil, constants.ErrNotFound
	}
	return tagValueIDs, nil
}

// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
func (r *tagReader) GetTagValueIDsForTagKeyID(tagID uint32) (*roaring.Bitmap, error) {
	entrySets := r.filterEntrySets(tagID)
	if len(entrySets) == 0 {
		return nil, constants.ErrNotFound
	}
	return entrySets.GetTagValueIDs(), nil
}

// filterEntrySets filters the entry-sets by tag key id
func (r *tagReader) filterEntrySets(tagID uint32) (entrySets TagKVEntries) {
	for _, reader := range r.readers {
		value, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		entrySet, err := newTagKVEntrySetFunc(value)
		if err != nil {
			continue
		}
		entrySets = append(entrySets, entrySet)
	}
	return
}

// SuggestTagValues finds tagValues by prefix search
func (r *tagReader) SuggestTagValues(
	tagID uint32,
	tagValuePrefix string,
	limit int,
) (
	tagValues []string,
) {
	if limit > constants.MaxSuggestions {
		limit = constants.MaxSuggestions
	}
	for _, reader := range r.readers {
		value, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		entrySet, err := newTagKVEntrySetFunc(value)
		if err != nil {
			continue
		}
		q, err := entrySet.TrieTree()
		if err != nil {
			invertedIndexReaderLogger.Error("failed reading trie-tree block", logger.Error(err))
			continue
		}
		tagValues = append(tagValues, q.PrefixSearch(tagValuePrefix, limit-len(tagValues))...)
		if len(tagValues) >= limit {
			return tagValues
		}
	}
	return tagValues
}

func (r *tagReader) WalkTagValues(
	tagID uint32,
	tagValuePrefix string,
	fn func(tagValue []byte, tagValueID uint32) bool,
) error {
	for _, reader := range r.readers {
		value, ok := reader.Get(tagID)
		if !ok {
			continue
		}
		entrySet, err := newTagKVEntrySetFunc(value)
		if err != nil {
			continue
		}
		q, err := entrySet.TrieTree()
		if err != nil {
			continue
		}
		offsetsItr := q.Iterator(tagValuePrefix)
		for offsetsItr.HasNext() {
			tagValue, offset := offsetsItr.Next()
			if fn != nil && !fn(tagValue, entrySet.GetTagValueID(offset)) {
				return nil
			}
		}
	}
	return nil
}

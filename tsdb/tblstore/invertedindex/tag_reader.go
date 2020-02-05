package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

var invertedIndexReaderLogger = logger.GetLogger("tsdb", "InvertedIndexReader")

//go:generate mockgen -source ./tag_reader.go -destination=./tag_reader_mock.go -package invertedindex

// for testing
var (
	newTagKVEntrySetFunc = newTagKVEntrySet
)

const (
	tagFooterSize = 4 + // tag value ids position
		4 // crc32 checksum
)

// TagReader reads tag value data from tag-index-table
type TagReader interface {
	// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
	GetTagValueIDsForTagKeyID(tagID uint32) (tagValueIDs *roaring.Bitmap, err error)

	// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
	FindValueIDsByExprForTagKeyID(tagID uint32, expr stmt.TagFilter) (tagValueIDs *roaring.Bitmap, err error)

	// SuggestTagValues finds tag values by prefix search
	SuggestTagValues(tagID uint32, tagValuePrefix string, limit int) []string

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

// FindValueIDsByExprForTagKeyID finds tag values ids by tag filter expr and tag key id
func (r *tagReader) FindValueIDsByExprForTagKeyID(tagID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	entrySets := r.filterEntrySets(tagID)
	if len(entrySets) == 0 {
		return nil, series.ErrNotFound
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
			return nil, series.ErrNotFound
		}
		if len(offsets) == 0 {
			continue
		}
		for _, offset := range offsets {
			tagValueIDs.Add(entrySet.GetTagValueID(offset))
		}
	}
	if tagValueIDs.IsEmpty() {
		return nil, series.ErrNotFound
	}
	return tagValueIDs, nil
}

// GetTagValueIDsForTagKeyID get tag value ids for spec metric's tag key id
func (r *tagReader) GetTagValueIDsForTagKeyID(tagID uint32) (*roaring.Bitmap, error) {
	entrySets := r.filterEntrySets(tagID)
	if len(entrySets) == 0 {
		return nil, series.ErrNotFound
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

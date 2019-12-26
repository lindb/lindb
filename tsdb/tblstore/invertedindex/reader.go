package invertedindex

import (
	"sort"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/lindb/roaring"
)

var invertedIndexReaderLogger = logger.GetLogger("tsdb", "InvertedIndexReader")

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package invertedindex

const (
	invertedIndexTimeRangeSize = 8 + // int64, start-time
		8 // int64, end-time
	invertedIndexFooterSize = 4 + // offsets position
		4 // crc32 checksum
)

// Reader reads versioned seriesID bitmap from series-index-table
type Reader interface {
	// GetSeriesIDsForTagKeyID get series ids for spec metric's keyID
	GetSeriesIDsForTagKeyID(
		tagID uint32,
		timeRange timeutil.TimeRange,
	) (
		*series.MultiVerSeriesIDSet,
		error)

	// FindSeriesIDsByExprForTagKeyID finds series ids by tag filter expr and tagKeyID
	FindSeriesIDsByExprForTagKeyID(
		tagID uint32, expr stmt.TagFilter,
		timeRange timeutil.TimeRange,
	) (
		*series.MultiVerSeriesIDSet,
		error)

	// SuggestTagValues finds tagValues by prefix search
	SuggestTagValues(
		tagID uint32,
		tagValuePrefix string,
		limit int,
	) []string

	// WalkTagValues walks each tag value and bitmap via fn.
	// If fn returns false, the iteration is stopped.
	// The values are the raw byte slices and not the converted types.
	WalkTagValues(
		tagID uint32,
		tagValuePrefix string,
		fn func(tagValue []byte, dataIterator TagValueIterator) bool,
	) error
}

// reader implements Reader
type reader struct {
	readers []table.Reader
}

// NewReader returns a new Reader
func NewReader(readers []table.Reader) Reader {
	return &reader{readers: readers}
}

// FindSeriesIDsByExprForTagKeyID finds series ids by tag filter expr for tagId
func (r *reader) FindSeriesIDsByExprForTagKeyID(
	tagID uint32,
	expr stmt.TagFilter,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	entrySets := r.filterEntrySets(tagID, timeRange)
	if len(entrySets) == 0 {
		return nil, series.ErrNotFound
	}
	unionIDSet := series.NewMultiVerSeriesIDSet()
	for _, entrySet := range entrySets {
		var offsets []int
		q, err := entrySet.TrieTree()
		if err != nil {
			invertedIndexReaderLogger.Error("failed reading trie-tree block", logger.Error(err))
			continue
		}
		switch expression := expr.(type) {
		case *stmt.EqualsExpr:
			offsets = append(offsets, q.FindOffsetsByEqual(expression.Value)...)
		case *stmt.InExpr:
			offsets = append(offsets, q.FindOffsetsByIn(expression.Values)...)
		case *stmt.LikeExpr:
			offsets = append(offsets, q.FindOffsetsByLike(expression.Value)...)
		case *stmt.RegexExpr:
			offsets = append(offsets, q.FindOffsetsByRegex(expression.Regexp)...)
		default:
			return nil, series.ErrNotFound
		}
		if len(offsets) == 0 {
			continue
		}
		idSet, err := r.entrySetToIDSet(entrySet, timeRange, offsets)
		if err != nil {
			return nil, err
		}
		unionIDSet.Or(idSet)
	}
	if unionIDSet.IsEmpty() {
		return nil, series.ErrNotFound
	}
	return unionIDSet, nil
}

// GetSeriesIDsForTagKeyID get series ids for spec metric's tag keyID
func (r *reader) GetSeriesIDsForTagKeyID(
	tagID uint32,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	entrySets := r.filterEntrySets(tagID, timeRange)
	if len(entrySets) == 0 {
		return nil, series.ErrNotFound
	}
	unionIDSet := series.NewMultiVerSeriesIDSet()
	for _, entrySet := range entrySets {
		idSet, err := r.entrySetToIDSet(entrySet, timeRange, nil)
		if err != nil {
			return nil, err
		}
		unionIDSet.Or(idSet)
	}
	return unionIDSet, nil
}

// filterEntrySets filters the entry-sets which matches the time-range in the series-index-table
func (r *reader) filterEntrySets(
	tagID uint32,
	timeRange timeutil.TimeRange,
) (
	entrySets []tagKVEntrySetINTF,
) {
	for _, reader := range r.readers {
		entrySet, err := newTagKVEntrySet(reader.Get(tagID))
		if err != nil {
			continue
		}
		entrySetTimeRange := entrySet.TimeRange()
		if !timeRange.Overlap(&entrySetTimeRange) {
			continue
		}
		entrySets = append(entrySets, entrySet)
	}
	return
}

// entrySetToIDSet parses the entry-set block, then return the multi-versions seriesID bitmap
func (r *reader) entrySetToIDSet(
	entrySet tagKVEntrySetINTF,
	timeRange timeutil.TimeRange,
	offsets []int,
) (
	idSet *series.MultiVerSeriesIDSet,
	err error,
) {
	bitmap := roaring.New()

	retrieve := func(offset int) error {
		dataItr, err := entrySet.ReadTagValueDataBlock(offset)
		if err != nil {
			return err
		}
		for dataItr.HasNext() {
			dataTimeRange := dataItr.DataTimeRange()
			if !timeRange.Overlap(&dataTimeRange) {
				return nil
			}
			if err := bitmap.UnmarshalBinary(dataItr.Next()); err != nil {
				return err
			}
			if idSet == nil {
				idSet = series.NewMultiVerSeriesIDSet()
			}
			theBitMap, ok := idSet.Versions()[dataItr.DataVersion()]
			if ok {
				theBitMap.Or(bitmap)
			} else {
				theBitMap = bitmap.Clone()
			}
			idSet.Add(dataItr.DataVersion(), theBitMap)
		}
		return nil
	}
	// nil means reads data at all positions
	if len(offsets) == 0 {
		posItr := entrySet.PositionIterator()
		for posItr.HasNext() {
			offset, _ := posItr.Next()
			if err = retrieve(offset); err != nil {
				return nil, err
			}
		}
	} else {
		sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
		for _, offset := range offsets {
			if err = retrieve(offset); err != nil {
				return nil, err
			}
		}
	}
	if idSet == nil {
		return nil, series.ErrNotFound
	}
	return idSet, nil
}

// SuggestTagValues finds tagValues by prefix search
func (r *reader) SuggestTagValues(
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
		entrySet, err := newTagKVEntrySet(reader.Get(tagID))
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

func (r *reader) WalkTagValues(
	tagID uint32,
	tagValuePrefix string,
	fn func(tagValue []byte, dataIterator TagValueIterator) bool,
) error {
	for _, reader := range r.readers {
		entrySet, err := newTagKVEntrySet(reader.Get(tagID))
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
			dataItr, err := entrySet.ReadTagValueDataBlock(offset)
			if err != nil {
				return err
			}
			if !fn(tagValue, dataItr) {
				return nil
			}
		}
	}
	return nil
}

package invertedindex

import (
	"fmt"
	"math"
	"sort"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
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
	entrySets []*tagKVEntrySet,
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
	entrySet *tagKVEntrySet,
	timeRange timeutil.TimeRange,
	offsets []int,
) (
	idSet *series.MultiVerSeriesIDSet,
	err error,
) {
	positions, err := entrySet.OffsetsToPosition(offsets)
	if err != nil {
		return nil, err
	}
	// sort positions for continuously read
	var (
		count         = 0
		positionsList = make([]int, len(positions))
	)
	for _, position := range positions {
		positionsList[count] = position
		count++
	}
	sort.Slice(positionsList, func(i, j int) bool { return positionsList[i] < positionsList[j] })
	// read in order
	for _, position := range positions {
		tagValueData, err := entrySet.ReadTagValueDataBlock(position)
		if err != nil {
			return nil, err
		}
		for _, data := range tagValueData {
			dataTimeRange := data.TimeRange()
			if !timeRange.Overlap(&dataTimeRange) {
				continue
			}
			bitmap, err := data.Bitmap()
			if err != nil {
				return nil, err
			}
			if idSet == nil {
				idSet = series.NewMultiVerSeriesIDSet()
			}
			theBitMap, ok := idSet.Versions()[data.version]
			if ok {
				theBitMap.Or(bitmap)
			} else {
				theBitMap = bitmap
			}
			idSet.Add(data.version, theBitMap)
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

type tagKVEntrySet struct {
	sr            *stream.Reader
	startTime     int64
	endTime       int64
	tree          trieTreeQuerier
	offsetsBlock  []byte
	crc32CheckSum uint32
	// buffer
	tagValueData []versionedTagValueData
}

func newTagKVEntrySet(block []byte) (*tagKVEntrySet, error) {
	if len(block) <= invertedIndexTimeRangeSize+invertedIndexFooterSize {
		return nil, fmt.Errorf("block length no ok")
	}
	entrySet := &tagKVEntrySet{
		sr: stream.NewReader(block)}
	entrySet.startTime = entrySet.sr.ReadInt64()
	entrySet.endTime = entrySet.sr.ReadInt64()
	// read footer
	offsetsEndPos := len(block) - invertedIndexFooterSize
	_ = entrySet.sr.ReadSlice(offsetsEndPos - invertedIndexTimeRangeSize)
	offsetsStartPos := int(entrySet.sr.ReadUint32())
	entrySet.crc32CheckSum = entrySet.sr.ReadUint32()
	// validate offsets
	if !(invertedIndexTimeRangeSize < offsetsStartPos && offsetsStartPos < offsetsEndPos) {
		return nil, fmt.Errorf("bad offsets")
	}
	entrySet.offsetsBlock = block[offsetsStartPos:offsetsEndPos]
	return entrySet, nil
}

// TimeRange computes the timeRange from delta in seconds
func (entrySet *tagKVEntrySet) TimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{
		Start: entrySet.startTime,
		End:   entrySet.endTime}
}

// TrieTree builds the trie-tree block for querying
func (entrySet *tagKVEntrySet) TrieTree() (trieTreeQuerier, error) {
	var tree trieTreeBlock
	entrySet.sr.SeekStart()
	// read time-range
	_ = entrySet.sr.ReadSlice(invertedIndexTimeRangeSize)
	////////////////////////////////
	// Block: LOUDS Trie-Tree
	////////////////////////////////
	// read trie-tree length
	expectedTrieTreeLen := entrySet.sr.ReadUvarint64()
	startPosOfTree := entrySet.sr.Position()
	// read label length
	labelsLen := entrySet.sr.ReadUvarint64()
	// read labels block
	tree.labels = entrySet.sr.ReadSlice(int(labelsLen))
	// read isPrefix length
	isPrefixKeyLen := entrySet.sr.ReadUvarint64()
	// read isPrefixKey bitmap
	isPrefixBlock := entrySet.sr.ReadSlice(int(isPrefixKeyLen))
	// read LOUDS length
	loudsLen := entrySet.sr.ReadUvarint64()
	// read LOUDS block
	LOUDSBlock := entrySet.sr.ReadSlice(int(loudsLen))
	// validation of stream error
	if entrySet.sr.Error() != nil {
		return nil, entrySet.sr.Error()
	}
	// validation of length
	if entrySet.sr.Position()-startPosOfTree != int(expectedTrieTreeLen) {
		return nil, fmt.Errorf("failed validation of trie-tree")
	}
	// unmarshal LOUDS block to rank-select
	tree.LOUDS = NewRankSelect()
	if err := tree.LOUDS.UnmarshalBinary(LOUDSBlock); err != nil {
		return nil, err
	}
	// unmarshal isPrefixKey block to rank-select
	tree.isPrefixKey = NewRankSelect()
	if err := tree.isPrefixKey.UnmarshalBinary(isPrefixBlock); err != nil {
		return nil, err
	}
	entrySet.tree = &tree
	return entrySet.tree, nil
}

// ReadTagValueDataBlock reads tagValueDataBlocks at specified position
func (entrySet *tagKVEntrySet) ReadTagValueDataBlock(
	pos int,
) (
	[]versionedTagValueData,
	error,
) {
	// jump to target
	entrySet.sr.SeekStart()
	_ = entrySet.sr.ReadSlice(pos)
	// clear data
	entrySet.tagValueData = entrySet.tagValueData[:0]
	// reset buffer
	versionCount := entrySet.sr.ReadUvarint64()
	var counter = 0
	for !entrySet.sr.Empty() && counter < int(versionCount) {
		// read version
		version := series.Version(entrySet.sr.ReadInt64())
		// read start-time delta
		startTimeDelta := entrySet.sr.ReadVarint64()
		// read end-time delta
		endTimeDelta := entrySet.sr.ReadVarint64()
		// read bitmap length
		bitMapLen := int(entrySet.sr.ReadUvarint64())
		// read bitmap
		bitMapBlock := entrySet.sr.ReadSlice(bitMapLen)
		if entrySet.sr.Error() != nil {
			break
		}
		entrySet.tagValueData = append(entrySet.tagValueData, versionedTagValueData{
			version:        version,
			startTimeDelta: startTimeDelta,
			endTimeDelta:   endTimeDelta,
			bitMapData:     bitMapBlock})
		counter++
	}
	return entrySet.tagValueData, entrySet.sr.Error()
}

// OffsetsToPosition converts different offsets to positions
func (entrySet *tagKVEntrySet) OffsetsToPosition(
	offsets []int,
) (
	offsetsPos map[int]int,
	err error,
) {
	offsetsPos = make(map[int]int)
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	decoder := encoding.NewDeltaBitPackingDecoder(entrySet.offsetsBlock)
	var (
		maxOffset     = math.MaxInt32
		offsetCounter = 0
	)
	if len(offsets) != 0 {
		maxOffset = offsets[len(offsets)-1]
	}
	for decoder.HasNext() && offsetCounter <= maxOffset {
		pos := decoder.Next()
		if len(offsets) == 0 || intSliceContains(offsets, offsetCounter) {
			offsetsPos[offsetCounter] = int(pos)
		}
		offsetCounter++
	}
	// validate offset
	if len(offsets) != 0 {
		err = fmt.Errorf("read positions of offsets failure")
		if _, ok := offsetsPos[offsets[0]]; !ok {
			return nil, err
		}
		if _, ok := offsetsPos[offsets[len(offsets)-1]]; !ok {
			return nil, err
		}
	}
	return offsetsPos, nil
}

// intSliceContains detects if item is in the slice
func intSliceContains(slice []int, item int) bool {
	idx := sort.Search(len(slice), func(i int) bool { return slice[i] >= item })
	return idx < len(slice) && slice[idx] == item
}

type versionedTagValueData struct {
	version        series.Version
	startTimeDelta int64
	endTimeDelta   int64
	bitMapData     []byte
}

// TimeRange computes the timeRange from delta in seconds
func (data *versionedTagValueData) TimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{
		Start: data.startTimeDelta*1000 + data.version.Int64(),
		End:   data.endTimeDelta*1000 + data.version.Int64()}
}

// Bitmap unmarshals the binary to bitmap
func (data *versionedTagValueData) Bitmap() (*roaring.Bitmap, error) {
	bitmap := roaring.New()
	if err := bitmap.UnmarshalBinary(data.bitMapData); err != nil {
		return nil, err
	}
	return bitmap, nil
}

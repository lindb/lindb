package indextbl

import (
	"fmt"
	"sort"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/series"
)

var seriesIndexReaderLogger = logger.GetLogger("tsdb", "SeriesIndexTableReader")

//go:generate mockgen -source ./series_reader.go -destination=./series_reader_mock.go -package indextbl

const (
	timeRangeSize = 4 + // uint32, start-time
		4 // uint32, end-time
)

// SeriesIndexReader reads versioned seriesID bitmap from series-index-table
type SeriesIndexReader interface {
	// GetSeriesIDsForTagID get series ids for spec metric's keyID
	GetSeriesIDsForTagID(tagID uint32, timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
	// FindSeriesIDsByExprForTagID finds series ids by tag filter expr and tagID
	FindSeriesIDsByExprForTagID(tagID uint32, expr stmt.TagFilter,
		timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
	series.MetaGetter
}

// seriesIndexReader implements SeriesIndexReader
type seriesIndexReader struct {
	snapshot kv.Snapshot
}

// NewSeriesIndexReader returns a new SeriesIndexReader
func NewSeriesIndexReader(snapshot kv.Snapshot) SeriesIndexReader {
	return &seriesIndexReader{snapshot: snapshot}
}

// GetTagValues returns tag values by tag keys and spec version for metric level
func (r *seriesIndexReader) GetTagValues(metricID uint32, tagKeys []string, version uint32) (
	tagValues [][]string, err error) {
	// todo: @codingcrush, #92 forward index implementation
	return nil, nil
}

// FindSeriesIDsByExprForTagID finds series ids by tag filter expr for tagId
func (r *seriesIndexReader) FindSeriesIDsByExprForTagID(tagID uint32, expr stmt.TagFilter,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	entrySetBlocks := r.filterEntrySetBlocks(tagID, timeRange)
	if len(entrySetBlocks) == 0 {
		return nil, series.ErrNotFound
	}
	unionIDSet := series.NewMultiVerSeriesIDSet()
	for _, entrySetBlock := range entrySetBlocks {
		var offsets []int
		q, err := r.entrySetBlockToTreeQuerier(entrySetBlock)
		if err != nil {
			seriesIndexReaderLogger.Error("failed reading trie-tree block", logger.Error(err))
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
		idSet, err := r.entrySetBlockToIDSet(entrySetBlock, timeRange, offsets)
		if err != nil {
			return nil, err
		}
		if idSet == nil {
			continue
		}
		unionIDSet.Or(idSet)
	}
	if unionIDSet.IsEmpty() {
		return nil, series.ErrNotFound
	}
	return unionIDSet, nil
}

// GetSeriesIDsForTagID get series ids for spec metric's tag keyID
func (r *seriesIndexReader) GetSeriesIDsForTagID(tagID uint32,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	entrySetBlocks := r.filterEntrySetBlocks(tagID, timeRange)
	if len(entrySetBlocks) == 0 {
		return nil, series.ErrNotFound
	}
	unionIDSet := series.NewMultiVerSeriesIDSet()
	for _, entrySetBlock := range entrySetBlocks {
		idSet, err := r.entrySetBlockToIDSet(entrySetBlock, timeRange, nil)
		if err != nil {
			return nil, err
		}
		if idSet == nil {
			continue
		}
		unionIDSet.Or(idSet)
	}
	return unionIDSet, nil
}

// filterEntrySetBlocks filters the entry-set block which matches the time-range in the series-index-table
func (r *seriesIndexReader) filterEntrySetBlocks(tagID uint32, timeRange timeutil.TimeRange) (entrySetBlocks [][]byte) {
	for _, reader := range r.snapshot.Readers() {
		block := reader.Get(tagID)
		if len(block) <= timeRangeSize {
			continue
		}
		// read time-range of the total entry-set
		sr := stream.BinaryReader(block)
		startTime := sr.ReadUint32()
		endTime := sr.ReadUint32()
		blockTimeRange := timeutil.TimeRange{
			Start: int64(startTime) * 1000,
			End:   int64(endTime) * 1000}
		if !timeRange.Overlap(&blockTimeRange) {
			continue
		}
		entrySetBlocks = append(entrySetBlocks, block)
	}
	return
}

// entrySetBlockToTreeQuerier converts the binary block to a tire tree block querier
func (r *seriesIndexReader) entrySetBlockToTreeQuerier(block []byte) (trieTreeQuerier, error) {
	var tree trieTreeBlock
	sr := stream.BinaryReader(block)
	// read time-range
	_ = sr.ReadBytes(timeRangeSize)
	////////////////////////////////
	// Block: LOUDS Trie-Tree
	////////////////////////////////
	// read trie-tree length
	expectedTrieTreeLen := sr.ReadUvarint64()
	// read label length
	labelsLen := sr.ReadUvarint64()
	// read labels block
	tree.labels = sr.ReadBytes(int(labelsLen))
	// read isPrefix length
	isPrefixKeyLen := sr.ReadUvarint64()
	// read isPrefixKey bitmap
	isPrefixBlock := sr.ReadBytes(int(isPrefixKeyLen))
	// read LOUDS length
	loudsLen := sr.ReadUvarint64()
	// read LOUDS block
	LOUDSBlock := sr.ReadBytes(int(loudsLen))
	// validation of stream error
	if sr.Error() != nil {
		return nil, sr.Error()
	}
	// validation of length
	realTrieTreeBlockLen := bufioutil.GetUVariantLength(labelsLen) +
		len(tree.labels) +
		bufioutil.GetUVariantLength(isPrefixKeyLen) +
		len(isPrefixBlock) +
		bufioutil.GetUVariantLength(loudsLen) +
		len(LOUDSBlock)
	if realTrieTreeBlockLen != int(expectedTrieTreeLen) {
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
	return &tree, nil
}

// entrySetBlockToIDSet parses the entry-set block, then return the multi-versions seriesID bitmap
func (r *seriesIndexReader) entrySetBlockToIDSet(block []byte, timeRange timeutil.TimeRange,
	offsets []int) (*series.MultiVerSeriesIDSet, error) {

	// read trie-tree length
	sr := stream.BinaryReader(block)
	_ = sr.ReadBytes(timeRangeSize)
	trieTreeLen := sr.ReadUvarint64()
	cursorPos := timeRangeSize + // time range
		bufioutil.GetUVariantLength(trieTreeLen) + int(trieTreeLen) // tree block
	if len(block) <= cursorPos {
		return nil, fmt.Errorf("entrySet block length:%d validation failure", cursorPos)
	}
	// move to the end of trie-tree
	sr = stream.BinaryReader(block[cursorPos:])
	////////////////////////////////
	// Block: TagValue Info
	////////////////////////////////
	// read tag-value count
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	tagValueCount := sr.ReadUvarint64()
	if tagValueCount == 0 {
		return nil, fmt.Errorf("tagValueCount equals to 0")
	}
	// move to tagValueCount
	cursorPos += bufioutil.GetUVariantLength(tagValueCount)
	var (
		// offsets to the end of tagValueInfo block
		tagValueDataBlockOffsets []int
		offsetCounter            = 0
	)
	for i := 0; i < int(tagValueCount); i++ {
		dataLen := sr.ReadUvarint64()
		if sr.Error() != nil {
			return nil, sr.Error()
		}
		// offsets is nil means traversing all blocks of tagValueData
		// offsets contains i mean that is specified offset will be searched
		if len(offsets) == 0 || intSliceContains(offsets, i) {
			tagValueDataBlockOffsets = append(tagValueDataBlockOffsets, offsetCounter)
		}
		offsetCounter += int(dataLen)
		cursorPos += bufioutil.GetUVariantLength(dataLen)
	}
	////////////////////////////////
	// Block: Versioned TagValue Data
	////////////////////////////////
	if len(tagValueDataBlockOffsets) == 0 {
		return nil, series.ErrNotFound
	}
	idSet := series.NewMultiVerSeriesIDSet()
	for _, offset := range tagValueDataBlockOffsets {
		subIDSet, err := r.readTagValueDataBlock(block, offset+cursorPos, timeRange)
		if err != nil {
			return nil, err
		}
		if subIDSet == nil {
			continue
		}
		idSet.Or(subIDSet)
	}
	return idSet, nil
}

// readTagValueDataBlock parses the tagValueDataBlock, and return the the multi-versions seriesID bitmap
func (r *seriesIndexReader) readTagValueDataBlock(block []byte, pos int,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {
	if len(block) <= pos {
		return nil, fmt.Errorf("failed validation of tagValueData's length")
	}
	sr := stream.BinaryReader(block[pos:])
	// read VersionCount
	versionCount := sr.ReadUvarint64()
	if versionCount == 0 {
		return nil, fmt.Errorf("versionCount equals to 0")
	}
	var (
		idSet       *series.MultiVerSeriesIDSet
		readCounter = 0
	)
	for !sr.Empty() && readCounter < int(versionCount) {
		// read version
		version := sr.ReadUint32()
		// read start-time delta
		startTime := sr.ReadVarint64()*1000 + int64(version)*1000 // startTime in milliseconds
		// read end-time delta
		endTime := sr.ReadVarint64()*1000 + int64(version)*1000 // endTime in milliseconds
		// read bitmap length
		bitMapLen := int(sr.ReadUvarint64())
		// read bitmap
		bitMapBlock := sr.ReadBytes(bitMapLen)
		if sr.Error() != nil {
			return nil, sr.Error()
		}
		// finished read a full VersionedTagValue block
		// check time range
		if !timeRange.Overlap(&timeutil.TimeRange{Start: startTime, End: endTime}) {
			readCounter++
			continue
		}
		// unmarshal bitmap
		bitMap := roaring.New()
		if err := bitMap.UnmarshalBinary(bitMapBlock); err != nil {
			return nil, err
		}
		if idSet == nil {
			idSet = series.NewMultiVerSeriesIDSet()
		}
		idSet.Add(version, bitMap)
		readCounter++
	}
	return idSet, nil
}

// intSliceContains detects if item is in the slice
func intSliceContains(slice []int, item int) bool {
	idx := sort.Search(len(slice), func(i int) bool { return slice[i] >= item })
	return idx < len(slice) && slice[idx] == item
}

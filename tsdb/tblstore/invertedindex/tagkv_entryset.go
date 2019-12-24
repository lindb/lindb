package invertedindex

import (
	"fmt"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

type positionIteratorINTF interface {
	// HasNext checks if there are any remaining positions unread
	HasNext() bool

	// Next returns the offset and position read last time
	Next() (offset, position int)
}

type tagKVEntrySetINTF interface {
	// TimeRange computes the timeRange from delta in milliseconds
	TimeRange() timeutil.TimeRange

	// TrieTree builds the trie-tree block for querying
	TrieTree() (trieTreeQuerier, error)

	// ReadTagValueDataBlock iterate on tagValue with specified offset
	ReadTagValueDataBlock(
		offset int,
	) (
		tagValueIterator,
		error,
	)

	// PositionIterator iterates all positions without conditions
	PositionIterator() positionIteratorINTF
}

type tagValueIterator interface {
	// DataTimeRange computes the timeRange from delta in seconds
	DataTimeRange() timeutil.TimeRange

	// DataVersion returns the version of this data
	DataVersion() series.Version

	// HasNext returns if there are any version remaining
	HasNext() bool

	// Next returns the bitmap underlying data.
	Next() (bitmapData []byte)
}

// tagKVEntrySet implements tagKVEntrySetINTF
type tagKVEntrySet struct {
	sr            *stream.Reader
	startTime     int64
	endTime       int64
	tree          trieTreeQuerier
	offsetsBlock  []byte
	crc32CheckSum uint32
	// offsets to positions
	decoder *encoding.FixedOffsetDecoder // decoder for offset
	// tag value data iterator
	versionTotal                 int
	versionRead                  int
	version                      series.Version
	startTimeDelta, endTimeDelta int64
	bitmapData                   []byte
}

func newTagKVEntrySet(block []byte) (tagKVEntrySetINTF, error) {
	if len(block) <= invertedIndexTimeRangeSize+invertedIndexFooterSize {
		return nil, fmt.Errorf("block length no ok")
	}
	entrySet := &tagKVEntrySet{
		sr: stream.NewReader(block),
	}
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
	entrySet.decoder = encoding.NewFixedOffsetDecoder(entrySet.offsetsBlock)
	return entrySet, nil
}

func (entrySet *tagKVEntrySet) TimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{
		Start: entrySet.startTime,
		End:   entrySet.endTime}
}

func (entrySet *tagKVEntrySet) TrieTree() (trieTreeQuerier, error) {
	var tree trieTreeBlock
	// read time-range
	entrySet.sr.ReadAt(invertedIndexTimeRangeSize)
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

func (entrySet *tagKVEntrySet) ReadTagValueDataBlock(
	offset int,
) (
	tagValueIterator,
	error,
) {
	pos := entrySet.decoder.Get(offset)
	if pos < 0 {
		return nil, fmt.Errorf("read position failure")
	}
	// jump to target
	entrySet.sr.ReadAt(pos)
	// reset buffer
	entrySet.versionTotal = int(entrySet.sr.ReadUvarint64())
	entrySet.versionRead = 0
	entrySet.version = 0
	entrySet.startTimeDelta = 0
	entrySet.endTimeDelta = 0
	return entrySet, nil
}

func (entrySet *tagKVEntrySet) DataTimeRange() timeutil.TimeRange {
	return timeutil.TimeRange{
		Start: entrySet.startTimeDelta*1000 + entrySet.version.Int64(),
		End:   entrySet.endTimeDelta*1000 + entrySet.version.Int64()}
}
func (entrySet *tagKVEntrySet) DataVersion() series.Version { return entrySet.version }
func (entrySet *tagKVEntrySet) Next() (bitmapData []byte)   { return entrySet.bitmapData }
func (entrySet *tagKVEntrySet) HasNext() bool {
	if entrySet.sr.Empty() {
		return false
	}
	if entrySet.versionRead >= entrySet.versionTotal {
		return false
	}
	// read version
	entrySet.version = series.Version(entrySet.sr.ReadInt64())
	// read start-time delta
	entrySet.startTimeDelta = entrySet.sr.ReadVarint64()
	// read end-time delta
	entrySet.endTimeDelta = entrySet.sr.ReadVarint64()
	// read bitmap length
	bitMapLen := int(entrySet.sr.ReadUvarint64())
	// read bitmap data
	entrySet.bitmapData = entrySet.sr.ReadSlice(bitMapLen)
	entrySet.versionRead++
	return entrySet.sr.Error() == nil
}

func (entrySet *tagKVEntrySet) PositionIterator() positionIteratorINTF {
	return &positionIterator{decoder: entrySet.decoder}
}

type positionIterator struct {
	decoder  *encoding.FixedOffsetDecoder
	offset   int
	position int
}

func (itr *positionIterator) HasNext() bool {
	itr.position = itr.decoder.Get(itr.offset)
	itr.offset++
	return itr.position >= 0
}
func (itr *positionIterator) Next() (offset, pos int) {
	return itr.offset - 1, itr.position
}

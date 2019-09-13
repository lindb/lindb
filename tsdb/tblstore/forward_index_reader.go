package tblstore

import (
	"fmt"
	"math"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/snappy"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
)

const (
	forwardIndexTimeRangeSize = 4 + // startTime
		4 // endTime
	footerSizeAfterVersionEntries = 4 + // versionOffsetPos, uint32
		4 // CRC32 checksum, uint32
	footerSizeOfVersionEntry = 4 + // Offsets's Position of DictBlock of versionEntry
		4 + // OffsetsBlock's Position of versionEntry
		4 // bitmap's Position of versionEntry
)

var forwardIndexReaderLogger = logger.GetLogger("tsdb", "ForwardIndexReader")

//go:generate mockgen -source ./forward_index_reader.go -destination=./forward_index_reader_mock.go -package tblstore

// ForwardIndexReader reads tagKeys and tagValues from forward-index
type ForwardIndexReader interface {
	series.MetaGetter
}

// forwardIndexReader implements ForwardIndexReader
type forwardIndexReader struct {
	readers []table.Reader
	sr      *stream.Reader
}

type forwardIndexVersionEntry struct {
	startTime           uint32
	endTime             uint32
	tagKeys             []string // tagKeySeq -> tagKey
	tagKeysBitArraySize int
	offsets             []int32
	seriesIDOffsets     *encoding.DeltaBitPackingDecoder
	seriesIDBitmap      *roaring.Bitmap
	// positions
	posOfDictBlock       int
	posOfDictBlockOffset int
	posOfOffsets         int
	posOfSeriesIDBitmap  int
	// tools
	sr           *stream.Reader
	versionBlock []byte
	buffer       []byte
	dict         map[int]string // string index -> string value
	bitArray     *collections.BitArray
}

func newForwardIndexVersionEntry(
	versionBlock []byte,
) (
	versionEntry *forwardIndexVersionEntry,
	err error,
) {
	entry := &forwardIndexVersionEntry{
		versionBlock: versionBlock,
		sr:           stream.NewReader(versionBlock),
		dict:         make(map[int]string),
		bitArray:     collections.NewBitArray(nil),
	}
	// Read Footer
	if err := entry.readFooter(); err != nil {
		return nil, err
	}
	// Read TimeRange Block
	entry.sr.SeekStart()
	entry.startTime = entry.sr.ReadUint32()
	entry.endTime = entry.sr.ReadUint32()
	// Read TagKeys Block
	if err := entry.readTagKeys(); err != nil {
		return nil, err
	}
	// computeBitArraySize
	entry.tagKeysBitArraySize = int(math.Ceil(float64(len(entry.tagKeys)) / float64(8)))
	// Unmarshal offsets
	entry.seriesIDOffsets = encoding.NewDeltaBitPackingDecoder(
		entry.versionBlock[entry.posOfOffsets:entry.posOfSeriesIDBitmap])
	for entry.seriesIDOffsets.HasNext() {
		entry.offsets = append(entry.offsets, entry.seriesIDOffsets.Next())
	}
	// Unmarshal seriesIDBitmap
	entry.seriesIDBitmap = roaring.New()
	if err := entry.seriesIDBitmap.UnmarshalBinary(
		entry.versionBlock[entry.posOfSeriesIDBitmap : len(versionBlock)-footerSizeOfVersionEntry]); err != nil {
		return nil, err
	}
	if len(entry.offsets) != int(entry.seriesIDBitmap.GetCardinality()) {
		return nil, fmt.Errorf("num of offsets does not equal to bitmap's cardinality")
	}
	return entry, nil
}

// seriesID-> {idx1, idx2, idx3}
func (entry *forwardIndexVersionEntry) searchSeriesIDsTagValueIndexes(
	tagKeyIndexes []int,
	seriesIDs *roaring.Bitmap,
) (
	mappings map[uint32][]int,
	err error,
) {
	mappings = make(map[uint32][]int)
	itr := seriesIDs.Iterator()
	for itr.HasNext() {
		seriesID := itr.Next()
		if !entry.seriesIDBitmap.Contains(seriesID) {
			mappings[seriesID] = nil
			continue
		}
		idx := entry.seriesIDBitmap.Rank(seriesID)
		offset := entry.offsets[idx-1]
		indexes, err := entry.searchTagLUT(tagKeyIndexes, offset)
		if err != nil {
			return nil, err
		}
		mappings[seriesID] = indexes
	}
	return mappings, nil
}

func (entry *forwardIndexVersionEntry) searchTagLUT(
	tagKeyIndexes []int,
	offset int32,
) (
	indexes []int,
	err error,
) {
	// jump to the tags LUT block
	entry.sr.SeekStart()
	_ = entry.sr.ReadSlice(int(offset))
	// read bit-array
	bitArrayBuf := entry.sr.ReadSlice(entry.tagKeysBitArraySize)
	entry.bitArray.Reset(bitArrayBuf)
	// pre allocate space
	indexes = make([]int, len(tagKeyIndexes))
	for i := 0; i < len(indexes); i++ {
		indexes[i] = -1
	}
	// helper function
	searchTagKeyIndex := func(expectedTagKeyIndex int) (int, bool) {
		for idx, tagKeyIndex := range tagKeyIndexes {
			if tagKeyIndex == expectedTagKeyIndex {
				return idx, true
			}
		}
		return -1, false
	}

	for tagKeyIndex := range entry.tagKeys {
		// this tagKey exist
		if entry.bitArray.GetBit(uint16(tagKeyIndex)) {
			stringBlockIndex := entry.sr.ReadUvarint64()
			idx, found := searchTagKeyIndex(tagKeyIndex)
			if found {
				indexes[idx] = int(stringBlockIndex)
			}
		}
	}
	return indexes, entry.sr.Error()
}

// readTagKeys reads the tagKeys in order
func (entry *forwardIndexVersionEntry) readTagKeys() error {
	entry.sr.SeekStart()
	_ = entry.sr.ReadSlice(forwardIndexTimeRangeSize)
	tagKeyCount := entry.sr.ReadUvarint64()
	for i := 0; i < int(tagKeyCount); i++ {
		thisTagKeyLength := entry.sr.ReadUvarint64()
		thisTagKey := entry.sr.ReadSlice(int(thisTagKeyLength))
		if entry.sr.Error() != nil {
			return entry.sr.Error()
		}
		entry.tagKeys = append(entry.tagKeys, string(thisTagKey))
	}
	entry.posOfDictBlock = entry.sr.Position()
	return nil
}

// readFooter reads the positions in version entry block
func (entry *forwardIndexVersionEntry) readFooter() (err error) {
	if len(entry.versionBlock) <= footerSizeOfVersionEntry+forwardIndexTimeRangeSize {
		return fmt.Errorf("validation of versionEntrySize failed")
	}
	entry.sr.SeekStart()
	_ = entry.sr.ReadSlice(len(entry.versionBlock) - footerSizeOfVersionEntry)
	entry.posOfDictBlockOffset = int(entry.sr.ReadUint32())
	entry.posOfOffsets = int(entry.sr.ReadUint32())
	entry.posOfSeriesIDBitmap = int(entry.sr.ReadUint32())
	if entry.posOfSeriesIDBitmap >= len(entry.versionBlock) ||
		entry.posOfOffsets >= len(entry.versionBlock) ||
		entry.posOfDictBlockOffset >= len(entry.versionBlock) {
		return fmt.Errorf("position out of index")
	}
	return nil
}

// loadDictByIndexes decodes compressed string to the dict-block by specified indexes
func (entry *forwardIndexVersionEntry) loadDictByIndexes(strIndexes []int) error {
	// Read String Block Offsets In DictBlock
	entry.sr.SeekStart()
	_ = entry.sr.ReadSlice(entry.posOfDictBlockOffset)
	// string block index -> offsets
	var (
		offsets          []int
		lengths          []int
		movedOffset      int
		decodedBlockSeqs = make(map[int]struct{})
	)
	// read string block offsets to StartPosition of DictBlock
	// read stringBlock count
	stringBlockCount := entry.sr.ReadUvarint64()
	for i := 0; i < int(stringBlockCount); i++ {
		offsets = append(offsets, movedOffset)
		length := entry.sr.ReadUvarint64()
		lengths = append(lengths, int(length))
		if entry.sr.Error() != nil {
			return entry.sr.Error()
		}
		movedOffset += int(length)
	}
	// stringBlocksSegment in dict block
	stringBlocksSegment := entry.versionBlock[entry.posOfDictBlock:entry.posOfDictBlockOffset]
	// pick each block
	for _, strIndex := range strIndexes {
		if strIndex < 0 {
			continue
		}
		thisBlockSeq := strIndex / defaultStringBlockSize
		// this block has been decoded before
		if _, ok := decodedBlockSeqs[thisBlockSeq]; ok {
			continue
		}
		// get a uncompressed string block
		if thisBlockSeq >= len(offsets) {
			return fmt.Errorf("index cannot be found in dict block")
		}
		// mark this block is decoded
		decodedBlockSeqs[thisBlockSeq] = struct{}{}
		thisBlockStartPos := offsets[thisBlockSeq]
		thisBlockEndPos := thisBlockStartPos + lengths[thisBlockSeq]
		if thisBlockEndPos > len(stringBlocksSegment) {
			return fmt.Errorf("index string block failure")
		}
		// decode Snappy Compressed String Blocks
		if err := entry.decodeStringBlock(thisBlockSeq,
			stringBlocksSegment[thisBlockStartPos:thisBlockEndPos]); err != nil {
			return err
		}
	}
	return nil
}

// decodeStringBlock decodes the string block, then put it to the map.
func (entry *forwardIndexVersionEntry) decodeStringBlock(
	stringBlockSeq int,
	stringBlock []byte,
) (err error) {
	entry.buffer = entry.buffer[:0]
	if entry.buffer, err = snappy.Decode(entry.buffer, stringBlock); err != nil {
		return err
	}
	// read this decode string block
	entry.sr.Reset(entry.buffer)
	var offset = 0
	for !entry.sr.Empty() {
		tagValueLength := entry.sr.ReadUvarint64()
		tagValue := entry.sr.ReadSlice(int(tagValueLength))
		if entry.sr.Error() != nil {
			return entry.sr.Error()
		}
		entry.dict[stringBlockSeq*defaultStringBlockSize+offset] = string(tagValue)
		offset++
	}
	return nil
}

func (entry *forwardIndexVersionEntry) getTagKeysOrder(
	tagKeys []string,
) (
	tagKeyIndexes []int,
	err error,
) {
	for _, tagKey := range tagKeys {
		matched := false
		for existedTagKeyIndex, existedTagKey := range entry.tagKeys {
			if tagKey == existedTagKey {
				tagKeyIndexes = append(tagKeyIndexes, existedTagKeyIndex)
				matched = true
				break
			}
		}
		if !matched {
			return nil, fmt.Errorf("tagKey: %s not exist", tagKey)
		}
	}
	return tagKeyIndexes, nil
}

// NewForwardIndexReader returns a new ForwardIndexReader
func NewForwardIndexReader(readers []table.Reader) ForwardIndexReader {
	return &forwardIndexReader{
		readers: readers,
		sr:      stream.NewReader(nil)}
}

// GetTagValues returns tag values by tag keys and spec version for metric level
func (r *forwardIndexReader) GetTagValues(
	metricID uint32,
	tagKeys []string,
	version series.Version,
	seriesIDs *roaring.Bitmap,
) (
	seriesID2TagValues map[uint32][]string, // seriesID->
	err error,
) {
	if len(tagKeys) == 0 || seriesIDs.IsEmpty() {
		return nil, series.ErrNotFound
	}
	// get version Block
	versionBlock := r.getVersionBlock(metricID, version)
	if len(versionBlock) == 0 {
		return nil, series.ErrNotFound
	}
	versionEntry, err := newForwardIndexVersionEntry(versionBlock)
	if err != nil {
		return nil, err
	}
	// check if this index does not contains any tagKey in the list
	tagKeyIndexes, err := versionEntry.getTagKeysOrder(tagKeys)
	if err != nil {
		return nil, err
	}

	mappings, err := versionEntry.searchSeriesIDsTagValueIndexes(tagKeyIndexes, seriesIDs)
	if err != nil {
		return nil, err
	}
	// get all string indexes
	var strIndexes []int
	for _, indexes := range mappings {
		strIndexes = append(strIndexes, indexes...)
	}
	if err = versionEntry.loadDictByIndexes(strIndexes); err != nil {
		return nil, err
	}
	seriesID2TagValues = make(map[uint32][]string)
	// assemble the result
	for seriesID, indexes := range mappings {
		tagValues, ok := seriesID2TagValues[seriesID]
		if !ok {
			tagValues = []string{}
		}
		// length=0, means this seriesID inexist
		if len(indexes) == 0 {
			continue
		}
		for _, index := range indexes {
			// index<0, means the tagValue inexist
			if index >= 0 {
				tagValue, ok := versionEntry.dict[index]
				if ok {
					tagValues = append(tagValues, tagValue)
					continue
				}
			}
			tagValues = append(tagValues, "")
		}
		seriesID2TagValues[seriesID] = tagValues
	}
	return seriesID2TagValues, nil
}

// getVersionBlock gets the latest block from snapshot which matches the version in forward-index-table
func (r *forwardIndexReader) getVersionBlock(metricID uint32, version series.Version) (versionBlock []byte) {
	// if we get it from the latest reader, ignore the elder readers
	for i := len(r.readers) - 1; i >= 0; i-- {
		reader := r.readers[i]
		versionBlockItr := newForwardIndexVersionBlockIterator(reader.Get(metricID))
		for versionBlockItr.HasNext() {
			thisVersion, thisVersionBlock := versionBlockItr.Next()
			if thisVersion == version {
				return thisVersionBlock
			}
		}
	}
	return nil
}

type forwardIndexVersionBlockIterator struct {
	block              []byte
	sr                 *stream.Reader
	totalVersions      int // total
	haveReadVersions   int // accumulative
	versionBlockCursor int
}

func newForwardIndexVersionBlockIterator(block []byte) *forwardIndexVersionBlockIterator {
	itr := &forwardIndexVersionBlockIterator{
		block: block,
		sr:    stream.NewReader(block)}
	itr.readTotalVersions()
	return itr
}

func (fii *forwardIndexVersionBlockIterator) readTotalVersions() {
	//////////////////////////////////////////////////
	// Read VersionOffSetsBlock
	//////////////////////////////////////////////////
	_ = fii.sr.ReadSlice(len(fii.block) - footerSizeAfterVersionEntries)
	versionOffsetPos := fii.sr.ReadUint32()
	// shift to Start Position of the VersionOffsetsBlock
	fii.sr.SeekStart()
	_ = fii.sr.ReadSlice(int(versionOffsetPos))
	// read version count
	fii.totalVersions = int(fii.sr.ReadUvarint64())
}

func (fii *forwardIndexVersionBlockIterator) HasNext() bool {
	if len(fii.block) <= footerSizeAfterVersionEntries {
		return false
	}
	if fii.haveReadVersions >= fii.totalVersions {
		return false
	}
	return !fii.sr.Empty() && fii.sr.Error() == nil
}

func (fii *forwardIndexVersionBlockIterator) Next() (version series.Version, versionBlock []byte) {
	defer func() { fii.haveReadVersions++ }()
	// read version
	thisVersion := series.Version(fii.sr.ReadInt64())
	// read version length
	versionLength := fii.sr.ReadUvarint64()
	if fii.sr.Error() != nil {
		forwardIndexReaderLogger.Error("read error occurred", logger.Error(fii.sr.Error()))
		return thisVersion, nil
	}
	versionEntryStartPos := fii.versionBlockCursor
	versionEntryEndPos := versionEntryStartPos + int(versionLength)
	fii.versionBlockCursor += int(versionLength)
	if versionEntryEndPos < len(fii.block) {
		return thisVersion, fii.block[versionEntryStartPos:versionEntryEndPos]
	}
	return thisVersion, nil
}

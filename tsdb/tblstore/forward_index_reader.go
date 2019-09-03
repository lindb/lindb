package tblstore

import (
	"fmt"
	"sort"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"

	"github.com/golang/snappy"
)

const (
	footerSizeAfterVersionEntries = 4 + // versionOffsetPos, uint32
		4 // CRC32 checksum, uint32
	footerSizeOfVersionEntry = 4 + // Offsets's Position of DictBlock of versionEntry
		4 + // Keys Position of versionEntry
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
	buffer  []byte
	dict    map[int]string
}

// NewForwardIndexReader returns a new ForwardIndexReader
func NewForwardIndexReader(readers []table.Reader) ForwardIndexReader {
	return &forwardIndexReader{
		readers: readers,
		dict:    make(map[int]string),
		sr:      stream.NewReader(nil)}
}

// GetTagValues returns tag values by tag keys and spec version for metric level
func (r *forwardIndexReader) GetTagValues(metricID uint32, tagKeys []string, version series.Version) (
	tagValues [][]string, err error) {
	if len(tagKeys) == 0 {
		return nil, nil
	}
	// get version Block
	versionBlock := r.getVersionBlock(metricID, version)
	if len(versionBlock) == 0 {
		return nil, series.ErrNotFound
	}
	// {ip:0, zone:1, host:2}, tagKey -> sequence
	tagKey2Seq, err := r.readTagKeysBlock(versionBlock)
	if err != nil {
		return nil, err
	}
	// indexes of tagKeys
	var existedTagKeysSeq []int
	for _, tagKey := range tagKeys {
		if seq, ok := tagKey2Seq[tagKey]; ok {
			existedTagKeysSeq = append(existedTagKeysSeq, seq)
		}
	}
	// ip, cluster, zone -> 0, 1
	if len(existedTagKeysSeq) == 0 {
		return nil, series.ErrNotFound
	}
	//////////////////////////////////////////////////
	// Read Keys LOOKUP-TABLE Block
	//////////////////////////////////////////////////
	// kvIndexes: {k1seq:[v1,v2,v3], k2Seq: [v2,v4,v6]
	// allValueIndexes: [v1,v2,v3,v4,v6]
	kvIndexes, allValueIndexes, err := r.readKeysLUTBlock(versionBlock, existedTagKeysSeq)
	if err != nil {
		return nil, err
	}
	//////////////////////////////////////////////////
	// Read Dict Block
	//////////////////////////////////////////////////
	if err := r.readDictBlockByIndexes(versionBlock, allValueIndexes); err != nil {
		return nil, err
	}
	// construct tagValues
	for _, tagKey := range tagKeys {
		var thisTagValues []string
		if seq, ok := tagKey2Seq[tagKey]; ok {
			for _, valueIndex := range kvIndexes[seq] {
				if item, ok := r.dict[valueIndex]; ok {
					thisTagValues = append(thisTagValues, item)
				}
			}
		}
		tagValues = append(tagValues, thisTagValues)
	}
	return tagValues, nil
}

// readKeysLUTBlock reads tagValue indexes from the specified tagKeys seq in the LUT block
func (r *forwardIndexReader) readKeysLUTBlock(
	versionBlock []byte,
	tagKeysSeq []int,
) (
	map[int][]int,
	[]int,
	error,
) {
	_, posOfKeysLUT, _, _, err := r.readFooter(versionBlock)
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(tagKeysSeq, func(i, j int) bool { return tagKeysSeq[i] < tagKeysSeq[j] })
	//////////////////////////////////////////////////
	// Keys LOOKUP-TABLE Block
	//////////////////////////////////////////////////
	kvIndexes := make(map[int][]int)
	r.sr.Reset(versionBlock)
	r.sr.ShiftAt(posOfKeysLUT)
	lastSeq := tagKeysSeq[len(tagKeysSeq)-1]
	for seq := 0; seq <= lastSeq && r.sr.Error() == nil; seq++ {
		// jump to the end if we do not need to maps this tagKey
		thisKeyValuesBlockLength := r.sr.ReadUvarint64()
		if !intSliceContains(tagKeysSeq, seq) {
			r.sr.ShiftAt(uint32(thisKeyValuesBlockLength))
			continue
		}
		var thisIndexes []int
		tagValueCount := r.sr.ReadUvarint64()
		for i := 0; i < int(tagValueCount) && r.sr.Error() == nil; i++ {
			thisIndexes = append(thisIndexes, int(r.sr.ReadUvarint64()))
		}
		kvIndexes[seq] = thisIndexes
	}
	// get all indexes of values
	var indexesList []int
	for _, tagValueIndexes := range kvIndexes {
		indexesList = append(indexesList, tagValueIndexes...)
	}
	sort.Slice(indexesList, func(i, j int) bool { return indexesList[i] < indexesList[j] })
	return kvIndexes, indexesList, r.sr.Error()
}

// readTagKeysBlock return a map mapping from tagKey -> tagKey sequence
func (r *forwardIndexReader) readTagKeysBlock(block []byte) (map[string]int, error) {
	r.sr.Reset(block)
	// read time-range
	_ = r.sr.ReadBytes(timeRangeSize)
	//////////////////////////////////////////////////
	// Read TagKeys Block
	//////////////////////////////////////////////////
	tagKey2Seq := make(map[string]int)
	tagKeyCount := r.sr.ReadUvarint64()
	for keySeq := 0; keySeq < int(tagKeyCount); keySeq++ {
		thisTagKeyLength := r.sr.ReadUvarint64()
		thisTagKey := r.sr.ReadBytes(int(thisTagKeyLength))
		if r.sr.Error() != nil {
			return nil, r.sr.Error()
		}
		tagKey2Seq[string(thisTagKey)] = keySeq
	}
	return tagKey2Seq, nil
}

// readFooter reads the positions in version entry block
func (r *forwardIndexReader) readFooter(
	block []byte,
) (
	posOfDictBlockOffset,
	posOfKeysLUT,
	posOfOffsets,
	posOfBitmap uint32,
	err error,
) {
	if len(block) <= footerSizeOfVersionEntry+timeRangeSize {
		err = fmt.Errorf("validation of versionEntrySize failed")
		return
	}
	r.sr.Reset(block)
	r.sr.ShiftAt(uint32(len(block) - footerSizeOfVersionEntry))
	posOfDictBlockOffset = r.sr.ReadUint32()
	posOfKeysLUT = r.sr.ReadUint32()
	posOfOffsets = r.sr.ReadUint32()
	posOfBitmap = r.sr.ReadUint32()
	return
}

// readStringByIndexes reads string from the dict-block by specified indexes
func (r *forwardIndexReader) readDictBlockByIndexes(block []byte, strIndexes []int) error {
	// read PosOfDictBlock Offset in footer
	posOfDictBlockOffset, _, _, _, err := r.readFooter(block)
	if err != nil {
		return err
	}
	//////////////////////////////////////////////////
	// Read String Block Offsets In DictBlock
	//////////////////////////////////////////////////
	r.sr.Reset(block)
	r.sr.ShiftAt(posOfDictBlockOffset)
	// string block index -> offsets
	var (
		offsets     []int
		lengths     []int
		movedOffset int
	)
	// read string block offsets to StartPosition of DictBlock
	// read stringBlock count
	stringBlockCount := r.sr.ReadUvarint64()
	for i := 0; i < int(stringBlockCount); i++ {
		offsets = append(offsets, movedOffset)
		length := r.sr.ReadUvarint64()
		lengths = append(lengths, int(length))
		if r.sr.Error() != nil {
			return r.sr.Error()
		}
		movedOffset += int(length)
	}
	//////////////////////////////////////////////////
	// Read Snappy Compressed String Blocks
	//////////////////////////////////////////////////
	stringBlockStartPos := int(posOfDictBlockOffset) - movedOffset
	stringBlocsEndPos := int(posOfDictBlockOffset)
	if len(block) <= stringBlocsEndPos || stringBlockStartPos < 0 {
		return fmt.Errorf("get string blocks failure")
	}
	return r.readStringBlockByOffsets(block[stringBlockStartPos:stringBlocsEndPos], offsets, lengths, strIndexes)
}

// readStringBlockByOffsets reads different strings in different offsets from the string blocks
// then updates them to the found map.
func (r *forwardIndexReader) readStringBlockByOffsets(
	stringBlocks []byte,
	offsets,
	lengths,
	strIndexes []int,
) error {
	sort.Slice(strIndexes, func(i, j int) bool { return strIndexes[i] < strIndexes[j] })
	// read each block
	lastDecodedBlockSeq := -1
	var err error
	for _, strIndex := range strIndexes {
		thisBlockSeq := strIndex / defaultStringBlockSize
		// this block has been decoded before
		if thisBlockSeq == lastDecodedBlockSeq {
			continue
		}
		// get a uncompressed string block
		if thisBlockSeq >= len(offsets) {
			return fmt.Errorf("index cannot be found in dict block")
		}
		// this block is decoded
		lastDecodedBlockSeq = thisBlockSeq
		thisBlockStartPos := offsets[thisBlockSeq]
		thisBlockEndPos := thisBlockStartPos + lengths[thisBlockSeq]
		if thisBlockEndPos > len(stringBlocks) {
			return fmt.Errorf("index string block failure")
		}
		// decode this string block
		r.buffer = r.buffer[:0]
		if r.buffer, err = snappy.Decode(r.buffer, stringBlocks[thisBlockStartPos:thisBlockEndPos]); err != nil {
			return err
		}
		// read this decode string block
		r.sr.Reset(r.buffer)
		var offset = 0
		for !r.sr.Empty() {
			tagValueLength := r.sr.ReadUvarint64()
			tagValue := r.sr.ReadBytes(int(tagValueLength))
			if r.sr.Error() != nil {
				return r.sr.Error()
			}
			r.dict[thisBlockSeq*defaultStringBlockSize+offset] = string(tagValue)
			offset++
		}
	}
	return nil
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
	fii.sr.ShiftAt(uint32(len(fii.block) - footerSizeAfterVersionEntries))
	versionOffsetPos := fii.sr.ReadUint32()
	// shift to Start Position of the VersionOffsetsBlock
	fii.sr.Reset(fii.block)
	fii.sr.ShiftAt(versionOffsetPos)
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

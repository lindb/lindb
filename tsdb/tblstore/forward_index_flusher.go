package tblstore

import (
	"hash/crc32"
	"math"
	"sync"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bufpool"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/snappy"
)

//go:generate mockgen -source ./forward_index_flusher.go -destination=./forward_index_flusher_mock.go -package tblstore

const (
	// stringBlockSize is the size of a compressed string block
	defaultStringBlockSize = 500
)

var forwardIndexFlusherLogger = logger.GetLogger("tsdb", "ForwardIndexFlusher")

// FlushVersion is a wrapper of kv.Builder, provides the ability to build a forward-index table.
// The layout is available in `tsdb/doc.go`
type ForwardIndexFlusher interface {
	// FlushTagValue flushes a tagValue and the related bitmap
	FlushTagValue(tagValue string, bitmap *roaring.Bitmap)
	// FlushTagKey ends writing the tagValues
	FlushTagKey(tagKey string)
	// FlushVersion ends writes a version block
	FlushVersion(version uint32, startTime, endTime uint32)
	// FlushMetricID ends write a full metric-block
	FlushMetricID(metricID uint32) error
	// Commit closes the writer, this will be called after writing all tagKeys.
	Commit() error
}

// forwardIndexFlusher implements ForwardIndexFlusher
type forwardIndexFlusher struct {
	// build version block
	offsets           *encoding.DeltaBitPackingEncoder // offsets
	keys              *roaring.Bitmap                  // keys
	tagKeysList       []string                         // tagKeys in order
	tagKeysMap        map[string]int                   // tagKey -> index in tagKeysList
	tagValuesList     []string                         // tagValues in order
	tagValuesMap      map[string]int                   // tagValue -> index in tagValuesList
	seriesID2TagValue map[uint32]*[]int                // seriesID -> tagValue index in order
	seriesID2TagKey   map[uint32]*[]int                // seriesID -> tagKey index in order
	tagKeyValuesMap   map[string]*[]int                // tagKey -> tagValue indexes
	tagValuesOfKey    map[int]struct{}                 // tagValue indexes of this tagKey in order
	// build metric block
	metricBlockWriter *stream.BufferWriter // writer for build metric-block
	versionBlocksLen  []int                // length of all flushed version blocks
	versions          []uint32             // all flushed versions
	// common elements
	tmpWriter *stream.BufferWriter // temporary writer
	dstSlice  []byte               // snappy dst slice
	intPool   sync.Pool            // storing *[]int
	kvFlusher kv.Flusher           // real underlying flusher
	bitArray  *collections.BitArray
	// used for mock, default false
	resetDisabled bool
}

// NewForwardIndexFlusher returns a new ForwardIndexFlusher.
func NewForwardIndexFlusher(flusher kv.Flusher) ForwardIndexFlusher {
	bitArray, _ := collections.NewBitArray(nil)
	return &forwardIndexFlusher{
		kvFlusher:         flusher,
		tagKeysMap:        make(map[string]int),
		tagValuesMap:      make(map[string]int),
		seriesID2TagValue: make(map[uint32]*[]int),
		seriesID2TagKey:   make(map[uint32]*[]int),
		tagKeyValuesMap:   make(map[string]*[]int),
		tagValuesOfKey:    make(map[int]struct{}),
		intPool: sync.Pool{New: func() interface{} {
			return &[]int{}
		}},
		tmpWriter:         stream.NewBufferWriter(nil),
		metricBlockWriter: stream.NewBufferWriter(nil),
		keys:              roaring.New(),
		offsets:           encoding.NewDeltaBitPackingEncoder(),
		bitArray:          bitArray}
}

func (flusher *forwardIndexFlusher) getSlice() *[]int {
	return flusher.intPool.Get().(*[]int)
}
func (flusher *forwardIndexFlusher) putSlice(s *[]int) {
	*s = (*s)[:0]
	flusher.intPool.Put(s)
}

// FlushTagKey ends writing the tagValues
func (flusher *forwardIndexFlusher) FlushTagValue(tagValue string, bitmap *roaring.Bitmap) {
	// do not insert a same tagValue twice
	idxOfTagValuesList, ok := flusher.tagValuesMap[tagValue]
	if !ok {
		idxOfTagValuesList = len(flusher.tagValuesList)
		flusher.tagValuesMap[tagValue] = idxOfTagValuesList
		flusher.tagValuesList = append(flusher.tagValuesList, tagValue)
	}
	// record the newly tagValues of the tagKey
	flusher.tagValuesOfKey[idxOfTagValuesList] = struct{}{}

	iterator := bitmap.Iterator()
	for iterator.HasNext() {
		seriesID := iterator.Next()
		// record newly written index of tagValues
		orderedTagValues, ok := flusher.seriesID2TagValue[seriesID]
		if !ok {
			orderedTagValues = flusher.getSlice()
		}
		*orderedTagValues = append(*orderedTagValues, idxOfTagValuesList)
		flusher.seriesID2TagValue[seriesID] = orderedTagValues
		// record newly written index of tagKeys
		orderedTagKeys, ok := flusher.seriesID2TagKey[seriesID]
		if !ok {
			orderedTagKeys = flusher.getSlice()
		}
		*orderedTagKeys = append(*orderedTagKeys, len(flusher.tagKeysList))
		flusher.seriesID2TagKey[seriesID] = orderedTagKeys
	}
}

// FlushTagKey ends writing the tagValues
func (flusher *forwardIndexFlusher) FlushTagKey(tagKey string) {
	flusher.tagKeysList = append(flusher.tagKeysList, tagKey)
	flusher.tagKeysMap[tagKey] = len(flusher.tagKeysList) - 1

	// record all tagValues of this Key
	newSlice := flusher.getSlice()
	for tagValueIndex := range flusher.tagValuesOfKey {
		*newSlice = append(*newSlice, tagValueIndex)
		// clear the container
		delete(flusher.tagValuesOfKey, tagValueIndex)
	}
	flusher.tagKeyValuesMap[tagKey] = newSlice
}

// reset resets the internal containers for build next version block
func (flusher *forwardIndexFlusher) resetVersionContext() {
	// reset tagKeys related
	for _, tagKey := range flusher.tagKeysList {
		delete(flusher.tagKeysMap, tagKey)
	}
	flusher.tagKeysList = flusher.tagKeysList[:0]
	// reset tag values
	for _, tagValue := range flusher.tagValuesList {
		delete(flusher.tagValuesMap, tagValue)
	}
	flusher.tagValuesList = flusher.tagValuesList[:0]
	// reset seriesID2TagValue, seriesID2TagKey
	for seriesID, sl := range flusher.seriesID2TagValue {
		flusher.putSlice(sl)
		delete(flusher.seriesID2TagValue, seriesID)
	}
	for seriesID, sl := range flusher.seriesID2TagKey {
		flusher.putSlice(sl)
		delete(flusher.seriesID2TagKey, seriesID)
	}
	// reset tagKeyValuesMap
	for tagKey, tagValueIndexes := range flusher.tagKeyValuesMap {
		flusher.putSlice(tagValueIndexes)
		delete(flusher.tagValuesMap, tagKey)
	}
	// reset keys and offsets
	flusher.offsets.Reset()
	flusher.keys.Clear()
}

// FlushVersion ends writes a version block
func (flusher *forwardIndexFlusher) FlushVersion(version uint32, startTime, endTime uint32) {
	//////////////////////////////////////////////////
	// Reset
	//////////////////////////////////////////////////
	defer flusher.resetVersionContext()
	// record the start position of this entry
	startPosOfThisEntry := flusher.metricBlockWriter.Len()
	// record the length of this entry
	defer func() {
		endPosOfThisEntry := flusher.metricBlockWriter.Len()
		flusher.versionBlocksLen = append(flusher.versionBlocksLen, endPosOfThisEntry-startPosOfThisEntry)
		flusher.versions = append(flusher.versions, version)
	}()
	//////////////////////////////////////////////////
	// build Time Range Block
	//////////////////////////////////////////////////
	// write start-time
	flusher.metricBlockWriter.PutUint32(startTime)
	// write end-time
	flusher.metricBlockWriter.PutUint32(endTime)
	//////////////////////////////////////////////////
	// build TagKeys Block
	//////////////////////////////////////////////////
	// write tag-key count
	flusher.metricBlockWriter.PutUvarint64(uint64(len(flusher.tagKeysList)))
	// write tagKey length and tagKey
	for _, tagKey := range flusher.tagKeysList {
		flusher.metricBlockWriter.PutUvarint64(uint64(len(tagKey)))
		flusher.metricBlockWriter.PutBytes([]byte(tagKey))
	}
	//////////////////////////////////////////////////
	// build Dict Block
	//////////////////////////////////////////////////
	dictBlockOffsetPos := flusher.buildDictBlocks()
	//////////////////////////////////////////////////
	// TagKeys LUT Block
	//////////////////////////////////////////////////
	tagKeysLUTBlockPos := flusher.metricBlockWriter.Len()
	flusher.buildKeysLUTBlock()
	//////////////////////////////////////////////////
	// build Series Tags Block's BitArray for TagKeys
	//////////////////////////////////////////////////
	for seriesID, tagKeyIndexes := range flusher.seriesID2TagKey {
		tagsBlockPosition := flusher.metricBlockWriter.Len()
		flusher.bitArray.Reset()
		for _, idx := range *tagKeyIndexes {
			flusher.bitArray.SetBit(uint16(idx))
		}
		// write bit-array
		flusher.metricBlockWriter.PutBytes(flusher.bitArray.Bytes())
		// get tagValue indexes
		tagValueIndexes := flusher.seriesID2TagValue[seriesID]
		for _, idx := range *tagValueIndexes {
			// write index of tagValue in String Block
			flusher.metricBlockWriter.PutUvarint64(uint64(idx))
		}
		// write offset of tags block in the version block
		flusher.offsets.Add(int32(tagsBlockPosition))
		// write seriesID
		flusher.keys.Add(seriesID)
	}
	//////////////////////////////////////////////////
	// build offsets, keys, footer
	//////////////////////////////////////////////////
	flusher.finishVersion(startPosOfThisEntry, dictBlockOffsetPos, tagKeysLUTBlockPos)
}

// finishVersion writes the version
func (flusher *forwardIndexFlusher) finishVersion(startPos, dictBlockOffsetPos, tagKeysLUTBlockPos int) {
	offsets := flusher.offsets.Bytes()
	// position of the offset block
	offsetsPosition := flusher.metricBlockWriter.Len()
	// write offsets
	flusher.metricBlockWriter.PutBytes(offsets)

	// write keys
	flusher.keys.RunOptimize()
	keys, err := flusher.keys.MarshalBinary()
	if err != nil {
		forwardIndexFlusherLogger.Error("marshal keys error", logger.Error(err))
	}
	// position of the keys block
	keysPosition := flusher.metricBlockWriter.Len()
	flusher.metricBlockWriter.PutBytes(keys)
	//////////////////////////////////////////////////
	// build Footer
	//////////////////////////////////////////////////
	// write pos of dict block offset
	flusher.metricBlockWriter.PutUint32(uint32(dictBlockOffsetPos - startPos))
	// write pos of keys LUT
	flusher.metricBlockWriter.PutUint32(uint32(tagKeysLUTBlockPos - startPos))
	// write pos of offset blocks
	flusher.metricBlockWriter.PutUint32(uint32(offsetsPosition - startPos))
	// write pos of keys block
	flusher.metricBlockWriter.PutUint32(uint32(keysPosition - startPos))
}

// buildKeysLUTBlock writes the keys LUT block to the writer
func (flusher *forwardIndexFlusher) buildKeysLUTBlock() {
	for _, tagKey := range flusher.tagKeysList {
		flusher.tmpWriter.Reset()
		valuesIndexes := flusher.tagKeyValuesMap[tagKey]
		// write key values count
		flusher.tmpWriter.PutUvarint64(uint64(len(*valuesIndexes)))
		// write each tag value index
		for _, tagValueIndex := range *valuesIndexes {
			flusher.tmpWriter.PutUvarint64(uint64(tagValueIndex))
		}
		// write this keyBlock to the real writer
		data, _ := flusher.tmpWriter.Bytes()
		// write this keyBlock length
		flusher.metricBlockWriter.PutUvarint64(uint64(len(data)))
		// write this keyBlock
		flusher.metricBlockWriter.PutBytes(data)
	}
}

// buildDictBlocks writes the dict block to the writer
func (flusher *forwardIndexFlusher) buildDictBlocks() (offsetPos int) {
	tagValuesCount := len(flusher.tagValuesList)
	blockCount := int(math.Ceil(float64(tagValuesCount) / float64(defaultStringBlockSize)))
	//////////////////////////////////////////////////
	// build Snappy Compressed String block
	//////////////////////////////////////////////////
	// get a slice for writing all block length
	blockLengths := flusher.getSlice()
	defer flusher.putSlice(blockLengths)
	// get a stream writer to build string block
	buf := bufpool.GetBuffer()
	sw := stream.NewBufferWriter(buf)
	defer bufpool.PutBuffer(buf)

	for i := 0; i < blockCount; i++ {
		start := i * defaultStringBlockSize
		end := (i + 1) * defaultStringBlockSize
		if end > tagValuesCount {
			end = tagValuesCount
		}
		// clean the slice before use
		sw.Reset()
		flusher.dstSlice = flusher.dstSlice[:0]
		// build src slice
		for j := start; j < end; j++ {
			data := []byte(flusher.tagValuesList[j])
			sw.PutUvarint64(uint64(len(data)))
			sw.PutBytes(data)
		}
		thisBlock, _ := sw.Bytes()
		// encode to dst slice
		flusher.dstSlice = snappy.Encode(flusher.dstSlice, thisBlock)
		// record the length
		*blockLengths = append(*blockLengths, len(flusher.dstSlice))
		// write this block
		flusher.metricBlockWriter.PutBytes(flusher.dstSlice)
	}
	//////////////////////////////////////////////////
	// build String Block Offsets
	//////////////////////////////////////////////////
	// report the start position of the offsets block
	offsetPos = flusher.metricBlockWriter.Len()
	// write string block count
	flusher.metricBlockWriter.PutUvarint64(uint64(blockCount))
	// write each block length
	for _, l := range *blockLengths {
		flusher.metricBlockWriter.PutUvarint64(uint64(l))
	}
	return offsetPos
}

// resetMetricBlockContext resets the internal containers for build next metric block
func (flusher *forwardIndexFlusher) resetMetricBlockContext() {
	// reset writer
	flusher.metricBlockWriter.Reset()
	// reset version block meta info
	flusher.versionBlocksLen = flusher.versionBlocksLen[:0]
	flusher.versions = flusher.versions[:0]
}

// FlushMetricID ends write a full metric-block
func (flusher *forwardIndexFlusher) FlushMetricID(metricID uint32) error {
	//////////////////////////////////////////////////
	// Reset
	//////////////////////////////////////////////////
	if !flusher.resetDisabled {
		defer flusher.resetMetricBlockContext()
	}
	//////////////////////////////////////////////////
	// build Version Offsets Block
	//////////////////////////////////////////////////
	// start position of the offsets block
	posOfVersionOffsets := flusher.metricBlockWriter.Len()
	// write versions count
	flusher.metricBlockWriter.PutUvarint64(uint64(len(flusher.versionBlocksLen)))
	// write all versions and version lengths
	for idx, version := range flusher.versions {
		// write version
		flusher.metricBlockWriter.PutUint32(version)
		// write version block length
		flusher.metricBlockWriter.PutUvarint64(uint64(flusher.versionBlocksLen[idx]))
	}
	//////////////////////////////////////////////////
	// build Footer
	//////////////////////////////////////////////////
	// write position of the offsets block
	flusher.metricBlockWriter.PutUint32(uint32(posOfVersionOffsets))
	// write CRC32 checksum
	data, _ := flusher.metricBlockWriter.Bytes()
	flusher.metricBlockWriter.PutUint32(crc32.ChecksumIEEE(data))
	// real flush process
	data, _ = flusher.metricBlockWriter.Bytes()
	return flusher.kvFlusher.Add(metricID, data)
}

// Commit closes the writer, this will be called after writing all metrics.
func (flusher *forwardIndexFlusher) Commit() error {
	return flusher.kvFlusher.Commit()
}

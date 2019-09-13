package tblstore

import (
	"hash/crc32"
	"math"
	"sort"
	"sync"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/snappy"
)

//go:generate mockgen -source ./forward_index_flusher.go -destination=./forward_index_flusher_mock.go -package tblstore

const (
	// stringBlockSize is the size of a compressed string block
	defaultStringBlockSize = 300
)

var (
	forwardIndexFlusherLogger = logger.GetLogger("tsdb", "ForwardIndexFlusher")
	intPool                   = sync.Pool{New: func() interface{} {
		return &[]int{} // storing *[]int

	}}
)

// FlushVersion is a wrapper of kv.Builder, provides the ability to build a forward-index table.
// The layout is available in `tsdb/doc.go`
type ForwardIndexFlusher interface {
	// FlushTagValue flushes a tagValue and the related bitmap
	FlushTagValue(tagValue string, bitmap *roaring.Bitmap)
	// FlushTagKey ends writing the tagValues
	FlushTagKey(tagKey string)
	// FlushVersion ends writes a version block
	FlushVersion(version series.Version, startTime, endTime uint32)
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
	sortedSeriesIDs   []uint32                         // used for sort
	// build metric block
	metricBlockWriter *stream.BufferWriter // writer for build metric-block
	versionBlocks     []struct {
		length  int            // length of flushed version blocks
		version series.Version // flushed version
	}
	// common elements
	tmpWriter *stream.BufferWriter // temporary writer
	dstSlice  []byte               // snappy dst slice
	kvFlusher kv.Flusher           // real underlying flusher
	bitArray  *collections.BitArray
}

// NewForwardIndexFlusher returns a new ForwardIndexFlusher.
func NewForwardIndexFlusher(flusher kv.Flusher) ForwardIndexFlusher {
	return &forwardIndexFlusher{
		kvFlusher:         flusher,
		tagKeysMap:        make(map[string]int),
		tagValuesMap:      make(map[string]int),
		seriesID2TagValue: make(map[uint32]*[]int),
		seriesID2TagKey:   make(map[uint32]*[]int),
		metricBlockWriter: stream.NewBufferWriter(nil),
		tmpWriter:         stream.NewBufferWriter(nil),
		keys:              roaring.New(),
		offsets:           encoding.NewDeltaBitPackingEncoder(),
		bitArray:          collections.NewBitArray(nil)}
}

func (flusher *forwardIndexFlusher) getSlice() *[]int {
	return intPool.Get().(*[]int)
}

func (flusher *forwardIndexFlusher) putSlice(s *[]int) {
	*s = (*s)[:0]
	intPool.Put(s)
}

func (flusher *forwardIndexFlusher) sortSeriesIDs() {
	flusher.sortedSeriesIDs = flusher.sortedSeriesIDs[:0]
	for seriesID := range flusher.seriesID2TagKey {
		flusher.sortedSeriesIDs = append(flusher.sortedSeriesIDs, seriesID)
	}
	sort.Slice(flusher.sortedSeriesIDs, func(i, j int) bool {
		return flusher.sortedSeriesIDs[i] < flusher.sortedSeriesIDs[j]
	})
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
	// reset keys and offsets
	flusher.offsets.Reset()
	flusher.keys.Clear()
}

// FlushVersion ends writes a version block
func (flusher *forwardIndexFlusher) FlushVersion(version series.Version, startTime, endTime uint32) {
	//////////////////////////////////////////////////
	// Reset
	//////////////////////////////////////////////////
	defer flusher.resetVersionContext()
	// record the start position of this entry
	startPosOfThisEntry := flusher.metricBlockWriter.Len()
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
	// write Dict Block
	//////////////////////////////////////////////////
	dictBlockOffsetPos := flusher.writeDictBlocks()
	//////////////////////////////////////////////////
	// build Series Tags Block's BitArray for TagKeys
	//////////////////////////////////////////////////
	flusher.sortSeriesIDs()
	for _, seriesID := range flusher.sortedSeriesIDs {
		tagKeyIndexes := flusher.seriesID2TagKey[seriesID]
		tagsBlockPosition := flusher.metricBlockWriter.Len()
		flusher.bitArray.Reset(nil)
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
		flusher.offsets.Add(int32(tagsBlockPosition - startPosOfThisEntry))
		// write seriesID
		flusher.keys.Add(seriesID)
	}
	//////////////////////////////////////////////////
	// build offsets, keys, footer
	//////////////////////////////////////////////////
	flusher.finishVersion(startPosOfThisEntry, dictBlockOffsetPos)
	// record the length of the entry
	flusher.RecordVersionOffset(version, startPosOfThisEntry)
}

func (flusher *forwardIndexFlusher) RecordVersionOffset(version series.Version, startPos int) {
	endPosOfThisEntry := flusher.metricBlockWriter.Len()
	flusher.versionBlocks = append(flusher.versionBlocks, struct {
		length  int
		version series.Version
	}{length: endPosOfThisEntry - startPos, version: version})
}

// finishVersion writes the version
func (flusher *forwardIndexFlusher) finishVersion(startPos, dictBlockOffsetPos int) {
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
	// write pos of offset blocks
	flusher.metricBlockWriter.PutUint32(uint32(offsetsPosition - startPos))
	// write pos of keys block
	flusher.metricBlockWriter.PutUint32(uint32(keysPosition - startPos))
}

// writeDictBlocks writes the dict block to the writer
func (flusher *forwardIndexFlusher) writeDictBlocks() (offsetPos int) {
	tagValuesCount := len(flusher.tagValuesList)
	blockCount := int(math.Ceil(float64(tagValuesCount) / float64(defaultStringBlockSize)))
	//////////////////////////////////////////////////
	// build Snappy Compressed String block
	//////////////////////////////////////////////////
	// get a slice for writing all block length
	blockLengths := flusher.getSlice()
	defer flusher.putSlice(blockLengths)

	for i := 0; i < blockCount; i++ {
		start := i * defaultStringBlockSize
		end := (i + 1) * defaultStringBlockSize
		if end > tagValuesCount {
			end = tagValuesCount
		}
		// clean the slice before use
		flusher.tmpWriter.Reset()
		flusher.dstSlice = flusher.dstSlice[:0]
		// build src slice
		for j := start; j < end; j++ {
			data := []byte(flusher.tagValuesList[j])
			flusher.tmpWriter.PutUvarint64(uint64(len(data)))
			flusher.tmpWriter.PutBytes(data)
		}
		thisBlock, _ := flusher.tmpWriter.Bytes()
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

// Reset resets the all containers for build next metric block
func (flusher *forwardIndexFlusher) Reset() {
	flusher.resetVersionContext()
	// reset writer
	flusher.metricBlockWriter.Reset()
	flusher.tmpWriter.Reset()
	// reset version block meta info
	flusher.versionBlocks = flusher.versionBlocks[:0]
}

// FlushMetricID ends write a full metric-block
func (flusher *forwardIndexFlusher) FlushMetricID(metricID uint32) error {
	defer flusher.Reset()
	//////////////////////////////////////////////////
	// build Version Offsets Block
	//////////////////////////////////////////////////
	// start position of the offsets block
	posOfVersionOffsets := flusher.metricBlockWriter.Len()
	// write versions count
	flusher.metricBlockWriter.PutUvarint64(uint64(len(flusher.versionBlocks)))
	// write all versions and version lengths
	for _, versionBlock := range flusher.versionBlocks {
		// write version
		flusher.metricBlockWriter.PutInt64(versionBlock.version.Int64())
		// write version block length
		flusher.metricBlockWriter.PutUvarint64(uint64(versionBlock.length))
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

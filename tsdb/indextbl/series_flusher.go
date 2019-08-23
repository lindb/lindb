package indextbl

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sync"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/logger"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./series_flusher.go -destination=./series_flusher_mock.go -package indextbl

var seriesIndexFlusherLogger = logger.GetLogger("tsdb", "SeriesIndexTableFlusher")

// SeriesIndexFlusher is a wrapper of kv.Builder, provides the ability to build a versioned series-id index table.
// The layout is available in `tsdb/doc.go`
type SeriesIndexFlusher interface {
	// FlushVersion writes a versioned bitmap to index table.
	FlushVersion(version uint32, startTime, endTime uint32, bitmap *roaring.Bitmap)
	// FlushTagValue ends writing VersionedTagValueBlock in index table.
	FlushTagValue(tagValue string)
	// FlushTagKey ends writing entrySetBlock in index table.
	FlushTagKey(tagID uint32) error
	// Commit closes the writer, this will be called after writing all tagKeys.
	Commit() error
}

// NewSeriesIndexFlusher returns a new SeriesIndexFlusher
func NewSeriesIndexFlusher(flusher kv.Flusher) SeriesIndexFlusher {
	return &seriesIndexFlusher{
		flusher:    flusher,
		bufferPool: sync.Pool{New: func() interface{} { return &bytes.Buffer{} }},
		trie:       newTrieTree(),
	}
}

// seriesIndexFlusher implements SeriesIndexFlusher.
type seriesIndexFlusher struct {
	flusher        kv.Flusher
	trie           trieTreeBuilder
	entrySetBuffer bytes.Buffer
	bufferPool     sync.Pool
	VariableBuf    [8]byte
	// time range
	minStartTime uint32
	maxEndTime   uint32
	// for tagValue data builder
	versionCount   int
	tagValueBuffer *bytes.Buffer
	// used for mock
	resetDisabled bool
}

func (w *seriesIndexFlusher) getBuffer() *bytes.Buffer {
	return w.bufferPool.Get().(*bytes.Buffer)
}

// FlushVersion writes a versioned bitmap to index table.
func (w *seriesIndexFlusher) FlushVersion(version uint32, startTime, endTime uint32, bitmap *roaring.Bitmap) {
	if w.tagValueBuffer == nil {
		w.tagValueBuffer = w.getBuffer()
	}
	// count flushed versions
	w.versionCount++
	// update time range
	if startTime < w.minStartTime || w.minStartTime == 0 {
		w.minStartTime = startTime
	}
	if endTime > w.maxEndTime || w.maxEndTime == 0 {
		w.maxEndTime = endTime
	}
	// write version
	binary.LittleEndian.PutUint32(w.VariableBuf[:], version)
	_, _ = w.tagValueBuffer.Write(w.VariableBuf[:4])
	// write startTime delta
	startTimeDelta := int64(startTime) - int64(version)
	size := binary.PutVarint(w.VariableBuf[:], startTimeDelta)
	_, _ = w.tagValueBuffer.Write(w.VariableBuf[:size])
	// write endTime delta
	endTimeDelta := int64(endTime) - int64(version)
	size = binary.PutVarint(w.VariableBuf[:], endTimeDelta)
	_, _ = w.tagValueBuffer.Write(w.VariableBuf[:size])
	// write bitmap length
	out, err := bitmap.MarshalBinary()
	if err != nil {
		seriesIndexFlusherLogger.Error("marshal bitmap failure", logger.Error(err))
	}
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(out)))
	_, _ = w.tagValueBuffer.Write(w.VariableBuf[:size])
	// write bitmap data
	_, _ = w.tagValueBuffer.Write(out)
}

// bufferWithVersionCount is the value of trie-tree node
type bufferWithVersionCount struct {
	versionCount int
	buffer       *bytes.Buffer
}

// FlushTagValue indicate a VersionedTagValueDataBlock is done.
func (w *seriesIndexFlusher) FlushTagValue(tagValue string) {
	w.trie.Add(tagValue, bufferWithVersionCount{
		versionCount: w.versionCount,
		buffer:       w.tagValueBuffer})

	w.tagValueBuffer = nil
	w.versionCount = 0
}

func (w *seriesIndexFlusher) FlushTagKey(tagID uint32) error {
	if !w.resetDisabled {
		defer w.reset()
	}

	treeDataBlock := w.trie.MarshalBinary()
	// write startTime
	binary.LittleEndian.PutUint32(w.VariableBuf[:], w.minStartTime)
	w.entrySetBuffer.Write(w.VariableBuf[:4])
	// write endTime
	binary.LittleEndian.PutUint32(w.VariableBuf[:], w.maxEndTime)
	w.entrySetBuffer.Write(w.VariableBuf[:4])
	// build isPrefixKey
	isPrefixBlock, err := treeDataBlock.isPrefixKey.MarshalBinary()
	if err != nil {
		return err
	}
	// build LOUDS length
	LOUDSBlock, err := treeDataBlock.LOUDS.MarshalBinary()
	if err != nil {
		return err
	}
	treeLength := bufioutil.GetUVariantLength(uint64(len(treeDataBlock.labels))) + // labels length uvariant size
		len(treeDataBlock.labels) + // labels length
		bufioutil.GetUVariantLength(uint64(len(isPrefixBlock))) + // isPrefixKey length uvariant size
		len(isPrefixBlock) + // isPrefixKey length
		bufioutil.GetUVariantLength(uint64(len(LOUDSBlock))) + // LOUDSBlock length uvariantsize
		len(LOUDSBlock) // LOUDSBlock length
	// write tree length
	size := binary.PutUvarint(w.VariableBuf[:], uint64(treeLength))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	// write labels length & labels
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(treeDataBlock.labels)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	_, _ = w.entrySetBuffer.Write(treeDataBlock.labels)
	// write isPrefixKey length & bitmap
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(isPrefixBlock)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	_, _ = w.entrySetBuffer.Write(isPrefixBlock)
	// write LOUDS length & bitmap
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(LOUDSBlock)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	_, _ = w.entrySetBuffer.Write(LOUDSBlock)

	// write tagValueCount
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(treeDataBlock.values)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])

	// write all data length and versioned tagValue data blocks
	w.writeTagValueDataBlockTo(&w.entrySetBuffer, treeDataBlock)

	// write crc32 checksum
	binary.LittleEndian.PutUint32(w.VariableBuf[0:4], crc32.ChecksumIEEE(w.entrySetBuffer.Bytes()))
	w.entrySetBuffer.Write(w.VariableBuf[:4])

	return w.flusher.Add(tagID, w.entrySetBuffer.Bytes())
}

// writeTagValueDataBlockTo write tagValueDataBlocks to the buffer.
func (w *seriesIndexFlusher) writeTagValueDataBlockTo(buffer *bytes.Buffer, treeDataBlock *trieTreeData) {
	// write lengths of all versioned tagValue data block
	for _, item := range treeDataBlock.values {
		it := item.(bufferWithVersionCount)
		// write all data length
		dataBlockLen := bufioutil.GetUVariantLength(uint64(it.versionCount)) + // version count size
			len(it.buffer.Bytes()) // versionedTagValue blocks
		size := binary.PutUvarint(w.VariableBuf[:], uint64(dataBlockLen))
		_, _ = buffer.Write(w.VariableBuf[:size])
	}
	// write all versioned tagValue data block
	for _, item := range treeDataBlock.values {
		it := item.(bufferWithVersionCount)
		// write version count
		size := binary.PutUvarint(w.VariableBuf[:], uint64(it.versionCount))
		_, _ = buffer.Write(w.VariableBuf[:size])
		// write all versions of tagValue bitmaps
		_, _ = buffer.Write(it.buffer.Bytes())
		// put buffer back to pool
		it.buffer.Reset()
		w.bufferPool.Put(it.buffer)
	}
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (w *seriesIndexFlusher) Commit() error {
	w.reset()
	return w.flusher.Commit()
}

// reset resets the trie and buf
func (w *seriesIndexFlusher) reset() {
	w.trie.Reset()
	w.entrySetBuffer.Reset()
}

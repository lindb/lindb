package indextbl

import (
	"bytes"
	"hash/crc32"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bufpool"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./inverted_index_flusher.go -destination=./inverted_index_flusher_mock.go -package indextbl

var invertedIndexFlusherLogger = logger.GetLogger("tsdb", "InvertedIndexFlusher")

// InvertedIndexFlusher is a wrapper of kv.Builder, provides the ability to build a versioned series-id index table.
// The layout is available in `tsdb/doc.go`
type InvertedIndexFlusher interface {
	// FlushVersion writes a versioned bitmap to index table.
	FlushVersion(version uint32, startTime, endTime uint32, bitmap *roaring.Bitmap)
	// FlushTagValue ends writing VersionedTagValueBlock in index table.
	FlushTagValue(tagValue string)
	// FlushTagKey ends writing entrySetBlock in index table.
	FlushTagKey(tagID uint32) error
	// Commit closes the writer, this will be called after writing all tagKeys.
	Commit() error
}

// NewInvertedIndexFlusher returns a new InvertedIndexFlusher
func NewInvertedIndexFlusher(flusher kv.Flusher) InvertedIndexFlusher {
	return &invertedIndexFlusher{
		flusher:        flusher,
		entrySetWriter: stream.NewBufferWriter(nil),
		trie:           newTrieTree(),
		tagValueWriter: stream.NewBufferWriter(nil),
	}
}

// invertedIndexFlusher implements InvertedIndexFlusher.
type invertedIndexFlusher struct {
	flusher        kv.Flusher
	trie           trieTreeBuilder
	entrySetWriter *stream.BufferWriter
	// time range
	minStartTime uint32
	maxEndTime   uint32
	// for tagValue data builder
	versionCount   int
	tagValueWriter *stream.BufferWriter
	tagValueBuffer *bytes.Buffer
	// used for mock
	resetDisabled bool
}

// FlushVersion writes a versioned bitmap to index table.
func (w *invertedIndexFlusher) FlushVersion(version uint32, startTime, endTime uint32, bitmap *roaring.Bitmap) {
	if w.tagValueBuffer == nil {
		w.tagValueBuffer = bufpool.GetBuffer()
		w.tagValueWriter.SwitchBuffer(w.tagValueBuffer)
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
	w.tagValueWriter.PutUint32(version)
	// write startTime delta
	startTimeDelta := int64(startTime) - int64(version)
	w.tagValueWriter.PutVarint64(startTimeDelta)
	// write endTime delta
	endTimeDelta := int64(endTime) - int64(version)
	w.tagValueWriter.PutVarint64(endTimeDelta)
	// write bitmap length
	out, err := bitmap.MarshalBinary()
	if err != nil {
		invertedIndexFlusherLogger.Error("marshal bitmap failure", logger.Error(err))
	}
	w.tagValueWriter.PutUvarint64(uint64(len(out)))
	// write bitmap data
	w.tagValueWriter.PutBytes(out)
}

// bufferWithVersionCount is the value of trie-tree node
type bufferWithVersionCount struct {
	versionCount int
	buffer       *bytes.Buffer
}

// FlushTagValue indicate a VersionedTagValueDataBlock is done.
func (w *invertedIndexFlusher) FlushTagValue(tagValue string) {
	w.trie.Add(tagValue, bufferWithVersionCount{
		versionCount: w.versionCount,
		buffer:       w.tagValueBuffer})

	w.tagValueBuffer = nil
	w.versionCount = 0
}

func (w *invertedIndexFlusher) FlushTagKey(tagID uint32) error {
	if !w.resetDisabled {
		defer w.reset()
	}

	treeDataBlock := w.trie.MarshalBinary()
	// write startTime
	w.entrySetWriter.PutUint32(w.minStartTime)
	// write endTime
	w.entrySetWriter.PutUint32(w.maxEndTime)
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
	treeLength := stream.GetUVariantLength(uint64(len(treeDataBlock.labels))) + // labels length uvariant size
		len(treeDataBlock.labels) + // labels length
		stream.GetUVariantLength(uint64(len(isPrefixBlock))) + // isPrefixKey length uvariant size
		len(isPrefixBlock) + // isPrefixKey length
		stream.GetUVariantLength(uint64(len(LOUDSBlock))) + // LOUDSBlock length uvariantsize
		len(LOUDSBlock) // LOUDSBlock length
	// write tree length
	w.entrySetWriter.PutUvarint64(uint64(treeLength))
	// write labels length & labels
	w.entrySetWriter.PutUvarint64(uint64(len(treeDataBlock.labels)))
	w.entrySetWriter.PutBytes(treeDataBlock.labels)
	// write isPrefixKey length & bitmap
	w.entrySetWriter.PutUvarint64(uint64(len(isPrefixBlock)))
	w.entrySetWriter.PutBytes(isPrefixBlock)
	// write LOUDS length & bitmap
	w.entrySetWriter.PutUvarint64(uint64(len(LOUDSBlock)))
	w.entrySetWriter.PutBytes(LOUDSBlock)
	// write tagValueCount
	w.entrySetWriter.PutUvarint64(uint64(len(treeDataBlock.values)))
	// write all data length and versioned tagValue data blocks
	w.writeTagValueDataBlockTo(w.entrySetWriter, treeDataBlock)

	// write crc32 checksum
	data, _ := w.entrySetWriter.Bytes()
	w.entrySetWriter.PutUint32(crc32.ChecksumIEEE(data))
	data, _ = w.entrySetWriter.Bytes()
	return w.flusher.Add(tagID, data)
}

// writeTagValueDataBlockTo write tagValueDataBlocks to the writer.
func (w *invertedIndexFlusher) writeTagValueDataBlockTo(writer *stream.BufferWriter, treeDataBlock *trieTreeData) {
	// write lengths of all versioned tagValue data block
	for _, item := range treeDataBlock.values {
		it := item.(bufferWithVersionCount)
		// write all data length
		dataBlockLen := stream.GetUVariantLength(uint64(it.versionCount)) + // version count size
			len(it.buffer.Bytes()) // versionedTagValue blocks
		writer.PutUvarint64(uint64(dataBlockLen))
	}
	// write all versioned tagValue data block
	for _, item := range treeDataBlock.values {
		it := item.(bufferWithVersionCount)
		// write version count
		writer.PutUvarint64(uint64(it.versionCount))
		// write all versions of tagValue bitmaps
		writer.PutBytes(it.buffer.Bytes())
		// put buffer back to pool
		it.buffer.Reset()
		bufpool.PutBuffer(it.buffer)
	}
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (w *invertedIndexFlusher) Commit() error {
	w.reset()
	return w.flusher.Commit()
}

// reset resets the trie and buf
func (w *invertedIndexFlusher) reset() {
	w.trie.Reset()
	w.entrySetWriter.Reset()
}

package tblstore

import (
	"bytes"
	"hash/crc32"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bufpool"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./inverted_index_flusher.go -destination=./inverted_index_flusher_mock.go -package tblstore

var invertedIndexFlusherLogger = logger.GetLogger("tsdb", "InvertedIndexFlusher")

// InvertedIndexFlusher is a wrapper of kv.Builder, provides the ability to build a versioned series-id index table.
// The layout is available in `tsdb/doc.go`
type InvertedIndexFlusher interface {
	// FlushVersion writes a versioned bitmap to index table.
	FlushVersion(version series.Version, timeRange timeutil.TimeRange, bitmap *roaring.Bitmap)
	// FlushTagValue ends writing VersionedTagValueBlock in index table.
	FlushTagValue(tagValue string)
	// FlushTagKeyID ends writing entrySetBlock in index table.
	FlushTagKeyID(tagID uint32) error
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
		offsets:        encoding.NewDeltaBitPackingEncoder(),
	}
}

// invertedIndexFlusher implements InvertedIndexFlusher.
type invertedIndexFlusher struct {
	flusher        kv.Flusher
	trie           trieTreeBuilder
	entrySetWriter *stream.BufferWriter
	offsets        *encoding.DeltaBitPackingEncoder
	// time range
	minStartTime int64
	maxEndTime   int64
	// for tagValue data builder
	versionCount   int
	tagValueWriter *stream.BufferWriter
	tagValueBuffer *bytes.Buffer
}

// FlushVersion writes a versioned bitmap to index table.
func (w *invertedIndexFlusher) FlushVersion(
	version series.Version,
	timeRange timeutil.TimeRange,
	bitmap *roaring.Bitmap,
) {
	out, err := bitmap.MarshalBinary()
	if err != nil {
		invertedIndexFlusherLogger.Error("marshal bitmap failure", logger.Error(err))
	}
	w.flushVersion(version, timeRange, out)
}

// real flush-version method
func (w *invertedIndexFlusher) flushVersion(
	version series.Version,
	timeRange timeutil.TimeRange,
	bitmapData []byte,
) {
	if w.tagValueBuffer == nil {
		w.tagValueBuffer = bufpool.GetBuffer()
		w.tagValueWriter.SwitchBuffer(w.tagValueBuffer)
	}
	// count flushed versions
	w.versionCount++
	// update time range
	if timeRange.Start < w.minStartTime || w.minStartTime == 0 {
		w.minStartTime = timeRange.Start
	}
	if timeRange.End > w.maxEndTime || w.maxEndTime == 0 {
		w.maxEndTime = timeRange.End
	}
	// write version
	w.tagValueWriter.PutInt64(version.Int64())
	// write startTime delta
	startTimeDelta := (timeRange.Start - version.Int64()) / 1000 // seconds
	w.tagValueWriter.PutVarint64(startTimeDelta)
	// write endTime delta
	endTimeDelta := (timeRange.End - version.Int64()) / 1000 // seconds
	w.tagValueWriter.PutVarint64(endTimeDelta)
	// write bitmap length
	w.tagValueWriter.PutUvarint64(uint64(len(bitmapData)))
	// write bitmap data
	w.tagValueWriter.PutBytes(bitmapData)
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

// FlushTagKeyID ends writing entrySetBlock in index table.
func (w *invertedIndexFlusher) FlushTagKeyID(tagID uint32) error {
	defer w.reset()

	// write startTime
	w.entrySetWriter.PutInt64(w.minStartTime)
	// write endTime
	w.entrySetWriter.PutInt64(w.maxEndTime)

	treeDataBlock := w.trie.MarshalBinary()
	// write tree
	if err := w.writeTrieTree(treeDataBlock); err != nil {
		return err
	}
	// write tagValueData list
	w.writeTagValueData(treeDataBlock)
	// write offsets, footer
	w.writeOffsetsAndFooter()
	// write all
	data, _ := w.entrySetWriter.Bytes()
	return w.flusher.Add(tagID, data)
}

func (w *invertedIndexFlusher) writeTrieTree(treeDataBlock *trieTreeData) error {
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
	treeLength := stream.UvariantSize(uint64(len(treeDataBlock.labels))) + // labels length uvariant size
		len(treeDataBlock.labels) + // labels length
		stream.UvariantSize(uint64(len(isPrefixBlock))) + // isPrefixKey length uvariant size
		len(isPrefixBlock) + // isPrefixKey length
		stream.UvariantSize(uint64(len(LOUDSBlock))) + // LOUDSBlock length uvariantsize
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
	return nil
}

func (w *invertedIndexFlusher) writeTagValueData(treeDataBlock *trieTreeData) {
	// write tagValueCount
	w.entrySetWriter.PutUvarint64(uint64(len(treeDataBlock.values)))
	// write all data length and versioned tagValue data blocks
	for _, item := range treeDataBlock.values {
		it := item.(bufferWithVersionCount)
		// record this position
		w.offsets.Add(int32(w.entrySetWriter.Len()))
		// write version count
		w.entrySetWriter.PutUvarint64(uint64(it.versionCount))
		// write all versions of tagValue bitmaps
		w.entrySetWriter.PutBytes(it.buffer.Bytes())
		// put buffer back to pool
		it.buffer.Reset()
		bufpool.PutBuffer(it.buffer)
	}
}

func (w *invertedIndexFlusher) writeOffsetsAndFooter() {
	// offsets start position
	offsetsStartPos := w.entrySetWriter.Len()
	// write offsets
	w.entrySetWriter.PutBytes(w.offsets.Bytes())
	////////////////////////////////
	// footer
	////////////////////////////////
	// write offsets start position
	w.entrySetWriter.PutUint32(uint32(offsetsStartPos))
	// write crc32 checksum
	data, _ := w.entrySetWriter.Bytes()
	w.entrySetWriter.PutUint32(crc32.ChecksumIEEE(data))
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (w *invertedIndexFlusher) Commit() error {
	w.reset()
	return w.flusher.Commit()
}

// reset resets the trie and buf
func (w *invertedIndexFlusher) reset() {
	w.trie.Reset()
	w.offsets.Reset()
	w.entrySetWriter.Reset()
}

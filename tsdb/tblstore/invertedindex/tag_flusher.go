package invertedindex

import (
	"hash/crc32"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./tag_flusher.go -destination=./tag_flusher_mock.go -package invertedindex

// TagFlusher is a wrapper of kv.Builder, provides the ability to build a tag index table.
// The layout is available in `tsdb/doc.go`
type TagFlusher interface {
	// FlushTagValue ends writing trie tree in tag index table.
	FlushTagValue(tagValue string, tagValueID uint32)
	// FlushTagKeyID ends writing trie tree data in tag index table.
	FlushTagKeyID(tagID uint32, tagValueSeq uint32) error
	// Commit closes the writer, this will be called after writing all tagKeys.
	Commit() error
}

// NewTagFlusher returns a new TagFlusher
func NewTagFlusher(kvFlusher kv.Flusher) TagFlusher {
	return &tagFlusher{
		kvFlusher:      kvFlusher,
		entrySetWriter: stream.NewBufferWriter(nil),
		trie:           newTrieTree(),
		tagValueIDs:    encoding.NewFixedOffsetEncoder(),
	}
}

// tagFlusher implements TagFlusher.
type tagFlusher struct {
	kvFlusher      kv.Flusher
	trie           trieTreeBuilder
	entrySetWriter *stream.BufferWriter
	tagValueIDs    *encoding.FixedOffsetEncoder
}

// FlushTagValue writes the tag value into tag value prefix trie
func (w *tagFlusher) FlushTagValue(tagValue string, tagValueID uint32) {
	w.trie.Add(tagValue, tagValueID)
}

// FlushTagKeyID ends writing prefix trie in tag index table.
func (w *tagFlusher) FlushTagKeyID(tagID uint32, tagValueSeq uint32) error {
	defer w.reset()

	treeDataBlock := w.trie.MarshalBinary()
	// write tree
	if err := w.writeTrieTree(treeDataBlock); err != nil {
		return err
	}
	// write tagValueData list
	w.writeTagValueData(treeDataBlock)
	// write offsets, footer
	w.writeOffsetsAndFooter(tagValueSeq)
	// write all
	data, _ := w.entrySetWriter.Bytes()
	return w.kvFlusher.Add(tagID, data)
}

func (w *tagFlusher) writeTrieTree(treeDataBlock *trieTreeData) error {
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

func (w *tagFlusher) writeTagValueData(treeDataBlock *trieTreeData) {
	// write all data length and versioned tagValue data blocks
	for _, item := range treeDataBlock.values {
		tagValueID := item.(uint32)
		// record tag value id
		w.tagValueIDs.Add(int(tagValueID))
	}
}

func (w *tagFlusher) writeOffsetsAndFooter(tagValueSeq uint32) {
	// offsets start position
	offsetsStartPos := w.entrySetWriter.Len()
	// write all tag value ids
	w.entrySetWriter.PutBytes(w.tagValueIDs.MarshalBinary())
	////////////////////////////////
	// footer (tag value seq+tag value ids' offset+crc32 checksum)(4 bytes+4 bytes+4 bytes)
	////////////////////////////////
	w.entrySetWriter.PutUint32(tagValueSeq)
	// write tag value ids' start position
	w.entrySetWriter.PutUint32(uint32(offsetsStartPos))
	// write crc32 checksum
	data, _ := w.entrySetWriter.Bytes()
	w.entrySetWriter.PutUint32(crc32.ChecksumIEEE(data))
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (w *tagFlusher) Commit() error {
	w.reset()
	return w.kvFlusher.Commit()
}

// reset resets the trie and buf
func (w *tagFlusher) reset() {
	w.trie.Reset()
	w.tagValueIDs.Reset()
	w.entrySetWriter.Reset()
}

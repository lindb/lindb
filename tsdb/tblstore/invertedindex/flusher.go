package invertedindex

import (
	"hash/crc32"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package invertedindex

// Flusher is a wrapper of kv.Builder, provides the ability to build a inverted index table.
// The layout is available in `tsdb/doc.go`
type Flusher interface {
	// FlushInvertedIndex ends writing trie tree in tag index table.
	// !!!!NOTICE: need add tag value id in order. tag value id=0 store all series ids under this tag
	FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error
	// FlushTagKeyID ends writing trie tree data in tag index table.
	FlushTagKeyID(tagID uint32) error
	// Commit closes the writer, this will be called after writing all tag keys.
	Commit() error
}

// NewFlusher returns a new Flusher
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	return &flusher{
		kvFlusher:   kvFlusher,
		writer:      stream.NewBufferWriter(nil),
		tagValueIDs: roaring.New(),
		lowOffsets:  encoding.NewFixedOffsetEncoder(),
		highOffsets: encoding.NewFixedOffsetEncoder(),
	}
}

// flusher implements Flusher.
type flusher struct {
	kvFlusher   kv.Flusher
	tagValueIDs *roaring.Bitmap
	writer      *stream.BufferWriter
	highOffsets *encoding.FixedOffsetEncoder
	lowOffsets  *encoding.FixedOffsetEncoder
	highKey     uint16
}

// FlushInvertedIndex writes tag value id->series ids inverted index data
func (w *flusher) FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error {
	seriesData, err := encoding.BitmapMarshal(seriesIDs)
	if err != nil {
		return err
	}
	highKey := encoding.HighBits(tagValueID)
	if highKey != w.highKey {
		// flush data by diff high key
		w.flushTagValueBucket()
	}

	pos := w.writer.Len()
	// write series ids into data block
	w.writer.PutBytes(seriesData)
	w.lowOffsets.Add(pos)
	// add tag value id into index block
	w.tagValueIDs.Add(tagValueID)
	return nil
}

// flushTagValueBucket flushes data by bucket based on bitmap container
func (w *flusher) flushTagValueBucket() {
	if w.tagValueIDs.IsEmpty() {
		// maybe first high key not start with 0
		return
	}

	defer w.lowOffsets.Reset()

	pos := w.writer.Len()
	w.writer.PutBytes(w.lowOffsets.MarshalBinary())
	w.highOffsets.Add(pos)
}

// FlushTagKeyID ends writing tag inverted index data in index table.
func (w *flusher) FlushTagKeyID(tagID uint32) error {
	defer w.reset()

	// check if has pending tag value bucket not flush
	w.flushTagValueBucket()
	// write high offsets
	offsetPos := w.writer.Len()
	w.writer.PutBytes(w.highOffsets.MarshalBinary())
	// write tag value ids bitmap
	tagValueIDsBlock, err := encoding.BitmapMarshal(w.tagValueIDs)
	if err != nil {
		return err
	}
	tagValueIDsPos := w.writer.Len()
	w.writer.PutBytes(tagValueIDsBlock)
	////////////////////////////////
	// footer (tag value ids' offset+high level offsets+crc32 checksum)
	// (4 bytes + 4 bytes + 4 bytes)
	////////////////////////////////
	// write tag value ids' start position
	w.writer.PutUint32(uint32(tagValueIDsPos))
	// write offset block start position
	w.writer.PutUint32(uint32(offsetPos))
	// write crc32 checksum
	data, _ := w.writer.Bytes()
	w.writer.PutUint32(crc32.ChecksumIEEE(data))
	// write all
	data, _ = w.writer.Bytes()
	return w.kvFlusher.Add(tagID, data)
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (w *flusher) Commit() error {
	w.reset()
	return w.kvFlusher.Commit()
}

// reset resets the trie and buf
func (w *flusher) reset() {
	w.tagValueIDs.Clear()
	w.lowOffsets.Reset()
	w.highOffsets.Reset()
	w.writer.Reset()
	w.highKey = 0
}

package indextbl

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sync"

	"github.com/lindb/lindb/kv"

	"github.com/RoaringBitmap/roaring"
)

//go:generate mockgen -source ./series_flusher.go -destination=./series_flusher_mock.go -package indextbl

// VersionedTagKVEntrySet is a entrySet related a specific tagKey and version.
type VersionedTagKVEntrySet struct {
	Version  int64                      // series version
	EntrySet map[string]*roaring.Bitmap // tagValues bitmap
}

// SeriesIndexFlusher is a wrapper of kv.Builder, provides the ability to build a versioned series-id index table.
// The layout is available in `tsdb/doc.go`
type SeriesIndexFlusher interface {
	// FlushTagKey writes a version of the tagValues and related bitmap to the index table.
	FlushTagKey(tagID uint32, versionedTagKVSets []VersionedTagKVEntrySet) error
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
	trie           trieTreeINTF
	entrySetBuffer bytes.Buffer
	usingBuffers   []*bytes.Buffer
	bufferPool     sync.Pool
	VariableBuf    [8]byte
}

func (w *seriesIndexFlusher) getBuffer() *bytes.Buffer {
	return w.bufferPool.Get().(*bytes.Buffer)
}

// FlushTagKey writes a version of the tagValues and related bitmap to the index table.
func (w *seriesIndexFlusher) FlushTagKey(tagID uint32, versionedTagKVSets []VersionedTagKVEntrySet) error {
	defer w.reset()

	for _, tagKVSet := range versionedTagKVSets {
		for tagValue, bitmap := range tagKVSet.EntrySet {
			w.trie.Add(tagValue, tagKVSet.Version, bitmap)
		}
	}
	bin := w.trie.MarshalBinary()
	// write entrySet
	var (
		err error
		out []byte
	)
	// write labels length
	size := binary.PutUvarint(w.VariableBuf[:], uint64(len(bin.labels)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	// write labels
	_, _ = w.entrySetBuffer.Write(bin.labels)
	// write isPrefixKey length
	out, err = bin.isPrefixKey.MarshalBinary()
	if err != nil {
		return err
	}
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(out)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	// write isPrefixKey bitmap
	_, _ = w.entrySetBuffer.Write(out)
	// write LOUDS length
	out, err = bin.LOUDS.MarshalBinary()
	if err != nil {
		return err
	}
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(out)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	// write isPrefixKey bitmap
	_, _ = w.entrySetBuffer.Write(out)

	// write tagValueCount
	size = binary.PutUvarint(w.VariableBuf[:], uint64(len(bin.values)))
	_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])

	// write data lengths
	for _, versionedBitMapList := range bin.values {
		bufs, err := w.tagValueData2Buffer(versionedBitMapList)
		if err != nil {
			return err
		}
		var length = 0
		for _, buf := range bufs {
			length += buf.Len()
			w.usingBuffers = append(w.usingBuffers, buf)
		}
		// write data length
		size := binary.PutUvarint(w.VariableBuf[:], uint64(length))
		_, _ = w.entrySetBuffer.Write(w.VariableBuf[:size])
	}
	// compact buffers
	for _, usingBuffer := range w.usingBuffers {
		w.entrySetBuffer.Write(usingBuffer.Bytes())
	}
	// write crc32 checksum
	binary.BigEndian.PutUint32(w.VariableBuf[0:4], crc32.ChecksumIEEE(w.entrySetBuffer.Bytes()))
	w.entrySetBuffer.Write(w.VariableBuf[:4])

	return w.flusher.Add(tagID, w.entrySetBuffer.Bytes())
}

func (w *seriesIndexFlusher) tagValueData2Buffer(vbs versionedBitmaps) ([]*bytes.Buffer, error) {
	// storing tag value length data
	tagValueDataBuffer := w.getBuffer()
	// write version count
	size := binary.PutUvarint(w.VariableBuf[:], uint64(len(vbs)))
	_, _ = tagValueDataBuffer.Write(w.VariableBuf[:size])
	// write version and version delta
	lastVersion := int64(0)
	for _, vb := range vbs {
		size = binary.PutUvarint(w.VariableBuf[:], uint64(vb.version-lastVersion))
		_, _ = tagValueDataBuffer.Write(w.VariableBuf[:size])
		lastVersion = vb.version
	}
	// storing bitmap data
	bitMapBuffer := w.getBuffer()
	for _, vb := range vbs {
		output, err := vb.bitmap.MarshalBinary()
		if err != nil {
			return nil, err
		}
		// write version length
		size = binary.PutUvarint(w.VariableBuf[:], uint64(len(output)))
		_, _ = tagValueDataBuffer.Write(w.VariableBuf[:size])
		// write tag value bitmap
		bitMapBuffer.Write(output)
	}
	return []*bytes.Buffer{tagValueDataBuffer, bitMapBuffer}, nil
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
	for _, buffer := range w.usingBuffers {
		buffer.Reset()
		w.bufferPool.Put(buffer)
	}
	w.usingBuffers = w.usingBuffers[:0]
}

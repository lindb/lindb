package tagkeymeta

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"sort"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/trie"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package tagkeymeta

// Flusher is a wrapper of kv.Builder, provides the ability to build a tag index table.
// The layout is available in `tsdb/doc.go`
type Flusher interface {
	// EnsureSize ensures slice's capacity to meets the demand
	EnsureSize(size int)
	// FlushTagValue ends writing trie tree in tag index table.
	FlushTagValue(tagValue []byte, tagValueID uint32)
	// FlushTagKeyID ends writing trie tree data in tag index table.
	FlushTagKeyID(tagKeyID uint32, tagValueSeq uint32) error
	// Commit closes the writer, this will be called after writing all tagKeys.
	Commit() error
}

// NewFlusher returns a new TagFlusher
func NewFlusher(kvFlusher kv.Flusher) Flusher {
	return &flusher{
		kvFlusher:      kvFlusher,
		entrySetWriter: stream.NewBufferWriter(nil),
		idBitmap:       roaring.New(),
		rankOffsets:    encoding.NewFixedOffsetEncoder(),
		trieBuilder:    trie.NewBuilder(),
	}
}

// flusher implements Flusher.
type flusher struct {
	kvFlusher      kv.Flusher
	trieBuilder    trie.Builder
	entrySetWriter *stream.BufferWriter
	maxTagValueID  uint32
	// cached kv paris for building the fast succinct trie
	tagValueMapping tagValueMapping
	idBitmap        *roaring.Bitmap              // storing all tag-value ids
	rankOffsets     *encoding.FixedOffsetEncoder // storing all ranks of tag-ids on trie tree
}

// tagValueMapping sorts the ids based on the order in keys
type tagValueMapping struct {
	keys [][]byte // tag-values
	ids  [][]byte // tag-value ids([]byte)
	idRanks
}

// idRanks sorts the ranks based on the order in rawIDs
type idRanks struct {
	rawIDs []uint32 // tag-value ids(uint32)
	ranks  []int    // tag-value's rank list on the tree
}

func (ir idRanks) Len() int           { return len(ir.rawIDs) }
func (ir idRanks) Less(i, j int) bool { return ir.rawIDs[i] < ir.rawIDs[j] }
func (ir idRanks) Swap(i, j int) {
	ir.rawIDs[i], ir.rawIDs[j] = ir.rawIDs[j], ir.rawIDs[i]
	ir.ranks[i], ir.ranks[j] = ir.ranks[j], ir.ranks[i]
}

func (m tagValueMapping) Len() int { return len(m.keys) }
func (m tagValueMapping) Less(i, j int) bool {
	return bytes.Compare(m.keys[i], m.keys[j]) < 0
}
func (m tagValueMapping) Swap(i, j int) {
	m.keys[i], m.keys[j] = m.keys[j], m.keys[i]
	m.ids[i], m.ids[j] = m.ids[j], m.ids[i]
	m.rawIDs[i], m.rawIDs[j] = m.rawIDs[j], m.rawIDs[i]
}

func (m tagValueMapping) SortByKeys() {
	sort.Sort(m)
}

// SortByKeys should be called before
func (m *tagValueMapping) SortByRawIDs() {
	// insert ranks
	for i := 0; i < len(m.keys); i++ {
		m.ranks = append(m.ranks, i)
	}
	sort.Sort(m.idRanks)
}

func (m *tagValueMapping) reset() {
	m.keys = m.keys[:0]
	m.ids = m.ids[:0]
	m.rawIDs = m.rawIDs[:0]
	m.ranks = m.ranks[:0]
}

func (m *tagValueMapping) ensureSize(size int) {
	if cap(m.keys) < size {
		m.keys = make([][]byte, size)[:0]
		m.ids = make([][]byte, size)[:0]
		m.rawIDs = make([]uint32, size)[:0]
		m.ranks = make([]int, size)[:0]
	}
}

func (tf *flusher) EnsureSize(size int) { tf.tagValueMapping.ensureSize(size) }

func (tf *flusher) FlushTagValue(tagValue []byte, tagValueID uint32) {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, tagValueID)
	tf.tagValueMapping.keys = append(tf.tagValueMapping.keys, tagValue)
	tf.tagValueMapping.ids = append(tf.tagValueMapping.ids, buf)

	if tagValueID > tf.maxTagValueID {
		tf.maxTagValueID = tagValueID
	}
	tf.tagValueMapping.rawIDs = append(tf.tagValueMapping.rawIDs, tagValueID)
}

// FlushTagKeyID ends writing prefix trie in tag index table.
func (tf *flusher) FlushTagKeyID(tagKeyID uint32, tagValueSeq uint32) error {
	defer tf.reset()

	if len(tf.tagValueMapping.keys) == 0 {
		return nil
	}
	// pre-sort for building trie
	tf.tagValueMapping.SortByKeys()
	// build trie
	tree := tf.trieBuilder.Build(
		tf.tagValueMapping.keys,
		tf.tagValueMapping.ids,
		uint32(encoding.Uint32MinWidth(tf.maxTagValueID)))

	// writing to buffer in memory won't raise error
	_ = tree.WriteTo(tf.entrySetWriter)
	tf.tagValueMapping.SortByRawIDs()
	// remember bitmap position
	bitmapPosition := tf.entrySetWriter.Len()
	// flush bitmap
	tf.idBitmap.AddMany(tf.tagValueMapping.rawIDs)
	// writing to buffer in memory won't raise error
	_, _ = tf.idBitmap.WriteTo(tf.entrySetWriter)
	// flush offsets
	offsetsPosition := tf.entrySetWriter.Len()
	for _, rank := range tf.tagValueMapping.ranks {
		tf.rankOffsets.Add(rank)
	}

	// writing to buffer in memory won't raise error
	_, _ = tf.entrySetWriter.Write(tf.rankOffsets.MarshalBinary())

	// footer
	// flush bitmap position
	tf.entrySetWriter.PutUint32(uint32(bitmapPosition))
	// flush offsets position
	tf.entrySetWriter.PutUint32(uint32(offsetsPosition))
	// flush tag-value sequence
	tf.entrySetWriter.PutUint32(tagValueSeq)
	// write crc32 checksum
	data, _ := tf.entrySetWriter.Bytes()
	tf.entrySetWriter.PutUint32(crc32.ChecksumIEEE(data))

	data, _ = tf.entrySetWriter.Bytes()
	return tf.kvFlusher.Add(tagKeyID, data)
}

// Commit closes the writer, this will be called after writing all tagKeys.
func (tf *flusher) Commit() error {
	tf.reset()
	return tf.kvFlusher.Commit()
}

// reset resets the underlying data structures
func (tf *flusher) reset() {
	tf.maxTagValueID = 0
	tf.idBitmap.Clear()
	tf.entrySetWriter.Reset()
	tf.rankOffsets.Reset()

	tf.tagValueMapping.reset()
	tf.trieBuilder.Reset()
}

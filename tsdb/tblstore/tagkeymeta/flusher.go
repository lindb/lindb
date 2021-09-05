// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tagkeymeta

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
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
	// used for merging
	commitTagKeyID() error
	// Closer closes the writer, this will be called after writing all tagKeys.
	io.Closer
}

// NewFlusher returns a new TagFlusher
func NewFlusher(kvFlusher kv.Flusher) (Flusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	f := &flusher{
		kvFlusher: kvFlusher,
		kvWriter:  kvWriter,
	}
	f.Level2.trieBuilder = trie.NewBuilder()
	f.Level2.tagValueIDsBitmap = roaring.New()
	f.Level2.rankOffsets = encoding.NewFixedOffsetEncoder(false)
	return f, nil
}

// flusher implements Flusher.
type flusher struct {
	// Level1 flusher
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter

	//  ━━━━━━━━━━━━━━━━━━━━━━━Layout of TagKeys Meta Table━━━━━━━━━━━━━━━━━━━━━━━━
	//
	//                    Level1
	//                    +---------+---------+---------+---------+---------+---------+
	//                    │ TagKey  │ TagKey  │ TagKey  │ Offsets │ Bitmap  │ Footer  │
	//                    │  Meta   │  Meta   │  Meta   │         │         │         │
	//                    +---------+---------+---------+---------+---------+---------+
	//                   /           \                  |         |
	//                  /             \                 |          \
	//                 /                \              /            \
	//                /                   \           /               \
	//   +-----------+                     |        /                   \
	//  /                     Level2       |       |                     |
	// v--------+--------+--------+--------v       v--------+---+--------v
	// │  Trie  │TagValue│ Offsets│ Footer │       │ Offset │...│ Offset │
	// │  Tree  │IDBitmap│        │        │       │        │   │        │
	// +--------+--------+--------+--------+       +--------+---+--------+
	//
	//
	// Level1(KV table: TagKeyID -> TagKeyMeta data)
	//
	// Level2(Footer)
	//
	// ┌───────────────────────────────────────────┐
	// │                 Footer                    │
	// ├──────────┬──────────┬──────────┬──────────┤
	// │  BitMap  │  Offsets │ TagValue │  CRC32   │
	// │ Position │ Position │ Sequence │ CheckSum │
	// ├──────────┼──────────┼──────────┼──────────┤
	// │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │ 4 Bytes  │
	// └──────────┴──────────┴──────────┴──────────┘
	Level2 struct {
		trieBuilder       trie.Builder
		maxTagValueID     uint32
		tagValueMapping   tagValueMapping              // cached kv paris for building succinct trie
		tagValueIDsBitmap *roaring.Bitmap              // storing all tag-value ids
		rankOffsets       *encoding.FixedOffsetEncoder // storing all ranks of tag-ids on trie tree
		footer            [tagFooterSize]byte
	}
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

func (m tagValueMapping) Len() int           { return len(m.keys) }
func (m tagValueMapping) Less(i, j int) bool { return bytes.Compare(m.keys[i], m.keys[j]) < 0 }
func (m tagValueMapping) Swap(i, j int) {
	m.keys[i], m.keys[j] = m.keys[j], m.keys[i]
	m.ids[i], m.ids[j] = m.ids[j], m.ids[i]
	m.rawIDs[i], m.rawIDs[j] = m.rawIDs[j], m.rawIDs[i]
}

func (m tagValueMapping) SortByKeys() {
	sort.Sort(m)
}

// SortByRawIDs should be called before
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
		m.keys = make([][]byte, 0, size)
		m.ids = make([][]byte, 0, size)
		m.rawIDs = make([]uint32, 0, size)
		m.ranks = make([]int, 0, size)
	}
}

func (tf *flusher) EnsureSize(size int) { tf.Level2.tagValueMapping.ensureSize(size) }

func (tf *flusher) FlushTagValue(tagValue []byte, tagValueID uint32) {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, tagValueID)
	tf.Level2.tagValueMapping.keys = append(tf.Level2.tagValueMapping.keys, tagValue)
	tf.Level2.tagValueMapping.ids = append(tf.Level2.tagValueMapping.ids, buf)

	if tagValueID > tf.Level2.maxTagValueID {
		tf.Level2.maxTagValueID = tagValueID
	}
	tf.Level2.tagValueMapping.rawIDs = append(tf.Level2.tagValueMapping.rawIDs, tagValueID)
}

// FlushTagKeyID ends writing prefix trie in tag index table.
func (tf *flusher) FlushTagKeyID(tagKeyID uint32, tagValueSeq uint32) error {
	defer tf.resetLevel2()

	if len(tf.Level2.tagValueMapping.keys) == 0 {
		return nil
	}
	tf.kvWriter.Prepare(tagKeyID)

	// pre-sort for building trie
	tf.Level2.tagValueMapping.SortByKeys()
	// build trie
	tree := tf.Level2.trieBuilder.Build(
		tf.Level2.tagValueMapping.keys,
		tf.Level2.tagValueMapping.ids,
		uint32(encoding.Uint32MinWidth(tf.Level2.maxTagValueID)))

	if err := tree.Write(tf.kvWriter); err != nil {
		return err
	}

	tf.Level2.tagValueMapping.SortByRawIDs()
	// remember bitmap position
	tagValueBitmapAt := tf.kvWriter.Size()
	// flush bitmap
	tf.Level2.tagValueIDsBitmap.AddMany(tf.Level2.tagValueMapping.rawIDs)
	if _, err := tf.Level2.tagValueIDsBitmap.WriteTo(tf.kvWriter); err != nil {
		return err
	}
	// build offsets
	offsetsAt := tf.kvWriter.Size()
	for _, rank := range tf.Level2.tagValueMapping.ranks {
		tf.Level2.rankOffsets.Add(rank)
	}

	// write offsets
	if err := tf.Level2.rankOffsets.Write(tf.kvWriter); err != nil {
		return err
	}

	// footer
	// flush bitmap position
	binary.LittleEndian.PutUint32(tf.Level2.footer[0:4], tagValueBitmapAt)
	// flush offsets position
	binary.LittleEndian.PutUint32(tf.Level2.footer[4:8], offsetsAt)
	// flush tag-value sequence
	binary.LittleEndian.PutUint32(tf.Level2.footer[8:12], tagValueSeq)
	// write crc32 checksum
	binary.LittleEndian.PutUint32(tf.Level2.footer[12:16], tf.kvWriter.CRC32CheckSum())

	if _, err := tf.kvWriter.Write(tf.Level2.footer[:]); err != nil {
		return err
	}
	return tf.commitTagKeyID()
}

func (tf *flusher) commitTagKeyID() error { return tf.kvWriter.Commit() }

func (tf *flusher) Close() error {
	return tf.kvFlusher.Commit()
}

// reset resets the underlying data structures
func (tf *flusher) resetLevel2() {
	tf.Level2.trieBuilder.Reset()
	tf.Level2.maxTagValueID = 0
	tf.Level2.tagValueMapping.reset()
	tf.Level2.tagValueIDsBitmap.Clear()
	tf.Level2.rankOffsets.Reset()
}

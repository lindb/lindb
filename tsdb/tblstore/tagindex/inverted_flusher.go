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

package tagindex

import (
	"encoding/binary"
	"io"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./inverted_flusher.go -destination=./inverted_flusher_mock.go -package tagindex

// InvertedFlusher is a wrapper of kv.Builder, provides the ability to build a inverted index table.
// The layout is available in `tsdb/doc.go`
type InvertedFlusher interface {
	// PrepareTagKey should be called firstly
	PrepareTagKey(tagKeyID uint32)
	// FlushInvertedIndex writes tag value id->series ids inverted index data,
	// !!!!! NOTICE: need add tag value id in order.
	FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error
	// CommitTagKey ends writing tag inverted index data in index table.
	CommitTagKey() error
	// Close closes the writer, this will be called after writing all tag keys.
	io.Closer
}

// NewInvertedFlusher returns a new InvertedFlusher
func NewInvertedFlusher(kvFlusher kv.Flusher) (InvertedFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	iFlusher := &invertedFlusher{
		kvFlusher: kvFlusher,
		kvWriter:  kvWriter,
	}
	iFlusher.Level2.isHighKeySetEver = false
	iFlusher.Level2.highOffsets = encoding.NewFixedOffsetEncoder(true)
	iFlusher.Level2.tagValueIDs = roaring.New()

	iFlusher.Level3.lowOffsets = encoding.NewFixedOffsetEncoder(true)
	return iFlusher, nil
}

// invertedFlusher implements InvertedFlusher.
type invertedFlusher struct {
	// Level1
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter

	// Level2 (KV table: TagKey Inverted)
	// each entry is a tag value bucket ordered by roaring high key
	// Resets it after writing a tag key
	// v--------+--------+--------+--------+--------v
	// │TagValue│TagValue│TagValue|SeriesID│ Footer │
	// │ Bucket │ Bucket │Bitmap  |Offsets │        │
	// +--------+--------+--------+--------+--------+
	Level2 struct {
		// highKeySetEver symbols if highKey has been set before
		// set it to true after any tagvalue's bitmap flushed
		isHighKeySetEver bool
		highOffsets      *encoding.FixedOffsetEncoder
		// highKey is a the higher 16 bits of seriesIDs.
		// boundary for bucket
		highKey     uint16
		footer      [indexFooterSize]byte
		tagValueIDs *roaring.Bitmap
	}
	// TagValueBucket
	// v--------+--------+--------+--------v
	// │SeriesID│SeriesID│ LowKey │Offsets │
	// │ Bitmap │ Bitmap │ Offsets│Length  │
	// +--------+--------+--------+--------+
	// each entry is a series id bitmap ordered by low key
	// Offsets Length is a little-endian uvariant number
	Level3 struct {
		// startAt is the absolute position in Level2's SeriesEntry
		startAt int
		// scratch for uvariant encoding offsets marshal size
		scratch [binary.MaxVarintLen64]byte
		// lowOffsets holds distances between startAt and position of specified tag value id's bitmap
		lowOffsets *encoding.FixedOffsetEncoder
	}
}

func (w *invertedFlusher) PrepareTagKey(tagKeyID uint32) {
	w.kvWriter.Prepare(tagKeyID)
}

func (w *invertedFlusher) flushLevel2TagValueBucket() error {
	bucketHasData := int(w.kvWriter.Size())-w.Level3.startAt > 0
	if !bucketHasData {
		return nil
	}
	// start of offsets
	beforeLen := w.kvWriter.Size()
	if err := w.Level3.lowOffsets.Write(w.kvWriter); err != nil {
		return err
	}
	// write level3's length of low offsets
	writtenLen := stream.PutUvariantLittleEndian(w.Level3.scratch[:], uint64(w.kvWriter.Size()-beforeLen))
	_, err := w.kvWriter.Write(w.Level3.scratch[:writtenLen])
	return err
}

// FlushInvertedIndex writes tag value id->series ids inverted index data
func (w *invertedFlusher) FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error {
	// first occurrence
	highKey := encoding.HighBits(tagValueID)
	if !w.Level2.isHighKeySetEver {
		w.Level2.isHighKeySetEver = true
		w.Level2.highKey = highKey
		w.Level2.highOffsets.Add(0)
	}
	if highKey != w.Level2.highKey {
		if err := w.flushLevel2TagValueBucket(); err != nil {
			return err
		}
		w.Level2.highKey = highKey
		w.Level2.highOffsets.Add(int(w.kvWriter.Size()))
		// resets level3 context
		w.Level3.lowOffsets.Reset()
		w.Level3.startAt = int(w.kvWriter.Size())
	}
	// flush bitmap
	bitmapAt := w.kvWriter.Size()
	if _, err := seriesIDs.WriteTo(w.kvWriter); err != nil {
		return err
	}
	w.Level2.tagValueIDs.Add(tagValueID)
	w.Level3.lowOffsets.Add(int(bitmapAt) - w.Level3.startAt)
	return nil
}

// CommitTagKey ends writing tag inverted index data in index table.
func (w *invertedFlusher) CommitTagKey() error {
	defer w.reset()
	// empty tagvalue ids
	if w.Level2.tagValueIDs.IsEmpty() {
		return nil
	}
	if err := w.flushLevel2TagValueBucket(); err != nil {
		return err
	}
	// bitmap position
	tagValueBitmapAt := w.kvWriter.Size()
	if _, err := w.Level2.tagValueIDs.WriteTo(w.kvWriter); err != nil {
		return err
	}
	// offsets position
	offsetsAt := w.kvWriter.Size()
	// write offsets
	if err := w.Level2.highOffsets.Write(w.kvWriter); err != nil {
		return err
	}
	// footer (tag value bitmap
	//         high key offset position
	//         crc32 checksum)
	// (4 bytes + 4 bytes + 4 bytes)
	binary.LittleEndian.PutUint32(w.Level2.footer[0:4], tagValueBitmapAt)
	binary.LittleEndian.PutUint32(w.Level2.footer[4:8], offsetsAt)
	binary.LittleEndian.PutUint32(w.Level2.footer[8:12], w.kvWriter.CRC32CheckSum())
	// write footer
	if _, err := w.kvWriter.Write(w.Level2.footer[:]); err != nil {
		return err
	}
	return w.kvWriter.Commit()
}

// Close closes the writer, this will be called after writing all tagKeys.
func (w *invertedFlusher) Close() error {
	return w.kvFlusher.Commit()
}

// reset resets the trie and buf
func (w *invertedFlusher) reset() {
	w.Level2.isHighKeySetEver = false
	w.Level2.highOffsets.Reset()
	w.Level2.highKey = 0
	w.Level2.tagValueIDs.Clear()

	w.Level3.startAt = 0
	w.Level3.lowOffsets.Reset()
}

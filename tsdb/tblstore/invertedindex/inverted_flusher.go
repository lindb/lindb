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

package invertedindex

import (
	"hash/crc32"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./inverted_flusher.go -destination=./inverted_flusher_mock.go -package invertedindex

// InvertedFlusher is a wrapper of kv.Builder, provides the ability to build a inverted index table.
// The layout is available in `tsdb/doc.go`
type InvertedFlusher interface {
	// FlushInvertedIndex writes tag value id->series ids inverted index data,
	// !!!!! NOTICE: need add tag value id in order.
	FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error
	// FlushTagKeyID ends writing tag inverted index data in index table.
	FlushTagKeyID(tagID uint32) error
	// Commit closes the writer, this will be called after writing all tag keys.
	Commit() error
}

// NewInvertedFlusher returns a new InvertedFlusher
func NewInvertedFlusher(kvFlusher kv.Flusher) InvertedFlusher {
	return &invertedFlusher{
		kvFlusher:   kvFlusher,
		writer:      stream.NewBufferWriter(nil),
		tagValueIDs: roaring.New(),
		lowOffsets:  encoding.NewFixedOffsetEncoder(),
		highOffsets: encoding.NewFixedOffsetEncoder(),
	}
}

// invertedFlusher implements InvertedFlusher.
type invertedFlusher struct {
	kvFlusher   kv.Flusher
	tagValueIDs *roaring.Bitmap
	writer      *stream.BufferWriter
	highOffsets *encoding.FixedOffsetEncoder
	lowOffsets  *encoding.FixedOffsetEncoder
	highKey     uint16
}

// FlushInvertedIndex writes tag value id->series ids inverted index data
func (w *invertedFlusher) FlushInvertedIndex(tagValueID uint32, seriesIDs *roaring.Bitmap) error {
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
func (w *invertedFlusher) flushTagValueBucket() {
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
func (w *invertedFlusher) FlushTagKeyID(tagID uint32) error {
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
func (w *invertedFlusher) Commit() error {
	w.reset()
	return w.kvFlusher.Commit()
}

// reset resets the trie and buf
func (w *invertedFlusher) reset() {
	w.tagValueIDs.Clear()
	w.lowOffsets.Reset()
	w.highOffsets.Reset()
	w.writer.Reset()
	w.highKey = 0
}

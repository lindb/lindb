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

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./forward_flusher.go -destination=./forward_flusher_mock.go -package invertedindex

// ForwardFlusher represents forward index invertedFlusher which flushes series id => tag value id mapping
// The layout is available in `tsdb/doc.go`
type ForwardFlusher interface {
	// FlushForwardIndex flushes tag value ids by bitmap container
	FlushForwardIndex(tagValueIDs []uint32)
	// FlushTagKeyID ends writing series ids in tag index table.
	FlushTagKeyID(tagID uint32, seriesIDs *roaring.Bitmap) error
	// Commit closes the writer, this will be called after writing all tag keys.
	Commit() error
}

// forwardFlusher implements ForwardFlusher interface
type forwardFlusher struct {
	tagValueIDs *encoding.DeltaBitPackingEncoder // temp store tag value ids for encoding
	offsets     *encoding.FixedOffsetEncoder     // store offset that is tag value ids of one container
	writer      *stream.BufferWriter

	kvFlusher kv.Flusher
}

// NewForwardFlusher creates a forward index invertedFlusher
func NewForwardFlusher(kvFlusher kv.Flusher) ForwardFlusher {
	return &forwardFlusher{
		writer:      stream.NewBufferWriter(nil),
		tagValueIDs: encoding.NewDeltaBitPackingEncoder(),
		offsets:     encoding.NewFixedOffsetEncoder(true),
		kvFlusher:   kvFlusher,
	}
}

// FlushForwardIndex flushes tag value ids by bitmap container
func (f *forwardFlusher) FlushForwardIndex(tagValueIDs []uint32) {
	defer f.tagValueIDs.Reset()

	for _, tagValueID := range tagValueIDs {
		f.tagValueIDs.Add(int32(tagValueID))
	}
	offset := f.writer.Len()
	f.writer.PutBytes(f.tagValueIDs.Bytes()) // write tag value ids
	f.offsets.Add(offset)                    // add tag value ids' offset
}

// FlushTagKeyID ends writing series ids in tag index table.
func (f *forwardFlusher) FlushTagKeyID(tagID uint32, seriesIDs *roaring.Bitmap) error {
	defer f.reset()

	// write offsets
	offsetPos := f.writer.Len()
	f.writer.PutBytes(f.offsets.MarshalBinary())
	// write series ids bitmap
	seriesIDsBlock, err := encoding.BitmapMarshal(seriesIDs)
	if err != nil {
		return err
	}
	seriesIDsPos := f.writer.Len()
	f.writer.PutBytes(seriesIDsBlock)
	////////////////////////////////
	// footer (series ids' offset + offsets + crc32 checksum)
	// (4 bytes + 4 bytes + 4 bytes)
	////////////////////////////////
	// write tag value ids' start position
	f.writer.PutUint32(uint32(seriesIDsPos))
	// write offset block start position
	f.writer.PutUint32(uint32(offsetPos))
	// write crc32 checksum
	data, _ := f.writer.Bytes()
	f.writer.PutUint32(crc32.ChecksumIEEE(data))
	// write all
	data, _ = f.writer.Bytes()
	return f.kvFlusher.Add(tagID, data)
}

// Commit closes the writer, this will be called after writing all tag keys.
func (f *forwardFlusher) Commit() error {
	f.reset()
	return f.kvFlusher.Commit()
}

// reset resets the tag value ids and buf
func (f *forwardFlusher) reset() {
	f.tagValueIDs.Reset()
	f.offsets.Reset()
	f.writer.Reset()
}

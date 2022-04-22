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
)

//go:generate mockgen -source ./forward_flusher.go -destination=./forward_flusher_mock.go -package tagindex

// ForwardFlusher represents forward index invertedFlusher which flushes series id => tag value id mapping
// The layout is available in `tsdb/doc.go`
type ForwardFlusher interface {
	// PrepareTagKey should be called firstly
	PrepareTagKey(tagKeyID uint32)
	// FlushForwardIndex flushes tag value ids by bitmap container
	FlushForwardIndex(tagValueIDs []uint32) error
	// CommitTagKey ends writing series ids in tag index table.
	CommitTagKey(seriesIDs *roaring.Bitmap) error
	// Closer closes the writer, this will be called after writing all tag keys.
	io.Closer
}

// forwardFlusher implements ForwardFlusher interface
type forwardFlusher struct {
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter
	// level2
	tagValueIDs *encoding.DeltaBitPackingEncoder // temp store tag value ids for encoding
	offsets     *encoding.FixedOffsetEncoder     // store offset that is tag value ids of one container
	footer      [indexFooterSize]byte
}

// NewForwardFlusher creates a forward index invertedFlusher
func NewForwardFlusher(kvFlusher kv.Flusher) (ForwardFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &forwardFlusher{
		kvFlusher:   kvFlusher,
		kvWriter:    kvWriter,
		tagValueIDs: encoding.NewDeltaBitPackingEncoder(),
		offsets:     encoding.NewFixedOffsetEncoder(true),
	}, nil
}

func (f *forwardFlusher) PrepareTagKey(tagKeyID uint32) {
	f.kvWriter.Prepare(tagKeyID)
}

// FlushForwardIndex flushes tag value ids by bitmap container
func (f *forwardFlusher) FlushForwardIndex(tagValueIDs []uint32) error {
	defer f.tagValueIDs.Reset()

	for _, tagValueID := range tagValueIDs {
		f.tagValueIDs.Add(int32(tagValueID))
	}
	offset := f.kvWriter.Size()
	// write tag value ids
	if _, err := f.kvWriter.Write(f.tagValueIDs.Bytes()); err != nil {
		return err
	}
	// add tag value ids' offset
	f.offsets.Add(int(offset))
	return nil
}

func (f *forwardFlusher) CommitTagKey(seriesIDs *roaring.Bitmap) error {
	defer f.reset()

	// write series ids bitmap
	seriesIDAt := f.kvWriter.Size()
	if _, err := seriesIDs.WriteTo(f.kvWriter); err != nil {
		return err
	}
	// write offsets
	offsetsAt := f.kvWriter.Size()
	if err := f.offsets.Write(f.kvWriter); err != nil {
		return err
	}
	////////////////////////////////
	// footer (series ids' offset + offsets + crc32 checksum)
	// (4 bytes + 4 bytes + 4 bytes)
	////////////////////////////////
	// write tag value ids' start position
	binary.LittleEndian.PutUint32(f.footer[0:4], seriesIDAt)
	// write offset block start position
	binary.LittleEndian.PutUint32(f.footer[4:8], offsetsAt)
	// write crc32 checksum
	binary.LittleEndian.PutUint32(f.footer[8:12], f.kvWriter.CRC32CheckSum())
	if _, err := f.kvWriter.Write(f.footer[:]); err != nil {
		return err
	}
	return f.kvWriter.Commit()
}

func (f *forwardFlusher) Close() error {
	f.reset()
	return f.kvFlusher.Commit()
}

func (f *forwardFlusher) reset() {
	f.tagValueIDs.Reset()
	f.offsets.Reset()
}

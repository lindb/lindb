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

package v1

import (
	"io"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

//go:generate mockgen -source ./forward_flusher.go -destination=./forward_flusher_mock.go -package v1

// ForwardIndexFlusher represents forward index flusher.
type ForwardIndexFlusher interface {
	io.Closer

	// Prepare prepares tag key id for writing forward index.
	Prepare(tagKeyID uint32)
	// WriteSeriesIDs writes all series ids for tag key.
	WriteSeriesIDs(seriesIDs *roaring.Bitmap) error
	// WriteTagValueIDs writes tag value ids container based on bitmap container.
	WriteTagValueIDs(tagValueIDs []uint32) error
	// Commit completes index write.
	Commit() error
}

// forwardIndexFlusher implements ForwardIndexFlusher interface.
type forwardIndexFlusher struct {
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter
}

// NewForwardIndexFlusher creates a ForwardIndexFlusher instance.
func NewForwardIndexFlusher(kvFlusher kv.Flusher) (ForwardIndexFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &forwardIndexFlusher{
		kvWriter:  kvWriter,
		kvFlusher: kvFlusher,
	}, nil
}

// Prepare prepares tag key id for writing forward index.
func (f *forwardIndexFlusher) Prepare(tagKeyID uint32) {
	f.kvWriter.Prepare(tagKeyID)
}

// WriteSeriesIDs writes all series ids for tag key.
func (f *forwardIndexFlusher) WriteSeriesIDs(seriesIDs *roaring.Bitmap) error {
	if _, err := seriesIDs.WriteTo(f.kvWriter); err != nil {
		return err
	}
	return nil
}

// WriteTagValueIDs writes tag value ids container based on bitmap container.
func (f *forwardIndexFlusher) WriteTagValueIDs(tagValueIDs []uint32) error {
	// write tag value ids container
	if _, err := f.kvWriter.Write(encoding.U32SliceToBytes(tagValueIDs)); err != nil {
		return err
	}
	return nil
}

// Commit completes index write.
func (f *forwardIndexFlusher) Commit() error {
	return f.kvWriter.Commit()
}

// Close closes kv flusher.
func (f *forwardIndexFlusher) Close() error {
	return f.kvFlusher.Commit()
}

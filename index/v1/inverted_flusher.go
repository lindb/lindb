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
)

//go:generate mockgen -source ./inverted_flusher.go -destination=./inverted_flusher_mock.go -package v1

// InvertedIndexFlusher represents inverted index flusher.
type InvertedIndexFlusher interface {
	io.Closer

	// Prepare prepares inverted index key.
	Prepare(key uint32)
	// Write writes series ids into index.
	Write(seriesIDs *roaring.Bitmap) error
	// Commit completes index write.
	Commit() error
}

// invertedIndexFlusher implements InvertedIndexFlusher interface.
type invertedIndexFlusher struct {
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter
}

// NewInvertedIndexFlusher creates an InvertedIndexFlusher instance.
func NewInvertedIndexFlusher(kvFlusher kv.Flusher) (InvertedIndexFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &invertedIndexFlusher{
		kvWriter:  kvWriter,
		kvFlusher: kvFlusher,
	}, nil
}

// Prepare prepares inverted index key.
func (f *invertedIndexFlusher) Prepare(key uint32) {
	f.kvWriter.Prepare(key)
}

// Write writes series ids into index.
func (f *invertedIndexFlusher) Write(seriesIDs *roaring.Bitmap) error {
	_, err := seriesIDs.WriteTo(f.kvWriter)
	return err
}

// Commit completes index write.
func (f *invertedIndexFlusher) Commit() error {
	return f.kvWriter.Commit()
}

// Close close kv flusher.
func (f *invertedIndexFlusher) Close() error {
	return f.kvFlusher.Commit()
}

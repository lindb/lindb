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

	"github.com/lindb/lindb/index/model"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
)

//go:generate mockgen -source ./index_kv_flusher.go -destination=./index_kv_flusher_mock.go -package v1

// IndexKVFlusher represents index kv store flusher.
type IndexKVFlusher interface {
	io.Closer

	// PrepareBucket prepares bucket for writing.
	PrepareBucket(bucket uint32)
	// WriteKVs writes key/value pairs into kv store.
	WriteKVs(keys [][]byte, ids []uint32) error
	// CommitBucket completes bucket write.
	CommitBucket() error
}

// indexKVFlusher implements IndexKVFlusher interface.
type indexKVFlusher struct {
	kvFlusher kv.Flusher
	kvWriter  table.StreamWriter

	builder *model.TrieBucketBuilder
}

// NewIndexKVFlusher creates an IndexKVFlusher instance.
func NewIndexKVFlusher(blockSize int, kvFlusher kv.Flusher) (IndexKVFlusher, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &indexKVFlusher{
		kvFlusher: kvFlusher,
		kvWriter:  kvWriter,
		builder:   model.NewTrieBucketBuilder(blockSize, kvWriter),
	}, nil
}

// PrepareBucket prepares bucket for writing.
func (f *indexKVFlusher) PrepareBucket(bucket uint32) {
	f.kvWriter.Prepare(bucket)
}

// WriteKVs writes key/value pairs into kv store.
func (f *indexKVFlusher) WriteKVs(keys [][]byte, ids []uint32) error {
	return f.builder.Write(keys, ids)
}

// CommitBucket completes bucket write.
func (f *indexKVFlusher) CommitBucket() error {
	return f.kvWriter.Commit()
}

// Close completes kv flush.
func (f *indexKVFlusher) Close() error {
	return f.kvFlusher.Commit()
}

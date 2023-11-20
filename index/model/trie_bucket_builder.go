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

package model

import (
	"encoding/binary"
	"io"
	"sort"

	"github.com/lindb/lindb/pkg/trie"
)

// for testing
var (
	newTrieBuilder = trie.NewBuilder
)

// TrieBucketBuilder represents trie bucket builder.
type TrieBucketBuilder struct {
	blockSize int
	writer    io.Writer
	// write context(reuse), need reset after write entry
	builder trie.Builder
	sizeBuf []byte
}

// NewTrieBucketBuilder creates a TrieBucketBuilder instance.
func NewTrieBucketBuilder(blockSize int, writer io.Writer) *TrieBucketBuilder {
	return &TrieBucketBuilder{
		blockSize: blockSize,
		writer:    writer,
		builder:   newTrieBuilder(),
		sizeBuf:   make([]byte, 4),
	}
}

// Write writes key/value pairs into trie bucket.
func (b *TrieBucketBuilder) Write(keys [][]byte, ids []uint32) error {
	// NOTE: need sort keys for building trie
	kvs := &KVs{Keys: keys, IDs: ids}
	sort.Sort(kvs)
	// split keys based on block size
	numBlocks := len(keys) / b.blockSize
	if len(keys)%b.blockSize != 0 {
		numBlocks++
	}

	for i := 0; i < numBlocks; i++ {
		start := i * b.blockSize
		end := start + b.blockSize
		if end > len(keys) {
			// last block, check end pos
			end = len(keys)
		}
		// build trie
		b.builder.Reset()
		b.builder.Build(kvs.Keys[start:end], kvs.IDs[start:end])
		// write trie size
		size := b.builder.MarshalSize()
		binary.LittleEndian.PutUint32(b.sizeBuf[0:4], uint32(size))
		if _, err := b.writer.Write(b.sizeBuf); err != nil {
			return err
		}
		// write trie data
		if err := b.builder.Write(b.writer); err != nil {
			return err
		}
	}
	return nil
}

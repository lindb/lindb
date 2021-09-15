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

package replica

import (
	"bytes"
	"sync"

	"github.com/golang/snappy"

	"github.com/lindb/lindb/pkg/ltoml"
)

//go:generate mockgen -source=./chunk.go -destination=./chunk_mock.go -package=replica

// Chunk represents the writeTask buffer chunk for compressing the metric list
type Chunk interface {
	// Compress marshals and compresses the data, then resets the context,
	Compress() (*compressedChunk, error)
	// IsFull checks the chunk if is full
	IsFull() bool
	// IsEmpty checks the chunk if is empty
	IsEmpty() bool
	// Size returns the size of chunk
	Size() ltoml.Size
	// Write writes the metric into buffer
	Write([]byte) (n int, err error)
}

// chunk represents the buffer with snappy compress
type chunk struct {
	buffer   bytes.Buffer
	capacity ltoml.Size // use bytes capacity instead of lines-num
	size     ltoml.Size // chunk size and append index
}

// newChunk creates a new chunk
func newChunk(capacity ltoml.Size) Chunk {
	return &chunk{capacity: capacity}
}

// IsEmpty checks the chunk if is empty
func (c *chunk) IsEmpty() bool {
	return c.size == 0
}

// IsFull checks the chunk if is full
func (c *chunk) IsFull() bool {
	return c.size >= c.capacity
}

// Size returns the size of chunk
func (c *chunk) Size() ltoml.Size {
	return c.size
}

// Append appends the metric into buffer
func (c *chunk) Write(row []byte) (n int, err error) {
	n, err = c.buffer.Write(row)
	c.size += ltoml.Size(n)
	return n, err
}

// Compress marshals the data, then resets the context,
func (c *chunk) Compress() (*compressedChunk, error) {
	// if chunk is empty, return nil,nil
	if c.size == 0 {
		return nil, nil
	}

	defer func() {
		// if error, will ignore buffer data
		c.size = 0
		// reset for re-use
		c.buffer.Reset()
	}()

	// we use snappy block format here
	ck := newCompressedChunk(len(c.buffer.Bytes()))
	ck.Encode(c.buffer.Bytes())
	return ck, nil
}

var (
	compressedChunkPool sync.Pool
)

type compressedChunk []byte

// Release put the compressed chunk into sync pool
func (cc *compressedChunk) Release() {
	*cc = (*cc)[:0]
	compressedChunkPool.Put(cc)
}

// Encode compresses source block
func (cc *compressedChunk) Encode(block []byte) {
	*cc = snappy.Encode(*cc, block)
}

// newCompressedChunk picks a fixed sized buffer from pool
// expected compress ratio for snappy is 0.6 under of test
func newCompressedChunk(originalSize int) *compressedChunk {
	item := compressedChunkPool.Get()
	if item == nil {
		ck := make(compressedChunk, int(float64(originalSize)*0.6))
		return &ck
	}
	return item.(*compressedChunk)
}

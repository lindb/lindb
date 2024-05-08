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
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/pkg/compress"
)

//go:generate mockgen -source=./chunk.go -destination=./chunk_mock.go -package=replica

// Chunk represents the writeTask buffer chunk for compressing the metric list
type Chunk interface {
	// Compress marshals and compresses the data, then resets the context,
	Compress() ([]byte, error)
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
	writer   compress.Writer
	capacity ltoml.Size
	size     ltoml.Size
}

// newChunk creates a new chunk
func newChunk(capacity ltoml.Size) Chunk {
	c := &chunk{capacity: capacity}
	c.writer = compress.NewSnappyWriter()
	return c
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
	n, err = c.writer.Write(row)
	c.size += ltoml.Size(n)
	return n, err
}

// Compress marshals the data, then resets the context,
func (c *chunk) Compress() ([]byte, error) {
	defer func() {
		// if error, will ignore buffer data
		c.size = 0
	}()

	// if chunk is empty, return nil,nil
	if c.size == 0 {
		return nil, nil
	}

	if err := c.writer.Close(); err != nil {
		return nil, err
	}
	return c.writer.Bytes(), nil
}

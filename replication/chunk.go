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

package replication

import (
	"bytes"
	"encoding"

	"github.com/golang/snappy"

	"github.com/lindb/lindb/rpc/proto/field"
)

//go:generate mockgen -source=./chunk.go -destination=./chunk_mock.go -package=replication

// Chunk represents the write buffer chunk for compressing the metric list
type Chunk interface {
	// IsFull checks the chunk if is full
	IsFull() bool
	// IsEmpty checks the chunk if is empty
	IsEmpty() bool
	// Size returns the size of chunk
	Size() int
	// Append appends the metric into buffer
	Append(metric *field.Metric)
	// BinaryMarshaler marshals the data
	encoding.BinaryMarshaler
}

// chunk represents the buffer with snappy compress
type chunk struct {
	buf      *bytes.Buffer
	writer   *snappy.Writer
	buffer   field.MetricList
	capacity int
	size     int // chunk size and append index
}

// newChunk creates a new chunk
func newChunk(capacity int) Chunk {
	buf := &bytes.Buffer{}
	return &chunk{
		capacity: capacity,
		buf:      buf,
		buffer: field.MetricList{
			Metrics: make([]*field.Metric, capacity),
		},
		writer: snappy.NewBufferedWriter(buf),
	}
}

// IsEmpty checks the chunk if is empty
func (c *chunk) IsEmpty() bool {
	return c.size == 0
}

// IsFull checks the chunk if is full
func (c *chunk) IsFull() bool {
	return c.size == c.capacity
}

// Size returns the size of chunk
func (c *chunk) Size() int {
	return c.size
}

// Append appends the metric into buffer
func (c *chunk) Append(metric *field.Metric) {
	c.buffer.Metrics[c.size] = metric
	c.size++
}

// MarshalBinary marshals the data, then resets the context,
func (c *chunk) MarshalBinary() ([]byte, error) {
	// if chunk is empty, return nil,nil
	if c.size == 0 {
		return nil, nil
	}

	defer func() {
		// if error, will ignore buffer data
		c.size = 0
		// reset for re-use
		c.buffer.Metrics = make([]*field.Metric, c.capacity)
		c.buf = &bytes.Buffer{}
		c.writer.Reset(c.buf)
	}()

	// 1. if chunk not full, need truncate metric buffer list by the size
	if c.size < c.capacity {
		c.buffer.Metrics = c.buffer.Metrics[0:c.size]
	}
	// 2. marshal metric list
	data, err := c.buffer.Marshal()
	if err != nil {
		return nil, err
	}
	// 3. compress the data
	_, err = c.writer.Write(data)
	if err != nil {
		return nil, err
	}
	// 4. flush data
	if err := c.writer.Flush(); err != nil {
		return nil, err
	}
	// 5. return the binary data
	return c.buf.Bytes(), nil
}

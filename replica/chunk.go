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
	"encoding"
	"sync"

	"github.com/golang/snappy"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source=./chunk.go -destination=./chunk_mock.go -package=replica

// Chunk represents the writeTask buffer chunk for compressing the metric list
type Chunk interface {
	encoding.BinaryMarshaler

	// IsFull checks the chunk if is full
	IsFull() bool
	// IsEmpty checks the chunk if is empty
	IsEmpty() bool
	// Size returns the size of chunk
	Size() int
	// Append appends the metric into buffer
	Append(metric *protoMetricsV1.Metric)
}

// chunk represents the buffer with snappy compress
type chunk struct {
	buffer       *bytes.Buffer
	protoMetrics protoMetricsV1.MetricList
	capacity     int
	size         int // chunk size and append index
}

// newChunk creates a new chunk
func newChunk(capacity int) Chunk {
	return &chunk{
		capacity: capacity,
		buffer:   &bytes.Buffer{},
		protoMetrics: protoMetricsV1.MetricList{
			Metrics: make([]*protoMetricsV1.Metric, capacity),
		},
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
func (c *chunk) Append(metric *protoMetricsV1.Metric) {
	c.protoMetrics.Metrics[c.size] = metric
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
		c.protoMetrics.Metrics = make([]*protoMetricsV1.Metric, c.capacity)
		// TODO:  use flat metric
		c.buffer.Reset()
	}()

	// 1. if chunk not full, need truncate metric buffer list by the size
	if c.size < c.capacity {
		c.protoMetrics.Metrics = c.protoMetrics.Metrics[0:c.size]
	}

	// 2. marshal and compress metric list
	_, err := metric.MarshalProtoMetricsV1ListTo(c.protoMetrics, c.buffer)
	if err != nil {
		return nil, err
	}
	// we use snappy block format here
	var block = *getMarshalBlock()
	block = snappy.Encode(block, c.buffer.Bytes())
	return block, nil
}

var marshalBlockPool sync.Pool

func getMarshalBlock() *[]byte {
	item := marshalBlockPool.Get()
	if item == nil {
		var buf []byte
		return &buf
	}
	return item.(*[]byte)
}

func putMarshalBlock(b *[]byte) {
	*b = (*b)[:0]
	marshalBlockPool.Put(b)
}

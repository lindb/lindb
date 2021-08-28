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

package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lindb/lindb/pkg/stream"
)

// FixedOffsetEncoder represents the offset encoder with fixed length
// Make sure that added offset is increasing
type FixedOffsetEncoder struct {
	values           []int
	max              int
	ensureIncreasing bool
}

// NewFixedOffsetEncoder creates the fixed length offset encoder
// ensureIncreasing=true ensure that added offsets are increasing, panic when value is smaller than before
// ensureIncreasing=false suppresses the increasing check.
// Offset must >= 0
func NewFixedOffsetEncoder(ensureIncreasing bool) *FixedOffsetEncoder {
	return &FixedOffsetEncoder{ensureIncreasing: ensureIncreasing}
}

// IsEmpty returns if is empty
func (e *FixedOffsetEncoder) IsEmpty() bool {
	return len(e.values) == 0
}

// Size returns the size
func (e *FixedOffsetEncoder) Size() int {
	return len(e.values)
}

// Reset resets the encoder context for reuse
func (e *FixedOffsetEncoder) Reset() {
	e.max = 0
	e.values = e.values[:0]
}

// Add adds the start offset value,
func (e *FixedOffsetEncoder) Add(v int) {
	if e.ensureIncreasing && len(e.values) > 0 && e.values[len(e.values)-1] > v {
		panic("value added to FixedOffsetEncoder must be increasing")
	}
	if v < 0 {
		panic("value add be FixedOffsetEncoder must > 0")
	}
	e.values = append(e.values, v)
	if e.max < v {
		e.max = v
	}
}

// FromValues resets the encoder, then init it with multi values.
func (e *FixedOffsetEncoder) FromValues(values []int) {
	e.Reset()
	e.values = values
	for _, value := range values {
		if e.max < value {
			e.max = value
		}
	}
}

// MarshalBinary marshals the values to binary
func (e *FixedOffsetEncoder) MarshalBinary() []byte {
	var buf bytes.Buffer
	buf.Grow(e.MarshalSize())
	_ = e.Write(&buf)
	return buf.Bytes()
}

func (e *FixedOffsetEncoder) MarshalSize() int {
	return 1 + // width flag
		stream.UvariantSize(uint64(len(e.values))) + // size
		len(e.values)*e.width() // values
}

func (e *FixedOffsetEncoder) width() int {
	return Uint32MinWidth(uint32(e.max))
}

// Write writes the data to the writer.
func (e *FixedOffsetEncoder) Write(writer io.Writer) error {
	if len(e.values) == 0 {
		return nil
	}
	width := e.width()
	// fixed value width
	if _, err := writer.Write([]byte{uint8(width)}); err != nil {
		return err
	}
	// put all values with fixed length
	var buf [binary.MaxVarintLen64]byte
	// write size
	sizeFlagWidth := binary.PutUvarint(buf[:], uint64(len(e.values)))
	if _, err := writer.Write(buf[:sizeFlagWidth]); err != nil {
		return err
	}
	// write values
	for _, value := range e.values {
		binary.LittleEndian.PutUint32(buf[:], uint32(value))
		if _, err := writer.Write(buf[:width]); err != nil {
			return err
		}
	}
	return nil
}

// FixedOffsetDecoder represents the fixed offset decoder,
// supports random reads offset by index
type FixedOffsetDecoder struct {
	offsetsBlock []byte
	width        int
	size         int
}

// NewFixedOffsetDecoder creates the fixed offset decoder
func NewFixedOffsetDecoder() *FixedOffsetDecoder {
	return &FixedOffsetDecoder{}
}

// ValueWidth returns the width of all stored values
func (d *FixedOffsetDecoder) ValueWidth() int {
	return d.width
}

// Size returns the size of  offset values
func (d *FixedOffsetDecoder) Size() int {
	if d.width == 0 {
		return 0
	}
	return d.size
}

// Unmarshal unmarshals from data block, then return the remaining buffer.
func (d *FixedOffsetDecoder) Unmarshal(data []byte) (left []byte, err error) {
	d.offsetsBlock = d.offsetsBlock[:0]
	d.width = 0
	d.size = 0
	if len(data) < 2 {
		return nil, fmt.Errorf("length too short of FixedOffsetDecoder: %d", len(data))
	}
	d.width = int(data[0])
	if d.width < 0 || d.width > 4 {
		return nil, fmt.Errorf("ivalid width of FixedOffsetDecoder: %d", d.width)
	}
	size, readBytes := binary.Uvarint(data[1:])
	if readBytes <= 0 {
		return nil, fmt.Errorf("invalid uvariant of FixedOffsetDecoder")
	}
	d.size = int(size)
	wantLen := 1 + readBytes + d.width*d.size
	if wantLen > len(data) || wantLen < 0 || 1+readBytes > wantLen {
		return nil, fmt.Errorf("cannot unmarshal FixedOffsetDecoder with a invalid buffer: %d, want: %d",
			len(data), wantLen)
	}
	d.offsetsBlock = data[1+readBytes : wantLen]
	return data[wantLen:], nil
}

func (d *FixedOffsetDecoder) Get(index int) (int, bool) {
	start := index * d.width
	if start < 0 || len(d.offsetsBlock) == 0 || start >= len(d.offsetsBlock) || d.width > 4 {
		return 0, false
	}
	end := start + d.width
	if end > len(d.offsetsBlock) {
		return 0, false
	}
	var scratch [4]byte
	copy(scratch[:], d.offsetsBlock[start:end])
	offset := int(binary.LittleEndian.Uint32(scratch[:]))
	// on x32, data may overflow
	if offset < 0 {
		return 0, false
	}
	return offset, true
}

// GetBlock returns the block by offset range(start -> end) with index
// GetBlock is only supported when Offsets are increasing encoded.
func (d *FixedOffsetDecoder) GetBlock(index int, dataBlock []byte) (block []byte, err error) {
	startOffset, ok := d.Get(index)
	if !ok {
		return nil, fmt.Errorf("corrupted FixedOffsetDecoder block, length: %d, startOffset: %d",
			len(d.offsetsBlock), startOffset)
	}
	endOffset, ok := d.Get(index + 1)
	if !ok {
		endOffset = len(dataBlock)
	}

	if startOffset < 0 || endOffset < 0 || endOffset < startOffset || endOffset > len(dataBlock) {
		return nil, fmt.Errorf("corrupted FixedOffsetDecoder block, "+
			"data block length: %d, data range: [%d, %d]", len(dataBlock), startOffset, endOffset,
		)
	}
	return dataBlock[startOffset:endOffset], nil
}

func ByteSlice2Uint32(slice []byte) uint32 {
	var buf = make([]byte, 4)
	copy(buf, slice)
	return binary.LittleEndian.Uint32(buf)
}

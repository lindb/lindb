package encoding

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/lindb/lindb/pkg/stream"
)

const (
	adjustValue    = 1
	lengthOfHeader = 1
	EmptyOffset    = -1
)

// FixedOffsetEncoder represents the offset encoder with fixed length
type FixedOffsetEncoder struct {
	values []int
	buf    *bytes.Buffer
	max    int
	bw     *stream.BufferWriter
}

// NewFixedOffsetEncoder creates the fixed length offset encoder
func NewFixedOffsetEncoder() *FixedOffsetEncoder {
	var buf bytes.Buffer
	bw := stream.NewBufferWriter(&buf)
	return &FixedOffsetEncoder{
		buf: &buf,
		bw:  bw,
		max: EmptyOffset,
	}
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
	e.bw.Reset()
	e.max = 0
	e.values = e.values[:0]
}

// Add adds the offset value,
func (e *FixedOffsetEncoder) Add(v int) {
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
	_ = e.WriteTo(e.buf)
	return e.buf.Bytes()
}

// WriteTo writes the data to the writer.
func (e *FixedOffsetEncoder) WriteTo(writer io.Writer) error {
	if len(e.values) == 0 {
		return nil
	}
	if e.max < 0 {
		e.max = EmptyOffset
	}
	width := Uint32MinWidth(uint32(e.max + adjustValue))
	// fixed value width
	e.bw.PutByte(byte(width))
	// put all values with fixed length
	buf := make([]byte, 4)
	for _, value := range e.values {
		if value < 0 {
			value = EmptyOffset
		}
		binary.LittleEndian.PutUint32(buf, uint32(value+adjustValue))
		if _, err := writer.Write(buf[:width]); err != nil {
			return err
		}
	}
	return nil
}

// FixedOffsetDecoder represents the fixed offset decoder, supports random reads offset by index
type FixedOffsetDecoder struct {
	buf     []byte
	width   int
	scratch []byte
}

// NewFixedOffsetDecoder creates the fixed offset decoder
func NewFixedOffsetDecoder(buf []byte) *FixedOffsetDecoder {
	if len(buf) == 0 {
		return &FixedOffsetDecoder{
			buf: nil,
		}
	}
	return &FixedOffsetDecoder{
		buf:     buf[lengthOfHeader:],
		width:   int(buf[0]),
		scratch: make([]byte, 4),
	}
}

// Header returns the length of header
func (d *FixedOffsetDecoder) Header() int {
	return lengthOfHeader
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
	return len(d.buf) / d.width
}

// Get gets the offset value by index
func (d *FixedOffsetDecoder) Get(index int) (int, bool) {
	start := index * d.width
	if start < 0 || len(d.buf) == 0 || start >= len(d.buf) || d.width > 4 {
		return 0, false
	}
	end := start + d.width
	if end > len(d.buf) {
		return 0, false
	}
	copy(d.scratch, d.buf[start:end])
	offset := int(binary.LittleEndian.Uint32(d.scratch)) - adjustValue
	if offset < 0 {
		return 0, false
	}
	return offset, true
}

func ByteSlice2Uint32(slice []byte) uint32 {
	var buf = make([]byte, 4)
	copy(buf, slice)
	return binary.LittleEndian.Uint32(buf)
}

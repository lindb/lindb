package encoding

import (
	"bytes"

	"github.com/lindb/lindb/pkg/stream"
)

// FixedOffsetEncoder represents the offset encoder with fixed length
type FixedOffsetEncoder struct {
	values []int
	buf    *bytes.Buffer
	bw     *stream.BufferWriter
}

// NewFixedOffsetEncoder creates the fixed length offset encoder
func NewFixedOffsetEncoder() *FixedOffsetEncoder {
	var buf bytes.Buffer
	bw := stream.NewBufferWriter(&buf)
	return &FixedOffsetEncoder{
		buf: &buf,
		bw:  bw,
	}
}

// Reset resets the encoder context for reuse
func (e *FixedOffsetEncoder) Reset() {
	e.bw.Reset()
	e.values = e.values[:0]
}

// Add adds the offset value,
// NOTICE: value need keep in sort
func (e *FixedOffsetEncoder) Add(v int) {
	e.values = append(e.values, v)
}

// MarshalBinary marshals the values to binary
func (e *FixedOffsetEncoder) MarshalBinary() []byte {
	length := len(e.values)
	if length == 0 {
		return nil
	}
	maxLength := getValueLen(e.values[length-1])
	e.bw.PutByte(byte(maxLength)) // max fixed length
	// put all values with fixed length
	buf := make([]byte, maxLength)
	for i := 0; i < length; i++ {
		putInt32(buf, e.values[i], maxLength)
		e.bw.PutBytes(buf)
	}
	return e.buf.Bytes()
}

// FixedOffsetDecoder represents the fixed offset decoder, supports random reads offset by index
type FixedOffsetDecoder struct {
	buf                 []byte
	valueLength, length int
}

// NewFixedOffsetDecoder creates the fixed offset decoder
func NewFixedOffsetDecoder(buf []byte) *FixedOffsetDecoder {
	return &FixedOffsetDecoder{
		buf:         buf,
		valueLength: int(buf[0]),
		length:      len(buf),
	}
}

// Get gets the offset value by index, if offset > buffer length or index <0 returns -1
func (d *FixedOffsetDecoder) Get(index int) int {
	if index < 0 {
		return -1
	}
	// offset = index * length + 1 (1 is max length)
	offset := index*d.valueLength + 1
	if offset+d.valueLength > d.length {
		return -1
	}
	return getInt(d.buf, index*d.valueLength+1, d.valueLength)
}

// getInt32 gets value from buf with fixed length
func getInt(buf []byte, offset, length int) int {
	var x uint32
	switch {
	case length == 1:
		x = uint32(buf[offset])
	case length == 2:
		x = uint32(buf[offset])
		x |= uint32(buf[offset+1]) << 8
	case length == 3:
		x = uint32(buf[offset])
		x |= uint32(buf[offset+1]) << 8
		x |= uint32(buf[offset+2]) << 16
	case length == 4:
		x = uint32(buf[offset])
		x |= uint32(buf[offset+1]) << 8
		x |= uint32(buf[offset+2]) << 16
		x |= uint32(buf[offset+3]) << 24
	}
	return int(x)
}

// putInt32 puts the value into buf with fixed length
func putInt32(buf []byte, value int, length int) {
	x := uint32(value)
	switch {
	case length == 1:
		buf[0] = uint8(x & 0xff)
	case length == 2:
		buf[0] = uint8(x)
		buf[1] = uint8(x >> 8)
	case length == 3:
		buf[0] = uint8(x)
		buf[1] = uint8(x >> 8)
		buf[2] = uint8(x >> 16)
	case length == 4:
		buf[0] = uint8(x)
		buf[1] = uint8(x >> 8)
		buf[2] = uint8(x >> 16)
		buf[3] = uint8(x >> 24)
	}
}

// getValueLen returns the value of store min length
func getValueLen(value int) int {
	switch {
	case value < 1<<8:
		return 1
	case value < 1<<16:
		return 2
	case value < 1<<24:
		return 3
	default:
		return 4
	}
}

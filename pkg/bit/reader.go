package bit

import (
	"bytes"
	"io"
)

// NewReader crate bit reader
func NewReader(data []byte) *Reader {
	return &Reader{
		reader: bytes.NewReader(data),
		count:  0}
}

// Reader reads bits from buffer
type Reader struct {
	reader *bytes.Reader
	b      byte
	count  uint8

	err error
}

// ReadBit reads a bit, if failure return error
func (r *Reader) ReadBit() (Bit, error) {
	if r.count == 0 {
		r.b, r.err = r.reader.ReadByte()
		r.count = 8
	}
	r.count--
	d := r.b & 0x80
	r.b <<= 1
	return d != 0, r.err
}

// ReadBits read number of bits
func (r *Reader) ReadBits(numBits int) (uint64, error) {
	var u uint64

	for numBits >= 8 {
		byt, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		numBits -= 8
	}

	var err error
	for numBits > 0 && err != io.EOF {
		byt, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		numBits--
	}

	return u, nil
}

// ReadByte reads a byte
func (r *Reader) ReadByte() (byte, error) {
	if r.count == 0 {
		r.b, r.err = r.reader.ReadByte()
		return r.b, r.err
	}

	byt := r.b

	r.b, r.err = r.reader.ReadByte()
	byt |= r.b >> r.count
	r.b <<= 8 - r.count
	return byt, r.err
}

// Reset resets the reader to read from a new slice
func (r *Reader) Reset(data []byte) {
	r.reader.Reset(data)
	r.err = nil
	r.count = 0
	r.b = 0
}

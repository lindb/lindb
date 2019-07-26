package bit

import (
	"bytes"
	"io"
)

// A Bit is a zero or a one
type Bit bool

const (
	// Zero is our exported type for '0' bits
	Zero Bit = false
	// One is our exported type for '1' bits
	One Bit = true
)

// Writer writes bits to an io.Writer
type Writer struct {
	w     io.Writer
	b     [1]byte
	count uint8
}

// Reader reads bits from buffer
type Reader struct {
	buf   *bytes.Buffer
	b     byte
	count uint8

	err error
}

// NewWriter create bit writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     w,
		count: 8,
	}
}

// NewReader crate bit reader
func NewReader(buf *bytes.Buffer) *Reader {
	return &Reader{
		buf:   buf,
		count: 0,
	}
}

// WriteBit writes a bit value
func (w *Writer) WriteBit(bit Bit) error {
	if bit {
		w.b[0] |= 1 << (w.count - 1)
	}

	w.count--

	if w.count == 0 {
		// fill byte to io.Writer
		if n, err := w.w.Write(w.b[:]); n != 1 || err != nil {
			return err
		}
		w.b[0] = 0
		w.count = 8
	}

	return nil
}

// WriteBits writes number of bits
func (w *Writer) WriteBits(u uint64, numBits int) error {
	u <<= 64 - uint(numBits)

	for numBits >= 8 {
		byt := byte(u >> 56)
		err := w.WriteByte(byt)
		if err != nil {
			return err
		}
		u <<= 8
		numBits -= 8
	}

	for numBits > 0 {
		err := w.WriteBit((u >> 63) == 1)
		if err != nil {
			return err
		}
		u <<= 1
		numBits--
	}

	return nil
}

// WriteByte write a byte
func (w *Writer) WriteByte(b byte) error {
	w.b[0] |= b >> (8 - w.count)

	if n, err := w.w.Write(w.b[:]); n != 1 || err != nil {
		return err
	}

	w.b[0] = b << w.count

	return nil
}

//Flush flushes the currently in-process byte
func (w *Writer) Flush() error {
	if w.count != 8 {
		_, err := w.w.Write(w.b[:])
		return err
	}
	return nil
}

// ReadBit reads a bit, if failure return error
func (r *Reader) ReadBit() (Bit, error) {
	if r.count == 0 {
		r.b, r.err = r.buf.ReadByte()
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
		r.b, r.err = r.buf.ReadByte()
		return r.b, r.err
	}

	byt := r.b

	r.b, r.err = r.buf.ReadByte()
	byt |= r.b >> r.count
	r.b <<= 8 - r.count
	return byt, r.err
}

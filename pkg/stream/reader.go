package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ErrUnexpectedRead is raised when reading negative length
var ErrUnexpectedRead = fmt.Errorf("unexpected read")

// Reader is a stream reader
type Reader struct {
	original []byte        // the original data block
	reader   *bytes.Reader // Reader of sub-slice
	err      error
}

// NewReader read data from binary stream
func NewReader(data []byte) *Reader {
	return &Reader{
		original: data,
		reader:   bytes.NewReader(data)}
}

// ReadVarint32 reads int32 from buffer
func (r *Reader) ReadVarint32() int32 {
	return int32(r.ReadVarint64())
}

// ReadVarint64 reads int64 from buffer
func (r *Reader) ReadVarint64() int64 {
	var v int64
	v, r.err = readVarint(r.reader)
	return v
}

// ReadUvarint32 reads uint32 from buffer
func (r *Reader) ReadUvarint32() uint32 {
	return uint32(r.ReadUvarint64())
}

// ReadUvarint64 reads uint64 from buffer
func (r *Reader) ReadUvarint64() uint64 {
	var v uint64
	v, r.err = readUvarint(r.reader)
	return v
}

// ReadUint16 read 2 bytes from buf as uint16
func (r *Reader) ReadUint16() uint16 {
	buf := r.ReadSlice(2)
	if len(buf) != 2 {
		return 0
	}
	return binary.LittleEndian.Uint16(buf)
}

// ReadInt16 read 2 bytes from buf as int16
func (r *Reader) ReadInt16() int16 {
	return int16(r.ReadUint16())
}

// ReadUint32 read 4 bytes from buf as uint32
func (r *Reader) ReadUint32() uint32 {
	buf := r.ReadSlice(4)
	if len(buf) != 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(buf)
}

// ReadUint64 read 8 bytes from buf as uint64
func (r *Reader) ReadUint64() uint64 {
	buf := r.ReadSlice(8)
	if len(buf) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(buf)
}

// ReadInt32 read 4 bytes from buf as int32
func (r *Reader) ReadInt32() int32 {
	return int32(r.ReadUint32())
}

// ReadInt64 read 8 bytes from buf as int64
func (r *Reader) ReadInt64() int64 {
	return int64(r.ReadUint64())
}

// ReadByte reads 1 byte
func (r *Reader) ReadByte() byte {
	var b byte
	b, r.err = r.reader.ReadByte()
	return b
}

// ReadBytes reads n len bytes
func (r *Reader) ReadBytes(n int) []byte {
	if n < 0 {
		r.err = ErrUnexpectedRead
		return nil
	}
	block := make([]byte, n)
	for i := 0; i < n; i++ {
		block[i], r.err = r.reader.ReadByte()
		if r.err != nil {
			return block[:i]
		}
	}
	return block
}

// ReadSlice returns a sub-slice.
// make sure that the sub-slice is not in use before calling Reset.
func (r *Reader) ReadSlice(n int) []byte {
	if n < 0 {
		r.err = ErrUnexpectedRead
		return nil
	}
	if r.err != nil {
		return nil
	}
	startPos, endPos := r.Position(), r.Position()+n
	if endPos > len(r.original) {
		endPos = len(r.original)
		r.err = io.EOF
	}
	r.reader.Reset(r.original[endPos:])
	return r.original[startPos:endPos]
}

// Empty reports whether the unread portion of the buffer is empty.
func (r *Reader) Empty() bool {
	return r.reader.Len() <= 0
}

// Position returns the position where reader at
func (r *Reader) Position() int {
	return len(r.original) - r.reader.Len()
}

// Reset resets the Reader, then reads from the buffer
func (r *Reader) Reset(buf []byte) {
	r.original = buf
	r.reader.Reset(buf)
	r.err = nil
}

// SeekStart seeks to the start of the underlying slice.
func (r *Reader) SeekStart() {
	r.err = nil
	r.reader.Reset(r.original)
}

// Error return binary err
func (r *Reader) Error() error {
	return r.err
}

var errOverflow = errors.New("varint overflows a 64-bit integer")

// copy from binary
// readUvarint reads an encoded unsigned integer from bytes.Reader and returns it as a uint64.
func readUvarint(r *bytes.Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, errOverflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

// copy from binary
// readVarint reads an encoded signed integer from bytes.Reader and returns it as an int64.
func readVarint(r *bytes.Reader) (int64, error) {
	ux, err := readUvarint(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
}

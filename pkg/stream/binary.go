package stream

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Binary is stream for writing data into memory buffer
type Binary struct {
	buf *bytes.Buffer

	scratch [binary.MaxVarintLen64]byte
	err     error
}

// BinaryWriter create binary stream for writing data
func BinaryWriter() *Binary {
	var v []byte
	return &Binary{buf: bytes.NewBuffer(v)}
}

// BinaryBufWriter creates a binary stream for writing with provided buffer(writing start from offset 0).
// The caller is responsible to make sure the wrote bytes do not exceed the buffer size,
// otherwise the exceeding data will be missing.
func BinaryBufWriter(buffer []byte) *Binary {
	return &Binary{buf: bytes.NewBuffer(buffer[:0])}
}

// BinaryReader create binary stream for reading data
func BinaryReader(v []byte) *Binary {
	return &Binary{
		buf: bytes.NewBuffer(v),
	}
}

// PutBytes encodes bytes into buf
func (b *Binary) PutBytes(v []byte) {
	length := len(v)
	if length > 0 {
		n, err := b.buf.Write(v)
		b.err = err
		if n != length {
			b.err = fmt.Errorf("write len not eqauls value's len")
		}
	}
}

// PutByte encodes a byte into buf
func (b *Binary) PutByte(v byte) {
	b.buf.WriteByte(v)
}

// PutLenBytes encodes data as length(Uvarint64) +  data
func (b *Binary) PutLenBytes(data []byte) {
	b.PutUvarint64(uint64(len(data)))
	b.PutBytes(data)
}

// PutVarint32 encodes a int32 into buf
func (b *Binary) PutVarint32(v int32) {
	b.PutUvarint64(uint64(v))
}

// PutVarint64 encodes a int64 into buf
func (b *Binary) PutVarint64(v int64) {
	b.PutUvarint64(uint64(v))
}

// PutUvarint32 encodes a uint32 into buf
func (b *Binary) PutUvarint32(v uint32) {
	b.PutUvarint64(uint64(v))
}

// PutUvarint64 encodes a uint64 into buf
func (b *Binary) PutUvarint64(v uint64) {
	if b.err != nil {
		return
	}
	n := binary.PutUvarint(b.scratch[:], v)
	b.PutBytes(b.scratch[:n])
}

// PutUint32 encodes a uint32 as 4 bytes into buf
func (b *Binary) PutUint32(v uint32) {
	binary.LittleEndian.PutUint32(b.scratch[:], v)
	b.PutBytes(b.scratch[:4])
}

// PutUint64 encodes a int64 as 8 bytes into buf
func (b *Binary) PutUint64(v uint64) {
	binary.LittleEndian.PutUint64(b.scratch[:], v)
	b.PutBytes(b.scratch[:8])
}

// PutInt32 encodes a int32 as 4 bytes into buf
func (b *Binary) PutInt32(v int32) {
	b.PutUint32(uint32(v))
}

// PutInt64 encodes a int64 as 8 bytes into buf
func (b *Binary) PutInt64(v int64) {
	b.PutUint64(uint64(v))
}

// Bytes returns memory buffer data, if error return err
func (b *Binary) Bytes() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.buf.Bytes(), nil
}

//Len returns the number of bytes of the unread portion of the buffer;
func (b *Binary) Len() int {
	return b.buf.Len()
}

// ReadVarint32 reads int32 from buffer
func (b *Binary) ReadVarint32() int32 {
	return int32(b.ReadUvarint64())
}

// ReadVarint64 reads int64 from buffer
func (b *Binary) ReadVarint64() int64 {
	return int64(b.ReadUvarint64())
}

// ReadUvarint32 reads uint32 from buffer
func (b *Binary) ReadUvarint32() uint32 {
	return uint32(b.ReadUvarint64())
}

// ReadUvarint64 reads uint64 from buffer
func (b *Binary) ReadUvarint64() uint64 {
	v, err := binary.ReadUvarint(b.buf)
	if err != nil {
		b.err = err
	}
	return v
}

// ReadUint32 read 4 bytes from buf as uint32
func (b *Binary) ReadUint32() uint32 {
	buf := b.ReadBytes(4)
	return binary.LittleEndian.Uint32(buf)
}

// ReadUint64 read 8 bytes from buf as uint64
func (b *Binary) ReadUint64() uint64 {
	buf := b.ReadBytes(8)
	return binary.LittleEndian.Uint64(buf)
}

// ReadInt32 read 4 bytes from buf as int32
func (b *Binary) ReadInt32() int32 {
	return int32(b.ReadUint32())
}

// ReadInt64 read 8 bytes from buf as int64
func (b *Binary) ReadInt64() int64 {
	return int64(b.ReadUint64())
}

// ReadBytes reads n len bytes, use buf.Next()
func (b *Binary) ReadBytes(n int) []byte {
	return b.buf.Next(n)
}

// Empty reports whether the unread portion of the buffer is empty.
func (b *Binary) Empty() bool {
	return b.buf.Len() <= 0
}

// Error return binary err
func (b *Binary) Error() error {
	return b.err
}

//TODO need refactor??
//ByteBufReader provides methods to read specific values from a byte array
type ByteBufReader struct {
	buf      []byte
	position int //current reader position
	length   int
}

//NewBufReader create binary stream for reading data
func NewBufReader(bufArray []byte) *ByteBufReader {
	return &ByteBufReader{
		buf:    bufArray,
		length: len(bufArray),
	}
}

// ReadLenBytes reads a key containing length and bytes from the buffer
func (r *ByteBufReader) ReadLenBytes() (length int, key []byte) {
	length = int(r.ReadUvarint64())
	if length > 0 {
		key = r.ReadBytes(length)
	}
	return length, key
}

// ReadUvarint64 reads variable-length positive int from the buffer
func (r *ByteBufReader) ReadUvarint64() uint64 {
	v, length := readUvarint(r.buf[r.position:])
	r.position += length
	return v
}

// readUvarint reads variable-length positive int from the buffer
func readUvarint(buf []byte) (v uint64, len int) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		b := buf[i]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, i + 1
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

//ReadUint32 reads uint32 from the buffer
func (r *ByteBufReader) ReadUint32() uint32 {
	b := r.ReadBytes(4)
	return binary.LittleEndian.Uint32(b)
}

//ReadBytes reads fixed-length bytes from the buffer
func (r *ByteBufReader) ReadBytes(length int) []byte {
	if length == 0 {
		return nil
	}
	b := r.buf[r.position : r.position+length]
	r.position += length
	return b
}

//SubArray means to get the byte array from start
func (r *ByteBufReader) SubArray(start int) []byte {
	sub := r.buf[start:]
	return sub
}

//ReadByte reads a byte from the buffer
func (r *ByteBufReader) ReadByte() byte {
	b := r.buf[r.position]
	r.position++
	return b
}

// NewPosition indicates the reset position
func (r *ByteBufReader) NewPosition(newPos int) {
	r.position = newPos
}

//IsEnd indicates whether the end bit is read
func (r *ByteBufReader) IsEnd() bool {
	return r.position == r.length
}

//GetPosition indicates the current read position
func (r *ByteBufReader) GetPosition() int {
	return r.position
}

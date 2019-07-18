package stream

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/eleme/lindb/pkg/util"
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

// BinaryReader create binary stream for reading data
func BinaryReader(v []byte) *Binary {
	return &Binary{
		buf: bytes.NewBuffer(v),
	}
}

//=================================binary writer
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

// PutKey encodes a key into buf
func (b *Binary) PutKey(k []byte) {
	b.PutUvarint64(uint64(len(k)))
	b.PutBytes(k)
}

// PutInt32 encodes a int32 into buf
func (b *Binary) PutInt32(v int32) {
	b.PutUvarint64(uint64(v))
}

//PutUInt32 encodes a uint32 into buf
func (b *Binary) PutUInt32(v uint32) {
	by := util.Uint32ToBytes(v)
	b.PutBytes(by)
}

// PutInt64 encodes a int64 into buf
func (b *Binary) PutInt64(v int64) {
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

//=================================binary reader
// ReadInt32 reads int32 from buffer
func (b *Binary) ReadInt32() int32 {
	return int32(b.ReadUvarint64())
}

// ReadInt64 reads int64 from buffer
func (b *Binary) ReadInt64() int64 {
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

//=================================buf reader
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

//ReadKey reads a key containing length and bytes from the buffer
func (r *ByteBufReader) ReadKey() (length int, key []byte) {
	length = int(r.ReadInt())
	if length > 0 {
		key = r.ReadBytes(length)
	}
	return length, key
}

//ReadInt reads variable-length positive int from the buffer
func (r *ByteBufReader) ReadInt() uint64 {
	v, length := readVInt(r.buf[r.position:])
	r.position += length
	return v
}

//ReadInt reads variable-length positive int from the buffer
func readVInt(buf []byte) (v uint64, len int) {
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

//ReadUInt32 reads uint32 from the buffer
func (r *ByteBufReader) ReadUInt32() uint32 {
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

//ReadBytes reads a byte from the buffer
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

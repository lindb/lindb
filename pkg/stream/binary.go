package stream

import (
	"bytes"
	"encoding/binary"
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
		_, b.err = b.buf.Write(v)
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

// ReadByte reads 1 byte
func (b *Binary) ReadByte() byte {
	data := b.buf.Next(1)
	if len(data) == 1 {
		return data[0]
	}
	return byte(0)
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

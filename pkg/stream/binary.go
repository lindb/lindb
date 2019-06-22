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

// BinaryReader create binary stream for reading data
func BinaryReader(v []byte) *Binary {
	return &Binary{
		buf: bytes.NewBuffer(v),
	}
}

// PutBytes encodes bytes into buf
func (b *Binary) PutBytes(v []byte) {
	len := len(v)

	n, err := b.buf.Write(v)
	b.err = err
	if n != len {
		b.err = fmt.Errorf("write len not eqauls value's len")
	}
}

// PutInt32 encodes a int32 into buf
func (b *Binary) PutInt32(v int32) {
	b.PutUvarint64(uint64(v))
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

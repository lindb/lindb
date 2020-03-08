package stream

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

////////////////////////////////////////////////////////
//      Base Writer
////////////////////////////////////////////////////////

// writer is base writer for writing data
type writer struct {
	buf *bytes.Buffer

	scratch [binary.MaxVarintLen64]byte
	err     error
}

// PutBytes encodes bytes into buf
func (w *writer) PutBytes(v []byte) {
	_, w.err = w.buf.Write(v)
}

// PutByte encodes a byte into buf
func (w *writer) PutByte(v byte) {
	w.err = w.buf.WriteByte(v)
}

// Write implements io.Writer
func (w *writer) Write(p []byte) (n int, err error) {
	n, w.err = w.buf.Write(p)
	return n, w.err
}

// PutVarint32 encodes a int32 into buf
func (w *writer) PutVarint32(v int32) {
	w.PutVarint64(int64(v))
}

// PutVarint64 encodes a int64 into buf
func (w *writer) PutVarint64(v int64) {
	n := binary.PutVarint(w.scratch[:], v)
	_, w.err = w.buf.Write(w.scratch[:n])
}

// PutUvarint32 encodes a uint32 into buf
func (w *writer) PutUvarint32(v uint32) {
	w.PutUvarint64(uint64(v))
}

// PutUvarint64 encodes a uint64 into buf
func (w *writer) PutUvarint64(v uint64) {
	n := binary.PutUvarint(w.scratch[:], v)
	_, w.err = w.buf.Write(w.scratch[:n])
}

// PutUint32 encodes a uint32 as 4 bytes into buf
func (w *writer) PutUint32(v uint32) {
	binary.LittleEndian.PutUint32(w.scratch[:], v)
	_, w.err = w.buf.Write(w.scratch[:4])
}

// PutUint64 encodes a uint64 as 8 bytes into buf
func (w *writer) PutUint64(v uint64) {
	binary.LittleEndian.PutUint64(w.scratch[:], v)
	_, w.err = w.buf.Write(w.scratch[:8])
}

// PutInt32 encodes a int32 as 4 bytes into buf
func (w *writer) PutInt32(v int32) {
	w.PutUint32(uint32(v))
}

// PutInt64 encodes a int64 as 8 bytes into buf
func (w *writer) PutInt64(v int64) {
	w.PutUint64(uint64(v))
}

// PutUInt16 encodes a uint16 as 2 bytes into buf
func (w *writer) PutUInt16(v uint16) {
	binary.LittleEndian.PutUint16(w.scratch[:], v)
	_, w.err = w.buf.Write(w.scratch[:2])
}

// PutInt16 encodes a int16 as 2 bytes into buf
func (w *writer) PutInt16(v int16) {
	w.PutUInt16(uint16(v))
}

// Len returns the size of bytes of the written data of the buffer;
func (w *writer) Len() int {
	return w.buf.Len()
}

// BufferWriter is a writer for writing data into a buffer
type BufferWriter struct {
	writer
}

////////////////////////////////////////////////////////
//      Buffer Writer
////////////////////////////////////////////////////////

// NewBufferWriter creates a binary stream for writing with provided buffer(append write).
func NewBufferWriter(buffer *bytes.Buffer) *BufferWriter {
	if buffer == nil {
		return &BufferWriter{writer{buf: &bytes.Buffer{}}}
	}
	return &BufferWriter{writer{buf: buffer}}
}

// Reset resets the underling buffer
func (bw *BufferWriter) Reset() {
	bw.writer.err = nil
	bw.buf.Reset()
}

// SwitchBuffer switches to write a new buffer
func (bw *BufferWriter) SwitchBuffer(newBuffer *bytes.Buffer) {
	bw.writer.err = nil
	bw.buf = newBuffer
}

// Error returns the error of BufferWriter
func (bw *BufferWriter) Error() error {
	return bw.err
}

// Bytes returns memory buffer data, if error return err
func (bw *BufferWriter) Bytes() ([]byte, error) {
	return bw.buf.Bytes(), bw.Error()
}

////////////////////////////////////////////////////////
//      Slice Writer
////////////////////////////////////////////////////////

// SliceWriter is a writer for writing data into a slice
type SliceWriter struct {
	writer
	maxLen int
}

// NewSliceWriter creates a binary stream for writing with provided slice(writing start from offset 0).
// The caller is responsible to make sure the wrote bytes do not exceed the buffer size,
// otherwise the exceeding data will be missing.
// and if it is a large buf, buffer grow will cause a expensive reallocation
func NewSliceWriter(buffer []byte) *SliceWriter {
	length := len(buffer)
	return &SliceWriter{
		writer: writer{buf: bytes.NewBuffer(buffer[:0])},
		maxLen: length}
}

// Error returns the error of BufferWriter
func (sw *SliceWriter) Error() error {
	if len(sw.buf.Bytes()) > sw.maxLen {
		return fmt.Errorf("write longer than fixed size, %d > %d", len(sw.buf.Bytes()), sw.maxLen)
	}
	return sw.err
}

// Bytes returns memory buffer data, if error return err
func (sw *SliceWriter) Bytes() ([]byte, error) {
	return sw.buf.Bytes(), sw.Error()
}

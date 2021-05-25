// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bufioutil

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

//go:generate mockgen -source=./bufio_writer.go -destination=./bufio_writer_mock.go -package=bufioutil

const (
	defaultWriteBufferSize = 256 * 1024 // 256KB
)

// BufioWriter writes entries to a specified file by buffered I/O. Not thread-safe.
type BufioWriter interface {
	// Write writes a new entry containing logs in order.
	// Close syncs data to disk, then closes the opened file.
	io.WriteCloser
	// Reset switches the buffered writer to write to a new file:
	// open the new file; close the old opening file;
	// discards any buffered data and reset the states of bufio.Writer
	Reset(fileName string) error
	// Sync flushes data first, then calls syscall.sync.
	Sync() error
	// Flush flushes data to io.Writer.
	Flush() error
	// Size returns the length of written data.
	Size() int64
}

// bufioWriter implements BufioWriter.
type bufioWriter struct {
	fileName string
	w        *bufio.Writer
	f        *os.File
	size     int64
}

// NewBufioWriter returns a new BufioWriter from fileName.
func NewBufioWriter(fileName string) (BufioWriter, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &bufioWriter{
		fileName: fileName,
		w:        bufio.NewWriterSize(f, defaultWriteBufferSize),
		f:        f,
	}, nil
}

// Reset switches the buffered writer to write to a new file.
func (bw *bufioWriter) Reset(fileName string) error {
	newF, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if err = bw.Close(); err != nil {
		return err
	}
	bw.f = newF
	bw.w = bufio.NewWriterSize(newF, defaultWriteBufferSize)
	bw.size = 0
	bw.fileName = fileName
	return nil
}

// Write writes byte-slice to the buffer.
func (bw *bufioWriter) Write(content []byte) (int, error) {
	var buf [8]byte // buf for store length
	variantLength := binary.PutUvarint(buf[:], uint64(len(content)))
	// write length
	n1, err := bw.w.Write(buf[:variantLength])
	if err != nil {
		return 0, err
	}
	bw.size += int64(n1)
	// write content
	n2, err := bw.w.Write(content)
	if err != nil {
		return 0, err
	}
	bw.size += int64(n2)
	return n1 + n2, nil
}

// Sync flushes the buffered data to the write-queue of the disk.
// It does not wait for the end of the actual write operation of disk.
func (bw *bufioWriter) Sync() error {
	// Flush just flushes data to io.Writer
	if err := bw.w.Flush(); err != nil {
		return err
	}
	// sync syscall
	return bw.f.Sync()
}

// Flush flushes buffered data to the underlying io.Writer.
func (bw *bufioWriter) Flush() error {
	return bw.w.Flush()
}

// Size returns the length of all written data.
func (bw *bufioWriter) Size() int64 {
	return bw.size
}

// Close closes the writer after flushing the buffered data.
func (bw *bufioWriter) Close() error {
	if err := bw.w.Flush(); err != nil {
		return err
	}
	return bw.f.Close()
}

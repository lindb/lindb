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

package compress

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/snappy"
)

// Writer represents compress writer.
type Writer interface {
	io.Writer
	io.Closer
	// Bytes returns compressed binary data.
	Bytes() []byte
}

// Reader represents uncompress reader.
type Reader interface {
	// Uncompress de-compresses compressed binary data.
	Uncompress(compressData []byte) ([]byte, error)
}

type snappyWriter struct {
	writer *snappy.Writer
	buffer bytes.Buffer
}

// NewBufferedWriter creates a snappy compress writer.
func NewSnappyWriter() Writer {
	w := &snappyWriter{}
	w.writer = snappy.NewBufferedWriter(&w.buffer)
	return w
}

func (w *snappyWriter) Write(row []byte) (n int, err error) {
	n, err = w.writer.Write(row)
	return n, err
}

func (w *snappyWriter) Bytes() (r []byte) {
	// TODO: opts bytes copy?
	data := w.buffer.Bytes()
	r = make([]byte, len(data))
	copy(r, data)

	// reset compress context
	w.buffer.Reset()
	w.writer.Reset(&w.buffer)
	return
}

func (w *snappyWriter) Close() error {
	return w.writer.Close()
}

type snappyReader struct {
	reader       *snappy.Reader
	compressed   bytes.Buffer
	decompressed bytes.Buffer
}

// NewSnappyReader creates a snappy uncompress reader.
func NewSnappyReader() Reader {
	r := &snappyReader{}
	r.reader = snappy.NewReader(&r.compressed)
	return r
}

func (r *snappyReader) Uncompress(compressData []byte) ([]byte, error) {
	defer func() {
		// reset uncompress context
		r.compressed.Reset()
		r.decompressed.Reset()
		r.reader.Reset(&r.compressed)
	}()
	_, err := r.compressed.Write(compressData)

	if err == nil {
		_, err = io.Copy(&r.decompressed, r.reader)
	}

	if err != nil {
		return nil, err
	}
	return r.decompressed.Bytes(), nil
}

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

package encoding

import (
	"bytes"

	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./tsd_stream.go -destination=./tsd_stream_mock.go -package encoding

// TSDStreamWriter represents the tsd data points stream writer for multi-fields
type TSDStreamWriter interface {
	// WriteField writes the field data
	WriteField(pFieldID uint16, data []byte)
	// Bytes returns the binary data with time range and fields' data
	Bytes() ([]byte, error)
}

// TSDStreamReader represents the tsd data points reader for multi-fields
type TSDStreamReader interface {
	// TimeRange returns the time slot range
	TimeRange() (start, end uint16)
	// HasNext returns the if has more field data
	HasNext() bool
	// Next returns the fieldID and field data
	Next() (pFieldID uint16, fieldData *TSDDecoder)
	// Close closes the reader and releases the resource
	Close()
}

// tsdStreamWriter implements the TSDStreamWriter interface
type tsdStreamWriter struct {
	writer *stream.BufferWriter
}

// NewTSDStreamWriter creates the tsd stream writer
func NewTSDStreamWriter(start, end uint16) TSDStreamWriter {
	var buf bytes.Buffer
	writer := stream.NewBufferWriter(&buf)
	writer.PutUInt16(start)
	writer.PutUInt16(end)
	return &tsdStreamWriter{writer: writer}
}

// WriteField writes the field data
func (sw *tsdStreamWriter) WriteField(pFieldID uint16, data []byte) {
	sw.writer.PutUInt16(pFieldID)
	sw.writer.PutUvarint32(uint32(len(data)))
	sw.writer.PutBytes(data)
}

// Bytes returns the binary data with time range and fields' data
func (sw *tsdStreamWriter) Bytes() ([]byte, error) {
	return sw.writer.Bytes()
}

// tsdStreamReader implements the TSDStreamReader interface
type tsdStreamReader struct {
	reader    *stream.Reader
	fieldData *TSDDecoder
	startTime uint16
	endTime   uint16
}

// NewTSDStreamReader creates the tsd stream reader
func NewTSDStreamReader(data []byte) TSDStreamReader {
	reader := stream.NewReader(data)
	startTime := reader.ReadUint16()
	endTime := reader.ReadUint16()
	return &tsdStreamReader{
		startTime: startTime,
		endTime:   endTime,
		reader:    reader,
		fieldData: GetTSDDecoder(),
	}
}

// TimeRange returns the time slot range
func (sr *tsdStreamReader) TimeRange() (start, end uint16) {
	return sr.startTime, sr.endTime
}

// HasNext returns the if it has more field data
func (sr *tsdStreamReader) HasNext() bool {
	return !sr.reader.Empty()
}

// Next returns the fieldID and field data
func (sr *tsdStreamReader) Next() (pFieldID uint16, fieldData *TSDDecoder) {
	pFieldID = sr.reader.ReadUint16()
	data := sr.reader.ReadSlice(int(sr.reader.ReadUvarint32()))
	sr.fieldData.ResetWithTimeRange(data, sr.startTime, sr.endTime)
	return pFieldID, sr.fieldData
}

// Close closes the reader and releases the resource
func (sr *tsdStreamReader) Close() {
	ReleaseTSDDecoder(sr.fieldData)
}

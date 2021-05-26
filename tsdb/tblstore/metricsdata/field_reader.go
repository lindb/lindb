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

package metricsdata

import (
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./field_reader.go -destination=./field_reader_mock.go -package metricsdata

// FieldReader represents the field metricReader when does metric data merge.
// !!!!NOTICE: need get field value in order by field
type FieldReader interface {
	// slotRange returns the time slot range of metric level
	slotRange() (start, end uint16)
	// getFieldData returns the field data by field id,
	// if metricReader is completed, return nil, if found data returns field data else returns nil
	getFieldData(fieldID field.ID) []byte
	// reset resets the field data for reading
	reset(buf []byte, position int, start, end uint16)
	// close closes the metricReader
	close()
}

// fieldReader implements FieldReader
type fieldReader struct {
	start, end   uint16
	seriesData   []byte
	fieldOffsets *encoding.FixedOffsetDecoder
	fieldIndexes map[field.ID]int
	fieldCount   int

	completed bool // !!!!NOTICE: need reset completed
}

// newFieldReader creates the field metricReader
func newFieldReader(fieldIndexes map[field.ID]int, buf []byte, position int, start, end uint16) FieldReader {
	r := &fieldReader{
		fieldIndexes: fieldIndexes,
		fieldCount:   len(fieldIndexes),
	}
	r.reset(buf, position, start, end)
	return r
}

// reset resets the field data for reading
func (r *fieldReader) reset(buf []byte, position int, start, end uint16) {
	r.completed = false
	r.start = start
	r.end = end
	if r.fieldCount == 1 {
		r.seriesData = buf
		return
	}
	data := buf[position:]
	r.fieldOffsets = encoding.NewFixedOffsetDecoder(data)
	r.seriesData = data[r.fieldOffsets.Header()+r.fieldCount*r.fieldOffsets.ValueWidth():]
	r.fieldOffsets = encoding.NewFixedOffsetDecoder(buf[position:])
}

// slotRange returns the time slot range of metric level
func (r *fieldReader) slotRange() (start, end uint16) {
	return r.start, r.end
}

// getFieldData returns the field data by field id,
// if metricReader is completed, return nil, if found data returns field data else returns nil
func (r *fieldReader) getFieldData(fieldID field.ID) []byte {
	if r.completed {
		return nil
	}
	idx, ok := r.fieldIndexes[fieldID]
	if !ok {
		return nil
	}
	if r.fieldCount == 1 {
		return r.seriesData
	}
	offset, ok := r.fieldOffsets.Get(idx)
	if !ok {
		return nil
	}
	return r.seriesData[offset:]
}

// close marks the metricReader completed
func (r *fieldReader) close() {
	r.completed = true
}

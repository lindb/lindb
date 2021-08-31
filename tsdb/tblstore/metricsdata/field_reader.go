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
	"encoding/binary"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./field_reader.go -destination=./field_reader_mock.go -package metricsdata

// FieldReader represents the field metricReader when does metric data merge.
// !!!!NOTICE: need get field value in order by field
type FieldReader interface {
	// SlotRange returns the time slot range of metric level
	SlotRange() timeutil.SlotRange
	// GetFieldData returns the field data by field id,
	// if metricReader is completed, return nil, if found data returns field data else returns nil
	GetFieldData(fieldID field.ID) []byte
	// Reset resets the field data for reading
	Reset(seriesEntry []byte, slotRange timeutil.SlotRange)
	// Close closes the metricReader
	Close()
}

// fieldReader implements FieldReader
type fieldReader struct {
	slotRange    timeutil.SlotRange
	seriesEntry  []byte
	fieldOffsets *encoding.FixedOffsetDecoder
	fieldDatas   []byte
	fieldIndexes map[field.ID]int
	fieldCount   int

	completed bool // !!!!NOTICE: need reset completed
}

// newFieldReader creates the field metricReader
func newFieldReader(fieldIndexes map[field.ID]int, seriesEntry []byte, slotRange timeutil.SlotRange) FieldReader {
	r := &fieldReader{
		fieldIndexes: fieldIndexes,
		fieldCount:   len(fieldIndexes),
		seriesEntry:  seriesEntry,
		slotRange:    slotRange,
		fieldOffsets: encoding.NewFixedOffsetDecoder(),
	}
	r.Reset(seriesEntry, slotRange)
	return r
}

// Reset resets the field data for reading
func (r *fieldReader) Reset(seriesEntry []byte, slotRange timeutil.SlotRange) {
	r.completed = false
	r.slotRange = slotRange
	if r.fieldCount == 1 {
		r.seriesEntry = seriesEntry
		return
	}
	if len(seriesEntry) <= 1 {
		r.completed = true
		return
	}
	// reversed decoding binary.Uvariant
	var (
		fieldOffsetsLen uint64
		readBytes       uint
	)
	{
		for j := len(seriesEntry) - 1; j >= 0; j-- {
			b := seriesEntry[j]
			i := len(seriesEntry) - j - 1
			if b < 0x80 {
				if i >= binary.MaxVarintLen64 || i == binary.MaxVarintLen64-1 && b > 1 {
					// overflow
					r.completed = true
					return
				}
				fieldOffsetsLen |= uint64(b) << readBytes
				readBytes = uint(i + 1)
				break
			}
			fieldOffsetsLen |= uint64(b&0x7f) << readBytes
			readBytes += 7
		}
	}
	fieldOffsetsAt := len(seriesEntry) - int(readBytes) - int(fieldOffsetsLen)
	if fieldOffsetsAt < 0 || fieldOffsetsAt >= len(seriesEntry) {
		r.completed = true
		return
	}
	if _, err := r.fieldOffsets.Unmarshal(seriesEntry[fieldOffsetsAt:]); err != nil {
		r.completed = true
	}
	r.fieldDatas = seriesEntry[:fieldOffsetsAt]
}

// SlotRange returns the time slot range of metric level
func (r *fieldReader) SlotRange() timeutil.SlotRange {
	return r.slotRange
}

// GetFieldData returns the field data by field id,
// if metricReader is completed, return nil, if found data returns field data else returns nil
func (r *fieldReader) GetFieldData(fieldID field.ID) []byte {
	if r.completed {
		return nil
	}
	idx, ok := r.fieldIndexes[fieldID]
	if !ok {
		return nil
	}
	if r.fieldCount == 1 {
		return r.seriesEntry
	}
	fieldBlock, err := r.fieldOffsets.GetBlock(idx, r.fieldDatas)
	if err != nil {
		return nil
	}
	return fieldBlock
}

// Close marks the metricReader completed
func (r *fieldReader) Close() {
	r.completed = true
}

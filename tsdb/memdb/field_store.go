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

package memdb

import (
	"encoding/binary"
	"math"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./field_store.go -destination=./field_store_mock.go -package memdb

// memory layout as below:
// header: field id[2bytes]
//        + start time[2byte] + end time(delta of start time)[1byte] + mark container[2byte]
// body: data points(value.....)
// last mark flag of container marks the buf if has data written

const (
	fieldOffset = 0
	startOffset = fieldOffset + 2
	endOffset   = startOffset + 2
	markOffset  = endOffset + 1
	bodyOffset  = markOffset + 2
	headLen     = 8
	valueSize   = 8

	emptyFieldStoreSize = 24 + // empty buf slice cost
		24 // empty compress slice cost
)

// fStoreINTF represents field-store,
// which abstracts a store for storing field data based on family start time + field id
type fStoreINTF interface {
	// Capacity returns the size usage
	Capacity() int
	// GetFieldID returns the field id of metric level
	GetFieldID() field.ID
	// Write writes the field data into current buffer
	// if time slot out of current time window, need compress time window then resets the current buffer
	// if it has same time slot in current buffer, need do rollup operation by field type
	Write(fieldType field.Type, slotIndex uint16, value float64)
	// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
	FlushFieldTo(tableFlusher metricsdata.Flusher, fieldMeta field.Meta, flushCtx *flushContext) error
	// Load loads field series data.
	Load(fieldType field.Type, slotRange timeutil.SlotRange) []byte
}

// fieldStore implements fStoreINTF interface
type fieldStore struct {
	buf      []byte // current write buffer, accept write data
	compress []byte // immutable compress data
}

// newFieldStore creates a new field store
func newFieldStore(buf []byte, fieldID field.ID) fStoreINTF {
	stream.PutUint16(buf, fieldOffset, uint16(fieldID))
	return &fieldStore{
		buf: buf,
	}
}

// GetFieldID returns the field id of metric level
func (fs *fieldStore) GetFieldID() field.ID {
	return field.ID(stream.ReadUint16(fs.buf, fieldOffset))
}

func (fs *fieldStore) Write(fieldType field.Type, slotIndex uint16, value float64) {
	if fs.buf[markOffset+1] == 0 {
		// no data written before
		fs.writeFirstPoint(slotIndex, value)
		return
	}

	startTime := fs.getStart()
	if slotIndex < startTime || slotIndex > startTime+fs.timeWindow()-1 {
		// if current slot time out of current time window, need compress block data, start new time window
		fs.compact(fieldType, startTime)

		// write first point after compact
		fs.writeFirstPoint(slotIndex, value)
		return
	}

	// write data in current write buffer
	delta := slotIndex - startTime
	pos, markIdx, flagIdx := fs.position(delta)
	if fs.buf[markOffset+markIdx]&flagIdx != 0 {
		// there is same point of same time slot
		oldValue := math.Float64frombits(binary.LittleEndian.Uint64(fs.buf[pos:]))
		value = fieldType.AggType().Aggregate(oldValue, value)
	} else {
		// new data for time slot
		fs.buf[endOffset] = byte(delta)
		fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
	}
	// finally, write value into the body of current write buffer
	binary.LittleEndian.PutUint64(fs.buf[pos:], math.Float64bits(value))
}

// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
func (fs *fieldStore) FlushFieldTo(tableFlusher metricsdata.Flusher, fieldMeta field.Meta, flushCtx *flushContext) error {
	var decoder *encoding.TSDDecoder
	if len(fs.compress) > 0 {
		// calc new start/end based on old compress values
		decoder = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(decoder)
		decoder.Reset(fs.compress)
	}

	encoder := tableFlusher.GetEncoder(flushCtx.fieldIdx)
	encoder.RestWithStartTime(flushCtx.SlotRange.Start)

	data, err := fs.merge(fieldMeta.Type, encoder, decoder, fs.getStart(), flushCtx.SlotRange, false)
	if err != nil {
		memDBLogger.Error("flush field store err, data lost", logger.Error(err))
		return nil
	}
	return tableFlusher.FlushField(data)
}

// writeFirstPoint writes first point in current write buffer
func (fs *fieldStore) writeFirstPoint(slotIndex uint16, value float64) {
	pos, markIdx, flagIdx := fs.position(0)
	binary.LittleEndian.PutUint16(fs.buf[startOffset:], slotIndex) // write start time
	fs.buf[endOffset] = 0
	fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
	fs.buf[markOffset+1] |= 1             // last mark flag marks if buf has data written
	binary.LittleEndian.PutUint64(fs.buf[pos:], math.Float64bits(value))
}

// timeWindow returns the time window of current write buffer
func (fs *fieldStore) timeWindow() uint16 {
	return uint16((len(fs.buf) - headLen) / valueSize)
}

// resetBuf resets the writer buffer mark, makes the current buffer is new
func (fs *fieldStore) resetBuf() {
	fs.buf[markOffset] = 0
	fs.buf[markOffset+1] = 0
}

func (fs *fieldStore) Capacity() int {
	// notice: do not use cap as it's a allocated page
	return cap(fs.compress) + len(fs.buf) + emptyFieldStoreSize
}

// compact the current write buffer,
// new compress operation will be executed when it's necessary
func (fs *fieldStore) compact(fieldType field.Type, startTime uint16) {
	length := len(fs.compress)
	thisSlotRange := fs.slotRange(startTime)

	var decoder *encoding.TSDDecoder
	if length > 0 {
		// if has compress data, create tsd decoder for merge compress
		decoder = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(decoder)
		decoder.Reset(fs.compress)
	}
	encoder := encoding.TSDEncodeFunc(thisSlotRange.Start)
	defer encoding.ReleaseTSDEncoder(encoder)

	data, err := fs.merge(fieldType, encoder, decoder, startTime, thisSlotRange, true)
	if err != nil {
		memDBLogger.Error("compact field store data err", logger.Error(err))
	}

	fs.compress = encoding.MustCopy(fs.compress, data)
	// !!!!! IMPORTANT: need reset current write buffer
	fs.resetBuf()
}

// position returns the point write position/mark index/flag index
func (fs *fieldStore) position(deltaOfTime uint16) (pos, markIdx uint16, flagIdx uint8) {
	pos = bodyOffset + valueSize*deltaOfTime
	markIdx = deltaOfTime / valueSize
	flagIdx = uint8(1 << (valueSize - deltaOfTime%valueSize - 1))
	return
}

// getStart returns the start time in current write buffer
func (fs *fieldStore) getStart() uint16 {
	return stream.ReadUint16(fs.buf, startOffset)
}

// getEnd returns the delta time of start time in current write buffer
func (fs *fieldStore) getEnd() uint16 {
	return uint16(fs.buf[endOffset])
}

// merge the current and compress data based on field aggregate function,
// startTime => current write start time
// start/end slot => target compact time slot
func (fs *fieldStore) merge(
	fieldType field.Type,
	encoder *encoding.TSDEncoder,
	decoder *encoding.TSDDecoder,
	startTime uint16,
	thisSlotRange timeutil.SlotRange,
	withTimeRange bool,
) (compress []byte, err error) {
	for i := thisSlotRange.Start; i <= thisSlotRange.End; i++ {
		newValue, hasNewValue := fs.getCurrentValue(startTime, i)
		oldValue, hasOldValue := getOldFloatValue(decoder, i)
		switch {
		case hasNewValue && !hasOldValue:
			// just compress current block value with pos
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(newValue))
		case hasNewValue && hasOldValue:
			// merge and compress
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(fieldType.AggType().Aggregate(newValue, oldValue)))
		case !hasNewValue && hasOldValue:
			// compress old value
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(oldValue))
		default:
			// append empty value
			encoder.AppendTime(bit.Zero)
		}
	}
	if withTimeRange {
		compress, err = encoder.Bytes()
		if err != nil {
			return nil, err
		}
		return compress, err
	}
	// get compress data without time slot range
	compress, err = encoder.BytesWithoutTime()
	if err != nil {
		return nil, err
	}
	return compress, err
}

// Load loads field series data.
func (fs *fieldStore) Load(fieldType field.Type, slotRange timeutil.SlotRange) []byte {
	var tsd *encoding.TSDDecoder
	size := len(fs.compress)
	if size > 0 {
		// calc new start/end based on old compress values
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	// todo: pool encoder after loading ?
	encoder := encoding.NewTSDEncoder(slotRange.Start)
	data, err := fs.merge(fieldType, encoder, tsd, fs.getStart(), slotRange, false)
	if err != nil {
		memDBLogger.Error("load field store err", logger.Error(err))
		return nil
	}
	return data
}

// slotRange returns time slot range in current/compress buffer
func (fs *fieldStore) slotRange(currentStart uint16) timeutil.SlotRange {
	startSlot := currentStart
	endSlot := currentStart + fs.getEnd()
	if len(fs.compress) == 0 {
		return timeutil.NewSlotRange(startSlot, endSlot)
	}
	start, end := encoding.DecodeTSDTime(fs.compress)
	return getTimeSlotRange(start, end, startSlot, endSlot)
}

// getTimeSlotRange returns the final time slot range based on start/end
func getTimeSlotRange(startSlot1, endSlot1, startSlot2, endSlot2 uint16) timeutil.SlotRange {
	sr := timeutil.NewSlotRange(startSlot1, endSlot1)
	if sr.End < endSlot2 {
		sr.End = endSlot2
	}
	if sr.Start > startSlot2 {
		sr.Start = startSlot2
	}
	return sr
}

// getCurrentValue returns the value in current write buffer
func (fs *fieldStore) getCurrentValue(startTime, timeSlot uint16) (value float64, hasValue bool) {
	if timeSlot < startTime || timeSlot > startTime+fs.getEnd() {
		return
	}
	delta := timeSlot - startTime
	pos, markIdx, flagIdx := fs.position(delta)
	if fs.buf[markOffset+markIdx]&flagIdx == 0 {
		return
	}
	hasValue = true
	value = math.Float64frombits(binary.LittleEndian.Uint64(fs.buf[pos:]))
	return
}

// getOldFloatValue returns the value in compress buffer
func getOldFloatValue(tsd *encoding.TSDDecoder, timeSlot uint16) (value float64, hasValue bool) {
	if tsd == nil {
		return
	}
	if !tsd.HasValueWithSlot(timeSlot) {
		return
	}
	hasValue = true
	value = math.Float64frombits(tsd.Value())
	return
}

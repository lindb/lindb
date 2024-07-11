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

	commonencoding "github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

// memory layout as below:
// header: start time[2byte] + end time(delta of start time)[1byte] + mark container[2byte]
// body: data points(value.....)
// last mark flag of container marks the buf if it has data written

const (
	startOffset   = 0 // start time
	endOffset     = startOffset + 2
	markOffset    = endOffset + 1
	bodyOffset    = markOffset + 2
	headLen       = 8
	valueSize     = 8
	markContainer = 8
)

// write writes field data into temp buffer.
// if time slot out of current time window, need compress time window then resets the current buffer
// if it has same time slot in current buffer, need do rollup operation by field type
func write(md *memoryDatabase, buf []byte, memTimeSeries uint32, fieldIndex uint8, fieldType field.Type, slotIndex uint16, value float64) {
	if buf[markOffset+1] == 0 {
		// no data written before
		writeFirstPoint(buf, slotIndex, value)
		return
	}

	startTime := getStart(buf)
	if slotIndex < startTime || slotIndex > startTime+timeWindow(buf)-1 {
		// if current slot time out of current time window, need compress block data, start new time window
		compact(md, buf, memTimeSeries, fieldIndex, fieldType, startTime)

		// write first point after compact
		writeFirstPoint(buf, slotIndex, value)
		return
	}

	// write data in current write buffer
	delta := slotIndex - startTime
	pos, markIdx, flagIdx := position(delta)
	if buf[markOffset+markIdx]&flagIdx != 0 {
		// there is same point of same time slot
		oldValue := encoding.BytesToFloat64(buf[pos : pos+8])
		value = fieldType.AggType().Aggregate(oldValue, value)
	} else {
		// new data for time slot
		buf[endOffset] = byte(delta)
		buf[markOffset+markIdx] |= flagIdx // mark value exist
	}
	// finally, write value into the body of current write buffer
	copy(buf[pos:], encoding.Float64ToBytes(value))
}

// compact the current write buffer,
// new compress operation will be executed when it's necessary
func compact(md *memoryDatabase, buf []byte, memTimeSeries uint32, fieldIndex uint8, fieldType field.Type, startTime uint16) {
	compress := md.getFieldCompressBuffer(memTimeSeries, fieldIndex)
	length := len(compress)
	thisSlotRange := slotRange(startTime, buf, compress)

	var decoder *encoding.TSDDecoder
	if length > 0 {
		// if has compress data, create tsd decoder for merge compress
		decoder = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(decoder)
		decoder.Reset(compress)
	}
	encoder := encoding.TSDEncodeFunc(thisSlotRange.Start)
	defer encoding.ReleaseTSDEncoder(encoder)

	data, err := merge(fieldType, buf, encoder, decoder, startTime, thisSlotRange, true)
	if err != nil {
		// NOTE: lost data
		memDBLogger.Error("compact field store data err, lost data", logger.Error(err))
	} else {
		compress = commonencoding.MustCopy(compress, data)
		md.storeFieldComressBuffer(memTimeSeries, fieldIndex, compress)
	}

	// NOTE: !!!!! IMPORTANT: need reset current write buffer
	resetBuf(buf)
}

// writeFirstPoint writes first point in current write buffer.
func writeFirstPoint(buf []byte, slotIndex uint16, value float64) {
	pos, markIdx, flagIdx := position(0)
	binary.LittleEndian.PutUint16(buf[startOffset:], slotIndex) // write start time
	buf[endOffset] = 0
	buf[markOffset+markIdx] |= flagIdx // mark value exist
	buf[markOffset+1] |= 1             // last mark flag marks if buf has data written
	copy(buf[pos:], encoding.Float64ToBytes(value))
}

// merge the current and compress data based on field aggregate function,
// startTime => current write start time
// start/end slot => target compact time slot
func merge(
	fieldType field.Type,
	buf []byte,
	encoder *encoding.TSDEncoder,
	decoder *encoding.TSDDecoder,
	startTime uint16,
	thisSlotRange timeutil.SlotRange,
	withTimeRange bool,
) (compress []byte, err error) {
	for i := thisSlotRange.Start; i <= thisSlotRange.End; i++ {
		newValue, hasNewValue := getCurrentValue(buf, startTime, i)
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

func timeWindow(buf []byte) uint16 {
	return uint16((len(buf) - headLen) / valueSize)
}

// position returns the point write position/mark index/flag index.
// position: write value
// markIdx: mark container index
// flagIdx: flag if pos has value
func position(deltaOfTime uint16) (pos, markIdx uint16, flagIdx uint8) {
	pos = bodyOffset + valueSize*deltaOfTime
	markIdx = deltaOfTime / valueSize
	flagIdx = uint8(1 << (markContainer - deltaOfTime%markContainer - 1))
	return
}

// getStart returns the start time in current write buffer
func getStart(buf []byte) uint16 {
	return stream.ReadUint16(buf, startOffset)
}

// getEnd returns the delta time of start time in current write buffer
func getEnd(buf []byte) uint16 {
	return uint16(buf[endOffset])
}

// slotRange returns time slot range in current/compress buffer
func slotRange(currentStart uint16, buf, compress []byte) timeutil.SlotRange {
	startSlot := currentStart
	endSlot := currentStart + getEnd(buf)
	if len(compress) == 0 {
		return timeutil.NewSlotRange(startSlot, endSlot)
	}
	start, end := encoding.DecodeTSDTime(compress)
	return getTimeSlotRange(start, end, startSlot, endSlot)
}

// resetBuf resets the writer buffer mark, makes the current buffer is new
func resetBuf(buf []byte) {
	buf[markOffset] = 0
	buf[markOffset+1] = 0
}

// getCurrentValue returns the value in current write buffer
func getCurrentValue(buf []byte, startTime, timeSlot uint16) (value float64, hasValue bool) {
	if timeSlot < startTime || timeSlot > startTime+getEnd(buf) {
		return
	}
	delta := timeSlot - startTime
	pos, markIdx, flagIdx := position(delta)
	if buf[markOffset+markIdx]&flagIdx == 0 {
		return
	}
	hasValue = true
	value = math.Float64frombits(binary.LittleEndian.Uint64(buf[pos:]))
	return
}

// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
func flushFieldTo(md *memoryDatabase, memTimeSeries uint32, buf []byte,
	slotRange timeutil.SlotRange, tableFlusher metricsdata.Flusher, flushIndex int, fieldMeta field.Meta,
) error {
	compress := md.getFieldCompressBuffer(memTimeSeries, fieldMeta.Index)
	var decoder *encoding.TSDDecoder
	if len(compress) > 0 {
		// calc new start/end based on old compress values
		decoder = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(decoder)
		decoder.Reset(compress)
	}

	encoder := tableFlusher.GetEncoder(flushIndex)
	encoder.RestWithStartTime(slotRange.Start)

	data, err := merge(fieldMeta.Type, buf, encoder, decoder, getStart(buf), slotRange, false)
	if err != nil {
		memDBLogger.Error("flush field store err, data lost", logger.Error(err))
		return nil
	}
	return tableFlusher.FlushField(data)
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

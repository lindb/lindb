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

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./field_store.go -destination=./field_store_mock.go -package memdb

var (
	fieldStoreMergeFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mem_field_store_merge_fail",
			Help: "Field Store merge fail when flush.",
		},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(fieldStoreMergeFailCounter)
}

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
	// GetFieldID returns the field id of metric level
	GetFieldID() field.ID
	// Write writes the field data into current buffer, returns the written size.
	// if time slot out of current time window, need compress time window then resets the current buffer
	// if has same time slot in current buffer, need do rollup operation by field type
	Write(fieldType field.Type, slotIndex uint16, value float64) (writtenSize int)
	// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
	FlushFieldTo(tableFlusher metricsdata.Flusher, fieldMeta field.Meta, flushCtx flushContext)
	// Load loads field series data.
	Load(fieldType field.Type) []byte
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

// Write writes the field data into current buffer, returns the written size.
// if time slot out of current time window, need compress time window then resets the current buffer
// if has same time slot in current buffer, need do rollup operation by field type
func (fs *fieldStore) Write(fieldType field.Type, slotIndex uint16, value float64) (writtenSize int) {
	if fs.buf[markOffset+1] == 0 {
		// no data written before
		return fs.writeFirstPoint(slotIndex, value)
	}

	startTime := fs.getStart()
	if slotIndex < startTime || slotIndex > startTime+fs.timeWindow()-1 {
		// if current slot time out of current time window, need compress block data, start new time window
		writtenSize = fs.compact(fieldType, startTime)

		// write first point after compact
		writtenSize += fs.writeFirstPoint(slotIndex, value)
		return writtenSize
	}

	// write data in current write buffer
	delta := slotIndex - startTime
	pos, markIdx, flagIdx := fs.position(delta)
	if fs.buf[markOffset+markIdx]&flagIdx != 0 {
		// has same point of same time slot
		aggFunc := fieldType.GetAggFunc()
		oldValue := math.Float64frombits(binary.LittleEndian.Uint64(fs.buf[pos:]))
		value = aggFunc.Aggregate(oldValue, value)
	} else {
		// new data for time slot
		fs.buf[endOffset] = byte(delta)
		fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
		writtenSize += valueSize
	}
	// finally write value into the body of current write buffer
	binary.LittleEndian.PutUint64(fs.buf[pos:], math.Float64bits(value))
	return writtenSize
}

// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
func (fs *fieldStore) FlushFieldTo(tableFlusher metricsdata.Flusher, fieldMeta field.Meta, flushCtx flushContext) {
	aggFunc := fieldMeta.Type.GetAggFunc()
	var tsd *encoding.TSDDecoder
	size := len(fs.compress)
	if size > 0 {
		// calc new start/end based on old compress values
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	data, _, err := fs.merge(aggFunc, tsd, fs.getStart(), flushCtx.SlotRange, false)
	if err != nil {
		fieldStoreMergeFailCounter.Inc()
		memDBLogger.Error("flush field store err, data lost", logger.Error(err))
		return
	}

	tableFlusher.FlushField(data)
}

// writeFirstPoint writes first point in current write buffer
func (fs *fieldStore) writeFirstPoint(slotIndex uint16, value float64) (writtenSize int) {
	pos, markIdx, flagIdx := fs.position(0)
	binary.LittleEndian.PutUint16(fs.buf[startOffset:], slotIndex) // write start time
	fs.buf[endOffset] = 0
	fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
	fs.buf[markOffset+1] |= 1             // last mark flag marks if buf has data written
	binary.LittleEndian.PutUint64(fs.buf[pos:], math.Float64bits(value))
	return valueSize + headLen
}

// timeWindow returns the time window of current write buffer
func (fs *fieldStore) timeWindow() uint16 {
	return uint16((len(fs.buf) - headLen) / valueSize)
}

// resetBuf resets the write buffer mark, makes the current buffer is new
func (fs *fieldStore) resetBuf() {
	fs.buf[markOffset] = 0
	fs.buf[markOffset+1] = 0
}

// compact compacts the current write buffer,
// new compress operation will be executed when it's necessary
func (fs *fieldStore) compact(fieldType field.Type, startTime uint16) (size int) {
	length := len(fs.compress)
	thisSlotRange := fs.slotRange(startTime)

	aggFunc := fieldType.GetAggFunc()
	var tsd *encoding.TSDDecoder
	if length > 0 {
		// if has compress data, create tsd decoder for merge compress
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	data, freeSize, err := fs.merge(aggFunc, tsd, startTime, thisSlotRange, true)
	if err != nil {
		memDBLogger.Error("compact field store data err", logger.Error(err))
	}

	fs.compress = data
	// !!!!! IMPORTANT: need reset current write buffer
	fs.resetBuf()
	return len(fs.compress) - length - freeSize
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

// merge merges the current and compress data based on field aggregate function,
// startTime => current write start time
// start/end slot => target compact time slot
func (fs *fieldStore) merge(
	aggFunc field.AggFunc,
	tsd *encoding.TSDDecoder,
	startTime uint16,
	thisSlotRange timeutil.SlotRange,
	withTimeRange bool,
) (compress []byte, freeSize int, err error) {
	encode := encoding.TSDEncodeFunc(thisSlotRange.Start)
	for i := thisSlotRange.Start; i <= thisSlotRange.End; i++ {
		newValue, hasNewValue := fs.getCurrentValue(startTime, i)
		oldValue, hasOldValue := getOldFloatValue(tsd, i)
		switch {
		case hasNewValue && !hasOldValue:
			// just compress current block value with pos
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(newValue))
		case hasNewValue && hasOldValue:
			// merge and compress
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(aggFunc.Aggregate(newValue, oldValue)))
		case !hasNewValue && hasOldValue:
			// compress old value
			encode.AppendTime(bit.One)
			encode.AppendValue(math.Float64bits(oldValue))
		default:
			// append empty value
			encode.AppendTime(bit.Zero)
		}
		// if has value calc the free size
		if hasNewValue {
			freeSize += valueSize
		}
	}
	if withTimeRange {
		compress, err = encode.Bytes()
		if err != nil {
			return nil, 0, err
		}
		return compress, freeSize, err
	}
	// get compress data without time slot range
	compress, err = encode.BytesWithoutTime()
	if err != nil {
		return nil, 0, err
	}
	return compress, freeSize, err
}

// Load loads field series data.
func (fs *fieldStore) Load(fieldType field.Type) []byte {
	//TODO check if need do compact
	_ = fs.compact(fieldType, fs.getStart())
	rs := make([]byte, len(fs.compress))
	copy(rs, fs.compress)
	//TODO remove time range???
	return rs
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
func getTimeSlotRange(startSlot1, endSlot1 uint16, startSlot2, endSlot2 uint16) timeutil.SlotRange {
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
func (fs *fieldStore) getCurrentValue(startTime uint16, timeSlot uint16) (value float64, hasValue bool) {
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

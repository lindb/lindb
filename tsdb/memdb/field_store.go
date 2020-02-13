package memdb

import (
	"encoding/binary"
	"math"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./field_store.go -destination=./field_store_mock.go -package memdb

// for testing
var (
	encodeFunc = encoding.NewTSDEncoder
)

// memory layout as below:
// header: family[1byte] + field id[1byte] + primitive id[1byte]
//        + start time[2byte] + end time(delta of start time)[1byte] + mark container[2byte]
// body: data points(value.....)
// last mark flag of container marks the buf if has data written

const (
	familyOffset    = 0
	fieldOffset     = familyOffset + 1
	primitiveOffset = fieldOffset + 1
	startOffset     = primitiveOffset + 1
	endOffset       = startOffset + 2
	markOffset      = endOffset + 1
	bodyOffset      = markOffset + 2
	headLen         = 8
	valueSize       = 8

	emptyPrimitiveFieldStoreSize = 8 + // buf pointer
		8 // compress pointer
)

// fStoreINTF represents field-store,
// which abstracts a store for storing field data based on family start time + field id + primitive field id
type fStoreINTF interface {
	// GetKey returns the field store key, sorts in field list will use this key for sorting
	// field key = family id + field id + primitive field id
	GetKey() uint32
	// GetFamilyID returns the family time mapping id
	GetFamilyID() familyID
	// GetFieldKey returns the field key
	// field key = field id + primitive field id
	GetFieldKey() field.Key
	// GetFieldID returns the field id of metric level
	GetFieldID() field.ID
	// GetPrimitiveID returns the primitive field id of field level, some complex field will have many primitive fields
	GetPrimitiveID() field.PrimitiveID
	// Write writes the field data into current buffer, returns the written size.
	// if time slot out of current time window, need compress time window then resets the current buffer
	// if has same time slot in current buffer, need do rollup operation by field type
	Write(fieldType field.Type, slotIndex uint16, value float64) (writtenSize int)
	// FlushFieldTo flushes field store data into kv store, need align slot range in metric level
	FlushFieldTo(tableFlusher metricsdata.Flusher, flushCtx flushContext)
	// Load loads field store data based on query time range, then aggregate the result
	Load(fieldType field.Type, agg aggregation.PrimitiveAggregator, memScanCtx *memScanContext)
}

// fieldStore implements fStoreINTF interface
type fieldStore struct {
	buf      []byte // current write buffer, accept write data
	compress []byte // immutable compress data
}

// newFieldStore creates a new field store
func newFieldStore(buf []byte, familyID familyID, fieldID field.ID, pFieldID field.PrimitiveID) fStoreINTF {
	buf[familyOffset] = byte(familyID)
	buf[fieldOffset] = byte(fieldID)
	buf[primitiveOffset] = byte(pFieldID)
	return &fieldStore{
		buf: buf,
	}
}

// GetKey returns the field store key, sorts in field list will use this key for sorting
// field key = family id + field id + primitive field id
func (fs *fieldStore) GetKey() uint32 {
	return uint32(fs.buf[primitiveOffset]) | uint32(fs.buf[fieldOffset])<<8 | uint32(fs.buf[familyOffset])<<16
}

// GetFamilyID returns the family time mapping id
func (fs *fieldStore) GetFamilyID() familyID {
	return familyID(fs.buf[familyOffset])
}

// GetFieldKey returns the field key
// field key = field id + primitive field id
func (fs *fieldStore) GetFieldKey() field.Key {
	return field.Key(binary.LittleEndian.Uint16(fs.buf[fieldOffset:]))
}

// GetFieldID returns the field id of metric level
func (fs *fieldStore) GetFieldID() field.ID {
	return field.ID(fs.buf[fieldOffset])
}

// GetPrimitiveID returns the primitive field id of field level, some complex field will have many primitive fields
func (fs *fieldStore) GetPrimitiveID() field.PrimitiveID {
	return field.PrimitiveID(fs.buf[primitiveOffset])
}

// Write writes the field data into current buffer, returns the written size.
// if time slot out of current time window, need compress time window then resets the current buffer
// if has same time slot in current buffer, need do rollup operation by field type
func (fs *fieldStore) Write(fieldType field.Type, slotIndex uint16, value float64) (writtenSize int) {
	if fs.buf[markOffset+2] == 0 {
		// no data written before
		return fs.writeFirstPoint(slotIndex, value)
	}

	startTime := fs.getStart()
	if slotIndex < startTime || slotIndex > startTime+fs.timeWindow()-1 {
		// if current slot time out of current time window, need compress block data, start new time window
		writtenSize = fs.compact(fieldType, startTime)
		// !!!!! IMPORTANT: need reset current write buffer
		fs.resetBuf()
		// write first point after compact
		writtenSize += fs.writeFirstPoint(slotIndex, value)
		return writtenSize
	}

	// write data in current write buffer
	delta := slotIndex - startTime
	pos, markIdx, flagIdx := fs.position(delta)
	if fs.buf[markOffset+markIdx]&flagIdx != 0 {
		// has same point of same time slot
		aggFunc := fieldType.GetSchema().GetAggFunc(fs.GetPrimitiveID())
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
func (fs *fieldStore) FlushFieldTo(tableFlusher metricsdata.Flusher, flushCtx flushContext) {
	fieldMeta, ok := tableFlusher.GetFieldMeta(fs.GetFieldID())
	if !ok {
		memDBLogger.Error("field meta not exist in flush context when flush field store data")
		return
	}
	aggFunc := fieldMeta.Type.GetSchema().GetAggFunc(fs.GetPrimitiveID())
	var tsd *encoding.TSDDecoder
	size := len(fs.compress)
	if size > 0 {
		// calc new start/end based on old compress values
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	data, _, err := fs.merge(aggFunc, tsd, fs.getStart(), flushCtx.start, flushCtx.end, false)
	if err != nil {
		//FIXME stone100 add metric
		memDBLogger.Error("flush field store err, data lost", logger.Error(err))
		return
	}

	tableFlusher.FlushField(fs.GetFieldKey(), data)
}

// writeFirstPoint writes first point in current write buffer
func (fs *fieldStore) writeFirstPoint(slotIndex uint16, value float64) (writtenSize int) {
	pos, markIdx, flagIdx := fs.position(0)
	binary.LittleEndian.PutUint16(fs.buf[startOffset:], slotIndex) // write start time
	fs.buf[endOffset] = 0
	fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
	fs.buf[markOffset+2] |= 1             // last mark flag marks if buf has data written
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
	fs.buf[markOffset+2] = 0
}

// compact compacts the current write buffer, if has compress need do merge operation generate new compress data
func (fs *fieldStore) compact(fieldType field.Type, startTime uint16) (memSize int) {
	size := len(fs.compress)
	start, end := fs.slotRange(startTime)

	aggFunc := fieldType.GetSchema().GetAggFunc(fs.GetPrimitiveID())
	var tsd *encoding.TSDDecoder
	if size > 0 {
		// if has compress data, create tsd decoder for merge compress
		tsd = encoding.GetTSDDecoder()
		defer encoding.ReleaseTSDDecoder(tsd)
		tsd.Reset(fs.compress)
	}
	data, freeSize, err := fs.merge(aggFunc, tsd, startTime, start, end, true)
	if err != nil {
		memDBLogger.Error("compact primitive field store data err", logger.Error(err))
	}

	fs.compress = data
	return len(fs.compress) - size - freeSize
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
	return binary.LittleEndian.Uint16(fs.buf[startOffset:])
}

// getEnd returns the delta time of start time in current write buffer
func (fs *fieldStore) getEnd() uint16 {
	return uint16(fs.buf[endOffset])
}

// merge merges the current and compress data based on primitive field aggregate function,
// startTime => current write start time
// start/end slot => target compact time slot
func (fs *fieldStore) merge(aggFunc field.AggFunc, tsd *encoding.TSDDecoder,
	startTime, startSlot, endSlot uint16, withTimeRange bool,
) (compress []byte, freeSize int, err error) {
	encode := encodeFunc(startSlot)
	for i := startSlot; i <= endSlot; i++ {
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

// Load loads the field data based on query time range, then aggregates the data
func (fs *fieldStore) Load(
	fieldType field.Type,
	agg aggregation.PrimitiveAggregator,
	memScanCtx *memScanContext,
) {
	hasOld := len(fs.compress) > 0
	aggFunc := fieldType.GetSchema().GetAggFunc(fs.GetPrimitiveID())

	var tsd *encoding.TSDDecoder
	if hasOld {
		// calc new start/end based on old compress values
		tsd = memScanCtx.tsd
		tsd.Reset(fs.compress)
	}
	startTime := fs.getStart()
	start, end := fs.slotRange(startTime)
	value := 0.0
	for i := start; i <= end; i++ {
		newValue, hasNewValue := fs.getCurrentValue(startTime, i)
		oldValue, hasOldValue := getOldFloatValue(tsd, i)

		switch {
		case hasNewValue && !hasOldValue:
			// get value from new block buffer
			value = newValue
		case hasNewValue && hasOldValue:
			// merge data from new and old
			value = aggFunc.Aggregate(newValue, oldValue)
		case !hasNewValue && hasOldValue:
			// get old value from compress data
			value = oldValue
		}
		if hasNewValue || hasOldValue {
			// aggregate the data
			if agg.Aggregate(int(i), value) {
				return
			}
		}
	}
}

// slotRange returns time slot range in current/compress buffer
func (fs *fieldStore) slotRange(currentStart uint16) (startSlot, endSlot uint16) {
	startSlot = currentStart
	endSlot = currentStart + fs.getEnd()
	if len(fs.compress) == 0 {
		return
	}
	start, end := encoding.DecodeTSDTime(fs.compress)
	return getTimeSlotRange(start, end, startSlot, endSlot)
}

// getTimeSlotRange returns the final time slot range based on start/end
func getTimeSlotRange(startSlot1, endSlot1 uint16, startSlot2, endSlot2 uint16) (start, end uint16) {
	start = startSlot1
	end = endSlot1
	if end < endSlot2 {
		end = endSlot2
	}
	if start > startSlot2 {
		start = startSlot2
	}
	return
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

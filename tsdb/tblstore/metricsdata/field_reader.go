package metricsdata

import (
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./field_reader.go -destination=./field_reader_mock.go -package metricsdata

// FieldReader represents the field reader when does metric data merge.
// !!!!NOTICE: need get field value in order by field/primitive
type FieldReader interface {
	// slotRange returns the time slot range of metric level
	slotRange() (start, end uint16)
	// getPrimitiveData returns the primitive data by field/primitive,
	// if reader is completed, return nil, if found data returns primitive data else returns nil
	getPrimitiveData(fieldID field.ID, primitiveID field.PrimitiveID) []byte
	// reset resets the field data for reading
	reset(buf []byte, position int, start, end uint16)
	// close closes the reader
	close()
}

// fieldReader implements FieldReader
type fieldReader struct {
	start, end   uint16
	buf          []byte
	fieldCount   int
	fieldOffsets *encoding.FixedOffsetDecoder

	offset    int
	ok        bool
	idx       int
	completed bool // !!!!NOTICE: need reset completed
}

// newFieldReader creates the field reader
func newFieldReader(buf []byte, position int, start, end uint16) FieldReader {
	r := &fieldReader{}
	r.reset(buf, position, start, end)
	return r
}

// reset resets the field data for reading
func (r *fieldReader) reset(buf []byte, position int, start, end uint16) {
	r.start = start
	r.end = end
	r.buf = buf
	r.fieldCount = int(stream.ReadUint16(buf, position))
	r.fieldOffsets = encoding.NewFixedOffsetDecoder(buf[position+2:])
	r.offset = 0
	r.ok = false
	r.idx = 0
	r.completed = false
}

// slotRange returns the time slot range of metric level
func (r *fieldReader) slotRange() (start, end uint16) {
	return r.start, r.end
}

// getPrimitiveData returns the primitive data by field/primitive,
// if reader is completed, return nil, if found data returns primitive data else returns nil
func (r *fieldReader) getPrimitiveData(fieldID field.ID, primitiveID field.PrimitiveID) []byte {
	if r.completed {
		return nil
	}
	if !r.ok {
		if !r.nextField() {
			return nil
		}
	}
	fID1 := r.buf[r.offset]
	fID2 := byte(fieldID)
	if fID1 < fID2 {
		if !r.nextField() {
			return nil
		}
		fID1 = r.buf[r.offset]
	}
	pID1 := r.buf[r.offset+1]
	pID2 := byte(primitiveID)
	if pID1 < pID2 {
		if !r.nextField() {
			return nil
		}
		pID1 = r.buf[r.offset+1]
	}
	if fID1 == fID2 && pID1 == pID2 {
		return r.buf[r.offset+2:]
	}
	return nil
}

// nextField goto next valid field, if on data return false and marks completed
func (r *fieldReader) nextField() bool {
	if r.idx >= r.fieldCount {
		r.completed = true
		return false
	}
	r.offset, r.ok = r.fieldOffsets.Get(r.idx)
	r.idx++
	return true
}

// close marks the reader completed
func (r *fieldReader) close() {
	r.completed = true
}

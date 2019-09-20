package series

import (
	"math"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////////
// binaryFieldIterator implements FieldIterator
//////////////////////////////////////////////////////
type binaryFieldIterator struct {
	field  field.Meta
	reader *stream.Reader

	segmentStartTime int64
}

// NewFieldIterator create field iterator based on binary data
func NewFieldIterator(fieldName string, data []byte) FieldIterator {
	it := &binaryFieldIterator{
		field:  field.Meta{Name: fieldName},
		reader: stream.NewReader(data),
	}
	it.segmentStartTime = it.reader.ReadVarint64()
	return it
}

func (it *binaryFieldIterator) FieldMeta() field.Meta   { return it.field }
func (it *binaryFieldIterator) SegmentStartTime() int64 { return it.segmentStartTime }
func (it *binaryFieldIterator) HasNext() bool           { return !it.reader.Empty() }

func (it *binaryFieldIterator) Next() PrimitiveIterator {
	fieldID := it.reader.ReadUint16()
	aggType := field.AggType(it.reader.ReadByte())
	length := it.reader.ReadVarint32()
	data := it.reader.ReadBytes(int(length))

	return NewPrimitiveIterator(fieldID, aggType, data)
}

func (it *binaryFieldIterator) Bytes() ([]byte, error) {
	//FIXME stone1100
	return nil, nil
}

//////////////////////////////////////////////////////
// primitiveIterator implements PrimitiveIterator
//////////////////////////////////////////////////////
type primitiveIterator struct {
	fieldID uint16
	aggType field.AggType
	tsd     *encoding.TSDDecoder
}

func NewPrimitiveIterator(fieldID uint16, aggType field.AggType, data []byte) PrimitiveIterator {
	return &primitiveIterator{
		fieldID: fieldID,
		aggType: aggType,
		tsd:     encoding.NewTSDDecoder(data),
	}
}

func (pi *primitiveIterator) FieldID() uint16 {
	return pi.fieldID
}

func (pi *primitiveIterator) AggType() field.AggType {
	return pi.aggType
}

func (pi *primitiveIterator) HasNext() bool {
	if pi.tsd.Error() != nil {
		return false
	}
	for pi.tsd.Next() {
		if pi.tsd.HasValue() {
			return true
		}
	}
	return false
}

func (pi *primitiveIterator) Next() (timeSlot int, value float64) {
	timeSlot = pi.tsd.Slot()
	val := pi.tsd.Value()
	value = math.Float64frombits(val)
	return
}

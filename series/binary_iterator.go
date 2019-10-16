package series

import (
	"fmt"
	"math"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////////
// binaryGroupedIterator implements GroupedIterator
//////////////////////////////////////////////////////
type binaryGroupedIterator struct {
	tags       map[string]string
	fields     map[string][]byte
	fieldNames []string

	it *BinaryIterator

	idx int
}

func NewGroupedIterator(tags map[string]string, fields map[string][]byte) GroupedIterator {
	it := &binaryGroupedIterator{tags: tags, fields: fields}
	for fieldName := range fields {
		it.fieldNames = append(it.fieldNames, fieldName)
	}
	return it
}

func (g *binaryGroupedIterator) Tags() map[string]string {
	return g.tags
}

func (g *binaryGroupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldNames) {
		return false
	}
	g.idx++
	return true
}

func (g *binaryGroupedIterator) Next() Iterator {
	fieldName := g.fieldNames[g.idx-1]
	if g.it == nil {
		g.it = NewIterator(fieldName, g.fields[fieldName])
	} else {
		g.it.Reset(fieldName, g.fields[fieldName])
	}
	return g.it
}

//////////////////////////////////////////////////////
// BinaryIterator implements Iterator
//////////////////////////////////////////////////////
type BinaryIterator struct {
	fieldName string
	fieldType field.Type
	reader    *stream.Reader
	fieldIt   *BinaryFieldIterator
}

func NewIterator(fieldName string, data []byte) *BinaryIterator {
	it := &BinaryIterator{fieldName: fieldName, reader: stream.NewReader(data)}
	it.fieldType = field.Type(it.reader.ReadByte())
	return it
}

func (b *BinaryIterator) Reset(fieldName string, data []byte) {
	b.fieldName = fieldName
	b.reader.Reset(data)
	b.fieldType = field.Type(b.reader.ReadByte())
}

func (b *BinaryIterator) FieldName() string {
	return b.fieldName
}

func (b *BinaryIterator) FieldType() field.Type {
	return b.fieldType
}

func (b *BinaryIterator) HasNext() bool {
	return !b.reader.Empty()
}

func (b *BinaryIterator) Next() (startTime int64, fieldIt FieldIterator) {
	startTime = b.reader.ReadVarint64()
	length := b.reader.ReadVarint32()
	if length == 0 {
		return
	}
	data := b.reader.ReadBytes(int(length))
	if b.fieldIt == nil {
		b.fieldIt = NewFieldIterator(data)
	} else {
		b.fieldIt.reset(data)
	}
	fieldIt = b.fieldIt
	return
}

func (b *BinaryIterator) MarshalBinary() ([]byte, error) {
	return MarshalIterator(b)
}

//////////////////////////////////////////////////////
// binaryFieldIterator implements FieldIterator
//////////////////////////////////////////////////////
type BinaryFieldIterator struct {
	reader *stream.Reader
	pIt    *BinaryPrimitiveIterator
}

// NewFieldIterator create field iterator based on binary data
func NewFieldIterator(data []byte) *BinaryFieldIterator {
	it := &BinaryFieldIterator{
		reader: stream.NewReader(data),
	}
	return it
}

func (it *BinaryFieldIterator) reset(data []byte) {
	it.reader.Reset(data)
}

func (it *BinaryFieldIterator) HasNext() bool { return !it.reader.Empty() }

func (it *BinaryFieldIterator) Next() PrimitiveIterator {
	fieldID := it.reader.ReadUint16()
	aggType := field.AggType(it.reader.ReadByte())
	length := it.reader.ReadVarint32()
	data := it.reader.ReadBytes(int(length))

	if it.pIt == nil {
		it.pIt = NewPrimitiveIterator(fieldID, aggType, encoding.NewTSDDecoder(data))
	} else {
		it.pIt.Reset(fieldID, aggType, data)
	}
	return it.pIt
}

func (it *BinaryFieldIterator) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("not support")
}

//////////////////////////////////////////////////////
// primitiveIterator implements PrimitiveIterator
//////////////////////////////////////////////////////
type BinaryPrimitiveIterator struct {
	fieldID uint16
	aggType field.AggType
	tsd     *encoding.TSDDecoder
}

func NewPrimitiveIterator(fieldID uint16, aggType field.AggType, tsd *encoding.TSDDecoder) *BinaryPrimitiveIterator {
	return &BinaryPrimitiveIterator{
		fieldID: fieldID,
		aggType: aggType,
		tsd:     tsd,
	}
}

func (pi *BinaryPrimitiveIterator) Reset(fieldID uint16, aggType field.AggType, data []byte) {
	pi.fieldID = fieldID
	pi.aggType = aggType
	pi.tsd.Reset(data)
}

func (pi *BinaryPrimitiveIterator) FieldID() uint16 {
	return pi.fieldID
}

func (pi *BinaryPrimitiveIterator) AggType() field.AggType {
	return pi.aggType
}

func (pi *BinaryPrimitiveIterator) HasNext() bool {
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

func (pi *BinaryPrimitiveIterator) Next() (timeSlot int, value float64) {
	timeSlot = pi.tsd.Slot()
	val := pi.tsd.Value()
	value = math.Float64frombits(val)
	return
}

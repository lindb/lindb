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
	tags       string
	fields     map[field.Name][]byte
	fieldNames []field.Name

	it *BinaryIterator

	idx int
}

func NewGroupedIterator(tags string, fields map[field.Name][]byte) GroupedIterator {
	it := &binaryGroupedIterator{tags: tags, fields: fields}
	for fieldName := range fields {
		it.fieldNames = append(it.fieldNames, fieldName)
	}
	return it
}

func (g *binaryGroupedIterator) Tags() string {
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
	fieldName field.Name
	fieldType field.Type
	reader    *stream.Reader
	fieldIt   *BinaryFieldIterator
	data      []byte
}

func NewIterator(fieldName field.Name, data []byte) *BinaryIterator {
	it := &BinaryIterator{fieldName: fieldName, reader: stream.NewReader(data), data: data}
	it.fieldType = field.Type(it.reader.ReadByte())
	return it
}

func (b *BinaryIterator) Reset(fieldName field.Name, data []byte) {
	b.fieldName = fieldName
	b.reader.Reset(data)
	b.fieldType = field.Type(b.reader.ReadByte())
}

func (b *BinaryIterator) FieldName() field.Name {
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
	aggType := field.AggType(b.reader.ReadByte())
	length := b.reader.ReadVarint32()
	if length <= 0 {
		return
	}

	data := b.reader.ReadBytes(int(length))
	if b.fieldIt == nil {
		b.fieldIt = NewFieldIterator(aggType, encoding.NewTSDDecoder(data))
	} else {
		b.fieldIt.reset(aggType, data)
	}
	fieldIt = b.fieldIt
	return
}

func (b *BinaryIterator) MarshalBinary() ([]byte, error) {
	return b.data, nil
}

//////////////////////////////////////////////////////
// binaryFieldIterator implements FieldIterator
//////////////////////////////////////////////////////
type BinaryFieldIterator struct {
	aggType field.AggType
	tsd     *encoding.TSDDecoder
}

// NewFieldIterator create field iterator based on binary data
func NewFieldIterator(aggType field.AggType, tsd *encoding.TSDDecoder) *BinaryFieldIterator {
	it := &BinaryFieldIterator{
		aggType: aggType,
		tsd:     tsd,
	}
	return it
}

func (it *BinaryFieldIterator) reset(aggType field.AggType, data []byte) {
	it.aggType = aggType
	it.tsd.Reset(data)
}

func (it *BinaryFieldIterator) AggType() field.AggType {
	return it.aggType
}

func (it *BinaryFieldIterator) HasNext() bool {
	if it.tsd.Error() != nil {
		return false
	}
	for it.tsd.Next() {
		if it.tsd.HasValue() {
			return true
		}
	}
	return false
}

func (it *BinaryFieldIterator) Next() (timeSlot int, value float64) {
	timeSlot = int(it.tsd.Slot())
	val := it.tsd.Value()
	value = math.Float64frombits(val)
	return
}

func (it *BinaryFieldIterator) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("not support")
}

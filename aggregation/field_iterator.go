package aggregation

import (
	"math"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

type fieldIterator struct {
	startSlot int

	length int
	idx    int
	its    []series.PrimitiveIterator
}

func newFieldIterator(startSlot int,
	its []series.PrimitiveIterator) series.FieldIterator {
	return &fieldIterator{
		startSlot: startSlot,
		its:       its,
		length:    len(its),
	}
}

func (it *fieldIterator) HasNext() bool {
	return it.idx < it.length
}

func (it *fieldIterator) Next() series.PrimitiveIterator {
	if it.idx >= it.length {
		return nil
	}
	primitiveIt := it.its[it.idx]
	it.idx++
	return primitiveIt
}

func (it *fieldIterator) Bytes() ([]byte, error) {
	if it.length == 0 {
		return nil, nil
	}
	//need reset idx
	it.idx = 0
	writer := stream.NewBufferWriter(nil)
	for it.HasNext() {
		primitiveIt := it.Next()
		encoder := encoding.NewTSDEncoder(it.startSlot)
		idx := 0
		for primitiveIt.HasNext() {
			slot, value := primitiveIt.Next()
			for slot > idx {
				encoder.AppendTime(bit.Zero)
				idx++
			}
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(value))
			idx++
		}
		data, err := encoder.Bytes()
		if err != nil {
			return nil, err
		}
		writer.PutUInt16(primitiveIt.FieldID())
		writer.PutByte(byte(primitiveIt.AggType()))
		writer.PutVarint32(int32(len(data)))
		writer.PutBytes(data)
	}
	return writer.Bytes()
}

// primitiveIterator represents primitive iterator using array
type primitiveIterator struct {
	start   int
	id      uint16
	aggType field.AggType
	it      collections.FloatArrayIterator
}

// newPrimitiveIterator create primitive iterator using array
func newPrimitiveIterator(id uint16, start int, aggType field.AggType, values collections.FloatArray) series.PrimitiveIterator {
	it := &primitiveIterator{
		start:   start,
		id:      id,
		aggType: aggType,
	}
	if values != nil {
		it.it = values.Iterator()
	}
	return it
}

// ID returns the primitive field id
func (it *primitiveIterator) FieldID() uint16 {
	return it.id
}

// AggType returns the primitive field's agg type
func (it *primitiveIterator) AggType() field.AggType {
	return it.aggType
}

// HasNext returns if the iteration has more data points
func (it *primitiveIterator) HasNext() bool {
	if it.it == nil {
		return false
	}
	return it.it.HasNext()
}

// Next returns the data point in the iteration
func (it *primitiveIterator) Next() (timeSlot int, value float64) {
	if it.it == nil {
		return -1, 0
	}
	timeSlot, value = it.it.Next()
	if timeSlot == -1 {
		return
	}
	timeSlot += it.start
	return
}

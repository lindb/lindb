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

// fieldIterator implements series.FieldIterator interface
type fieldIterator struct {
	startSlot int
	aggType   field.AggType
	it        collections.FloatArrayIterator
}

// newFieldIterator creates a field iterator
func newFieldIterator(startSlot int, aggType field.AggType, values collections.FloatArray) series.FieldIterator {
	it := &fieldIterator{
		startSlot: startSlot,
		aggType:   aggType,
	}
	if values != nil {
		it.it = values.Iterator()
	}
	return it
}

func (it *fieldIterator) AggType() field.AggType {
	return it.aggType
}

// HasNext returns if the iteration has more fields
func (it *fieldIterator) HasNext() bool {
	if it.it == nil {
		return false
	}
	return it.it.HasNext()
}

// Next returns the data point in the iteration
func (it *fieldIterator) Next() (timeSlot int, value float64) {
	if it.it == nil {
		return -1, 0
	}
	timeSlot, value = it.it.Next()
	if timeSlot == -1 {
		return -1, 0
	}
	timeSlot += it.startSlot
	return
}

// MarshalBinary marshals the data
func (it *fieldIterator) MarshalBinary() ([]byte, error) {
	if it.it == nil {
		return nil, nil
	}
	//FIXME reuse encoder???
	encoder := encoding.TSDEncodeFunc(uint16(it.startSlot))
	idx := it.startSlot
	writer := stream.NewBufferWriter(nil)
	for it.HasNext() {
		slot, value := it.Next()
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
	if idx == it.startSlot {
		// maybe field data already read
		return nil, nil
	}
	writer.PutByte(byte(it.AggType()))   // agg type
	writer.PutVarint32(int32(len(data))) // length of field data
	writer.PutBytes(data)                // field data
	return writer.Bytes()
}

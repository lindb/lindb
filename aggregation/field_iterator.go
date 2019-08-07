package aggregation

import (
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/field"
)

type fieldIterator struct {
	id        uint16
	fieldType field.Type

	length int
	idx    int
	its    []field.PrimitiveIterator
}

func newFieldIterator(id uint16, fieldType field.Type, its []field.PrimitiveIterator) field.Iterator {
	return &fieldIterator{
		id:        id,
		fieldType: fieldType,
		its:       its,
		length:    len(its),
	}
}

func (it *fieldIterator) Name() string {
	//TODO need impl
	return ""
}

func (it *fieldIterator) ID() uint16 {
	return it.id
}

func (it *fieldIterator) FieldType() field.Type {
	return it.fieldType
}

func (it *fieldIterator) HasNext() bool {
	return it.idx < it.length
}

func (it *fieldIterator) Next() field.PrimitiveIterator {
	if it.idx >= it.length {
		return nil
	}
	primitiveIt := it.its[it.idx]
	it.idx++
	return primitiveIt
}

// primitiveIterator represents primitive iterator using array
type primitiveIterator struct {
	id uint16
	it collections.FloatArrayIterator
}

// newPrimitiveIterator create primitive iterator using array
func newPrimitiveIterator(id uint16, values collections.FloatArray) field.PrimitiveIterator {
	it := &primitiveIterator{
		id: id,
	}
	if values != nil {
		it.it = values.Iterator()
	}
	return it
}

// ID returns the primitive field id
func (it *primitiveIterator) ID() uint16 {
	return it.id
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
	return it.it.Next()
}

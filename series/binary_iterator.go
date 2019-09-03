package series

import (
	"math"

	"github.com/lindb/lindb/pkg/encoding"
)

//////////////////////////////////////////////////////
// primitiveIterator implements PrimitiveIterator
//////////////////////////////////////////////////////
type primitiveIterator struct {
	fieldID uint16
	tsd     *encoding.TSDDecoder
}

func NewPrimitiveIterator(fieldID uint16, data []byte) PrimitiveIterator {
	return &primitiveIterator{
		fieldID: fieldID,
		tsd:     encoding.NewTSDDecoder(data),
	}
}

func (pi *primitiveIterator) FieldID() uint16 {
	return pi.fieldID
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

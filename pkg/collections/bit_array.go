package collections

import (
	"fmt"
	"math"
	"strings"
)

// BitArray is a simple struct for maintaining state of a bit array,
// which is useful for tracking bool type values efficiently.
// Not thread-safe.
type BitArray struct {
	payload []byte
	builder strings.Builder
}

// NewBitArray returns a new BitArray from buf.
func NewBitArray(buf []byte) (*BitArray, error) {
	if len(buf) > math.MaxUint16 {
		return nil, fmt.Errorf("%d is too long for a bit array", len(buf))
	}
	return &BitArray{payload: buf}, nil
}

// Reset resets all payload to zero.
func (ba *BitArray) Reset() {
	ba.payload = ba.payload[:0]
	ba.builder.Reset()
}

// SetBit sets a bit at the given index.
func (ba *BitArray) SetBit(k uint16) {
	for int(math.Ceil(float64(k+1)/float64(8))) > ba.Len() {
		ba.payload = append(ba.payload, 0)
	}
	idx := int(k / 8)
	offset := k % 8

	ba.payload[idx] |= 1 << offset
}

// Bytes return the underlying bytes slice
func (ba *BitArray) Bytes() []byte {
	return ba.payload
}

// GetBit returns a bool which indicates given index has been set before.
func (ba *BitArray) GetBit(k uint16) bool {
	if int(k) >= ba.Len()*8 {
		return false
	}
	idx := int(k / 8)
	offset := k % 8

	return ba.payload[idx]&(1<<offset) != 0
}

// Len returns the length of the bit-array.
func (ba *BitArray) Len() int {
	return len(ba.payload)
}

// String implements stringer.
// Inefficient function, just for test and debug
func (ba *BitArray) String() string {
	ba.builder.Reset()
	for _, val := range ba.payload {
		section := []byte(fmt.Sprintf("%08b", val))
		for i := 0; i < len(section); i++ {
			ba.builder.WriteByte(section[len(section)-i-1])
		}
	}
	return ba.builder.String()
}

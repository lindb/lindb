package metrictbl

import (
	"fmt"
	"math"
	"strings"
)

// bitArray is a simple struct for maintaining state of a bit array,
// which is useful for tracking bool type values efficiently.
// Not thread-safe.
type bitArray struct {
	payload []byte
}

// newBitArray returns a new bitArray from buf.
func newBitArray(buf []byte) (*bitArray, error) {
	if len(buf) > math.MaxUint16 {
		return nil, fmt.Errorf("%d is too long for a bit array", len(buf))
	}
	return &bitArray{payload: buf}, nil
}

// reset resets all payload to zero.
func (ba *bitArray) reset() {
	ba.payload = ba.payload[:0]
}

// setBit sets a bit at the given index.
func (ba *bitArray) setBit(k uint16) {
	for int(math.Ceil(float64(k+1)/float64(8))) > ba.getLen() {
		ba.payload = append(ba.payload, 0)
	}
	idx := int(k / 8)
	offset := k % 8

	ba.payload[idx] |= 1 << offset
}

// getBit returns if a bool indicating if given index has been set before.
func (ba *bitArray) getBit(k uint16) bool {
	if int(k) >= ba.getLen()*8 {
		return false
	}
	idx := int(k / 8)
	offset := k % 8

	return ba.payload[idx]&(1<<offset) != 0
}

// getLen returns the length of the bit-array.
func (ba *bitArray) getLen() int {
	return len(ba.payload)
}

// String implements stringer.
// Inefficient function, just for test and debug
func (ba *bitArray) String() string {
	var b strings.Builder
	for _, val := range ba.payload {
		section := []byte(fmt.Sprintf("%08b", val))
		for i := 0; i < len(section); i++ {
			b.WriteByte(section[len(section)-i-1])
		}
	}
	return b.String()
}

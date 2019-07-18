package util

import (
	"encoding/binary"
	"math"
)

//ShortToInt represents convert a 2-short  value to  an uint32
func ShortToInt(high uint16, low uint16) uint32 {
	return uint32(high)<<16 + uint32(low)
}

//IntToShort represents convert an uint32 value to a 2-short
func IntToShort(value uint32) (high, low uint16) {
	high = uint16(value >> 16 & math.MaxUint16)
	low = uint16(value & math.MaxUint16)
	return
}

//Uint32ToBytes represents converts 32-bits positive integer to 4-bytes
func Uint32ToBytes(v uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, v)
	return bs
}

//BytesToUint32 represents converts bytes to 32-bits positive integer
func BytesToUint32(bs []byte) uint32 {
	return binary.LittleEndian.Uint32(bs)
}

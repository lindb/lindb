package encoding

const maxLowBit = 0xFFFF

// ZigZagEncode converts a int64 to a uint64 by zig zagging negative and positive values
// across even and odd numbers.  Eg. [0,-1,1,-2] becomes [0, 1, 2, 3].
func ZigZagEncode(x int64) uint64 {
	return uint64(x<<1) ^ uint64(x>>63)
}

// ZigZagDecode converts a previously zigzag encoded uint64 back to a int64.
func ZigZagDecode(v uint64) int64 {
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63))
}

// HighBits returns the high 16 bits of value
func HighBits(x uint32) uint16 {
	return uint16(x >> 16)
}

// LowBits returns the low 16 bits of value
func LowBits(x uint32) uint16 {
	return uint16(x & maxLowBit)
}

// ValueWithHighLowBits returns the value with high/low 16 bits
func ValueWithHighLowBits(high uint32, low uint16) uint32 {
	return uint32(low&maxLowBit) | high
}

// GetMinLength returns the min length of value
func GetMinLength(value int) int {
	switch {
	case value < 1<<8:
		return 1
	case value < 1<<16:
		return 2
	case value < 1<<24:
		return 3
	default:
		return 4
	}
}

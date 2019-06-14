package encoding

// ZigZagEncode converts a int64 to a uint64 by zig zagging negative and positive values
// across even and odd numbers.  Eg. [0,-1,1,-2] becomes [0, 1, 2, 3].
func ZigZagEncode(x int64) uint64 {
	return uint64(uint64(x<<1) ^ uint64(int64(x) >> 63))
}

// ZigZagDecode converts a previously zigzag encoded uint64 back to a int64.
func ZigZagDecode(v uint64) int64 {
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63))
}

// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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

// Uint32MinWidth returns the min length of uint32
func Uint32MinWidth(value uint32) int {
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

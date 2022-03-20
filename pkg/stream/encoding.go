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

package stream

import "encoding/binary"

// UvariantSize returns the bytes-size of a uint64 uvariant encoded number.
func UvariantSize(value uint64) int {
	i := 0
	for value >= 0x80 {
		value >>= 7
		i++
	}
	return i + 1
}

// VariantSize returns the bytes-size of a int64 variant encoded number.
func VariantSize(value int64) int {
	ux := uint64(value) << 1
	if value < 0 {
		ux = ^ux
	}
	return UvariantSize(ux)
}

// PutUvariantLittleEndian encodes a uint64 into buf and returns the number of bytes written.
// The default PutUvarint use big-endian
func PutUvariantLittleEndian(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	for x := 0; x < (i+1)/2; x++ {
		buf[x], buf[i-x] = buf[i-x], buf[x]
	}
	return i + 1
}

// UvarintLittleEndian decodes a uint64 from buf's tail with little endian
// and returns that value and the number of bytes read (> 0).
// If an error occurred, the value is 0 and the number of bytes n is <= 0 meaning:
//
// 	n == 0: buf too small
// 	n  < 0: value larger than 64 bits (overflow)
// 	        and -n is the number of bytes read
func UvarintLittleEndian(buf []byte) (x uint64, s int) {
	for cursor := len(buf) - 1; cursor >= 0; cursor-- {
		b := buf[cursor]
		i := len(buf) - 1 - cursor
		if b < 0x80 {
			if i >= binary.MaxVarintLen64 || i == binary.MaxVarintLen64-1 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

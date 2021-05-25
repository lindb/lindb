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

import (
	"encoding/binary"
)

// PutUint64 puts uint64 into buf with offset
func PutUint64(buf []byte, offset int, value uint64) {
	binary.LittleEndian.PutUint64(buf[offset:], value)
}

// ReadUint64 reads 8 bytes from buf as uint64
func ReadUint64(buf []byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(buf[offset : offset+8])
}

// PutUint32 puts uint32 into buf with offset
func PutUint32(buf []byte, offset int, value uint32) {
	binary.LittleEndian.PutUint32(buf[offset:], value)
}

// ReadUint32 reads 4 bytes from buf as uint32
func ReadUint32(buf []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(buf[offset : offset+4])
}

// PutUint16 reads 2 bytes from buf as uint16
func PutUint16(buf []byte, offset int, value uint16) {
	binary.LittleEndian.PutUint16(buf[offset:], value)
}

// ReadUint16 reads 2 bytes from buf as uint16
func ReadUint16(buf []byte, offset int) uint16 {
	return binary.LittleEndian.Uint16(buf[offset : offset+2])
}

// ReadUvarint reads an encoded unsigned integer from bytes.Reader and returns it as a uint64.
func ReadUvarint(buf []byte, offset int) (value uint64, readLen int, err error) {
	var s uint
	for i := 0; ; i++ {
		b := buf[offset]
		offset++
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return value, i + 1, errOverflow
			}
			return value | uint64(b)<<s, i + 1, nil
		}
		value |= uint64(b&0x7f) << s
		s += 7
	}
}

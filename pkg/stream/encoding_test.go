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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UvariantLittleEndian(t *testing.T) {
	for _, i := range []int64{
		-1 << 63,
		-1<<63 + 1,
		-1,
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		1<<63 - 1,
	} {
		assertValueEqual(t, uint64(i), UvariantSize(uint64(i)))
	}
}

func Test_UvarintLittleEndian_special_cases(t *testing.T) {
	// nil buffer
	value, size := UvarintLittleEndian(nil)
	assert.Zero(t, value)
	assert.Zero(t, size)

	// overflow
	var buf2 = []byte{
		1, 1,
		0x80, 0x80, 0x80, 0x80, 0x80,
		0x80, 0x80, 0x80, 0x80, 0x80,
	}
	value, size = UvarintLittleEndian(buf2)
	assert.Zero(t, value)
	assert.Equal(t, -11, size)
}

func assertValueEqual(t *testing.T, value uint64, size int) {
	var buf [binary.MaxVarintLen64]byte
	realSize := PutUvariantLittleEndian(buf[:], value)
	assert.Equal(t, size, realSize)

	// put into the tail
	decodedValue, decodedSize := UvarintLittleEndian(buf[:realSize])
	assert.Equal(t, decodedSize, realSize)
	assert.Equal(t, decodedValue, value)
}

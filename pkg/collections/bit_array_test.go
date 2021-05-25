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

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BitArray(t *testing.T) {
	ba := NewBitArray(nil)
	assert.Equal(t, "", ba.String())
	assert.False(t, ba.GetBit(0))

	ba.SetBit(uint16(0))
	ba.SetBit(uint16(2))
	assert.Len(t, ba.Bytes(), 1)

	assert.Equal(t, "10100000", ba.String())
	ba.SetBit(uint16(4))
	assert.Equal(t, "10101000", ba.String())

	ba.SetBit(uint16(8))
	assert.Equal(t, "1010100010000000", ba.String())
	ba.SetBit(uint16(9))
	ba.SetBit(uint16(9))
	assert.Equal(t, "1010100011000000", ba.String())
	ba.SetBit(uint16(16))
	assert.Equal(t, "101010001100000010000000", ba.String())

	assert.True(t, ba.GetBit(0))
	assert.False(t, ba.GetBit(1))
	assert.True(t, ba.GetBit(8))
	assert.True(t, ba.GetBit(9))
	assert.False(t, ba.GetBit(23))
	assert.False(t, ba.GetBit(24))
	assert.False(t, ba.GetBit(800))

	ba.Reset(nil)
	assert.False(t, ba.GetBit(0))

	ba2 := NewBitArray(nil)
	ba2.Reset([]byte{255, 255})
	assert.NotNil(t, ba2)
	assert.True(t, ba2.GetBit(0))

	ba3 := NewBitArray([]byte{})
	assert.NotNil(t, ba3)
	assert.False(t, ba3.GetBit(23))

}

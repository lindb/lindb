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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZigZag(t *testing.T) {
	v := ZigZagEncode(1)
	assert.Equal(t, int64(1), ZigZagDecode(v))

	v = ZigZagEncode(-99999)
	assert.Equal(t, int64(-99999), ZigZagDecode(v))
}

func TestHighLowBits(t *testing.T) {
	v := uint32(67043434)
	high := HighBits(v)
	low := LowBits(v)
	assert.Equal(t, v, ValueWithHighLowBits(uint32(high)<<16, low))
}

func TestUint32MinWidth(t *testing.T) {
	assert.Equal(t, 1, Uint32MinWidth(0))
	assert.Equal(t, 1, Uint32MinWidth(1))
	assert.Equal(t, 2, Uint32MinWidth(1<<8))
	assert.Equal(t, 2, Uint32MinWidth((1<<8)+1))
	assert.Equal(t, 3, Uint32MinWidth(1<<16))
	assert.Equal(t, 3, Uint32MinWidth((1<<16)+1))
	assert.Equal(t, 4, Uint32MinWidth(1<<24))
	assert.Equal(t, 4, Uint32MinWidth((1<<24)+1))
}

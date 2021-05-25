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

package bufpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BufferPool(t *testing.T) {
	buf1 := GetBuffer()
	assert.Equal(t, 0, buf1.Len())

	buf1.WriteByte(byte(66))
	assert.Equal(t, 1, buf1.Len())
	assert.Equal(t, 64, buf1.Cap())

	PutBuffer(buf1)

	buf2 := GetBuffer()
	assert.Equal(t, 0, buf2.Len())
}

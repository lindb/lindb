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

func Test_MustCopy(t *testing.T) {
	dst1 := make([]byte, 100)
	src1 := make([]byte, 50)
	assert.Len(t, dst1, 100)
	assert.Equal(t, 100, cap(dst1))
	assert.Len(t, src1, 50)
	assert.Equal(t, 50, cap(src1))

	dst1Copied := MustCopy(dst1, src1)
	assert.Len(t, dst1Copied, 50)
	assert.Equal(t, 100, cap(dst1Copied))

	dst1 = dst1[:0]
	dst1Copied = MustCopy(dst1, src1)
	assert.Len(t, dst1Copied, 50)
	assert.Equal(t, 100, cap(dst1Copied))

	src2 := make([]byte, 120)
	dst1Copied = MustCopy(dst1, src2)
	assert.Len(t, dst1Copied, 120)
	assert.Equal(t, 200, cap(dst1Copied))
}

func Test_growSlice(t *testing.T) {
	dst1 := make([]byte, 100)
	assert.Equal(t, 100, cap(growSlice(dst1, 99)))
	assert.Equal(t, 200, cap(growSlice(dst1, 101)))
	assert.Equal(t, 1025, cap(growSlice(dst1, 1025)))

	dst2 := make([]byte, 1024)
	assert.Equal(t, 1280, cap(growSlice(dst2, 1025)))
	assert.Equal(t, 1290, cap(growSlice(dst2, 1290)))
}

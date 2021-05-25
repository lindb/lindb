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

package bit

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOkWriter struct{}

func (w *mockOkWriter) Write(p []byte) (n int, err error) {
	return 1, nil
}

type mockErrWriter struct{}

func (w *mockErrWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error")
}

var (
	okWriter  = NewWriter(&mockOkWriter{})
	badWriter = NewWriter(&mockErrWriter{})
)

func Test_Writer_WriteBit(t *testing.T) {
	for range [10]struct{}{} {
		assert.Nil(t, okWriter.WriteBit(Zero))
		assert.Nil(t, okWriter.WriteBit(One))
	}

	for range [7]struct{}{} {
		assert.Nil(t, badWriter.WriteBit(Zero))
	}
	assert.NotNil(t, badWriter.WriteBit(One))
}

func Test_Writer_WriteBytes(t *testing.T) {
	assert.Nil(t, okWriter.WriteBits(math.MaxUint64, 63))

	assert.NotNil(t, badWriter.WriteBits(math.MaxUint64, 63))
}

func Test_Writer_Flush(t *testing.T) {
	okWriter.count = 8
	assert.Nil(t, okWriter.Flush())

	badWriter.count = 1
	assert.NotNil(t, badWriter.Flush())
	assert.NotNil(t, badWriter.Flush())
}

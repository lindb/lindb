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

package common

import (
	"github.com/stretchr/testify/assert"

	"io"
	"testing"
)

type mockReader struct {
	n int
}

func (mr *mockReader) Read(p []byte) (n int, err error) {
	defer func() { mr.n++ }()

	switch mr.n {
	case 2:
		return 2, io.EOF
	default:
		p = append(p[:0], []byte{0x1f, 0x8b, 8, 0, 1, 2, 3, 4, 5, 6}...)
		_ = p
		return 20, nil
	}
}

func Test_GetGzipReader(t *testing.T) {
	PutGzipReader(nil)

	mr := &mockReader{}
	r, err := GetGzipReader(mr)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	PutGzipReader(r)

	r, err = GetGzipReader(mr)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	PutGzipReader(r)

	r, err = GetGzipReader(mr)
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

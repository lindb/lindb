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

package compress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Snappy(t *testing.T) {
	data := []byte("hello snappy")
	w := NewSnappyWriter()
	_, err := w.Write(data)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	compress := w.Bytes()

	r := NewSnappyReader()
	dst, err := r.Uncompress(compress)
	assert.NoError(t, err)
	assert.Equal(t, data, dst)
}

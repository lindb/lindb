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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bufioutil"
)

func Test_Reader(t *testing.T) {
	var data []byte
	buf := bufioutil.NewBuffer(data)
	reader := NewReader(buf)

	_, err := reader.ReadBit()
	assert.NotNil(t, err)
	_, err = reader.ReadByte()
	assert.NotNil(t, err)
	_, err = reader.ReadBits(10)
	assert.NotNil(t, err)
	_, err = reader.ReadBits(1)
	assert.NotNil(t, err)

	data = append(data, []byte{1, 2, 3, 4, 5, 6, 7, 8}...)
	buf.SetBuf(data)
	reader.Reset()
	_, err = reader.ReadBits(10)
	assert.Nil(t, err)
}

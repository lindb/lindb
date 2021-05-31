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

	"strings"
	"testing"
)

const influxText = `
# comment
 a1,location=us-midwest temperature=82 1465839830100400200

 a2,location=us-midwest temperature=1 1465839830100400200

a3,location=us-midwest temperature=100 1465839830100400200`

func assertReadAll(t *testing.T, cr *ChunkReader) {
	assert.True(t, cr.HasNext())
	assert.Equal(t, "# comment", string(cr.Next()))
	assert.Nil(t, cr.Error())

	assert.True(t, cr.HasNext())
	assert.Equal(t, "a1,location=us-midwest temperature=82 1465839830100400200", string(cr.Next()))

	assert.True(t, cr.HasNext())
	assert.Equal(t, "a2,location=us-midwest temperature=1 1465839830100400200", string(cr.Next()))

	assert.True(t, cr.HasNext())
	assert.Equal(t, "a3,location=us-midwest temperature=100 1465839830100400200", string(cr.Next()))

	assert.False(t, cr.HasNext())
	assert.NotNil(t, cr.Error())
}

func Test_ChunkReader(t *testing.T) {
	assertReadAll(t, newChunkReader(strings.NewReader(influxText)))

	assertReadAll(t, newChunkReaderWithSize(strings.NewReader(influxText), 64))
	assertReadAll(t, newChunkReaderWithSize(strings.NewReader(influxText), 128))
	assertReadAll(t, newChunkReaderWithSize(strings.NewReader(influxText), 256))
}

func Test_ChunkReader_TooLongLine(t *testing.T) {
	cr := newChunkReaderWithSize(strings.NewReader(influxText), 16)

	assert.True(t, cr.HasNext())
	assert.Equal(t, "# comment", string(cr.Next()))
	assert.Nil(t, cr.Error())

	assert.False(t, cr.HasNext())
	assert.Equal(t, "a1,location=us-", string(cr.Next()))
	assert.NotNil(t, cr.Error())
}

func Test_ChunkReaderPool(t *testing.T) {
	PutChunkReader(nil)

	cr1 := GetChunkReader(strings.NewReader(influxText))
	assertReadAll(t, cr1)
	PutChunkReader(cr1)

	cr2 := GetChunkReader(strings.NewReader(influxText))
	assertReadAll(t, cr2)
}

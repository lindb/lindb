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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/bufioutil"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	bitWriter := bit.NewWriter(&buf)
	e := NewXOREncoder(bitWriter)
	_ = e.Write(uint64(76))
	_ = e.Write(uint64(50))
	_ = e.Write(uint64(50))
	_ = e.Write(uint64(999999999))
	_ = e.Write(uint64(100))

	err := bitWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()

	reader := bit.NewReader(bufioutil.NewBuffer(data))
	d := NewXORDecoder(reader)
	exceptIntValue(d, t, uint64(76))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(50))
	exceptIntValue(d, t, uint64(999999999))
	exceptIntValue(d, t, uint64(100))
}

func exceptIntValue(d *XORDecoder, t *testing.T, except uint64) {
	assert.True(t, d.Next())
	assert.Equal(t, except, d.Value())
}

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

	"github.com/lindb/lindb/pkg/bit"
)

func TestTSDStream(t *testing.T) {
	// time range => [10,13]
	encoder := NewTSDEncoder(10)

	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(100))
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(50))

	data, err := encoder.BytesWithoutTime()
	assert.NoError(t, err)

	writer := NewTSDStreamWriter(10, 13)
	writer.WriteField(1, data)
	writer.WriteField(10, data)
	writer.WriteField(15, data)
	writer.WriteField(16, data)
	fieldsData, err := writer.Bytes()
	assert.NoError(t, err)
	start, end := DecodeTSDTime(fieldsData)
	assert.Equal(t, uint16(10), start)
	assert.Equal(t, uint16(13), end)

	reader := NewTSDStreamReader(fieldsData)
	defer reader.Close()
	start, end = reader.TimeRange()
	assert.Equal(t, uint16(10), start)
	assert.Equal(t, uint16(13), end)
	fieldIDs := []uint16{1, 10, 15, 16}
	for _, fieldID := range fieldIDs {
		assert.True(t, reader.HasNext())
		field, fieldData := reader.Next()
		assert.Equal(t, fieldID, field)
		assertFieldData(t, fieldData)
	}

	assert.False(t, reader.HasNext())
}

func assertFieldData(t *testing.T, decoder *TSDDecoder) {
	assert.Equal(t, uint16(10), decoder.StartTime())
	assert.Equal(t, uint16(13), decoder.EndTime())

	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(10), decoder.Value())
	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(100), decoder.Value())
	assert.True(t, decoder.Next())
	assert.False(t, decoder.HasValue())
	assert.True(t, decoder.Next())
	assert.True(t, decoder.HasValue())
	assert.Equal(t, uint64(50), decoder.Value())

	assert.False(t, decoder.Next())
}

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
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
)

var f = flushFunc

func TestCodec(t *testing.T) {
	encoder := GetTSDEncoder(10)
	defer ReleaseTSDEncoder(encoder)
	data, err := encoder.Bytes()
	assert.NoError(t, err)
	assert.Nil(t, data)
	encoder.Reset()
	encoder.RestWithStartTime(10)

	encoder.EmitDownSamplingValue(0, math.Float64frombits(10))
	encoder.EmitDownSamplingValue(1, math.Float64frombits(100))
	encoder.EmitDownSamplingValue(2, math.Inf(1))
	encoder.EmitDownSamplingValue(3, math.Float64frombits(50))

	data, err = encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder := NewTSDDecoder(data)
	assert.Equal(t, uint16(10), decoder.StartTime())
	assert.Equal(t, uint16(13), decoder.EndTime())
	startTime, endTime := DecodeTSDTime(data)
	assert.Equal(t, uint16(10), startTime)
	assert.Equal(t, uint16(13), endTime)

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
	assert.False(t, decoder.HasValueWithSlot(10))
	assert.False(t, decoder.HasValueWithSlot(13))

	assert.False(t, decoder.Next())
	data, err = encoder.BytesWithoutTime()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	decoder.ResetWithTimeRange(data, 10, 13)
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

	encoder.Reset()
	data, _ = encoder.Bytes()
	assert.Len(t, data, 4)

	decoder.Reset(nil)
	assert.Error(t, decoder.Error())
}

func TestTsdEncoder_Err(t *testing.T) {
	defer func() {
		flushFunc = f
	}()
	encoder := NewTSDEncoder(10)
	// case 1: encode with err
	encoder.err = fmt.Errorf("err")
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	data, err := encoder.Bytes()
	assert.Error(t, err)
	assert.Nil(t, data)
	data, err = encoder.BytesWithoutTime()
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 2: flush err
	encoder = NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	flushFunc = func(writer *bit.Writer) error {
		return fmt.Errorf("err")
	}
	data, err = encoder.Bytes()
	assert.Error(t, err)
	assert.Nil(t, data)
	data, err = encoder.BytesWithoutTime()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestHasValueWithSlot(t *testing.T) {
	encoder := NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(100))
	encoder.AppendTime(bit.Zero)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(50))

	data, err := encoder.Bytes()
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	// seek test
	decoder0 := NewTSDDecoder(data)
	assert.False(t, decoder0.Seek(15))
	assert.True(t, decoder0.Seek(12))
	assert.False(t, decoder0.HasValueWithSlot(12))
	assert.True(t, decoder0.HasValueWithSlot(13))

	decoder := NewTSDDecoder(data)
	assert.False(t, decoder.HasValueWithSlot(9))
	assert.True(t, decoder.HasValueWithSlot(10))
	assert.False(t, decoder.HasValueWithSlot(13))
	assert.Equal(t, uint64(10), decoder.Value())
	assert.True(t, decoder.HasValueWithSlot(11))
	assert.Equal(t, uint64(100), decoder.Value())
	assert.False(t, decoder.HasValueWithSlot(12))
	assert.True(t, decoder.HasValueWithSlot(13))
	assert.Equal(t, uint64(50), decoder.Value())
	// out of range
	assert.False(t, decoder.HasValueWithSlot(9))
	assert.False(t, decoder.HasValueWithSlot(100))

	decoder.Reset(data)
	result := map[uint16]uint64{
		10: uint64(10),
		11: uint64(100),
		13: uint64(50),
	}
	c := 0
	total := 0
	for decoder.Next() {
		if decoder.HasValue() {
			assert.Equal(t, result[decoder.Slot()], decoder.Value())
			c++
		}
		total++
	}
	assert.Equal(t, 3, c)
	assert.Equal(t, 4, total)
}

func Test_Empty_TSDDecoder(t *testing.T) {
	decoder := NewTSDDecoder(nil)
	assert.Nil(t, decoder.Error())
	assert.Equal(t, uint64(0), decoder.Value())
	assert.False(t, decoder.HasValue())
}

func TestGetTSDDecoder(t *testing.T) {
	decoder := GetTSDDecoder()
	assert.NotNil(t, decoder)
	ReleaseTSDDecoder(decoder)
}

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
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/stream"
)

// for testing
var (
	TSDEncodeFunc = GetTSDEncoder
	flushFunc     = flush
)

var (
	decoderPool = sync.Pool{
		New: func() interface{} {
			return NewTSDDecoder(nil)
		},
	}
	encoderPool = sync.Pool{}
)

func GetTSDDecoder() *TSDDecoder {
	decoder := decoderPool.Get()
	return decoder.(*TSDDecoder)
}

func ReleaseTSDDecoder(decoder *TSDDecoder) {
	if decoder != nil {
		decoderPool.Put(decoder)
	}
}

func GetTSDEncoder(startTime uint16) *TSDEncoder {
	encoderIntf := encoderPool.Get()
	if encoderIntf == nil {
		return NewTSDEncoder(startTime)
	}
	encoder := encoderIntf.(*TSDEncoder)
	encoder.RestWithStartTime(startTime)
	return encoder
}

func ReleaseTSDEncoder(encoder *TSDEncoder) {
	if encoder != nil {
		encoderPool.Put(encoder)
	}
}

// TSDEncoder encodes time series data point
type TSDEncoder struct {
	startTime  uint16
	bitBuffer  bytes.Buffer
	bitWriter  *bit.Writer
	values     *XOREncoder
	count      uint16
	err        error
	timeBitBuf bytes.Buffer // time + bitBuffer
}

// NewTSDEncoder creates tsd encoder instance
func NewTSDEncoder(startTime uint16) *TSDEncoder {
	e := &TSDEncoder{startTime: startTime}
	e.bitWriter = bit.NewWriter(&e.bitBuffer)
	e.values = NewXOREncoder(e.bitWriter)
	return e
}

// Reset resets the underlying bytes.Buffer
func (e *TSDEncoder) Reset() {
	e.bitBuffer.Reset()
	e.bitWriter.Reset(&e.bitBuffer)
	e.values.Reset()
	e.timeBitBuf.Reset()
}

// RestWithStartTime resets the buffer and slot info
func (e *TSDEncoder) RestWithStartTime(startTime uint16) {
	e.startTime = startTime
	e.count = 0
	e.err = nil
	e.Reset()
}

// EmitDownSamplingValue appends the value after down sampling
// Inf value symbols a empty value to omit
func (e *TSDEncoder) EmitDownSamplingValue(pos int, value float64) {
	_ = pos
	if math.IsInf(value, 1) {
		e.AppendTime(bit.Zero)
		return
	}
	e.AppendTime(bit.One)
	e.AppendValue(math.Float64bits(value))
}

// AppendTime appends time slot, marks time slot if has data point
func (e *TSDEncoder) AppendTime(slot bit.Bit) {
	if e.err != nil {
		return
	}
	e.err = e.bitWriter.WriteBit(slot)
	e.count++
}

// AppendValue appends data point value
func (e *TSDEncoder) AppendValue(value uint64) {
	if e.err != nil {
		return
	}
	e.err = e.values.Write(value)
}

// Bytes returns binary which compress time series data point
func (e *TSDEncoder) Bytes() ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	if err := flushFunc(e.bitWriter); err != nil {
		return nil, err
	}
	if e.count == 0 {
		// if no data add in tsd stream, return nil,
		// if return data with empty data, will get wrong start/end time range(because end is negative)
		return nil, nil
	}

	e.timeBitBuf.Reset()
	var scratch [4]byte
	stream.PutUint16(scratch[:], 0, e.startTime)
	stream.PutUint16(scratch[:], 2, e.startTime+e.count-1)
	_, _ = e.timeBitBuf.Write(scratch[:])
	_, _ = e.timeBitBuf.Write(e.bitBuffer.Bytes())
	return e.timeBitBuf.Bytes(), nil
}

// BytesWithoutTime returns binary which compress time series data point without time slot range
func (e *TSDEncoder) BytesWithoutTime() ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	if err := flushFunc(e.bitWriter); err != nil {
		return nil, err
	}
	return e.bitBuffer.Bytes(), nil
}

func flush(writer *bit.Writer) error {
	return writer.Flush()
}

// TSDDecoder decodes time series compress data
type TSDDecoder struct {
	startTime, endTime uint16

	reader *bit.Reader
	values *XORDecoder
	buf    *bufioutil.Buffer

	idx uint16

	err error
}

// NewTSDDecoder create tsd decoder instance
func NewTSDDecoder(data []byte) *TSDDecoder {
	decoder := &TSDDecoder{}
	if len(data) > 4 {
		decoder.Reset(data)
	}
	return decoder
}

// ResetWithTimeRange resets tsd data and reads the meta info from the data with time range
func (d *TSDDecoder) ResetWithTimeRange(data []byte, start, end uint16) {
	d.reset(data)

	d.startTime = start
	d.endTime = end

	d.reader.Reset()
}

// Reset resets tsd data and reads the meta info from the data
func (d *TSDDecoder) Reset(data []byte) {
	if len(data) <= 4 {
		d.err = fmt.Errorf("TSDDecoder resets with bad data")
		return
	}

	d.reset(data)

	d.startTime = binary.LittleEndian.Uint16(data[0:2])
	d.endTime = binary.LittleEndian.Uint16(data[2:4])
	d.buf.SetIdx(4)

	d.reader.Reset()
}

func (d *TSDDecoder) reset(data []byte) {
	if d.buf == nil {
		d.buf = bufioutil.NewBuffer(data)
		d.reader = bit.NewReader(d.buf)
		d.values = NewXORDecoder(d.reader)
	} else {
		d.values.Reset()
		d.buf.SetBuf(data)
	}
	d.idx = 0
	d.err = nil
}

// Error returns decode error
func (d *TSDDecoder) Error() error {
	return d.err
}

// StartTime returns tsd start time slot
func (d *TSDDecoder) StartTime() uint16 {
	return d.startTime
}

// EndTime returns tsd end time slot
func (d *TSDDecoder) EndTime() uint16 {
	return d.endTime
}

// Next returns if has next slot data
func (d *TSDDecoder) Next() bool {
	if d.startTime+d.idx <= d.endTime {
		d.idx++
		return true
	}
	return false
}

// Seek seeks and reads at the specified slot
func (d *TSDDecoder) Seek(slot uint16) bool {
	if slot > d.endTime || slot < d.startTime {
		return false
	}
	for d.idx+d.startTime < slot {
		if d.HasValueWithSlot(d.idx + d.startTime) {
			_ = d.Value()
		} else {
			return false
		}
	}
	return d.idx+d.startTime == slot
}

// HasValue returns slot value if exist
func (d *TSDDecoder) HasValue() bool {
	if d.reader == nil {
		return false
	}
	b, err := d.reader.ReadBit()
	if err != nil {
		d.err = err
		return false
	}
	return b == bit.One
}

// HasValueWithSlot returns value if exist by given time slot
func (d *TSDDecoder) HasValueWithSlot(slot uint16) bool {
	if slot < d.startTime || slot > d.endTime {
		return false
	}
	if slot == d.idx+d.startTime {
		d.idx++
		return d.HasValue()
	}
	return false
}

func (d *TSDDecoder) Slot() uint16 {
	return d.startTime + d.idx - 1
}

// Value returns value of time slot
func (d *TSDDecoder) Value() uint64 {
	if d.values == nil {
		return 0
	}
	if d.values.Next() {
		return d.values.Value()
	}
	return 0
}

// DecodeTSDTime decodes start-time-slot and end-time-slot of tsd.
// a simple method extracted from NewTSDDecoder to reduce gc pressure.
func DecodeTSDTime(data []byte) (startTime, endTime uint16) {
	startTime = binary.LittleEndian.Uint16(data[0:2])
	endTime = binary.LittleEndian.Uint16(data[2:4])
	return
}

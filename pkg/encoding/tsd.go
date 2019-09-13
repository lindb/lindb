package encoding

import (
	"bytes"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/stream"
)

// TSDEncoder encodes time series data point
type TSDEncoder struct {
	startTime    int
	values       *XOREncoder
	bitBuffer    bytes.Buffer
	timeSlots    *bit.Writer
	count        int
	streamBuffer bytes.Buffer
	sw           *stream.BufferWriter
	err          error
}

// NewTSDEncoder creates tsd encoder instance
func NewTSDEncoder(startTime int) *TSDEncoder {
	e := &TSDEncoder{
		startTime: startTime,
		values:    NewXOREncoder(),
	}
	e.timeSlots = bit.NewWriter(&e.bitBuffer)
	e.sw = stream.NewBufferWriter(&e.streamBuffer)
	return e
}

// Reset resets the underlying bytes.Buffer
func (e *TSDEncoder) Reset() {
	e.bitBuffer.Reset()
	e.timeSlots.Reset(&e.bitBuffer)

	e.streamBuffer.Reset()
	e.values.Reset()
}

// AppendTime appends time slot, marks time slot if has data point
func (e *TSDEncoder) AppendTime(slot bit.Bit) {
	if e.err != nil {
		return
	}
	e.err = e.timeSlots.WriteBit(slot)
	e.count++
}

// AppendValue appends data point value
func (e *TSDEncoder) AppendValue(value uint64) {
	if e.err != nil {
		return
	}
	e.err = e.values.Write(value)
}

// Error returns tsd encode error
func (e *TSDEncoder) Error() error {
	return e.err
}

// Bytes returns binary which compress time series data point
func (e *TSDEncoder) Bytes() ([]byte, error) {
	valueBuf, err := e.values.Bytes()
	if err != nil {
		e.err = err
		return nil, e.err
	}
	e.err = e.timeSlots.Flush()
	if e.err != nil {
		return nil, e.err
	}

	e.sw.PutUvarint32(uint32(e.startTime))
	e.sw.PutUvarint32(uint32(e.count))

	windowBuf := e.bitBuffer.Bytes()
	e.sw.PutUvarint32(uint32(len(windowBuf)))
	e.sw.PutBytes(windowBuf)
	e.sw.PutUvarint32(uint32(len(valueBuf)))
	e.sw.PutBytes(valueBuf)

	return e.sw.Bytes()
}

// TSDDecoder decodes time series compress data
type TSDDecoder struct {
	reader    *stream.Reader
	startTime int
	endTime   int
	count     int

	timeSlots *bit.Reader
	values    *XORDecoder

	idx int

	err error
}

// NewTSDDecoder create tsd decoder instance
func NewTSDDecoder(data []byte) *TSDDecoder {
	decoder := &TSDDecoder{
		reader:    stream.NewReader(nil),
		values:    NewXORDecoder(nil),
		timeSlots: bit.NewReader(nil),
	}
	decoder.Reset(data)
	return decoder
}

// readMeta reads the meta info from the data
func (d *TSDDecoder) Reset(data []byte) {
	d.reader.Reset(data)

	d.startTime = int(d.reader.ReadUvarint32())
	d.count = int(d.reader.ReadUvarint32())
	d.endTime = d.startTime + d.count - 1

	length := d.reader.ReadUvarint32()
	d.timeSlots.Reset(d.reader.ReadSlice(int(length)))

	length = d.reader.ReadUvarint32()
	d.values.Reset(d.reader.ReadSlice(int(length)))
}

// Error returns decode error
func (d *TSDDecoder) Error() error {
	return d.reader.Error()
}

// StartTime returns tsd start time slot
func (d *TSDDecoder) StartTime() int {
	return d.startTime
}

// EndTime returns tsd end time slot
func (d *TSDDecoder) EndTime() int {
	return d.endTime
}

// Next returns if has next slot data
func (d *TSDDecoder) Next() bool {
	if d.count > d.idx {
		d.idx++
		return true
	}
	return false
}

// HasValue returns slot value if exist
func (d *TSDDecoder) HasValue() bool {
	b, err := d.timeSlots.ReadBit()
	if err != nil {
		d.err = err
		return false
	}
	return b == bit.One
}

// HasValueWithSlot returns value if exist by given time slot
func (d *TSDDecoder) HasValueWithSlot(slot int) bool {
	if slot < 0 || slot > d.count {
		return false
	}
	if slot == d.idx {
		d.idx++
		return d.HasValue()
	}
	return false
}

func (d *TSDDecoder) Slot() int {
	return d.startTime + d.idx - 1
}

// Value returns value of time slot
func (d *TSDDecoder) Value() uint64 {
	if d.values.Next() {
		return d.values.Value()
	}
	return 0
}

// DecodeTSDTime decodes start-time-slot and end-time-slot of tsd.
// a simple method extracted from NewTSDDecoder to reduce gc pressure.
func DecodeTSDTime(data []byte) (startTime, endTime int) {
	reader := stream.NewReader(data)
	startTime = int(reader.ReadUvarint32())
	count := int(reader.ReadUvarint32())
	endTime = startTime + count - 1
	return
}

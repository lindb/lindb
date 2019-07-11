package encoding

import (
	"bytes"

	"github.com/eleme/lindb/pkg/bit"
	"github.com/eleme/lindb/pkg/stream"
)

// TSDEncoder encodes time series data point
type TSDEncoder struct {
	startTime int
	values    *XOREncoder
	timeSlots *bit.Writer
	count     int

	buf bytes.Buffer
	err error
}

// NewTSDEncoder creates tsd encoder instance
func NewTSDEncoder(startTime int) *TSDEncoder {
	e := &TSDEncoder{
		startTime: startTime,
		values:    NewXOREncoder(),
	}
	e.timeSlots = bit.NewWriter(&e.buf)
	return e
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

	binary := stream.BinaryWriter()
	binary.PutUvarint32(uint32(e.startTime))
	binary.PutUvarint32(uint32(e.count))

	windowBuf := e.buf.Bytes()
	binary.PutUvarint32(uint32(len(windowBuf)))
	binary.PutBytes(windowBuf)
	binary.PutUvarint32(uint32(len(valueBuf)))
	binary.PutBytes(valueBuf)

	return binary.Bytes()
}

// TSDDecoder decodes time series compress data
type TSDDecoder struct {
	binary    *stream.Binary
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
	binary := stream.BinaryReader(data)
	startTime := int(binary.ReadInt32())
	count := int(binary.ReadInt32())
	length := binary.ReadUvarint32()
	buf := binary.ReadBytes(int(length))
	timeSlots := bit.NewReader(bytes.NewBuffer(buf))
	length = binary.ReadUvarint32()
	buf = binary.ReadBytes(int(length))
	return &TSDDecoder{
		startTime: startTime,
		endTime:   startTime + count - 1,
		count:     count,
		timeSlots: timeSlots,
		values:    NewXORDecoder(buf),
		binary:    binary,
	}
}

// Error returns decode error
func (d *TSDDecoder) Error() error {
	return d.binary.Error()
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
	binary := stream.BinaryReader(data)
	startTime = int(binary.ReadInt32())
	count := int(binary.ReadInt32())
	endTime = startTime + count - 1
	return
}

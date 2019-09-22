package encoding

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/stream"
)

var decoderPool = sync.Pool{
	New: func() interface{} {
		return NewTSDDecoder(nil)
	},
}

func GetTSDDecoder() *TSDDecoder {
	decoder := decoderPool.Get()
	return decoder.(*TSDDecoder)
}

func ReleaseTSDDecoder(decoder *TSDDecoder) {
	decoderPool.Put(decoder)
}

// TSDEncoder encodes time series data point
type TSDEncoder struct {
	startTime int
	bitBuffer bytes.Buffer
	bitWriter *bit.Writer
	values    *XOREncoder
	count     int
	err       error
}

// NewTSDEncoder creates tsd encoder instance
func NewTSDEncoder(startTime int) *TSDEncoder {
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

// Error returns tsd encode error
func (e *TSDEncoder) Error() error {
	return e.err
}

// Bytes returns binary which compress time series data point
func (e *TSDEncoder) Bytes() ([]byte, error) {
	e.err = e.bitWriter.Flush()
	if e.err != nil {
		return nil, e.err
	}
	var buf bytes.Buffer
	writer := stream.NewBufferWriter(&buf)
	writer.PutUInt16(uint16(e.startTime))
	writer.PutUInt16(uint16(e.count))
	writer.PutBytes(e.bitBuffer.Bytes())
	return writer.Bytes()
}

// TSDDecoder decodes time series compress data
type TSDDecoder struct {
	startTime int
	endTime   int
	count     int

	reader *bit.Reader
	values *XORDecoder
	buf    *bufioutil.Buffer

	idx int

	err error
}

// NewTSDDecoder create tsd decoder instance
func NewTSDDecoder(data []byte) *TSDDecoder {
	decoder := &TSDDecoder{}
	if data != nil {
		decoder.Reset(data)
	}
	return decoder
}

// readMeta reads the meta info from the data
func (d *TSDDecoder) Reset(data []byte) {
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

	d.startTime = int(binary.LittleEndian.Uint16(data[0:2]))
	d.count = int(binary.LittleEndian.Uint16(data[2:4]))
	d.endTime = d.startTime + d.count - 1
	d.buf.SetIdx(4)

	d.reader.Reset()
}

// Error returns decode error
func (d *TSDDecoder) Error() error {
	return d.err
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
	startTime = int(reader.ReadUint16())
	count := int(reader.ReadUint16())
	endTime = startTime + count - 1
	return
}

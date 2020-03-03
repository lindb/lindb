package encoding

import (
	"bytes"
	"math"
	"math/bits"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/stream"
)

// reference:
// parquet delta encoding https://github.com/apache/parquet-format/blob/master/Encodings.md#RLE
// <num of values -1(exclude first value)><min delta><bit width of max value(delta of min delta)><list of deltas>
// for singed values, use zigzag encoding(https://developers.google.com/protocol-buffers/docs/encoding#signed-integers)

// DeltaBitPackingEncoder represents a delta encoding for int32
type DeltaBitPackingEncoder struct {
	first    int32
	previous int32
	minDelta int32
	deltas   []int32
	buffer   *bytes.Buffer
	sw       *stream.BufferWriter
	bw       *bit.Writer
	hasFirst bool
}

// NewDeltaBitPackingEncoder creates a delta encoder
func NewDeltaBitPackingEncoder() *DeltaBitPackingEncoder {
	var buffer bytes.Buffer
	return &DeltaBitPackingEncoder{
		buffer: &buffer,
		sw:     stream.NewBufferWriter(&buffer),
		bw:     bit.NewWriter(&buffer)}
}

// Reset clears the underlying data structure to prepare for next use
func (p *DeltaBitPackingEncoder) Reset() {
	p.buffer.Reset()
	p.sw.Reset()
	p.bw.Reset(p.buffer)

	p.hasFirst = false
	p.first = 0
	p.previous = 0
	p.minDelta = int32(math.MaxInt32)
	p.deltas = p.deltas[:0]
}

// Add adds a new int value
func (p *DeltaBitPackingEncoder) Add(v int32) {
	if !p.hasFirst {
		p.hasFirst = true
		p.first = v
		p.previous = v
		return
	}

	delta := p.previous - v

	p.deltas = append(p.deltas, delta)

	if delta < p.minDelta {
		p.minDelta = delta
	}
	p.previous = v
}

// Bytes returns binary data
func (p *DeltaBitPackingEncoder) Bytes() []byte {
	max := uint32(0)
	p.buffer.Reset()

	p.sw.PutVarint32(int32(len(p.deltas))) // write deltas length
	for _, v := range p.deltas {
		deltaDelta := uint32(v - p.minDelta)
		if max < deltaDelta {
			max = deltaDelta
		}
	}
	width := 32 - bits.LeadingZeros32(max)
	p.sw.PutByte(byte(width))                                // width
	p.sw.PutVarint64(int64(ZigZagEncode(int64(p.minDelta)))) // min delta
	p.sw.PutVarint32(p.first)                                // first value

	for _, v := range p.deltas {
		deltaDelta := v - p.minDelta
		_ = p.bw.WriteBits(uint64(deltaDelta), width)
	}

	_ = p.bw.Flush()
	return p.buffer.Bytes()
}

// DeltaBitPackingDecoder represents a delta decoding for int32
type DeltaBitPackingDecoder struct {
	sr       *stream.Reader
	br       *bit.Reader
	buf      *bufioutil.Buffer
	count    int32
	pos      int32
	width    int
	previous int32
	minDelta int32
}

// NewDeltaBitPackingDecoder creates a delta decoder
func NewDeltaBitPackingDecoder(buf []byte) *DeltaBitPackingDecoder {
	d := &DeltaBitPackingDecoder{
		sr: stream.NewReader(nil),
	}
	d.buf = bufioutil.NewBuffer(buf)
	d.br = bit.NewReader(d.buf)

	d.Reset(buf)
	return d
}

func (d *DeltaBitPackingDecoder) Reset(buf []byte) {
	d.sr.Reset(buf)
	x := d.sr.ReadVarint32() // deltas length
	d.count = x + 1
	d.pos = d.count
	w := d.sr.ReadByte() // width
	d.width = int(w)
	min := d.sr.ReadVarint64()
	d.minDelta = int32(ZigZagDecode(uint64(min))) // min delta

	// need read first value
	d.previous = d.sr.ReadVarint32()
	pos := d.sr.Position()

	d.buf.SetBuf(buf[pos:])

	// reset bit stream
	d.br.Reset()
}

// HasNext tests if has more int32 value
func (d *DeltaBitPackingDecoder) HasNext() bool {
	return d.pos > 0
}

// Next returns next value if exist
func (d *DeltaBitPackingDecoder) Next() int32 {
	if d.pos == d.count {
		d.pos--
		return d.previous
	}
	x, _ := d.br.ReadBits(d.width)
	d.pos--
	v := int32(x) + d.minDelta
	vv := d.previous - v
	d.previous = vv
	return vv
}

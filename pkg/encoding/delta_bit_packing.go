package encoding

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	"github.com/eleme/lindb/pkg/bit"
)

//parquet delta encoding https://github.com/apache/parquet-format/blob/master/Encodings.md#RLE
//<num of values -1(exclude first value)><min delta><bit width of max value(delta of min delta)><list of deltas>
//for singed values, use zigzag encoding(https://developers.google.com/protocol-buffers/docs/encoding#signed-integers)

type DeltaBitPackingEncoder struct {
	first    int32
	previous int32
	minDelta int32
	deltas   []int32

	hasFirst bool
}

type DeltaBitPackingDecoder struct {
	buf      *bytes.Buffer
	br       *bit.Reader
	count    int64
	pos      int64
	width    int
	previous int32
	minDelta int32
}

func NewDeltaBitPackingEncoder() *DeltaBitPackingEncoder {
	return &DeltaBitPackingEncoder{}
}

func NewDeltaBitPackingDecoder(buf *[]byte) *DeltaBitPackingDecoder {
	d := &DeltaBitPackingDecoder{
		buf: bytes.NewBuffer(*buf),
	}
	x, _ := binary.ReadVarint(d.buf)
	d.count = x + 1
	d.pos = d.count
	w, _ := d.buf.ReadByte()
	d.width = int(w)

	min, _ := binary.ReadVarint(d.buf)
	d.minDelta = int32(ZigZagDecode(uint64(min)))
	d.br = bit.NewReader(d.buf)
	return d
}

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

func (p *DeltaBitPackingEncoder) Bytes() ([]byte, error) {
	var scratch [binary.MaxVarintLen64]byte
	var buf bytes.Buffer
	n := binary.PutVarint(scratch[:], int64(len(p.deltas)))
	if _, err := buf.Write(scratch[:n]); err != nil {
		return nil, err
	}

	var max int32
	for _, v := range p.deltas {
		deltaDelta := v - p.minDelta
		if max < deltaDelta {
			max = deltaDelta
		}
	}
	width := 32 - bits.LeadingZeros32(uint32(max))
	buf.WriteByte(byte(width))
	n1 := binary.PutVarint(scratch[:], int64(ZigZagEncode(int64(p.minDelta))))
	if _, err := buf.Write(scratch[:n1]); err != nil {
		return nil, err
	}

	n2 := binary.PutVarint(scratch[:], int64(p.first))
	if _, err := buf.Write(scratch[:n2]); err != nil {
		return nil, err
	}

	bw := bit.NewWriter(&buf)
	for _, v := range p.deltas {
		deltaDelta := v - p.minDelta
		err := bw.WriteBits(uint64(deltaDelta), width)
		if err != nil {
			return nil, err
		}
	}

	bw.Flush()
	return buf.Bytes(), nil
}

func (d *DeltaBitPackingDecoder) HasNext() bool {
	return d.pos > 0
}

func (d *DeltaBitPackingDecoder) Next() int32 {
	if d.pos == d.count {
		x, _ := binary.ReadVarint(d.buf)
		d.pos--
		v := int32(x)
		d.previous = v
		return v
	}
	x, _ := d.br.ReadBits(d.width)
	d.pos--
	v := int32(x) + d.minDelta
	vv := d.previous - v
	//fmt.Printf("jj==%d:%d:%d\n", x, v, d.minDelta)
	d.previous = vv
	return vv
}

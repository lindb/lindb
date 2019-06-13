package encoding

import (
	"bytes"
	"math"
	"fmt"
	"math/bits"
	"github.com/eleme/lindb/pkg/bit"
)

type FloatEncoder struct {
	previousVal uint64

	buf bytes.Buffer
	bw  *bit.Writer

	leading  int
	trailing int

	first  bool
	finish bool
}

type FloatDecoder struct {
	val uint64

	b  []byte
	br *bit.Reader

	leading  uint64
	trailing uint64

	first bool
	err   error
}

func NewFloatEncoder() *FloatEncoder {
	s := &FloatEncoder{
		first: true,
		leading: int(^uint8(0)),
	}
	// new bit writer
	s.bw = bit.NewWriter(&s.buf)
	return s
}

func NewFloatDecoder(b []byte) *FloatDecoder {
	s := &FloatDecoder{
		b:     b,
		first: true,
	}
	s.br = bit.NewReader(bytes.NewBuffer(b))
	return s
}

// write float64 v to underlying buffer, using xor compress
// Value is encoded by XOR then with previous value.
// If XOR result is a zero value(value is the same as the previous value),
// only a single '0' bit is stored, otherwise '1' bit is stored.
// For non-zero XOR result, there are two choices:
// 1). If the block of meaningful bits falls in between the block of previous meaningful bits,
//     i.e., there are at least as many leading and trailing zeros as with the previous value,
//           use that information for the block position and just store the XOR value.
// 2). Length of the number of leading zeros is stored in the next 6 bits,
//     then length of the XOR value is stored in the next 6 bits,
//     finally the XOR value is stored.
func (e *FloatEncoder) Write(v float64) error {
	if math.IsNaN(v) {
		return fmt.Errorf("unspoorted value for NaN")
	}

	val := math.Float64bits(v)
	if e.first {
		e.first = false
		e.previousVal = val
		e.bw.WriteBits(val, 64)
		return nil
	}

	delta := val ^ e.previousVal
	if delta == 0 {
		// write '0' bit
		e.bw.WriteBit(bit.Zero)
	} else {
		// write '1' bit
		e.bw.WriteBit(bit.One)

		leading := bits.LeadingZeros64(delta)
		trailing := bits.TrailingZeros64(delta)

		if leading >= e.leading && trailing >= e.trailing {
			// write control bit('1') for using previous block information
			e.bw.WriteBit(bit.One)
			e.bw.WriteBits(delta>>uint(e.trailing), 64-e.leading-e.trailing)
		} else {
			// write control bit('0') for not using previous block information
			e.bw.WriteBit(bit.Zero)
			blockSize := 64 - leading - trailing
			/*
			 * Store the length of the number of leading zeros in the next 6 bits.
			 * Store the length of the meaningful XORed value in the next 6 bits.
			 * Store the meaningful bits of the XOR value.
			 */
			e.bw.WriteBits(uint64(leading), 6)
			e.bw.WriteBits(uint64(blockSize-1), 6)
			e.bw.WriteBits(delta>>uint(trailing), blockSize)

			e.leading = leading
			e.trailing = trailing
		}
	}

	e.previousVal = val

	return nil
}

func (e *FloatEncoder) Flush() error {
	e.finish = true
	err := e.bw.Flush()
	if nil != err {
		return err
	}
	return nil
}

// Bytes returns a copy of the underlying byte buffer
func (e *FloatEncoder) Bytes() []byte {
	return e.buf.Bytes()
}

func (d *FloatDecoder) Next() bool {
	if d.err != nil {
		return false
	}
	if d.first {
		d.first = false
		d.val, d.err = d.br.ReadBits(64)
		return true
	}

	var b bit.Bit
	b, d.err = d.br.ReadBit()
	if d.err != nil {
		return false
	}
	if b == bit.Zero {
		//same as previous
	} else {
		// read control bit
		b, d.err = d.br.ReadBit()
		if d.err != nil {
			return false
		}

		var blockSize uint64
		if b == bit.Zero {
			d.leading, d.err = d.br.ReadBits(6)
			if d.err != nil {
				return false
			}
			blockSize, d.err = d.br.ReadBits(6)
			if d.err != nil {
				return false
			}
			blockSize = blockSize + 1
			d.trailing = 64 - d.leading - blockSize
		} else {
			//reuse previous leading and trailing
			blockSize = 64 - d.leading - d.trailing
		}
		delta, err := d.br.ReadBits(int(blockSize))
		if err != nil {
			d.err = err
			return false
		}
		val := delta << d.trailing
		d.val = d.val ^ val
	}
	return true
}

func (d *FloatDecoder) Value() float64 {
	return math.Float64frombits(d.val)
}

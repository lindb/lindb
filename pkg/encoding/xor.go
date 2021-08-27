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
	"math/bits"

	"github.com/lindb/lindb/pkg/bit"
)

// reference facebook gorilla paper(https://www.vldb.org/pvldb/vol8/p1816-teller.pdf)

const blockSizeAdjustment = 1
const firstValueLen = 64

// XOREncoder encodes uint64 value using xor compress
type XOREncoder struct {
	bw *bit.Writer

	previousVal uint64
	leading     int
	trailing    int

	first bool
	err   error
}

// NewXOREncoder creates xor encoder for compressing uint64 data
func NewXOREncoder(bw *bit.Writer) *XOREncoder {
	return &XOREncoder{
		bw:    bw,
		first: true,
	}
}

func (e *XOREncoder) Reset() {
	e.previousVal = 0
	e.leading = 0
	e.trailing = 0
	e.first = true
	e.err = nil
}

// Write writs uint64 v to underlying buffer, using xor compress
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
func (e *XOREncoder) Write(val uint64) error {
	if e.first {
		// write first value
		e.first = false
		e.previousVal = val
		e.err = e.bw.WriteBits(val, firstValueLen)
		return nil
	}

	// xor using previous value
	delta := val ^ e.previousVal
	if delta == 0 {
		// write '0' bit, same with previous value
		e.err = e.bw.WriteBit(bit.Zero)
	} else {
		// write '1' bit, diff with preivous value
		e.err = e.bw.WriteBit(bit.One)

		leading := bits.LeadingZeros64(delta)
		trailing := bits.TrailingZeros64(delta)

		if leading >= e.leading && trailing >= e.trailing {
			// write control bit('1') for using previous block information
			e.err = e.bw.WriteBit(bit.One)
			e.err = e.bw.WriteBits(delta>>uint(e.trailing), 64-e.leading-e.trailing)
		} else {
			// write control bit('0') for not using previous block information
			e.err = e.bw.WriteBit(bit.Zero)
			blockSize := 64 - leading - trailing
			/*
			 * Store the length of the number of leading zeros in the next 6 bits.
			 * Store the length of the meaningful XORed value in the next 6 bits.
			 * Store the meaningful bits of the XOR value.
			 */
			e.err = e.bw.WriteBits(uint64(leading), 6)
			e.err = e.bw.WriteBits(uint64(blockSize-blockSizeAdjustment), 6)
			e.err = e.bw.WriteBits(delta>>uint(trailing), blockSize)

			e.leading = leading
			e.trailing = trailing
		}
	}

	e.previousVal = val

	return nil
}

// XORDecoder decodes buffer to uint64 values using xor compress
type XORDecoder struct {
	val uint64

	br *bit.Reader

	leading  uint64
	trailing uint64

	first bool
	err   error
}

// NewXORDecoder create decoder uncompress buffer using xor
func NewXORDecoder(br *bit.Reader) *XORDecoder {
	s := &XORDecoder{
		br:    br,
		first: true,
	}
	return s
}

// Reset resets the underlying buffer to decode
func (d *XORDecoder) Reset() {
	d.first = true
	d.leading = 0
	d.trailing = 0
	d.err = nil
	d.val = 0
}

// Next return if has value in buffer using xor, do uncompress logic in next method,
// data format reference encoder format
func (d *XORDecoder) Next() bool {
	// if has err, always return false
	if d.err != nil {
		return false
	}
	if d.first {
		// read first value
		d.first = false
		d.val, d.err = d.br.ReadBits(firstValueLen)
		return true
	}

	var b bit.Bit
	// read delta control bit
	b, d.err = d.br.ReadBit()
	if d.err != nil {
		return false
	}
	if b == bit.Zero {
		//same as previous, do nothing, use previous value directly
	} else {
		// read control bit
		b, d.err = d.br.ReadBit()
		if d.err != nil {
			return false
		}

		var blockSize uint64
		if b == bit.Zero {
			// read leading and trailing, because block is diff with previous
			d.leading, d.err = d.br.ReadBits(6)
			if d.err != nil {
				return false
			}
			blockSize, d.err = d.br.ReadBits(6)
			if d.err != nil {
				return false
			}
			blockSize += blockSizeAdjustment
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
		// calc value
		val := delta << d.trailing
		d.val ^= val
	}
	return true
}

// Value returns uint64 from buffer
func (d *XORDecoder) Value() uint64 {
	return d.val
}

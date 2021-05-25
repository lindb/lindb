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

package bit

import (
	"io"

	"github.com/lindb/lindb/pkg/bufioutil"
)

// Reader reads bits from buffer
type Reader struct {
	buf   *bufioutil.Buffer
	b     byte
	count uint8

	err error
}

// NewReader crate bit reader
func NewReader(buf *bufioutil.Buffer) *Reader {
	return &Reader{
		buf:   buf,
		count: 0}
}

// ReadBit reads a bit, if failure return error
func (r *Reader) ReadBit() (Bit, error) {
	if r.count == 0 {
		r.b, r.err = r.buf.GetByte()
		r.count = 8
	}
	r.count--
	d := r.b & 0x80
	r.b <<= 1
	return d != 0, r.err
}

// ReadBits read number of bits
func (r *Reader) ReadBits(numBits int) (uint64, error) {
	var u uint64

	for numBits >= 8 {
		byt, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		numBits -= 8
	}

	var err error
	for numBits > 0 && err != io.EOF {
		byt, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		numBits--
	}

	return u, nil
}

// ReadByte reads a byte
func (r *Reader) ReadByte() (byte, error) {
	if r.count == 0 {
		r.b, r.err = r.buf.GetByte()
		return r.b, r.err
	}

	byt := r.b

	r.b, r.err = r.buf.GetByte()
	byt |= r.b >> r.count
	r.b <<= 8 - r.count
	return byt, r.err
}

// Reset resets the reader to read from a new slice
func (r *Reader) Reset() {
	r.err = nil
	r.count = 0
	r.b = 0
}

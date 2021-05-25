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
)

// A Bit is a zero or a one
type Bit bool

const (
	// Zero is our exported type for '0' bits
	Zero Bit = false
	// One is our exported type for '1' bits
	One Bit = true
)

// Writer writes bits to an io.Writer
type Writer struct {
	w     io.Writer
	b     [1]byte
	count uint8
}

// NewWriter create bit writer
func NewWriter(w io.Writer) *Writer {
	var bw Writer
	bw.Reset(w)
	return &bw
}

// Reset writes to a new writer
func (w *Writer) Reset(writer io.Writer) {
	w.w = writer
	w.b[0] = 0
	w.count = 8
}

// WriteBit writes a bit value
func (w *Writer) WriteBit(bit Bit) error {
	if bit {
		w.b[0] |= 1 << (w.count - 1)
	}

	w.count--

	if w.count == 0 {
		// fill byte to io.Writer
		if n, err := w.w.Write(w.b[:]); n != 1 || err != nil {
			return err
		}
		w.b[0] = 0
		w.count = 8
	}

	return nil
}

// WriteBits writes number of bits
func (w *Writer) WriteBits(u uint64, numBits int) error {
	u <<= 64 - uint(numBits)

	for numBits >= 8 {
		byt := byte(u >> 56)
		err := w.WriteByte(byt)
		if err != nil {
			return err
		}
		u <<= 8
		numBits -= 8
	}

	for numBits > 0 {
		err := w.WriteBit((u >> 63) == 1)
		if err != nil {
			return err
		}
		u <<= 1
		numBits--
	}

	return nil
}

// WriteByte write a byte
func (w *Writer) WriteByte(b byte) error {
	w.b[0] |= b >> (8 - w.count)

	if n, err := w.w.Write(w.b[:]); n != 1 || err != nil {
		return err
	}

	w.b[0] = b << w.count

	return nil
}

// Flush flushes the currently in-process byte
func (w *Writer) Flush() error {
	if w.count != 8 {
		_, err := w.w.Write(w.b[:])
		return err
	}
	return nil
}

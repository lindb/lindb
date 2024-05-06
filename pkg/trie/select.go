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

package trie

import (
	"fmt"
	"io"
	"math/bits"

	"github.com/lindb/lindb/pkg/encoding"
)

const selectSampleInterval = 64

type selectVector struct {
	bitVector
	numOnes   uint32
	selectLut []uint32
}

func (v *selectVector) Init(levels []*Level, bitmapType BitmapType) *selectVector {
	v.bitVector.Init(levels, bitmapType)
	v.selectLut = v.selectLut[:0]
	v.selectLut = append(v.selectLut, 0)
	sampledOnes := selectSampleInterval
	onesUptoWord := 0
	for i, w := range v.bits {
		ones := bits.OnesCount64(w)
		for sampledOnes <= onesUptoWord+ones {
			diff := sampledOnes - onesUptoWord
			targetPos := i*wordSize + int(select64(w, int64(diff)))
			v.selectLut = append(v.selectLut, uint32(targetPos))
			sampledOnes += selectSampleInterval
		}
		onesUptoWord += ones
	}

	v.numOnes = uint32(onesUptoWord)
	return v
}

func (v *selectVector) lutSize() uint32 {
	return (v.numOnes/selectSampleInterval + 1) * 4
}

// Select returns the position of the rank-th 1 bit.
// position is zero-based; rank is one-based.
// E.g., for bitvector: 100101000, select(3) = 5
func (v *selectVector) Select(rank uint32) uint32 {
	lutIdx := rank / selectSampleInterval
	rankLeft := rank % selectSampleInterval
	if lutIdx == 0 {
		rankLeft--
	}

	pos := v.selectLut[lutIdx]
	if rankLeft == 0 {
		return pos
	}

	wordOff := pos / wordSize
	bitsOff := pos % wordSize
	if bitsOff == wordSize-1 {
		wordOff++
		bitsOff = 0
	} else {
		bitsOff++
	}

	w := v.bits[wordOff] >> bitsOff << bitsOff
	ones := uint32(bits.OnesCount64(w))
	for ones < rankLeft {
		wordOff++
		w = v.bits[wordOff]
		rankLeft -= ones
		ones = uint32(bits.OnesCount64(w))
	}

	return wordOff*wordSize + uint32(select64(w, int64(rankLeft)))
}

func (v *selectVector) MarshalSize() int {
	return 4 + 4 + int(v.bitsSize()+v.lutSize())
}

func (v *selectVector) Write(w io.Writer) error {
	if err := v.write(w); err != nil {
		return err
	}
	var buf [4]byte
	endian.PutUint32(buf[:], v.numOnes)
	_, err := w.Write(buf[:])
	if err != nil {
		return err
	}
	lutBlk := v.numOnes/selectSampleInterval + 1
	_, err = w.Write(encoding.U32SliceToBytes(v.selectLut[:lutBlk]))
	if err != nil {
		return err
	}
	return nil
}

func (v *selectVector) Unmarshal(buf []byte) ([]byte, error) { //nolint:dupl
	if len(buf) < 8 {
		return nil, fmt.Errorf("cannot read numBits and numOnes of selectVector")
	}
	buf, err := v.unmarshal(buf)
	if err != nil {
		return nil, err
	}
	v.numOnes = endian.Uint32(buf[:4])
	buf = buf[4:]
	// read lut
	lutSize := int(v.lutSize())
	if len(buf) < lutSize {
		return nil, fmt.Errorf("cannot read lut: %d from selectVector:%d", lutSize, len(buf))
	}
	v.selectLut = encoding.BytesToU32Slice(buf[:lutSize])
	buf = buf[lutSize:]
	return buf, nil
}

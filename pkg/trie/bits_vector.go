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
	"strings"
)

type bitVector struct {
	numBits uint32
	bits    []uint64
	words   uint32
}

func (v *bitVector) numWords() uint32 {
	wordSz := v.numBits / wordSize
	if v.numBits%wordSize != 0 {
		wordSz++
	}
	return wordSz
}

func (v *bitVector) String() string {
	var s strings.Builder
	for i := uint32(0); i < v.numBits; i++ {
		if readBit(v.bits, i) {
			s.WriteString("1")
		} else {
			s.WriteString("0")
		}
	}
	return s.String()
}

func (v *bitVector) bitsSize() uint32 {
	return v.words * 8
}

func (v *bitVector) Init(levels []*Level, bitmapType BitmapType) {
	numBits := 0
	for _, n := range levels {
		if bitmapType == HasPrefix {
			numBits += n.nodeCount
		} else {
			numBits += len(n.lsLabels)
		}
	}
	v.numBits = uint32(numBits)
	v.words = v.numWords()
	if uint32(len(v.bits)) < v.words {
		v.bits = make([]uint64, v.words)
	}

	var wordID, bitShift, n uint32
	for _, level := range levels {
		bitsBlock := level.GetBitmap(bitmapType)
		if bitmapType == HasPrefix {
			n = uint32(level.nodeCount)
		} else {
			n = uint32(len(level.lsLabels))
		}
		if n == 0 {
			continue
		}

		nCompleteWords := n / wordSize
		for word := 0; uint32(word) < nCompleteWords; word++ {
			v.bits[wordID] |= bitsBlock[word] << bitShift
			wordID++
			if bitShift > 0 {
				v.bits[wordID] |= bitsBlock[word] >> (wordSize - bitShift)
			}
		}

		remain := n % wordSize
		if remain > 0 {
			lastWord := bitsBlock[nCompleteWords]
			v.bits[wordID] |= lastWord << bitShift
			if bitShift+remain <= wordSize {
				bitShift = (bitShift + remain) % wordSize
				if bitShift == 0 {
					wordID++
				}
			} else {
				wordID++
				v.bits[wordID] |= lastWord >> (wordSize - bitShift)
				bitShift = bitShift + remain - wordSize
			}
		}
	}
}

func (v *bitVector) IsSet(pos uint32) bool {
	return readBit(v.bits, pos)
}

func (v *bitVector) DistanceToNextSetBit(pos uint32) uint32 {
	var distance uint32 = 1
	wordOff := (pos + 1) / wordSize
	bitsOff := (pos + 1) % wordSize

	if wordOff >= uint32(len(v.bits)) {
		return 0
	}

	testBits := v.bits[wordOff] >> bitsOff
	if testBits > 0 {
		return distance + uint32(bits.TrailingZeros64(testBits))
	}

	numWords := v.words
	if wordOff == numWords-1 {
		return v.numBits - pos
	}
	distance += wordSize - bitsOff

	for wordOff < numWords-1 {
		wordOff++
		testBits = v.bits[wordOff]
		if testBits > 0 {
			return distance + uint32(bits.TrailingZeros64(testBits))
		}
		distance += wordSize
	}

	if wordOff == numWords-1 && v.numBits%64 != 0 {
		distance -= wordSize - v.numBits%64
	}

	return distance
}

func (v *bitVector) write(w io.Writer) error {
	var buf [4]byte
	endian.PutUint32(buf[:], v.numBits)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := w.Write(u64SliceToBytes(v.bits[:v.words])); err != nil {
		return err
	}
	return nil
}

func (v *bitVector) unmarshal(buf []byte) ([]byte, error) {
	v.numBits = endian.Uint32(buf[:4])
	v.words = v.numWords()
	buf = buf[4:]
	// read bits
	bitsSize := int(v.bitsSize())
	if len(buf) < bitsSize {
		return nil, fmt.Errorf("cannot read bits: %d from rankVector: %d", bitsSize, len(buf))
	}
	v.bits = bytesToU64Slice(buf[:bitsSize])
	buf = buf[bitsSize:]
	return buf, nil
}

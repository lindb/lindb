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

	"github.com/lindb/lindb/pkg/encoding"
)

const (
	rankSparseBlockSize = 512
)

type rankVector struct {
	bitVector
	blockSize uint32
	rankLut   []uint32
}

func (v *rankVector) init(blockSize uint32, levels []*Level, bitmapType BitmapType) *rankVector {
	v.Init(levels, bitmapType)
	v.blockSize = blockSize
	wordPerBlk := v.blockSize / wordSize
	nblks := v.numBits/v.blockSize + 1
	if len(v.rankLut) < int(nblks) {
		v.rankLut = make([]uint32, nblks)
	}

	var totalRank, i uint32
	for i = 0; i < nblks-1; i++ {
		v.rankLut[i] = totalRank
		totalRank += popcountBlock(v.bits, i*wordPerBlk, v.blockSize)
	}
	v.rankLut[nblks-1] = totalRank
	return v
}

func (v *rankVector) lutSize() uint32 {
	return (v.numBits/v.blockSize + 1) * 4
}

func (v *rankVector) MarshalSize() int {
	return 4 + 4 + int(v.bitsSize()+v.lutSize())
}

func (v *rankVector) Write(w io.Writer) error {
	if err := v.write(w); err != nil {
		return err
	}
	var buf [4]byte
	endian.PutUint32(buf[:], v.blockSize)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	nblks := v.numBits/v.blockSize + 1
	if _, err := w.Write(encoding.U32SliceToBytes(v.rankLut[:nblks])); err != nil {
		return err
	}
	return nil
}

func (v *rankVector) Unmarshal(buf []byte) ([]byte, error) { //nolint:dupl
	if len(buf) < 8 {
		return nil, fmt.Errorf("cannot read numBits and blockSize of rankVector")
	}
	buf, err := v.unmarshal(buf)
	if err != nil {
		return nil, err
	}
	v.blockSize = endian.Uint32(buf[:4])
	buf = buf[4:]
	// reading lut
	lutSize := int(v.lutSize())
	if len(buf) < lutSize {
		return nil, fmt.Errorf("cannot read lut: %d from rankVector: %d", lutSize, len(buf))
	}
	v.rankLut = encoding.BytesToU32Slice(buf[:lutSize])
	buf = buf[lutSize:]
	return buf, nil
}

type rankVectorSparse struct {
	rankVector
}

func (v *rankVectorSparse) Init(levels []*Level, bitmapType BitmapType) {
	v.init(rankSparseBlockSize, levels, bitmapType)
}

func (v *rankVectorSparse) Rank(pos uint32) uint32 {
	wordPreBlk := uint32(rankSparseBlockSize / wordSize)
	blockOff := pos / rankSparseBlockSize
	bitsOff := pos % rankSparseBlockSize

	return v.rankLut[blockOff] + popcountBlock(v.bits, blockOff*wordPreBlk, bitsOff+1)
}

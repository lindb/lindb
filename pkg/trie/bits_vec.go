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
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"sort"
	"strings"
)

type bitVector struct {
	numBits uint32
	bits    []uint64
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
	return v.numWords() * 8
}

func (v *bitVector) Init(bitsPerLevel [][]uint64, numBitsPerLevel []uint32) {
	for _, n := range numBitsPerLevel {
		v.numBits += n
	}

	v.bits = make([]uint64, v.numWords())

	var wordID, bitShift uint32
	for level, bitsBlock := range bitsPerLevel {
		n := numBitsPerLevel[level]
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

	numWords := v.numWords()
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

type valueVector struct {
	bytes      []byte
	valueWidth uint32
}

func (v *valueVector) Init(valuesPerLevel [][]byte, valueWidth uint32) {
	var size int
	for l := range valuesPerLevel {
		size += len(valuesPerLevel[l])
	}
	v.valueWidth = valueWidth
	v.bytes = make([]byte, size)

	var pos uint32
	for _, val := range valuesPerLevel {
		copy(v.bytes[pos:], val)
		pos += uint32(len(val))
	}
}

func (v *valueVector) Get(pos uint32) []byte {
	off := pos * v.valueWidth
	return v.bytes[off : off+v.valueWidth]
}

func (v *valueVector) MarshalSize() int64 {
	return align(v.rawMarshalSize())
}

func (v *valueVector) rawMarshalSize() int64 {
	return 8 + int64(len(v.bytes))
}

func (v *valueVector) Write(w io.Writer) error {
	var bs [4]byte
	endian.PutUint32(bs[:], uint32(len(v.bytes)))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}

	endian.PutUint32(bs[:], v.valueWidth)
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}

	if _, err := w.Write(v.bytes); err != nil {
		return err
	}

	var zeros [8]byte
	padding := v.MarshalSize() - v.rawMarshalSize()
	_, err := w.Write(zeros[:padding])
	return err
}

func (v *valueVector) Unmarshal(buf []byte) ([]byte, error) {
	if len(buf) < 8 {
		return nil, fmt.Errorf("cannot read valueWidth and bytes of valueVector")
	}
	var position int64
	sz := endian.Uint32(buf[:4])
	v.valueWidth = endian.Uint32(buf[4:8])
	buf = buf[8:]
	position += 8
	// read bytes
	if uint32(len(buf)) < sz {
		return nil, fmt.Errorf("cannot read bytes: %d from valueVector: %d", sz, len(buf))
	}
	v.bytes = buf[:sz]
	buf = buf[sz:]
	position += int64(sz)

	// read padding
	paddingWidth := align(position) - position
	if int64(len(buf)) < paddingWidth {
		return nil, fmt.Errorf("cannot read padding: %d from valueVector: %d", paddingWidth, len(buf))
	}
	return buf[paddingWidth:], nil
}

const selectSampleInterval = 64

type selectVector struct {
	bitVector
	numOnes   uint32
	selectLut []uint32
}

func (v *selectVector) Init(bitsPerLevel [][]uint64, numBitsPerLevel []uint32) *selectVector {
	v.bitVector.Init(bitsPerLevel, numBitsPerLevel)
	lut := []uint32{0}
	sampledOnes := selectSampleInterval
	onesUptoWord := 0
	for i, w := range v.bits {
		ones := bits.OnesCount64(w)
		for sampledOnes <= onesUptoWord+ones {
			diff := sampledOnes - onesUptoWord
			targetPos := i*wordSize + int(select64(w, int64(diff)))
			lut = append(lut, uint32(targetPos))
			sampledOnes += selectSampleInterval
		}
		onesUptoWord += ones
	}

	v.numOnes = uint32(onesUptoWord)
	v.selectLut = make([]uint32, len(lut))
	copy(v.selectLut, lut)

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

func (v *selectVector) MarshalSize() int64 {
	return align(v.rawMarshalSize())
}

func (v *selectVector) rawMarshalSize() int64 {
	return 4 + 4 + int64(v.bitsSize()) + int64(v.lutSize())
}

func (v *selectVector) Write(w io.Writer) error {
	var buf [4]byte
	endian.PutUint32(buf[:], v.numBits)
	_, err := w.Write(buf[:])
	if err != nil {
		return err
	}
	endian.PutUint32(buf[:], v.numOnes)
	_, err = w.Write(buf[:])
	if err != nil {
		return err
	}
	if _, err := w.Write(u64SliceToBytes(v.bits)); err != nil {
		return err
	}
	if _, err := w.Write(u32SliceToBytes(v.selectLut)); err != nil {
		return err
	}

	var zeros [8]byte
	padding := v.MarshalSize() - v.rawMarshalSize()
	_, err = w.Write(zeros[:padding])
	return err
}

func (v *selectVector) Unmarshal(buf []byte) ([]byte, error) {
	if len(buf) < 8 {
		return nil, fmt.Errorf("cannot read numBits and numOnes of selectVector")
	}
	var position int64
	v.numBits = endian.Uint32(buf[:4])
	v.numOnes = endian.Uint32(buf[4:8])
	buf = buf[8:]
	position += 8
	// read bits
	bitsSize := int(v.bitsSize())
	if len(buf) < bitsSize {
		return nil, fmt.Errorf("cannot read bitsSize: %d from selectVector:%d", bitsSize, len(buf))
	}
	v.bits = bytesToU64Slice(buf[:bitsSize])
	buf = buf[bitsSize:]
	position += int64(bitsSize)

	// read lut
	lutSize := int(v.lutSize())
	if len(buf) < lutSize {
		return nil, fmt.Errorf("cannot read lut: %d from selectVector:%d", lutSize, len(buf))
	}
	v.selectLut = bytesToU32Slice(buf[:lutSize])
	buf = buf[lutSize:]
	position += int64(lutSize)

	// read padding
	paddingWidth := align(position) - position

	if int64(len(buf)) < paddingWidth {
		return nil, fmt.Errorf("cannot read padding: %d from selectVector: %d", paddingWidth, len(buf))
	}
	return buf[paddingWidth:], nil
}

const (
	rankSparseBlockSize = 512
)

type rankVector struct {
	bitVector
	blockSize uint32
	rankLut   []uint32
}

func (v *rankVector) init(blockSize uint32, bitsPerLevel [][]uint64, numBitsPerLevel []uint32) *rankVector {
	v.bitVector.Init(bitsPerLevel, numBitsPerLevel)
	v.blockSize = blockSize
	wordPerBlk := v.blockSize / wordSize
	nblks := v.numBits/v.blockSize + 1
	v.rankLut = make([]uint32, nblks)

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

func (v *rankVector) MarshalSize() int64 {
	return align(v.rawMarshalSize())
}

func (v *rankVector) rawMarshalSize() int64 {
	return 4 + 4 + int64(v.bitsSize()) + int64(v.lutSize())
}

func (v *rankVector) Write(w io.Writer) error {
	var buf [4]byte
	endian.PutUint32(buf[:], v.numBits)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	endian.PutUint32(buf[:], v.blockSize)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := w.Write(u64SliceToBytes(v.bits)); err != nil {
		return err
	}
	if _, err := w.Write(u32SliceToBytes(v.rankLut)); err != nil {
		return err
	}

	var zeros [8]byte
	padding := v.MarshalSize() - v.rawMarshalSize()
	_, err := w.Write(zeros[:padding])
	return err
}

func (v *rankVector) Unmarshal(buf []byte) ([]byte, error) {
	if len(buf) < 8 {
		return nil, fmt.Errorf("cannot read numBits and blockSize of rankVector")
	}
	v.numBits = endian.Uint32(buf[:4])
	v.blockSize = endian.Uint32(buf[4:8])
	buf = buf[8:]
	// read bits
	bitsSize := int(v.bitsSize())
	if len(buf) < bitsSize {
		return nil, fmt.Errorf("cannot read bits: %d from rankVector: %d", bitsSize, len(buf))
	}
	v.bits = bytesToU64Slice(buf[:bitsSize])
	buf = buf[bitsSize:]

	// reading lut
	lutSize := int(v.lutSize())
	if len(buf) < lutSize {
		return nil, fmt.Errorf("cannot read lut: %d from rankVector: %d", lutSize, len(buf))
	}
	v.rankLut = bytesToU32Slice(buf[:lutSize])
	buf = buf[lutSize:]

	// read padding
	position := 8 + bitsSize + lutSize
	paddingWidth := align(int64(position)) - int64(position)
	if int64(len(buf)) < paddingWidth {
		return nil, fmt.Errorf("cannot read padding: %d from rankVector: %d", paddingWidth, len(buf))
	}
	return buf[paddingWidth:], nil
}

type rankVectorSparse struct {
	rankVector
}

func (v *rankVectorSparse) Init(bitsPerLevel [][]uint64, numBitsPerLevel []uint32) {
	v.rankVector.init(rankSparseBlockSize, bitsPerLevel, numBitsPerLevel)
}

func (v *rankVectorSparse) Rank(pos uint32) uint32 {
	wordPreBlk := uint32(rankSparseBlockSize / wordSize)
	blockOff := pos / rankSparseBlockSize
	bitsOff := pos % rankSparseBlockSize

	return v.rankLut[blockOff] + popcountBlock(v.bits, blockOff*wordPreBlk, bitsOff+1)
}

const labelTerminator = 0xff

type labelVector struct {
	labels []byte
}

func (v *labelVector) Init(labelsPerLevel [][]byte, endLevel uint32) {
	numBytes := 1
	for l := uint32(0); l < endLevel; l++ {
		numBytes += len(labelsPerLevel[l])
	}
	v.labels = make([]byte, numBytes)

	var pos uint32
	for l := uint32(0); l < endLevel; l++ {
		copy(v.labels[pos:], labelsPerLevel[l])
		pos += uint32(len(labelsPerLevel[l]))
	}
}

func (v *labelVector) GetLabel(pos uint32) byte {
	return v.labels[pos]
}

func (v *labelVector) Search(k byte, off, size uint32) (uint32, bool) {
	start := off
	if size > 1 && v.labels[start] == labelTerminator {
		start++
		size--
	}

	end := start + size
	if end > uint32(len(v.labels)) {
		end = uint32(len(v.labels))
	}
	result := bytes.IndexByte(v.labels[start:end], k)
	if result < 0 {
		return off, false
	}
	return start + uint32(result), true
}

func (v *labelVector) SearchGreaterThan(label byte, pos, size uint32) (uint32, bool) {
	if size > 1 && v.labels[pos] == labelTerminator {
		pos++
		size--
	}

	result := sort.Search(int(size), func(i int) bool { return v.labels[pos+uint32(i)] > label })
	if uint32(result) == size {
		return pos + uint32(result) - 1, false
	}
	return pos + uint32(result), true
}

func (v *labelVector) MarshalSize() int64 {
	return align(v.rawMarshalSize())
}

func (v *labelVector) rawMarshalSize() int64 {
	return 4 + int64(len(v.labels))
}

func (v *labelVector) Write(w io.Writer) error {
	var bs [4]byte
	endian.PutUint32(bs[:], uint32(len(v.labels)))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	if _, err := w.Write(v.labels); err != nil {
		return err
	}

	padding := v.MarshalSize() - v.rawMarshalSize()
	var zeros [8]byte
	_, err := w.Write(zeros[:padding])
	return err
}

func (v *labelVector) Unmarshal(buf []byte) ([]byte, error) {
	if len(buf) < 4 {
		return nil, fmt.Errorf("failed unmarsh labelVector, length: %d is too short", len(buf))
	}
	l := endian.Uint32(buf)
	tail := align(int64(4 + l))
	if tail > int64(len(buf)) {
		return nil, fmt.Errorf("failed unmarsh labelVector, offset:%d > %d ", tail, len(buf))
	}
	v.labels = buf[4 : 4+l]
	return buf[tail:], nil
}

type compressPathVector struct {
	hasPathVector rankVectorSparse
	offsets       []uint32
	data          []byte
}

func (pv *compressPathVector) Init(hasPathBits [][]uint64, numNodesPerLevel []uint32, paths [][][]byte) {
	pv.hasPathVector.Init(hasPathBits, numNodesPerLevel)
	var offset uint32
	for _, level := range paths {
		for idx := range level {
			pv.offsets = append(pv.offsets, offset)
			offset += uint32(len(level[idx]))
			pv.data = append(pv.data, level[idx]...)
		}
	}
}

func (pv *compressPathVector) rawMarshalSize() int64 {
	return pv.hasPathVector.MarshalSize() + 8 + int64(len(pv.offsets)*4+len(pv.data))
}

func (pv *compressPathVector) MarshalSize() int64 {
	return align(pv.rawMarshalSize())
}

func (pv *compressPathVector) GetPath(nodeID uint32) []byte {
	if !pv.hasPathVector.IsSet(nodeID) {
		return nil
	}
	pathID := pv.hasPathVector.Rank(nodeID) - 1
	start := pv.offsets[pathID]
	end := uint32(len(pv.data))
	if int(pathID+1) < len(pv.offsets) {
		end = pv.offsets[pathID+1]
	}
	return pv.data[start:end]
}

func (pv *compressPathVector) Write(w io.Writer) error {
	if err := pv.hasPathVector.Write(w); err != nil {
		return err
	}

	var length [8]byte
	endian.PutUint32(length[:4], uint32(len(pv.offsets)*4))
	endian.PutUint32(length[4:], uint32(len(pv.data)))

	if _, err := w.Write(length[:]); err != nil {
		return err
	}
	if _, err := w.Write(u32SliceToBytes(pv.offsets)); err != nil {
		return err
	}
	if _, err := w.Write(pv.data); err != nil {
		return err
	}

	padding := pv.MarshalSize() - pv.rawMarshalSize()
	var zeros [8]byte
	_, err := w.Write(zeros[:padding])
	return err
}

func (pv *compressPathVector) Unmarshal(b []byte) ([]byte, error) {
	buf1, err := pv.hasPathVector.Unmarshal(b)
	if err != nil {
		return buf1, err
	}
	var position int64 = 0
	if len(buf1) < 8 {
		return nil, fmt.Errorf("cannot read offsetsLen and dataLen of compressPathVector")
	}
	offsetsLen := endian.Uint32(buf1[0:4])
	dataLen := endian.Uint32(buf1[4:8])
	buf1 = buf1[8:]
	position += 8
	// read offsets
	if uint32(len(buf1)) < offsetsLen+dataLen {
		return nil, fmt.Errorf("offsets+data: %d is longer than remaining buf: %d of compressPathVector", offsetsLen+dataLen, len(buf1))
	}
	pv.offsets = bytesToU32Slice(buf1[:offsetsLen])
	buf1 = buf1[offsetsLen:]
	position += int64(offsetsLen)
	// read data
	pv.data = buf1[:dataLen]
	buf1 = buf1[dataLen:]
	position += int64(dataLen)

	// read padding
	paddingWidth := align(position) - position
	if int64(len(buf1)) < paddingWidth {
		return nil, fmt.Errorf("paddingWidth: %d is longer than remaining buf: %d of compressPathVector", paddingWidth, len(buf1))
	}
	return buf1[paddingWidth:], nil
}

type prefixVector struct {
	compressPathVector
}

func (v *prefixVector) CheckPrefix(key []byte, depth uint32, nodeID uint32) (uint32, bool) {
	prefix := v.GetPrefix(nodeID)
	if len(prefix) == 0 {
		return 0, true
	}
	if int(depth)+len(prefix) > len(key) {
		return 0, false
	}
	if !bytes.Equal(key[depth:depth+uint32(len(prefix))], prefix) {
		return 0, false
	}
	return uint32(len(prefix)), true
}

func (v *prefixVector) GetPrefix(nodeID uint32) []byte {
	return v.GetPath(nodeID)
}

type suffixVector struct {
	compressPathVector
}

func (v *suffixVector) CheckSuffix(key []byte, depth uint32, nodeID uint32) bool {
	suffix := v.GetSuffix(nodeID)
	if depth+1 >= uint32(len(key)) {
		return len(suffix) == 0
	}
	return bytes.Equal(suffix, key[depth+1:])
}

func (v *suffixVector) GetSuffix(nodeID uint32) []byte {
	return v.GetPath(nodeID)
}

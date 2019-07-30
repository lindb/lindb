// Package rsdic provides a rank/select dictionary
// supporting many basic operations in constant time
// using very small working space (smaller than original).
package rsdic

// RSDic provides rank/select operations.
//
// Conceptually RSDic represents a bit vector B[0...num), B[i] = 0 or 1,
// and these bits are set by PushBack (Thus RSDic can handle growing bits).
// All operations (Bit, Rank, Select) are supported in O(1) time.
// (also called as fully indexable dictionary in CS literatures (FID)).
//
// In RSDic, a bit vector is stored in compressed (Note, we don't need to decode all at operations)
// A bit vector is divided into small blocks of length 64, and each small block
// is compressed using enum coding. For example, if a small block contains 10 ones
// and 54 zeros, the block is compressed in 38 bits (See enumCode.go for detail)
// This achieves not only its information theoretic bound, but also achieves more compression
// if same bits appeared togather (e.g. 000...000111...111000...000)
//
// See performance in readme.md
//
// C++ version https://code.google.com/p/rsdic/
// [1] "Fast, Small, Simple Rank/Select on Bitmaps", Gonzalo Navarro and Eliana Providel, SEA 2012

import (
	"github.com/ugorji/go/codec"
)

type RSDic struct {
	bits            []uint64
	pointerBlocks   []uint64
	rankBlocks      []uint64
	selectOneInds   []uint64
	selectZeroInds  []uint64
	rankSmallBlocks []uint8
	num             uint64
	oneNum          uint64
	zeroNum         uint64
	lastBlock       uint64
	lastOneNum      uint64
	lastZeroNum     uint64
	codeLen         uint64
}

// Num returns the number of bits
func (rs RSDic) Num() uint64 {
	return rs.num
}

// OneNum returns the number of ones in bits
func (rs RSDic) OneNum() uint64 {
	return rs.oneNum
}

// ZeroNum returns the number of zeros in bits
func (rs RSDic) ZeroNum() uint64 {
	return rs.zeroNum
}

// PushBack appends the bit to the end of B
func (rs *RSDic) PushBack(bit bool) {
	if (rs.num % kSmallBlockSize) == 0 {
		rs.writeBlock()
	}
	if bit {
		rs.lastBlock |= (1 << (rs.num % kSmallBlockSize))
		if (rs.oneNum % kSelectBlockSize) == 0 {
			rs.selectOneInds = append(rs.selectOneInds, rs.num/kLargeBlockSize)
		}
		rs.oneNum++
		rs.lastOneNum++
	} else {
		if (rs.zeroNum % kSelectBlockSize) == 0 {
			rs.selectZeroInds = append(rs.selectZeroInds, rs.num/kLargeBlockSize)
		}
		rs.zeroNum++
		rs.lastZeroNum++
	}
	rs.num++
}

func (rs *RSDic) writeBlock() {
	if rs.num > 0 {
		rankSB := uint8(rs.lastOneNum)
		rs.rankSmallBlocks = append(rs.rankSmallBlocks, rankSB)
		codeLen := kEnumCodeLength[rankSB]
		code := enumEncode(rs.lastBlock, rankSB)
		newSize := floor(rs.codeLen+uint64(codeLen), kSmallBlockSize)
		if newSize > uint64(len(rs.bits)) {
			rs.bits = append(rs.bits, 0)
		}
		setSlice(rs.bits, rs.codeLen, codeLen, code)
		rs.lastBlock = 0
		rs.lastZeroNum = 0
		rs.lastOneNum = 0
		rs.codeLen += uint64(codeLen)
	}
	if (rs.num % kLargeBlockSize) == 0 {
		rs.rankBlocks = append(rs.rankBlocks, rs.oneNum)
		rs.pointerBlocks = append(rs.pointerBlocks, rs.codeLen)
	}
}

func (rs RSDic) lastBlockInd() uint64 {
	if rs.num == 0 {
		return 0
	}
	return ((rs.num - 1) / kSmallBlockSize) * kSmallBlockSize
}

func (rs RSDic) isLastBlock(pos uint64) bool {
	return pos >= rs.lastBlockInd()
}

// Bit returns the (pos+1)-th bit in bits, i.e. bits[pos]
func (rs RSDic) Bit(pos uint64) bool {
	if rs.isLastBlock(pos) {
		return getBit(rs.lastBlock, uint8(pos%kSmallBlockSize))
	}
	lblock := pos / kLargeBlockSize
	pointer := rs.pointerBlocks[lblock]
	sblock := pos / kSmallBlockSize
	for i := lblock * kSmallBlockPerLargeBlock; i < sblock; i++ {
		pointer += uint64(kEnumCodeLength[rs.rankSmallBlocks[i]])
	}
	rankSB := rs.rankSmallBlocks[sblock]
	code := getSlice(rs.bits, pointer, kEnumCodeLength[rankSB])
	return enumBit(code, rankSB, uint8(pos%kSmallBlockSize))
}

// Rank returns the number of bit's in B[0...pos)
func (rs RSDic) Rank(pos uint64, bit bool) uint64 {
	if pos >= rs.num {
		return bitNum(rs.oneNum, rs.num, bit)
	}
	if rs.isLastBlock(pos) {
		afterRank := popCount(rs.lastBlock >> (pos % kSmallBlockSize))
		return bitNum(rs.oneNum-uint64(afterRank), pos, bit)
	}
	lblock := pos / kLargeBlockSize
	pointer := rs.pointerBlocks[lblock]
	sblock := pos / kSmallBlockSize
	rank := rs.rankBlocks[lblock]
	for i := lblock * kSmallBlockPerLargeBlock; i < sblock; i++ {
		rankSB := rs.rankSmallBlocks[i]
		pointer += uint64(kEnumCodeLength[rankSB])
		rank += uint64(rankSB)
	}
	if pos%kSmallBlockSize == 0 {
		return bitNum(rank, pos, bit)
	}
	rankSB := rs.rankSmallBlocks[sblock]
	code := getSlice(rs.bits, pointer, kEnumCodeLength[rankSB])
	rank += uint64(enumRank(code, rankSB, uint8(pos%kSmallBlockSize)))
	return bitNum(rank, pos, bit)
}

// Select returns the position of (rank+1)-th occurence of bit in B
// Select returns num if rank+1 is larger than the possible range.
// (i.e. Select(oneNum, true) = num, Select(zeroNum, false) = num)
func (rs RSDic) Select(rank uint64, bit bool) uint64 {
	if bit {
		return rs.Select1(rank)
	} else {
		return rs.Select0(rank)
	}
}

func (rs RSDic) Select1(rank uint64) uint64 {
	if rank >= rs.oneNum {
		return rs.num
	} else if rank >= rs.oneNum-rs.lastOneNum {
		lastBlockRank := uint8(rank - (rs.oneNum - rs.lastOneNum))
		return rs.lastBlockInd() + uint64(selectRaw(rs.lastBlock, lastBlockRank+1))
	}
	selectInd := rank / kSelectBlockSize
	lblock := rs.selectOneInds[selectInd]
	for ; lblock < uint64(len(rs.rankBlocks)); lblock++ {
		if rank < rs.rankBlocks[lblock] {
			break
		}
	}
	lblock--
	sblock := lblock * kSmallBlockPerLargeBlock
	pointer := rs.pointerBlocks[lblock]
	remain := rank - rs.rankBlocks[lblock] + 1
	for ; sblock < uint64(len(rs.rankSmallBlocks)); sblock++ {
		rankSB := rs.rankSmallBlocks[sblock]
		if remain <= uint64(rankSB) {
			break
		}
		remain -= uint64(rankSB)
		pointer += uint64(kEnumCodeLength[rankSB])
	}
	rankSB := rs.rankSmallBlocks[sblock]
	code := getSlice(rs.bits, pointer, kEnumCodeLength[rankSB])
	return sblock*kSmallBlockSize + uint64(enumSelect1(code, rankSB, uint8(remain)))
}

func (rs RSDic) Select0(rank uint64) uint64 {
	if rank >= rs.zeroNum {
		return rs.num
	}
	if rank >= rs.zeroNum-rs.lastZeroNum {
		lastBlockRank := uint8(rank - (rs.zeroNum - rs.lastZeroNum))
		return rs.lastBlockInd() + uint64(selectRaw(^rs.lastBlock, lastBlockRank+1))
	}
	selectInd := rank / kSelectBlockSize
	lblock := rs.selectZeroInds[selectInd]
	for ; lblock < uint64(len(rs.rankBlocks)); lblock++ {
		if rank < lblock*kLargeBlockSize-rs.rankBlocks[lblock] {
			break
		}
	}
	lblock--
	sblock := lblock * kSmallBlockPerLargeBlock
	pointer := rs.pointerBlocks[lblock]
	remain := rank - lblock*kLargeBlockSize + rs.rankBlocks[lblock] + 1
	for ; sblock < uint64(len(rs.rankSmallBlocks)); sblock++ {
		rankSB := kSmallBlockSize - rs.rankSmallBlocks[sblock]
		if remain <= uint64(rankSB) {
			break
		}
		remain -= uint64(rankSB)
		pointer += uint64(kEnumCodeLength[rankSB])
	}
	rankSB := rs.rankSmallBlocks[sblock]
	code := getSlice(rs.bits, pointer, kEnumCodeLength[rankSB])
	return sblock*kSmallBlockSize + uint64(enumSelect0(code, rankSB, uint8(remain)))
}

// BitAndRank returns the (pos+1)-th bit (=b) and Rank(pos, b)
// Although this is equivalent to b := Bit(pos), r := Rank(pos, b),
// BitAndRank is faster.
func (rs RSDic) BitAndRank(pos uint64) (bool, uint64) {
	if rs.isLastBlock(pos) {
		offset := uint8(pos % kSmallBlockSize)
		bit := getBit(rs.lastBlock, offset)
		afterRank := uint64(popCount(rs.lastBlock >> offset))
		return bit, bitNum(rs.oneNum-afterRank, pos, bit)
	}
	lblock := pos / kLargeBlockSize
	pointer := rs.pointerBlocks[lblock]
	sblock := pos / kSmallBlockSize
	rank := rs.rankBlocks[lblock]
	for i := lblock * kSmallBlockPerLargeBlock; i < sblock; i++ {
		rankSB := rs.rankSmallBlocks[i]
		pointer += uint64(kEnumCodeLength[rankSB])
		rank += uint64(rankSB)
	}
	rankSB := rs.rankSmallBlocks[sblock]
	code := getSlice(rs.bits, pointer, kEnumCodeLength[rankSB])
	rank += uint64(enumRank(code, rankSB, uint8(pos%kSmallBlockSize)))
	bit := enumBit(code, rankSB, uint8(pos%kSmallBlockSize))
	return bit, bitNum(rank, pos, bit)
}

// AllocSize returns the allocated size in bytes.
func (rsd RSDic) AllocSize() int {
	return len(rsd.bits)*8 +
		len(rsd.pointerBlocks)*8 +
		len(rsd.rankBlocks)*8 +
		len(rsd.selectOneInds)*8 +
		len(rsd.selectZeroInds)*8 +
		len(rsd.rankSmallBlocks)*1
}

// MarshalBinary encodes the RSDic into a binary form and returns the result.
func (rsd RSDic) MarshalBinary() (out []byte, err error) {
	var bh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&out, &bh)
	err = enc.Encode(rsd.bits)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.pointerBlocks)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.rankBlocks)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.selectOneInds)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.selectZeroInds)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.rankSmallBlocks)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.num)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.oneNum)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.zeroNum)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.lastBlock)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.lastOneNum)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.lastZeroNum)
	if err != nil {
		return
	}
	err = enc.Encode(rsd.codeLen)
	if err != nil {
		return
	}
	return
}

// UnmarshalBinary decodes the RSDic from a binary from generated MarshalBinary
func (rsd *RSDic) UnmarshalBinary(in []byte) (err error) {
	var bh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(in, &bh)
	err = dec.Decode(&rsd.bits)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.pointerBlocks)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.rankBlocks)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.selectOneInds)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.selectZeroInds)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.rankSmallBlocks)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.num)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.oneNum)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.zeroNum)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.lastBlock)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.lastOneNum)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.lastZeroNum)
	if err != nil {
		return
	}
	err = dec.Decode(&rsd.codeLen)
	if err != nil {
		return
	}
	return nil
}

// New returns RSDic with a bit array of length 0.
func New() *RSDic {
	return &RSDic{
		bits:            make([]uint64, 0),
		pointerBlocks:   make([]uint64, 0),
		rankBlocks:      make([]uint64, 0),
		selectOneInds:   make([]uint64, 0),
		selectZeroInds:  make([]uint64, 0),
		rankSmallBlocks: make([]uint8, 0),
		num:             0,
		oneNum:          0,
		zeroNum:         0,
		lastBlock:       0,
		lastOneNum:      0,
		lastZeroNum:     0,
		codeLen:         0,
	}
}

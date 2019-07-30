package rsdic

import (
	"fmt"
)

func floor(num uint64, div uint64) uint64 {
	return (num + div - 1) / div
}

func decompose(x uint64, y uint64) (uint64, uint64) {
	return x / y, x % y
}

func setSlice(bits []uint64, pos uint64, codeLen uint8, val uint64) {
	if codeLen == 0 {
		return
	}
	block, offset := decompose(pos, kSmallBlockSize)
	bits[block] |= val << offset
	if offset+uint64(codeLen) > kSmallBlockSize {
		bits[block+1] |= (val >> (kSmallBlockSize - offset))
	}
}

func getBit(x uint64, pos uint8) bool {
	return ((x >> pos) & 1) == 1
}

func getSlice(bits []uint64, pos uint64, codeLen uint8) uint64 {
	if codeLen == 0 {
		return 0
	}
	block, offset := decompose(pos, kSmallBlockSize)
	ret := (bits[block] >> offset)
	if offset+uint64(codeLen) > kSmallBlockSize {
		ret |= (bits[block+1] << (kSmallBlockSize - offset))
	}
	if codeLen == 64 {
		return ret
	}
	return ret & ((1 << codeLen) - 1)
}

func bitNum(x uint64, n uint64, b bool) uint64 {
	if b {
		return x
	} else {
		return n - x
	}
}

func printBit(x uint64) {
	for i := 0; i < 64; i++ {
		fmt.Printf("%d", i%10)
	}
	fmt.Printf("\n")
	for i := uint8(0); i < 64; i++ {
		if getBit(x, i) {
			fmt.Printf("1")
		} else {
			fmt.Printf("0")
		}
	}
	fmt.Printf("\n")
}

func popCount(x uint64) uint8 {
	x = x - ((x & 0xAAAAAAAAAAAAAAAA) >> 1)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0F0F0F0F0F0F0F0F
	return uint8(x * 0x0101010101010101 >> 56)
}

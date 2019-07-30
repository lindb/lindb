RSDic
=====

rsdic is a Go package for rank/select dictionary supporting rank/select operations efficiently,
and space efficient for both sparse and dense bit arrays.
(Such data structures are also called as fully indexable dictionary in CS literatures (FID)).

Conceptually, rsdic represents a bit vector B[0...num), B[i] = 0 or 1,
and bits are provided PushBack operation (Thus RSDic supports dynamic addition).

All operations (Bit, Rank, Select, BitAndRank, RunZeros) are supported in O(1) time.
For example, rsdic can solve all above operation within 1 micro second for a bit vector of length 10^8 (100M bit).

RSDic combines the idea of on-the-fly decoding of enumrative code,
and utilization of rank/select samplings, which is discussed as a future work in the paper [1].

In RSDic, a bit vector is stored in compressed form (Note, we don't need to decode all at operations).
To achieve this, a bit vector is divided into small blocks of length 64, and each small block
is compressed using enumurative code. For example, a small block contains 10 ones
and 54 zeros will be compressed in 38 bits (See enumCode.go for detail).
This achieves not only its information theoretic bound, but also achieves more compression
if bits are clusterd.
rsdic stores information at most 1.3 bit per original bit including its indicies, and compress more if bit vector
is "compresible".

This Go version is based on the C++ implementation [2].
But this Go version supports PushBack so that it can support dynamic addition.

- [1] "Fast, Small, Simple Rank/Select on Bitmaps", Gonzalo Navarro and Eliana Providel, SEA 2012 [pdf](http://www.dcc.uchile.cl/~gnavarro/ps/sea12.1.pdf)
- [2] C++ version https://code.google.com/p/rsdic/ (the same author)

Usage
=====

	import "github.com/hillbig/rsdic"

	rsd := rsdic.New()
	rsd.PushBack(true)
	rsd.PushBack(false)
	rsd.PushBack(true)
	rsd.PushBack(true)
	// rsd = 1011
	fmt.Printf("%d %d %d\n", rsd.Num(), rsd.OneNum(), rsd.ZeroNum()) // 4 3 1

	// Bit(pos uint64) returns B[pos]
	fmt.Printf("%v\n", rsd.Bit(2)) // true

	// Rank(pos uint64, bit bool) returns the number of bit's in B[0...pos)
	fmt.Printf("%d %d\n", rsd.Rank(2, false), rsd.Rank(4, true)) // 1 3

	// Select(rank uint64, bit bool) returns the position of (rank+1)-th occurence of bit in B.
	oneNum := rsd.OneNum()
	for i := uint64(0); i < oneNum; i++ {
		fmt.Printf("%d:%d\n", rsd.Select(i, true))
	}
	// 0:0
	// 1:2
	// 2:3

	rsd.PushBack(false) // You can add anytime

	// Use MarshalBinary() and UnmarshalBinary() for serialize/deserialize RSDic.
	bytes, err := rsd.MarshalBinary()
	newrsd := rsdic.NewRSDic()
	err := newrsd.UnmarshalBinary(bytes)

	// Enjoy !


Benchmark
=========

- 1.7 GHz Intel Core i7
- OS X 10.9.2
- 8GB 1600 MHz DDR3
- go version go1.3 darwin/amd64

The results shows that RSDic operations require always
(almost) constant time with regard to the length and one's ratio.

	go test -bench=.

	// RSDic
	// A bit vector of length 10^6 with one's ratio = 0.5
	// Allocated size: 1.27 bits per original bit.
	BenchmarkDenseRSDicBit	20000000	       129 ns/op
	BenchmarkDenseRSDicRank	20000000	       127 ns/op
	BenchmarkDenseRSDicSelect	 5000000	       315 ns/op

	// RSDic
	// A bit vector of length 10^8 with one's ratio = 0.5
	BenchmarkBit	 5000000	       274 ns/op
	BenchmarkDenseRSDicRank	 5000000	       283 ns/op
	BenchmarkDenseRSDicSelect	 5000000	       551 ns/op
	BenchmarkSparseRSDicBit	10000000	       238 ns/op
	BenchmarkSparseRSDicRank	10000000	       279 ns/op
	BenchmarkSparseRSDicSelect	 5000000	       593 ns/op

	// RSDic
	// A bit vector of length 10^6 with one's ratios = 0.01
	// Allocated size: 0.32 bits per original bit.
	BenchmarkSparseRSDicBit	10000000	       190 ns/op
	BenchmarkSparseRSDicRank	10000000	       197 ns/op
	BenchmarkSparseRSDicSelect	 5000000	       495 ns/op

	// RSDic
	// A bit vector of length 10^8 with one's ratios = 0.01
	// Allocated size: 0.32 bits per original bit.
	BenchmarkSparseRSDicBit	10000000	       247 ns/op
	BenchmarkSparseRSDicRank	10000000	       264 ns/op
	BenchmarkSparseRSDicSelect	 5000000	       655 ns/op

	// []uint8 (for comparison)
	// A raw bit vector of length 10^6 with one's ratio = 0.5
	// Bit requires O(1) time and Rank, Select require O(n) time
	BenchmarkDenseRawBit	2000000000	         1.86 ns/op
	BenchmarkDenseRawRank	    5000	    629768 ns/op
	BenchmarkDenseRawSelect	     500	   5089415 ns/op

	// []uint8 (for comparison)
	// A raw bit vector of length 10^8 with one's ratio = 0.5
	BenchmarkDenseRawBit	2000000000	         1.82 ns/op
	BenchmarkDenseRawRank	      50	  67224404 ns/op
	BenchmarkDenseRawSelect	       5	 477187054 ns/op



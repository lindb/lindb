package hashers

import (
	"hash/fnv"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _testString = "abcdefghijklmnopqrstuvwxyz"

func Test_Fnv32a(t *testing.T) {
	for i := 0; i < 100; i++ {
		assert.Equal(t, Fnv32a(strconv.Itoa(i)), fnv32a(strconv.Itoa(i)))
	}
}

func Test_Fnv64a(t *testing.T) {
	for i := 0; i < 100; i++ {
		assert.Equal(t, Fnv64a(strconv.Itoa(i)), fnv64a(strconv.Itoa(i)))
	}
}

func Benchmark_Allocate_Free_FNV32a(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fnv32a(_testString)
	}
}

func Benchmark_fnv32a(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fnv32a(_testString)
	}
}

func Benchmark_Allocate_Free_FNV64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fnv64a(_testString)
	}
}

func Benchmark_fnv64a(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fnv64a(_testString)
	}
}

func fnv32a(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func fnv64a(s string) (hash uint64) {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

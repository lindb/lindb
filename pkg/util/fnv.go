package util

import (
	"hash/fnv"
)

// todo: @codingcrush, allocate-free

// Fnv32 returns a 32-bit FNV-1 hash of a string.
func Fnv32(s string) (hash uint32) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// Fnv64 returns a 64-bit FNV-1 hash of a string.
func Fnv64(s string) (hash uint64) {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

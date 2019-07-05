package hashers

import (
	"unsafe"
)

// copy from hash/fnv/fnv.go
const (
	offset32 uint32 = 2166136261
	offset64 uint64 = 14695981039346656037
	prime32  uint32 = 16777619
	prime64  uint64 = 1099511628211
)

// Fnv32a returns a 32-bit FNV-1a hash of a string.
func Fnv32a(s string) uint32 {
	hash := offset32
	for _, b := range *(*[]byte)(unsafe.Pointer(&s)) {
		hash = (hash ^ uint32(b)) * prime32
	}
	return hash
}

// Fnv64a returns a 64-bit FNV-1a hash of a string.
func Fnv64a(s string) uint64 {
	hash := offset64
	for _, b := range *(*[]byte)(unsafe.Pointer(&s)) {
		hash = (hash ^ uint64(b)) * prime64
	}
	return hash
}

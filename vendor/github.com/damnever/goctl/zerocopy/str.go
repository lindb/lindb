package zerocopy

import "unsafe"

// Ref: https://golang.org/src/strings/builder.go?s=1379:1412#L36

// UnsafeAtob converts string into bytes without copy.
func UnsafeAtob(s string) []byte {
	if s == "" {
		return []byte{}
	}
	return *(*[]byte)(unsafe.Pointer(&s))
}

// UnsafeBtoa converts bytes into string without copy.
func UnsafeBtoa(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

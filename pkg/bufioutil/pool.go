package bufioutil

import "sync"

// bufferPool is a set of temporary buffers for storing byte slice.
var bufferPool sync.Pool

// GetBuffer picks a buffer from the pool and then returns it.
func GetBuffer(length uint32) *[]byte {
	item := bufferPool.Get()
	if item == nil {
		buf := make([]byte, length)
		return &buf
	}
	buf := item.(*[]byte)
	// cap is smaller than required size.
	if uint32(cap(*buf)) < length {
		bufferPool.Put(item)
		buf := make([]byte, length)
		return &buf
	}
	*buf = (*buf)[:length]
	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *[]byte) {
	bufferPool.Put(buf)
}

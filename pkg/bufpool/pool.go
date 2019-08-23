package bufpool

import (
	"bytes"
	"sync"
)

// bufferPool is a set of temporary buffers for storing bytes.Buffer.
var bufferPool = &sync.Pool{New: func() interface{} {
	return &bytes.Buffer{}
}}

// GetBuffer picks a buffer from the pool and then returns it.
func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

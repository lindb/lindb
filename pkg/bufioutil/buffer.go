package bufioutil

import (
	"errors"
)

var errOutOfRange = errors.New("index out of range")

// A Buffer is a buffer manager for marshaling and unmarshaling
// protocol buffers.  It may be reused between invocations to
// reduce memory usage.  It is not necessary to use a Buffer;
// the global functions Marshal and Unmarshal create a
// temporary Buffer and are fine for most applications.
type Buffer struct {
	buf   []byte // encode/decode byte stream
	index int    // read point
}

// NewBuffer allocates a new Buffer and initializes its internal data to
// the contents of the argument slice.
func NewBuffer(e []byte) *Buffer {
	return &Buffer{buf: e}
}

// SetBuf replaces the internal buffer with the slice,
// ready for unmarshaling the contents of the slice.
func (p *Buffer) SetBuf(s []byte) {
	p.buf = s
	p.index = 0
}

func (p *Buffer) SetIdx(idx int) {
	p.index = idx
}

func (p *Buffer) GetByte() (b byte, err error) {
	if p.index >= len(p.buf) {
		err = errOutOfRange
		return
	}
	b = p.buf[p.index]
	p.index++
	return
}

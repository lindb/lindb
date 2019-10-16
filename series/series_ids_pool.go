package series

import "sync"

const ScanBufSize = 4096

// Uint32Pool is a singleton pool for reusing []uint32 with length as 4096
var Uint32Pool = _uint32Pool{
	Pool: sync.Pool{
		New: func() interface{} {
			item := make([]uint32, ScanBufSize)
			return &item
		}}}

type _uint32Pool struct {
	sync.Pool
}

func (p *_uint32Pool) Get() *[]uint32 {
	return p.Pool.Get().(*[]uint32)
}

func (p *_uint32Pool) Put(item *[]uint32) {
	p.Pool.Put(item)
}

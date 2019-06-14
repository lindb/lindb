package table

import (
	"github.com/eleme/lindb/pkg/mmap"
)

type Reader interface {
	Get(key uint32) ([]byte, error)
}

type StoreReader struct {
	buf []byte
	len int
}

func NewReader(path string) (*StoreReader, error) {
	buf, err := mmap.Map(path)
	if err != nil {
		return nil, err
	}
	r := &StoreReader{
		buf: buf,
		len: len(buf),
	}

	r.initialize()

	return r, nil
}

// initialize store reader, read index block(keys,offset etc.), then cache it
func (r *StoreReader) initialize() {
	//footer := r.buf[r.len-8:r.len]
}

func (r *StoreReader) Get(key uint32) ([]byte, error) {

	return nil, nil
}

// close store reader, release resource
func (r *StoreReader) Close() error {
	return mmap.Unmap(r.buf)
}

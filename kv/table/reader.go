package table

import (
	"encoding/binary"
	"fmt"

	"github.com/RoaringBitmap/roaring"

	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/mmap"
)

// Reader reads k/v pair from store file
type Reader interface {
	// Get returns value for giving key
	Get(key uint32) []byte
	// Iterator iterates over a store's key/value pairs in key order.
	Iterator() Iterator
	// Close closes reader, release related resources
	Close() error
}

// storeMMapReader is mmap store file reader
type storeMMapReader struct {
	path string
	buf  []byte
	len  int

	keys    *roaring.Bitmap
	offsets []int32
}

// newMMapStoreReader creates mmap store file reader
func newMMapStoreReader(path string) (Reader, error) {
	buf, err := mmap.Map(path)
	if err != nil {
		return nil, fmt.Errorf("create mmap store reader error:%s", err)
	}
	r := &storeMMapReader{
		path: path,
		buf:  buf,
		len:  len(buf),
		keys: roaring.New(),
	}

	if err := r.initialize(); err != nil {
		return nil, err
	}

	return r, nil
}

// initialize initializes store reader, reads index block(keys,offset etc.), then caches it
func (r *storeMMapReader) initialize() error {
	fileLen := r.len
	offsetOfOffset := int(r.readUint32(fileLen - 8))
	keyOfOffset := int(r.readUint32(fileLen - 4))

	if err := r.keys.UnmarshalBinary(r.readBytes(keyOfOffset)); err != nil {
		return fmt.Errorf("unmarshal keys data from file[%s] error:%s", r.path, err)
	}
	offset := r.readBytes(offsetOfOffset)
	d := encoding.NewDeltaBitPackingDecoder(&offset)

	for d.HasNext() {
		r.offsets = append(r.offsets, d.Next())
	}

	if len(r.offsets) != int(r.keys.GetCardinality()) {
		return fmt.Errorf("num. of keys != num. of offsets in file[%s]", r.path)
	}
	return nil
}

// Get return value for key, if not exist return nil
func (r *storeMMapReader) Get(key uint32) []byte {
	if !r.keys.Contains(key) {
		return nil
	}
	// bitmap data's index from 1, so idx=get index -1
	idx := r.keys.Rank(key)
	offset := r.offsets[idx-1]
	return r.readBytes(int(offset))
}

// Iterator iterates over a store's key/value pairs in key order.
func (r *storeMMapReader) Iterator() Iterator {
	return newMMapIterator(r)
}

// close store reader, release resource
func (r *storeMMapReader) Close() error {
	return mmap.Unmap(r.buf)
}

// readBytes reads bytes from buffer, read length+data format
func (r *storeMMapReader) readBytes(offset int) []byte {
	length := int(r.readUint32(offset))
	return r.buf[offset+4 : offset+4+length]
}

// readUint32 reads uint32 from buffer
func (r *storeMMapReader) readUint32(offset int) uint32 {
	return binary.BigEndian.Uint32(r.buf[offset : offset+4])
}

// storeMMapIterator iterates k/v pair using mmap store reader
type storeMMapIterator struct {
	reader *storeMMapReader
	keyIt  roaring.IntIterable

	idx int
}

// newMMapIterator creates store iterator using mmap store reader
func newMMapIterator(reader *storeMMapReader) Iterator {
	return &storeMMapIterator{
		reader: reader,
		keyIt:  reader.keys.Iterator(),
	}
}

// Next moves the iterator to the next key/value pair.
// It returns false if the iterator is exhausted.
func (it *storeMMapIterator) Next() bool {
	return it.keyIt.HasNext()
}

// Key returns the key of the current key/value pair
func (it *storeMMapIterator) Key() uint32 {
	key := it.keyIt.Next()
	return key
}

// Value returns the value of the current key/value pair
func (it *storeMMapIterator) Value() []byte {
	offset := it.reader.offsets[it.idx]
	it.idx++

	return it.reader.readBytes(int(offset))
}

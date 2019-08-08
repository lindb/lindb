package table

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
)

const (
	sstFileFooterSize = 1 + // entry length wrote by bufioutil
		4 + // posOfOffset(4)
		4 + // posOfKeys(4)
		1 + // version(1)
		8 // magicNumber(8)
	// footer-size, offset(1), keys(1)
	sstFileMinLength = sstFileFooterSize + 2
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
	path    string          // path of sst-file
	data    []byte          // mmaped  file content
	len     int             // length of the file
	keys    *roaring.Bitmap // bitmap of keys
	offsets []int32         // offset of values
}

// newMMapStoreReader creates mmap store file reader
func newMMapStoreReader(path string) (Reader, error) {
	data, err := fileutil.Map(path)
	if err != nil {
		return nil, fmt.Errorf("create mmap store reader error:%s", err)
	}
	if len(data) < sstFileMinLength {
		return nil, fmt.Errorf("length of sstfile:%s length is too short", path)
	}
	r := &storeMMapReader{
		path: path,
		data: data,
		len:  len(data),
		keys: roaring.New(),
	}

	if err := r.initialize(); err != nil {
		return nil, err
	}

	return r, nil
}

// initialize initializes store reader, reads index block(keys,offset etc.), then caches it
func (r *storeMMapReader) initialize() error {
	buf := r.readBytes(r.len - sstFileFooterSize)
	if (len(buf)) != sstFileFooterSize-1 {
		return fmt.Errorf("read sstfile:%s footer error", r.path)
	}
	// validate magic-number
	if binary.BigEndian.Uint64(buf[9:]) != magicNumberOffsetFile {
		return fmt.Errorf("verify magic-number of sstfile:%s failure", r.path)
	}
	posOfOffset := int(binary.BigEndian.Uint32(buf[:4]))
	posOfKeys := int(binary.BigEndian.Uint32(buf[4:8]))
	if err := r.keys.UnmarshalBinary(r.readBytes(posOfKeys)); err != nil {
		return fmt.Errorf("unmarshal keys data from file[%s] error:%s", r.path, err)
	}
	offset := r.readBytes(posOfOffset)
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
	return fileutil.Unmap(r.data)
}

// readBytes reads bytes from buffer, read length+data format
func (r *storeMMapReader) readBytes(offset int) []byte {
	length, err := binary.ReadUvarint(bytes.NewReader(r.data[offset:]))
	if err != nil {
		return nil
	}
	bytesCount := int(bufioutil.GetVariantLength(length))
	start := offset + bytesCount
	end := start + int(length)
	if end > len(r.data) {
		return nil
	}
	return r.data[start:end]
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

// HasNext returns if the iteration has more element.
// It returns false if the iterator is exhausted.
func (it *storeMMapIterator) HasNext() bool {
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

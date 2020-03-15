package table

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package table

// for testing
var (
	mapFunc     = fileutil.Map
	unmapFunc   = fileutil.Unmap
	uvarintFunc = binary.ReadUvarint
	uint64Func  = binary.LittleEndian.Uint64
)

// Reader represents reader which reads k/v pair from store file
type Reader interface {
	// Path returns the file path
	Path() string
	// Get returns value for giving key, if key not exist, return nil, false
	Get(key uint32) ([]byte, bool)
	// Iterator iterates over a store's key/value pairs in key order.
	Iterator() Iterator
	// Close closes reader, release related resources
	Close() error
}

// storeMMapReader represents mmap store file reader
type storeMMapReader struct {
	path    string                       // path of sst-file
	data    []byte                       // mmaped file content
	len     int                          // length of the file
	keys    *roaring.Bitmap              // bitmap of keys
	offsets *encoding.FixedOffsetDecoder // offset of values
}

// newMMapStoreReader creates mmap store file reader
func newMMapStoreReader(path string) (r Reader, err error) {
	data, err := mapFunc(path)
	defer func() {
		if err != nil && len(data) > 0 {
			// if init err and map data exist, need unmap it
			if e := unmapFunc(data); e != nil {
				tableLogger.Warn("unmap error when new store reader fail",
					logger.String("path", path), logger.Error(err))
			}
		}
	}()
	if err != nil {
		return
	}
	if len(data) < sstFileMinLength {
		err = fmt.Errorf("length of sstfile:%s length is too short", path)
		return
	}
	reader := &storeMMapReader{
		path: path,
		data: data,
		len:  len(data),
		keys: roaring.New(),
	}

	if err := reader.initialize(); err != nil {
		return nil, err
	}

	return reader, nil
}

// initialize initializes store reader, reads index block(keys,offset etc.), then caches it
func (r *storeMMapReader) initialize() error {
	buf := r.readBytes(r.len - sstFileFooterSize)
	if (len(buf)) != sstFileFooterSize-1 {
		return fmt.Errorf("read sstfile:%s footer error", r.path)
	}
	// validate magic-number
	if uint64Func(buf[9:]) != magicNumberOffsetFile {
		return fmt.Errorf("verify magic-number of sstfile:%s failure", r.path)
	}
	posOfOffset := int(binary.LittleEndian.Uint32(buf[:4]))
	posOfKeys := int(binary.LittleEndian.Uint32(buf[4:8]))
	if err := encoding.BitmapUnmarshal(r.keys, r.readBytes(posOfKeys)); err != nil {
		return fmt.Errorf("unmarshal keys data from file[%s] error:%s", r.path, err)
	}
	offset := r.readBytes(posOfOffset)
	r.offsets = encoding.NewFixedOffsetDecoder(offset)

	if r.offsets.Size() != int(r.keys.GetCardinality()) {
		return fmt.Errorf("num. of keys != num. of offsets in file[%s]", r.path)
	}
	return nil
}

// Path returns the file path
func (r *storeMMapReader) Path() string {
	return r.path
}

// Get return value for key, if not exist return nil,false
func (r *storeMMapReader) Get(key uint32) ([]byte, bool) {
	if !r.keys.Contains(key) {
		return nil, false
	}
	// bitmap data's index from 1, so idx=get index -1
	idx := r.keys.Rank(key)
	offset, _ := r.offsets.Get(int(idx) - 1)
	return r.readBytes(int(offset)), true
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
	length, err := uvarintFunc(bytes.NewReader(r.data[offset:]))
	if err != nil {
		return nil
	}
	bytesCount := stream.UvariantSize(length)
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
	offset, _ := it.reader.offsets.Get(it.idx)
	it.idx++
	return it.reader.readBytes(int(offset))
}

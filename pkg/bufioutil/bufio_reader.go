package bufioutil

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

const (
	defaultReadBufferSize = 256 * 1024 // 256KB
)

/*

The entries are encoded as follows,
of which the length flag is encoded with binary.variant
┌───────────────────┬───────────────────┐
│         Entry     │    Entry          │
├─────────┬─────────┼─────────┬─────────┤
│ Length  │ Content │ Length  │ Content │
│ uvariant│ N bytes │ uvariant│ N bytes │
└─────────┴─────────┴─────────┴─────────┘

mapping of len(content)(uint32) and bytes-count:
0 -> 0
[1, 2<<6) -> 1
[2<<6, 2<<13) -> 2
[2<<13, 2<<20) -> 3
[2<<20, 2<<27) -> 4

*/

// BufioReader read entries from a specified file by buffered I/O. Not thread-safe
type BufioReader interface {
	// Read reads a new entry's content.
	Read() (content []byte, err error)
	// Next reads from Reader and records the content and error.
	// returns true if not exceed the end of file.
	// returns false means there is no more data can be read from the reader.
	Next() bool
	// Reset switches the buffered reader to read from a new file:
	// open the new file; close the old opening file;
	// discards any buffered data and reset the states of bufio.Reader
	// reset the content-buffer and count.
	Reset(fileName string) error
	// Count returns the total size of bytes read successfully, including length cost.
	Count() int64
	// Size returns the total size of the file.
	Size() (int64, error)
	// Close closes the underlying file.
	Close() error
}

// bufioReader implements BufioReader.
type bufioReader struct {
	fileName string
	r        *bufio.Reader
	f        *os.File
	count    int64
	content  []byte
	err      error
}

// NewBufioReader returns a new BufioReader from fileName.
func NewBufioReader(fileName string) (BufioReader, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return &bufioReader{
		fileName: fileName,
		r:        bufio.NewReaderSize(f, defaultReadBufferSize),
		f:        f,
	}, nil
}

// GetVariantLength returns the length of variant-encoded-bytes.
func GetVariantLength(value uint64) int64 {
	i := uint8(1)
	for ; ; i++ {
		if value < 2<<(i*7-1) {
			break
		}
	}
	return int64(i)
}

// Next detects if there is data to read.
func (br *bufioReader) Next() bool {
	length, err := binary.ReadUvarint(br.r)
	if err == io.EOF {
		return false
	} else if err != nil {
		br.err = err
		return true
	}
	br.count += GetVariantLength(length)
	// expand the cap or not
	if uint64(cap(br.content)) < length {
		br.content = make([]byte, length)
	}
	// shrink the length
	br.content = br.content[:length]
	// read content
	n, err := io.ReadFull(br.r, br.content)
	if err == io.EOF {
		return false
	}
	br.err = err
	br.count += int64(n)
	return true
}

// Read returns content from next entry, the underlying buffer is reusable.
func (br *bufioReader) Read() (content []byte, err error) {
	// read length
	return br.content, br.err
}

// Reset switches the buffered reader to read from a new file.
func (br *bufioReader) Reset(fileName string) error {
	newF, err := os.Open(fileName)
	if err != nil {
		return err
	}
	if err = br.Close(); err != nil {
		return err
	}
	br.f = newF
	br.r.Reset(newF)
	br.count = 0
	// keep the underling array to avoid next memory allocation.
	br.content = br.content[:0]
	return nil
}

// Count returns the size of read bytes.
func (br *bufioReader) Count() int64 {
	return br.count
}

// Close closes the opened file.
func (br *bufioReader) Close() error {
	if br.f == nil {
		return nil
	}
	return br.f.Close()
}

// Size returns the stat of the file.
func (br *bufioReader) Size() (int64, error) {
	stat, err := br.f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

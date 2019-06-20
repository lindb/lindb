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

// The entries are encoded as follows:
// ┌────────────────┬─────────────────┐
// │         Entry  │    Entry        │
// ├──────┬─────────┬───────┬─────────┤
// │ Len  │ Content │  Len  │ Content │
// │4 byte│ N bytes │4 byte │ N bytes │
// └──────┴─────────┴───────┴─────────┘

// BufioReader read entries from a specified file by buffered I/O. Not thread-safe
type BufioReader interface {
	// Read reads a new entry's content.
	// After calling Read, the `eof` flag must be checked.
	// `eof = true`: there is no more data can be read from the reader.
	// `eof = false`: just check error and read the slice.
	Read() (eof bool, content []byte, err error)
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

// Read returns content from next entry, the underlying buffer is reusable.
func (br *bufioReader) Read() (eof bool, content []byte, err error) {
	var lenBuf [4]byte // buffer for store uint32
	// read length
	_, err = io.ReadFull(br.r, lenBuf[:])
	if err == io.EOF {
		return true, nil, err
	} else if err != nil {
		return false, nil, err
	}
	// got length
	length := binary.BigEndian.Uint32(lenBuf[:])
	// expand the cap or not
	if uint32(cap(br.content)) < length {
		br.content = make([]byte, length)
	}
	// shrink the length
	br.content = br.content[:length]
	// read content
	n, err := io.ReadFull(br.r, br.content)
	if err != nil {
		return false, nil, err
	}
	br.count += int64(n) + 4
	return false, br.content, nil
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

// Size return the stat of the file.
func (br *bufioReader) Size() (int64, error) {
	stat, err := br.f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

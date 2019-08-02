package segment

import (
	"sync/atomic"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
	"github.com/lindb/lindb/pkg/stream"
)

const (
	// indexItemSize is the size in bytes for a index item
	// (dataFileOffset(4 bytes int32) + dataFileLen(4 bytes int32)).
	indexItemSize = 4 + 4
)

// ErrExceedPageSize returns when appending data exceeding the mapped page size.
var ErrExceedPageSize = errors.New("exceed data page size")

// ErrOutOfRange returns when retrieving data with seq not in the segment range.
var ErrOutOfRange = errors.New("seq out of segment range")

// Segment represents a sequence of consecutive messages in a queue.
type Segment interface {
	// Begin returns sequence begin, inclusive.
	Begin() int64
	// End returns sequence end, exclusive.
	End() int64
	// Contains checks if sequence seq lies in segment sequence range [begin, end).
	Contains(seq int64) bool
	// Reads returns the message with sequence seq, if seq is not in sequence range, errorOutOfRange returns.
	Read(seq int64) ([]byte, error)
	// Append appends the message at the end of sequence,
	// if success returns the sequence to retrieve the message, otherwise returns the error.
	Append(message []byte) (int64, error)
	// Close releases the underlying resources.
	Close()
}

// segment implements Segment
type segment struct {
	// mmap page for index file
	indexPage page.MappedPage
	// mmap page for data file
	dataPage page.MappedPage
	// sequence begin, inclusive
	begin int64
	// sequence end, exclusive
	end int64
	// current dataOffset for data page
	dataOffset int
	// writer for index bytes
	indexWriter *stream.Binary
	// writer for data bytes
	dataWriter *stream.Binary
	logger     *logger.Logger
}

// NewSegment returns a Segment with provided index, data mmap page.
// Sequence range[begin, end) is used to reconstruct the state when loading from file.
func NewSegment(indexPage, dataPage page.MappedPage, begin, end int64) (Segment, error) {
	seg := &segment{
		indexPage: indexPage,
		dataPage:  dataPage,
		begin:     begin,
		end:       end,
		logger:    logger.GetLogger("pkg/segment"),
	}

	if err := seg.adjustOffset(); err != nil {
		return nil, err
	}

	return seg, nil
}

// Append appends the message at the end of sequence,
// if success returns the sequence to retrieve the message, otherwise returns the error.
func (seg *segment) Append(message []byte) (int64, error) {
	dataLen := len(message)
	if seg.dataOffset+dataLen > seg.dataPage.Size() {
		return 0, ErrExceedPageSize
	}

	// append message, preCheck ensures dataPage has enough space
	seg.dataWriter.PutBytes(message)

	// append index
	seg.indexWriter.PutInt32(int32(seg.dataOffset))
	seg.indexWriter.PutInt32(int32(dataLen))

	// advance dataOffset
	seg.dataOffset += dataLen

	seq := seg.End()
	atomic.AddInt64(&seg.end, 1)
	return seq, nil
}

// adjustOffset adjusts dataOffset by sequence range.
func (seg *segment) adjustOffset() error {
	// new segment
	if seg.begin == seg.end {
		seg.dataOffset = 0
		seg.indexWriter = stream.BinaryBufWriter(seg.indexPage.Buffer(0))
		seg.dataWriter = stream.BinaryBufWriter(seg.dataPage.Buffer(0))
		return nil
	}
	// restore segment from file
	dataOffset, dataLen, err := seg.calDataOffsetAndLen(seg.end - 1)
	if err != nil {
		return err
	}

	seg.dataOffset = dataOffset + dataLen
	seg.indexWriter = stream.BinaryBufWriter(seg.indexPage.Buffer(int(seg.end-seg.begin) * indexItemSize))
	seg.dataWriter = stream.BinaryBufWriter(seg.dataPage.Buffer(seg.dataOffset))
	return nil
}

// Begin returns sequence begin, inclusive.
func (seg *segment) Begin() int64 {
	return atomic.LoadInt64(&seg.begin)
}

// End returns sequence end, exclusive.
func (seg *segment) End() int64 {
	return atomic.LoadInt64(&seg.end)
}

// Reads returns the message with sequence seq, if seq is not in sequence range, errorOutOfRange returns.
func (seg *segment) Read(seq int64) ([]byte, error) {
	dataOffset, dataLen, err := seg.calDataOffsetAndLen(seq)
	if err != nil {
		return nil, err
	}
	return seg.dataPage.Data(dataOffset, dataLen), nil
}

// calDataOffsetAndLen returns the offset and length for message with sequence seq in data file.
func (seg *segment) calDataOffsetAndLen(seq int64) (int, int, error) {
	if !seg.Contains(seq) {
		return 0, 0, ErrOutOfRange
	}
	indexOffset := seq - seg.Begin()

	bin := stream.BinaryReader(seg.indexPage.Data(int(indexOffset)*indexItemSize, indexItemSize))

	return int(bin.ReadInt32()), int(bin.ReadInt32()), nil
}

// Contains checks if sequence seq lies in segment sequence range [begin, end).
func (seg *segment) Contains(seq int64) bool {
	return seg.Begin() <= seq && seq < seg.End()
}

// Close releases the underlying resources.
func (seg *segment) Close() {
	err := seg.indexPage.Close()
	if err != nil {
		seg.logger.Error("error close mmap file", zap.Error(err))
	}
	err = seg.dataPage.Close()
	if err != nil {
		seg.logger.Error("error close mmap file", zap.Error(err))
	}
}

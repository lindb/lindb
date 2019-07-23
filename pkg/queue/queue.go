package queue

import (
	"path"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/eleme/lindb/pkg/fileutil"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/queue/segment"
)

const (
	segmentDirName = "segment"
	metaFileName   = "queue.meta"
	// headSeq(int64), tailSeq(int64)
	queueMetaSize      = 8 + 8
	queueHeadSeqOffset = 0
	queueTailSeqOffset = 8
)

// ErrExceedingMessageSizeLimit returns when appending message exceeds the max size limit.
var ErrExceedingMessageSizeLimit = errors.New("message exceeds the max size limit")

// Queue represents a sequence of segments, new data is appended at headSeq.
// Segments with all message seqNum < tailSeq will be removed by ticker task.
type Queue interface {
	// Append appends data to the end of the queue,
	// if successes, returns the seq to retrieve data, otherwise returns err.
	// Concurrent unsafe.
	Append(message []byte) (int64, error)
	// GetSegment returns segment containing seq, returns error when not found.
	GetSegment(seq int64) (segment.Segment, error)
	// Size returns the total size of message.
	Size() int64
	// HeadSeq returns the head seq which stands for the latest read barrier.
	// New message is appended at head seq.
	HeadSeq() int64
	// TailSeq returns the tail seq which stands for the oldest read barrier.
	// Message with req less than tailSeq would be deleted at some point.
	TailSeq() int64
	// Ack advances the tailSeq to seq.
	Ack(seq int64)
	// Close closes the queue.
	Close()
}

// queue implements queue.
type queue struct {
	// dirPath for queue file
	dirPath string
	// the max size limit in bytes for data file
	dataFileSizeLimit int
	// segment factory
	fct segment.Factory
	// head segment for writing
	headSeg segment.Segment
	// queue meta, headSeq and tailSeq
	meta    Meta
	headSeq int64
	tailSeq int64
	// ticker to remove segments
	rmSegmentsTicker *time.Ticker
	logger           *logger.Logger
}

// NewQueue returns Queue based on dirPath, dataFileSizeLimit is used to limit the segment file size,
// removeTaskInterval specifics the interval to remove expired segments.
func NewQueue(dirPath string, dataFileSizeLimit int, removeTaskInterval time.Duration) (Queue, error) {
	if err := fileutil.MkDir(dirPath); err != nil {
		return nil, err
	}

	metaPath := path.Join(dirPath, metaFileName)
	meta, err := loadOrCreateMeta(metaPath)
	if err != nil {
		return nil, err
	}

	headSeq, tailSeq := meta.ReadInt64(queueHeadSeqOffset), meta.ReadInt64(queueTailSeqOffset)
	fct, err := segment.NewFactory(path.Join(dirPath, segmentDirName), dataFileSizeLimit, headSeq, tailSeq)
	if err != nil {
		return nil, err
	}

	headSeg, err := fct.GetSegment(headSeq)
	if err == segment.ErrSegmentNotFound {
		// only occurs when inits new queue
		headSeg, err = fct.NewSegment(headSeq)
	}

	if err != nil {
		return nil, err
	}

	q := &queue{
		dirPath:           dirPath,
		dataFileSizeLimit: dataFileSizeLimit,
		fct:               fct,
		headSeg:           headSeg,
		meta:              meta,
		headSeq:           headSeq,
		tailSeq:           tailSeq,
		rmSegmentsTicker:  time.NewTicker(removeTaskInterval),
		logger:            logger.GetLogger("pkg/queue"),
	}

	q.initRemoveSegmentsTask()

	return q, nil
}

// Append appends data to the end of the queue,
// if successes, returns the seq to retrieve data, otherwise returns err.
// Concurrent unsafe.
func (q *queue) Append(data []byte) (int64, error) {
	if len(data) > q.dataFileSizeLimit {
		return -1, ErrExceedingMessageSizeLimit
	}

	seq, err := q.headSeg.Append(data)
	if err == segment.ErrExceedPageSize {
		// rotate
		var newHeadSeg segment.Segment
		newHeadSeg, err = q.fct.NewSegment(q.HeadSeq())
		if err != nil {
			return -1, err
		}

		q.headSeg = newHeadSeg
		seq, err = q.headSeg.Append(data)
		if err != nil {
			return -1, err
		}
	}

	// assert
	if seq != q.headSeq {
		q.logger.Error("seq num and head seq not equal",
			logger.Int64("seq", seq), logger.Int64("headSeq", q.headSeq))
		panic("append error")
	}

	atomic.AddInt64(&q.headSeq, 1)

	q.meta.WriteInt64(queueHeadSeqOffset, q.HeadSeq())
	q.meta.WriteInt64(queueTailSeqOffset, q.TailSeq())

	return seq, nil
}

// GetSegment returns segment containing seq, returns error when not found.
func (q *queue) GetSegment(index int64) (segment.Segment, error) {
	return q.fct.GetSegment(index)
}

// Size returns the total size of message.
func (q *queue) Size() int64 {
	return q.HeadSeq() - q.TailSeq()
}

// HeadSeq returns the head seq which stands for the latest read barrier.
// New message is appended at head seq.
func (q *queue) HeadSeq() int64 {
	return atomic.LoadInt64(&q.headSeq)
}

// TailSeq returns the tail seq which stands for the oldest read barrier.
// Message with req less than tailSeq would be deleted at some point.
func (q *queue) TailSeq() int64 {
	return atomic.LoadInt64(&q.tailSeq)
}

func (q *queue) setTailSeq(seq int64) {
	atomic.StoreInt64(&q.tailSeq, seq)
}

// Ack advances the tailSeq to seq.
func (q *queue) Ack(seq int64) {
	if seq > q.TailSeq() && seq < q.HeadSeq() {
		q.setTailSeq(seq)
		q.meta.WriteInt64(queueTailSeqOffset, seq)
		if err := q.meta.Sync(); err != nil {
			q.logger.Error("sync queue meta error", logger.Error(err))
		}
	}
}

// Close closes the queue.
func (q *queue) Close() {
	if q.rmSegmentsTicker != nil {
		q.rmSegmentsTicker.Stop()
	}

	q.fct.Close()
	if err := q.meta.Close(); err != nil {
		q.logger.Error("close queue meta error", logger.Error(err))
	}
}

// RemoveSegments removes segments before TailSeq.
func (q *queue) initRemoveSegmentsTask() {
	go func() {
		q.logger.Info("initRemoveSegmentsTask")
		for range q.rmSegmentsTicker.C {
			if err := q.fct.RemoveSegments(q.TailSeq()); err != nil {
				q.logger.Error("remove segments error", logger.String("dirPath", q.dirPath), logger.Error(err))
			}
		}
	}()
}

// loadOrCreateMeta returns queueMeta, loads if metaFile exists, otherwise creates metaFile.
func loadOrCreateMeta(metaPath string) (Meta, error) {
	return NewMeta(metaPath, queueMetaSize)
}

package queue

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/segment"
)

//go:generate mockgen -source ./fanout.go -destination ./fanout_mock.go -package queue

const (
	fanOutDirName    = "fanOut"
	fanOutMetaSuffix = ".meta"
	// headSeq(int64), tailSeq(int64)
	fanOutMetaSize      = 8 + 8
	fanOutHeadSeqOffset = 0
	fanOutTailSeqOffset = 8
	// SeqNoNewMessageAvailable is the seqNum returned when no new message available
	SeqNoNewMessageAvailable = int64(-1)
)

// FanOutQueue represents a queue "produce once, consume multiple times".
// FanOut represents a individual consumer with own consume seq and ack seq.
type FanOutQueue interface {
	// Append appends data to tail of the queue,
	// if successes returns the seq to retrieve the data, otherwise returns err.
	// Concurrent unsafe.
	Append(data []byte) (int64, error)

	// GetOrCreateFanOut returns the FanOut if exists,
	// otherwise creates a new FanOut with consume seq and ack seq == queue tail seq.
	GetOrCreateFanOut(name string) (FanOut, error)

	// FanOutNames returns all fanOut names.
	FanOutNames() []string

	// Sync checks all the FanOuts tailSeqs, update the tailSeq as the smallest one.
	// Then syncs meta data to storage.
	Sync()

	// HeadSeq returns the headSeq which is the next seq for appending data.
	HeadSeq() int64

	// TailSeq returns the tailSeq which is the smallest seq among all the fanOut tailSeq.
	TailSeq() int64

	// GetSegment returns the segment contains seq.
	GetSegment(seq int64) (segment.Segment, error)

	// Close persists Seq meta, FanOut seq meta, release resources.
	Close()
}

// fanOutQueue implements FanOutQueue.
type fanOutQueue struct {
	// dir path for persistence file
	dirPath string
	// dir path for fanOut seqs
	fanOutDir string
	// underlying queue
	queue Queue
	// name -> FanOut mapping
	fanOutMap map[string]FanOut
	// lock for fanOutMap
	lock4map sync.RWMutex
	// 0 -> running, 1 -> closed
	closed int32
}

// NewFanOutQueue returns a FanOutQueue persisted in dirPath.
func NewFanOutQueue(dirPath string, dataFileSize int, removeTaskInterval time.Duration) (FanOutQueue, error) {
	// loads queue
	q, err := NewQueue(dirPath, dataFileSize, removeTaskInterval)
	if err != nil {
		return nil, err
	}

	foDir := path.Join(dirPath, fanOutDirName)
	if err := fileutil.MkDir(foDir); err != nil {
		return nil, err
	}

	fileNames, err := fileutil.ListDir(foDir)
	if err != nil {
		return nil, err
	}

	fq := &fanOutQueue{
		dirPath:   dirPath,
		fanOutDir: foDir,
		queue:     q,
		fanOutMap: make(map[string]FanOut),
	}

	// restores fanOut map
	for _, fn := range fileNames {
		if path.Ext(fn) == fanOutMetaSuffix {
			name := removeSuffix(fn)
			fo, err := NewFanOut(path.Join(foDir, fn), fq)
			if err != nil {
				return nil, err
			}
			fq.fanOutMap[name] = fo
		}
	}

	return fq, nil
}

// Append appends data to tail of the queue,
// if successes returns the seq to retrieve the data, otherwise returns err.
// Concurrent unsafe.
func (fq *fanOutQueue) Append(data []byte) (int64, error) {
	return fq.queue.Append(data)
}

// GetOrCreateFanOut returns the FanOut if exists,
// otherwise creates a new FanOut with consume seq and ack seq == queue tail seq.
func (fq *fanOutQueue) GetOrCreateFanOut(name string) (FanOut, error) {
	fq.lock4map.Lock()
	defer fq.lock4map.Unlock()

	fo, ok := fq.fanOutMap[name]
	if ok {
		return fo, nil
	}

	fo, err := NewFanOut(path.Join(fq.fanOutDir, name+fanOutMetaSuffix), fq)
	if err != nil {
		return nil, err
	}

	fq.fanOutMap[name] = fo

	return fo, nil
}

// FanOutNames returns all fanOut names
func (fq *fanOutQueue) FanOutNames() []string {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()
	names := make([]string, 0, len(fq.fanOutMap))
	for name := range fq.fanOutMap {
		names = append(names, name)
	}
	return names
}

// HeadSeq returns the headSeq which is the next seq for appending data.
func (fq *fanOutQueue) HeadSeq() int64 {
	return fq.queue.HeadSeq()
}

// TailSeq returns the tailSeq which is the smallest seq among all the fanOut tailSeq.
func (fq *fanOutQueue) TailSeq() int64 {
	return fq.queue.TailSeq()
}

// Sync checks all the FanOuts tailSeqs, update the tailSeq as the smallest one.
// Then syncs meta data to storage.
func (fq *fanOutQueue) Sync() {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()

	// no fanOuts
	if len(fq.fanOutMap) == 0 {
		return
	}

	// use the queue headSeq as the init value
	ackSeq := fq.queue.HeadSeq()
	for _, fo := range fq.fanOutMap {
		ts := fo.TailSeq()
		if ts < ackSeq {
			ackSeq = ts
		}
	}
	fq.queue.Ack(ackSeq)
}

// GetSegment returns the segment contains seq.
func (fq *fanOutQueue) GetSegment(index int64) (segment.Segment, error) {
	return fq.queue.GetSegment(index)
}

// Close persists Seq meta, FanOut seq meta, release resources.
func (fq *fanOutQueue) Close() {
	if atomic.CompareAndSwapInt32(&fq.closed, 0, 1) {
		fq.lock4map.RLock()
		defer fq.lock4map.RUnlock()

		for _, fo := range fq.fanOutMap {
			fo.Close()
		}

		fq.queue.Close()
	}
}

func removeSuffix(base string) string {
	return base[:strings.LastIndex(base, fanOutMetaSuffix)]
}

// FanOut represents a individual consumer with own consume seq and ack seq.
// The typical way to use FanOut is using a single go-routine to consume message,
// and using other go-routine to ack the messages which have been processed successfully.
type FanOut interface {
	// Name returns a unique name for FanOut in a FanOutQueue.
	Name() string
	// Consume returns the seq for the next data to consume.
	// If no new data is available, SeqNoNewMessageAvailable is returned.
	Consume() int64
	// SetHeadSeq sets the HeadSeq to seq, this is useful when re-consume message.
	// error returns when seq is invalidate(less than ackSeq or greater than the read barrier).
	SetHeadSeq(seq int64) error
	// Get retrieves the data for seq.
	// The seq must bu a valid sequence num returned by consume.
	// Call with seq less than ackSeq has undefined result.
	// Concurrent unsafe.
	Get(seq int64) ([]byte, error)
	// Ack mark the data processed with sequence less than or equals to ackSeq.
	Ack(ackSeq int64)
	// HeadSeq represents the next seq Consume returns.
	HeadSeq() int64
	// TailSeq returns the seq acked.
	TailSeq() int64
	// Pending returns the offset between FanOut HeadSeq and FanOutQueue HeadSeq.
	Pending() int64
	// Close persists  headSeq, tailSeq.
	Close()
}

// fanOut implements FanOut.
type fanOut struct {
	// unique name
	name string
	// underlying queue for retrieving data
	q FanOutQueue
	// persists meta
	meta Meta
	// the current segment for reading
	seg segment.Segment
	// consume seq
	headSeq int64
	// ack seq
	tailSeq int64
	// 0 -> running, 1 -> closed
	closed int32
	// lock to protect headSeq
	lock4headSeq sync.RWMutex
	logger       *logger.Logger
}

// NewFanOut builds a FanOut from metaPath.
func NewFanOut(metaPath string, q FanOutQueue) (FanOut, error) {
	meta, err := NewMeta(metaPath, fanOutMetaSize)
	if err != nil {
		return nil, err
	}

	base := path.Base(metaPath)
	name := removeSuffix(base)

	headSeq, tailSeq := meta.ReadInt64(fanOutHeadSeqOffset), meta.ReadInt64(fanOutTailSeqOffset)
	//reset to queue tailSeq
	if headSeq == 0 && tailSeq == 0 {
		tailSeq = q.TailSeq()
		headSeq = tailSeq
	}
	seg, err := q.GetSegment(headSeq)
	if err != nil {
		return nil, err
	}

	return &fanOut{
		name:    name,
		q:       q,
		meta:    meta,
		seg:     seg,
		headSeq: headSeq,
		tailSeq: tailSeq,
		logger:  logger.GetLogger("pkg/fanout"),
	}, nil
}

// Name returns a unique name for FanOut in a FanOutQueue.
func (f *fanOut) Name() string {
	return f.name
}

// Consume returns the seq for the next data to consume.
// If no new data is available, SeqNoNewMessageAvailable is returned.
func (f *fanOut) Consume() int64 {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	headSeq := f.headSeq
	if headSeq < f.q.HeadSeq() {
		f.headSeq = headSeq + 1
		return headSeq
	}
	return SeqNoNewMessageAvailable
}

// SetHeadSeq sets the HeadSeq to seq.
func (f *fanOut) SetHeadSeq(seq int64) error {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	hs := f.q.HeadSeq()
	ts := f.TailSeq()

	if seq > hs || seq < ts {
		return fmt.Errorf("set headSeq failed, %d not in the range [%d,%d]", seq, ts, hs)
	}

	f.headSeq = seq
	return nil
}

// Get retrieves the data for seq.
// The seq must bu a valid sequence num returned by consume.
// Call with seq less than ackSeq has undefined result.
// Concurrent unsafe.
func (f *fanOut) Get(seq int64) ([]byte, error) {
	bys, err := f.seg.Read(seq)
	if err == segment.ErrOutOfRange {
		var newSeg segment.Segment
		// try to locate segment
		newSeg, err = f.q.GetSegment(seq)
		if err != nil {
			return nil, err
		}
		f.seg = newSeg
		bys, err = f.seg.Read(seq)
	}
	return bys, err
}

// Ack mark the data with seq less than or equals to ackSeq.
func (f *fanOut) Ack(ackSeq int64) {
	ts := f.TailSeq()
	hs := f.HeadSeq()
	if ackSeq > ts && ackSeq < hs {
		f.setTailSeq(ackSeq)
		ts = ackSeq
		f.meta.WriteInt64(fanOutHeadSeqOffset, hs)
		f.meta.WriteInt64(fanOutTailSeqOffset, ts)
		if err := f.meta.Sync(); err != nil {
			f.logger.Error("sync fanOut meta error", logger.Error(err))
		}

		// update FanOutQueue ackSeq
		f.q.Sync()
	}
}

// HeadSeq represents the next seq Consume returns.
func (f *fanOut) HeadSeq() int64 {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()
	return f.headSeq
}

// TailSeq returns the seq acked.
func (f *fanOut) TailSeq() int64 {
	return atomic.LoadInt64(&f.tailSeq)
}

func (f *fanOut) setTailSeq(seq int64) {
	atomic.StoreInt64(&f.tailSeq, seq)
}

// Pending returns the offset between FanOut HeadSeq and FanOutQueue HeadSeq.
func (f *fanOut) Pending() int64 {
	fh := f.HeadSeq()
	qh := f.q.HeadSeq()
	return qh - fh
}

// Close persists  headSeq, tailSeq.
func (f *fanOut) Close() {
	if atomic.CompareAndSwapInt32(&f.closed, 0, 1) {
		if err := f.meta.Close(); err != nil {
			f.logger.Error("close fanOut meta error", logger.String("fanOut", f.name), logger.Error(err))
		}
	}
}

// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package queue

import (
	"fmt"
	"path"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source ./fanout.go -destination ./fanout_mock.go -package queue

// for testing
var (
	newQueueFunc  = NewQueue
	listDirFunc   = fileutil.ListDir
	newFanOutFunc = NewFanOut
)

// FanOutQueue represents a queue "produce once, consume multiple times".
// FanOut represents a individual consumer with own consume seq and ack seq.
type FanOutQueue interface {
	// Put puts data to tail of the queue,
	Put(data []byte) error
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
	//SetAppendSeq sets append seq(head/tail seq)
	SetAppendSeq(seq int64)
	// Close persists Seq meta, FanOut seq meta, release resources.
	Close()
	// get gets the message data by spec consume sequence
	get(sequence int64) ([]byte, error)
}

// fanOutQueue implements FanOutQueue.
type fanOutQueue struct {
	// dir path for persistence file
	dirPath string
	// dir path for storing fanOut consume sequence
	fanOutDir string
	// underlying queue
	queue Queue
	// name -> FanOut mapping
	fanOutMap map[string]FanOut
	// lock for fanOutMap
	lock4map sync.RWMutex
	// false -> running, true -> closed
	closed atomic.Bool
}

// NewFanOutQueue returns a FanOutQueue persisted in dirPath.
func NewFanOutQueue(dirPath string, dataSizeLimit int64, removeTaskInterval time.Duration) (FanOutQueue, error) {
	var err error

	fq := &fanOutQueue{
		dirPath:   dirPath,
		fanOutDir: path.Join(dirPath, fanOutDirName),
		fanOutMap: make(map[string]FanOut),
	}

	defer func() {
		if err != nil {
			// if initialize fanOut queue failure, need release the resource
			fq.Close()
		}
	}()

	// create underlying queue
	fq.queue, err = newQueueFunc(dirPath, dataSizeLimit, removeTaskInterval)
	if err != nil {
		return nil, err
	}
	// init fanOut sequence
	if err = fq.initFanOut(); err != nil {
		return nil, err
	}

	return fq, nil
}

// Put puts data to tail of the queue,
func (fq *fanOutQueue) Put(data []byte) error {
	return fq.queue.Put(data)
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

	fo, err := newFanOutFunc(fq.fanOutDir, name, fq)
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
	return fq.queue.HeadSeq() + 1
}

// TailSeq returns the tailSeq which is the smallest seq among all the fanOut tailSeq.
func (fq *fanOutQueue) TailSeq() int64 {
	return fq.queue.TailSeq()
}

// SetAppendSeq sets append seq(head/tail) underlying queue
func (fq *fanOutQueue) SetAppendSeq(seq int64) {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()

	fq.queue.SetAppendSeq(seq - 1)

	for _, fo := range fq.fanOutMap {
		fo.SetSeq(seq - 1)
	}
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

	if ackSeq >= 0 {
		fq.queue.Ack(ackSeq)
	}
}

// Close persists Seq meta, FanOut seq meta, release resources.
func (fq *fanOutQueue) Close() {
	if fq.closed.CAS(false, true) {
		fq.lock4map.RLock()
		defer fq.lock4map.RUnlock()

		for _, fo := range fq.fanOutMap {
			fo.Close()
		}

		if fq.queue != nil {
			fq.queue.Close()
		}
	}
}

// get gets the message data by spec consume sequence
func (fq *fanOutQueue) get(sequence int64) ([]byte, error) {
	return fq.queue.Get(sequence)
}

// initFanOut initializes exist fanOut consume sequence
func (fq *fanOutQueue) initFanOut() error {
	if err := mkDirFunc(fq.fanOutDir); err != nil {
		return err
	}

	fileNames, err := listDirFunc(fq.fanOutDir)
	if err != nil {
		return err
	}

	// load exist fanOut consume sequence
	for _, fn := range fileNames {
		fo, err := newFanOutFunc(fq.fanOutDir, fn, fq)
		if err != nil {
			return err
		}

		fq.fanOutMap[fn] = fo
	}

	return nil
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
	// Queue returns underlying queue.
	Queue() FanOutQueue
	//SetSeq sets append seq(head/tail seq)
	SetSeq(seq int64)
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
	q           FanOutQueue
	metaPageFct page.Factory
	// persists meta
	metaPage page.MappedPage
	// consume seq
	headSeq *atomic.Int64
	// ack seq
	tailSeq *atomic.Int64
	// false -> running, true -> closed
	closed atomic.Bool
	// lock to protect headSeq
	lock4headSeq sync.RWMutex
}

// NewFanOut builds a FanOut from metaPath.
func NewFanOut(parent, path string, q FanOutQueue) (FanOut, error) {
	name := filepath.Join(parent, path)
	var err error
	metaPageFct, err := newPageFactoryFunc(name, fanOutMetaSize)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if err1 := metaPageFct.Close(); err1 != nil {
				queueLogger.Error("close meta page factory when create fanOut",
					logger.String("fanOut", name), logger.Error(err))
			}
		}
	}()

	metaPage, err := metaPageFct.AcquirePage(metaPageIndex)
	if err != nil {
		return nil, err
	}

	headSeq := int64(metaPage.ReadUint64(fanOutHeadSeqOffset))
	tailSeq := int64(metaPage.ReadUint64(fanOutTailSeqOffset))
	// reset to queue tailSeq
	if headSeq == 0 && tailSeq == 0 {
		tailSeq = q.TailSeq()
		headSeq = tailSeq
	}

	return &fanOut{
		name:        name,
		q:           q,
		metaPageFct: metaPageFct,
		metaPage:    metaPage,
		headSeq:     atomic.NewInt64(headSeq),
		tailSeq:     atomic.NewInt64(tailSeq),
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

	headSeq := f.headSeq.Load() + 1
	if headSeq < f.q.HeadSeq() {
		f.headSeq.Store(headSeq)
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

	f.headSeq.Store(seq)

	return nil
}

// Get retrieves the data for seq.
// The seq must bu a valid sequence num returned by consume.
// Call with seq less than ackSeq has undefined result.
// Concurrent unsafe.
func (f *fanOut) Get(seq int64) ([]byte, error) {
	return f.q.get(seq)
}

// Queue returns underlying queue.
func (f *fanOut) Queue() FanOutQueue {
	return f.q
}

// Ack mark the data with seq less than or equals to ackSeq.
func (f *fanOut) Ack(ackSeq int64) {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	ts := f.TailSeq()
	hs := f.headSeq.Load()
	// In the initial condition, ts == 0, if the first ackSeq == 0, it would be ignore.
	// Since ack is always in batch mode and the following ack will ack the previous data, it's not big problem.
	if ackSeq > ts && ackSeq <= hs {
		f.setTailSeq(ackSeq)

		ts = ackSeq

		f.metaPage.PutUint64(uint64(hs), fanOutHeadSeqOffset)
		f.metaPage.PutUint64(uint64(ts), fanOutTailSeqOffset)

		if err := f.metaPage.Sync(); err != nil {
			queueLogger.Error("sync fanOut meta page error", logger.String("fanOut", f.name), logger.Error(err))
		}
	}
}

// HeadSeq represents the next seq Consume returns.
func (f *fanOut) HeadSeq() int64 {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	return f.headSeq.Load() + 1
}

// TailSeq returns the seq acked.
func (f *fanOut) TailSeq() int64 {
	return f.tailSeq.Load()
}

// AppendSeq returns the seq for next append.
func (f *fanOut) AppendSeq() int64 {
	return f.q.HeadSeq()
}

// SetAppendSeq sets append seq(head/tail) underlying queue
func (f *fanOut) SetSeq(seq int64) {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	f.headSeq.Store(seq)
	f.tailSeq.Store(seq)
	f.metaPage.PutUint64(uint64(seq), fanOutHeadSeqOffset)
	f.metaPage.PutUint64(uint64(seq-1), fanOutTailSeqOffset)
}

func (f *fanOut) setTailSeq(seq int64) {
	f.tailSeq.Store(seq)
}

// Pending returns the offset between FanOut HeadSeq and FanOutQueue HeadSeq.
func (f *fanOut) Pending() int64 {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	fh := f.HeadSeq() - 1
	qh := f.q.HeadSeq() - 1

	return qh - fh
}

// Close persists headSeq, tailSeq.
func (f *fanOut) Close() {
	if f.closed.CAS(false, true) {
		if err := f.metaPageFct.Close(); err != nil {
			queueLogger.Error("close fanOut meta error", logger.String("fanOut", f.name), logger.Error(err))
		}
	}
}

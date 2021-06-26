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
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source ./queue.go -destination ./queue_mock.go -package queue

// for testing
var (
	mkDirFunc          = fileutil.MkDirIfNotExist
	newPageFactoryFunc = page.NewFactory
)

// ErrExceedingMessageSizeLimit returns when appending message exceeds the max size limit.
var ErrExceedingMessageSizeLimit = errors.New("message exceeds the max page size limit")
var ErrOutOfSequenceRange = errors.New("out of sequence range")
var ErrExceedingTotalSizeLimit = errors.New("queue data size exceeds the max size limit")
var ErrMsgNotFound = errors.New("message not found")

var queueLogger = logger.GetLogger("queue", "FanOutQueue")

// Queue represents a sequence of segments, new data is appended at headSeq.
// Segments with all message seqNum < tailSeq will be removed by ticker task.
type Queue interface {
	// Put puts data to the end of the queue, if puts failure return err
	Put(message []byte) error
	// Get gets the message data at specific index
	Get(sequence int64) (message []byte, err error)
	// Size returns the total size of message.
	Size() int64
	// IsEmpty returns if queue is empty
	IsEmpty() bool
	// HeadSeq returns the head seq which stands for the latest read barrier.
	// New message is appended at head seq.
	HeadSeq() int64
	// TailSeq returns the tail seq which stands for the oldest read barrier.
	// Message with req less than tailSeq would be deleted at some point.
	TailSeq() int64
	// SetAppendSeq sets head/tail seq.
	SetAppendSeq(seq int64)
	// Ack advances the tailSeq to seq.
	Ack(seq int64)
	// Close closes the queue.
	Close()
}

// queue implements queue.
type queue struct {
	ctx    context.Context
	cancel context.CancelFunc
	// dirPath for queue file
	dirPath string
	// the max size limit in bytes for data file
	dataSizeLimit int64

	indexPageFct page.Factory // index page factory
	dataPageFct  page.Factory // data page factory
	metaPageFct  page.Factory // meta page factory

	// queue meta with headSeq and tailSeq
	metaPage page.MappedPage // meta buffer
	headSeq  atomic.Int64    // current written sequence
	tailSeq  atomic.Int64    // current acked sequence

	indexPage      page.MappedPage // index buffer
	indexPageIndex int64

	// message data write context
	dataPageIndex int64
	dataPage      page.MappedPage
	messageOffset int

	// ticker to remove acked data/index page
	removeTaskTicker *time.Ticker
	expireDataPage   atomic.Int64
	expireIndexPage  atomic.Int64
	closed           atomic.Bool
	rwMutex          sync.RWMutex
}

// NewQueue returns Queue based on dirPath, dataSizeLimit is used to limit the total data/index size,
// removeTaskInterval specifics the interval to remove expired segments.
func NewQueue(dirPath string, dataSizeLimit int64, removeTaskInterval time.Duration) (Queue, error) {
	var err error
	if err = mkDirFunc(dirPath); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	q := &queue{
		ctx:           ctx,
		cancel:        cancel,
		dirPath:       dirPath,
		dataSizeLimit: dataSizeLimit,
	}

	// if data size limit < default limit, need reset
	if q.dataSizeLimit < defaultDataSizeLimit {
		q.dataSizeLimit = defaultDataSizeLimit
	}

	defer func() {
		// if init queue failure, need release resource(like file/map file etc.)
		if err != nil {
			q.Close()
		}
	}()

	// init data page factory
	fct, err := newPageFactoryFunc(filepath.Join(dirPath, dataPath), dataPageSize)
	if err != nil {
		return nil, err
	}

	q.dataPageFct = fct

	// init index page factory
	fct, err = newPageFactoryFunc(filepath.Join(dirPath, indexPath), indexPageSize)
	if err != nil {
		return nil, err
	}

	q.indexPageFct = fct

	hasMeta := false
	if fileutil.Exist(filepath.Join(dirPath, metaPath, fmt.Sprintf("%d.bat", metaPageIndex))) {
		hasMeta = true
	}

	// init meta page factory
	fct, err = newPageFactoryFunc(filepath.Join(dirPath, metaPath), metaPageSize)
	if err != nil {
		return nil, err
	}

	q.metaPageFct = fct

	q.metaPage, err = q.metaPageFct.AcquirePage(metaPageIndex)
	if err != nil {
		return nil, err
	}

	if hasMeta {
		// initialize sequence
		q.initSequence()
	} else {
		q.headSeq.Store(-1)
		q.tailSeq.Store(-1)
		q.expireDataPage.Store(-1)
		q.expireIndexPage.Store(-1)
		// persist metadata
		q.metaPage.PutUint64(uint64(q.HeadSeq()), queueHeadSeqOffset)
		q.metaPage.PutUint64(uint64(q.TailSeq()), queueTailSeqOffset)
		q.metaPage.PutUint64(uint64(q.expireDataPage.Load()), queueExpireDataOffset)
		q.metaPage.PutUint64(uint64(q.expireIndexPage.Load()), queueExpireIndexOffset)
		if err = q.metaPage.Sync(); err != nil {
			return nil, err
		}
	}

	// initialize data page indexes
	if err = q.initDataPageIndex(); err != nil {
		return nil, err
	}

	q.removeTaskTicker = time.NewTicker(removeTaskInterval)
	q.initRemoveTask()

	return q, nil
}

// Put puts data to the end of the queue, if puts failure return err
func (q *queue) Put(data []byte) error {
	dataLength := len(data)
	if dataLength > dataPageSize {
		// if message size > data page size, return err
		return ErrExceedingMessageSizeLimit
	}

	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()

	dataPage, offset, err := q.alloc(dataLength)
	if err != nil {
		return err
	}
	// write message data
	dataPage.WriteBytes(data, offset)

	return nil
}

// Get gets the message data at specific index
func (q *queue) Get(sequence int64) (data []byte, err error) {
	if err = q.validateSequence(sequence); err != nil {
		return
	}

	indexPageID := sequence / indexItemsPerPage
	indexPage, ok := q.indexPageFct.GetPage(indexPageID)

	if !ok {
		return nil, ErrMsgNotFound
	}

	// calculate index offset of previous sequence
	indexOffset := int((sequence % indexItemsPerPage) * indexItemLength)
	dataPageID := int64(indexPage.ReadUint64(indexOffset + queueDataPageIndexOffset))

	dataPage, ok := q.dataPageFct.GetPage(dataPageID)
	if !ok {
		return nil, ErrMsgNotFound
	}

	messageOffset := int(indexPage.ReadUint32(indexOffset + messageOffsetOffset))
	messageLength := int(indexPage.ReadUint32(indexOffset + messageLengthOffset))

	return dataPage.ReadBytes(messageOffset, messageLength), nil
}

// Size returns the total size of message.
func (q *queue) Size() int64 {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	return q.HeadSeq() - q.TailSeq()
}

// HeadSeq returns the head seq which stands for the latest read barrier.
// New message is appended at head seq.
func (q *queue) HeadSeq() int64 {
	return q.headSeq.Load()
}

// TailSeq returns the tail seq which stands for the oldest read barrier.
// Message with req less than tailSeq would be deleted at some point.
func (q *queue) TailSeq() int64 {
	return q.tailSeq.Load()
}

// SetAppendSeq sets head/tail seq.
func (q *queue) SetAppendSeq(seq int64) {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	indexPageIndex := seq / indexItemsPerPage
	if indexPageIndex != q.indexPageIndex {
		// sync previous data page
		if err := q.indexPage.Sync(); err != nil {
			queueLogger.Error("sync index page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
		}
		indexPage, err := q.indexPageFct.AcquirePage(indexPageIndex)
		if err != nil {
			queueLogger.Error("sync index page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
			return
		}

		q.indexPage = indexPage
		q.indexPageIndex++
	}

	head := seq

	q.headSeq.Store(head)
	q.metaPage.PutUint64(uint64(head), queueHeadSeqOffset)
	tail := head - 1
	q.tailSeq.Store(tail)
	q.metaPage.PutUint64(uint64(tail), queueTailSeqOffset)
	q.metaPage.PutUint64(uint64(q.HeadSeq()), queueHeadSeqOffset)
	q.metaPage.PutUint64(uint64(q.TailSeq()), queueTailSeqOffset)
	if err := q.metaPage.Sync(); err != nil {
		queueLogger.Error("sync queue meta page error, when set append seq",
			logger.String("path", q.dirPath), logger.Error(err))
	}
}

// Ack advances the tailSeq to seq.
func (q *queue) Ack(seq int64) {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	if seq > q.TailSeq() && seq <= q.HeadSeq() {
		q.tailSeq.Store(seq)
		q.metaPage.PutUint64(uint64(seq), queueTailSeqOffset)

		if err := q.metaPage.Sync(); err != nil {
			queueLogger.Error("sync queue meta page error, when ack seq",
				logger.String("path", q.dirPath), logger.Error(err))
		}
	}
}

// IsEmpty returns if queue is empty
func (q *queue) IsEmpty() bool {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	return q.HeadSeq() == q.TailSeq()
}

// Close closes the queue.
func (q *queue) Close() {
	if q.closed.CAS(false, true) {
		q.rwMutex.RLock()
		defer q.rwMutex.RUnlock()

		q.cancel()
		if q.removeTaskTicker != nil {
			q.removeTaskTicker.Stop()
		}

		if q.dataPageFct != nil {
			if err := q.dataPageFct.Close(); err != nil {
				queueLogger.Error("close data page factory error",
					logger.String("queue", q.dirPath), logger.Error(err))
			}
		}

		if q.indexPageFct != nil {
			if err := q.indexPageFct.Close(); err != nil {
				queueLogger.Error("close index page factory error",
					logger.String("queue", q.dirPath), logger.Error(err))
			}
		}

		if q.metaPageFct != nil {
			if err := q.metaPageFct.Close(); err != nil {
				queueLogger.Error("close meta page factory error",
					logger.String("queue", q.dirPath), logger.Error(err))
			}
		}
	}
}

// RemoveSegments removes segments before TailSeq.
func (q *queue) initRemoveTask() {
	go func() {
		defer queueLogger.Info("exist remove ack queue task")
		queueLogger.Info("start remove ack queue task")
		for {
			select {
			case <-q.removeTaskTicker.C:
				q.removeExpirePage()
			case <-q.ctx.Done():
				return
			}
		}
	}()
}

func (q *queue) removeExpirePage() {
	ackSeq := q.TailSeq() // get current acked sequence
	if ackSeq < 0 {
		return
	}
	indexPageID := ackSeq / indexItemsPerPage
	indexPage, ok := q.indexPageFct.GetPage(indexPageID)
	if !ok {
		return
	}
	// calculate index offset of ack sequence
	indexOffset := int((ackSeq % indexItemsPerPage) * indexItemLength)
	dataPageID := int64(indexPage.ReadUint64(indexOffset + queueDataPageIndexOffset))
	lastDataPageID := q.expireDataPage.Load()
	for i := lastDataPageID + 1; i < dataPageID; i++ {
		if err := q.dataPageFct.ReleasePage(i); err != nil {
			queueLogger.Error("remove expire data page error",
				logger.String("queue", q.dirPath), logger.Any("page", i), logger.Error(err))
			break
		}
		queueLogger.Info("remove expire data page",
			logger.String("queue", q.dirPath), logger.Any("page", i))
		q.expireDataPage.Store(i)
		q.metaPage.PutUint64(uint64(q.expireDataPage.Load()), queueExpireDataOffset)
	}
	lastIndexPageID := q.expireIndexPage.Load()
	for i := lastIndexPageID + 1; i < indexPageID; i++ {
		if err := q.indexPageFct.ReleasePage(i); err != nil {
			queueLogger.Error("remove expire index page error",
				logger.String("queue", q.dirPath), logger.Any("page", i), logger.Error(err))
			break
		}
		queueLogger.Info("remove expire index page",
			logger.String("queue", q.dirPath), logger.Any("page", i))
		q.expireIndexPage.Store(i)
		q.metaPage.PutUint64(uint64(q.expireIndexPage.Load()), queueExpireIndexOffset)
	}

	if err := q.metaPage.Sync(); err != nil {
		queueLogger.Error("sync meta page error when do expire page",
			logger.String("queue", q.dirPath), logger.Error(err))
	}
}

// alloc allocates the data page and offset for message writing
func (q *queue) alloc(dataLen int) (dataPage page.MappedPage, offset int, err error) {
	// prepare the data pointer
	if q.messageOffset+dataLen > dataPageSize {
		// check size limit before data page acquire
		if err = q.checkDataSize(); err != nil {
			return nil, 0, err
		}
		// sync previous data page
		if err = q.dataPage.Sync(); err != nil {
			queueLogger.Error("sync data page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
		}
		// not enough space in current data page, need create new page
		dataPage, err := q.dataPageFct.AcquirePage(q.dataPageIndex + 1)
		if err != nil {
			return nil, 0, err
		}

		q.dataPage = dataPage
		q.dataPageIndex++
		q.messageOffset = 0 // need reset message offset for new data page
	}

	seq := q.headSeq.Load() + 1
	indexPageIndex := seq / indexItemsPerPage
	if indexPageIndex != q.indexPageIndex {
		// check size limit before index page acquire
		if err = q.checkDataSize(); err != nil {
			return nil, 0, err
		}
		// sync previous data page
		if err = q.indexPage.Sync(); err != nil {
			queueLogger.Error("sync index page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
		}
		indexPage, err := q.indexPageFct.AcquirePage(indexPageIndex)
		if err != nil {
			return nil, 0, err
		}

		q.indexPage = indexPage
		q.indexPageIndex++
	}
	// advance dataOffset
	messageOffset := q.messageOffset

	// save index data
	indexOffset := int((seq % indexItemsPerPage) * indexItemLength)
	q.indexPage.PutUint64(uint64(q.dataPageIndex), indexOffset+queueDataPageIndexOffset)
	q.indexPage.PutUint32(uint32(messageOffset), indexOffset+messageOffsetOffset)
	q.indexPage.PutUint32(uint32(dataLen), indexOffset+messageLengthOffset)

	q.messageOffset += dataLen

	// save metadata
	q.headSeq.Store(seq)
	q.metaPage.PutUint64(uint64(q.HeadSeq()), queueHeadSeqOffset)
	q.metaPage.PutUint64(uint64(q.TailSeq()), queueTailSeqOffset)

	return q.dataPage, messageOffset, nil
}

// initSequence initializes head/tail from the meta data
func (q *queue) initSequence() {
	q.headSeq.Store(int64(q.metaPage.ReadUint64(queueHeadSeqOffset)))
	q.tailSeq.Store(int64(q.metaPage.ReadUint64(queueTailSeqOffset)))
	q.expireDataPage.Store(int64(q.metaPage.ReadUint64(queueExpireDataOffset)))
	q.expireIndexPage.Store(int64(q.metaPage.ReadUint64(queueExpireIndexOffset)))
}

// initDataPageIndex finds out data page head index and message offset
func (q *queue) initDataPageIndex() (err error) {
	if q.IsEmpty() {
		// if queue is empty, start with new empty queue
		q.dataPageIndex = 0
		q.messageOffset = 0

		if q.dataPage, err = q.dataPageFct.AcquirePage(0); err != nil {
			return err
		}

		if q.indexPage, err = q.indexPageFct.AcquirePage(0); err != nil {
			return err
		}

		return nil
	}

	previousSeq := q.HeadSeq() // get previous sequence
	q.indexPageIndex = previousSeq / indexItemsPerPage

	if q.indexPage, err = q.indexPageFct.AcquirePage(q.indexPageIndex); err != nil {
		return err
	}

	// calculate index offset of previous sequence
	indexOffset := int((previousSeq % indexItemsPerPage) * indexItemLength)
	q.dataPageIndex = int64(q.indexPage.ReadUint64(indexOffset + queueDataPageIndexOffset))
	previousMessageOffset := q.indexPage.ReadUint32(indexOffset + messageOffsetOffset)
	previousMessageLength := q.indexPage.ReadUint32(indexOffset + messageLengthOffset)
	// calculate next message offset
	q.messageOffset = int(previousMessageOffset + previousMessageLength)

	if q.dataPage, err = q.dataPageFct.AcquirePage(q.dataPageIndex); err != nil {
		return err
	}

	return nil
}

// validateSequence validates the sequence if in range
func (q *queue) validateSequence(sequence int64) error {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	if sequence <= q.TailSeq() || sequence > q.HeadSeq() {
		return ErrOutOfSequenceRange
	}

	return nil
}

// checkDataSize checks the data size if exceeds the size limit
func (q *queue) checkDataSize() error {
	if q.dataPageFct.Size()+q.indexPageFct.Size() > q.dataSizeLimit {
		return ErrExceedingTotalSizeLimit
	}
	return nil
}

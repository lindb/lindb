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
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source ./queue.go -destination ./queue_mock.go -package queue

// for testing
var (
	mkDirFunc          = fileutil.MkDirIfNotExist
	newPageFactoryFunc = page.NewFactory
	existFunc          = fileutil.Exist
)

var (
	// ErrExceedingMessageSizeLimit returns when appending message exceeds the max size limit.
	ErrExceedingMessageSizeLimit = errors.New("message exceeds the max page size limit")
	// ErrOutOfSequenceRange returns sequence out of range.
	ErrOutOfSequenceRange = errors.New("out of sequence range")
	// ErrExceedingTotalSizeLimit returns total size limit.
	ErrExceedingTotalSizeLimit = errors.New("queue data size exceeds the max size limit")
	// ErrMsgNotFound returns message not found.
	ErrMsgNotFound = errors.New("message not found")
)

var queueLogger = logger.GetLogger("Queue", "FanOutQueue")

// Queue represents a sequence of segments, new data is appended at append sequence.
// Segments with all message will be removed by gc which sequence < acknowledged sequence.
type Queue interface {
	// Put puts data to the end of the queue, if puts failure return err.
	Put(message []byte) error
	// Get gets the message data at specific index.
	Get(sequence int64) (message []byte, err error)
	// AppendedSeq returns the written sequence which stands for the latest write barrier.
	// New message is appended at append sequence.
	AppendedSeq() int64
	// SetAppendedSeq sets appended sequence.
	SetAppendedSeq(seq int64)
	// AcknowledgedSeq returns the acknowledged sequence which stands for the oldest read barrier.
	// Message with req less than acknowledged sequence would be deleted at some point.
	AcknowledgedSeq() int64
	// SetAcknowledgedSeq sets acknowledged sequence.
	SetAcknowledgedSeq(seq int64)
	// NotEmpty checks queue if empty, waiting until new data written.
	NotEmpty(consumeHead int64, checkClosed func() bool) bool
	// Signal signals waiting consumers.
	Signal()
	// GC removes all message which sequence <= acknowledged sequence.
	GC()
	// Close closes the queue.
	Close()
}

// queue implements queue.
type queue struct {
	indexPage       page.MappedPage // index buffer
	dataPage        page.MappedPage // data buffer
	indexPageFct    page.Factory    // index page factory
	dataPageFct     page.Factory    // data page factory
	metaPage        page.MappedPage // meta buffer
	metaPageFct     page.Factory    // meta page factory
	notEmpty        *sync.Cond      // not empty condition
	rwMutex         *sync.RWMutex
	dirPath         string       // path for queue file
	appendedSeq     atomic.Int64 // current written sequence
	dataPageIndex   int64
	indexPageIndex  int64
	messageOffset   int
	closed          atomic.Bool
	acknowledgedSeq atomic.Int64 // acknowledged sequence
	pageSize        int64        // data page size
}

// NewQueue returns Queue based on dirPath, pageSize is used to limit the data page size,
func NewQueue(dirPath string, pageSize int64) (Queue, error) {
	if err := mkDirFunc(dirPath); err != nil {
		return nil, err
	}
	lock := &sync.RWMutex{}
	q := &queue{
		dirPath:  dirPath,
		pageSize: pageSize,
		rwMutex:  lock,
		notEmpty: sync.NewCond(lock),
	}

	// if data size limit < default limit, need reset
	if q.pageSize < dataPageSize {
		q.pageSize = dataPageSize
	}

	var err error

	defer func() {
		// if init queue failure, need release resource(like file/map file etc.)
		if err != nil {
			q.Close()
		}
	}()

	// init data page factory
	var dataPageFct page.Factory
	dataPageFct, err = newPageFactoryFunc(filepath.Join(dirPath, dataPath), int(q.pageSize))
	if err != nil {
		return nil, err
	}

	q.dataPageFct = dataPageFct

	// init index page factory
	var indexPageFct page.Factory
	indexPageFct, err = newPageFactoryFunc(filepath.Join(dirPath, indexPath), indexPageSize)
	if err != nil {
		return nil, err
	}

	q.indexPageFct = indexPageFct

	hasMeta := fileutil.Exist(filepath.Join(dirPath, metaPath, fmt.Sprintf("%d.bat", metaPageIndex)))

	// init meta page factory
	var metaPageFct page.Factory
	metaPageFct, err = newPageFactoryFunc(filepath.Join(dirPath, metaPath), metaPageSize)
	if err != nil {
		return nil, err
	}

	q.metaPageFct = metaPageFct

	q.metaPage, err = q.metaPageFct.AcquirePage(metaPageIndex)
	if err != nil {
		return nil, err
	}

	if hasMeta {
		// initialize sequence
		q.initSequence()
	} else {
		q.appendedSeq.Store(SeqNoNewMessageAvailable)
		q.acknowledgedSeq.Store(SeqNoNewMessageAvailable)

		// persist metadata
		q.metaPage.PutUint64(uint64(q.AppendedSeq()), queueAppendedSeqOffset)
		q.metaPage.PutUint64(uint64(q.AcknowledgedSeq()), queueAcknowledgedSeqOffset)

		err = q.metaPage.Sync()
		if err != nil {
			return nil, err
		}
	}

	// initialize data page indexes
	err = q.initDataPageIndex()
	if err != nil {
		return nil, err
	}
	return q, nil
}

// Put puts data to the end of the queue, if puts failure return err
func (q *queue) Put(data []byte) error {
	dataLength := len(data)
	if dataLength > dataPageSize {
		// if message size > data page size, return err
		return ErrExceedingMessageSizeLimit
	}

	dataPageIndex, dataPage, offset, err := q.alloc(dataLength)
	if err != nil {
		return err
	}

	// write message data
	dataPage.WriteBytes(data, offset)

	// persist metadata of message after write data
	return q.persistMetaOfMessage(dataPageIndex, dataLength, offset)
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

// AppendedSeq returns the written sequence which stands for the latest write barrier.
// New message is appended at append sequence.
func (q *queue) AppendedSeq() int64 {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	return q.appendedSeq.Load()
}

// SetAppendedSeq sets appended sequence.
func (q *queue) SetAppendedSeq(seq int64) {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()

	q.appendedSeq.Store(seq)
	q.acknowledgedSeq.Store(seq)

	q.metaPage.PutUint64(uint64(q.appendedSeq.Load()), queueAppendedSeqOffset)
	q.metaPage.PutUint64(uint64(q.acknowledgedSeq.Load()), queueAcknowledgedSeqOffset)
	if err := q.metaPage.Sync(); err != nil {
		queueLogger.Error("sync queue meta page error, when set append seq",
			logger.String("path", q.dirPath), logger.Error(err))
	}
}

// AcknowledgedSeq returns the acknowledged sequence which stands for the oldest read barrier.
// Message with req less than acknowledged sequence would be deleted at some point.
func (q *queue) AcknowledgedSeq() int64 {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()

	return q.acknowledgedSeq.Load()
}

// SetAcknowledgedSeq sets acknowledged sequence.
func (q *queue) SetAcknowledgedSeq(seq int64) {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()

	if seq > q.acknowledgedSeq.Load() && seq <= q.appendedSeq.Load() {
		q.acknowledgedSeq.Store(seq)
		q.metaPage.PutUint64(uint64(seq), queueAcknowledgedSeqOffset)

		if err := q.metaPage.Sync(); err != nil {
			queueLogger.Error("sync queue meta page error, when ack seq",
				logger.String("path", q.dirPath), logger.Error(err))
		}
	}
}

// NotEmpty checks queue if empty, waiting until new data written.
func (q *queue) NotEmpty(consumeHead int64, checkClosed func() bool) bool {
	q.notEmpty.L.Lock()
	for consumeHead > q.appendedSeq.Load() && !q.closed.Load() && !checkClosed() {
		q.notEmpty.Wait()
	}
	q.notEmpty.L.Unlock()

	return !q.closed.Load() && !checkClosed()
}

// Signal signals waiting consumers.
func (q *queue) Signal() {
	q.notEmpty.Broadcast()
}

// Close closes the queue.
func (q *queue) Close() {
	if q.closed.CompareAndSwap(false, true) {
		q.rwMutex.RLock()
		defer q.rwMutex.RUnlock()

		q.notEmpty.Broadcast()

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

// GC removes all message which sequence < acknowledged sequence.
func (q *queue) GC() {
	// get current acknowledged sequence.
	ackSeq := q.AcknowledgedSeq()
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

	q.dataPageFct.TruncatePages(dataPageID)
	q.indexPageFct.TruncatePages(indexPageID)
}

// alloc allocates the data page and offset for message writing
func (q *queue) alloc(dataLen int) (dataPageIndex int64, dataPage page.MappedPage, offset int, err error) {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()

	// prepare the data pointer
	if q.messageOffset+dataLen > dataPageSize {
		// sync previous data page
		if err := q.dataPage.Sync(); err != nil {
			queueLogger.Error("sync data page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
		}
		nextDataPageIndex := q.dataPageIndex + 1
		// not enough space in current data page, need create new page
		dataPage, err := q.dataPageFct.AcquirePage(nextDataPageIndex)
		if err != nil {
			return 0, nil, 0, err
		}

		q.dataPage = dataPage
		q.dataPageIndex = nextDataPageIndex
		q.messageOffset = 0 // need reset message offset for new data page
	}
	// advance dataOffset
	messageOffset := q.messageOffset
	q.messageOffset += dataLen // set next message offset
	return q.dataPageIndex, q.dataPage, messageOffset, nil
}

// persistMetaOfMessage persists metadata of message after write data
func (q *queue) persistMetaOfMessage(dataPageIndex int64, dataLen, messageOffset int) error {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()

	seq := q.appendedSeq.Load() + 1 // append sequence
	indexPageIndex := seq / indexItemsPerPage
	if indexPageIndex != q.indexPageIndex {
		// sync previous data page
		if err := q.indexPage.Sync(); err != nil {
			queueLogger.Error("sync index page err when alloc",
				logger.String("queue", q.dirPath), logger.Error(err))
		}
		indexPage, err := q.indexPageFct.AcquirePage(indexPageIndex)
		if err != nil {
			return err
		}

		q.indexPage = indexPage
		q.indexPageIndex = indexPageIndex
	}

	// save index data
	indexOffset := int((seq % indexItemsPerPage) * indexItemLength)
	q.indexPage.PutUint64(uint64(dataPageIndex), indexOffset+queueDataPageIndexOffset)
	q.indexPage.PutUint32(uint32(messageOffset), indexOffset+messageOffsetOffset)
	q.indexPage.PutUint32(uint32(dataLen), indexOffset+messageLengthOffset)

	// save metadata
	q.metaPage.PutUint64(uint64(seq), queueAppendedSeqOffset)
	q.appendedSeq.Store(seq)

	// new data written, notify all waiting consumer groups can consume data
	q.notEmpty.Broadcast()
	return nil
}

// initSequence initializes sequences from the metadata.
func (q *queue) initSequence() {
	q.appendedSeq.Store(int64(q.metaPage.ReadUint64(queueAppendedSeqOffset)))
	q.acknowledgedSeq.Store(int64(q.metaPage.ReadUint64(queueAcknowledgedSeqOffset)))
}

// initDataPageIndex finds out data page head index and message offset
func (q *queue) initDataPageIndex() (err error) {
	if q.appendedSeq.Load() == SeqNoNewMessageAvailable {
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

	previousSeq := q.appendedSeq.Load() // get previous sequence
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

	if sequence > q.appendedSeq.Load() || sequence <= q.acknowledgedSeq.Load() {
		return fmt.Errorf("%w: get %d, range [%d~%d]", ErrOutOfSequenceRange,
			sequence, q.appendedSeq.Load(), q.acknowledgedSeq.Load())
	}

	return nil
}

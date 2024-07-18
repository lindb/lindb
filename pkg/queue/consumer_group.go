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
	"path/filepath"
	"sync"

	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source ./consumer_group.go -destination ./consumer_group_mock.go -package queue

// ConsumerGroup represents an individual consumer with own consume and ack sequence.
// The typical way to use ConsumerGroup is using a single go-routine to consume message,
// and using other go-routine to ack the messages which have been processed successfully.
type ConsumerGroup interface {
	// Name returns a unique name for ConsumerGroup in a FanOutQueue.
	Name() string
	// Consume returns the seq for the next data to consume.
	// If no new data is available, SeqNoNewMessageAvailable is returned.
	Consume() int64
	// SetConsumedSeq sets the consumed sequence, this is useful when re-consume message.
	// error returns when seq is invalid(less than acknowledged seq or greater than the read barrier).
	SetConsumedSeq(seq int64)
	// Ack mark the data processed with sequence less than or equals to acknowledged sequence.
	Ack(ackSeq int64)
	// ConsumedSeq returns the sequence of consumed.
	ConsumedSeq() int64
	// AcknowledgedSeq returns the acknowledged sequence.
	AcknowledgedSeq() int64
	// Queue returns underlying queue.
	Queue() FanOutQueue
	// Pause pauses consume data.
	Pause()
	// SetSeq sets consumed/acknowledged sequence.
	SetSeq(seq int64)
	// Pending returns the offset between ConsumerGroup consumed sequence and FanOutQueue appended sequence.
	Pending() int64
	// IsEmpty returns if fan out consumer cannot consume any data.
	IsEmpty() bool
	// Close persists  headSeq, tailSeq.
	Close()
	// consume returns the seq for the next data to consume.
	consume() int64
}

// consumerGroup implements ConsumerGroup.
type consumerGroup struct {
	q               FanOutQueue // underlying query for retreving data
	metaPageFct     page.Factory
	metaPage        page.MappedPage // persists meta
	consumedSeq     *atomic.Int64   // consumed sequence
	acknowledgedSeq *atomic.Int64   // acknowledged sequence
	name            string          // unique name
	closed          atomic.Bool     // false -> running, true -> closed
	paused          atomic.Bool
	lock4headSeq    sync.RWMutex // lock to protect head seq
}

// NewConsumerGroup builds a ConsumerGroup from metaPath.
func NewConsumerGroup(parent, fanOutPath string, q FanOutQueue) (ConsumerGroup, error) {
	name := filepath.Join(parent, fanOutPath)
	var err error
	metaPageFct, err := newPageFactoryFunc(name, consumerGroupMetaSize)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if err1 := metaPageFct.Close(); err1 != nil {
				queueLogger.Error("close meta page factory when create consumer group",
					logger.String("name", name), logger.Error(err))
			}
		}
	}()

	hasMeta := existFunc(filepath.Join(name, fmt.Sprintf("%d.bat", metaPageIndex)))

	metaPage, err := metaPageFct.AcquirePage(metaPageIndex)
	if err != nil {
		return nil, err
	}

	consumedSeq := int64(-1)
	ackSeq := int64(-1)

	if hasMeta {
		consumedSeq = int64(metaPage.ReadUint64(consumerGroupConsumedSeqOffset))
		ackSeq = int64(metaPage.ReadUint64(consumerGroupAcknowledgedSeqOffset))
		ackOfQueue := q.Queue().AcknowledgedSeq()
		// if queue ack > consume group ack, need reset use queue ack
		if ackSeq < ackOfQueue {
			ackSeq = ackOfQueue
		}
	}
	// persist metadata
	metaPage.PutUint64(uint64(consumedSeq), consumerGroupConsumedSeqOffset)
	metaPage.PutUint64(uint64(ackSeq), consumerGroupAcknowledgedSeqOffset)

	return &consumerGroup{
		name:            name,
		q:               q,
		metaPageFct:     metaPageFct,
		metaPage:        metaPage,
		consumedSeq:     atomic.NewInt64(consumedSeq),
		acknowledgedSeq: atomic.NewInt64(ackSeq),
	}, nil
}

// Name returns a unique name for ConsumerGroup in a FanOutQueue.
func (f *consumerGroup) Name() string {
	return f.name
}

// Pause pauses consume data.
func (f *consumerGroup) Pause() {
	f.paused.Store(true)
	f.Queue().Queue().Signal()
}

// isPause returns if consumer group is paused.
func (f *consumerGroup) isPause() bool {
	return f.closed.Load() || f.paused.Load()
}

// Consume returns the seq for the next data to consume.
// If no new data is available, SeqNoNewMessageAvailable is returned.
func (f *consumerGroup) Consume() int64 {
	headSeq := f.consumedSeq.Load() + 1

	// check queue if empty using current consume head without lock,
	// if queue is empty will waiting new data write or closed/paused.
	if !f.Queue().Queue().NotEmpty(headSeq, f.isPause) {
		return SeqNoNewMessageAvailable
	}
	return f.consume()
}

// consume returns the seq for the next data to consume.
func (f *consumerGroup) consume() int64 {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	headSeq := f.consumedSeq.Load() + 1
	if headSeq <= f.q.Queue().AppendedSeq() {
		f.consumedSeq.Store(headSeq)
		f.metaPage.PutUint64(uint64(headSeq), consumerGroupConsumedSeqOffset)
		return headSeq
	}

	return SeqNoNewMessageAvailable
}

// SetConsumedSeq sets the sequence of consumed.
func (f *consumerGroup) SetConsumedSeq(seq int64) {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	f.consumedSeq.Store(seq)
	f.metaPage.PutUint64(uint64(f.ConsumedSeq()), consumerGroupConsumedSeqOffset)
}

// Queue returns underlying queue.
func (f *consumerGroup) Queue() FanOutQueue {
	return f.q
}

// Ack mark the data with seq less than or equals to acknowledgedSeq.
func (f *consumerGroup) Ack(ackSeq int64) {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	ts := f.AcknowledgedSeq()
	hs := f.ConsumedSeq()
	// In the initial condition, ts == 0, if the first acknowledgedSeq == 0, it would be ignored.
	// Since ack is always in batch mode and the following ack will ack the previous data, it's not big problem.
	if ackSeq >= ts && ackSeq <= hs {
		f.acknowledgedSeq.Store(ackSeq)

		f.metaPage.PutUint64(uint64(f.ConsumedSeq()), consumerGroupConsumedSeqOffset)
		f.metaPage.PutUint64(uint64(f.AcknowledgedSeq()), consumerGroupAcknowledgedSeqOffset)

		if err := f.metaPage.Sync(); err != nil {
			queueLogger.Error("sync consumerGroup meta page error", logger.String("consumerGroup", f.name), logger.Error(err))
		}
	} else {
		queueLogger.Warn("ack failure, invalid ack seq", logger.Int64("ack", ackSeq),
			logger.Int64("lastAck", ts), logger.Int64("consumedSeq", hs))
	}
}

// ConsumedSeq returns the next sequence of consumed.
func (f *consumerGroup) ConsumedSeq() int64 {
	return f.consumedSeq.Load()
}

// AcknowledgedSeq returns the acknowledged sequence.
func (f *consumerGroup) AcknowledgedSeq() int64 {
	return f.acknowledgedSeq.Load()
}

// SetSeq sets consumed/acknowledged sequence.
func (f *consumerGroup) SetSeq(seq int64) {
	f.lock4headSeq.Lock()
	defer f.lock4headSeq.Unlock()

	f.consumedSeq.Store(seq)
	f.acknowledgedSeq.Store(seq)
	f.metaPage.PutUint64(uint64(f.ConsumedSeq()), consumerGroupConsumedSeqOffset)
	f.metaPage.PutUint64(uint64(f.AcknowledgedSeq()), consumerGroupAcknowledgedSeqOffset)
}

// Pending returns the offset between ConsumerGroup HeadSeq and FanOutQueue HeadSeq.
func (f *consumerGroup) Pending() int64 {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	fh := f.ConsumedSeq()
	qh := f.q.Queue().AppendedSeq()

	pending := qh - fh
	if pending < 0 {
		return 0
	}
	return pending
}

// IsEmpty returns if fan out consumer cannot consume any data.
func (f *consumerGroup) IsEmpty() bool {
	f.lock4headSeq.RLock()
	defer f.lock4headSeq.RUnlock()

	qh := f.q.Queue().AppendedSeq()

	return qh <= f.AcknowledgedSeq()
}

// Close persists headSeq, tailSeq.
func (f *consumerGroup) Close() {
	if f.closed.CompareAndSwap(false, true) {
		f.Queue().Queue().Signal()

		if err := f.metaPageFct.Close(); err != nil {
			queueLogger.Error("close consumerGroup meta error", logger.String("consumerGroup", f.name), logger.Error(err))
		}
	}
}

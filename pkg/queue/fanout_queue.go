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
	"path"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
)

//go:generate mockgen -source ./fanout_queue.go -destination ./fanout_queue_mock.go -package queue

// for testing
var (
	newQueueFunc         = NewQueue
	listDirFunc          = fileutil.ListDir
	newConsumerGroupFunc = NewConsumerGroup
)

// FanOutQueue represents a queue "produce once, consume multiple times".
// ConsumerGroup represents an individual consumer with own consume/acknowledge sequence.
type FanOutQueue interface {
	// Path returns path for persistence files.
	Path() string
	// Queue returns underlying queue.
	Queue() Queue
	// GetOrCreateConsumerGroup returns the ConsumerGroup if exists,
	// otherwise creates a new ConsumerGroup with consume seq and ack seq == queue ack seq.
	GetOrCreateConsumerGroup(name string) (ConsumerGroup, error)
	// ConsumerGroupNames returns all names of ConsumerGroup.
	ConsumerGroupNames() []string
	// Sync checks the acknowledged sequence of each ConsumerGroup, update the acknowledged sequence as the smallest one.
	// Then syncs metadata to storage.
	Sync()
	// SetAppendedSeq sets appended sequence underlying queue, then set consumed/acknowledged sequence for each ConsumerGroup.
	SetAppendedSeq(seq int64)
	// Close persists Seq meta, ConsumerGroup seq meta, release resources.
	Close()
}

// fanOutQueue implements FanOutQueue.
type fanOutQueue struct {
	dirPath          string                   // dir path for persistence file
	consumerGroupDir string                   // dir path for storing ConsumerGroup
	queue            Queue                    // underlying queue
	consumerGroups   map[string]ConsumerGroup // name -> ConsumerGroup

	lock4map sync.RWMutex // lock for fanOutMap
	closed   atomic.Bool  // false -> running, true -> closed
}

// NewFanOutQueue returns a FanOutQueue persisted in dirPath.
func NewFanOutQueue(dirPath string, dataSizeLimit int64) (q FanOutQueue, err error) {
	fq := &fanOutQueue{
		dirPath:          dirPath,
		consumerGroupDir: path.Join(dirPath, consumerGroupDirName),
		consumerGroups:   make(map[string]ConsumerGroup),
	}

	defer func() {
		if err != nil {
			// if initialize consumerGroup queue failure, need release the resource
			fq.Close()
		}
	}()

	// create underlying queue
	fq.queue, err = newQueueFunc(dirPath, dataSizeLimit)
	if err != nil {
		return nil, err
	}
	// init consumerGroup sequence
	err = fq.initConsumerGroups()
	if err != nil {
		return nil, err
	}

	return fq, nil
}

// Path returns path for persistence files.
func (fq *fanOutQueue) Path() string {
	return fq.dirPath
}

// Queue returns underlying queue.
func (fq *fanOutQueue) Queue() Queue {
	return fq.queue
}

// GetOrCreateConsumerGroup returns the ConsumerGroup if exists,
// otherwise creates a new ConsumerGroup with consume seq and ack seq == queue ack seq.
func (fq *fanOutQueue) GetOrCreateConsumerGroup(name string) (ConsumerGroup, error) {
	fq.lock4map.Lock()
	defer fq.lock4map.Unlock()

	if fo, ok := fq.consumerGroups[name]; ok {
		return fo, nil
	}

	fo, err := newConsumerGroupFunc(fq.consumerGroupDir, name, fq)
	if err != nil {
		return nil, err
	}

	fq.consumerGroups[name] = fo

	return fo, nil
}

// ConsumerGroupNames returns all names of ConsumerGroup.
func (fq *fanOutQueue) ConsumerGroupNames() (names []string) {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()

	for name := range fq.consumerGroups {
		names = append(names, name)
	}
	return names
}

// SetAppendedSeq sets appended sequence underlying queue, then set consumed/acknowledged sequence for each ConsumerGroup.
func (fq *fanOutQueue) SetAppendedSeq(seq int64) {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()

	fq.queue.SetAppendedSeq(seq)

	for _, fo := range fq.consumerGroups {
		fo.SetSeq(seq)
	}
}

// Sync checks the acknowledged sequence of each ConsumerGroup, update the acknowledged sequence as the smallest one.
// Then syncs metadata to storage.
func (fq *fanOutQueue) Sync() {
	fq.lock4map.RLock()
	defer fq.lock4map.RUnlock()

	// no consumer group
	if len(fq.consumerGroups) == 0 {
		return
	}

	// use the queue appended sequence as the init value
	ackSeq := fq.queue.AppendedSeq()

	for _, fo := range fq.consumerGroups {
		ts := fo.AcknowledgedSeq()
		if ts < ackSeq {
			ackSeq = ts
		}
	}

	if ackSeq >= 0 {
		fq.queue.SetAcknowledgedSeq(ackSeq)
	}
}

// Close persists Seq meta, ConsumerGroup seq meta, release resources.
func (fq *fanOutQueue) Close() {
	if fq.closed.CAS(false, true) {
		fq.lock4map.RLock()
		defer fq.lock4map.RUnlock()

		for _, fo := range fq.consumerGroups {
			fo.Close()
		}

		if fq.queue != nil {
			fq.queue.Close()
		}
	}
}

// initConsumerGroups initializes exist ConsumerGroup.
func (fq *fanOutQueue) initConsumerGroups() error {
	if err := mkDirFunc(fq.consumerGroupDir); err != nil {
		return err
	}

	fileNames, err := listDirFunc(fq.consumerGroupDir)
	if err != nil {
		return err
	}

	// load exist ConsumerGroup
	for _, fn := range fileNames {
		fo, err := newConsumerGroupFunc(fq.consumerGroupDir, fn, fq)
		if err != nil {
			return err
		}

		fq.consumerGroups[fn] = fo
	}

	return nil
}

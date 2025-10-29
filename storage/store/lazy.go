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

package store

import (
	"sync"

	"go.uber.org/atomic"
)

type CreatePartitionFn func(timestamp int64) (Partition, error)

type LazyPartition struct {
	parition  Partition
	loaded    *atomic.Bool
	timestamp int64
	err       error

	lock sync.Mutex

	createPartitionFn func(timestamp int64) (Partition, error)
}

func NewPartition(partition Partition) *LazyPartition {
	return &LazyPartition{
		loaded:   atomic.NewBool(true),
		parition: partition,
	}
}

func NewLazyPartition(tiemestamp int64, createPartitionFn func(timestamp int64) (Partition, error)) *LazyPartition {
	return &LazyPartition{
		loaded:            atomic.NewBool(false),
		timestamp:         tiemestamp,
		createPartitionFn: createPartitionFn,
	}
}

func (l *LazyPartition) Loaded() bool {
	return l.loaded.Load()
}

func (l *LazyPartition) Get() (Partition, error) {
	if !l.loaded.Load() {
		l.lock.Lock()
		defer func() {
			l.lock.Unlock()
			l.loaded.Store(true)
		}()

		if l.parition == nil {
			l.parition, l.err = l.createPartitionFn(l.timestamp)
		}
	}
	return l.parition, l.err
}

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

package memdb

import (
	"fmt"
	"path/filepath"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source ./buffer_manager.go -destination=./buffer_manager_mock.go -package memdb

// BufferManager represents data points write buffer manager, maintains all buffers for spec shard.
type BufferManager interface {
	// AllocBuffer allocates a new DataPointBuffer.
	AllocBuffer(familyTime int64) (buf DataPointBuffer, err error)
	// GarbageCollect cleans all dirty buffers.
	GarbageCollect()
	// Cleanup cleans all history buffers.
	Cleanup()
}

// bufferManager implements BufferManager.
type bufferManager struct {
	path string

	value atomic.Value // []DataPointBuffer

	logger *logger.Logger
}

// NewBufferManager creates a BufferManager instance.
func NewBufferManager(path string) BufferManager {
	mgr := &bufferManager{
		path:   path,
		logger: logger.GetLogger("TSDB", "BufferManager"),
	}
	mgr.value.Store(make([]DataPointBuffer, 0))
	return mgr
}

// AllocBuffer allocates a new DataPointBuffer.
func (b *bufferManager) AllocBuffer(familyTime int64) (buf DataPointBuffer, err error) {
	// path = root path + family time + create time(nano)
	path := filepath.Join(b.path,
		timeutil.FormatTimestamp(familyTime, timeutil.DataTimeFormat4),
		fmt.Sprintf("%d", timeutil.NowNano()))
	buf, err = newDataPointBuffer(path)
	if err != nil {
		return nil, err
	}

	// copy for write
	set := b.value.Load().([]DataPointBuffer)
	newSet := make([]DataPointBuffer, 0, len(set)+1)
	newSet = append(newSet, set...)
	newSet = append(newSet, buf)
	b.value.Store(newSet)

	return buf, err
}

// GarbageCollect cleans all dirty buffers.
func (b *bufferManager) GarbageCollect() {
	oldSet := b.value.Load().([]DataPointBuffer)
	newSet := make([]DataPointBuffer, 0)

	// gc all dirty buffer, then remove it from alive list
	for idx := range oldSet {
		buf := oldSet[idx]
		needClean := false
		if buf.IsDirty() {
			// close buf(remove tmp files) if buf is dirty
			if err := buf.Close(); err != nil {
				b.logger.Error("close data write buffer", logger.Error(err))
			} else {
				needClean = true
			}
		}

		if !needClean {
			newSet = append(newSet, buf)
		}
	}

	b.value.Store(newSet)
}

// Cleanup cleans all history buffers.
func (b *bufferManager) Cleanup() {
	err := removeFunc(b.path)
	if err != nil {
		b.logger.Error("clean up data point write buffer", logger.String("path", b.path), logger.Error(err))
	}
}

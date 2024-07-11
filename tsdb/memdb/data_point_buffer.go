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
	"io"
	"os"
	"path/filepath"
	"sync"

	commonfileutil "github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/imap"
)

//go:generate mockgen -source ./data_point_buffer.go -destination=./data_point_buffer_mock.go -package memdb

// for testing
var (
	closeFileFunc = closeFile
	mkdirFunc     = commonfileutil.MkDirIfNotExist
	mapFunc       = fileutil.RWMap
	unmapFunc     = fileutil.Unmap
	removeFunc    = commonfileutil.RemoveDir
	openFileFunc  = os.OpenFile
)

const (
	// TODO: add db config
	regionSize = 128 * 1024 * 1024 // 128M
	pageSize   = 128
	pageCount  = regionSize / pageSize
)

// DataPointBuffer represents data point buffer write buffer based on memory map file
type DataPointBuffer interface {
	io.Closer
	// GetOrCreatePage returns write page buffer, if not exist create new page buffer.
	GetOrCreatePage(memSeriesID uint32) ([]byte, error)
	// GetPage returns write page buffer, if not exist returns nil.
	GetPage(memSeriesID uint32) ([]byte, bool)
	// Release marks data point buffer is dirty.
	Release()
	// IsDirty returns data point buffer if dirty, dirty buffer can be collect.
	IsDirty() bool
}

// dataPointBuffer implements DataPointBuffer interface
type dataPointBuffer struct {
	ids       *imap.IntMap[int32] // store all time series ids(memory time series id => page id)
	path      string
	buf       [][]byte
	files     []*os.File
	dirty     atomic.Bool
	lock      sync.RWMutex
	pageIDSeq int32
}

// newDataPointBuffer creates data point buffer for writing points of metric.
func newDataPointBuffer(path string) (DataPointBuffer, error) {
	if err := mkdirFunc(path); err != nil {
		return nil, err
	}
	return &dataPointBuffer{
		path:      path,
		pageIDSeq: 0,
		ids:       imap.NewIntMap[int32](),
	}, nil
}

// GetOrCreatePage returns write page buffer, if not exist create new page buffer.
func (d *dataPointBuffer) GetOrCreatePage(memSeriesID uint32) ([]byte, error) {
	var (
		pageID int32
		ok     bool
	)

	d.lock.RLock()
	pageID, ok = d.ids.Get(memSeriesID)
	d.lock.RUnlock()
	if !ok {
		// generate a new page id
		// NOTE: single goroutine write family data, so can read directly
		pageID = d.pageIDSeq
	}
	region := pageID / pageCount
	rOffset := pageID % pageCount
	if !ok && rOffset == 0 {
		// if page id is new and region not found, then create a new temp region buffer
		if err := mkdirFunc(d.path); err != nil {
			return nil, err
		}
		f, err := openFileFunc(filepath.Join(d.path, fmt.Sprintf("%d.tmp", region)), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		buf, err := mapFunc(f, regionSize)
		if err != nil {
			_ = f.Close()
			return nil, err
		}
		d.files = append(d.files, f)
		d.buf = append(d.buf, buf)
	}
	offset := pageSize * rOffset

	if !ok {
		d.lock.Lock()
		d.ids.PutIfNotExist(memSeriesID, pageID)
		d.lock.Unlock()

		// increase page id sequence if all operators successfully
		d.pageIDSeq++
	}
	return d.buf[region][offset : offset+pageSize], nil
}

// GetPage returns write page buffer, if not exist returns nil.
func (d *dataPointBuffer) GetPage(memSeriesID uint32) ([]byte, bool) {
	var (
		pageID int32
		ok     bool
	)
	d.lock.RLock()
	pageID, ok = d.ids.Get(memSeriesID) // find page id by memory time series
	d.lock.RUnlock()

	if !ok {
		return nil, false
	}
	region := pageID / pageCount
	rOffset := pageID % pageCount
	offset := pageSize * rOffset
	return d.buf[region][offset : offset+pageSize], true
}

// Release marks data point buffer is dirty.
func (d *dataPointBuffer) Release() {
	d.dirty.Store(true)
}

// IsDirty returns data point buffer if dirty, dirty buffer can be collect.
func (d *dataPointBuffer) IsDirty() bool {
	return d.dirty.Load()
}

// Close closes data point buffer, unmap memory map file
func (d *dataPointBuffer) Close() error {
	if !d.dirty.Load() {
		memDBLogger.Error("buffer is not dirty, cannot close it",
			logger.String("file", d.path))
		return nil
	}
	d.closeBuffer()
	if err := removeFunc(d.path); err != nil {
		memDBLogger.Error("remove buffer file in memory database err",
			logger.String("file", d.path), logger.Error(err))
	}
	return nil
}

// closeBuffer just closes file for unix.
func (d *dataPointBuffer) closeBuffer() {
	for i, buf := range d.buf {
		if err := unmapFunc(d.files[i], buf); err != nil {
			memDBLogger.Error("unmap file in memory database err",
				logger.String("file", d.path), logger.Error(err))
		}
	}
	for _, f := range d.files {
		if err := closeFileFunc(f); err != nil {
			memDBLogger.Error("close file in memory database err",
				logger.String("file", d.path), logger.Error(err))
		}
	}
}

// closeFile closes file.
func closeFile(f *os.File) error {
	return f.Close()
}

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

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./data_point_buffer.go -destination=./data_point_buffer_mock.go -package memdb

var (
	mkdirFunc    = fileutil.MkDirIfNotExist
	mapFunc      = fileutil.RWMap
	removeFunc   = fileutil.RemoveDir
	openFileFunc = os.OpenFile
)

const (
	regionSize = 128 * 1024 * 1024 // 128M
	pageSize   = 128
	pageCount  = regionSize / pageSize
)

// DataPointBuffer represents data point buffer write buffer based on memory map file
type DataPointBuffer interface {
	io.Closer
	// AllocPage allocates the page buffer for writing data point.
	AllocPage() (buf []byte, err error)
	// Release marks data point buffer is dirty.
	Release()
	// IsDirty returns data point buffer if dirty, dirty buffer can be collect.
	IsDirty() bool
}

// dataPointBuffer implements DataPointBuffer interface
type dataPointBuffer struct {
	path      string
	buf       [][]byte
	files     []*os.File
	pageIDSeq atomic.Int32
	dirty     atomic.Bool
}

// newDataPointBuffer creates data point buffer for writing points of metric.
func newDataPointBuffer(path string) (DataPointBuffer, error) {
	if err := mkdirFunc(path); err != nil {
		return nil, err
	}
	return &dataPointBuffer{
		path:      path,
		pageIDSeq: *atomic.NewInt32(-1),
	}, nil
}

// AllocPage allocates the page buffer for writing data point
func (d *dataPointBuffer) AllocPage() (buf []byte, err error) {
	pageID := d.pageIDSeq.Inc()
	if pageID%pageCount == 0 {
		if err := mkdirFunc(d.path); err != nil {
			return nil, err
		}
		f, err := openFileFunc(filepath.Join(d.path, fmt.Sprintf("%d.tmp", pageID/pageCount)), os.O_CREATE|os.O_RDWR, 0644)
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
	region := uint16(pageID / pageCount)
	if d.buf == nil || uint16(len(d.buf)) <= region {
		return nil, fmt.Errorf("wrong region in memory buffer")
	}
	offset := pageSize * (int(pageID) % pageCount)
	return d.buf[region][offset : offset+pageSize], nil
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

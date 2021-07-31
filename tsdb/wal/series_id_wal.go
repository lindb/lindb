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

package wal

import (
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source=./series_id_wal.go -destination=./series_id_wal_mock.go -package=wal

// for testing
var (
	mkDirFunc          = fileutil.MkDirIfNotExist
	newPageFactoryFunc = page.NewFactory
)

var (
	recoverSeriesFailCounter = walScope.NewDeltaCounter("wal_recovery_series_fail")
)

const (
	seriesEntryLength = 4 + 8 + 4                      // metric id + tags hash + series id
	seriesPageSize    = seriesEntryLength * 512 * 1024 // series wal page size
	metricIDOffset    = 0                              // metric id offset
	tagsHashOffset    = metricIDOffset + 4             // tags hash offset
	seriesIDOffset    = tagsHashOffset + 8             // series id offset
)

// SeriesWAL represents write ahead log which stores series data for index database
type SeriesWAL interface {
	// Append appends metricID/tagsHash/seriesID into wal log
	Append(metricID uint32, tagsHash uint64, seriesID uint32) error
	// NeedRecovery checks if wal log need to recover
	NeedRecovery() bool
	// Recovery recoveries wal log, then writes data via recovery function
	Recovery(recovery SeriesRecoveryFunc, commit CommitFunc)
	// Sync flushes data into disk
	Sync() error
	// Close closes the wal log
	Close() error
}

// seriesWAL implements SeriesWAL interface
type seriesWAL struct {
	base *baseWAL
}

// NewSeriesWAL creates a new series write ahead log
func NewSeriesWAL(path string) (SeriesWAL, error) {
	base, err := newBaseWAL(path, metricMetaPageSize)
	if err != nil {
		return nil, err
	}
	return &seriesWAL{base: base}, nil
}

// Append appends metricID/tagsHash/seriesID into wal log
func (wal *seriesWAL) Append(metricID uint32, tagsHash uint64, seriesID uint32) (err error) {
	if err := wal.base.checkPage(seriesEntryLength); err != nil {
		return err
	}
	wal.base.putUint32(metricID)
	wal.base.putUint64(tagsHash)
	wal.base.putUint32(seriesID)

	return nil
}

// NeedRecovery checks if wal log need to recover
func (wal *seriesWAL) NeedRecovery() bool {
	return wal.base.needRecovery()
}

// Recovery recoveries wal log, then writes data via recovery function
func (wal *seriesWAL) Recovery(recovery SeriesRecoveryFunc, commit CommitFunc) {
	current := wal.base.pageIndex.Load()
	committed := wal.base.commitPageIndex.Load()
	for i := committed; i < current; i++ {
		walPage, ok := wal.base.walFactory.GetPage(i)
		if !ok {
			continue
		}
		offset := 0
		for offset < seriesPageSize {
			metricID := walPage.ReadUint32(offset + metricIDOffset)
			if metricID == 0 {
				break
			}

			if err := recovery(metricID,
				walPage.ReadUint64(offset+tagsHashOffset),
				walPage.ReadUint32(offset+seriesIDOffset)); err != nil {
				recoverSeriesFailCounter.Incr()

				walLogger.Error("invoke recovery func error",
					logger.String("wal", wal.base.path), logger.Error(err))
				return
			}
			offset += seriesEntryLength
		}

		if err := commit(); err != nil {
			recoveryCommitFailCounter.Incr()

			walLogger.Error("invoke commit func error",
				logger.String("wal", wal.base.path), logger.Error(err))
			return
		}

		if err := wal.base.walFactory.ReleasePage(i); err != nil {
			releaseWALPageFailCounter.Incr()

			walLogger.Error("release series wal page error",
				logger.String("wal", wal.base.path), logger.Error(err))
		}

		wal.base.commitPageIndex.Inc()
	}
}

// Sync flushes data into disk
func (wal *seriesWAL) Sync() error {
	return wal.base.sync()
}

// Close closes the wal log
func (wal *seriesWAL) Close() error {
	return wal.base.close()
}

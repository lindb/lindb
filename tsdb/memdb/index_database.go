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
	"context"
	"sync"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/series/metric"
)

// IndexDatabase represents memory index database for storing metric index(shard level).
//
//go:generate mockgen -source ./index_database.go -destination=./index_database_mock.go -package=memdb
type IndexDatabase interface {
	// GetOrCreateTimeSeriesIndex returns time series index if not exist create new.
	GetOrCreateTimeSeriesIndex(row *metric.StorageRow) TimeSeriesIndex
	// GenMemSeriesID generates memory time series id.
	GenMemSeriesID() uint32
	// GetMetadataDatabase returns memory metadata database.
	GetMetadataDatabase() MetadataDatabase
	// GetTimeSeriesIndex returns memory time series index by memory metric id.
	GetTimeSeriesIndex(memMetricID uint64) (TimeSeriesIndex, bool)
	// Cleanup cleanups index data for inactive memory database.
	Cleanup(db MemoryDatabase)
	// Notify notifies update or flush metric index.
	Notify(event any)
	// Close closed index database.
	Close()
}

// indexDatabase implements IndexDatabase interface.
type indexDatabase struct {
	metaDB  MetadataDatabase
	indexDB index.MetricIndexDatabase

	ctx    context.Context
	cancel context.CancelFunc

	ch                chan any
	timeSeriesIndexes sync.Map // hash(ns + metirc name) => metric index store(map[uint64]TimeSeriesIndex)

	timeSeriesSeq atomic.Uint32 // like db primary key sequence(memory level)

	lock sync.RWMutex
}

// NewIndexDatabase creates IndexDatabase instance.
func NewIndexDatabase(metaDB MetadataDatabase, indexDB index.MetricIndexDatabase) IndexDatabase {
	ctx, cacnel := context.WithCancel(context.TODO())
	db := &indexDatabase{
		metaDB:  metaDB,
		indexDB: indexDB,
		ch:      make(chan any, 100), // TODO: add config
		ctx:     ctx,
		cancel:  cacnel,
	}
	go db.handle()
	return db
}

// Notify notifies update or flush metric index.
func (idb *indexDatabase) Notify(event any) {
	idb.ch <- event
}

// GetMetadataDatabase returns memory metadata database.
func (idb *indexDatabase) GetMetadataDatabase() MetadataDatabase {
	return idb.metaDB
}

// GetOrCreateTimeSeriesIndex returns time series index if not exist create new.
func (idb *indexDatabase) GetOrCreateTimeSeriesIndex(row *metric.StorageRow) TimeSeriesIndex {
	nameHash := row.NameHash()
	timeSeriesIndex, ok := idb.timeSeriesIndexes.Load(nameHash)
	if ok {
		return timeSeriesIndex.(TimeSeriesIndex)
	}

	idb.lock.Lock()
	defer idb.lock.Unlock()

	return idb.getOrCreateTimeSeriesIndex(nameHash)
}

func (idb *indexDatabase) getOrCreateTimeSeriesIndex(nameHash uint64) TimeSeriesIndex {
	timeSeriesIndex, ok := idb.timeSeriesIndexes.Load(nameHash)
	if ok {
		return timeSeriesIndex.(TimeSeriesIndex)
	}
	newTimeSeriesIndex := NewTimeSeriesIndex()
	// store time series index
	idb.timeSeriesIndexes.Store(nameHash, newTimeSeriesIndex)
	return newTimeSeriesIndex
}

// GenMemSeriesID generates memory time series id.
func (idb *indexDatabase) GenMemSeriesID() uint32 {
	return idb.timeSeriesSeq.Inc()
}

// GetTimeSeriesIndex returns memory time series index by memory metric id.
func (idb *indexDatabase) GetTimeSeriesIndex(memMetricID uint64) (TimeSeriesIndex, bool) {
	timeSeriesIndex, ok := idb.timeSeriesIndexes.Load(memMetricID)
	if ok {
		return timeSeriesIndex.(TimeSeriesIndex), ok
	}
	return nil, false
}

// Cleanup cleanups index data for inactive memory database.
func (idb *indexDatabase) Cleanup(db MemoryDatabase) {
	familyCreateTime := db.CreatedTime()
	expiredTimestamp := timeutil.Now()
	memTimeSeriesIDs := db.MemTimeSeriesIDs()
	gcTimestamp := timeutil.Now() - 3*timeutil.OneHour // TODO: add config?
	idb.timeSeriesIndexes.Range(func(key, value any) bool {
		timeSeriesIndex := (value.(TimeSeriesIndex))
		timeSeriesIndex.ClearTimeRange(familyCreateTime)
		timeSeriesIndex.ExpireTimeSeriesIDs(memTimeSeriesIDs, expiredTimestamp)
		timeSeriesIndex.GC(gcTimestamp)

		// if no time series undex index, remove it from metric index store
		if timeSeriesIndex.NumOfSeries() == 0 {
			idb.timeSeriesIndexes.Delete(key)
		}
		return true
	})
}

// Close closed index database.
func (idb *indexDatabase) Close() {
	close(idb.ch)
}

func (idb *indexDatabase) indexTimeSeries(row *metric.StorageRow, seriesID uint32) {
	nameHash := row.NameHash()
	timeSeriesIndexObj, _ := idb.timeSeriesIndexes.Load(nameHash)
	timeSeriesIndex := timeSeriesIndexObj.(TimeSeriesIndex)
	idb.lock.Lock()
	timeSeriesIndex.IndexTimeSeries(seriesID, row.MemSeriesID)
	idb.lock.Unlock()
}

func (idb *indexDatabase) handle() {
	for e := range idb.ch {
		switch event := e.(type) {
		case *metric.StorageRow:
			idb.handleRow(event)
		case *FlushEvent:
			idb.indexDB.PrepareFlush()
			// flush data background
			go idb.handleFlush(event)
		}
	}
}

func (idb *indexDatabase) handleFlush(event *FlushEvent) {
	err := idb.indexDB.Flush()
	event.Callback(err)
}

func (idb *indexDatabase) handleRow(row *metric.StorageRow) {
	defer row.Done()

	metricID, err := idb.metaDB.GetMetaDB().GenMetricID(row.NameSpace(), row.Name())
	if err != nil {
		memDBLogger.Warn("generate metric id error", logger.String("namespace", string(row.NameSpace())),
			logger.String("metric", string(row.Name())), logger.Error(err))
		return
	}

	seriesID := idb.indexDB.GenSeriesID(metricID, row)
	idb.indexTimeSeries(row, seriesID)
}

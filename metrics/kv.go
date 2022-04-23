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

package metrics

import "github.com/lindb/lindb/internal/linmetric"

var (
	// table cache
	tableCacheScope = linmetric.StorageRegistry.NewScope("lindb.kv.table.cache")
	// TableCacheStatistics represents table reader cache statistics.
	TableCacheStatistics = struct {
		Evict         *linmetric.BoundCounter
		Hit           *linmetric.BoundCounter
		Miss          *linmetric.BoundCounter
		Close         *linmetric.BoundCounter
		CloseErr      *linmetric.BoundCounter
		ActiveReaders *linmetric.BoundGauge
	}{
		Evict:         tableCacheScope.NewCounter("evict_counts"),
		Hit:           tableCacheScope.NewCounter("cache_hits"),
		Miss:          tableCacheScope.NewCounter("cache_misses"),
		Close:         tableCacheScope.NewCounter("close_counts"),
		CloseErr:      tableCacheScope.NewCounter("close_errors"),
		ActiveReaders: tableCacheScope.NewGauge("active_readers"),
	}

	// table write
	tableWriteScope = linmetric.StorageRegistry.NewScope("lindb.kv.table.write")
	// TableWriteStatistics represents table file write statistics.
	TableWriteStatistics = struct {
		AddBadKeys *linmetric.BoundCounter
		AddKeys    *linmetric.BoundCounter
		WriteBytes *linmetric.BoundCounter
	}{
		AddBadKeys: tableWriteScope.NewCounter("bad_keys"),
		AddKeys:    tableWriteScope.NewCounter("add_keys"),
		WriteBytes: tableWriteScope.NewCounter("write_bytes"),
	}

	// table read
	tableReadScope = linmetric.StorageRegistry.NewScope("lindb.kv.table.read")
	// TableReadStatistics represents table file read statistics.
	TableReadStatistics = struct {
		GetErrors    *linmetric.BoundCounter
		GetCounts    *linmetric.BoundCounter
		ReadBytes    *linmetric.BoundCounter
		MMapCounts   *linmetric.BoundCounter
		MMapErrors   *linmetric.BoundCounter
		UnMMapCounts *linmetric.BoundCounter
		UnMMapErrors *linmetric.BoundCounter
	}{
		GetErrors:    tableReadScope.NewCounter("get_errors"),
		GetCounts:    tableReadScope.NewCounter("get_counts"),
		ReadBytes:    tableReadScope.NewCounter("read_bytes"),
		MMapCounts:   tableReadScope.NewCounter("mmap_counts"),
		MMapErrors:   tableReadScope.NewCounter("mmap_errors"),
		UnMMapCounts: tableReadScope.NewCounter("unmmap_counts"),
		UnMMapErrors: tableReadScope.NewCounter("unmmap_errors"),
	}

	// compact job
	compactScope = linmetric.StorageRegistry.NewScope("lindb.kv.compaction")
	// CompactStatistics represents compact job statistics.
	CompactStatistics = struct {
		Compacting *linmetric.GaugeVec
		Failure    *linmetric.DeltaCounterVec
		Duration   *linmetric.DeltaHistogramVec
	}{
		Compacting: flushScope.NewGaugeVec("compacting", "type"),
		Failure:    flushScope.NewCounterVec("failure", "type"),
		Duration:   compactScope.Scope("duration").NewHistogramVec("type"),
	}

	// flush job
	flushScope = linmetric.StorageRegistry.NewScope("lindb.kv.flush")
	// FlushStatistics represents flush job statistics.
	FlushStatistics = struct {
		Flushing *linmetric.BoundGauge
		Failure  *linmetric.BoundCounter
		Duration *linmetric.BoundHistogram
	}{
		Flushing: flushScope.NewGauge("flushing"),
		Failure:  flushScope.NewCounter("failure"),
		Duration: flushScope.Scope("duration").NewHistogram(),
	}
)

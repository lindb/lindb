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
		Evict                *linmetric.BoundCounter // evict reader from cache
		Hit                  *linmetric.BoundCounter // get reader hit cache
		Miss                 *linmetric.BoundCounter // get reader miss cache
		Close                *linmetric.BoundCounter // close reader success
		CloseFailures        *linmetric.BoundCounter // close reader failure
		CreateReaderFailures *linmetric.BoundCounter // create read failure
		ActiveReaders        *linmetric.BoundGauge   // number of active reader in cache
	}{
		Evict:         tableCacheScope.NewCounter("evicts"),
		Hit:           tableCacheScope.NewCounter("cache_hits"),
		Miss:          tableCacheScope.NewCounter("cache_misses"),
		Close:         tableCacheScope.NewCounter("closes"),
		CloseFailures: tableCacheScope.NewCounter("close_failuress"),
		ActiveReaders: tableCacheScope.NewGauge("active_readers"),
	}

	// table write
	tableWriteScope = linmetric.StorageRegistry.NewScope("lindb.kv.table.write")
	// TableWriteStatistics represents table file write statistics.
	TableWriteStatistics = struct {
		AddBadKeys *linmetric.BoundCounter // add bad key count
		AddKeys    *linmetric.BoundCounter // add key success
		WriteBytes *linmetric.BoundCounter // write data bytes
	}{
		AddBadKeys: tableWriteScope.NewCounter("bad_keys"),
		AddKeys:    tableWriteScope.NewCounter("add_keys"),
		WriteBytes: tableWriteScope.NewCounter("write_bytes"),
	}

	// table read
	tableReadScope = linmetric.StorageRegistry.NewScope("lindb.kv.table.read")
	// TableReadStatistics represents table file read statistics.
	TableReadStatistics = struct {
		Gets           *linmetric.BoundCounter // get data by key success
		GetFailures    *linmetric.BoundCounter // get data by key failure
		ReadBytes      *linmetric.BoundCounter // bytes of read data
		MMaps          *linmetric.BoundCounter // map file success
		MMapFailures   *linmetric.BoundCounter // map file failure
		UnMMaps        *linmetric.BoundCounter // unmap file success
		UnMMapFailures *linmetric.BoundCounter // unmap file failures
	}{
		Gets:           tableReadScope.NewCounter("gets"),
		GetFailures:    tableReadScope.NewCounter("get_failures"),
		ReadBytes:      tableReadScope.NewCounter("read_bytes"),
		MMaps:          tableReadScope.NewCounter("mmaps"),
		MMapFailures:   tableReadScope.NewCounter("mmap_failures"),
		UnMMaps:        tableReadScope.NewCounter("unmmaps"),
		UnMMapFailures: tableReadScope.NewCounter("unmmap_failures"),
	}

	// compact job
	compactScope = linmetric.StorageRegistry.NewScope("lindb.kv.compaction")
	// CompactStatistics represents compact job statistics.
	CompactStatistics = struct {
		Compacting *linmetric.GaugeVec          // number of compacting jobs
		Failure    *linmetric.DeltaCounterVec   // compact failure
		Duration   *linmetric.DeltaHistogramVec // compact duration(include count)
	}{
		Compacting: flushScope.NewGaugeVec("compacting", "type"),
		Failure:    flushScope.NewCounterVec("failure", "type"),
		Duration:   compactScope.Scope("duration").NewHistogramVec("type"),
	}

	// flush job
	flushScope = linmetric.StorageRegistry.NewScope("lindb.kv.flush")
	// FlushStatistics represents flush job statistics.
	FlushStatistics = struct {
		Flushing *linmetric.BoundGauge     // number of flushing jobs
		Failure  *linmetric.BoundCounter   // flush job failure
		Duration *linmetric.BoundHistogram // flush duration(include count)
	}{
		Flushing: flushScope.NewGauge("flushing"),
		Failure:  flushScope.NewCounter("failure"),
		Duration: flushScope.Scope("duration").NewHistogram(),
	}
)

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
	tableCacheScope      = linmetric.StorageRegistry.NewScope("lindb.kv.table.cache")
	TableCacheStatistics = struct {
		EvictCounter    *linmetric.BoundCounter
		HitCounter      *linmetric.BoundCounter
		MissCounter     *linmetric.BoundCounter
		CloseCounter    *linmetric.BoundCounter
		CloseErrCounter *linmetric.BoundCounter
		ActiveReaders   *linmetric.BoundGauge
	}{
		EvictCounter:    tableCacheScope.NewCounter("evict_counts"),
		HitCounter:      tableCacheScope.NewCounter("cache_hits"),
		MissCounter:     tableCacheScope.NewCounter("cache_misses"),
		CloseCounter:    tableCacheScope.NewCounter("close_counts"),
		CloseErrCounter: tableCacheScope.NewCounter("close_errors"),
		ActiveReaders:   tableCacheScope.NewGauge("active_readers"),
	}

	// table write
	tableWriteScope      = linmetric.StorageRegistry.NewScope("lindb.kv.table.write")
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
	tableReadScope      = linmetric.StorageRegistry.NewScope("lindb.kv.table.read")
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
)

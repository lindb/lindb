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
	//  database metric
	dbScope            = linmetric.StorageRegistry.NewScope("lindb.tsdb.database")
	DatabaseStatistics = struct {
		MetaDBFlushFailures *linmetric.DeltaCounterVec
		MetaDBFlushDuration *linmetric.DeltaHistogramVec
	}{
		MetaDBFlushFailures: dbScope.NewCounterVec("metadb_flush_failures", "db"),
		MetaDBFlushDuration: dbScope.Scope("metadb_flush_duration").NewHistogramVec("db"),
	}

	// memory database metric
	memDBScope      = linmetric.StorageRegistry.NewScope("lindb.tsdb.memdb")
	MemDBStatistics = struct {
		AllocatedPages       *linmetric.DeltaCounterVec
		AllocatePageFailures *linmetric.DeltaCounterVec
	}{
		AllocatedPages:       memDBScope.NewCounterVec("allocated_pages", "db"),
		AllocatePageFailures: memDBScope.NewCounterVec("allocate_page_failures", "db"),
	}

	// index database metric
	indexDBScope      = linmetric.StorageRegistry.NewScope("lindb.tsdb.indexdb")
	IndexDBStatistics = struct {
		BuildInvertedIndex *linmetric.DeltaCounterVec
	}{
		BuildInvertedIndex: indexDBScope.NewCounterVec("build_inverted_index", "db"),
	}

	// mete database metric
	metaDBScope      = linmetric.StorageRegistry.NewScope("lindb.tsdb.metadb")
	MetaDBStatistics = struct {
		GenMetricIDs          *linmetric.DeltaCounterVec
		GenMetricIDFailures   *linmetric.DeltaCounterVec
		GenFieldIDs           *linmetric.DeltaCounterVec
		GenFieldIDFailures    *linmetric.DeltaCounterVec
		GenTagKeyIDs          *linmetric.DeltaCounterVec
		GenTagKeyIDFailures   *linmetric.DeltaCounterVec
		GenTagValueIDs        *linmetric.DeltaCounterVec
		GenTagValueIDFailures *linmetric.DeltaCounterVec
	}{
		GenMetricIDs:          metaDBScope.NewCounterVec("gen_metric_ids", "db"),
		GenMetricIDFailures:   metaDBScope.NewCounterVec("gen_metric_id_failures", "db"),
		GenTagKeyIDs:          metaDBScope.NewCounterVec("gen_tag_key_ids", "db"),
		GenTagKeyIDFailures:   metaDBScope.NewCounterVec("gen_tag_key_id_failures", "db"),
		GenTagValueIDs:        metaDBScope.NewCounterVec("gen_tag_value_ids", "db"),
		GenTagValueIDFailures: metaDBScope.NewCounterVec("gen_tag_value_id_failures", "db"),
		GenFieldIDs:           metaDBScope.NewCounterVec("gen_field_ids", "db"),
		GenFieldIDFailures:    metaDBScope.NewCounterVec("gen_field_id_failures", "db"),
	}

	// shard metric
	shardScope      = linmetric.StorageRegistry.NewScope("lindb.tsdb.shard")
	ShardStatistics = struct {
		LookupMetricMetaFailures *linmetric.DeltaCounterVec
		ActiveFamilies           *linmetric.GaugeVec
		WriteBatches             *linmetric.DeltaCounterVec
		WriteMetrics             *linmetric.DeltaCounterVec
		WriteFields              *linmetric.DeltaCounterVec
		WriteMetricFailures      *linmetric.DeltaCounterVec
		FlushInFlight            *linmetric.GaugeVec
		MemDBTotalSize           *linmetric.GaugeVec
		ActiveMemDBs             *linmetric.GaugeVec
		MemDBFlushFailures       *linmetric.DeltaCounterVec
		MemDBFlushDuration       *linmetric.DeltaHistogramVec
		IndexDBFlushDuration     *linmetric.DeltaHistogramVec
		IndexDBFlushFailures     *linmetric.DeltaCounterVec
	}{
		LookupMetricMetaFailures: shardScope.NewCounterVec("lookup_metric_meta_failures", "db", "shard"),
		ActiveFamilies:           shardScope.NewGaugeVec("active_families", "db", "shard"),
		WriteBatches:             shardScope.NewCounterVec("write_batches", "db", "shard"),
		WriteMetrics:             shardScope.NewCounterVec("write_metrics", "db", "shard"),
		WriteFields:              shardScope.NewCounterVec("write_fields", "db", "shard"),
		WriteMetricFailures:      shardScope.NewCounterVec("write_metrics_failures", "db", "shard"),
		FlushInFlight:            shardScope.NewGaugeVec("flush_inflight", "db", "shard"),
		MemDBTotalSize:           shardScope.NewGaugeVec("memdb_total_size", "db", "shard"),
		ActiveMemDBs:             shardScope.NewGaugeVec("active_memdbs", "db", "shard"),
		MemDBFlushFailures:       shardScope.NewCounterVec("memdb_flush_failures", "db", "shard"),
		MemDBFlushDuration:       shardScope.Scope("memdb_flush_duration").NewHistogramVec("db", "shard"),
		IndexDBFlushFailures:     shardScope.NewCounterVec("indexdb_flush_failures", "db", "shard"),
		IndexDBFlushDuration:     shardScope.Scope("indexdb_flush_duration").NewHistogramVec("db", "shard"),
	}
)

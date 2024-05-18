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
	// mete database metric
	metaDBScope = linmetric.StorageRegistry.NewScope("lindb.tsdb.metadb")
	// shard metric
	shardScope = linmetric.StorageRegistry.NewScope("lindb.tsdb.shard")

	// FlushCheckerStatistics represents flush checker statistics.
	FlushCheckerStatistics = struct {
		FlushInFlight *linmetric.GaugeVec // number of family flushing
	}{
		FlushInFlight: shardScope.NewGaugeVec("flush_inflight", "db", "shard"),
	}
)

// IndexDBStatistics represents index database statistics.
type IndexDBStatistics = struct {
	BuildInvertedIndex *linmetric.BoundCounter // build inverted index count
}

// MemDBStatistics represents memory database statistics.
type MemDBStatistics = struct {
	AllocatedPages       *linmetric.BoundCounter // allocate temp memory page success
	AllocatePageFailures *linmetric.BoundCounter // allocate temp memory page failures
}

// DatabaseStatistics represents database statistics.
type DatabaseStatistics struct {
	MetaDBFlushFailures *linmetric.BoundCounter   // flush metadata database failure
	MetaDBFlushDuration *linmetric.BoundHistogram // flush metadata database duration(include count)
}

// TagMetaStatistics represents tag metadata statistics.
type TagMetaStatistics struct {
	GenTagValueIDs        *linmetric.BoundCounter // generate tag value id success
	GenTagValueIDFailures *linmetric.BoundCounter // generate tag value id failure
}

// MetaDBStatistics represents metadata database statistics.
type MetaDBStatistics struct {
	GenMetricIDs        *linmetric.BoundCounter // generate metric id success
	GenMetricIDFailures *linmetric.BoundCounter // generate metric id failure
	GenFieldIDs         *linmetric.BoundCounter // generate field id success
	GenFieldIDFailures  *linmetric.BoundCounter // generate field id failure
	GenTagKeyIDs        *linmetric.BoundCounter // generate tag key id success
	GenTagKeyIDFailures *linmetric.BoundCounter // generate tag key id failure
}

// ShardStatistics represents shard statistics.
type ShardStatistics struct {
	LookupMetricMetaFailures *linmetric.BoundCounter   // lookup meta of metric failure
	IndexDBFlushDuration     *linmetric.BoundHistogram // flush index database duration(include count)
	IndexDBFlushFailures     *linmetric.BoundCounter   // flush index database failure
}

// SegmentStatistics represents segment statistics.
type SegmentStatistics struct {
	IndexDBFlushDuration *linmetric.BoundHistogram // flush index database duration(include count)
	IndexDBFlushFailures *linmetric.BoundCounter   // flush index database failure
}

// FamilyStatistics represents family statistics.
type FamilyStatistics struct {
	ActiveFamilies      *linmetric.BoundGauge     // number of current active families
	WriteBatches        *linmetric.BoundCounter   // write batch count
	WriteMetrics        *linmetric.BoundCounter   // write metric success count
	WriteFields         *linmetric.BoundCounter   // write field data point success count
	WriteMetricFailures *linmetric.BoundCounter   // write metric failures
	MemDBTotalSize      *linmetric.BoundGauge     // total memory size of memory database
	ActiveMemDBs        *linmetric.BoundGauge     // number of current active memory database
	MemDBFlushFailures  *linmetric.BoundCounter   // flush memory database failure
	MemDBFlushDuration  *linmetric.BoundHistogram // flush memory database duration(include count)
}

// NewFamilyStatistics creates a family statistics.
func NewFamilyStatistics(database, shard string) *FamilyStatistics {
	return &FamilyStatistics{
		ActiveFamilies: shardScope.NewGaugeVec("active_families", "db", "shard").
			WithTagValues(database, shard),
		WriteBatches: shardScope.NewCounterVec("write_batches", "db", "shard").
			WithTagValues(database, shard),
		WriteMetrics: shardScope.NewCounterVec("write_metrics", "db", "shard").
			WithTagValues(database, shard),
		WriteFields: shardScope.NewCounterVec("write_fields", "db").
			WithTagValues(database),
		WriteMetricFailures: shardScope.NewCounterVec("write_metrics_failures", "db", "shard").
			WithTagValues(database, shard),
		MemDBTotalSize: shardScope.NewGaugeVec("memdb_total_size", "db", "shard").
			WithTagValues(database, shard),
		ActiveMemDBs: shardScope.NewGaugeVec("active_memdbs", "db", "shard").
			WithTagValues(database, shard),
		MemDBFlushFailures: shardScope.NewCounterVec("memdb_flush_failures", "db", "shard").
			WithTagValues(database, shard),
		MemDBFlushDuration: shardScope.Scope("memdb_flush_duration").NewHistogramVec("db", "shard").
			WithTagValues(database, shard),
	}
}

// NewShardStatistics creates a shard statistics.
func NewShardStatistics(database, shard string) *ShardStatistics {
	return &ShardStatistics{
		LookupMetricMetaFailures: shardScope.NewCounterVec("lookup_metric_meta_failures", "db", "shard").
			WithTagValues(database, shard),
		IndexDBFlushFailures: shardScope.NewCounterVec("indexdb_flush_failures", "db", "shard").
			WithTagValues(database, shard),
		IndexDBFlushDuration: shardScope.Scope("indexdb_flush_duration").NewHistogramVec("db", "shard").
			WithTagValues(database, shard),
	}
}

// NewSegmentStatistics creates a segment statistics.
func NewSegmentStatistics(database, shard, segmentName string) *SegmentStatistics {
	return &SegmentStatistics{
		IndexDBFlushFailures: shardScope.NewCounterVec("indexdb_segment_flush_failures", "db", "shard", "segment").
			WithTagValues(database, shard, segmentName),
		IndexDBFlushDuration: shardScope.Scope("indexdb_segment_flush_duration").NewHistogramVec("db", "shard", "segment").
			WithTagValues(database, shard, segmentName),
	}
}

// NewMetaDBStatistics create a metadata database statistics.
func NewMetaDBStatistics(database string) *MetaDBStatistics {
	return &MetaDBStatistics{
		GenMetricIDs:        metaDBScope.NewCounterVec("gen_metric_ids", "db").WithTagValues(database),
		GenMetricIDFailures: metaDBScope.NewCounterVec("gen_metric_id_failures", "db").WithTagValues(database),
		GenTagKeyIDs:        metaDBScope.NewCounterVec("gen_tag_key_ids", "db").WithTagValues(database),
		GenTagKeyIDFailures: metaDBScope.NewCounterVec("gen_tag_key_id_failures", "db").WithTagValues(database),
		GenFieldIDs:         metaDBScope.NewCounterVec("gen_field_ids", "db").WithTagValues(database),
		GenFieldIDFailures:  metaDBScope.NewCounterVec("gen_field_id_failures", "db").WithTagValues(database),
	}
}

// NewTagMetaStatistics creates a tag metadata statistics.
func NewTagMetaStatistics(database string) *TagMetaStatistics {
	return &TagMetaStatistics{
		GenTagValueIDs:        metaDBScope.NewCounterVec("gen_tag_value_ids", "db").WithTagValues(database),
		GenTagValueIDFailures: metaDBScope.NewCounterVec("gen_tag_value_id_failures", "db").WithTagValues(database),
	}
}

// NewDatabaseStatistics creates a database statistics.
func NewDatabaseStatistics(database string) *DatabaseStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.tsdb.database")
	return &DatabaseStatistics{
		MetaDBFlushFailures: scope.NewCounterVec("metadb_flush_failures", "db").WithTagValues(database),
		MetaDBFlushDuration: scope.Scope("metadb_flush_duration").NewHistogramVec("db").WithTagValues(database),
	}
}

// NewMemDBStatistics create a memory database statistics.
func NewMemDBStatistics(database string) *MemDBStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.tsdb.memdb")
	return &MemDBStatistics{
		AllocatedPages:       scope.NewCounterVec("allocated_pages", "db").WithTagValues(database),
		AllocatePageFailures: scope.NewCounterVec("allocate_page_failures", "db").WithTagValues(database),
	}
}

// NewIndexDBStatistics creates an index database statistics.
func NewIndexDBStatistics(database string) *IndexDBStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.tsdb.indexdb")
	return &IndexDBStatistics{
		BuildInvertedIndex: scope.NewCounterVec("build_inverted_index", "db").WithTagValues(database),
	}
}

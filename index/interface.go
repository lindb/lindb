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

package index

import (
	"io"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=index

// FlushLifeCycle represents kv store flush lifecycle.
type FlushLifeCycle interface {
	// PrepareFlush prepares flush data.
	PrepareFlush()
	// Flush does flush data.
	Flush() error
}

// IndexKVStore represents common kv store for metadata/index data.
type IndexKVStore interface {
	FlushLifeCycle

	// GetOrCreateValue returns unique id for key, if key not exist, creates a new unique id.
	GetOrCreateValue(bucketID uint32, key []byte, createFn func() uint32) (uint32, error)
	// GetValue returns value based on bucket and key.
	GetValue(bucketID uint32, key []byte) (uint32, bool, error)
	// GetValues returns all values for bucket.
	GetValues(bucketID uint32) (ids []uint32, err error)
	// FindValuesByExpr returns values based on filter expr.
	FindValuesByExpr(bucketID uint32, expr stmt.TagFilter) (ids []uint32, err error)
	// CollectKVs collects all keys based on bucket and values.
	CollectKVs(bucketID uint32, values *roaring.Bitmap, result map[uint32]string) error
	// Suggest suggests the kv pairs by prefix.
	Suggest(bucketID uint32, prefix string, limit int) ([]string, error)
}

// MetricSchemaStore represents metric schema(tags/fields etc.) store.
type MetricSchemaStore interface {
	FlushLifeCycle

	// GetSchema returns metric schema by metric id, return nil if not exist.
	GetSchema(id metric.ID) (*metric.Schema, error)
	// genFieldID generates field id if field not exist.
	genFieldID(id metric.ID, fm field.Meta) (field.ID, error)
	// genTagKeyID generates tag key id if tag key not exist.
	genTagKeyID(id metric.ID, tagKey []byte, createFn func() uint32) (tag.KeyID, error)
}

type Notifier interface{}

type Notify interface {
	Notify(notifier Notifier)
}

// MetricMetaDatabase represents metric metadata store.
type MetricMetaDatabase interface {
	io.Closer
	FlushLifeCycle
	Notify
	series.MetricMetaSuggester
	series.TagValueSuggester

	// genMetricID generates metric id if not exist, else return it.
	genMetricID(ns, metricName []byte) (metric.ID, error)
	// genTagKeyID generates tag key id if not exist, else returns it.
	genTagKeyID(metricID metric.ID, tagKey []byte) (tag.KeyID, error)
	// genTagValueID generates tag value id if not exist, else returns it.
	genTagValueID(tagKeyID tag.KeyID, tagValue []byte) (uint32, error)

	// GetSchema returns metric schame by metric id.
	GetSchema(metricID metric.ID) (*metric.Schema, error)
	// GetMetricID returns metric id by namespace and metric name.
	GetMetricID(namespace, metricName string) (metric.ID, error)
	// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	FindTagValueDsByExpr(tagKeyID tag.KeyID, expr stmt.TagFilter) (*roaring.Bitmap, error)
	// FindTagValueIDsForTag get tag value ids for spec tag key of metric,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	FindTagValueIDsForTag(tagKeyID tag.KeyID) (tagValueIDs *roaring.Bitmap, err error)
	// CollectTagValues collects the tag values by tag value ids,
	CollectTagValues(
		tagKeyID tag.KeyID,
		tagValueIDs *roaring.Bitmap,
		tagValues map[uint32]string,
	) error
}

// MetricIndexDatabase represents metric index store.
type MetricIndexDatabase interface {
	io.Closer
	FlushLifeCycle
	Notify
	series.Filter
	flow.GroupingBuilder
}

// MetricIndexSegment represents metric index segment, which manages multiple index database stores.
type MetricIndexSegment interface {
	io.Closer
	FlushLifeCycle
	series.FilterTimeRange

	// GetOrCreateIndex returns the corresponding index database based on familyTime
	GetOrCreateIndex(familyTime int64) (MetricIndexDatabase, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(ctx *flow.ShardExecuteContext) error
}

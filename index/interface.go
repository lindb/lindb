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
	"github.com/lindb/lindb/models"
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
	GetOrCreateValue(bucketID uint32, key []byte, createFn func() (uint32, error)) (id uint32, isNew bool, err error)
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
	genFieldID(id metric.ID, fm field.Meta, limits *models.Limits) (field.ID, error)
	// genTagKeyID generates tag key id if tag key not exist.
	genTagKeyID(id metric.ID, tagKey []byte, limits *models.Limits, createFn func() uint32) (tag.KeyID, error)
}

// MetricMetaDatabase represents metric metadata store.
type MetricMetaDatabase interface {
	io.Closer
	FlushLifeCycle
	series.MetricMetaSuggester
	series.TagValueSuggester

	// Name returns database's name.
	Name() string
	// GenMetricID generates metric id if not exist, else return it.
	GenMetricID(ns, metricName []byte) (metric.ID, error)
	// GenTagKeyID generates tag key id if not exist, else returns it.
	GenTagKeyID(metricID metric.ID, tagKey []byte) (tag.KeyID, error)
	// GenTagValueID generates tag value id if not exist, else returns it.
	GenTagValueID(tagKeyID tag.KeyID, tagValue []byte) (uint32, error)
	// GenFieldID generates field id for metric.
	GenFieldID(metricID metric.ID, f field.Meta) (field.ID, error)

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
	series.Filter
	flow.GroupingBuilder

	// GenSeriesID generates time series id based on tags hash.
	GenSeriesID(metricID metric.ID, row *metric.StorageRow) (seriesID uint32, err error)
}

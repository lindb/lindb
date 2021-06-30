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

package metadb

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=metadb

var metaLogger = logger.GetLogger("tsdb", "MetaDB")

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates the metric id in the memory
	GenMetricID(namespace, metricName string) (metricID uint32, err error)
	// GenFieldID generates the field id in the memory
	// error-case1: field type doesn't matches to before
	// error-case2: there are too many fields
	GenFieldID(namespace, metricName string, fieldName field.Name, fieldType field.Type) (field.ID, error)
	// GenTagKeyID generates the tag key id in the memory
	GenTagKeyID(namespace, metricName, tagKey string) (uint32, error)
}

// IDGetter represents the query ability for metric level, such as metric id, field meta etc.
type IDGetter interface {
	// GetMetricID gets the metric id by namespace and metric name,
	// if not exist return constants.ErrMetricIDNotFound
	GetMetricID(namespace, metricName string) (metricID uint32, err error)
	// GetTagKeyID gets the tag key id by namespace/metric name/tag key key,
	// if not exist return constants.ErrTagKeyIDNotFound
	GetTagKeyID(namespace, metricName, tagKey string) (tagKeyID uint32, err error)
	// GetAllTagKeys returns the all tag keys by namespace/metric name,
	// if not exist return  constants.ErrMetricIDNotFound, constants.ErrMetricBucketNotFound
	GetAllTagKeys(namespace, metricName string) (tags []tag.Meta, err error)
	// GetField gets the field meta by namespace/metric name/field name, if not exist return series.ErrNotFound
	GetField(namespace, metricName string, fieldName field.Name) (field field.Meta, err error)
	// GetAllFields returns the  all fields by namespace/metric name, if not exist return series.ErrNotFound
	GetAllFields(namespace, metricName string) (fields []field.Meta, err error)
}

// Metadata represents all metadata of tsdb, like metric/tag metadata
type Metadata interface {
	io.Closer
	DatabaseName() string
	// MetadataDatabase returns the metric level metadata
	MetadataDatabase() MetadataDatabase
	// TagMetadata returns the tag metadata
	TagMetadata() TagMetadata
	// Flush flushes the metadata to disk
	Flush() error
}

// MetadataDatabase represents the metadata storage includes namespace/metric metadata
type MetadataDatabase interface {
	io.Closer
	IDGetter
	IDGenerator
	series.MetricMetaSuggester

	// SuggestNamespace suggests the namespace by namespace's prefix
	SuggestNamespace(prefix string, limit int) (namespaces []string, err error)
	// Sync syncs the pending metadata update event
	Sync() error
}

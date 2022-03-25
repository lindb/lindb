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
	"context"
	"fmt"
	"sync"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

// for testing
var (
	createMetadataBackendFn = newMetadataBackend
)

var (
	metaDBScope             = linmetric.StorageRegistry.NewScope("lindb.tsdb.metadb")
	getMetricIDCounterVec   = metaDBScope.NewCounterVec("get_metric_ids", "db")
	genMetricIDCounterVec   = metaDBScope.NewCounterVec("gen_metric_ids", "db")
	getTagKeysCounterVec    = metaDBScope.NewCounterVec("get_tag_keys", "db")
	genTagKeyIDCounterVec   = metaDBScope.NewCounterVec("gen_tag_key_ids", "db")
	getFieldsCounterVec     = metaDBScope.NewCounterVec("get_fields", "db")
	genFieldIDCounterVec    = metaDBScope.NewCounterVec("gen_field_ids", "db")
	recoveryMetaWALTimerVec = metaDBScope.Scope("recovery_wal_duration").NewHistogramVec("db")
)

// metadataDatabase implements the MetadataDatabase interface,
// !!!!NOTICE: need cache all tag keys/fields of metric.
type metadataDatabase struct {
	databaseName string
	path         string
	ctx          context.Context
	cancel       context.CancelFunc
	backend      MetadataBackend
	metrics      map[string]MetricMetadata // metadata cache(key: namespace + delimiter + metric-name, value: metric metadata)

	rwMux sync.RWMutex

	statistics struct {
		genMetricIDCounter   *linmetric.BoundCounter
		getMetricIDCounter   *linmetric.BoundCounter
		genTagKeyIDCounter   *linmetric.BoundCounter
		getTagKeysCounter    *linmetric.BoundCounter
		genFieldIDCounter    *linmetric.BoundCounter
		getFieldsCounter     *linmetric.BoundCounter
		recoveryMetaWALTimer *linmetric.BoundHistogram
	}
}

// NewMetadataDatabase creates new metadata database
func NewMetadataDatabase(ctx context.Context, databaseName, parent string) (MetadataDatabase, error) {
	backend, err := createMetadataBackendFn(parent)
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithCancel(ctx)
	mdb := &metadataDatabase{
		databaseName: databaseName,
		path:         parent,
		ctx:          c,
		cancel:       cancel,
		backend:      backend,
		metrics:      make(map[string]MetricMetadata),
	}
	mdb.statistics.genMetricIDCounter = genMetricIDCounterVec.WithTagValues(databaseName)
	mdb.statistics.getMetricIDCounter = getMetricIDCounterVec.WithTagValues(databaseName)
	mdb.statistics.genFieldIDCounter = genFieldIDCounterVec.WithTagValues(databaseName)
	mdb.statistics.getFieldsCounter = getFieldsCounterVec.WithTagValues(databaseName)
	mdb.statistics.genTagKeyIDCounter = genTagKeyIDCounterVec.WithTagValues(databaseName)
	mdb.statistics.getTagKeysCounter = getTagKeysCounterVec.WithTagValues(databaseName)
	mdb.statistics.recoveryMetaWALTimer = recoveryMetaWALTimerVec.WithTagValues(databaseName)

	return mdb, nil
}

// SuggestNamespace suggests the namespace by namespace's prefix
func (mdb *metadataDatabase) SuggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	return mdb.backend.suggestNamespace(prefix, limit)
}

// SuggestMetrics returns suggestions from a given prefix of metricName
func (mdb *metadataDatabase) SuggestMetrics(namespace, metricPrefix string, limit int) ([]string, error) {
	return mdb.backend.suggestMetricName(namespace, metricPrefix, limit)
}

// GetMetricID gets the metric id by namespace and metric name, if not exist return constants.ErrMetricIDNotFound.
func (mdb *metadataDatabase) GetMetricID(namespace, metricName string) (metricID metric.ID, err error) {
	metricMetadata, ok := mdb.getMetricMetadataFromCache(namespace, metricName)
	if ok {
		return metricMetadata.getMetricID(), nil
	}

	// read from meta db
	return mdb.backend.getMetricID(namespace, metricName)
}

// GetAllTagKeys returns the all tag keys by namespace/metric name,
// if not exist return constants.ErrMetricIDNotFound.
func (mdb *metadataDatabase) GetAllTagKeys(namespace, metricName string) (tags tag.Metas, err error) {
	mdb.statistics.getTagKeysCounter.Incr()

	metricMetadata, ok := mdb.getMetricMetadataFromCache(namespace, metricName)

	if ok {
		// need add read lock for getting tag keys from metric metadata.
		mdb.rwMux.RLock()
		tags = metricMetadata.getAllTagKeys()
		mdb.rwMux.RUnlock()
		return
	}

	metricID, err := mdb.backend.getMetricID(namespace, metricName)
	if err != nil {
		return nil, fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
	}

	tags, err = mdb.backend.getAllTagKeys(metricID)
	return
}

// GetTagKeyID gets the tag key id by namespace/metric name/tag key, if not exist return constants.ErrTagKeyIDNotFound
func (mdb *metadataDatabase) GetTagKeyID(namespace, metricName, tagKey string) (tagKeyID tag.KeyID, err error) {
	tagKeys, err := mdb.GetAllTagKeys(namespace, metricName)
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	t, ok := tagKeys.Find(tagKey)
	if ok {
		return t.ID, nil
	}
	return tag.EmptyTagKeyID, constants.ErrTagKeyIDNotFound
}

// GetAllFields returns the all visible fields by namespace/metric name,
// if not exist return series.ErrNotFound
func (mdb *metadataDatabase) GetAllFields(namespace, metricName string) (fields field.Metas, err error) {
	mdb.statistics.getFieldsCounter.Incr()

	metricMetadata, ok := mdb.getMetricMetadataFromCache(namespace, metricName)
	if ok {
		// need add read lock for getting fields from metric metadata.
		mdb.rwMux.RLock()
		fields = metricMetadata.getAllFields()
		mdb.rwMux.RUnlock()
		return
	}

	metricID, err := mdb.GetMetricID(namespace, metricName)
	if err != nil {
		return nil, fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
	}
	fields, _, err = mdb.backend.getAllFields(metricID)
	return
}

// GetAllHistogramFields returns histogram-fields namespace/metric name,
// if not exist return series.ErrNotFound
func (mdb *metadataDatabase) GetAllHistogramFields(namespace, metricName string) (rs field.Metas, err error) {
	fields, err := mdb.GetAllFields(namespace, metricName)
	if err != nil {
		return nil, err
	}
	// with format like __bucket_${boundary}
	for idx := range fields {
		if fields[idx].Type == field.HistogramField {
			rs = append(rs, fields[idx])
		}
	}
	return rs, nil
}

// GetField gets the field meta by namespace/metric name/field name, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetField(namespace, metricName string, fieldName field.Name) (f field.Meta, err error) {
	fields, err := mdb.GetAllFields(namespace, metricName)
	if err != nil {
		return field.Meta{}, err
	}
	f, ok := fields.Find(fieldName)
	if ok {
		return f, nil
	}
	return field.Meta{}, constants.ErrNotFound
}

// GenMetricID generates the metric id in the memory.
// 1) get metric id from memory if existed, if not exist goto 2
// 2) get metric metadata from backend storage, if not exist need create new metric metadata
func (mdb *metadataDatabase) GenMetricID(namespace, metricName string) (metricID metric.ID, err error) {
	// get metric id from memory
	metricMetadata, ok := mdb.getMetricMetadataFromCache(namespace, metricName)
	if ok {
		return metricMetadata.getMetricID(), nil
	}

	key := metric.JoinNamespaceMetric(namespace, metricName)

	// assign metric id from memory, add write lock
	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// double check with memory
	metricMetadata, ok = mdb.metrics[key]
	if ok {
		return metricMetadata.getMetricID(), nil
	}

	metricMetadata, err = mdb.backend.getOrCreateMetricMetadata(namespace, metricName)
	mdb.statistics.genMetricIDCounter.Incr()

	if err != nil {
		return
	}
	mdb.metrics[key] = metricMetadata

	return metricMetadata.getMetricID(), nil
}

// GenFieldID generates the field id in the memory,
// !!!!! NOTICE: metric metadata must be existed in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenFieldID(
	namespace, metricName string,
	fieldName field.Name, fieldType field.Type,
) (fieldID field.ID, err error) {
	if fieldType == field.Unknown {
		return field.EmptyFieldID, series.ErrFieldTypeUnspecified
	}
	metricMetadata, _ := mdb.getMetricMetadataFromCache(namespace, metricName)

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	f, ok := metricMetadata.getField(fieldName)
	if ok {
		mdb.statistics.getFieldsCounter.Incr()
		if f.Type == fieldType {
			return f.ID, nil
		}
		return field.EmptyFieldID, series.ErrWrongFieldType
	}
	// assign new field id, then add field into metric metadata
	fieldMeta, err := metricMetadata.createField(fieldName, fieldType)
	mdb.statistics.genFieldIDCounter.Incr()
	if err != nil {
		return field.EmptyFieldID, err
	}
	// TODO need change?
	err = mdb.backend.saveField(metricMetadata.getMetricID(), fieldMeta)
	if err != nil {
		return field.EmptyFieldID, err
	}
	return fieldMeta.ID, nil
}

// GenTagKeyID generates the tag key id in the memory
// !!!!! NOTICE: metric metadata must be existed in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenTagKeyID(namespace, metricName, tagKey string) (tagKeyID tag.KeyID, err error) {
	metricMetadata, _ := mdb.getMetricMetadataFromCache(namespace, metricName)

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	tagKeyID, ok := metricMetadata.getTagKeyID(tagKey)
	if ok {
		mdb.statistics.genTagKeyIDCounter.Incr()
		return tagKeyID, nil
	}

	err = metricMetadata.checkTagKey(tagKey)
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	// assign new tag key id
	tagKeyID, err = mdb.backend.saveTagKey(metricMetadata.getMetricID(), tagKey)
	mdb.statistics.genTagKeyIDCounter.Incr()
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	metricMetadata.createTagKey(tagKey, tagKeyID)
	return tagKeyID, nil
}

// Sync syncs the backend storage.
func (mdb *metadataDatabase) Sync() error {
	return mdb.backend.sync()
}

// Close closes the resources
func (mdb *metadataDatabase) Close() error {
	mdb.cancel()
	return mdb.backend.Close()
}

// getMetricMetadataFromCache gets metric metadata from memory cache.
func (mdb *metadataDatabase) getMetricMetadataFromCache(namespace, metricName string) (MetricMetadata, bool) {
	mdb.statistics.getMetricIDCounter.Incr()
	key := metric.JoinNamespaceMetric(namespace, metricName)

	// read from memory
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	mdb.rwMux.RUnlock()

	return metricMetadata, ok
}

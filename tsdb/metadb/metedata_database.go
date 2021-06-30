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
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/wal"
)

// for testing
var (
	createMetadataBackend = newMetadataBackend
	createMetaWAL         = wal.NewMetricMetaWAL
)

var (
	genMetricIDCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_metric_id",
			Help: "Generate metric id counter.",
		},
		[]string{"db"},
	)
	genTagKeyIDCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_tag_key_id",
			Help: "Generate tag key id counter.",
		},
		[]string{"db"},
	)
	genFieldIDCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_field_id",
			Help: "Generate field id counter.",
		},
		[]string{"db"},
	)
	recoveryMetaWALTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "recovery_metadata_wal_duration",
			Help:    "Recovery metadata wal duration(ms).",
			Buckets: monitoring.DefaultHistogramBuckets,
		},
		[]string{"db"},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(genMetricIDCounter)
	monitoring.StorageRegistry.MustRegister(genTagKeyIDCounter)
	monitoring.StorageRegistry.MustRegister(genFieldIDCounter)
	monitoring.StorageRegistry.MustRegister(recoveryMetaWALTimer)
}

var (
	syncInterval       = 2 * timeutil.OneSecond
	ErrNeedRecoveryWAL = errors.New("need recovery meta wal")
)

const (
	walPath = "wal"
)

// metadataDatabase implements the MetadataDatabase interface,
// !!!!NOTICE: need cache all tag keys/fields of metric
type metadataDatabase struct {
	databaseName string
	path         string
	ctx          context.Context
	cancel       context.CancelFunc
	backend      MetadataBackend
	metrics      map[string]MetricMetadata // metadata cache(key: namespace + metric-name, value: metric metadata)

	metaWAL wal.MetricMetaWAL

	syncInterval int64

	rwMux sync.RWMutex
}

// NewMetadataDatabase creates new metadata database
func NewMetadataDatabase(ctx context.Context, databaseName, parent string) (MetadataDatabase, error) {
	var err error
	backend, err := createMetadataBackend(parent)
	if err != nil {
		return nil, err
	}
	defer func() {
		// if init metadata database err, need close backend
		if err != nil {
			if err1 := backend.Close(); err1 != nil {
				metaLogger.Info("close metadata backend error when init metadata database",
					logger.String("db", parent), logger.Error(err))
			}
		}
	}()

	metaWAL, err := createMetaWAL(filepath.Join(parent, walPath))
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
		metaWAL:      metaWAL,
		syncInterval: syncInterval,
	}
	// meta recovery
	mdb.metaRecovery()

	// if recovery meta wal fail, need return err
	if mdb.metaWAL.NeedRecovery() {
		err = ErrNeedRecoveryWAL
		return nil, err
	}
	go mdb.checkSync()
	return mdb, nil
}

// SuggestNamespace suggests the namespace by namespace's prefix
func (mdb *metadataDatabase) SuggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	return mdb.backend.suggestNamespace(prefix, limit)
}

// SuggestMetricName suggests the metric name by name's prefix
func (mdb *metadataDatabase) SuggestMetricName(namespace, prefix string, limit int) (metricNames []string, err error) {
	return mdb.backend.suggestMetricName(namespace, prefix, limit)
}

// GetMetricID gets the metric id by namespace and metric name, if not exist return constants.ErrMetricIDNotFound
func (mdb *metadataDatabase) GetMetricID(namespace, metricName string) (metricID uint32, err error) {
	mdb.rwMux.RLock()
	// read from memory
	key := namespace + metricName
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getMetricID(), nil
	}
	mdb.rwMux.RUnlock()

	// read from meta db
	return mdb.backend.getMetricID(namespace, metricName)
}

// GetTagKeyID gets the tag key id by namespace/metric name/tag key key, if not exist return constants.ErrTagKeyIDNotFound
func (mdb *metadataDatabase) GetTagKeyID(namespace, metricName string, tagKey string) (tagKeyID uint32, err error) {
	key := namespace + metricName

	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		tagKeyID, ok = metricMetadata.getTagKeyID(tagKey)
		if ok {
			return
		}
		return 0, fmt.Errorf("%w, tagKey: %s", constants.ErrTagKeyIDNotFound, tagKey)
	}
	mdb.rwMux.RUnlock()

	metricID, err := mdb.backend.getMetricID(namespace, metricName)
	if err != nil {
		return 0, err
	}

	return mdb.backend.getTagKeyID(metricID, tagKey)
}

// GetAllTagKeys returns the all tag keys by namespace/metric name,
// if not exist return constants.ErrMetricIDNotFound, constants.ErrMetricBucketNotFound
func (mdb *metadataDatabase) GetAllTagKeys(namespace, metricName string) (tags []tag.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getAllTagKeys(), nil
	}
	mdb.rwMux.RUnlock()

	metricID, err := mdb.backend.getMetricID(namespace, metricName)
	if err != nil {
		return nil, fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
	}

	return mdb.backend.getAllTagKeys(metricID)
}

// GetField gets the field meta by namespace/metric name/field name, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetField(namespace, metricName string, fieldName field.Name) (f field.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		f, ok = metricMetadata.getField(fieldName)
		if ok {
			return f, nil
		}
		return field.Meta{}, fmt.Errorf("%w ,namespace: %s, metricName: %s, fieldName: %s",
			constants.ErrFieldNotFound, namespace, metricName, fieldName)
	}
	mdb.rwMux.RUnlock()
	metricID, err := mdb.GetMetricID(namespace, metricName)
	if err != nil {
		return field.Meta{}, fmt.Errorf("%w, namespace: %s, metricName: %s, fieldName: %s",
			constants.ErrMetricIDNotFound, namespace, metricName, fieldName)
	}

	// read from db
	return mdb.backend.getField(metricID, fieldName)
}

func (mdb *metadataDatabase) GetAllFields(namespace, metricName string) (fields []field.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getAllFields(), nil
	}
	mdb.rwMux.RUnlock()
	metricID, err := mdb.GetMetricID(namespace, metricName)
	if err != nil {
		return nil, err
	}
	return mdb.backend.getAllFields(metricID)
}

// GenMetricID generates the metric id in the memory.
// 1) get metric id from memory if exist, if not exist goto 2
// 2) get metric metadata from backend storage, if not exist need create new metric metadata
func (mdb *metadataDatabase) GenMetricID(namespace, metricName string) (metricID uint32, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	// get metric id from memory, add read lock
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		mdb.rwMux.RUnlock()
		return metricMetadata.getMetricID(), nil
	}
	mdb.rwMux.RUnlock()

	// assign metric id from memory, add write lock
	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// double check with memory
	metricMetadata, ok = mdb.metrics[key]
	if ok {
		return metricMetadata.getMetricID(), nil
	}

	// load metric metadata from backend storage
	metricMetadata, err = mdb.backend.loadMetricMetadata(namespace, metricName)
	if err == nil {
		// get metric metadata from backend
		mdb.metrics[key] = metricMetadata
		return metricMetadata.getMetricID(), nil
	}
	// isn't not found, return err
	if !errors.Is(err, constants.ErrNotFound) {
		return
	}
	// assign new metric id
	metricID = mdb.backend.genMetricID()

	// append to wal
	if err = mdb.metaWAL.AppendMetric(namespace, metricName, metricID); err != nil {
		// if append wal fail, need rollback assigned metric id, then returns err
		mdb.backend.rollbackMetricID(metricID)
		return 0, err
	}

	mdb.metrics[key] = newMetricMetadata(metricID, 0)

	genMetricIDCounter.WithLabelValues(mdb.databaseName).Inc()

	return metricID, nil
}

// GenFieldID generates the field id in the memory,
// !!!!! NOTICE: metric metadata must be exist in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenFieldID(namespace, metricName string,
	fieldName field.Name, fieldType field.Type,
) (fieldID field.ID, err error) {
	key := namespace + metricName

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	metricMetadata := mdb.metrics[key]
	f, ok := metricMetadata.getField(fieldName)
	if ok {
		if f.Type == fieldType {
			return f.ID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// assign new field id
	fieldID, err = metricMetadata.createField(fieldName, fieldType)
	if err != nil {
		return 0, err
	}

	// append wal
	if err = mdb.metaWAL.AppendField(metricMetadata.getMetricID(), fieldID, fieldName, fieldType); err != nil {
		// if append wal fail, need rollback field id
		metricMetadata.rollbackFieldID(fieldID)
		return 0, err
	}
	// add field into metric metadata
	metricMetadata.addField(field.Meta{
		ID:   fieldID,
		Type: fieldType,
		Name: fieldName,
	})

	genFieldIDCounter.WithLabelValues(mdb.databaseName).Inc()

	return fieldID, nil
}

// GenTagKeyID generates the tag key id in the memory
// !!!!! NOTICE: metric metadata must be exist in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenTagKeyID(namespace, metricName, tagKey string) (tagKeyID uint32, err error) {
	key := namespace + metricName

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	metricMetadata := mdb.metrics[key]
	tagKeyID, ok := metricMetadata.getTagKeyID(tagKey)
	if ok {
		return tagKeyID, nil
	}
	// check tag keys count before create
	if err = metricMetadata.checkTagKeyCount(); err != nil {
		return 0, err
	}
	// assign new tag key id
	tagKeyID = mdb.backend.genTagKeyID()

	// append wal
	if err = mdb.metaWAL.AppendTagKey(metricMetadata.getMetricID(), tagKeyID, tagKey); err != nil {
		// if append wal fail, need rollback tag key id
		mdb.backend.rollbackTagKeyID(tagKeyID)
		return 0, err
	}

	metricMetadata.createTagKey(tagKey, tagKeyID)

	genTagKeyIDCounter.WithLabelValues(mdb.databaseName).Inc()
	return
}

// Sync syncs the bbolt.DB's data file and metadata write ahead log
func (mdb *metadataDatabase) Sync() error {
	if err := mdb.metaWAL.Sync(); err != nil {
		metaLogger.Error("sync meta wal err when invoke sync",
			logger.String("db", mdb.path), logger.Error(err))
	}
	return nil
}

// Close closes the resources
func (mdb *metadataDatabase) Close() error {
	mdb.cancel()

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()

	if err := mdb.metaWAL.Close(); err != nil {
		metaLogger.Error("sync meta wal err when close metadata database",
			logger.String("db", mdb.path), logger.Error(err))
	}

	return mdb.backend.Close()
}

// SuggestMetrics returns suggestions from a given prefix of metricName
func (mdb *metadataDatabase) SuggestMetrics(namespace, metricPrefix string, limit int) ([]string, error) {
	return mdb.SuggestMetricName(namespace, metricPrefix, limit)
}

// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
func (mdb *metadataDatabase) SuggestTagKeys(namespace, metricName, tagKeyPrefix string, limit int) ([]string, error) {
	tags, err := mdb.GetAllTagKeys(namespace, metricName)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	num := 0
	for _, tagMeta := range tags {
		if limit != 0 && num >= limit {
			break
		}
		if tagMeta.Key != "" && strings.HasPrefix(tagMeta.Key, tagKeyPrefix) {
			keys = append(keys, tagMeta.Key)
			num++
		}
	}
	return keys, nil
}

// checkSync checks if need sync pending metadata event in period
func (mdb *metadataDatabase) checkSync() {
	ticker := time.NewTicker(time.Duration(mdb.syncInterval * 1000000))
	for {
		select {
		case <-ticker.C:
			if mdb.metaWAL.NeedRecovery() {
				mdb.metaRecovery()
			}
		case <-mdb.ctx.Done():
			ticker.Stop()
			metaLogger.Info("check metadata event update goroutine exit...", logger.String("db", mdb.path))
			return
		}
	}
}

// metaRecovery recovers meta wal data
func (mdb *metadataDatabase) metaRecovery() {
	startTime := timeutil.Now()
	defer recoveryMetaWALTimer.WithLabelValues(mdb.databaseName).Observe(float64(timeutil.Now() - startTime))

	event := newMetadataUpdateEvent()
	mdb.metaWAL.Recovery(func(namespace, metricName string, metricID uint32) error {
		event.addMetric(namespace, metricName, metricID)

		if event.isFull() {
			if err := mdb.backend.saveMetadata(event); err != nil {
				return err
			}
			event = newMetadataUpdateEvent()
		}
		return nil
	}, func(metricID uint32, fID field.ID, fieldName field.Name, fType field.Type) error {
		event.addField(metricID, field.Meta{
			ID:   fID,
			Type: fType,
			Name: fieldName,
		})

		if event.isFull() {
			if err := mdb.backend.saveMetadata(event); err != nil {
				return err
			}
			event = newMetadataUpdateEvent()
		}
		return nil
	}, func(metricID uint32, tagKeyID uint32, tagKey string) error {
		event.addTagKey(metricID, tag.Meta{
			Key: tagKey,
			ID:  tagKeyID,
		})

		if event.isFull() {
			if err := mdb.backend.saveMetadata(event); err != nil {
				return err
			}
			event = newMetadataUpdateEvent()
		}
		return nil
	}, func() error {
		if !event.isEmpty() {
			if err := mdb.backend.saveMetadata(event); err != nil {
				return err
			}
		}
		return mdb.backend.sync()
	})
}

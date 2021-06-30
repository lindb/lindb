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

package indexdb

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"github.com/lindb/roaring"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/wal"
)

// for testing
var (
	createBackend   = newIDMappingBackend
	createSeriesWAL = wal.NewSeriesWAL
)

var (
	buildInvertedIndexCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "build_inverted_index_counter",
			Help: "Build inverted index counter.",
		},
		[]string{"db"},
	)
	recoverySeriesWALTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "recovery_series_wal_duration",
			Help:    "Recovery series wal duration(ms).",
			Buckets: monitoring.DefaultHistogramBuckets,
		},
		[]string{"db"},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(buildInvertedIndexCounter)
	monitoring.StorageRegistry.MustRegister(recoverySeriesWALTimer)
}

const (
	walPath       = "wal"
	seriesWALPath = "series"
)

var (
	syncInterval       = 2 * timeutil.OneSecond
	ErrNeedRecoveryWAL = errors.New("need recovery series wal")
)

// indexDatabase implements IndexDatabase interface
type indexDatabase struct {
	path             string
	ctx              context.Context
	cancel           context.CancelFunc
	backend          IDMappingBackend           // id mapping backend storage
	metricID2Mapping map[uint32]MetricIDMapping // key: metric id, value: metric id mapping
	metadata         metadb.Metadata            // the metadata for generating ID of metric, field
	index            InvertedIndex

	seriesWAL wal.SeriesWAL

	syncInterval int64

	rwMutex sync.RWMutex // lock of create metric index
}

// NewIndexDatabase creates a new index database
func NewIndexDatabase(ctx context.Context, parent string, metadata metadb.Metadata,
	forwardFamily kv.Family, invertedFamily kv.Family,
) (IndexDatabase, error) {
	var err error
	backend, err := createBackend(parent)
	if err != nil {
		return nil, err
	}
	defer func() {
		// if init index database err, need close backend
		if err != nil {
			if err1 := backend.Close(); err1 != nil {
				indexLogger.Info("close series id mapping backend error when init index database",
					logger.String("db", parent), logger.Error(err))
			}
		}
	}()
	seriesWAL, err := createSeriesWAL(filepath.Join(parent, walPath, seriesWALPath))
	if err != nil {
		return nil, err
	}
	c, cancel := context.WithCancel(ctx)
	db := &indexDatabase{
		path:             parent,
		ctx:              c,
		cancel:           cancel,
		backend:          backend,
		metadata:         metadata,
		metricID2Mapping: make(map[uint32]MetricIDMapping),
		index:            newInvertedIndex(metadata, forwardFamily, invertedFamily),
		seriesWAL:        seriesWAL,
		syncInterval:     syncInterval,
	}

	// series recovery
	db.seriesRecovery()

	// if recovery series wal fail, need return err
	if db.seriesWAL.NeedRecovery() {
		err = ErrNeedRecoveryWAL
		return nil, err
	}

	go db.checkSync()

	return db, nil
}

// SuggestTagValues returns suggestions from given tag key id and prefix of tagValue
func (db *indexDatabase) SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string {
	return db.metadata.TagMetadata().SuggestTagValues(tagKeyID, tagValuePrefix, limit)
}

// GetGroupingContext returns the context of group by
func (db *indexDatabase) GetGroupingContext(tagKeyIDs []uint32, seriesIDs *roaring.Bitmap) (series.GroupingContext, error) {
	return db.index.GetGroupingContext(tagKeyIDs, seriesIDs)
}

// GetOrCreateSeriesID gets series by tags hash, if not exist generate new series id in memory,
// if generate a new series id returns isCreate is true
// if generate fail return err
func (db *indexDatabase) GetOrCreateSeriesID(metricID uint32, tagsHash uint64,
) (seriesID uint32, isCreated bool, err error) {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	metricIDMapping, ok := db.metricID2Mapping[metricID]
	if ok {
		// get series id from memory cache
		seriesID, ok = metricIDMapping.GetSeriesID(tagsHash)
		if ok {
			return seriesID, false, nil
		}
	} else {
		// metric mapping not exist, need load from backend storage
		metricIDMapping, err = db.backend.loadMetricIDMapping(metricID)
		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			return 0, false, err
		}
		// if metric id not exist in backend storage
		if errors.Is(err, constants.ErrNotFound) {
			// create new metric id mapping with 0 sequence
			metricIDMapping = newMetricIDMapping(metricID, 0)
			// cache metric id mapping
			db.metricID2Mapping[metricID] = metricIDMapping
		} else {
			// cache metric id mapping
			db.metricID2Mapping[metricID] = metricIDMapping
			// metric id mapping exist, try get series id from backend storage
			seriesID, err = db.backend.getSeriesID(metricID, tagsHash)
			if err == nil {
				// cache load series id
				metricIDMapping.AddSeriesID(tagsHash, seriesID)
				return seriesID, false, nil
			}
		}
	}
	// throw err in backend storage
	if err != nil && !errors.Is(err, constants.ErrNotFound) {
		return 0, false, err
	}
	// generate new series id
	seriesID = metricIDMapping.GenSeriesID(tagsHash)

	// append to wal
	if err = db.seriesWAL.Append(metricID, tagsHash, seriesID); err != nil {
		// if append wal fail, need rollback assigned series id, then returns err
		metricIDMapping.RemoveSeriesID(tagsHash)
		return 0, false, err
	}
	return seriesID, true, nil
}

// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	return db.index.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
}

// GetSeriesIDsForTag gets series ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	return db.index.GetSeriesIDsForTag(tagKeyID)
}

// GetSeriesIDsForMetric gets series ids for spec metric name
func (db *indexDatabase) GetSeriesIDsForMetric(namespace, metricName string) (*roaring.Bitmap, error) {
	// get all tags under metric
	tags, err := db.metadata.MetadataDatabase().GetAllTagKeys(namespace, metricName)
	if err != nil {
		return nil, err
	}
	tagLength := len(tags)
	if tagLength == 0 {
		// if metric hasn't any tags, returns default series id(0)
		return roaring.BitmapOf(constants.SeriesIDWithoutTags), nil
	}
	tagKeyIDs := make([]uint32, tagLength)
	for idx, tag := range tags {
		tagKeyIDs[idx] = tag.ID
	}
	// get series ids under all tag key ids
	return db.index.GetSeriesIDsForTags(tagKeyIDs)
}

// BuildInvertIndex builds the inverted index for tag value => series ids,
// the tags is considered as a empty key-value pair while tags is nil.
func (db *indexDatabase) BuildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32) {
	db.index.buildInvertIndex(namespace, metricName, tags, seriesID)

	buildInvertedIndexCounter.WithLabelValues(db.metadata.DatabaseName()).Inc()
}

// Flush flushes index data to disk
func (db *indexDatabase) Flush() error {
	if err := db.seriesWAL.Sync(); err != nil {
		indexLogger.Error("sync series wal err when invoke flush",
			logger.String("db", db.path), logger.Error(err))
	}
	//fixme inverted index need add wal???
	return db.index.Flush()
}

// Close closes the database, releases the resources
func (db *indexDatabase) Close() error {
	db.cancel()
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	if err := db.seriesWAL.Close(); err != nil {
		indexLogger.Error("sync series wal err when close index database",
			logger.String("db", db.path), logger.Error(err))
	}
	if err := db.backend.Close(); err != nil {
		return err
	}
	return db.index.Flush()
}

// checkSync checks if need sync pending series event in period
func (db *indexDatabase) checkSync() {
	ticker := time.NewTicker(time.Duration(db.syncInterval * 1000000))
	for {
		select {
		case <-ticker.C:
			if db.seriesWAL.NeedRecovery() {
				db.seriesRecovery()
			}
		case <-db.ctx.Done():
			ticker.Stop()
			indexLogger.Info("check series event update goroutine exit...", logger.String("db", db.path))
			return
		}
	}
}

// seriesRecovery recovers series wal data
func (db *indexDatabase) seriesRecovery() {
	startTime := timeutil.Now()
	defer recoverySeriesWALTimer.WithLabelValues(db.metadata.DatabaseName()).Observe(float64(timeutil.Now() - startTime))

	event := newMappingEvent()
	db.seriesWAL.Recovery(func(metricID uint32, tagsHash uint64, seriesID uint32) error {
		event.addSeriesID(metricID, tagsHash, seriesID)
		if event.isFull() {
			if err := db.backend.saveMapping(event); err != nil {
				return err
			}
			event = newMappingEvent()
		}
		return nil
	}, func() error {
		if !event.isEmpty() {
			if err := db.backend.saveMapping(event); err != nil {
				return err
			}
		}
		return nil
	})
}

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
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/metadb"
)

// for testing
var (
	createBackendFn = newIDMappingBackend
)

var (
	indexDBScope                 = linmetric.NewScope("lindb.tsdb.indexdb")
	buildInvertedIndexCounterVec = indexDBScope.NewCounterVec("build_inverted_index_counter", "db")
)

// indexDatabase implements IndexDatabase interface
type indexDatabase struct {
	path             string
	ctx              context.Context
	cancel           context.CancelFunc
	backend          IDMappingBackend              // id mapping backend storage
	metricID2Mapping map[metric.ID]MetricIDMapping // key: metric id, value: metric id mapping
	metadata         metadb.Metadata               // the metadata for generating ID of metric, field
	index            InvertedIndex

	rwMutex sync.RWMutex // lock of create metric index
}

// NewIndexDatabase creates a new index database
func NewIndexDatabase(ctx context.Context, parent string, metadata metadb.Metadata,
	forwardFamily kv.Family, invertedFamily kv.Family,
) (IndexDatabase, error) {
	var err error
	backend, err := createBackendFn(parent)
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
		metricID2Mapping: make(map[metric.ID]MetricIDMapping),
		index:            newInvertedIndex(metadata, forwardFamily, invertedFamily),
	}

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
func (db *indexDatabase) GetOrCreateSeriesID(metricID metric.ID, tagsHash uint64,
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
		if err != nil {
			return series.EmptySeriesID, false, err
		}
		// cache metric id mapping
		db.metricID2Mapping[metricID] = metricIDMapping
	}
	// metric id mapping exist, try get series id from backend storage
	seriesID, err = db.backend.getSeriesID(metricID, tagsHash)
	if err == nil {
		// cache load series id
		metricIDMapping.AddSeriesID(tagsHash, seriesID)
		return seriesID, false, nil
	}
	// throw err in backend storage
	if err != nil && !errors.Is(err, constants.ErrNotFound) {
		return 0, false, err
	}
	// generate new series id
	// fixme: store seq
	seriesID = metricIDMapping.GenSeriesID(tagsHash)
	// save series id into backend
	if err := db.backend.genSeriesID(metricID, tagsHash, seriesID); err != nil {
		return series.EmptySeriesID, false, err
	}

	return seriesID, true, nil
}

// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec tag key of metric
func (db *indexDatabase) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	return db.index.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
}

// GetSeriesIDsForTag gets series ids for spec tag key of metric
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
// the tags is considered as an empty key-value pair while tags is nil.
func (db *indexDatabase) BuildInvertIndex(
	namespace, metricName string,
	tagIterator *metric.KeyValueIterator,
	seriesID uint32,
) {
	db.index.buildInvertIndex(namespace, metricName, tagIterator, seriesID)

	buildInvertedIndexCounterVec.WithTagValues(db.metadata.DatabaseName()).Incr()
}

// Flush flushes index data to disk
func (db *indexDatabase) Flush() error {
	if err := db.backend.sync(); err != nil {
		return err
	}
	//fixme inverted index need add wal???
	return db.index.Flush()
}

// Close closes the database, releases the resources
func (db *indexDatabase) Close() error {
	db.cancel()
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	if err := db.backend.Close(); err != nil {
		return err
	}
	return db.index.Flush()
}

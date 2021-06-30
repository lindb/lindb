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
	"sync"

	"github.com/lindb/roaring"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/query"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

//go:generate mockgen -source ./inverted_index.go -destination=./inverted_index_mock.go -package=indexdb

// for testing
var (
	newForwardReaderFunc   = invertedindex.NewForwardReader
	newInvertedReaderFunc  = invertedindex.NewInvertedReader
	newForwardFlusherFunc  = invertedindex.NewForwardFlusher
	newInvertedFlusherFunc = invertedindex.NewInvertedFlusher
)

var (
	genTagKeyFailCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_tag_key_id_fail",
			Help: "Generate tag key id fail when build inverted index.",
		},
		[]string{"db"},
	)
	genTagValueFailCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_tag_value_id_fail",
			Help: "Generate tag value id fail when build inverted index.",
		},
		[]string{"db"},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(genTagKeyFailCounter)
	monitoring.StorageRegistry.MustRegister(genTagValueFailCounter)
}

// InvertedIndex represents the tag's inverted index (tag values => series id list)
type InvertedIndex interface {
	// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec metric's tag key
	GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag gets series ids for spec metric's tag key
	GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)
	// GetSeriesIDsForTags gets series ids for spec metric's tag keys
	GetSeriesIDsForTags(tagKeyIDs []uint32) (*roaring.Bitmap, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(tagKeyIDs []uint32, seriesIDs *roaring.Bitmap) (series.GroupingContext, error)
	// buildInvertIndex builds the inverted index for tag value => series ids,
	// the tags is considered as a empty key-value pair while tags is nil.
	buildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32)
	// Flush flushes the inverted-index of tag value id=>series ids under tag key
	Flush() error
}

type invertedIndex struct {
	invertedFamily kv.Family // store tag value inverted index(tag value id=> series ids)
	forwardFamily  kv.Family // store tag value forward index(series id=>tag value id)
	metadata       metadb.Metadata

	mutable   *TagIndexStore
	immutable *TagIndexStore

	rwMutex sync.RWMutex
}

func newInvertedIndex(metadata metadb.Metadata, forwardFamily kv.Family, invertedFamily kv.Family) InvertedIndex {
	return &invertedIndex{
		invertedFamily: invertedFamily,
		forwardFamily:  forwardFamily,
		metadata:       metadata,
		mutable:        NewTagIndexStore(),
	}
}

// GetSeriesIDsByTagValueIDs finds series ids by tag filter expr
func (index *invertedIndex) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		seriesIDs := tagIndex.getSeriesIDsByTagValueIDs(tagValueIDs)
		if seriesIDs != nil {
			result.Or(seriesIDs)
		}
	})

	// read data from kv store
	if err := index.loadSeriesIDsInKV(tagKeyID, func(reader invertedindex.InvertedReader) error {
		seriesIDs, err := reader.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
		if err != nil {
			return err
		}
		result.Or(seriesIDs)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

// GetSeriesIDsForTag get series ids by tagKeyId
func (index *invertedIndex) GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	// get snapshot for getting data
	snapshot := index.forwardFamily.GetSnapshot()
	defer snapshot.Close()
	return index.getSeriesIDsForTag(tagKeyID, snapshot)
}

// getSeriesIDsForTag get series ids by tagKeyId and kv snapshot
func (index *invertedIndex) getSeriesIDsForTag(tagKeyID uint32, snapshot version.Snapshot) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		result.Or(tagIndex.getAllSeriesIDs())
	})

	// read data from kv store
	// try get tag key id from kv store
	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader invertedindex.ForwardReader

	if len(readers) > 0 {
		// found tag data in kv store, try load series ids data
		reader = newForwardReaderFunc(readers)
		seriesIDs, err := reader.GetSeriesIDsForTagKeyID(tagKeyID)
		if err != nil {
			return nil, err
		}
		result.Or(seriesIDs)
	}
	return result, nil
}

// GetSeriesIDsForTags gets series ids for spec metric's tag keys
func (index *invertedIndex) GetSeriesIDsForTags(tagKeyIDs []uint32) (*roaring.Bitmap, error) {
	// get kv store snapshot
	snapshot := index.forwardFamily.GetSnapshot()
	defer snapshot.Close()

	result := roaring.New()
	for _, tagKeyID := range tagKeyIDs {
		seriesIDs, err := index.getSeriesIDsForTag(tagKeyID, snapshot)
		if err != nil {
			return nil, err
		}
		result.Or(seriesIDs)
	}
	return result, nil
}

func (index *invertedIndex) GetGroupingContext(
	tagKeyIDs []uint32,
	seriesIDs *roaring.Bitmap,
) (series.GroupingContext, error) {
	// get kv store snapshot
	snapshot := index.forwardFamily.GetSnapshot()
	defer snapshot.Close()

	scannerMap := make(map[uint32][]series.GroupingScanner)
	for _, tagKeyID := range tagKeyIDs {
		// get grouping scanners by tag key
		scanners, err := index.getGroupingScanners(tagKeyID, seriesIDs, snapshot)
		if err != nil {
			return nil, err
		}
		scannerMap[tagKeyID] = scanners
	}
	return query.NewGroupContext(tagKeyIDs, scannerMap), nil
}

// getGroupingScanners returns the grouping scanner list for tag key, need match series ids
func (index *invertedIndex) getGroupingScanners(
	tagKeyID uint32,
	seriesIDs *roaring.Bitmap,
	snapshot version.Snapshot,
) ([]series.GroupingScanner, error) {
	var result []series.GroupingScanner
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		// get grouping scanner in memory, no err throw
		scanners, _ := tagIndex.GetGroupingScanner(seriesIDs)
		result = append(result, scanners...)
	})

	// read data from kv store
	// try get tag key id from kv store
	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader invertedindex.ForwardReader
	if len(readers) > 0 {
		// found tag data in kv store, try get grouping scanner
		reader = newForwardReaderFunc(readers)
		scanners, err := reader.GetGroupingScanner(tagKeyID, seriesIDs)
		if err != nil {
			return nil, err
		}
		result = append(result, scanners...)
	}
	return result, nil
}

// buildInvertIndex builds the inverted index for tag value => series ids,
// the tags is considered as a empty key-value pair while tags is nil.
func (index *invertedIndex) buildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32) {
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()

	metadataDB := index.metadata.MetadataDatabase()
	tagMetadata := index.metadata.TagMetadata()
	for tagKey, tagValue := range tags {
		tagKeyID, err := metadataDB.GenTagKeyID(namespace, metricName, tagKey)
		if err != nil {
			genTagKeyFailCounter.WithLabelValues(index.metadata.DatabaseName()).Inc()

			indexLogger.Error("gen tag key id fail, ignore index build for this tag key",
				logger.String("namespace", namespace), logger.String("metric", metricName),
				logger.String("key", tagKey), logger.Error(err))
			continue
		}
		tagIndex, ok := index.mutable.Get(tagKeyID)
		if !ok {
			tagIndex = newTagIndex()
			index.mutable.Put(tagKeyID, tagIndex)
		}
		tagValueID, err := tagMetadata.GenTagValueID(tagKeyID, tagValue)
		if err != nil {
			genTagValueFailCounter.WithLabelValues(index.metadata.DatabaseName()).Inc()

			indexLogger.Error("gen tag value id fail, ignore index build for this tag key",
				logger.String("namespace", namespace), logger.String("metric", metricName),
				logger.String("tagKey", tagKey), logger.String("tagValue", tagValue), logger.Error(err))
			continue
		}
		tagIndex.buildInvertedIndex(tagValueID, seriesID)
	}
}

// Flush flushes the inverted-index of tag value id=>series ids under tag key
func (index *invertedIndex) Flush() error {
	if !index.checkFlush() {
		return nil
	}

	// flush immutable data into kv store
	forwardFlusher := index.forwardFamily.NewFlusher()
	forward := newForwardFlusherFunc(forwardFlusher)
	invertedFlusher := index.invertedFamily.NewFlusher()
	inverted := newInvertedFlusherFunc(invertedFlusher)
	if err := index.immutable.WalkEntry(func(key uint32, value TagIndex) error {
		if err := value.flush(key, forward, inverted); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	// commit kv stone meta
	if err := forward.Commit(); err != nil {
		return err
	}
	if err := inverted.Commit(); err != nil {
		return err
	}
	// finally clear immutable
	index.rwMutex.Lock()
	index.immutable = nil
	index.rwMutex.Unlock()
	return nil
}

// checkFlush checks if need do flush job, if need, do switch mutable/immutable
func (index *invertedIndex) checkFlush() bool {
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()

	if index.mutable.Size() == 0 && index.immutable == nil {
		// no new data or immutable is not nil
		return false
	}
	if index.mutable.Size() > 0 && index.immutable == nil {
		// reset mutable, if flush fail immutable is not nil
		index.immutable = index.mutable
		index.mutable = NewTagIndexStore()
	}
	return true
}

// loadTagValueIDsInKV loads series ids in kv store
func (index *invertedIndex) loadSeriesIDsInKV(tagKeyID uint32, fn func(reader invertedindex.InvertedReader) error) error {
	// try get tag key id from kv store
	snapshot := index.invertedFamily.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return err
	}
	var reader invertedindex.InvertedReader
	if len(readers) > 0 {
		// found tag data in kv store, try load series ids data
		reader = newInvertedReaderFunc(readers)
		if err := fn(reader); err != nil {
			return err
		}
	}
	return nil
}

// loadSeriesIDsInMem loads series ids from mutable/immutable store
func (index *invertedIndex) loadSeriesIDsInMem(tagKeyID uint32, fn func(tagIndex TagIndex)) {
	// define get tag series ids func
	getSeriesIDsIDs := func(tagIndexStore *TagIndexStore) {
		tag, ok := tagIndexStore.Get(tagKeyID)
		if ok {
			fn(tag)
		}
	}

	// read data with read lock
	index.rwMutex.RLock()
	defer index.rwMutex.RUnlock()

	getSeriesIDsIDs(index.mutable)
	if index.immutable != nil {
		getSeriesIDsIDs(index.immutable)
	}
}

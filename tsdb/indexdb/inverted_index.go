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

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/tagindex"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./inverted_index.go -destination=./inverted_index_mock.go -package=indexdb

// for testing
var (
	newForwardReaderFunc   = tagindex.NewForwardReader
	newInvertedReaderFunc  = tagindex.NewInvertedReader
	newForwardFlusherFunc  = tagindex.NewForwardFlusher
	newInvertedFlusherFunc = tagindex.NewInvertedFlusher
)

var (
	genTagKeyFailCounterVec   = indexDBScope.NewCounterVec("gen_tag_key_id_fails", "db")
	genTagValueFailCounterVec = indexDBScope.NewCounterVec("gen_tag_value_id_fails", "db")
)

// InvertedIndex represents the tag's inverted index (tag values => series id list)
type InvertedIndex interface {
	// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec tag key of metric
	GetSeriesIDsByTagValueIDs(tagKeyID tag.KeyID, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag gets series ids for spec tag key of metric
	GetSeriesIDsForTag(tagKeyID tag.KeyID) (*roaring.Bitmap, error)
	// GetSeriesIDsForTags gets series ids for spec tag keys of metric
	GetSeriesIDsForTags(tagKeyIDs []tag.KeyID) (*roaring.Bitmap, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(ctx *flow.ShardExecuteContext) error
	// buildInvertIndex builds the inverted index for tag value => series ids,
	// the tags is considered as an empty key-value pair while tags is nil.
	buildInvertIndex(namespace, metricName string, tagIterator *metric.KeyValueIterator, seriesID uint32)
	// Flush flushes the inverted-index of tag value id=>series ids under tag key
	Flush() error
}

type invertedIndex struct {
	invertedFamily kv.Family // store tag value inverted index(tag value id=> series ids)
	forwardFamily  kv.Family // store tag value forward index(series id=>tag value id)
	metadata       metadb.Metadata

	mutable   *TagIndexStore
	immutable *TagIndexStore

	rwMutex                sync.RWMutex
	genTagKeyFailCounter   *linmetric.BoundCounter
	genTagValueFailCounter *linmetric.BoundCounter
}

func newInvertedIndex(metadata metadb.Metadata, forwardFamily, invertedFamily kv.Family) InvertedIndex {
	return &invertedIndex{
		invertedFamily:         invertedFamily,
		forwardFamily:          forwardFamily,
		metadata:               metadata,
		mutable:                NewTagIndexStore(),
		genTagKeyFailCounter:   genTagKeyFailCounterVec.WithTagValues(metadata.DatabaseName()),
		genTagValueFailCounter: genTagValueFailCounterVec.WithTagValues(metadata.DatabaseName()),
	}
}

// GetSeriesIDsByTagValueIDs finds series ids by tag filter expr
func (index *invertedIndex) GetSeriesIDsByTagValueIDs(tagKeyID tag.KeyID, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		seriesIDs := tagIndex.getSeriesIDsByTagValueIDs(tagValueIDs)
		if seriesIDs != nil {
			result.Or(seriesIDs)
		}
	})

	// read data from kv store
	if err := index.loadSeriesIDsInKV(tagKeyID, func(reader tagindex.InvertedReader) error {
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
func (index *invertedIndex) GetSeriesIDsForTag(tagKeyID tag.KeyID) (*roaring.Bitmap, error) {
	// get snapshot for getting data
	snapshot := index.forwardFamily.GetSnapshot()
	defer snapshot.Close()
	return index.getSeriesIDsForTag(tagKeyID, snapshot)
}

// getSeriesIDsForTag get series ids by tagKeyId and kv snapshot
func (index *invertedIndex) getSeriesIDsForTag(tagKeyID tag.KeyID, snapshot version.Snapshot) (*roaring.Bitmap, error) {
	result := roaring.New()
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		result.Or(tagIndex.getAllSeriesIDs())
	})

	// read data from kv store
	// try to get tag key id from kv store
	readers, err := snapshot.FindReaders(uint32(tagKeyID))
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader tagindex.ForwardReader

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

// GetSeriesIDsForTags gets series ids for spec tag keys of metric
func (index *invertedIndex) GetSeriesIDsForTags(tagKeyIDs []tag.KeyID) (*roaring.Bitmap, error) {
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

func (index *invertedIndex) GetGroupingContext(ctx *flow.ShardExecuteContext) error {
	// get kv store snapshot
	snapshot := index.forwardFamily.GetSnapshot()
	defer snapshot.Close()

	scannerMap := make(map[tag.KeyID][]flow.GroupingScanner)
	tagKeyIDs := ctx.StorageExecuteCtx.GroupByTagKeyIDs
	seriesIDs := ctx.SeriesIDsAfterFiltering
	for _, tagKeyID := range tagKeyIDs {
		// get grouping scanners by tag key
		scanners, err := index.getGroupingScanners(tagKeyID, seriesIDs, snapshot)
		if err != nil {
			return err
		}
		scannerMap[tagKeyID] = scanners
	}

	// set context for next execution stage of query
	ctx.GroupingContext = flow.NewGroupContext(tagKeyIDs, scannerMap)
	return nil
}

// getGroupingScanners returns the grouping scanner list for tag key, need match series ids
func (index *invertedIndex) getGroupingScanners(
	tagKeyID tag.KeyID,
	seriesIDs *roaring.Bitmap,
	snapshot version.Snapshot,
) ([]flow.GroupingScanner, error) {
	var result []flow.GroupingScanner
	// read data from mem
	index.loadSeriesIDsInMem(tagKeyID, func(tagIndex TagIndex) {
		// get grouping scanner in memory, no err throw
		scanners, _ := tagIndex.GetGroupingScanner(seriesIDs)
		result = append(result, scanners...)
	})

	// read data from kv store
	// try to get tag key id from kv store
	readers, err := snapshot.FindReaders(uint32(tagKeyID))
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader tagindex.ForwardReader
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
// the tags is considered as an empty key-value pair while tags is nil.
func (index *invertedIndex) buildInvertIndex(namespace, metricName string, tagIterator *metric.KeyValueIterator, seriesID uint32) {
	index.rwMutex.Lock()
	defer index.rwMutex.Unlock()

	metadataDB := index.metadata.MetadataDatabase()
	tagMetadata := index.metadata.TagMetadata()

	for tagIterator.HasNext() {
		tagKey := string(tagIterator.NextKey())
		tagValue := string(tagIterator.NextValue())

		tagKeyID, err := metadataDB.GenTagKeyID(namespace, metricName, tagKey)
		if err != nil {
			index.genTagKeyFailCounter.Incr()

			indexLogger.Error("gen tag key id fail, ignore index build for this tag key",
				logger.String("namespace", namespace), logger.String("metric", metricName),
				logger.String("key", tagKey), logger.Error(err))
			continue
		}
		tagIndex, ok := index.mutable.Get(uint32(tagKeyID))
		if !ok {
			tagIndex = newTagIndex()
			index.mutable.Put(uint32(tagKeyID), tagIndex)
		}
		tagValueID, err := tagMetadata.GenTagValueID(tagKeyID, tagValue)
		if err != nil {
			index.genTagValueFailCounter.Incr()

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
	defer forwardFlusher.Release()

	forward, err := newForwardFlusherFunc(forwardFlusher)
	if err != nil {
		return err
	}
	invertedFlusher := index.invertedFamily.NewFlusher()
	defer invertedFlusher.Release()

	inverted, err := newInvertedFlusherFunc(invertedFlusher)
	if err != nil {
		return err
	}
	if err := index.immutable.WalkEntry(func(key uint32, value TagIndex) error {
		if err := value.flush(key, forward, inverted); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	// commit kv stone meta
	if err := forward.Close(); err != nil {
		return err
	}
	if err := inverted.Close(); err != nil {
		return err
	}
	// finally, clear immutable
	index.rwMutex.Lock()
	index.immutable = nil
	index.rwMutex.Unlock()
	return nil
}

// checkFlush checks if it needs to do flush job, if it needs, do switch mutable/immutable
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
func (index *invertedIndex) loadSeriesIDsInKV(tagKeyID tag.KeyID, fn func(reader tagindex.InvertedReader) error) error {
	// try to get tag key id from kv store
	snapshot := index.invertedFamily.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(uint32(tagKeyID))
	if err != nil {
		// find table.Reader err, return it
		return err
	}
	var reader tagindex.InvertedReader
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
func (index *invertedIndex) loadSeriesIDsInMem(tagKeyID tag.KeyID, fn func(tagIndex TagIndex)) {
	// define get tag series ids func
	getSeriesIDsIDs := func(tagIndexStore *TagIndexStore) {
		tagIndex, ok := tagIndexStore.Get(uint32(tagKeyID))
		if ok {
			fn(tagIndex)
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

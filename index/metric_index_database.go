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
	"context"
	"encoding/binary"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

// for testing
var (
	newInvertedIndexFlusher = v1.NewInvertedIndexFlusher
	newForwardIndexFlusher  = v1.NewForwardIndexFlusher
	bitmapUnmarshal         = encoding.BitmapUnmarshal
	newForwardReader        = v1.NewForwardReader
)

const (
	metricSeriesFamilyName = "metric"
	seriesFamilyName       = "series"
	invertedFamilyName     = "inverted"
	forwardFamilyName      = "forward"
)

// metricIndexDatabase implements MetricIndexDatabase interface.
type metricIndexDatabase struct {
	ctx     context.Context
	cancel  context.CancelFunc
	kvStore kv.Store
	metaDB  MetricMetaDatabase
	series  IndexKVStore // tags => time series ids

	metricInverted *invertedIndex // metric id => time series ids
	inverted       *invertedIndex // tag value id => time series ids
	forward        *forwardIndex  // tag key id => [time seried ids, tag value ids]

	sequenceCache *expirable.LRU[metric.ID, uint32]
	statistics    *metrics.IndexDBStatistics
	logger        logger.Logger

	lock     sync.RWMutex
	flushing atomic.Bool
}

// NewMetricIndexDatabase creates an metric index store.
func NewMetricIndexDatabase(dir string, metaDB MetricMetaDatabase) (MetricIndexDatabase, error) {
	kvStore, err := kv.GetStoreManager().CreateStore(dir, kv.DefaultStoreOption())
	if err != nil {
		return nil, err
	}
	seriesFamily, err := kvStore.CreateFamily(seriesFamilyName, kv.FamilyOption{
		Merger: string(v1.IndexKVMerger),
	})
	if err != nil {
		return nil, err
	}
	invertedFamily, err := kvStore.CreateFamily(invertedFamilyName, kv.FamilyOption{
		Merger: string(v1.InvertedIndexMerger),
	})
	if err != nil {
		return nil, err
	}
	metricFamily, err := kvStore.CreateFamily(metricSeriesFamilyName, kv.FamilyOption{
		Merger: string(v1.InvertedIndexMerger),
	})
	if err != nil {
		return nil, err
	}
	forwadFamily, err := kvStore.CreateFamily(forwardFamilyName, kv.FamilyOption{
		Merger: string(v1.ForwardIndexMerger),
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	index := &metricIndexDatabase{
		ctx:            ctx,
		cancel:         cancel,
		kvStore:        kvStore,
		metaDB:         metaDB,
		series:         NewIndexKVStore(seriesFamily, 10000, 10*time.Minute),
		metricInverted: newInvertedIndex(metricFamily),
		inverted:       newInvertedIndex(invertedFamily),
		forward:        newForwardIndex(forwadFamily),
		statistics:     metrics.NewIndexDBStatistics(metaDB.Name()),
		sequenceCache:  expirable.NewLRU[metric.ID, uint32](100000, nil, time.Hour),
		logger:         logger.GetLogger("Index", "IndexDB"),
	}
	return index, nil
}

func (index *metricIndexDatabase) createSeriesID(metricID metric.ID) (seriesID uint32) {
	sequence, ok := index.sequenceCache.Get(metricID)
	if ok {
		return sequence + 1
	}
	seriesIDs, newSeriesErr := index.metricInverted.getSeriesIDs(uint32(metricID))
	if newSeriesErr == nil {
		if seriesIDs.IsEmpty() {
			return 0
		}
		return seriesIDs.Maximum() + 1
	}
	return 0
}

// GenSeriesID generates time series id based on tags hash.
func (index *metricIndexDatabase) GenSeriesID(metricID metric.ID, row *metric.StorageRow) (seriesID uint32, err error) {
	var isNewSeries bool
	var scratch [8]byte
	tagsHash := row.TagsHash()
	binary.LittleEndian.PutUint64(scratch[:], tagsHash)

	seriesID, isNewSeries, err = index.series.GetOrCreateValue(uint32(metricID), scratch[:], func() (uint32, error) {
		return index.createSeriesID(metricID), nil
	})
	if err == nil && isNewSeries {
		limits := models.GetDatabaseLimits(index.metaDB.Name())
		seriesLimit := limits.GetSeriesLimit(strutil.ByteSlice2String(row.NameSpace()), strutil.ByteSlice2String(row.Name()))
		if seriesLimit > 0 && seriesLimit < seriesID {
			return 0, constants.ErrTooManySeries
		}
		// if new series do inverted index build
		index.sequenceCache.Add(metricID, seriesID)

		// write metric inverted index
		index.lock.Lock()
		index.metricInverted.put(uint32(metricID), seriesID)
		index.lock.Unlock()

		if row.TagsLen() > 0 {
			// write tag related index
			index.buildInvertIndex(metricID, row.NewKeyValueIterator(), seriesID)
			index.statistics.BuildInvertedIndex.Incr()
		}
	}
	return
}

func (index *metricIndexDatabase) GetSeriesIDsForMetric(metricID metric.ID) (*roaring.Bitmap, error) {
	return index.metricInverted.getSeriesIDs(uint32(metricID))
}

func (index *metricIndexDatabase) GetSeriesIDsForTag(tagKeyID tag.KeyID) (*roaring.Bitmap, error) {
	return index.forward.findSeriesIDsForTag(tagKeyID)
}

func (index *metricIndexDatabase) GetSeriesIDsByTagValueIDs(tagKeyID tag.KeyID, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	return index.inverted.findSeriesIDsByKeys(tagValueIDs)
}

// GetGroupingContext returns the context of group by
func (index *metricIndexDatabase) GetGroupingContext(
	groupingTags tag.Metas, seriesIDs *roaring.Bitmap,
) (*roaring.Bitmap, flow.GroupingContext, error) {
	return index.forward.GetGroupingContext(groupingTags, seriesIDs)
}

func (index *metricIndexDatabase) PrepareFlush() {
	index.metricInverted.prepareFlush()
	index.forward.prepareFlush()
	index.inverted.prepareFlush()
	index.series.PrepareFlush()
}

func (index *metricIndexDatabase) Flush() error {
	if (&index.flushing).CompareAndSwap(false, true) {
		defer func() {
			index.flushing.Store(false)
		}()
		if err := index.metricInverted.flush(); err != nil {
			return err
		}
		if err := index.forward.flush(); err != nil {
			return err
		}
		if err := index.inverted.flush(); err != nil {
			return err
		}
		if err := index.series.Flush(); err != nil {
			return err
		}
	}
	// TODO: add wait?
	return nil
}

func (index *metricIndexDatabase) Close() error {
	index.cancel()
	return kv.GetStoreManager().CloseStore(index.kvStore.Name())
}

func (index *metricIndexDatabase) buildInvertIndex(metricID metric.ID,
	tags *metric.KeyValueIterator, seriesID uint32,
) {
	for tags.HasNext() {
		key := tags.NextKey()
		tagKeyID, err := index.metaDB.GenTagKeyID(metricID, key)
		if err != nil {
			index.logger.Error("gen tag key id error when build inverted index",
				logger.String("tagKey", string(key)), logger.Error(err))
			continue
		}
		value := tags.NextValue()
		tagValueID, err := index.metaDB.GenTagValueID(tagKeyID, value)
		if err != nil {
			index.logger.Error("gen tag value id error when build inverted index",
				logger.String("tagKey", string(key)), logger.String("tagValue", string(value)), logger.Error(err))
			continue
		}
		index.lock.Lock()
		// write tag value inverted index
		index.inverted.put(tagValueID, seriesID)
		// write tag key forward index
		index.forward.put(uint32(tagKeyID), tagValueID, seriesID)
		index.lock.Unlock()
	}
}

type invertedIndex struct {
	// tag value id => time series ids(bitmap)
	mutable   *imap.IntMap[*roaring.Bitmap]
	immutable *imap.IntMap[*roaring.Bitmap]

	family kv.Family
	lock   sync.RWMutex
}

func newInvertedIndex(family kv.Family) *invertedIndex {
	return &invertedIndex{
		family:  family,
		mutable: imap.NewIntMap[*roaring.Bitmap](),
	}
}

func (ii *invertedIndex) put(key, seriesID uint32) {
	ii.lock.Lock()
	defer ii.lock.Unlock()

	seriesIDs, ok := ii.mutable.Get(key)
	if !ok {
		// create new series ids for new tag value
		seriesIDs = roaring.NewBitmap()
		ii.mutable.Put(key, seriesIDs)
	}
	seriesIDs.Add(seriesID)
}

func (ii *invertedIndex) getSeriesIDs(key uint32) (*roaring.Bitmap, error) {
	snapshot := ii.family.GetSnapshot()
	defer snapshot.Close()

	result := roaring.New()
	seriesIDs := roaring.New()
	if err := snapshot.Load(key, func(value []byte) error {
		if _, err := bitmapUnmarshal(seriesIDs, value); err != nil {
			return err
		}
		result.Or(seriesIDs)
		seriesIDs.Clear()
		return nil
	}); err != nil {
		return nil, err
	}
	ii.findSeriesIDsByKeyFromMem(key, result)
	return result, nil
}

func (ii *invertedIndex) findSeriesIDsByKeys(keys *roaring.Bitmap) (*roaring.Bitmap, error) {
	snapshot := ii.family.GetSnapshot()
	defer snapshot.Close()

	result := roaring.New()
	seriesIDs := roaring.New()
	it := keys.Iterator()
	for it.HasNext() {
		key := it.Next()
		if err := snapshot.Load(key, func(value []byte) error {
			if _, err := bitmapUnmarshal(seriesIDs, value); err != nil {
				return err
			}
			result.Or(seriesIDs)
			seriesIDs.Clear()
			return nil
		}); err != nil {
			return nil, err
		}
		ii.findSeriesIDsByKeyFromMem(key, result)
	}
	return result, nil
}

func (ii *invertedIndex) findSeriesIDsByKeyFromMem(key uint32, seriesIDs *roaring.Bitmap) {
	findSeriesIDs := func(mem *imap.IntMap[*roaring.Bitmap]) {
		if mem == nil {
			return
		}
		sIDs, ok := mem.Get(key)
		if ok {
			seriesIDs.Or(sIDs)
		}
	}
	ii.lock.RLock()
	defer ii.lock.RUnlock()

	findSeriesIDs(ii.mutable)
	findSeriesIDs(ii.immutable)
}

func (ii *invertedIndex) prepareFlush() {
	ii.lock.Lock()
	defer ii.lock.Unlock()
	if ii.immutable == nil {
		ii.immutable = ii.mutable
		ii.mutable = imap.NewIntMap[*roaring.Bitmap]()
	}
}

func (ii *invertedIndex) needFlush() bool {
	ii.lock.RLock()
	defer ii.lock.RUnlock()

	return ii.immutable != nil && !ii.immutable.IsEmpty()
}

func (ii *invertedIndex) flush() (err error) {
	if !ii.needFlush() {
		return nil
	}

	kvFlusher := ii.family.NewFlusher()
	defer kvFlusher.Release()
	flusher, err := newInvertedIndexFlusher(kvFlusher)
	if err != nil {
		return err
	}

	err = ii.immutable.WalkEntry(func(tagValueID uint32, seriesIDs *roaring.Bitmap) error {
		if seriesIDs.IsEmpty() {
			return nil
		}
		flusher.Prepare(tagValueID)
		if err0 := flusher.Write(seriesIDs); err0 != nil {
			return err0
		}
		return flusher.Commit()
	})
	if err != nil {
		return err
	}
	err = flusher.Close()
	if err != nil {
		return err
	}

	ii.lock.Lock()
	ii.immutable = nil
	ii.lock.Unlock()
	return nil
}

type forwardIndex struct {
	mutable   *imap.IntMap[*imap.IntMap[uint32]]
	immutable *imap.IntMap[*imap.IntMap[uint32]]

	family kv.Family // tag key id => [time series ids -> tag value ids)

	lock sync.RWMutex
}

func newForwardIndex(family kv.Family) *forwardIndex {
	return &forwardIndex{
		mutable: imap.NewIntMap[*imap.IntMap[uint32]](),
		family:  family,
	}
}

func (fi *forwardIndex) put(tagKeyID, tagValueID, seriesID uint32) {
	fi.lock.Lock()
	defer fi.lock.Unlock()

	forwardEntry, ok := fi.mutable.Get(tagKeyID)
	if !ok {
		forwardEntry = imap.NewIntMap[uint32]()
		fi.mutable.Put(tagKeyID, forwardEntry)
	}
	// build forward index, because series id is a unique id, so just put into forward index
	forwardEntry.PutIfNotExist(seriesID, tagValueID)
}

func (fi *forwardIndex) findSeriesIDsForTag(tagKeyID tag.KeyID) (*roaring.Bitmap, error) {
	snapshot := fi.family.GetSnapshot()
	defer snapshot.Close()

	result := roaring.New()
	// read data from mem
	fi.loadSeriesIDsInMem(tagKeyID, func(tagIndex *imap.IntMap[uint32]) {
		result.Or(tagIndex.Keys())
	})

	// read data from kv store
	// try to get tag key id from kv store
	readers, err := snapshot.FindReaders(uint32(tagKeyID))
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader v1.ForwardReader

	if len(readers) > 0 {
		// found tag data in kv store, try load series ids data
		reader = newForwardReader(readers)
		seriesIDs, err := reader.GetSeriesIDsForTagKeyID(tagKeyID)
		if err != nil {
			return nil, err
		}
		result.Or(seriesIDs)
	}
	return result, nil
}

// GetGroupingContext returns the context of group by
func (fi *forwardIndex) GetGroupingContext(
	groupingTags tag.Metas, seriesIDs *roaring.Bitmap,
) (*roaring.Bitmap, flow.GroupingContext, error) {
	snapshot := fi.family.GetSnapshot()
	defer snapshot.Close()

	scannerMap := make(map[tag.KeyID][]flow.GroupingScanner)
	finalSeriesIDs := seriesIDs.Clone()

	for _, groupingTag := range groupingTags {
		// get grouping scanners by tag key
		scanners, err := fi.getGroupingScanners(groupingTag.ID, seriesIDs, snapshot)
		if err != nil {
			return nil, nil, err
		}
		seriesIDsForCurrentTagKey := roaring.New()
		for idx := range scanners {
			seriesIDsForCurrentTagKey.Or(scanners[idx].GetSeriesIDs())
		}
		finalSeriesIDs.And(seriesIDsForCurrentTagKey)
		if finalSeriesIDs.IsEmpty() {
			return finalSeriesIDs, nil, constants.ErrNotFound
		}
		scannerMap[groupingTag.ID] = scanners
	}

	return finalSeriesIDs, flow.NewGroupContext(groupingTags, scannerMap), nil
}

// getGroupingScanners returns the grouping scanner list for tag key, need match series ids
func (fi *forwardIndex) getGroupingScanners(
	tagKeyID tag.KeyID,
	seriesIDs *roaring.Bitmap,
	snapshot version.Snapshot,
) ([]flow.GroupingScanner, error) {
	var result []flow.GroupingScanner
	// read data from mem
	fi.loadSeriesIDsInMem(tagKeyID, func(tagIndex *imap.IntMap[uint32]) {
		// check reader if it has series ids(after filtering)
		finalSeriesIDs := roaring.FastAnd(seriesIDs, tagIndex.Keys())
		if finalSeriesIDs.IsEmpty() {
			// not found
			return
		}
		result = append(result, &memGroupingScanner{forward: tagIndex, withLock: fi.withLock})
	})

	// read data from kv store
	// try to get tag key id from kv store
	readers, err := snapshot.FindReaders(uint32(tagKeyID))
	if err != nil {
		// find table.Reader err, return it
		return nil, err
	}
	var reader v1.ForwardReader
	if len(readers) > 0 {
		// found tag data in kv store, try get grouping scanner
		reader = newForwardReader(readers)
		scanners, err := reader.GetGroupingScanner(tagKeyID, seriesIDs)
		if err != nil {
			return nil, err
		}
		result = append(result, scanners...)
	}
	return result, nil
}

// loadSeriesIDsInMem loads series ids from mutable/immutable store
func (fi *forwardIndex) loadSeriesIDsInMem(tagKeyID tag.KeyID, fn func(tagIndex *imap.IntMap[uint32])) {
	// define get tag series ids func
	getSeriesIDsIDs := func(mem *imap.IntMap[*imap.IntMap[uint32]]) {
		if mem == nil {
			return
		}
		if tagIndex, ok := mem.Get(uint32(tagKeyID)); ok {
			fn(tagIndex)
		}
	}

	// read data with read lock
	fi.lock.RLock()
	defer fi.lock.RUnlock()

	getSeriesIDsIDs(fi.mutable)
	getSeriesIDsIDs(fi.immutable)
}

// withLock retrieves the lock of inverted index, and returns the release function.
func (fi *forwardIndex) withLock() (release func()) {
	fi.lock.RLock()

	return fi.lock.RUnlock
}

func (fi *forwardIndex) prepareFlush() {
	fi.lock.Lock()
	defer fi.lock.Unlock()
	if fi.immutable == nil {
		fi.immutable = fi.mutable
		fi.mutable = imap.NewIntMap[*imap.IntMap[uint32]]()
	}
}

func (fi *forwardIndex) needFlush() bool {
	fi.lock.RLock()
	defer fi.lock.RUnlock()

	return fi.immutable != nil && !fi.immutable.IsEmpty()
}

func (fi *forwardIndex) flush() (err error) {
	if !fi.needFlush() {
		return nil
	}
	kvFlusher := fi.family.NewFlusher()
	defer kvFlusher.Release()

	flusher, err := newForwardIndexFlusher(kvFlusher)
	if err != nil {
		return err
	}

	err = fi.immutable.WalkEntry(func(tagKeyID uint32, forward *imap.IntMap[uint32]) error {
		if forward.Keys().IsEmpty() {
			return nil
		}
		flusher.Prepare(tagKeyID)
		if err0 := flusher.WriteSeriesIDs(forward.Keys()); err0 != nil {
			return err0
		}
		values := forward.Values()
		for _, tagValueIDs := range values {
			if err0 := flusher.WriteTagValueIDs(tagValueIDs); err0 != nil {
				return err0
			}
		}
		return flusher.Commit()
	})
	if err != nil {
		return err
	}
	err = flusher.Close()
	if err != nil {
		return err
	}

	fi.lock.Lock()
	fi.immutable = nil
	fi.lock.Unlock()
	return nil
}

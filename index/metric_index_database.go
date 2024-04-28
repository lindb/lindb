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
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/imap"
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
	series  IndexKVStore // tags -> time series ids

	metricInverted *invertedIndex // metric id -> time series ids
	inverted       *invertedIndex // tag value id -> time series ids
	forward        *forwardIndex  // tag key id -> [time seried ids, tag value ids]

	sequenceCache *expirable.LRU[metric.ID, uint32]

	lock     sync.RWMutex
	worker   *NotifyWorker
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
		sequenceCache:  expirable.NewLRU[metric.ID, uint32](100000, nil, time.Hour),
	}
	index.worker = NewWorker(ctx, dir, 100*time.Millisecond, index.handle)
	return index, nil
}

func (index *metricIndexDatabase) Notify(n Notifier) {
	if n == nil {
		return
	}
	switch notifier := n.(type) {
	case *MetaNotifier:
		// notify need generate metric id
		metaNotify := GetMetaNotifier()
		metaNotify.Namespace = notifier.Namespace
		metaNotify.MetricName = notifier.MetricName
		metaNotify.Callback = func(mid uint32, err error) {
			if err != nil {
				notifier.Callback(0, err)
				return
			}
			notifier.MetricID = metric.ID(mid)
			// after get metric id, notify generate series id/index
			index.worker.Notify(notifier)
		}

		index.metaDB.Notify(metaNotify)
	default:
		index.worker.Notify(notifier)
	}
}

func (index *metricIndexDatabase) handle(n Notifier) {
	switch notifier := n.(type) {
	case *MetaNotifier:
		metricID := notifier.MetricID

		var isNewSeries bool
		var scratch [8]byte
		var seriesIDs *roaring.Bitmap
		var err error
		var newSeriesErr error
		var seriesID uint32
		binary.LittleEndian.PutUint64(scratch[:], notifier.TagHash)
		seriesID, err = index.series.GetOrCreateValue(uint32(metricID), scratch[:], func() uint32 {
			isNewSeries = true
			sequence, ok := index.sequenceCache.Get(metricID)
			if ok {
				return sequence + 1
			}
			seriesIDs, newSeriesErr = index.metricInverted.getSeriesIDs(uint32(metricID))
			if newSeriesErr == nil {
				if seriesIDs.IsEmpty() {
					return 0
				}
				return seriesIDs.Maximum() + 1
			}
			isNewSeries = false
			return 0
		})
		if err == nil {
			// if err is nil, try set err using new series error
			err = newSeriesErr
		}
		if err == nil && isNewSeries {
			// if new series do inverted index build
			index.sequenceCache.Add(metricID, seriesID)
			// TODO: add limit

			// write metric inverted index
			index.lock.Lock()
			index.metricInverted.put(uint32(metricID), seriesID)
			index.lock.Unlock()

			if len(notifier.Tags) > 0 {
				// write tag related index
				index.buildInvertIndex(metricID, notifier.Tags, seriesID, models.NewDefaultLimits()) // FIXME: add limit
			}
		}
		notifier.Callback(seriesID, err)
	case *FlushNotifier:
		if index.flushing.CompareAndSwap(false, true) {
			index.PrepareFlush()
			go func() {
				notifier.Callback(index.Flush())
			}()
		} else {
			notifier.Callback(nil)
		}
	}
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
func (index *metricIndexDatabase) GetGroupingContext(ctx *flow.ShardExecuteContext) (map[tag.KeyID][]flow.GroupingScanner, error) {
	return index.forward.GetGroupingContext(ctx)
}

func (index *metricIndexDatabase) PrepareFlush() {
	index.metricInverted.prepareFlush()
	index.forward.prepareFlush()
	index.inverted.prepareFlush()
	index.series.PrepareFlush()
}

func (index *metricIndexDatabase) Flush() error {
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
	return nil
}

func (index *metricIndexDatabase) Close() error {
	index.cancel()
	index.worker.Shutdown()
	return kv.GetStoreManager().CloseStore(index.kvStore.Name())
}

func (index *metricIndexDatabase) buildInvertIndex(metricID metric.ID,
	tags tag.Tags, seriesID uint32, _ *models.Limits) {
	tagNotifier := GetTagNotifier()
	tagNotifier.tags = tags
	tagNotifier.metricID = metricID
	// callback function after generate tag meta(meta tag goroutine under metric meta database)
	tagNotifier.buildIndex = func(tagKeyID, tagValueID uint32) {
		index.lock.Lock()
		// write tag value inverted index
		index.inverted.put(tagValueID, seriesID)
		// write tag key forward index
		index.forward.put(tagKeyID, tagValueID, seriesID)
		index.lock.Unlock()
	}
	// notify meta db to generate tag meta
	index.metaDB.Notify(tagNotifier)
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

func (ii *invertedIndex) flush() (err error) {
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
func (fi *forwardIndex) GetGroupingContext(ctx *flow.ShardExecuteContext) (map[tag.KeyID][]flow.GroupingScanner, error) {
	snapshot := fi.family.GetSnapshot()
	defer snapshot.Close()
	scannerMap := make(map[tag.KeyID][]flow.GroupingScanner)
	tagKeyIDs := ctx.StorageExecuteCtx.GroupByTagKeyIDs
	seriesIDs := ctx.SeriesIDsAfterFiltering
	finalSeriesIDs := seriesIDs.Clone()
	defer func() {
		// maybe filtering some series ids that is result of filtering.
		// if not found, return empty series ids.
		ctx.SeriesIDsAfterFiltering = finalSeriesIDs
	}()
	for _, tagKeyID := range tagKeyIDs {
		// get grouping scanners by tag key
		scanners, err := fi.getGroupingScanners(tagKeyID, seriesIDs, snapshot)
		if err != nil {
			return nil, err
		}
		seriesIDsForCurrentTagKey := roaring.New()
		for idx := range scanners {
			seriesIDsForCurrentTagKey.Or(scanners[idx].GetSeriesIDs())
		}
		finalSeriesIDs.And(seriesIDsForCurrentTagKey)
		if finalSeriesIDs.IsEmpty() {
			return nil, constants.ErrNotFound
		}
		scannerMap[tagKeyID] = scanners
	}

	// set context for next execution stage of query
	return scannerMap, nil
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

func (fi *forwardIndex) flush() (err error) {
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

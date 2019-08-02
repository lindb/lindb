package memdb

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metrictbl"
	"github.com/lindb/lindb/tsdb/series"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock_test.go -package memdb

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	// getMetricID returns the metricID
	getMetricID() uint32
	// write writes the metric
	write(metric *pb.Metric, writeCtx writeContext) error
	// setMaxTagsLimit sets the max tags-limit
	setMaxTagsLimit(limit uint32)
	// isEmpty detects whether if tags number is empty or not.
	isEmpty() bool
	// evict scans all tsStore and removes which are not in use for a while.
	evict()
	// getTagsCount return the tags count
	getTagsCount() int
	// flushMetricsTo flushes metric-block of mStore to the writer.
	flushMetricsTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) error
	// resetVersion moves the current running mutable index to immutable list,
	// then creates a new mutable map.
	resetVersion() error
	// findSeriesIDsByExpr finds series ids by tag filter expr and timeRange
	findSeriesIDsByExpr(expr stmt.TagFilter, timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
	// getSeriesIDsForTag get series ids by tagKey and timeRange
	getSeriesIDsForTag(tagKey string, timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
}

// metricStore is composed of the immutable part and mutable part of indexes.
// evictor scans the index to check which of them should be purged from the mutable part.
// flusher flushes both the immutable and mutable index to disk,
// after flushing, the immutable part will be removed.
type metricStore struct {
	immutable       []tagIndexINTF // immutable index that has not been flushed to disk
	mutex4Immutable sync.RWMutex   // sync.RWMutex for immutable index
	mutable         tagIndexINTF   // current mutable index in use
	mutex4Mutable   sync.RWMutex   // sync.RWMutex for mutable index
	maxTagsLimit    uint32         // maximum number of combinations of tags
	metricID        uint32         // persistent on the disk
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore(metricID uint32) mStoreINTF {
	ms := metricStore{
		metricID:     metricID,
		mutable:      newTagIndex(),
		maxTagsLimit: defaultMaxTagsLimit}
	return &ms
}

// getMetricID returns the metricID
func (ms *metricStore) getMetricID() uint32 {
	return atomic.LoadUint32(&ms.metricID)
}

// write writes the metric to the tStore
func (ms *metricStore) write(metric *pb.Metric, writeCtx writeContext) error {
	if ms.isFull() {
		return models.ErrTooManyTags
	}
	var err error
	tStore, ok := ms.getTStore(metric.Tags)
	if !ok {
		ms.mutex4Mutable.Lock()
		tStore, err = ms.mutable.getOrCreateTStore(metric.Tags)
		if err != nil {
			ms.mutex4Mutable.Unlock()
			return err
		}
		ms.mutex4Mutable.Unlock()
	}
	return tStore.write(metric, writeCtx)
}

// setMaxTagsLimit sets the max tags-limit of the metricStore
func (ms *metricStore) setMaxTagsLimit(limit uint32) {
	atomic.StoreUint32(&ms.maxTagsLimit, limit)
}

// getMaxTagsLimit return the max tags limit without race condition.
func (ms *metricStore) getMaxTagsLimit() uint32 {
	return atomic.LoadUint32(&ms.maxTagsLimit)
}

// getTStore returns timeSeriesStore, return false when not exist.
func (ms *metricStore) getTStore(tags string) (tStore tStoreINTF, ok bool) {
	ms.mutex4Mutable.RLock()
	tStore, ok = ms.mutable.getTStore(tags)
	ms.mutex4Mutable.RUnlock()
	return
}

// getTagsCount return the map's length.
func (ms *metricStore) getTagsCount() int {
	ms.mutex4Mutable.RLock()
	length := ms.mutable.len()
	ms.mutex4Mutable.RUnlock()
	return length
}

// isFull detects if timeSeriesMap exceeds the tagsID limitation.
func (ms *metricStore) isFull() bool {
	return uint32(ms.getTagsCount()) >= ms.getMaxTagsLimit()
}

// isEmpty detects if timeSeriesMap is empty or not.
func (ms *metricStore) isEmpty() bool {
	return ms.getTagsCount() == 0
}

// evict scans all tsStore and removes which are not in use for a while.
func (ms *metricStore) evict() {
	var (
		evictList            []uint32
		doubleCheckEvictList []uint32
	)
	// first check
	ms.mutex4Mutable.RLock()
	for seriesID, tStore := range ms.mutable.allTStores() {
		if tStore.isNoData() && tStore.isExpired() {
			evictList = append(evictList, seriesID)
		}
	}
	ms.mutex4Mutable.RUnlock()
	// double check
	ms.mutex4Mutable.Lock()
	for _, seriesID := range evictList {
		tStore, ok := ms.mutable.getTStoreBySeriesID(seriesID)
		if !ok {
			continue
		}
		if tStore.isNoData() && tStore.isExpired() {
			doubleCheckEvictList = append(doubleCheckEvictList, seriesID)
		}
	}
	ms.mutable.removeTStores(doubleCheckEvictList...)
	ms.mutex4Mutable.Unlock()
}

// resetVersion moves the mutable index to immutable list, then creates a new active index.
func (ms *metricStore) resetVersion() error {
	ms.mutex4Mutable.Lock()
	if ms.mutable.getVersion()+minIntervalForResetMetricStore > timeutil.Now() {
		ms.mutex4Mutable.Unlock()
		return fmt.Errorf("reset version too frequently")
	}
	oldMutable := ms.mutable
	ms.mutable = newTagIndex()
	ms.mutex4Mutable.Unlock()

	ms.mutex4Immutable.Lock()
	ms.immutable = append(ms.immutable, oldMutable)
	ms.mutex4Immutable.Unlock()
	return nil
}

// flushMetricsTo writes metric-blocks to the writer.
func (ms *metricStore) flushMetricsTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) error {
	var err error
	// pick immutable
	ms.mutex4Immutable.Lock()
	// build immutable metric-blocks
	for _, tagIdx := range ms.immutable {
		if err = tagIdx.flushMetricTo(tableFlusher, flushCtx); err != nil {
			ms.mutex4Immutable.Unlock()
			return err
		}
	}
	// reset the immutable part
	ms.immutable = nil
	ms.mutex4Immutable.Unlock()

	// reset the mutable part
	ms.mutex4Mutable.RLock()
	err = ms.mutable.flushMetricTo(tableFlusher, flushCtx)
	ms.mutex4Mutable.RUnlock()
	return err
}

// findSeriesIDsByExpr finds series ids by tag filter expr and timeRange
func (ms *metricStore) findSeriesIDsByExpr(expr stmt.TagFilter,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {

	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()

	ms.mutex4Immutable.RLock()
	for _, tagIdx := range ms.immutable {
		if bitMap := tagIdx.findSeriesIDsByExpr(expr, timeRange); bitMap != nil {
			multiVerSeriesIDSet.Add(tagIdx.getVersion(), bitMap)
		}
	}
	ms.mutex4Immutable.RUnlock()

	ms.mutex4Mutable.RLock()
	if bitMap := ms.mutable.findSeriesIDsByExpr(expr, timeRange); bitMap != nil {
		multiVerSeriesIDSet.Add(ms.mutable.getVersion(), bitMap)
	}
	ms.mutex4Mutable.RUnlock()

	return multiVerSeriesIDSet, nil
}

// getSeriesIDsForTag get series ids by tagKey and timeRange
func (ms *metricStore) getSeriesIDsForTag(tagKey string,
	timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error) {

	multiVerSeriesIDSet := series.NewMultiVerSeriesIDSet()

	ms.mutex4Immutable.RLock()
	for _, tagIdx := range ms.immutable {
		if bitMap := tagIdx.getSeriesIDsForTag(tagKey, timeRange); bitMap != nil {
			multiVerSeriesIDSet.Add(tagIdx.getVersion(), bitMap)
		}
	}
	ms.mutex4Immutable.RUnlock()

	ms.mutex4Mutable.RLock()
	if bitMap := ms.mutable.getSeriesIDsForTag(tagKey, timeRange); bitMap != nil {
		multiVerSeriesIDSet.Add(ms.mutable.getVersion(), bitMap)
	}
	ms.mutex4Mutable.RUnlock()

	return multiVerSeriesIDSet, nil
}

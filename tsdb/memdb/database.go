package memdb

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"

	"github.com/cespare/xxhash"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"
)

var memDBLogger = logger.GetLogger("tsdb", "MemDB")

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// WithMaxTagsLimit spawn a goroutine to receives limitation from this channel
	// The producer shall send the config periodically
	// key: metric-name, value: max-limit
	WithMaxTagsLimit(<-chan map[string]uint32)
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(metric *pb.Metric) error
	// ResetMetricStore reassigns a new version to metricStore
	// This method provides the ability to reset the tsStore in memory for skipping the tsID-limitation
	ResetMetricStore(metricName string) error
	// CountMetrics returns the metrics-count of the memory-database
	CountMetrics() int
	// CountTags returns the tags-count of the metricName, return -1 if not exist
	CountTags(metricName string) int
	// Families returns the families in memory which has not been flushed yet
	Families() []int64
	// FlushInvertedIndexTo flushes the inverted-index of series to the kv builder
	FlushInvertedIndexTo(flusher invertedindex.Flusher) error
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int
	// series.Filter contains the methods for filtering seriesIDs from memDB
	series.Filter
	// series.MetaGetter returns tag values by tag keys and spec version for metric level
	series.MetaGetter
	// series.Suggester returns the suggestions from prefix string
	series.MetricMetaSuggester
	series.TagValueSuggester
	// flow.DataFilter filters the data based on condition
	flow.DataFilter
	// series.Storage returns the high level function of storage
	series.Storage
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	TimeWindow int
	Interval   timeutil.Interval
	Generator  metadb.IDGenerator
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	timeWindow          int                // rollup window of memory-database
	interval            timeutil.Interval  // time interval of rollup
	blockStore          *blockStore        // reusable pool
	ctx                 context.Context    // used for exiting goroutines
	evictNotifier       chan struct{}      // notifying evictor to evict
	once4Syncer         sync.Once          // once for tags-limitation syncer
	metricHash2ID       sync.Map           // key: FNV64a(metric-name), value: metric global unique id(metric-id)
	mStores             *metricBucket      // metric-id -> *metricStore
	generator           metadb.IDGenerator // the generator for generating ID of metric, field
	size                atomic.Int32       // memdb's size
	lastWroteFamilyTime atomic.Int64       // prevents familyTime inserting repeatedly
	familyTimes         sync.Map           // familyTime(int64) -> struct{}

	lock sync.Mutex //lock of create metric store
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(ctx context.Context, cfg MemoryDatabaseCfg) MemoryDatabase {
	md := memoryDatabase{
		timeWindow:          cfg.TimeWindow,
		interval:            cfg.Interval,
		generator:           cfg.Generator,
		blockStore:          newBlockStore(cfg.TimeWindow),
		mStores:             newMetricBucket(),
		ctx:                 ctx,
		evictNotifier:       make(chan struct{}),
		size:                *atomic.NewInt32(0),
		lastWroteFamilyTime: *atomic.NewInt64(0),
	}
	go md.evictor(ctx)
	return &md
}

// getMStore returns the mStore by metric-name.
func (md *memoryDatabase) getMStore(metricName string) (mStore mStoreINTF, ok bool) {
	metricIDINTF, ok := md.metricHash2ID.Load(xxhash.Sum64String(metricName))
	if !ok {
		return nil, ok
	}
	metricID := metricIDINTF.(uint32)
	return md.mStores.get(metricID)
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricName string, hash uint64) (metricID uint32, mStore mStoreINTF) {
	metricIDINTF, ok := md.metricHash2ID.Load(hash)
	if !ok {
		// gen new metric id
		metricID = md.generator.GenMetricID(metricName)
		md.metricHash2ID.Store(hash, metricID)
	} else {
		metricID = metricIDINTF.(uint32)
	}

	mStore, ok = md.mStores.get(metricID)
	if !ok {
		// not found need create new metric store
		md.lock.Lock()
		// double check mStore if exist
		mStore, ok = md.mStores.get(metricID)
		if !ok {
			mStore = newMetricStore()
			md.size.Add(int32(mStore.MemSize()))
			md.mStores.put(metricID, mStore)
		}
		md.lock.Unlock()
	}
	// found metric store in current memory database
	return
}

// WithMaxTagsLimit syncs the limitation for different metrics.
func (md *memoryDatabase) WithMaxTagsLimit(limitationCh <-chan map[string]uint32) {
	md.once4Syncer.Do(func() {
		go func() {
			for {
				select {
				case <-md.ctx.Done():
					return
				case limitations, ok := <-limitationCh:
					if !ok {
						return
					}
					if limitations == nil {
						continue
					}
					md.setLimitations(limitations)
				}
			}
		}()
	})
}

// setLimitations set max-count limitation of tagID.
func (md *memoryDatabase) setLimitations(limitations map[string]uint32) {
	for metricName, limit := range limitations {
		mStore, ok := md.getMStore(metricName)
		if !ok {
			continue
		}
		mStore.SetMaxTagsLimit(limit)
	}
}

// writeContext holds the context for writing
type writeContext struct {
	blockStore   *blockStore
	generator    metadb.IDGenerator
	metricID     uint32
	familyTime   int64
	slotIndex    int
	timeInterval int64
	mStoreFieldIDGetter
}

// PointTime returns the point time
func (writeCtx writeContext) PointTime() int64 {
	return writeCtx.familyTime + writeCtx.timeInterval*int64(writeCtx.slotIndex)
}

func (md *memoryDatabase) addFamilyTime(familyTime int64) {
	if md.lastWroteFamilyTime.Swap(familyTime) == familyTime {
		return
	}
	md.familyTimes.Store(familyTime, struct{}{})
}

// Write writes metric-point to database.
func (md *memoryDatabase) Write(metric *pb.Metric) error {
	timestamp := metric.Timestamp
	// calculate family start time and slot index
	intervalCalc := md.interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)                         // day
	family := intervalCalc.CalcFamily(timestamp, segmentTime)                      // hours
	familyTime := intervalCalc.CalcFamilyStartTime(segmentTime, family)            // family timestamp
	slotIndex := intervalCalc.CalcSlot(timestamp, familyTime, md.interval.Int64()) // slot offset of family

	hash := xxhash.Sum64String(metric.Name)
	metricID, mStore := md.getOrCreateMStore(metric.Name, hash)

	writtenSize, err := mStore.Write(metric, writeContext{
		metricID:            metricID,
		blockStore:          md.blockStore,
		generator:           md.generator,
		familyTime:          familyTime,
		slotIndex:           slotIndex,
		timeInterval:        md.interval.Int64(),
		mStoreFieldIDGetter: mStore})
	if err == nil {
		md.addFamilyTime(familyTime)
	}
	md.size.Add(int32(writtenSize))
	return err
}

// evictor do evict periodically.
func (md *memoryDatabase) evictor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-md.evictNotifier:
			md.evict(md.mStores.getAllMetricIDs())
			//FIXME need evict metricHash2ID
			//md.metricID2Hash.Delete(mStore.GetMetricID())
		}
	}
}

// evict evicts tsStore of mStore concurrently,
// and delete metricStore whose timeSeriesMap is empty.
func (md *memoryDatabase) evict(metricIDs *roaring.Bitmap) {
	it := metricIDs.Iterator()
	for it.HasNext() {
		metricID := it.Next()
		mStore, ok := md.mStores.get(metricID)
		if !ok {
			continue
		}
		// delete tag of tStore which has not been used for a while
		evictedSize := mStore.Evict()
		// reduce evicted size
		md.size.Sub(int32(evictedSize))
		// delete mStore whose tags is empty now.
		if mStore.IsEmpty() {
			_ = md.mStores.delete(metricID)
			// reduce empty mStore size
			md.size.Sub(int32(mStore.MemSize()))
		}
	}
}

// ResetMetricStore assigns a new version to the specified metric.
func (md *memoryDatabase) ResetMetricStore(metricName string) error {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return fmt.Errorf("metric: %s doesn't exist", metricName)
	}
	createdSize, err := mStore.ResetVersion()
	md.size.Add(int32(createdSize))
	return err
}

// CountMetrics returns count of metrics in all buckets.
func (md *memoryDatabase) CountMetrics() int {
	return md.mStores.size()
}

// CountTags returns count of tags of a specified metricName, return -1 when metric not exist.
func (md *memoryDatabase) CountTags(metricName string) int {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return -1
	}
	return mStore.GetTagsUsed()
}

// Families returns the families in memory which has not been flushed yet.
func (md *memoryDatabase) Families() []int64 {
	var families []int64
	md.familyTimes.Range(func(key, value interface{}) bool {
		familyTime := key.(int64)
		families = append(families, familyTime)
		return true
	})
	sort.Slice(families, func(i, j int) bool {
		return families[i] < families[j]
	})
	return families
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID     uint32
	familyTime   int64
	timeInterval int64
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder,
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error {
	defer func() {
		// non-block notifying evictor
		select {
		case md.evictNotifier <- struct{}{}:
		default:
			memDBLogger.Warn("flusher is working, concurrently flushing is not allowed")
		}
	}()

	md.familyTimes.Delete(familyTime)
	md.lastWroteFamilyTime.Store(0)

	metricIDs := md.mStores.getAllMetricIDs()
	it := metricIDs.Iterator()
	for it.HasNext() {
		metricID := it.Next()
		mStore, ok := md.mStores.get(metricID)
		if ok {
			flushedSize, err := mStore.FlushMetricsDataTo(flusher, flushContext{
				metricID:     metricID,
				familyTime:   familyTime,
				timeInterval: md.interval.Int64(),
			})
			md.size.Sub(int32(flushedSize))
			if err != nil {
				return err
			}
		}
	}
	//FIXME stone1100 remove it, and test family.deleteObsoleteFiles
	return flusher.Commit()
}

// FlushInvertedIndexTo flushes the series data to a inverted-index file.
func (md *memoryDatabase) FlushInvertedIndexTo(flusher invertedindex.Flusher) (err error) {
	metricIDs := md.mStores.getAllMetricIDs()
	it := metricIDs.Iterator()
	for it.HasNext() {
		metricID := it.Next()
		mStore, ok := md.mStores.get(metricID)
		if ok {
			if err = mStore.FlushInvertedIndexTo(metricID, flusher, md.generator); err != nil {
				return
			}
		}
	}
	return flusher.Commit()
}

// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id from mStore.
func (md *memoryDatabase) FindSeriesIDsByExpr(
	metricID uint32,
	expr stmt.TagFilter,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	mStore, ok := md.mStores.get(metricID)
	if !ok {
		return nil, series.ErrNotFound
	}
	return mStore.FindSeriesIDsByExpr(expr)
}

// GetSeriesIDsForTag get series ids for spec metric's tag key from mStore.
func (md *memoryDatabase) GetSeriesIDsForTag(
	metricID uint32,
	tagKey string,
	timeRange timeutil.TimeRange,
) (
	*series.MultiVerSeriesIDSet,
	error,
) {
	mStore, ok := md.mStores.get(metricID)
	if !ok {
		return nil, series.ErrNotFound
	}
	return mStore.GetSeriesIDsForTag(tagKey)
}

// GetGroupingContext returns the context of group by from memory database
func (md *memoryDatabase) GetGroupingContext(metricID uint32, tagKeys []string,
	version series.Version,
) (series.GroupingContext, error) {
	mStore, ok := md.mStores.get(metricID)
	if !ok {
		return nil, series.ErrNotFound
	}
	return mStore.GetGroupingContext(tagKeys, version)
}

// SuggestMetrics returns nil, as the index-db contains all metricNames
func (md *memoryDatabase) SuggestMetrics(prefix string, limit int) (suggestions []string) {
	return nil
}

// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
func (md *memoryDatabase) SuggestTagKeys(metricName, tagKeyPrefix string, limit int) []string {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return nil
	}
	return mStore.SuggestTagKeys(tagKeyPrefix, limit)
}

// SuggestTagValues returns suggestions from given metricName, tagKey and prefix of tagValue
func (md *memoryDatabase) SuggestTagValues(metricName, tagKey, tagValuePrefix string, limit int) []string {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return nil
	}
	return mStore.SuggestTagValues(tagKey, tagValuePrefix, limit)
}

// Filter filters the data based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(metricID uint32, fieldIDs []uint16,
	version series.Version, seriesIDs *roaring.Bitmap,
) []flow.FilterResultSet {
	mStore, ok := md.mStores.get(metricID)
	if !ok {
		return nil
	}
	return mStore.Filter(metricID, fieldIDs, version, seriesIDs)
}

// Interval return the interval of memory database
func (md *memoryDatabase) Interval() int64 {
	return md.interval.Int64()
}

func (md *memoryDatabase) MemSize() int {
	return int(md.size.Load())
}

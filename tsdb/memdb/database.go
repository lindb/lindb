package memdb

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/eleme/lindb/kv"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/hashers"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

var memDBLogger = logger.GetLogger("memdb")

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// WithMaxTagsLimit spawn a goroutine to receives limitation from this channel
	// The producer shall send the config periodically
	// key: metric-name, value: max-limit
	WithMaxTagsLimit(<-chan map[string]uint32)
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(point models.Point) error
	// ResetMetricStore reassigns a new version to metricStore
	// This method provides the ability to reset the tsStore in memory for skipping the tsID-limitation
	ResetMetricStore(metricName string) error
	// CountMetrics returns the metrics-count of the memory-database
	CountMetrics() int
	// CountTags returns the tags-count of the metricName, return -1 if not exist
	CountTags(metricName string) int
	// Families returns the families in memory which has not been flushed yet
	Families() []int64
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher kv.Flusher, familyTime int64) error
	// todo: @codingcrush, query
}

// mStoresBucket is a simple rwMutex locked map of metricStore.
type mStoresBucket struct {
	rwLock sync.RWMutex
	m      map[string]*metricStore // key: metric-name
}

// allMetricStores returns a clone of metricNames and pointer of mstores in bucket.
func (bkt *mStoresBucket) allMetricStores() (metricNames *[]string, stores *[]*metricStore, release func()) {
	// get all mStores
	length := len(bkt.m)
	metricNames = stringListPool.get(length)
	stores = metricStoresListPool.get(length)

	release = func() {
		stringListPool.put(metricNames)
		metricStoresListPool.put(stores)
	}

	bkt.rwLock.RLock()
	idx := 0
	for metricName, mStore := range bkt.m {
		// delete tag of tStore which has not been used for a while
		(*metricNames)[idx] = metricName
		(*stores)[idx] = mStore
		idx++
	}
	bkt.rwLock.RUnlock()
	return metricNames, stores, release
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	timeWindow    int                                    // rollup window of memory-database
	interval      int64                                  // time interval of rollup
	intervalType  interval.Type                          // month, day, hour
	intervalCalc  interval.Calculator                    // helper function for calculating interval
	blockStore    *blockStore                            // reusable pool
	ctx           context.Context                        // used for exiting goroutines
	evictNotifier chan struct{}                          // notifying evictor to evict
	once4Syncer   sync.Once                              // once for tags-limitation syncer
	mStoresList   [shardingCountOfMStores]*mStoresBucket // metric-name -> *metricStore
	generator     index.IDGenerator                      // the generator for generating ID of metric, field
}

// NewMemoryDatabase returns a new memoryDatabase.
func NewMemoryDatabase(ctx context.Context, timeWindow int,
	interval int64, intervalType interval.Type) (MemoryDatabase, error) {
	return newMemoryDatabase(ctx, timeWindow, interval, intervalType)
}

// newMemoryDatabase is the new method.
func newMemoryDatabase(ctx context.Context, timeWindow int,
	intervalValue int64, intervalType interval.Type) (*memoryDatabase, error) {
	timeCalc, err := interval.GetCalculator(intervalType)
	if err != nil {
		return nil, err
	}
	md := memoryDatabase{
		timeWindow:    timeWindow,
		interval:      intervalValue,
		intervalType:  intervalType,
		intervalCalc:  timeCalc,
		blockStore:    newBlockStore(timeWindow),
		ctx:           ctx,
		evictNotifier: make(chan struct{})}
	for i := range md.mStoresList {
		md.mStoresList[i] = &mStoresBucket{
			m: make(map[string]*metricStore)}
	}
	go md.evictor(ctx)
	// todo: go md.IDSyncer(), initialize it by calling NewMemoryDatabase?
	return &md, nil
}

// getBucket returns the mStoresBucket by metric-name.
func (md *memoryDatabase) getBucket(metricName string) *mStoresBucket {
	return md.mStoresList[shardingCountMask&hashers.Fnv32a(metricName)]
}

// getMStore returns the mStore by metric-name.
func (md *memoryDatabase) getMStore(metricName string) (mStore *metricStore, ok bool) {
	bkt := md.getBucket(metricName)
	bkt.rwLock.RLock()
	mStore, ok = bkt.m[metricName]
	bkt.rwLock.RUnlock()
	return
}

// getOrCreateMStore returns a TimeSeriesStore by metric + tags.
func (md *memoryDatabase) getOrCreateMStore(metricName string) *metricStore {
	bucket := md.getBucket(metricName)
	var mStore *metricStore
	mStore, ok := md.getMStore(metricName)
	if !ok {
		bucket.rwLock.Lock()
		mStore, ok = bucket.m[metricName]
		if !ok {
			mStore = newMetricStore()
			bucket.m[metricName] = mStore
		}
		bucket.rwLock.Unlock()
	}
	return mStore
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
		mStore.setMaxTagsLimit(limit)
	}
}

// Write writes metric-point to database.
func (md *memoryDatabase) Write(point models.Point) error {
	if point == nil {
		return fmt.Errorf("point is nil")
	}
	if point.Fields() == nil {
		return fmt.Errorf("fields is nil")
	}

	mStore := md.getOrCreateMStore(point.Name())
	if mStore.isFull() {
		return models.ErrTooManyTags
	}
	timestamp := point.Timestamp()

	// calculate family start time and slot index
	segmentTime := md.intervalCalc.CalSegmentTime(timestamp)                      // day
	family := md.intervalCalc.CalFamily(timestamp, segmentTime)                   // hours
	familyStartTime := md.intervalCalc.CalFamilyStartTime(segmentTime, family)    // family timestamp
	slotIndex := md.intervalCalc.CalSlot(timestamp, familyStartTime, md.interval) // slot offset of family
	tsStore := mStore.getOrCreateTSStore(point.Tags())
	if tsStore.isFull() {
		return models.ErrTooManyFields
	}

	for fieldName, f := range point.Fields() {
		fieldStore, err := tsStore.getOrCreateFStore(fieldName, f.Type())
		// field type do not match before
		if err != nil {
			return err
		}
		// write data
		fieldStore.write(md.blockStore, familyStartTime, slotIndex, f)
	}
	mStore.addFamilyTime(familyStartTime)
	return nil
}

// evictor do evict periodically.
func (md *memoryDatabase) evictor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-md.evictNotifier:
			for i := 0; i < shardingCountOfMStores; i++ {
				md.evict(md.mStoresList[i&shardingCountMask])
			}
		}
	}
}

// evict evicts tsStore of mStore concurrently,
// and delete metricStore whose timeSeriesMap is empty.
func (md *memoryDatabase) evict(bucket *mStoresBucket) {
	// get all allMStores
	metricNames, allMStores, release := bucket.allMetricStores()
	defer release()

	for idx, mStore := range *allMStores {
		// delete tag of tStore which has not been used for a while
		mStore.evict()
		// delete mStore whose tags is empty now.
		if mStore.isEmpty() {
			bucket.rwLock.Lock()
			if mStore.isEmpty() {
				delete(bucket.m, (*metricNames)[idx])
			}
			bucket.rwLock.Unlock()
		}
	}
}

// ResetMetricStore flushes the specified metricStore, then a new version will be assigned.
func (md *memoryDatabase) ResetMetricStore(metricName string) error {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return fmt.Errorf("metric: %s doesn't exist", metricName)
	}
	return mStore.assignNewVersion()
}

// CountMetrics returns count of metrics in all buckets.
func (md *memoryDatabase) CountMetrics() int {
	var counter = 0
	for bucketIndex := 0; bucketIndex < shardingCountOfMStores; bucketIndex++ {
		md.mStoresList[bucketIndex].rwLock.RLock()
		counter += len(md.mStoresList[bucketIndex].m)
		md.mStoresList[bucketIndex].rwLock.RUnlock()
	}
	return counter
}

// CountTags returns count of tags of a specified metricName, return -1 when metric not exist.
func (md *memoryDatabase) CountTags(metricName string) int {
	mStore, ok := md.getMStore(metricName)
	if !ok {
		return -1
	}
	return mStore.getTagsCount()
}

// Families returns the families in memory which has not been flushed yet.
func (md *memoryDatabase) Families() []int64 {
	families := make(map[int64]struct{})
	for bucketIndex := 0; bucketIndex < shardingCountOfMStores; bucketIndex++ {
		bkt := md.mStoresList[bucketIndex]
		bkt.rwLock.RLock()
		for _, mStore := range bkt.m {
			mStore.unionFamilyTimesTo(families)
		}
		bkt.rwLock.RUnlock()
	}
	var list []int64
	for familyTime := range families {
		list = append(list, familyTime)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return list
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder,
// this method must be called before the cancellation.
func (md *memoryDatabase) FlushFamilyTo(flusher kv.Flusher, familyTime int64) error {
	writer := metrictbl.NewTableWriter(flusher, md.interval)
	return md.flushFamilyTo(writer, familyTime)
}

// flushFamilyTo is the real flush method, used for mock-test
func (md *memoryDatabase) flushFamilyTo(writer metrictbl.TableWriter, familyTime int64) error {
	defer func() {
		// non-block notifying evictor
		select {
		case md.evictNotifier <- struct{}{}:
		default:
		}
	}()

	var err error
	for bucketIndex := 0; bucketIndex < shardingCountOfMStores; bucketIndex++ {
		metricNames, allMetricStores, release := md.mStoresList[bucketIndex].allMetricStores()
		for idx, mStore := range *allMetricStores {
			metricName := (*metricNames)[idx]
			if err = mStore.flushMetricBlocksTo(writer, familyTime, md.generator, metricName); err != nil {
				return err
			}
		}
		// put back to pool
		release()
	}
	return nil
}

// IDSyncer updates the metricID, tsID and fieldID periodically.
func (md *memoryDatabase) IDSyncer(ctx context.Context, syncInterval time.Duration) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			md.syncID()
		}
	}
}

// syncID is the real syncID method.
func (md *memoryDatabase) syncID() {
	for bucketIndex := 0; bucketIndex < shardingCountOfMStores; bucketIndex++ {
		metricNames, metricStores, release := md.mStoresList[bucketIndex].allMetricStores()
		for idx, mStore := range *metricStores {
			metricName := (*metricNames)[idx]
			metricID := mStore.mustGetMetricID(md.generator, metricName)
			mStore.updateTSIDAndFieldID(md.generator, metricID)
		}
		release()
	}
}

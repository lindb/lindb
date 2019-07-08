package memdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	radix "github.com/armon/go-radix"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/hashers"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/timeutil"
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
// A version is assigned since start.
type MemoryDatabase interface {
	// PrefixSearchMetricNames returns metric names by prefix search.
	PrefixSearchMetricNames(prefix string, maxCount uint16) []string
	// RegexSearchTags returns tag-values which matches the pattern.
	RegexSearchTags(metric string, tagsIDPattern string) []string
	// WithMaxTagsLimit spawn a goroutine to receives limitation from this channel.
	// The producer shall send the config periodically.
	// key: metric-name, value: max-limit
	WithMaxTagsLimit(<-chan map[string]uint32)
	// Write writes metrics to the memory-database,
	// return error when exceed max count of tagsIdentifier
	Write(point models.Point) error
	// GetVersion returns the version of memory-database,
	// which is actually the uptime in millisecond.
	GetVersion() int64
}

// mStoresBucket is a simple rwMutex locked map of metricStore.
type mStoresBucket struct {
	rwLock sync.RWMutex
	m      map[string]*metricStore
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	timeWindow   int
	interval     int64
	intervalType interval.Type
	intervalCalc interval.Calculator

	blockStore *blockStore

	ctx         context.Context
	once4Syncer sync.Once                              // once for tags-limitation syncer
	mStoresList [shardingCountOfMStores]*mStoresBucket // metric-name -> *metricStore
	mu4Tree     sync.RWMutex                           // rwMutex of radix.Tree
	tree        *radix.Tree                            // radix tree for prefix string searching
	version     int64                                  // start-time in seconds
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
		timeWindow:   timeWindow,
		interval:     intervalValue,
		intervalType: intervalType,
		intervalCalc: timeCalc,
		blockStore:   newBlockStore(timeWindow),
		ctx:          ctx,
		tree:         radix.New(),
		version:      timeutil.Now()}
	for i := range md.mStoresList {
		md.mStoresList[i] = &mStoresBucket{
			m: make(map[string]*metricStore)}
	}
	go md.evictorRunner(ctx)
	return &md, nil
}

// getBucket returns the mStoresBucket by metric-name.
func (md *memoryDatabase) getBucket(metric string) *mStoresBucket {
	return md.mStoresList[shardingCountMask&hashers.Fnv32a(metric)]
}

// getMetricStore returns a TimeSeriesStore by metric + tags.
func (md *memoryDatabase) getMetricStore(metricName string) *metricStore {
	bucket := md.getBucket(metricName)

	var mStore *metricStore
	bucket.rwLock.RLock()
	mStore, ok := bucket.m[metricName]
	bucket.rwLock.RUnlock()

	if !ok {
		bucket.rwLock.Lock()
		mStore, ok = bucket.m[metricName]
		if !ok {
			mStore = newMetricStore(metricName)
			md.mu4Tree.Lock()
			md.tree.Insert(metricName, nil)
			md.mu4Tree.Unlock()
			bucket.m[metricName] = mStore
		}
		bucket.rwLock.Unlock()
	}
	return mStore
}

// PrefixSearchMetricNames returns the metric names by prefix search.
func (md *memoryDatabase) PrefixSearchMetricNames(prefix string, maxCount uint16) []string {
	if prefix == "" {
		return nil
	}
	md.mu4Tree.RLock()
	defer md.mu4Tree.RUnlock()

	var metricNames []string
	md.tree.WalkPrefix(prefix, func(s string, v interface{}) bool {
		if len(metricNames) >= int(maxCount) {
			return true
		}
		metricNames = append(metricNames, s)
		return false
	})
	return metricNames
}

// RegexSearchTags returns tag-values which matches the pattern.
func (md *memoryDatabase) RegexSearchTags(metric string, tagsIDPattern string) []string {
	bucket := md.getBucket(metric)

	bucket.rwLock.RLock()
	mStore, ok := bucket.m[metric]
	bucket.rwLock.RUnlock()
	if !ok {
		return nil
	}
	return mStore.regexSearchTags(tagsIDPattern)
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
		bucket := md.getBucket(metricName)
		bucket.rwLock.RLock()
		mStore, ok := bucket.m[metricName]
		bucket.rwLock.RUnlock()
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

	mStore := md.getMetricStore(point.Name())
	if mStore.isFull() {
		return models.ErrTooManyTags
	}
	timestamp := point.Timestamp()

	// calc family start time and slot index
	segmentTime := md.intervalCalc.CalSegmentTime(timestamp)
	family := md.intervalCalc.CalFamily(timestamp, segmentTime)
	familyStartTime := md.intervalCalc.CalFamilyStartTime(segmentTime, family)
	slotIndex := md.intervalCalc.CalSlot(timestamp, familyStartTime, md.interval)

	tsStore := mStore.getTimeSeries(point.TagsID())
	for fieldName, field := range point.Fields() {
		fieldStore := tsStore.getFieldStore(fieldName)
		// write data
		fieldStore.write(md.blockStore, familyStartTime, slotIndex, field)
	}
	return nil
}

// evictorRunner runs evictor periodically.
func (md *memoryDatabase) evictorRunner(ctx context.Context) {
	// delay starting evictor in evictingDuration.
	delayRandom := time.NewTimer(time.Duration(getTagsIDTTL()) * time.Millisecond)
	select {
	case <-ctx.Done():
		delayRandom.Stop()
		return
	case <-delayRandom.C:
		delayRandom.Stop()
	}

	evictorTicker := time.NewTicker(time.Duration(getEvictInterval()) * time.Millisecond)
	var counter = 0
	defer evictorTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-evictorTicker.C:
			theStoreMap := md.mStoresList[counter&shardingCountMask]
			md.evict(theStoreMap)
			counter++
		}
	}
}

// evict evicts tsStore of mStore concurrently,
// and delete metricStore whose timeSeriesMap is empty.
func (md *memoryDatabase) evict(theStoreMap *mStoresBucket) {
	// get all mStores
	var allMstore []*metricStore
	theStoreMap.rwLock.RLock()
	for _, mStore := range theStoreMap.m {
		allMstore = append(allMstore, mStore)
	}
	theStoreMap.rwLock.RUnlock()

	for _, mStore := range allMstore {
		// delete tag of tStore which has not been used for a while
		mStore.evict()
		// delete mStore whose tags is empty now.
		if mStore.isEmpty() {
			theStoreMap.rwLock.Lock()
			if mStore.isEmpty() {
				delete(theStoreMap.m, mStore.name)
				md.mu4Tree.Lock()
				md.tree.Delete(mStore.name)
				md.mu4Tree.Unlock()
			}
			theStoreMap.rwLock.Unlock()
		}
	}
}

// GetVersion returns the version.
func (md *memoryDatabase) GetVersion() int64 {
	return md.version
}

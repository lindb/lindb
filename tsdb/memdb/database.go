package memdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eleme/lindb/pkg/hashers"

	radix "github.com/armon/go-radix"

	"github.com/eleme/lindb/models"
)

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
// A version is assigned since start.
type MemoryDatabase interface {
	// PrefixSearchMeasurements returns measurements by prefix search.
	PrefixSearchMeasurements(measurementPrefix string, maxCount uint16) []string
	// RegexSearchTags returns tag-values which matches the pattern.
	RegexSearchTags(measurement string, tagsIDPattern string) []string
	// WithMaxTagsLimit spawn a goroutine to receives limitation from this channel.
	// The producer shall send the config periodically.
	// key: measurement-name, value: max-limit
	WithMaxTagsLimit(<-chan map[string]uint32)
	// Write writes metrics to the memory-database,
	// return error when exceed max count of tagsIdentifier
	Write(point models.Point, segmentTime int64, slotTime int32) error
	// GetVersion returns the version of memory-database,
	// which is actually the uptime in millisecond.
	GetVersion() int64
}

// mStoresBucket is a simple rwMutex locked map of MeasurementStore.
type mStoresBucket struct {
	rwLock sync.RWMutex
	m      map[string]*measurementStore
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	ctx         context.Context
	once4Syncer sync.Once                              // once for tags-limitation syncer
	mStoresList [shardingCountOfMStores]*mStoresBucket // measurement-name -> *MeasurementStore
	mu4Tree     sync.RWMutex                           // rwMutex of radix.Tree
	tree        *radix.Tree                            // radix tree for prefix string searching
	version     int64                                  // start-time in seconds
}

// NewMemoryDatabase returns a new memoryDatabase.
func NewMemoryDatabase(ctx context.Context) MemoryDatabase {
	return newMemoryDatabase(ctx)
}

// newMemoryDatabase is the new method.
func newMemoryDatabase(ctx context.Context) *memoryDatabase {
	md := memoryDatabase{
		ctx:     ctx,
		tree:    radix.New(),
		version: time.Now().Unix()}
	for i := range md.mStoresList {
		md.mStoresList[i] = &mStoresBucket{
			m: make(map[string]*measurementStore)}
	}
	go md.evictorRunner(ctx)
	return &md
}

// getBucket returns the mStoresBucket by measurement-name.
func (md *memoryDatabase) getBucket(measurement string) *mStoresBucket {
	return md.mStoresList[shardingCountMask&hashers.Fnv32a(measurement)]
}

// getMeasurementStore returns a TimeSeriesStore by measurement + tags.
func (md *memoryDatabase) getMeasurementStore(measurement string) *measurementStore {
	bucket := md.getBucket(measurement)

	var mStore *measurementStore
	bucket.rwLock.RLock()
	mStore, ok := bucket.m[measurement]
	bucket.rwLock.RUnlock()

	if !ok {
		bucket.rwLock.Lock()
		mStore, ok = bucket.m[measurement]
		if !ok {
			mStore = newMeasurementStore(measurement)
			md.mu4Tree.Lock()
			md.tree.Insert(measurement, nil)
			md.mu4Tree.Unlock()
			bucket.m[measurement] = mStore
		}
		bucket.rwLock.Unlock()
	}
	return mStore
}

// PrefixSearchMeasurements returns the measurements by prefix search.
func (md *memoryDatabase) PrefixSearchMeasurements(measurementPrefix string, maxCount uint16) []string {
	if measurementPrefix == "" {
		return nil
	}
	md.mu4Tree.RLock()
	defer md.mu4Tree.RUnlock()

	var measurements []string
	md.tree.WalkPrefix(measurementPrefix, func(s string, v interface{}) bool {
		if len(measurements) >= int(maxCount) {
			return true
		}
		measurements = append(measurements, s)
		return false
	})
	return measurements
}

// RegexSearchTags returns tag-values which matches the pattern.
func (md *memoryDatabase) RegexSearchTags(measurement string, tagsIDPattern string) []string {
	bucket := md.getBucket(measurement)

	bucket.rwLock.RLock()
	mStore, ok := bucket.m[measurement]
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
	for measurementName, limit := range limitations {
		bucket := md.getBucket(measurementName)
		bucket.rwLock.RLock()
		mStore, ok := bucket.m[measurementName]
		bucket.rwLock.RUnlock()
		if !ok {
			continue
		}
		mStore.setMaxTagsLimit(limit)
	}
}

// Write writes metric-point to database.
func (md *memoryDatabase) Write(point models.Point, segmentTime int64, slotTime int32) error {
	if point == nil {
		return fmt.Errorf("point is nil")
	}
	if point.Fields() == nil {
		return fmt.Errorf("fields is nil")
	}

	mStore := md.getMeasurementStore(point.Name())
	if mStore.isFull() {
		return models.ErrTooManyTags
	}

	tsStore := mStore.getTimeSeries(point.TagsID())
	var err error
	for fieldName, field := range point.Fields() {
		fieldStore := tsStore.getFieldStore(fieldName)
		segmentStore := fieldStore.getSegmentStore(segmentTime)
		if err = segmentStore.Write(slotTime, field); err != nil {
			return err
		}
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
// and delete MeasurementStore whose timeSeriesMap is empty.
func (md *memoryDatabase) evict(theStoreMap *mStoresBucket) {
	// get all mStores
	var allMstore []*measurementStore
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
				delete(theStoreMap.m, mStore.measurementName)
				md.mu4Tree.Lock()
				md.tree.Delete(mStore.measurementName)
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

package memdb

import (
	"context"
	"sort"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("tsdb", "MemDB")

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(metric *pb.Metric) error
	// Families returns the families in memory which has not been flushed yet
	Families() []int64
	// FlushFamilyTo flushes the corresponded family data to builder.
	// Close is not in the flushing process.
	FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error
	// MemSize returns the memory-size of this metric-store
	MemSize() int32
	// flow.DataFilter filters the data based on condition
	flow.DataFilter
	// series.Storage returns the high level function of storage
	series.Storage
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	TimeWindow uint16
	Interval   timeutil.Interval
	Generator  metadb.IDGenerator
	Index      indexdb.MemoryIndexDatabase
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	timeWindow uint16            // rollup window of memory-database
	interval   timeutil.Interval // time interval of rollup
	generator  metadb.IDGenerator
	index      indexdb.MemoryIndexDatabase // memory index database for assign metric id/series id
	ctx        context.Context

	blockStore  *blockStore   // reusable pool
	mStores     *metricBucket // metric-id -> *metricStore
	size        atomic.Int32  // memory database's size
	familyTimes sync.Map      // familyTime(int64) -> struct{}

	lock sync.Mutex //lock of create metric store
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(ctx context.Context, cfg MemoryDatabaseCfg) MemoryDatabase {
	return &memoryDatabase{
		timeWindow: cfg.TimeWindow,
		interval:   cfg.Interval,
		generator:  cfg.Generator,
		index:      cfg.Index,
		ctx:        ctx,
		blockStore: newBlockStore(cfg.TimeWindow),
		mStores:    newMetricBucket(),
		size:       *atomic.NewInt32(0),
	}
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricID uint32) (mStore mStoreINTF) {
	mStore, ok := md.mStores.get(metricID)

	if !ok {
		// not found need create new metric store
		md.lock.Lock()
		// double check mStore if exist
		mStore, ok = md.mStores.get(metricID)
		if !ok {
			mStore = newMetricStore()
			md.size.Add(emptyMStoreSize)
			md.mStores.put(metricID, mStore)
		}
		md.lock.Unlock()
	}
	// found metric store in current memory database
	return
}

// writeContext holds the context for writing
type writeContext struct {
	blockStore *blockStore
	generator  metadb.IDGenerator
	familyTime int64
	metricID   uint32
	slotIndex  uint16
	mStoreFieldIDGetter
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

	metricID, seriesID := md.index.GetTimeSeriesID(metric.Name, metric.Tags, metric.TagsHash)
	mStore := md.getOrCreateMStore(metricID)

	writtenSize, err := mStore.Write(seriesID, metric.Fields, writeContext{
		metricID:            metricID,
		familyTime:          familyTime,
		blockStore:          md.blockStore,
		generator:           md.generator,
		slotIndex:           uint16(slotIndex), //FIXME
		mStoreFieldIDGetter: mStore})
	if err == nil {
		// set family times
		md.familyTimes.Store(familyTime, struct{}{})
	}
	md.size.Add(int32(writtenSize))
	return err
}

// CountMetrics returns count of metrics in all buckets.
func (md *memoryDatabase) CountMetrics() int {
	return md.mStores.size()
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

	start, end uint16 // start/end time slot, metric level flush context
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder,
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error {
	metricIDs := md.mStores.getAllMetricIDs()
	it := metricIDs.Iterator()
	for it.HasNext() {
		metricID := it.Next()
		mStore, ok := md.mStores.get(metricID)
		if ok {
			err := mStore.FlushMetricsDataTo(flusher, flushContext{
				metricID:     metricID,
				familyTime:   familyTime,
				timeInterval: md.interval.Int64(),
			})
			if err != nil {
				return err
			}
		}
	}
	//FIXME stone1100 remove it, and test family.deleteObsoleteFiles
	return flusher.Commit()
}

// Filter filters the data based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(metricID uint32, fieldIDs []uint16,
	version series.Version, seriesIDs *roaring.Bitmap,
) ([]flow.FilterResultSet, error) {
	mStore, ok := md.mStores.get(metricID)
	if !ok {
		return nil, nil
	}
	return mStore.Filter(metricID, fieldIDs, version, seriesIDs)
}

// Interval return the interval of memory database
func (md *memoryDatabase) Interval() int64 {
	return md.interval.Int64()
}

// MemSize returns the time series database memory size
func (md *memoryDatabase) MemSize() int32 {
	return md.size.Load()
}

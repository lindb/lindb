package memdb

import (
	"io"
	"sort"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("tsdb", "MemDB")

type familyID uint8

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// AcquireWrite acquires writing data points
	AcquireWrite()
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(namespace, metricName string, metricID, seriesID uint32, timestamp int64, fields []*pb.Field) (err error)
	// CompleteWrite completes writing data points
	CompleteWrite()
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
	// io.Closer closes the memory database resource
	io.Closer
}

// MemoryDatabaseCfg represents the memory database config
type MemoryDatabaseCfg struct {
	Interval timeutil.Interval
	Metadata metadb.Metadata
	TempPath string
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID     uint32
	familyID     familyID
	timeInterval int64

	slotRange // start/end time slot, metric level flush context
}

// familyTimeIDEntry keeps the mapping of familyTime and familyID
type familyTimeIDEntry struct {
	time int64
	id   familyID
}

// familyTimeIDEntries implements sort.Interface
type familyTimeIDEntries []familyTimeIDEntry

func (e familyTimeIDEntries) Len() int           { return len(e) }
func (e familyTimeIDEntries) Less(i, j int) bool { return e[i].time < e[j].time }
func (e familyTimeIDEntries) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e familyTimeIDEntries) GetID(time int64) (familyID, bool) {
	idx := sort.Search(e.Len(), func(i int) bool {
		return e[i].time >= time
	})
	if idx >= e.Len() || e[idx].time != time {
		return 0, false
	}
	return e[idx].id, true
}

func (e familyTimeIDEntries) AddID(time int64, id familyID) familyTimeIDEntries {
	newE := append(e, familyTimeIDEntry{id: id, time: time})
	sort.Sort(newE)
	return newE
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	interval timeutil.Interval // time interval of rollup
	metadata metadb.Metadata   // metadata for assign metric id/field id

	mStores *MetricBucketStore // metric id => mStoreINTF
	buf     DataPointBuffer

	allocSize           atomic.Int32        // allocated size
	familyTimeIDEntries familyTimeIDEntries // familyTime(int64) -> family time id
	familyIDSeq         uint8

	writeCondition sync.WaitGroup
	rwMutex        sync.RWMutex // lock of create metric store
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg MemoryDatabaseCfg) (MemoryDatabase, error) {
	buf, err := newDataPointBuffer(cfg.TempPath)
	if err != nil {
		return nil, err
	}
	return &memoryDatabase{
		interval:  cfg.Interval,
		metadata:  cfg.Metadata,
		buf:       buf,
		mStores:   NewMetricBucketStore(),
		allocSize: *atomic.NewInt32(0),
	}, err
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricID uint32) (mStore mStoreINTF) {
	mStore, ok := md.mStores.Get(metricID)
	if !ok {
		// not found need create new metric store
		mStore = newMetricStore()
		md.allocSize.Add(emptyMStoreSize)
		md.mStores.Put(metricID, mStore)
	}
	// found metric store in current memory database
	return
}

// AcquireWrite acquires writing data points
func (md *memoryDatabase) AcquireWrite() {
	md.writeCondition.Add(1)
}

// CompleteWrite completes writing data points
func (md *memoryDatabase) CompleteWrite() {
	md.writeCondition.Done()
}

// Write writes metric-point to database.
func (md *memoryDatabase) Write(
	namespace, metricName string,
	metricID, seriesID uint32,
	timestamp int64,
	fields []*pb.Field,
) (err error) {
	// calculate family start time and slot index
	intervalCalc := md.interval.Calculator()
	familyTime := md.getFamilyTime(timestamp)
	slotIndex := uint16(intervalCalc.CalcSlot(timestamp, familyTime, md.interval.Int64())) // slot offset of family

	md.rwMutex.Lock()
	defer md.rwMutex.Unlock()

	mStore := md.getOrCreateMStore(metricID)
	// assign family id for family time
	fi := md.assignFamilyID(familyTime)
	fID := fi

	tStore, size := mStore.GetOrCreateTStore(seriesID)
	written := false
	for _, f := range fields {
		fieldType := getFieldType(f)
		if fieldType == field.Unknown {
			//FIXME add log or metric
			continue
		}
		fieldID, err := md.metadata.MetadataDatabase().GenFieldID(namespace, metricName, f.Name, fieldType)
		if err != nil {
			//FIXME stone1100 add metric
			continue
		}
		for _, pField := range f.Fields {
			pFieldID := field.PrimitiveID(pField.PrimitiveID)
			pStore, ok := tStore.GetFStore(fID, fieldID, pFieldID)
			if !ok {
				buf, err := md.buf.AllocPage()
				if err != nil {
					return err
				}
				pStore = newFieldStore(buf, fID, fieldID, pFieldID)
				size += emptyPrimitiveFieldStoreSize + 8
				tStore.InsertFStore(pStore)
			}
			size += pStore.Write(fieldType, slotIndex, pField.Value)
		}

		// if write data success, add field into metric level for cache
		mStore.AddField(fieldID, fieldType)
		written = true
	}
	if written {
		mStore.SetTimestamp(fi, slotIndex)
	}
	md.allocSize.Add(int32(size))
	return nil
}

// Families returns the families in memory which has not been flushed yet.
func (md *memoryDatabase) Families() []int64 {
	var families []int64
	for _, entry := range md.familyTimeIDEntries {
		families = append(families, entry.time)
	}
	return families
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder,
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error {
	// waiting current writing complete
	md.writeCondition.Wait()

	familyID, _ := md.familyTimeIDEntries.GetID(familyTime)
	if err := md.mStores.WalkEntry(func(key uint32, value mStoreINTF) error {
		if err := value.FlushMetricsDataTo(flusher, flushContext{
			metricID:     key,
			familyID:     familyID,
			timeInterval: md.interval.Int64(),
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	//FIXME stone1100 remove it, and test family.deleteObsoleteFiles
	return flusher.Commit()
}

// Filter filters the data based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (md *memoryDatabase) Filter(metricID uint32, fieldIDs []field.ID,
	seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
) ([]flow.FilterResultSet, error) {
	// get family tine query range
	familyTimeRange := timeutil.TimeRange{
		Start: md.getFamilyTime(timeRange.Start),
		End:   md.getFamilyTime(timeRange.End),
	}

	md.rwMutex.RLock()
	defer md.rwMutex.RUnlock()

	// find if has match family id based on family time range
	familyIDs := make(map[familyID]int64)
	for _, entry := range md.familyTimeIDEntries {
		if familyTimeRange.Contains(entry.time) {
			familyIDs[entry.id] = entry.time
		}
	}
	if len(familyIDs) == 0 {
		return nil, constants.ErrNotFound
	}

	mStore, ok := md.mStores.Get(metricID)
	if !ok {
		return nil, constants.ErrNotFound
	}
	return mStore.Filter(fieldIDs, seriesIDs, familyIDs)
}

// Interval return the interval of memory database
func (md *memoryDatabase) Interval() int64 {
	return md.interval.Int64()
}

// MemSize returns the time series database memory size
func (md *memoryDatabase) MemSize() int32 {
	return md.allocSize.Load()
}

// Close closes memory data point buffer
func (md *memoryDatabase) Close() error {
	return md.buf.Close()
}

// assignFamily assigns family id for family time
func (md *memoryDatabase) assignFamilyID(familyTime int64) familyID {
	fID, ok := md.familyTimeIDEntries.GetID(familyTime)
	if ok {
		return fID
	}
	fID = familyID(md.familyIDSeq)
	md.familyIDSeq++
	md.familyTimeIDEntries = md.familyTimeIDEntries.AddID(familyTime, fID)
	return fID
}

func (md *memoryDatabase) getFamilyTime(timestamp int64) (familyTime int64) {
	intervalCalc := md.interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)             // day
	family := intervalCalc.CalcFamily(timestamp, segmentTime)          // hours
	familyTime = intervalCalc.CalcFamilyStartTime(segmentTime, family) // family timestamp
	return
}

package memdb

import (
	"sort"
	"sync"

	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./database.go -destination=./database_mock.go -package memdb

var memDBLogger = logger.GetLogger("tsdb", "MemDB")

type familyID uint8

// MemoryDatabase is a database-like concept of Shard as memTable in cassandra.
type MemoryDatabase interface {
	// Write writes metrics to the memory-database,
	// return error on exceeding max count of tagsIdentifier or writing failure
	Write(namespace, metricName string, metricID, seriesID uint32, timestamp int64, fields []*pb.Field) (err error)
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
	Interval timeutil.Interval
	Metadata metadb.Metadata
	Index    indexdb.IndexDatabase
	TempPath string
}

// memoryDatabase implements MemoryDatabase.
type memoryDatabase struct {
	interval timeutil.Interval // time interval of rollup
	metadata metadb.Metadata   // metadata for assign metric id/field id

	mStores *MetricBucketStore // metric id => mStoreINTF
	buf     DataPointBuffer

	size        atomic.Int32    // memory database's size
	familyTimes map[int64]uint8 // familyTime(int64) -> family time id
	familyIDSeq uint8

	rwMutex sync.RWMutex //lock of create metric store
}

// NewMemoryDatabase returns a new MemoryDatabase.
func NewMemoryDatabase(cfg MemoryDatabaseCfg) MemoryDatabase {
	//FIXME check temp path is empty
	buf := newDataPointBuffer(cfg.TempPath)
	return &memoryDatabase{
		interval:    cfg.Interval,
		metadata:    cfg.Metadata,
		buf:         buf,
		mStores:     NewMetricBucketStore(),
		size:        *atomic.NewInt32(0),
		familyTimes: make(map[int64]uint8),
	}
}

// getOrCreateMStore returns the mStore by metricHash.
func (md *memoryDatabase) getOrCreateMStore(metricID uint32) (mStore mStoreINTF) {
	mStore, ok := md.mStores.Get(metricID)
	if !ok {
		// not found need create new metric store
		mStore = newMetricStore()
		md.size.Add(emptyMStoreSize)
		md.mStores.Put(metricID, mStore)
	}
	// found metric store in current memory database
	return
}

// flushContext holds the context for flushing
type flushContext struct {
	metricID     uint32
	familyID     uint8
	timeInterval int64

	start, end uint16 // start/end time slot, metric level flush context
}

// Write writes metric-point to database.
func (md *memoryDatabase) Write(namespace, metricName string, metricID,
	seriesID uint32,
	timestamp int64, fields []*pb.Field,
) (err error) {
	// calculate family start time and slot index
	intervalCalc := md.interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)                                 // day
	family := intervalCalc.CalcFamily(timestamp, segmentTime)                              // hours
	familyTime := intervalCalc.CalcFamilyStartTime(segmentTime, family)                    // family timestamp
	slotIndex := uint16(intervalCalc.CalcSlot(timestamp, familyTime, md.interval.Int64())) // slot offset of family

	md.rwMutex.Lock()
	defer md.rwMutex.Unlock()

	mStore := md.getOrCreateMStore(metricID)
	// assign family id for family time
	fi := md.assignFamilyID(familyTime)
	fID := familyID(fi)

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
		//fStore, writtenSize := tStore.GetFStore(fieldID)
		pStore, ok := tStore.GetFStore(fID, field.ID(fieldID), field.PrimitiveID(1))
		if !ok {
			buf, err := md.buf.AllocPage()
			if err != nil {
				return err
			}
			pStore = newFieldStore(buf, fID, field.ID(fieldID), field.PrimitiveID(1))
			size += emptyPrimitiveFieldStoreSize + 8
			tStore.InsertFStore(pStore)
		}
		value := md.getFieldValue(fieldType, f)
		size += pStore.Write(fieldType, slotIndex, value)

		// if write data success, add field into metric level for cache
		mStore.AddField(fieldID, fieldType)
		written = true
	}
	if written {
		mStore.SetTimestamp(fi, slotIndex)
	}
	md.size.Add(int32(size))
	return nil
}

// Families returns the families in memory which has not been flushed yet.
func (md *memoryDatabase) Families() []int64 {
	var families []int64
	for familyTime := range md.familyTimes {
		families = append(families, familyTime)
	}
	sort.Slice(families, func(i, j int) bool {
		return families[i] < families[j]
	})
	return families
}

// FlushFamilyTo flushes all data related to the family from metric-stores to builder,
func (md *memoryDatabase) FlushFamilyTo(flusher metricsdata.Flusher, familyTime int64) error {
	if err := md.mStores.WalkEntry(func(key uint32, value mStoreINTF) error {
		if err := value.FlushMetricsDataTo(flusher, flushContext{
			metricID:     key,
			familyID:     md.familyTimes[familyTime],
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
func (md *memoryDatabase) Filter(metricID uint32, fieldIDs []uint16,
	version series.Version, seriesIDs *roaring.Bitmap,
) ([]flow.FilterResultSet, error) {
	mStore, ok := md.mStores.Get(metricID)
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

// assignFamily assigns family id for family time
func (md *memoryDatabase) assignFamilyID(familyTime int64) uint8 {
	familyID, ok := md.familyTimes[familyTime]
	if ok {
		return familyID
	}
	familyID = md.familyIDSeq
	md.familyIDSeq++
	md.familyTimes[familyTime] = familyID
	return familyID
}

// getFieldValue returns the field value based on field type
func (md *memoryDatabase) getFieldValue(fieldType field.Type, f *pb.Field) float64 {
	switch fieldType {
	case field.SumField:
		return f.GetSum().Value
	case field.MinField:
		return f.GetMin().Value
	case field.MaxField:
		return f.GetMax().Value
	case field.GaugeField:
		return f.GetGauge().Value
	default:
		return 0
	}
}

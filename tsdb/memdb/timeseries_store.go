package memdb

import (
	"sort"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock_test.go -package memdb

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// GetHash returns the xxhash of the tags
	GetHash() uint64

	// GetFStore returns the fStore in this list from field-id.
	GetFStore(fieldID uint16) (fStoreINTF, bool)

	// Write writes the metric
	Write(metric *pb.Metric, writeCtx writeContext) error

	// FlushSeriesTo flushes the series data segment.
	FlushSeriesTo(
		flusher tblstore.MetricsDataFlusher,
		flushCtx flushContext,
		seriesID uint32,
	) (flushed bool)

	// IsExpired detects if this tStore has not been used for a TTL
	IsExpired() bool

	// IsNoData symbols if all data of this tStore has been flushed
	IsNoData() bool

	// Scan scans time series store, then finds data by field id and time range
	Scan(
		sCtx *series.ScanContext,
		version series.Version,
		seriesID uint32,
		fieldMetas map[uint16]*field.Meta)
}

// fStoreNodes implements sort.Interface
type fStoreNodes []fStoreINTF

func (f fStoreNodes) Len() int           { return len(f) }
func (f fStoreNodes) Less(i, j int) bool { return f[i].GetFieldID() < f[j].GetFieldID() }
func (f fStoreNodes) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	sl            lockers.SpinLock // spin-lock
	hash          uint64           // hash of tags
	lastWroteTime atomic.Uint32    // last Write-time in seconds
	fStoreNodes   fStoreNodes      // key: sorted fStore list by field-name, insert-only
}

// newTimeSeriesStore returns a new tStoreINTF.
func newTimeSeriesStore(tagsHash uint64) tStoreINTF {
	return &timeSeriesStore{
		hash:          tagsHash,
		lastWroteTime: *atomic.NewUint32(uint32(timeutil.Now() / 1000))}
}

// GetHash returns the xxhash of the tags
func (ts *timeSeriesStore) GetHash() uint64 {
	return ts.hash
}

// GetFStore returns the fStore in this list from field-id.
func (ts *timeSeriesStore) GetFStore(fieldID uint16) (fStoreINTF, bool) {
	idx := sort.Search(len(ts.fStoreNodes), func(i int) bool {
		return ts.fStoreNodes[i].GetFieldID() >= fieldID
	})
	if idx >= len(ts.fStoreNodes) || ts.fStoreNodes[idx].GetFieldID() != fieldID {
		return nil, false
	}
	return ts.fStoreNodes[idx], true
}

// insertFStore inserts a new fStore to field list.
func (ts *timeSeriesStore) insertFStore(fStore fStoreINTF) {
	ts.fStoreNodes = append(ts.fStoreNodes, fStore)
	sort.Sort(ts.fStoreNodes)
}

// IsNoData symbols if all data of this tStore has been flushed
func (ts *timeSeriesStore) IsNoData() bool {
	ts.sl.Lock()
	defer ts.sl.Unlock()

	for _, fStore := range ts.fStoreNodes {
		if fStore.SegmentsCount() != 0 {
			return false
		}
	}
	return true
}

// afterFlush checks if the tStore contains any data after flushing
func (ts *timeSeriesStore) afterFlush(flushCtx flushContext) {
	// update hasData flag
	var startTime, endTime int64
	for _, fStore := range ts.fStoreNodes {
		timeRange, ok := fStore.TimeRange(flushCtx.timeInterval)
		if !ok {
			continue
		}
		if startTime == 0 || timeRange.Start < startTime {
			startTime = timeRange.Start
		}
		if endTime == 0 || endTime < timeRange.End {
			endTime = timeRange.End
		}
	}
}

// IsExpired detects if this tStore has not been used for a TTL
func (ts *timeSeriesStore) IsExpired() bool {
	return time.Unix(int64(ts.lastWroteTime.Load()), 0).Add(seriesTTL.Load()).Before(time.Now())
}

// Write Write the data of metric to the fStore.
func (ts *timeSeriesStore) Write(
	metric *pb.Metric,
	writeCtx writeContext,
) error {
	ts.sl.Lock()
	defer ts.sl.Unlock()

	for _, f := range metric.Fields {
		// todo FieldType
		fieldType := getFieldType(f)
		if fieldType == field.Unknown {
			//TODO add log or metric
			continue
		}
		if fStore, err := ts.getOrCreateFStore(f.Name, fieldType, writeCtx); err == nil {
			fStore.Write(f, writeCtx)
			ts.lastWroteTime.Store(uint32(timeutil.Now() / 1000))
		} else {
			return err // field type do not match before, too many fields
		}
	}
	return nil
}

// getOrCreateFStore get a fieldStore by fieldName and fieldType
func (ts *timeSeriesStore) getOrCreateFStore(
	fieldName string,
	fieldType field.Type,
	writeCtx writeContext,
) (
	fStoreINTF,
	error,
) {
	fieldID, err := writeCtx.GetFieldIDOrGenerate(fieldName, fieldType, writeCtx.generator)
	if err != nil {
		return nil, err
	}

	fStore, ok := ts.GetFStore(fieldID)
	if !ok {
		fStore = newFieldStore(fieldID)
		ts.insertFStore(fStore)
	}
	return fStore, nil
}

// FlushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) FlushSeriesTo(
	flusher tblstore.MetricsDataFlusher,
	flushCtx flushContext,
	seriesID uint32,
) (flushed bool) {
	ts.sl.Lock()
	for _, fStore := range ts.fStoreNodes {
		fieldDataFlushed := fStore.FlushFieldTo(flusher, flushCtx.familyTime)
		flushed = flushed || fieldDataFlushed
	}
	if flushed {
		flusher.FlushSeries(seriesID)
		ts.afterFlush(flushCtx)
	}
	// update time range info
	ts.sl.Unlock()
	return
}

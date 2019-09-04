package memdb

import (
	"sort"

	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/tblstore"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock_test.go -package memdb

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// retain retains a lock of tStore
	retain() func()
	// getHash returns the FNV1a hash of the tags
	getHash() uint64
	// getFStore returns the fStore in this list from field-id.
	getFStore(fieldID uint16) (fStoreINTF, bool)
	// write writes the metric
	write(metric *pb.Metric, writeCtx writeContext) error
	// flushSeriesTo flushes the series data segment.
	flushSeriesTo(flusher tblstore.MetricsDataFlusher, flushCtx flushContext, seriesID uint32) (flushed bool)
	// isExpired detects if this tStore has not been used for a TTL
	isExpired() bool
	// isNoData symbols if all data of this tStore has been flushed
	isNoData() bool
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
	lastWroteTime uint32           // last write-time in seconds
	hasData       bool             // noData flags symbols if all data has been flushed
	fStoreNodes   fStoreNodes      // key: sorted fStore list by field-name, insert-only
}

// newTimeSeriesStore returns a new tStoreINTF.
func newTimeSeriesStore(tagsHash uint64) tStoreINTF {
	return &timeSeriesStore{
		hash:          tagsHash,
		lastWroteTime: uint32(timeutil.Now() / 1000)}
}

// retain retains a lock of tStore
func (ts *timeSeriesStore) retain() func() {
	ts.sl.Lock()
	return ts.sl.Unlock
}

// getHash returns the FNV1a hash of the tags
func (ts *timeSeriesStore) getHash() uint64 {
	return ts.hash
}

// getFStore returns the fStore in this list from field-id.
func (ts *timeSeriesStore) getFStore(fieldID uint16) (fStoreINTF, bool) {
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

// isNoData symbols if all data of this tStore has been flushed
func (ts *timeSeriesStore) isNoData() bool {
	return !ts.hasData
}

// afterWrite set hasData to true and updates the lastWroteTime
func (ts *timeSeriesStore) afterWrite() {
	// hasData is true now
	ts.hasData = true
	// update lastWroteTime
	ts.lastWroteTime = uint32(timeutil.Now() / 1000)
}

// afterFlush checks if the tStore contains any data after flushing
func (ts *timeSeriesStore) afterFlush(flushCtx flushContext) {
	// update hasData flag
	var startTime, endTime int64
	ts.hasData = false
	for _, fStore := range ts.fStoreNodes {
		timeRange, ok := fStore.TimeRange(flushCtx.timeInterval)
		if !ok {
			continue
		}
		ts.hasData = true
		if startTime == 0 || timeRange.Start < startTime {
			startTime = timeRange.Start
		}
		if endTime == 0 || endTime < timeRange.End {
			endTime = timeRange.End
		}
	}
}

// isExpired detects if this tStore has not been used for a TTL
func (ts *timeSeriesStore) isExpired() bool {
	return int64(ts.lastWroteTime)*1000+getTagsIDTTL() < timeutil.Now()
}

// write write the data of metric to the fStore.
func (ts *timeSeriesStore) write(metric *pb.Metric, writeCtx writeContext) error {
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
			ts.afterWrite()
		} else {
			return err // field type do not match before, too many fields
		}
	}
	return nil
}

// getOrCreateFStore get a fieldStore by fieldName and fieldType
func (ts *timeSeriesStore) getOrCreateFStore(fieldName string, fieldType field.Type,
	writeCtx writeContext) (fStoreINTF, error) {

	fieldID, err := writeCtx.getFieldIDOrGenerate(fieldName, fieldType, writeCtx.generator)
	if err != nil {
		return nil, err
	}

	fStore, ok := ts.getFStore(fieldID)
	if !ok {
		fStore = newFieldStore(fieldID)
		ts.insertFStore(fStore)
	}
	return fStore, nil
}

// flushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) flushSeriesTo(flusher tblstore.MetricsDataFlusher,
	flushCtx flushContext, seriesID uint32) (flushed bool) {
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

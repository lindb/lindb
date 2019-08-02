package memdb

import (
	"math"
	"sort"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/metrictbl"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock_test.go -package memdb

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// getHash returns the FNV1a hash of the tags
	getHash() uint64
	// write writes the metric
	write(metric *pb.Metric, writeCtx writeContext) error
	// flushSeriesTo flushes the series data segment.
	flushSeriesTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) (flushed bool)
	// timeRange returns the start-time and end-time of segment's data
	// ok means data is available
	timeRange() (timeRange timeutil.TimeRange, ok bool)
	// isExpired detects if this tStore has not been used for a TTL
	isExpired() bool
	// isNoData symbols if all data of this tStore has been flushed
	isNoData() bool
}

// fStoreNodes implements sort.Interface
type fStoreNodes []fStoreINTF

func (f fStoreNodes) Len() int           { return len(f) }
func (f fStoreNodes) Less(i, j int) bool { return f[i].getFieldName() < f[j].getFieldName() }
func (f fStoreNodes) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	seriesID      uint32           // series id
	hash          uint64           // hash of tags
	fStoreNodes   fStoreNodes      // key: sorted fStore list by field-name, insert-only
	lastWroteTime int64            // last write-time in milliseconds
	startDelta    int32            // startTime = lastWroteTime + startDelta
	endDelta      int32            // endTime = lastWroteTime + endDelta
	hasData       bool             // noData flags symbols if all data has been flushed
	sl            lockers.SpinLock // spin-lock
}

// newTimeSeriesStore returns a new tStoreINTF.
func newTimeSeriesStore(seriesID uint32, tagsHash uint64) tStoreINTF {
	return &timeSeriesStore{
		seriesID:      seriesID,
		hash:          tagsHash,
		lastWroteTime: timeutil.Now(),
		startDelta:    math.MaxInt32,
		endDelta:      math.MaxInt32}
}

// getHash returns the FNV1a hash of the tags
func (ts *timeSeriesStore) getHash() uint64 {
	return ts.hash
}

// getFStore returns the fStore in this list from field-name.
func (ts *timeSeriesStore) getFStore(fieldName string) (fStoreINTF, bool) {
	idx := sort.Search(len(ts.fStoreNodes), func(i int) bool {
		return ts.fStoreNodes[i].getFieldName() >= fieldName
	})
	if idx >= len(ts.fStoreNodes) || ts.fStoreNodes[idx].getFieldName() != fieldName {
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

// timeRange returns the start-time and end-time in milliseconds
func (ts *timeSeriesStore) timeRange() (timeRange timeutil.TimeRange, ok bool) {
	ts.sl.Lock()
	timeRange = timeutil.TimeRange{
		Start: ts.lastWroteTime + int64(ts.startDelta),
		End:   ts.lastWroteTime + int64(ts.endDelta)}
	ok = ts.startDelta != math.MaxInt32 && ts.endDelta != math.MaxInt32
	ts.sl.Unlock()
	return
}

// afterWrite set hasData, then update the start-time and end-time
func (ts *timeSeriesStore) afterWrite(writeCtx writeContext) {
	// hasData is true now
	ts.hasData = true
	// update time range
	now := timeutil.Now()
	pointTime := writeCtx.familyTime + writeCtx.timeInterval*int64(writeCtx.slotIndex)
	startTime := ts.lastWroteTime + int64(ts.startDelta)
	endTime := ts.lastWroteTime + int64(ts.endDelta)
	if ts.startDelta == math.MaxInt32 || pointTime < startTime {
		ts.startDelta = int32(pointTime - now)
	}
	if ts.endDelta == math.MaxInt32 || endTime < pointTime {
		ts.endDelta = int32(pointTime - now)
	}
	ts.lastWroteTime = now
}

// afterFlush update the flag of hasData, and the start-time and end-time
func (ts *timeSeriesStore) afterFlush(flushCtx flushContext) {
	// update hasData flag
	var startTime, endTime int64
	ts.hasData = false
	for _, fStore := range ts.fStoreNodes {
		timeRange, ok := fStore.timeRange(flushCtx.timeInterval)
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
	if ts.hasData {
		ts.startDelta = int32(startTime - ts.lastWroteTime)
		ts.endDelta = int32(endTime - ts.lastWroteTime)
	} else {
		ts.startDelta = math.MaxInt32
		ts.endDelta = math.MaxInt32
	}
}

// isExpired detects if this tStore has not been used for a TTL
func (ts *timeSeriesStore) isExpired() bool {
	return ts.lastWroteTime+getTagsIDTTL() < timeutil.Now()
}

// write write the data of metric to the fStore.
func (ts *timeSeriesStore) write(metric *pb.Metric, writeCtx writeContext) error {
	ts.sl.Lock()
	defer ts.sl.Unlock()

	for _, f := range metric.Fields {
		fieldType := getFieldType(f)
		if fieldType == field.Unknown {
			//TODO add log or metric
			continue
		}
		if fStore, err := ts.getOrCreateFStore(f.Name, fieldType, writeCtx); err == nil {
			fStore.write(f, writeCtx)
			ts.afterWrite(writeCtx)
		} else {
			return err // field type do not match before, too many fields
		}
	}
	return nil
}

// getOrCreateFStore get a fieldStore by fieldName and fieldType
func (ts *timeSeriesStore) getOrCreateFStore(fieldName string, fieldType field.Type,
	writeCtx writeContext) (fStoreINTF, error) {

	fStore, ok := ts.getFStore(fieldName)
	if ok {
		if fStore.getFieldType() != fieldType {
			return nil, models.ErrWrongFieldType
		}
	} else {
		// forbid creating new fStore when full
		if len(ts.fStoreNodes) >= maxFieldsLimit {
			return nil, models.ErrTooManyFields
		}
		// wrong field type
		fieldID, err := writeCtx.generator.GenFieldID(writeCtx.metricID, fieldName, fieldType)
		if err != nil {
			return nil, err
		}
		fStore = newFieldStore(fieldName, fieldID, fieldType)
		ts.insertFStore(fStore)
	}
	return fStore, nil
}

// flushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) flushSeriesTo(tableFlusher metrictbl.TableFlusher, flushCtx flushContext) (flushed bool) {
	ts.sl.Lock()
	for _, fStore := range ts.fStoreNodes {
		fieldDataFlushed := fStore.flushFieldTo(tableFlusher, flushCtx.familyTime)
		flushed = flushed || fieldDataFlushed
	}
	if flushed {
		tableFlusher.FlushSeries(ts.seriesID)
		ts.afterFlush(flushCtx)
	}
	// update time range info
	ts.sl.Unlock()
	return
}

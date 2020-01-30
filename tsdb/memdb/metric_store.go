package memdb

import (
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./metric_store.go -destination=./metric_store_mock.go -package memdb

const emptyMStoreSize = 8 + // immutable
	8 + // mutable
	24 + // rwmutex
	8 + // atomic.Value
	4 + // uint32
	4 + // uint32
	4 // int32

// mStoreINTF abstracts a metricStore
type mStoreINTF interface {
	mStoreFieldIDGetter

	// flow.DataFilter filters the data based on condition
	flow.DataFilter

	///////////////////////////////////
	// Methods below will change the memory size
	///////////////////////////////////
	// Write Writes the metric
	Write(seriesID uint32,
		fields []*pb.Field,
		writeCtx writeContext,
	) (writtenSize int, err error)

	// FlushMetricsDataTo flushes metric-block of mStore to the Writer.
	FlushMetricsDataTo(tableFlusher metricsdata.Flusher, flushCtx flushContext) (err error)
}

type mStoreFieldIDGetter interface {
	// GetFieldIDOrGenerate gets fieldID from fieldsMeta
	// and calls the id-generator when it's not exist
	GetFieldIDOrGenerate(
		metricID uint32,
		fieldName string,
		fieldType field.Type,
		generator metadb.IDGenerator,
	) (
		fieldID uint16, err error)
}

// metricStore is composed of the immutable part and mutable part of indexes.
// evictor scans the index to check which of them should be purged from the mutable part.
// flusher flushes both the immutable and mutable index to disk,
// after flushing, the immutable part will be removed.
type metricStore struct {
	immutable   atomic.Value // lock free immutable index that has not been flushed to disk
	mutable     tagIndexINTF // active mutable index in use
	mux         sync.RWMutex // read-Write lock for mutable index and fieldMetas
	fieldsMetas atomic.Value // read only, storing (field.Metas), hold mux before storing new value
	start, end  uint16       // time slot range
}

// newMetricStore returns a new mStoreINTF.
func newMetricStore() mStoreINTF {
	mutable := newTagIndex()
	ms := metricStore{
		mutable: mutable,
	}
	var fm field.Metas
	ms.fieldsMetas.Store(fm)
	return &ms
}

// getFieldIDOrGenerate gets fieldID from fieldsMeta, and calls the id-generator when not exist
func (ms *metricStore) GetFieldIDOrGenerate(
	metricID uint32,
	fieldName string,
	fieldType field.Type,
	generator metadb.IDGenerator,
) (
	fieldID uint16,
	err error,
) {
	fmList := ms.fieldsMetas.Load().(field.Metas)
	fm, ok := fmList.GetFromName(fieldName)
	// exist, check fieldType
	if ok {
		if fm.Type == fieldType {
			return fm.ID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// forbid creating new fStore when full
	if fmList.Len() >= constants.TStoreMaxFieldsCount {
		return 0, series.ErrTooManyFields
	}
	// not exist, create a new one
	ms.mux.Lock()
	defer ms.mux.Unlock()

	fmList = ms.fieldsMetas.Load().(field.Metas)
	fm, ok = fmList.GetFromName(fieldName)
	// double check
	if ok {
		return fm.ID, nil
	}
	// generate and check fieldType
	newFieldID, err := generator.GenFieldID(metricID, fieldName, fieldType)
	if err != nil { // fieldType not matches to the existed
		return 0, err
	}
	x2 := fmList.Clone()
	x2 = x2.Insert(field.Meta{
		Name: fieldName,
		ID:   newFieldID,
		Type: fieldType})
	// store the new clone
	ms.fieldsMetas.Store(x2)
	return newFieldID, nil
}

// Write Writes the metric to the tStore
func (ms *metricStore) Write(seriesID uint32,
	fields []*pb.Field,
	writeCtx writeContext,
) (
	writtenSize int,
	err error,
) {
	//FIXME stone100 add metric version store
	ms.mux.Lock()
	tStore, createdSize := ms.mutable.GetOrCreateTStore(seriesID)
	ms.mux.Unlock()

	writtenSize, err = tStore.Write(fields, writeCtx)
	if err == nil {
		slot := writeCtx.slotIndex
		ms.mux.Lock()
		// set metric level slot range
		if slot < ms.start {
			ms.start = slot
		}
		if slot > ms.end {
			ms.end = slot
		}
		ms.mux.Unlock()
	}
	return writtenSize + createdSize, err
}

func (ms *metricStore) atomicGetImmutable() tagIndexINTF {
	immutable, ok := ms.immutable.Load().(tagIndexINTF)
	// version zero is the placeholder tagIndexINTF stored in atomic.Value
	if ok && immutable.Version() != 0 {
		return immutable
	}
	return nil
}

// FlushMetricsTo Writes metric-data to the table.
// immutable tagIndex will be removed after call,
// index shall be flushed before flushing data.
func (ms *metricStore) FlushMetricsDataTo(
	flusher metricsdata.Flusher,
	flushCtx flushContext,
) (
	err error,
) {
	// flush field meta info
	fmList := ms.fieldsMetas.Load().(field.Metas)
	flusher.FlushFieldMetas(fmList)
	//FIXME stone1100, need refactor index/data store
	flushCtx.start, flushCtx.end = ms.start, ms.end

	// reset the mutable part
	ms.mux.RLock()
	ms.mutable.FlushVersionDataTo(flusher, flushCtx)
	immutable := ms.atomicGetImmutable()
	// remove the immutable, put the nopTagIndex into it
	ms.immutable.Store(staticNopTagIndex)
	ms.mux.RUnlock()

	if immutable != nil {
		immutable.FlushVersionDataTo(flusher, flushCtx)
	}
	return flusher.FlushMetric(flushCtx.metricID)
}

func (ms *metricStore) TimeSlotRange() (start, end uint16) {
	return ms.start, ms.end
}

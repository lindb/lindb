package mem

import (
	"sync"
	"github.com/eleme/lindb/models"
)

type PrimitiveField struct {
}

type SegmentStore struct {
	baseTime int64
}
type FieldStore struct {
	segments    map[int64]*SegmentStore
	lastUpdated int64
	mux         sync.Mutex
}

type TimeSeriesStore struct {
	measurementStore *MeasurementStore
	tsId             uint32
	fields           map[string]*FieldStore
	lastUpdated      int64
	mux              sync.Mutex
}

func newTimeSeriesStore(measurementStore *MeasurementStore) *TimeSeriesStore {
	return &TimeSeriesStore{
		tsId:             0,
		measurementStore: measurementStore,
	}
}

func newFieldStore() *FieldStore {
	return &FieldStore{}
}

func newSegmentStore(baseTime int64) *SegmentStore {
	return &SegmentStore{
		baseTime: baseTime,
	}
}

func (t *TimeSeriesStore) GetFieldStore(field string) *FieldStore {
	var store, ok = t.fields[field]
	if !ok {
		store = newFieldStore()
		t.mux.Lock()
		defer t.mux.Unlock()
		t.fields[field] = store
	}
	return store
}

func (f *FieldStore) GetSegmentStore(segmentTime int64) *SegmentStore {
	var store, ok = f.segments[segmentTime]
	if !ok {
		store = newSegmentStore(segmentTime)
		f.mux.Lock()
		defer f.mux.Unlock()
		f.segments[segmentTime] = store
	}
	return store
}

func (t *SegmentStore) Write(slotTime int32, field models.Field) {
}

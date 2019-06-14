package mem

import (
	"sync"
)

type MemoryDatabase struct {
	measurements map[string]*MeasurementStore
	mux          sync.Mutex
}

func NewMemDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		measurements: make(map[string]*MeasurementStore),
	}
}

// get time series store by measurement + tags
func (db *MemoryDatabase) GetTimeSeriesStore(measurement string, tags string) *TimeSeriesStore {
	measurementStore := db.getMeasurementStore(measurement)

	return measurementStore.getTimeSeries(tags)
}

// get measurement store by name
func (db *MemoryDatabase) getMeasurementStore(measurement string) *MeasurementStore {
	var store, ok = db.measurements[measurement]
	if !ok {
		store = newMeasurement()
		db.mux.Lock()
		defer db.mux.Unlock()
		db.measurements[measurement] = store
	}
	return store
}

package mem

import "sync"

type MeasurementStore struct {
	timeSeries map[string]*TimeSeriesStore
	//tsSeq      uint32
	mux sync.Mutex
}

func newMeasurement() *MeasurementStore {
	return &MeasurementStore{
		timeSeries: make(map[string]*TimeSeriesStore),
	}
}

func (m *MeasurementStore) getTimeSeries(tags string) *TimeSeriesStore {
	var store, ok = m.timeSeries[tags]
	if !ok {
		store = newTimeSeriesStore(m)
		m.mux.Lock()
		defer m.mux.Unlock()
		m.timeSeries[tags] = store
	}
	return store
}

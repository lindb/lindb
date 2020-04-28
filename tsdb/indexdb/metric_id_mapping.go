package indexdb

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
)

// MetricIDMapping represents the metric id mapping,
// tag hash code => series id
type MetricIDMapping interface {
	// GetMetricID return the metric id
	GetMetricID() uint32
	// GetSeriesID gets series id by tags hash, if exist return true
	GetSeriesID(tagsHash uint64) (seriesID uint32, ok bool)
	// GenSeriesID generates series id by tags hash, then cache new series id
	GenSeriesID(tagsHash uint64) (seriesID uint32)
	// RemoveSeriesID removes series id by tags hash
	RemoveSeriesID(tagsHash uint64)
	// AddSeriesID adds the series id init cache
	AddSeriesID(tagsHash uint64, seriesID uint32)
	// SetMaxSeriesIDsLimit sets the max series ids limit
	SetMaxSeriesIDsLimit(limit uint32)
	// GetMaxSeriesIDsLimit returns the max series ids limit
	GetMaxSeriesIDsLimit() uint32
}

// metricIDMapping implements MetricIDMapping interface
type metricIDMapping struct {
	metricID uint32
	// forwardIndex for storing a mapping from tag-hash to the seriesID,
	// purpose of this index is used for fast writing
	hash2SeriesID     map[uint64]uint32
	idSequence        atomic.Uint32
	maxSeriesIDsLimit atomic.Uint32 // maximum number of combinations of series ids
}

// newMetricIDMapping returns a new metric id mapping
func newMetricIDMapping(metricID, sequence uint32) MetricIDMapping {
	return &metricIDMapping{
		metricID:          metricID,
		hash2SeriesID:     make(map[uint64]uint32),
		idSequence:        *atomic.NewUint32(sequence), // first value is 1
		maxSeriesIDsLimit: *atomic.NewUint32(constants.DefaultMaxSeriesIDsCount),
	}
}

// GetMetricID return the metric id
func (mim *metricIDMapping) GetMetricID() uint32 {
	return mim.metricID
}

// GetSeriesID gets series id by tags hash, if exist return true
func (mim *metricIDMapping) GetSeriesID(tagsHash uint64) (seriesID uint32, ok bool) {
	seriesID, ok = mim.hash2SeriesID[tagsHash]
	return
}

// AddSeriesID adds the series id init cache
func (mim *metricIDMapping) AddSeriesID(tagsHash uint64, seriesID uint32) {
	mim.hash2SeriesID[tagsHash] = seriesID
}

// GenSeriesID generates series id by tags hash, then cache new series id
func (mim *metricIDMapping) GenSeriesID(tagsHash uint64) (seriesID uint32) {
	// generate new series id
	if mim.maxSeriesIDsLimit.Load() == mim.idSequence.Load() {
		//FIXME too many series id, use max limit????
		seriesID = mim.maxSeriesIDsLimit.Load()
	} else {
		seriesID = mim.idSequence.Inc()
	}
	// cache it
	mim.hash2SeriesID[tagsHash] = seriesID
	return seriesID
}

// RemoveSeriesID removes series id by tags hash
func (mim *metricIDMapping) RemoveSeriesID(tagsHash uint64) {
	seriesID, ok := mim.hash2SeriesID[tagsHash]
	if ok {
		if seriesID == mim.idSequence.Load() {
			mim.idSequence.Dec() // recycle series id
		}
		delete(mim.hash2SeriesID, tagsHash)
	}
}

// SetMaxSeriesIDsLimit sets the max series ids limit
func (mim *metricIDMapping) SetMaxSeriesIDsLimit(limit uint32) {
	mim.maxSeriesIDsLimit.Store(limit)
}

// GetMaxSeriesIDsLimit return the max series ids limit without race condition.
func (mim *metricIDMapping) GetMaxSeriesIDsLimit() uint32 {
	return mim.maxSeriesIDsLimit.Load()
}

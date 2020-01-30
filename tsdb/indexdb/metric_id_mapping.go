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
	// GetOrCreateSeriesID gets or creates series id by tags hash,
	// if is new series id return created is true
	GetOrCreateSeriesID(tagsHash uint64) (seriesID uint32, created bool)
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
func newMetricIDMapping(metricID uint32) MetricIDMapping {
	return &metricIDMapping{
		metricID:          metricID,
		hash2SeriesID:     make(map[uint64]uint32),
		idSequence:        *atomic.NewUint32(0), // first value is 1
		maxSeriesIDsLimit: *atomic.NewUint32(constants.DefaultMaxSeriesIDsCount),
	}
}

// GetMetricID return the metric id
func (index *metricIDMapping) GetMetricID() uint32 {
	return index.metricID
}

// GetOrCreateSeriesID gets or creates series id by tags hash
func (index *metricIDMapping) GetOrCreateSeriesID(tagsHash uint64) (seriesID uint32, created bool) {
	seriesID, ok := index.hash2SeriesID[tagsHash]
	if ok {
		return seriesID, false
	}

	if index.maxSeriesIDsLimit.Load() == index.idSequence.Load() {
		//FIXME too many series id, use max limit????
		seriesID = index.maxSeriesIDsLimit.Load()
	} else {
		seriesID = index.idSequence.Inc()
	}
	index.hash2SeriesID[tagsHash] = seriesID
	return seriesID, true
}

// SetMaxSeriesIDsLimit sets the max series ids limit
func (index *metricIDMapping) SetMaxSeriesIDsLimit(limit uint32) {
	index.maxSeriesIDsLimit.Store(limit)
}

// GetMaxSeriesIDsLimit return the max series ids limit without race condition.
func (index *metricIDMapping) GetMaxSeriesIDsLimit() uint32 {
	return index.maxSeriesIDsLimit.Load()
}

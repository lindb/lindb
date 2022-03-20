// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package indexdb

import (
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series/metric"

	"go.uber.org/atomic"
)

//go:generate mockgen -source ./metric_id_mapping.go -destination=./metric_id_mapping_mock.go -package=indexdb

// MetricIDMapping represents the metric id mapping,
// tag hash code => series id
type MetricIDMapping interface {
	// GetMetricID return the metric id.
	GetMetricID() metric.ID
	// GetSeriesID gets series id by tags hash, if exist return true.
	GetSeriesID(tagsHash uint64) (seriesID uint32, ok bool)
	// GenSeriesID generates series id by tags hash, then cache new series id.
	GenSeriesID(tagsHash uint64) (seriesID uint32)
	// AddSeriesID adds the series id init cache.
	AddSeriesID(tagsHash uint64, seriesID uint32)
	// SeriesSequence returns series sequence.
	SeriesSequence() unique.Sequence
	// SetMaxSeriesIDsLimit sets the max series ids limit.
	SetMaxSeriesIDsLimit(limit uint32)
	// GetMaxSeriesIDsLimit returns the max series ids limit.
	GetMaxSeriesIDsLimit() uint32
}

// metricIDMapping implements MetricIDMapping interface.
type metricIDMapping struct {
	metricID metric.ID
	// forwardIndex for storing a mapping from tag-hash to the seriesID,
	// purpose of this index is used for fast writing
	hash2SeriesID map[uint64]uint32
	idSequence    unique.Sequence // first value is 1
	// TODO need remove
	maxSeriesIDsLimit atomic.Uint32 // maximum number of combinations of series ids
}

// newMetricIDMapping returns a new metric id mapping.
func newMetricIDMapping(metricID metric.ID, sequence uint32) MetricIDMapping {
	return &metricIDMapping{
		metricID:          metricID,
		hash2SeriesID:     make(map[uint64]uint32),
		idSequence:        unique.NewSequence(sequence, config.GlobalStorageConfig().TSDB.SeriesSequenceCache),
		maxSeriesIDsLimit: *atomic.NewUint32(uint32(config.GlobalStorageConfig().TSDB.MaxSeriesIDsNumber)),
	}
}

// GetMetricID return the metric id.
func (mim *metricIDMapping) GetMetricID() metric.ID {
	return mim.metricID
}

// GetSeriesID gets series id by tags hash, if exist return true.
func (mim *metricIDMapping) GetSeriesID(tagsHash uint64) (seriesID uint32, ok bool) {
	seriesID, ok = mim.hash2SeriesID[tagsHash]
	return
}

// AddSeriesID adds the series id init cache.
func (mim *metricIDMapping) AddSeriesID(tagsHash uint64, seriesID uint32) {
	mim.hash2SeriesID[tagsHash] = seriesID
}

// GenSeriesID generates series id by tags hash, then cache new series id.
func (mim *metricIDMapping) GenSeriesID(tagsHash uint64) (seriesID uint32) {
	// generate new series id
	if mim.idSequence.Current() >= mim.maxSeriesIDsLimit.Load() {
		// FIXME too many series id, use max limit????
		seriesID = mim.maxSeriesIDsLimit.Load()
	} else {
		seriesID = mim.idSequence.Next()
	}
	// cache it
	mim.hash2SeriesID[tagsHash] = seriesID
	return seriesID
}

// SeriesSequence returns series sequence.
func (mim *metricIDMapping) SeriesSequence() unique.Sequence {
	return mim.idSequence
}

// SetMaxSeriesIDsLimit sets the max series ids limit.
func (mim *metricIDMapping) SetMaxSeriesIDsLimit(limit uint32) {
	mim.maxSeriesIDsLimit.Store(limit)
}

// GetMaxSeriesIDsLimit return the max series ids limit without race condition.
func (mim *metricIDMapping) GetMaxSeriesIDsLimit() uint32 {
	return mim.maxSeriesIDsLimit.Load()
}

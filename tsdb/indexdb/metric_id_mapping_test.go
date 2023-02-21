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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
)

func TestMetricIDMapping_GetMetricID(t *testing.T) {
	idMapping := newMetricIDMapping(10, 0)
	assert.Equal(t, metric.ID(10), idMapping.GetMetricID())
	assert.NotNil(t, idMapping.SeriesSequence())
}

func TestMetricIDMapping_GetOrCreateSeriesID(t *testing.T) {
	limits := models.NewDefaultLimits()
	idMapping := newMetricIDMapping(10, 0)
	seriesID, ok := idMapping.GetSeriesID(100)
	assert.False(t, ok)
	assert.Equal(t, uint32(0), seriesID)
	seriesID, err := idMapping.GenSeriesID("ns", "metric", 100, limits)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), seriesID)
	// get exist series id
	seriesID, ok = idMapping.GetSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	assert.True(t, ok)

	// add series id
	idMapping.AddSeriesID(300, 4)
	seriesID, ok = idMapping.GetSeriesID(300)
	assert.Equal(t, uint32(4), seriesID)
	assert.True(t, ok)
}

func TestMetricIDMapping_SeriesLimit(t *testing.T) {
	limits := models.NewDefaultLimits()
	limits.MaxSeriesPerMetric = 1
	idMapping := newMetricIDMapping(10, 0)
	seriesID, err := idMapping.GenSeriesID("ns", "metric", 100, limits)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), seriesID)
	// gt limit
	seriesID, err = idMapping.GenSeriesID("ns", "metric", 1023, limits)
	assert.Error(t, err)
	assert.Equal(t, series.EmptySeriesID, seriesID)
}

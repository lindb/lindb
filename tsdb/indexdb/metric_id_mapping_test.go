package indexdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
)

func TestMetricIDMapping_GetMetricID(t *testing.T) {
	idMapping := newMetricIDMapping(10, 0)
	assert.Equal(t, uint32(10), idMapping.GetMetricID())
}

func TestMetricIDMapping_GetOrCreateSeriesID(t *testing.T) {
	idMapping := newMetricIDMapping(10, 0)
	seriesID, ok := idMapping.GetSeriesID(100)
	assert.False(t, ok)
	assert.Equal(t, uint32(0), seriesID)
	seriesID = idMapping.GenSeriesID(100)
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

func TestMetricIDMapping_SetMaxTagsLimit(t *testing.T) {
	idMapping := newMetricIDMapping(10, 0)
	seriesID := idMapping.GenSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	assert.Equal(t, uint32(constants.DefaultMaxSeriesIDsCount), idMapping.GetMaxSeriesIDsLimit())
	idMapping.SetMaxSeriesIDsLimit(2)
	_ = idMapping.GenSeriesID(102)
	seriesID = idMapping.GenSeriesID(1020)
	assert.Equal(t, uint32(2), seriesID)
}

func TestMetricIDMapping_RemoveSeriesID(t *testing.T) {
	idMapping := newMetricIDMapping(10, 0)
	seriesID := idMapping.GenSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	idMapping.RemoveSeriesID(100)
	seriesID = idMapping.GenSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	idMapping.RemoveSeriesID(1200)
}

package indexdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
)

func TestMetricIDMapping_GetMetricID(t *testing.T) {
	idMapping := newMetricIDMapping(10)
	assert.Equal(t, uint32(10), idMapping.GetMetricID())
}

func TestMetricIDMapping_GetOrCreateSeriesID(t *testing.T) {
	idMapping := newMetricIDMapping(10)
	seriesID, created := idMapping.GetOrCreateSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	assert.True(t, created)
	// get exist series id
	seriesID, created = idMapping.GetOrCreateSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	assert.False(t, created)
}

func TestMetricIDMapping_SetMaxTagsLimit(t *testing.T) {
	idMapping := newMetricIDMapping(10)
	seriesID, created := idMapping.GetOrCreateSeriesID(100)
	assert.Equal(t, uint32(1), seriesID)
	assert.True(t, created)
	assert.Equal(t, uint32(constants.DefaultMaxSeriesIDsCount), idMapping.GetMaxSeriesIDsLimit())
	idMapping.SetMaxSeriesIDsLimit(2)
	seriesID, created = idMapping.GetOrCreateSeriesID(102)
	assert.Equal(t, uint32(2), seriesID)
	assert.True(t, created)
	seriesID, created = idMapping.GetOrCreateSeriesID(1020)
	assert.Equal(t, uint32(2), seriesID)
	assert.True(t, created)
}

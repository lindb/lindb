package indexdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/tsdb/metadb"
)

func TestMemoryIndexDatabase_GetTimeSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	generator := metadb.NewMockIDGenerator(ctrl)
	generator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()
	generator.EXPECT().GenMetricID("cpu").Return(uint32(1)).AnyTimes()
	db := NewMemoryIndexDatabase(generator)
	metricID, seriesID := db.GetTimeSeriesID("cpu", map[string]string{
		"host": "1,1,1,1",
	}, 10)
	assert.Equal(t, uint32(1), seriesID)
	assert.Equal(t, uint32(1), metricID)

	generator.EXPECT().GenMetricID("cpu").Return(uint32(1)).AnyTimes()
	metricID, seriesID = db.GetTimeSeriesID("cpu", map[string]string{
		"host": "1,1,1,1",
	}, 10)
	assert.Equal(t, uint32(1), seriesID)
	assert.Equal(t, uint32(1), metricID)
}

func TestMemoryIndexDatabase_FlushInvertedIndexTo(t *testing.T) {
	//FIXME stone1100 need impl
	db := NewMemoryIndexDatabase(nil)
	err := db.FlushInvertedIndexTo(nil)
	assert.NoError(t, err)
}

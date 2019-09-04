package tblstore

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_MetricsNameIDFlusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("write failure")).AnyTimes()
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit failure")).AnyTimes()

	nameIDFlusher := NewMetricsNameIDFlusher(mockFlusher)
	assert.NotNil(t, nameIDFlusher.FlushMetricsNS(1, nil, 1, 2))
	assert.NotNil(t, nameIDFlusher.Commit())
}

func Test_MetricsMetaFlusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Commit().Return(nil).AnyTimes()

	metaFlusher := NewMetricsMetaFlusher(mockFlusher)
	// write only tags
	mockFlusher.EXPECT().Add(uint32(1), []byte{
		14, 2, 107, 49, 1, 0, 0, 0, 2, 107, 50, 2, 0, 0, 0, 0}).
		Return(nil)
	metaFlusher.FlushTagKeyID("k1", 1)
	metaFlusher.FlushTagKeyID("k2", 2)
	metaFlusher.FlushMetricMeta(1)
	assert.Nil(t, metaFlusher.Commit())
	// write only fields
	metaFlusher.FlushFieldID("f3", field.SumField, 3)
	metaFlusher.FlushFieldID("f4", field.MinField, 4)
	mockFlusher.EXPECT().Add(uint32(2), []byte{
		0, 12, 2, 102, 51, 1, 3, 0, 2, 102, 52, 2, 4, 0}).
		Return(nil)
	metaFlusher.FlushMetricMeta(2)
	assert.Nil(t, metaFlusher.Commit())
	// write tags fields
	mockFlusher.EXPECT().Add(uint32(3), []byte{
		7, 2, 107, 49, 1, 0, 0, 0, 6, 2, 102, 51, 1, 3, 0}).
		Return(nil)
	metaFlusher.FlushTagKeyID("k1", 1)
	metaFlusher.FlushFieldID("f3", field.SumField, 3)
	metaFlusher.FlushMetricMeta(3)
	assert.Nil(t, metaFlusher.Commit())
}

func Test_flusher_invalid_input(t *testing.T) {
	badKey := ""
	for i := 0; i < 1000; i++ {
		badKey += "X"
	}

	metaFlusher := NewMetricsMetaFlusher(nil)
	metaFlusher.FlushTagKeyID("", 1)
	metaFlusher.FlushTagKeyID(badKey, 1)
	metaFlusher.FlushFieldID("", field.SumField, 1)
	metaFlusher.FlushFieldID(badKey, field.SumField, 1)
}

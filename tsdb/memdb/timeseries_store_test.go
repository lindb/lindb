package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestTimeSeriesStore_GetOrCreateFStore(t *testing.T) {
	tStore := newTimeSeriesStore()
	f, ok := tStore.GetFStore(1, 10, 11)
	assert.Nil(t, f)
	assert.False(t, ok)
	tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10, 11))
	// get field store
	f, ok = tStore.GetFStore(1, 10, 11)
	assert.NotNil(t, f)
	assert.True(t, ok)
	// field store not exist
	f, ok = tStore.GetFStore(1, 10, 10)
	assert.Nil(t, f)
	assert.False(t, ok)
	for i := 1; i < 10; i++ {
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), familyID(1*i), field.ID(10*i), field.PrimitiveID(11*i)))
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10, 11))
		f, ok = tStore.GetFStore(1, 10, 11)
		assert.NotNil(t, f)
		assert.True(t, ok)
	}
}

func TestTimeSeriesStore_FlushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)
	tStore := newTimeSeriesStore()
	s := tStore.(*timeSeriesStore)
	fStore := NewMockfStoreINTF(ctrl)
	s.InsertFStore(fStore)
	gomock.InOrder(
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()),
	)
	tStore.FlushSeriesTo(flusher, flushContext{})
}

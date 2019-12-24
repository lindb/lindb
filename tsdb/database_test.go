package tsdb

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/metadb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_Database_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)
	mockStore := kv.NewMockStore(ctrl)
	db := &database{
		idSequencer: mockIDSequencer,
		metaStore:   mockStore}
	mockStore.EXPECT().Close().Return(nil).AnyTimes()

	// mock flush metrics-meta error
	mockIDSequencer.EXPECT().FlushMetricsMeta().Return(fmt.Errorf("error"))
	assert.Nil(t, db.Close())
	// mock flush nameids error
	mockIDSequencer.EXPECT().FlushMetricsMeta().Return(nil)
	mockIDSequencer.EXPECT().FlushNameIDs().Return(fmt.Errorf("error"))
	assert.Nil(t, db.Close())
	// mock shard close error
	mockIDSequencer.EXPECT().FlushMetricsMeta().Return(nil)
	mockIDSequencer.EXPECT().FlushNameIDs().Return(nil)
	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(fmt.Errorf("error"))
	db.shards.Store(int32(1), mockShard)
	assert.Nil(t, db.Close())

	assert.NotNil(t, db.IDGetter())
	assert.NotNil(t, db.MetricMetaSuggester())
}

func Test_Database_FlushMeta(t *testing.T) {
	db := &database{
		isFlushing: *atomic.NewBool(false)}
	db.isFlushing.Store(true)
	assert.Nil(t, db.FlushMeta())

	db.Range(func(key, value interface{}) bool {
		assert.Fail(t, "")
		return true
	})
}

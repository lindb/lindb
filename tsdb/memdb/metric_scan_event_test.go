package memdb

import (
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/series"
)

func Test_pool(t *testing.T) {
	stores := getStores()
	for idx := range stores {
		stores[idx] = nil
	}
	putStores(stores)
	stores = getStores()
	fmt.Println(len(stores))
}

func TestMetricScanEvent_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStore := NewMocktStoreINTF(ctrl)
	sCtx := &series.ScanContext{
		FieldIDs: []uint16{3, 4, 5},
	}
	stores := getStores()
	stores[0] = tStore
	seriesIDs := *series.Uint32Pool.Get()
	seriesIDs[0] = uint32(1)
	// test not match aggregator
	event := newScanEvent(1, stores, seriesIDs, series.Version(1), sCtx)
	ok := event.Scan()
	assert.False(t, ok)
	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	sCtx.Aggregators = sync.Pool{
		New: func() interface{} {
			return aggregation.FieldAggregates{sAgg}
		},
	}

	// test normal case
	gomock.InOrder(
		tStore.EXPECT().scan(gomock.Any()),
	)
	stores = getStores()
	stores[0] = tStore
	seriesIDs = *series.Uint32Pool.Get()
	seriesIDs[0] = uint32(1)
	event = newScanEvent(1, stores, seriesIDs, series.Version(1), sCtx)
	ok = event.Scan()
	assert.True(t, ok)

	resultSet := event.ResultSet()
	assert.Equal(t, aggregation.FieldAggregates{sAgg}, resultSet)
	sAgg.EXPECT().Reset()
	event.Release()
}

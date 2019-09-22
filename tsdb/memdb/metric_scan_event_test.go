package memdb

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/series"
)

func TestMetricScanEvent_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStore := NewMocktStoreINTF(ctrl)
	sCtx := &series.ScanContext{
		FieldIDs: []uint16{3, 4, 5},
	}
	// test not match aggregator
	event := newScanEvent(1,
		[]tStoreINTF{tStore},
		[]uint32{1}, series.Version(1),
		sCtx)
	ok := event.Scan()
	assert.False(t, ok)
	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	sCtx.Aggregates = sync.Pool{
		New: func() interface{} {
			return aggregation.FieldAggregates{sAgg}
		},
	}

	// test normal case
	gomock.InOrder(
		tStore.EXPECT().scan(gomock.Any()),
	)
	event = newScanEvent(1,
		[]tStoreINTF{tStore},
		[]uint32{1}, series.Version(1),
		sCtx)
	ok = event.Scan()
	assert.True(t, ok)

	resultSet := event.ResultSet()
	assert.Equal(t, aggregation.FieldAggregates{sAgg}, resultSet)
	sAgg.EXPECT().Reset()
	event.Release()
}

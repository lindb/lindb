package query

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/series"
)

func TestScanWorker_Emit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), nil, nil, groupAgg, nil)
	event := series.NewMockScanEvent(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(false),
		groupAgg.EXPECT().ResultSet().Return(nil),
	)
	worker.Emit(event)
	worker.Emit(nil)
	worker.Close()

	worker = createScanWorker(context.TODO(), uint32(10), nil, nil, groupAgg, nil)
	gomock.InOrder(
		groupAgg.EXPECT().ResultSet().Return(nil),
	)
	w := worker.(*scanWorker)
	w.events <- nil
	worker.Close()
}

func TestScanWorker_handle_event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)

	worker := createScanWorker(context.TODO(), uint32(10), nil, nil, groupAgg, nil)
	event := series.NewMockScanEvent(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(true),
		event.EXPECT().ResultSet().Return(nil),
		event.EXPECT().Release(),
		groupAgg.EXPECT().ResultSet().Return(nil),
	)
	worker.Emit(event)
	worker.Close()

	worker = createScanWorker(context.TODO(), uint32(10), nil, nil, groupAgg, nil)
	event = series.NewMockScanEvent(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(true),
		event.EXPECT().ResultSet().Return("mock"),
		event.EXPECT().Release(),
		groupAgg.EXPECT().ResultSet().Return(nil),
	)
	worker.Emit(event)
	worker.Close()

	rs := make(chan *series.TimeSeriesEvent)
	c := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		<-rs
		c.Inc()
		wait.Done()
	}()

	worker = createScanWorker(context.TODO(), uint32(10), nil, nil, groupAgg, rs)
	event = series.NewMockScanEvent(ctrl)
	seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(true),
		event.EXPECT().ResultSet().Return(aggregation.FieldAggregates{seriesAgg}),
		groupAgg.EXPECT().Merge(gomock.Any(), gomock.Any()),
		event.EXPECT().Release(),
		groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{series.NewMockGroupedIterator(ctrl)}),
	)
	worker.Emit(event)
	worker.Close()
	wait.Wait()
	assert.Equal(t, int32(1), c.Load())
}

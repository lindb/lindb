package query

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb"
)

var execPool = &tsdb.ExecutePool{
	Scan:  concurrent.NewPool("test-pool1", 10 /*nRoutines*/, 10 /*queueSize*/),
	Merge: concurrent.NewPool("test-pool2", 10 /*nRoutines*/, 10 /*queueSize*/),
}

func TestScanWorker_Emit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
	exeCtx := parallel.NewMockExecuteContext(ctrl)

	worker := createScanWorker(exeCtx, uint32(10), nil, nil, groupAgg, execPool)
	event := series.NewMockScanEvent(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(false),
	)
	worker.Emit(event)
	worker.Emit(nil)
	time.Sleep(100 * time.Millisecond)
	worker.Close()
}

func TestScanWorker_handle_event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exeCtx := parallel.NewMockExecuteContext(ctrl)
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
	agg := aggregation.NewMockSeriesAggregator(ctrl)
	fieldAggregates := aggregation.FieldAggregates{agg}

	worker := createScanWorker(exeCtx, uint32(10), nil, nil, groupAgg, execPool)
	event := series.NewMockScanEvent(ctrl)
	gomock.InOrder(
		event.EXPECT().Scan().Return(true),
		event.EXPECT().ResultSet().Return(fieldAggregates),
		groupAgg.EXPECT().Aggregate(gomock.Any()),
		event.EXPECT().Release(),
		groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{nil}),
		exeCtx.EXPECT().Emit(gomock.Any()),
		exeCtx.EXPECT().Complete(nil),
	)
	worker.Emit(event)
	worker.Close()
	time.Sleep(100 * time.Millisecond)
}

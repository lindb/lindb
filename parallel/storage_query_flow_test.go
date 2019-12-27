package parallel

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb"
)

var testExecPool = &tsdb.ExecutorPool{
	Filtering: concurrent.NewPool(
		"test-filtering-pool",
		runtime.NumCPU(), /*nRoutines*/
		time.Second*5),
	Grouping: concurrent.NewPool(
		"test-grouping-pool",
		runtime.NumCPU(), /*nRoutines*/
		time.Second*5),
	Scanner: concurrent.NewPool(
		"test-scanner-pool",
		runtime.NumCPU(), /*nRoutines*/
		time.Second*5),
}

func TestStorageQueryFlow_GetAggregator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	streamHandler := pb.NewMockTaskService_HandleServer(ctrl)
	queryFlow := NewStorageQueryFlow(context.TODO(), &pb.TaskRequest{}, streamHandler, testExecPool,
		timeutil.TimeRange{}, timeutil.Interval(timeutil.OneSecond), 1)
	queryFlow.Prepare(nil)

	agg := queryFlow.GetAggregator()
	assert.NotNil(t, agg)

	qf := queryFlow.(*storageQueryFlow)
	qf.releaseAgg(agg)

	agg2 := queryFlow.GetAggregator()
	assert.NotNil(t, agg2)
	assert.Equal(t, agg, agg2)

	for i := 0; i < 100; i++ {
		qf.releaseAgg(agg)
	}
}

func TestStorageQueryFlow_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	streamHandler := pb.NewMockTaskService_HandleServer(ctrl)
	streamHandler.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	queryFlow := NewStorageQueryFlow(context.TODO(), &pb.TaskRequest{}, streamHandler, testExecPool,
		timeutil.TimeRange{}, timeutil.Interval(timeutil.OneSecond), 1)
	queryFlow.Prepare(nil)
	qf := queryFlow.(*storageQueryFlow)
	reduceAgg := aggregation.NewMockGroupingAggregator(ctrl)
	qf.reduceAgg = reduceAgg
	reduceAgg.EXPECT().ResultSet().Return(nil)
	reduceAgg.EXPECT().Aggregate(gomock.Any()).AnyTimes()

	var wait sync.WaitGroup
	wait.Add(6)
	queryFlow.Filtering(func() {
		wait.Done()
		queryFlow.Grouping(func() {
			wait.Done()
			queryFlow.Scanner(func() {
				seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
				seriesAgg.EXPECT().Reset()

				queryFlow.Reduce("1.1.1.1", aggregation.FieldAggregates{seriesAgg})
				wait.Done()
			})
		})
	})
	queryFlow.Filtering(func() {
		wait.Done()
		queryFlow.Grouping(func() {
			wait.Done()
			queryFlow.Scanner(func() {
				wait.Done()
			})
		})
	})
	wait.Wait()
	queryFlow.Complete(nil)
	seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
	seriesAgg.EXPECT().Reset()
	queryFlow.Reduce("1.1.1.1", aggregation.FieldAggregates{seriesAgg})
}

func TestStorageQueryFlow_completeTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	streamHandler := pb.NewMockTaskService_HandleServer(ctrl)
	streamHandler.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	queryFlow := NewStorageQueryFlow(context.TODO(), &pb.TaskRequest{}, streamHandler, testExecPool,
		timeutil.TimeRange{}, timeutil.Interval(timeutil.OneSecond), 1)
	queryFlow.Prepare(nil)
	qf := queryFlow.(*storageQueryFlow)
	// test error
	qf.err = fmt.Errorf("err")
	var wait sync.WaitGroup
	wait.Add(1)
	queryFlow.Filtering(func() {
		wait.Done()
	})
	wait.Wait()
	queryFlow.Complete(fmt.Errorf("err"))

	queryFlow.Filtering(func() {
		assert.Fail(t, "exec err")
	})

	// test reduce result send
	queryFlow = NewStorageQueryFlow(context.TODO(), &pb.TaskRequest{}, streamHandler, testExecPool,
		timeutil.TimeRange{}, timeutil.Interval(timeutil.OneSecond), 1)
	queryFlow.Prepare(nil)
	qf = queryFlow.(*storageQueryFlow)
	reduceAgg := aggregation.NewMockGroupingAggregator(ctrl)
	qf.reduceAgg = reduceAgg
	groupIt := series.NewMockGroupedIterator(ctrl)
	it := series.NewMockIterator(ctrl)
	groupIt.EXPECT().HasNext().Return(true)
	groupIt.EXPECT().Tags().Return("1.1.1.1").AnyTimes()
	groupIt.EXPECT().Next().Return(it)
	it.EXPECT().MarshalBinary().Return(nil, nil)
	groupIt.EXPECT().HasNext().Return(true)
	it.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("err"))
	groupIt.EXPECT().Next().Return(it)
	groupIt.EXPECT().HasNext().Return(true)
	groupIt.EXPECT().Next().Return(it)
	it.EXPECT().MarshalBinary().Return([]byte{1, 2, 3}, nil)
	it.EXPECT().FieldName().Return("f1")
	groupIt.EXPECT().HasNext().Return(false)
	reduceAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{groupIt})
	var wait1 sync.WaitGroup
	wait1.Add(1)
	queryFlow.Filtering(func() {
		wait1.Done()
	})
	wait1.Wait()
	time.Sleep(100 * time.Millisecond)
}

func TestStorageQueryFlow_Task_panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	streamHandler := pb.NewMockTaskService_HandleServer(ctrl)
	streamHandler.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	queryFlow := NewStorageQueryFlow(context.TODO(), &pb.TaskRequest{}, streamHandler, testExecPool,
		timeutil.TimeRange{}, timeutil.Interval(timeutil.OneSecond), 1)
	queryFlow.Prepare(nil)
	var wait sync.WaitGroup
	wait.Add(3)
	queryFlow.Filtering(func() {
		wait.Done()
		panic(fmt.Errorf("xxx"))
	})
	queryFlow.Filtering(func() {
		wait.Done()
		panic("err_str")
	})
	queryFlow.Filtering(func() {
		wait.Done()
		panic(12)
	})
	wait.Wait()
	time.Sleep(100 * time.Millisecond)
}

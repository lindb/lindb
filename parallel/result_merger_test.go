package parallel

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

func TestResultMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)
	groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{series.NewMockGroupedIterator(ctrl)})
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), groupAgg, ch)
	c := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		<-ch
		c.Inc()
		wait.Done()
	}()
	merger.merge(&pb.TaskResponse{TaskID: "taskID"})
	merger.close()
	wait.Wait()
	assert.Equal(t, int32(1), c.Load())
}

func TestResultMerger_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)
	groupAgg.EXPECT().ResultSet().Return(nil)
	ch := make(chan *series.TimeSeriesEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	merger := newResultMerger(ctx, groupAgg, ch)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		wait.Wait()
		cancel()
	}()
	wait.Done()
	time.Sleep(100 * time.Millisecond)
	merger.close()
}

func TestResultMerger_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), groupAgg, ch)
	c := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for rs := range ch {
			if rs.Err != nil {
				c.Inc()
				wait.Done()
			}
		}
	}()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: []byte{1, 2, 3}})
	merger.close()
	wait.Wait()
	assert.Equal(t, int32(1), c.Load())
}

func TestResultMerger_GroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	groupAgg := aggregation.NewMockGroupByAggregator(ctrl)
	groupAgg.EXPECT().Aggregate(gomock.Any()).AnyTimes()
	groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{series.NewMockGroupedIterator(ctrl)})
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), groupAgg, ch)
	c := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for rs := range ch {
			if rs.Err == nil {
				c.Inc()
				wait.Done()
			}
		}
	}()
	fields := make(map[string][]byte)
	fields["f1"] = []byte{}
	timeSeries := &pb.TimeSeries{
		Tags:   map[string]string{"host": "1.1.1.1"},
		Fields: fields,
	}
	seriesList := pb.TimeSeriesList{
		TimeSeriesList: []*pb.TimeSeries{timeSeries},
	}
	data, _ := seriesList.Marshal()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})
	timeSeries = &pb.TimeSeries{
		Tags: map[string]string{"host": "1.1.1.1"},
	}
	seriesList = pb.TimeSeriesList{
		TimeSeriesList: []*pb.TimeSeries{timeSeries},
	}
	data, _ = seriesList.Marshal()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})
	merger.close()
	wait.Wait()
	assert.Equal(t, int32(1), c.Load())
}

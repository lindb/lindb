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
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

func TestResultMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
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
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
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
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
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
	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
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
	tagValues := "1.1.1.1"
	timeSeries := &pb.TimeSeries{
		Tags:   tagValues,
		Fields: fields,
	}
	seriesList := pb.TimeSeriesList{
		TimeSeriesList: []*pb.TimeSeries{timeSeries},
	}
	data, _ := seriesList.Marshal()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})
	timeSeries = &pb.TimeSeries{
		Tags: tagValues,
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

func TestSuggestMerge_merge(t *testing.T) {
	ch := make(chan []string)
	merger := newSuggestResultMerger(ch)
	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		merger.merge(&pb.TaskResponse{
			Payload: []byte{1, 2, 3},
		})
		merger.merge(&pb.TaskResponse{
			Payload: encoding.JSONMarshal(&models.SuggestResult{Values: []string{"a"}}),
		})
		merger.merge(&pb.TaskResponse{
			Payload: encoding.JSONMarshal(&models.SuggestResult{Values: []string{"a"}}),
		})
		// close result chan
		merger.close()
	}()

	for rs := range ch {
		for _, value := range rs {
			assert.Equal(t, "a", value)
			wait.Done()
		}
	}
	wait.Wait()
}

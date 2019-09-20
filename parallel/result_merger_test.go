package parallel

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

func TestResultMerger_Merge(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(),
		&stmt.Query{
			Interval:  10000,
			TimeRange: timeutil.TimeRange{Start: 10, End: 12},
		}, ch)
	go func() {
		for range ch {
		}
	}()

	merger.merge(&pb.TaskResponse{TaskID: "taskID"})

	merger.close()
}

func TestResultMerger_Cancel(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	merger := newResultMerger(ctx,
		&stmt.Query{
			Interval:  10000,
			TimeRange: timeutil.TimeRange{Start: 10, End: 12},
		}, ch)
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
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(),
		&stmt.Query{
			Interval:  10000,
			TimeRange: timeutil.TimeRange{Start: 10, End: 12},
		}, ch)
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
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(),
		&stmt.Query{
			Interval:  10000,
			TimeRange: timeutil.TimeRange{Start: 10, End: 12},
			GroupBy:   []string{"host", "disk"},
		}, ch)
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
		Fields: fields,
	}
	data, _ := timeSeries.Marshal()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})

	timeSeries = &pb.TimeSeries{
		Tags:   map[string]string{"host": "1.1.1.1"},
		Fields: fields,
	}
	data, _ = timeSeries.Marshal()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Payload: data})
	merger.close()
	wait.Wait()
	assert.Equal(t, int32(1), c.Load())
}

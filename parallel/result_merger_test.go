// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package parallel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

func TestResultMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newGroupingAgg = aggregation.NewGroupingAggregator
		ctrl.Finish()
	}()

	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
	newGroupingAgg = func(interval timeutil.Interval,
		intervalRatio int,
		timeRange timeutil.TimeRange,
		aggSpecs aggregation.AggregatorSpecs) aggregation.GroupingAggregator {
		return groupAgg
	}
	groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{series.NewMockGroupedIterator(ctrl)})
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), &stmt.Query{}, ch)
	c := atomic.NewInt32(0)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		<-ch
		c.Inc()
		wait.Done()
	}()
	merger.merge(&pb.TaskResponse{TaskID: "taskID", Stats: encoding.JSONMarshal(models.NewStorageStats())})
	merger.close()
	fmt.Println("1")
	wait.Wait()
	fmt.Println("2")
	assert.Equal(t, int32(1), c.Load())
}

func TestResultMerger_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	ch := make(chan *series.TimeSeriesEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	merger := newResultMerger(ctx, &stmt.Query{}, ch)
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
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), &stmt.Query{}, ch)
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
	defer func() {
		newGroupingAgg = aggregation.NewGroupingAggregator
		ctrl.Finish()
	}()

	groupAgg := aggregation.NewMockGroupingAggregator(ctrl)
	newGroupingAgg = func(interval timeutil.Interval,
		intervalRatio int,
		timeRange timeutil.TimeRange,
		aggSpecs aggregation.AggregatorSpecs) aggregation.GroupingAggregator {
		return groupAgg
	}

	groupAgg.EXPECT().Aggregate(gomock.Any()).AnyTimes()
	groupAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{series.NewMockGroupedIterator(ctrl)})
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), &stmt.Query{}, ch)
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

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

package storagequery

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

var testExecPool = &tsdb.ExecutorPool{
	Filtering: concurrent.NewPool(
		"test-filtering-pool",
		runtime.GOMAXPROCS(-1), /*nRoutines*/
		time.Second*5,
		linmetric.StorageRegistry.NewScope("test-filtering-pool"),
	),
	Grouping: concurrent.NewPool(
		"test-grouping-pool",
		runtime.GOMAXPROCS(-1), /*nRoutines*/
		time.Second*5,
		linmetric.StorageRegistry.NewScope("test-filtering-pool"),
	),
	Scanner: concurrent.NewPool(
		"test-scanner-pool",
		runtime.GOMAXPROCS(-1), /*nRoutines*/
		time.Second*5,
		linmetric.StorageRegistry.NewScope("test-filtering-pool"),
	),
}

func TestStorageQueryFlow_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil).AnyTimes()
	execPool := concurrent.NewMockPool(ctrl)
	execPool.EXPECT().Submit(gomock.Any()).DoAndReturn(func(fn concurrent.Task) {
		fn()
	}).AnyTimes()
	pool := &tsdb.ExecutorPool{
		Filtering: execPool,
		Grouping:  execPool,
		Scanner:   execPool,
	}

	cases := []struct {
		name  string
		stage flow.Stage
	}{
		{
			name:  "filtering",
			stage: flow.FilteringStage,
		},
		{
			name:  "grouping",
			stage: flow.GroupingStage,
		},
		{
			name:  "scanning",
			stage: flow.ScannerStage,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			storageExecuteCtx := &flow.StorageExecuteContext{
				QueryInterval:      timeutil.Interval(timeutil.OneSecond),
				QueryIntervalRatio: 1,
				QueryTimeRange:     timeutil.TimeRange{},
				Query:              &stmt.Query{},
				TaskCtx:            flow.NewTaskContextWithTimeout(context.Background(), time.Second),
			}
			queryFlow := NewStorageQueryFlow(
				storageExecuteCtx,
				&protoCommonV1.TaskRequest{},
				taskServerFactory,
				&models.Leaf{Receivers: []models.StatelessNode{
					{HostIP: "1.1.1.1", GRPCPort: 1000},
					{HostIP: "1.1.1.2", GRPCPort: 2000},
				}},
				pool,
			)
			queryFlow.Prepare()
			qf := queryFlow.(*storageQueryFlow)
			reduceAgg := aggregation.NewMockGroupingAggregator(ctrl)
			qf.reduceAgg = reduceAgg
			reduceAgg.EXPECT().ResultSet().Return(nil).AnyTimes()
			reduceAgg.EXPECT().Aggregate(gomock.Any()).AnyTimes()
			var wait sync.WaitGroup
			wait.Add(1)
			go func() {
				queryFlow.Submit(tt.stage, func() {
					wait.Done()
				})
			}()
			wait.Wait()
		})
	}
}

func TestStorageQueryFlow_Submit_Fail(t *testing.T) {
	qf := &storageQueryFlow{}
	qf.completed.Store(true)
	qf.Submit(flow.FilteringStage, func() {
		panic("err")
	})
	qf = &storageQueryFlow{
		leafNode: &models.Leaf{},
		ctx:      context.TODO(),
	}
	qf.Submit(flow.DownSamplingStage, func() {
		panic("err")
	})
}

func TestStorageQueryFlow_completeTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := &flow.StorageExecuteContext{
		QueryInterval:       timeutil.Interval(timeutil.OneSecond),
		QueryIntervalRatio:  1,
		QueryTimeRange:      timeutil.TimeRange{},
		Query:               &stmt.Query{GroupBy: []string{"host"}},
		GroupingTagValueIDs: []*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)},
		Stats:               models.NewStorageStats(),
	}
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).AnyTimes()

	cases := []struct {
		name     string
		receives []models.StatelessNode
		prepare  func(qf *storageQueryFlow)
	}{
		{
			name: "test execute task after completed",
			prepare: func(qf *storageQueryFlow) {
				qf.completed.Store(true)
			},
		},
		{
			name: "test reduce result send, one receive",
			receives: []models.StatelessNode{
				{HostIP: "1.1.1.1", GRPCPort: 1000},
			},
			prepare: func(qf *storageQueryFlow) {
				mockBuildResultSet(qf, ctrl)
			},
		},
		{
			name: "test reduce result send, more receives",
			receives: []models.StatelessNode{
				{HostIP: "1.1.1.1", GRPCPort: 1000},
				{HostIP: "1.1.1.2", GRPCPort: 1000},
			},
			prepare: func(qf *storageQueryFlow) {
				mockBuildResultSet(qf, ctrl)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			storageExecuteCtx.TaskCtx = flow.NewTaskContextWithTimeout(context.Background(), time.Second)
			queryFlow := NewStorageQueryFlow(
				storageExecuteCtx,
				&protoCommonV1.TaskRequest{},
				taskServerFactory,
				&models.Leaf{Receivers: tt.receives},
				testExecPool,
			)

			queryFlow.Prepare()
			qf := queryFlow.(*storageQueryFlow)
			if tt.prepare != nil {
				tt.prepare(qf)
			}
			qf.completeTask(1)
		})
	}
}

func mockBuildResultSet(qf *storageQueryFlow, ctrl *gomock.Controller) {
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
	it.EXPECT().FieldName().Return(field.Name("f1"))
	groupIt.EXPECT().HasNext().Return(false)
	reduceAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{groupIt})
	go func() {
		qf.ReduceTagValues(0, map[uint32]string{100: "1.1.1.1"})
	}()
}

func TestStorageQueryFlow_getValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).AnyTimes()

	storageExecuteCtx := &flow.StorageExecuteContext{
		QueryInterval:      timeutil.Interval(timeutil.OneSecond),
		QueryIntervalRatio: 1,
		QueryTimeRange:     timeutil.TimeRange{},
		Query:              &stmt.Query{},
		TaskCtx:            flow.NewTaskContextWithTimeout(context.Background(), time.Second),
	}
	queryFlow := NewStorageQueryFlow(
		storageExecuteCtx,
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool,
	)

	queryFlow.Prepare()
	qf := queryFlow.(*storageQueryFlow)
	qf.tagValues = make([]string, 2)
	qf.tagsMap = make(map[string]string)
	qf.tagValuesMap = []map[uint32]string{{100: "1.1.1.1"}, {200: "1.1.1.2"}}
	// case 1: build new tag values str
	tagValueIDs := make([]byte, 2*4)
	binary.LittleEndian.PutUint32(tagValueIDs[0:], 100)
	binary.LittleEndian.PutUint32(tagValueIDs[4:], 200)
	tags := qf.getTagValues(string(tagValueIDs))
	assert.Equal(t, tag.ConcatTagValues([]string{"1.1.1.1", "1.1.1.2"}), tags)
	// case 2: get from cache
	tags = qf.getTagValues(string(tagValueIDs))
	assert.Equal(t, tag.ConcatTagValues([]string{"1.1.1.1", "1.1.1.2"}), tags)
}

func TestStorageQueryFlow_Task_panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := &flow.StorageExecuteContext{
		QueryInterval:      timeutil.Interval(timeutil.OneSecond),
		QueryIntervalRatio: 1,
		QueryTimeRange:     timeutil.TimeRange{},
		Query:              &stmt.Query{},
	}
	cases := []struct {
		name string
		in   func() (wait sync.WaitGroup, fn func())
	}{
		{
			name: "panic with err",
			in: func() (wait sync.WaitGroup, fn func()) {
				wait.Add(1)
				fn = func() {
					panic(fmt.Errorf("xxx"))
				}
				return
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			storageExecuteCtx.TaskCtx = flow.NewTaskContextWithTimeout(context.Background(), time.Second)
			queryFlow := NewStorageQueryFlow(
				storageExecuteCtx,
				&protoCommonV1.TaskRequest{},
				nil,
				&models.Leaf{},
				testExecPool)
			queryFlow.Prepare()

			wait, fn := tt.in()
			queryFlow.Submit(flow.FilteringStage, func() {
				wait.Done()
				fn()
			})
			wait.Wait()
		})
	}
}

func TestStorageQueryFlow_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := &flow.StorageExecuteContext{
		QueryInterval:      timeutil.Interval(timeutil.OneSecond),
		QueryIntervalRatio: 1,
		QueryTimeRange:     timeutil.TimeRange{},
		Query:              &stmt.Query{},
	}
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).AnyTimes()

	storageExecuteCtx.TaskCtx = flow.NewTaskContextWithTimeout(context.Background(), time.Second)
	queryFlow := NewStorageQueryFlow(
		storageExecuteCtx,
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool)

	queryFlow.Complete(nil) // err is nil, need not send err result
	server.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe).AnyTimes()
	queryFlow.Complete(fmt.Errorf("err")) // send err result
	queryFlow.Complete(fmt.Errorf("err")) // no send err result

	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil).AnyTimes()
	storageExecuteCtx.TaskCtx = flow.NewTaskContextWithTimeout(context.Background(), time.Second)
	queryFlow = NewStorageQueryFlow(
		storageExecuteCtx,
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool)
	queryFlow.Complete(fmt.Errorf("err")) // stream not found
}

func TestStorageQueryFlow_Complete_Wait_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	spec := aggregation.NewAggregatorSpec("f", field.SumField)
	spec.AddFunctionType(function.Sum)
	storageExecuteCtx := &flow.StorageExecuteContext{
		GroupingTagValueIDs: []*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)},
		Query:               &stmt.Query{GroupBy: []string{"db"}},
		AggregatorSpecs:     aggregation.AggregatorSpecs{spec},
		Stats:               models.NewStorageStats(),
		TaskCtx:             flow.NewTaskContextWithTimeout(context.Background(), time.Second),
	}
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).MaxTimes(2)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err")).MaxTimes(2)
	start := timeutil.Now()
	queryFlow := NewStorageQueryFlow(
		storageExecuteCtx,
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool)
	queryFlow.Prepare()
	qf := queryFlow.(*storageQueryFlow)
	qf.completeTask(12)
	assert.True(t, timeutil.Now()-start >= timeutil.OneSecond)
}

func TestStorageQueryFlow_Reduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reduceAgg := aggregation.NewMockGroupingAggregator(ctrl)
	qf := &storageQueryFlow{
		reduceAgg: reduceAgg,
	}
	// reduce agg result set
	reduceAgg.EXPECT().Aggregate(nil)
	qf.Reduce(nil)
	// qf is complete
	qf.completed.Store(true)
	qf.Reduce(nil)
}

func TestStorageQueryFlow_isCompleted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	qf := &storageQueryFlow{ctx: ctx, leafNode: &models.Leaf{}}
	assert.False(t, qf.isCompleted())
	cancel()
	assert.True(t, qf.isCompleted())
	qf.completed.Store(true)
	assert.True(t, qf.isCompleted())
}

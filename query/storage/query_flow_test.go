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
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
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
		linmetric.NewScope("test-filtering-pool"),
	),
	Grouping: concurrent.NewPool(
		"test-grouping-pool",
		runtime.GOMAXPROCS(-1), /*nRoutines*/
		time.Second*5,
		linmetric.NewScope("test-filtering-pool"),
	),
	Scanner: concurrent.NewPool(
		"test-scanner-pool",
		runtime.GOMAXPROCS(-1), /*nRoutines*/
		time.Second*5,
		linmetric.NewScope("test-filtering-pool"),
	),
}

func TestStorageQueryFlow_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := NewMockStorageExecuteContext(ctrl)
	storageExecuteCtx.EXPECT().QueryStats().Return(models.NewStorageStats()).AnyTimes()
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil)

	queryFlow := NewStorageQueryFlow(
		context.TODO(),
		storageExecuteCtx,
		&stmt.Query{},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool,
	)
	queryFlow.Prepare(timeutil.Interval(timeutil.OneSecond), 1, timeutil.TimeRange{}, nil)
	qf := queryFlow.(*storageQueryFlow)
	reduceAgg := aggregation.NewMockGroupingAggregator(ctrl)
	qf.reduceAgg = reduceAgg
	reduceAgg.EXPECT().ResultSet().Return(nil).AnyTimes()
	reduceAgg.EXPECT().Aggregate(gomock.Any()).AnyTimes()

	var wait sync.WaitGroup
	wait.Add(6)
	queryFlow.Filtering(func() {
		wait.Done()
		queryFlow.Grouping(func() {
			wait.Done()
			queryFlow.Load(func() {
				//seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
				//seriesAgg.EXPECT().Reset()

				//queryFlow.Reduce("1.1.1.1", aggregation.FieldAggregates{seriesAgg})
				wait.Done()
			})
		})
	})
	queryFlow.Filtering(func() {
		wait.Done()
		queryFlow.Grouping(func() {
			wait.Done()
			queryFlow.Load(func() {
				wait.Done()
			})
		})
	})
	wait.Wait()
	queryFlow.Complete(nil)
	time.Sleep(100 * time.Millisecond)
	//seriesAgg := aggregation.NewMockSeriesAggregator(ctrl)
	//seriesAgg.EXPECT().Reset()
	//reduce after query flow complete
	//queryFlow.Reduce("1.1.1.1", aggregation.FieldAggregates{seriesAgg})
}

func TestStorageQueryFlow_completeTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := NewMockStorageExecuteContext(ctrl)
	storageExecuteCtx.EXPECT().QueryStats().Return(nil).AnyTimes()
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).AnyTimes()

	queryFlow := NewStorageQueryFlow(
		context.TODO(),
		storageExecuteCtx,
		&stmt.Query{},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool,
	)

	queryFlow.Prepare(timeutil.Interval(timeutil.OneSecond), 1, timeutil.TimeRange{}, nil)
	qf := queryFlow.(*storageQueryFlow)
	// case 1: test execute task after completed
	qf.completed.Store(true)
	queryFlow.Filtering(func() {
		assert.Fail(t, "exec err")
	})

	// case 2: test reduce result send
	queryFlow = NewStorageQueryFlow(
		context.TODO(),
		storageExecuteCtx,
		&stmt.Query{GroupBy: []string{"host"}},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
		}},
		testExecPool,
	)

	queryFlow.Prepare(timeutil.Interval(timeutil.OneSecond), 1, timeutil.TimeRange{}, nil)
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
	it.EXPECT().FieldName().Return(field.Name("f1"))
	groupIt.EXPECT().HasNext().Return(false)
	reduceAgg.EXPECT().ResultSet().Return([]series.GroupedIterator{groupIt})
	var wait1 sync.WaitGroup
	wait1.Add(1)
	queryFlow.Filtering(func() {
		wait1.Done()
	})
	wait1.Wait()
	go func() {
		queryFlow.ReduceTagValues(0, map[uint32]string{100: "1.1.1.1"})
	}()
	time.Sleep(300 * time.Millisecond)
}

func TestStorageQueryFlow_getValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := NewMockStorageExecuteContext(ctrl)
	storageExecuteCtx.EXPECT().QueryStats().Return(nil).AnyTimes()
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).AnyTimes()

	queryFlow := NewStorageQueryFlow(
		context.TODO(),
		storageExecuteCtx,
		&stmt.Query{},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool,
	)

	queryFlow.Prepare(timeutil.Interval(timeutil.OneSecond), 1, timeutil.TimeRange{}, nil)
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

	storageExecuteCtx := NewMockStorageExecuteContext(ctrl)
	storageExecuteCtx.EXPECT().QueryStats().Return(nil).AnyTimes()
	queryFlow := NewStorageQueryFlow(context.TODO(),
		storageExecuteCtx, &stmt.Query{},
		&protoCommonV1.TaskRequest{},
		nil,
		&models.Leaf{},
		testExecPool)
	queryFlow.Prepare(timeutil.Interval(timeutil.OneSecond), 1, timeutil.TimeRange{}, nil)
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

func TestStorageQueryFlow_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageExecuteCtx := NewMockStorageExecuteContext(ctrl)
	storageExecuteCtx.EXPECT().QueryStats().Return(nil).AnyTimes()
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	server := protoCommonV1.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(server).Times(2)

	queryFlow := NewStorageQueryFlow(context.TODO(),
		storageExecuteCtx, &stmt.Query{},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool)

	queryFlow.Complete(nil) // err is nil, need not send err result
	server.EXPECT().Send(gomock.Any()).Return(io.ErrClosedPipe).Times(2)
	queryFlow.Complete(fmt.Errorf("err")) // send err result
	queryFlow.Complete(fmt.Errorf("err")) // no send err result

	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil).Times(2)
	queryFlow = NewStorageQueryFlow(context.TODO(),
		storageExecuteCtx, &stmt.Query{},
		&protoCommonV1.TaskRequest{},
		taskServerFactory,
		&models.Leaf{Receivers: []models.StatelessNode{
			{HostIP: "1.1.1.1", GRPCPort: 1000},
			{HostIP: "1.1.1.2", GRPCPort: 2000},
		}},
		testExecPool)
	queryFlow.Complete(fmt.Errorf("err")) // stream not found

}

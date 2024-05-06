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

package query

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	queryctx "github.com/lindb/lindb/query/context"
)

func TestTaskManager_ManageTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := NewTaskManager(nil, linmetric.BrokerRegistry)
	taskCtx := queryctx.NewMockTaskContext(ctrl)
	mgr.AddTask("1", taskCtx)
	mgr1 := mgr.(*taskManager)
	val := mgr1.statistics.AliveTask.Get()
	assert.Equal(t, float64(1), val)
	mgr.RemoveTask("1")
	val = mgr1.statistics.AliveTask.Get()
	assert.Equal(t, float64(0), val)
}

func TestTaskManager_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := NewTaskManager(
		concurrent.NewPool(
			"test", 10, time.Second,
			metrics.NewConcurrentStatistics("test", linmetric.BrokerRegistry)),
		linmetric.BrokerRegistry)

	taskCtx := queryctx.NewMockTaskContext(ctrl)
	mgr.AddTask("1", taskCtx)
	assert.Error(t, mgr.Receive(&protoCommonV1.TaskResponse{RequestID: "2"}, "test"))
	var wait sync.WaitGroup
	wait.Add(1)
	taskCtx.EXPECT().Context().Return(context.TODO())
	taskCtx.EXPECT().HandleResponse(gomock.Any(), "test").Do(func(_ *protoCommonV1.TaskResponse, _ string) {
		wait.Done()
	})
	assert.NoError(t, mgr.Receive(&protoCommonV1.TaskResponse{RequestID: "1"}, "test"))
	wait.Wait()
}

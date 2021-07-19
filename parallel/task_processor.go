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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./task_processor.go -destination=./task_processor_mock.go -package=parallel

// TaskDispatcher represents the task dispatcher
type TaskDispatcher interface {
	// Dispatch dispatches the task request based on task type
	Dispatch(ctx context.Context, stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest)
}

// TaskProcessor represents the task processor, all task processors are async
type TaskProcessor interface {
	// Process processes the task request
	Process(ctx context.Context, req *protoCommonV1.TaskRequest) error
}

// leafTaskDispatcher represents leaf task dispatcher for storage
type leafTaskDispatcher struct {
	processor TaskProcessor
	logger    *logger.Logger
}

// NewLeafTaskDispatcher creates a leaf task dispatcher
func NewLeafTaskDispatcher(currentNode models.Node,
	storageService service.StorageService,
	executorFactory ExecutorFactory, taskServerFactory rpc.TaskServerFactory) TaskDispatcher {
	return &leafTaskDispatcher{
		processor: newLeafTask(currentNode, storageService, executorFactory, taskServerFactory),
		logger:    logger.GetLogger("parallel", "LeafTaskDispatcher"),
	}
}

// Dispatch dispatches the request to storage engine query processor
func (d *leafTaskDispatcher) Dispatch(ctx context.Context, stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest) {
	err := d.processor.Process(ctx, req)
	if err != nil {
		if err1 := stream.Send(&protoCommonV1.TaskResponse{
			JobID:     req.JobID,
			TaskID:    req.ParentTaskID,
			Completed: true,
			ErrMsg:    err.Error(),
			SendTime:  timeutil.NowNano(),
		}); err1 != nil {
			d.logger.Error("send error message to target stream error", logger.Error(err))
		}
	}
}

// intermediateTaskDispatcher represents intermediate task dispatcher for broker
type intermediateTaskDispatcher struct {
}

// NewIntermediateTaskDispatcher create an intermediate task dispatcher
func NewIntermediateTaskDispatcher() TaskDispatcher {
	return &intermediateTaskDispatcher{}
}

// Dispatch dispatches the request to distribution query processor, merges the results
func (d *intermediateTaskDispatcher) Dispatch(ctx context.Context, stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest) {

}

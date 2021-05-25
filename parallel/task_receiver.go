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
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

// taskReceiver represents receive the task result from the sub tasks
type taskReceiver struct {
	jobManager JobManager
}

// NewTaskReceiver creates the task receiver
func NewTaskReceiver(jobManager JobManager) rpc.TaskReceiver {
	return &taskReceiver{jobManager: jobManager}
}

// Receive receives the task result, merges them and finally returns the final result
func (r *taskReceiver) Receive(resp *pb.TaskResponse) error {
	taskID := resp.TaskID
	taskManager := r.jobManager.GetTaskManager()
	taskCtx := taskManager.Get(taskID)
	if taskCtx == nil {
		return nil
	}

	taskCtx.ReceiveResult(resp)

	if taskCtx.Completed() {
		taskManager.Complete(taskID)

		if taskCtx.TaskType() == RootTask {
			jobCtx := r.jobManager.GetJob(resp.JobID)
			if jobCtx != nil && !jobCtx.Completed() {
				err := taskCtx.Error()
				if err != nil {
					jobCtx.Emit(&series.TimeSeriesEvent{Err: err})
				}
				jobCtx.Complete()
			}
		}
	}
	return nil
}

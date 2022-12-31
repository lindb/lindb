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

package stage

import (
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

// shardScanStage represents task request send stage.
type taskSendStage struct {
	baseStage
	taskCtx context.TaskContext
	target  string
	req     *protoCommonV1.TaskRequest
}

// NewTaskSendStage creates a taskSendStage instance.
func NewTaskSendStage(taskCtx context.TaskContext, target string, req *protoCommonV1.TaskRequest) Stage {
	return &taskSendStage{
		baseStage: baseStage{
			stageType: PhysicalPlan,
			//TODO: add async pool?
		},
		taskCtx: taskCtx,
		target:  target,
		req:     req,
	}
}

// Plan returns sub execution tree for task request send.
func (s *taskSendStage) Plan() PlanNode {
	return NewPlanNode(operator.NewTaskSender(s.taskCtx, s.target, s.req))
}

// Identifier returns identifier value of task request send stage.
func (s *taskSendStage) Identifier() string {
	return "TaskSend"
}

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

package operator

import (
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/query/context"
)

// taskSender represents task request send operator.
type taskSender struct {
	taskCtx context.TaskContext
	target  string
	req     *protoCommonV1.TaskRequest
}

// NewTaskSender creates a taskSender instance.
func NewTaskSender(taskCtx context.TaskContext, target string, req *protoCommonV1.TaskRequest) Operator {
	return &taskSender{
		taskCtx: taskCtx,
		target:  target,
		req:     req,
	}
}

// Execute returns send task request by given task context.
func (op *taskSender) Execute() error {
	return op.taskCtx.SendRequest(op.target, op.req)
}

// Identifier returns identifier string value of task send operator.
func (op *taskSender) Identifier() string {
	return "Task Sender"
}

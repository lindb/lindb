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

package rpc

import (
	context "context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/sql/execution"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
)

type CommandService struct {
	taskMgr execution.TaskManager
	logger  logger.Logger
}

func NewCommandService(taskMgr execution.TaskManager) protoCommandV1.CommandServiceServer {
	return &CommandService{
		taskMgr: taskMgr,
		logger:  logger.GetLogger("RPC", "resultSet"),
	}
}

func (srv *CommandService) Command(ctx context.Context, request *protoCommandV1.CommandRequest) (*protoCommandV1.CommandResponse, error) {
	switch request.Cmd {
	case protoCommandV1.Command_SubmitTask:
		req := &model.TaskRequest{}
		if err := encoding.JSONUnmarshal(request.Payload, req); err != nil {
			return nil, err
		}
		fragment := &plan.PlanFragment{}
		data, _ := req.Fragment.MarshalJSON()
		fmt.Println(string(data))
		err := encoding.JSONUnmarshal(data, fragment)
		if err != nil {
			return nil, err
		}
		fmt.Printf("task-req=%v\n", fragment)
		srv.taskMgr.SubmitTask(req, fragment)
	default:
		panic("not support cmd")
	}
	return &protoCommandV1.CommandResponse{}, nil
}

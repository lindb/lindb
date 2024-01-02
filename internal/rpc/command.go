package rpc

import (
	context "context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/execution"
	"github.com/lindb/lindb/execution/model"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
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
	}
	return &protoCommandV1.CommandResponse{}, nil
}

package parallel

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/common"
)

// TaskHandler represents the task rpc handler
type TaskHandler struct {
	fct        rpc.TaskServerFactory
	dispatcher TaskDispatcher

	logger *logger.Logger
}

// NewTaskHandler creates the task rpc handler
func NewTaskHandler(fct rpc.TaskServerFactory, dispatcher TaskDispatcher) *TaskHandler {
	return &TaskHandler{
		fct:        fct,
		dispatcher: dispatcher,
		logger:     logger.GetLogger("parallel/task/handler"),
	}
}

// Handle handles the task request based on grpc stream
func (q *TaskHandler) Handle(stream common.TaskService_HandleServer) error {
	clientLogicNode, err := rpc.GetLogicNodeFromContext(stream.Context())
	if err != nil {
		return err
	}

	nodeID := clientLogicNode.Indicator()

	q.fct.Register(nodeID, stream)
	q.logger.Info("register task stream", logger.String("client", nodeID))

	// when return, the stream is closed, Deregister the stream
	defer func() {
		q.fct.Deregister(nodeID)
		q.logger.Info("unregister task stream", logger.String("client", nodeID))
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			q.logger.Info("task server stream close")
			return nil
		}
		if err != nil {
			q.logger.Error("task server stream error", logger.Error(err))
			continue
		}
		q.dispatcher.Dispatch(req)
	}
}

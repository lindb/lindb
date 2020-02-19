package parallel

import (
	"context"
	"io"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/common"
)

// TaskHandler represents the task rpc handler
type TaskHandler struct {
	cfg        config.Query
	fct        rpc.TaskServerFactory
	dispatcher TaskDispatcher
	timeout    time.Duration

	taskPool concurrent.Pool

	logger *logger.Logger
}

// NewTaskHandler creates the task rpc handler
func NewTaskHandler(cfg config.Query, fct rpc.TaskServerFactory, dispatcher TaskDispatcher) *TaskHandler {
	return &TaskHandler{
		cfg:        cfg,
		timeout:    cfg.Timeout.Duration(),
		taskPool:   concurrent.NewPool("task-handle-pool", cfg.MaxWorkers, time.Second*5),
		fct:        fct,
		dispatcher: dispatcher,
		logger:     logger.GetLogger("parallel", "TaskHandler"),
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
		q.dispatch(req)
	}
}

// dispatch dispatches request with timeout
func (q *TaskHandler) dispatch(req *common.TaskRequest) {
	//FIXME add timeout????
	ctx, cancel := context.WithTimeout(context.TODO(), q.timeout)
	q.taskPool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				q.logger.Error("dispatch task request", logger.Any("err", err), logger.Stack())
			}
			cancel()
		}()
		q.dispatcher.Dispatch(ctx, req)
	})
}

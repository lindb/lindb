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
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/rpc"
)

// TaskHandler represents the task rpc handler
type TaskHandler struct {
	cfg       config.Query
	fct       rpc.TaskServerFactory
	processor TaskProcessor
	timeout   time.Duration

	taskPool concurrent.Pool

	logger *logger.Logger
}

// NewTaskHandler creates the task rpc handler
func NewTaskHandler(
	cfg config.Query,
	fct rpc.TaskServerFactory,
	processor TaskProcessor,
	pool concurrent.Pool,
) *TaskHandler {
	return &TaskHandler{
		cfg:       cfg,
		timeout:   cfg.Timeout.Duration(),
		taskPool:  pool,
		fct:       fct,
		processor: processor,
		logger:    logger.GetLogger("query", "TaskHandler"),
	}
}

// Handle handles the task request based on grpc stream
func (q *TaskHandler) Handle(stream protoCommonV1.TaskService_HandleServer) (err error) {
	clientLogicNode, err := rpc.GetLogicNodeFromContext(stream.Context())
	if err != nil {
		return err
	}

	nodeID := clientLogicNode.Indicator()

	epoch := q.fct.Register(nodeID, stream)
	q.logger.Info("register task stream",
		logger.String("client", nodeID), logger.Int64("epoch", epoch))

	// when return, the stream is closed, Deregister the stream
	defer func() {
		ok := q.fct.Deregister(epoch, nodeID)
		if ok {
			q.logger.Info("unregister task stream successfully",
				logger.String("client", nodeID), logger.Int64("epoch", epoch))
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			q.logger.Error("task server stream error", logger.Error(err))
			return err
		}
		q.process(stream, req)
	}
}

// dispatch dispatches request with timeout
func (q *TaskHandler) process(stream protoCommonV1.TaskService_HandleServer, req *protoCommonV1.TaskRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout)
	q.taskPool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				q.logger.Error("dispatch task request", logger.Any("err", err), logger.Stack())
			}
			cancel()
		}()
		q.processor.Process(ctx, stream, req)
	})
}

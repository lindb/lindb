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
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
)

//go:generate mockgen -source ./task_transport.go -destination=./task_transport_mock.go -package=rpc

// TaskClientFactory represents the task stream manage
type TaskClientFactory interface {
	// CreateTaskClient creates a task client stream if not exist
	CreateTaskClient(target models.Node) error
	// GetTaskClient returns the task client stream by target node
	GetTaskClient(target string) protoCommonV1.TaskService_HandleClient
	// CloseTaskClient closes the task client stream for target node
	CloseTaskClient(targetNodeID string) (closed bool, err error)
	// SetTaskReceiver set task receiver for handling task response
	SetTaskReceiver(taskReceiver TaskReceiver)
}

type taskClient struct {
	cli      protoCommonV1.TaskService_HandleClient
	targetID string
	target   models.Node
	running  atomic.Bool
	ready    atomic.Bool
}

// taskClientFactory implements TaskClientFactory interface
type taskClientFactory struct {
	ctx          context.Context
	currentNode  models.Node
	taskReceiver TaskReceiver
	// target node ID => client stream
	taskStreams map[string]*taskClient
	mutex       sync.RWMutex

	newTaskServiceClientFunc func(cc *grpc.ClientConn) protoCommonV1.TaskServiceClient
	connFct                  ClientConnFactory
	logger                   *logger.Logger
}

// NewTaskClientFactory creates a task client factory
func NewTaskClientFactory(ctx context.Context, currentNode models.Node, connFct ClientConnFactory) TaskClientFactory {
	return &taskClientFactory{
		ctx:                      ctx,
		currentNode:              currentNode,
		connFct:                  connFct,
		taskStreams:              make(map[string]*taskClient),
		newTaskServiceClientFunc: protoCommonV1.NewTaskServiceClient,
		logger:                   logger.GetLogger("RPC", "TaskClient"),
	}
}

// SetTaskReceiver set task receiver for handling task response
func (f *taskClientFactory) SetTaskReceiver(taskReceiver TaskReceiver) {
	f.taskReceiver = taskReceiver
}

// GetTaskClient returns the task client stream by target node
func (f *taskClientFactory) GetTaskClient(target string) protoCommonV1.TaskService_HandleClient {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	if stream, ok := f.taskStreams[target]; ok && stream != nil {
		return stream.cli
	}
	return nil
}

// CreateTaskClient creates a stream task client if not exist,
// then create a goroutine handle task response if created successfully.
func (f *taskClientFactory) CreateTaskClient(target models.Node) error {
	targetNodeID := target.Indicator()
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, ok := f.taskStreams[targetNodeID]; ok {
		return nil
	}

	taskClient := &taskClient{
		targetID: targetNodeID,
		target:   target,
	}
	taskClient.running.Store(true)

	go f.handleTaskResponse(taskClient)

	// cache task client stream
	f.taskStreams[targetNodeID] = taskClient
	return nil
}

// CloseTaskClient closes the task client stream for target node
func (f *taskClientFactory) CloseTaskClient(targetNodeID string) (closed bool, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if client, ok := f.taskStreams[targetNodeID]; ok && client.cli != nil {
		client.running.Store(false)
		err = client.cli.CloseSend()
		delete(f.taskStreams, targetNodeID)
		return closed, err
	}
	return false, nil
}

func (f *taskClientFactory) initTaskClient(client *taskClient) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if client.cli != nil {
		if err := client.cli.CloseSend(); err != nil {
			f.logger.Error("close task client error", logger.Error(err))
		}
		client.cli = nil
	}
	conn, err := f.connFct.GetClientConn(client.target)
	if err != nil {
		return err
	}

	// https://pkg.go.dev/google.golang.org/grpc#ClientConn.NewStream
	// context is the lifetime of stream
	ctx := CreateOutgoingContextWithPairs(f.ctx, constants.RPCMetaKeyLogicNode, f.currentNode.Indicator())
	cli, err := f.newTaskServiceClientFunc(conn).Handle(ctx)
	if err != nil {
		return err
	}
	client.cli = cli
	return nil
}

// handleTaskResponse handles task response loop, if stream closed exist loop
func (f *taskClientFactory) handleTaskResponse(client *taskClient) {
	var attempt int32 = 0
	for client.running.Load() {
		select {
		case <-f.ctx.Done():
			// if client is not ready, this goroutine may be blocked without ctx.Done()
			return
		default:
		}

		if !client.ready.Load() {
			attempt++
			f.logger.Info("initializing task client",
				logger.String("target", client.targetID),
				logger.Int32("attempt", attempt),
			)
			if err := f.initTaskClient(client); err != nil {
				f.logger.Error("failed to initialize task client",
					logger.Error(err),
					logger.String("target", client.targetID),
					logger.Int32("attempt", attempt),
				)
				time.Sleep(time.Second)
				continue
			} else {
				f.logger.Info("initialized task client successfully",
					logger.String("target", client.targetID),
					logger.Int32("attempt", attempt))
				client.ready.Store(true)
			}
		}
		var cli protoCommonV1.TaskService_HandleClient
		f.mutex.RLock()
		cli = client.cli
		f.mutex.RUnlock()
		resp, err := cli.Recv()
		if err != nil {
			client.ready.Store(false)
			// todo: suppress errors before shard assignment
			f.logger.Error("receive task error from stream", logger.Error(err))
			continue
		}

		if err = f.taskReceiver.Receive(resp, client.targetID); err != nil {
			f.logger.Error("receive task response",
				logger.String("taskID", resp.TaskID),
				logger.String("taskType", resp.Type.String()),
				logger.Error(err))
		}
	}
}

// TaskServerFactory represents a factory to get server stream.
type TaskServerFactory interface {
	// GetStream returns a ServerStream for a node.
	GetStream(node string) protoCommonV1.TaskService_HandleServer
	// Register registers a stream for a node.
	Register(node string, stream protoCommonV1.TaskService_HandleServer) (epoch int64)
	// Deregister unregisters a stream for node, if returns true, unregister successfully.
	Deregister(epoch int64, node string) bool
	// Nodes returns all registered nodes.
	Nodes() []models.Node
}

type taskService struct {
	handle protoCommonV1.TaskService_HandleServer
	epoch  int64
}

// taskServerFactory implements TaskServerFactory interface
type taskServerFactory struct {
	nodeMap map[string]*taskService
	epoch   atomic.Int64
	lock    sync.RWMutex
	logger  *logger.Logger
}

// NewTaskServerFactory returns the singleton server stream factory
func NewTaskServerFactory() TaskServerFactory {
	return &taskServerFactory{
		nodeMap: make(map[string]*taskService),
		logger:  logger.GetLogger("RPC", "TaskServer"),
	}
}

// GetStream returns a ServerStream for a node.
func (fct *taskServerFactory) GetStream(node string) protoCommonV1.TaskService_HandleServer {
	fct.lock.RLock()
	defer fct.lock.RUnlock()

	if st, ok := fct.nodeMap[node]; ok {
		return st.handle
	}
	return nil
}

// Register registers a stream for a node.
func (fct *taskServerFactory) Register(node string, stream protoCommonV1.TaskService_HandleServer) (epoch int64) {
	fct.lock.Lock()
	defer fct.lock.Unlock()
	epoch = fct.epoch.Inc()
	fct.nodeMap[node] = &taskService{
		epoch:  epoch,
		handle: stream,
	}
	return epoch
}

// Nodes returns all registered nodes.
func (fct *taskServerFactory) Nodes() []models.Node {
	fct.lock.RLock()
	defer fct.lock.RUnlock()

	nodes := make([]models.Node, 0, len(fct.nodeMap))
	for nodeID := range fct.nodeMap {
		node, err := models.ParseNode(nodeID)
		if err != nil {
			fct.logger.Warn("parse node error", logger.Error(err))
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// Deregister unregisters a stream for node.
func (fct *taskServerFactory) Deregister(epoch int64, node string) bool {
	fct.lock.Lock()
	defer fct.lock.Unlock()

	if st, ok := fct.nodeMap[node]; ok && st.epoch == epoch {
		delete(fct.nodeMap, node)
		return true
	}
	return false
}

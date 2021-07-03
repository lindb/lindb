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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc/proto/common"
)

//go:generate mockgen -source ./task_transport.go -destination=./task_transport_mock.go -package=rpc

var log = logger.GetLogger("rpc", "TaskClient")

// TaskClientFactory represents the task stream manage
type TaskClientFactory interface {
	// CreateTaskClient creates a task client stream if not exist
	CreateTaskClient(target models.Node) error
	// GetTaskClient returns the task client stream by target node
	GetTaskClient(target string) common.TaskService_HandleClient
	// CloseTaskClient closes the task client stream for target node
	CloseTaskClient(targetNodeID string) (closed bool, err error)
	// SetTaskReceiver set task receiver for handling task response
	SetTaskReceiver(taskReceiver TaskReceiver)
}

type taskClient struct {
	cli      common.TaskService_HandleClient
	targetID string
	target   models.Node
	running  atomic.Bool
	ready    atomic.Bool
}

// taskClientFactory implements TaskClientFactory interface
type taskClientFactory struct {
	currentNode  models.Node
	taskReceiver TaskReceiver
	// target node ID => client stream
	taskStreams map[string]*taskClient
	mutex       sync.RWMutex

	newTaskServiceClientFunc func(cc *grpc.ClientConn) common.TaskServiceClient
	connFct                  ClientConnFactory
}

// NewTaskClientFactory creates a task client factory
func NewTaskClientFactory(currentNode models.Node) TaskClientFactory {
	return &taskClientFactory{
		currentNode:              currentNode,
		connFct:                  GetClientConnFactory(),
		taskStreams:              make(map[string]*taskClient),
		newTaskServiceClientFunc: common.NewTaskServiceClient,
	}
}

// SetTaskReceiver set task receiver for handling task response
func (f *taskClientFactory) SetTaskReceiver(taskReceiver TaskReceiver) {
	f.taskReceiver = taskReceiver
}

// GetTaskClient returns the task client stream by target node
func (f *taskClientFactory) GetTaskClient(target string) common.TaskService_HandleClient {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	stream, ok := f.taskStreams[target]
	if ok && stream != nil {
		return stream.cli
	}
	return nil
}

// CreateTaskClient creates a stream task client if not exist,
// then create a goroutine handle task response if created successfully.
func (f *taskClientFactory) CreateTaskClient(target models.Node) error {
	targetNodeID := (&target).Indicator()
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, ok := f.taskStreams[targetNodeID]
	if ok {
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

	client, ok := f.taskStreams[targetNodeID]
	if ok && client.cli != nil {
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
			log.Error("close task client error", logger.Error(err))
		}
		client.cli = nil
	}
	conn, err := f.connFct.GetClientConn(client.target)
	if err != nil {
		return err
	}

	//TODO handle context?????
	ctx := createOutgoingContextWithPairs(context.TODO(), metaKeyLogicNode, (&f.currentNode).Indicator())
	cli, err := f.newTaskServiceClientFunc(conn).Handle(ctx)
	if err != nil {
		return err
	}
	client.cli = cli
	return nil
}

// handleTaskResponse handles task response loop, if stream closed exist loop
func (f *taskClientFactory) handleTaskResponse(client *taskClient) {
	var sequence int32 = 0
	for client.running.Load() {
		if !client.ready.Load() {
			sequence++
			log.Info("initializing task client",
				logger.String("target", client.targetID),
				logger.Int32("sequence", sequence),
			)
			if err := f.initTaskClient(client); err != nil {
				log.Error("failed to initialize task client",
					logger.Error(err),
					logger.String("target", client.targetID),
					logger.Int32("sequence", sequence),
				)
				time.Sleep(time.Second)
				continue
			} else {
				log.Info("initialized task client successfully",
					logger.String("target", client.targetID),
					logger.Int32("sequence", sequence))
				client.ready.Store(true)
			}
		}
		resp, err := client.cli.Recv()
		if err != nil {
			client.ready.Store(false)
			log.Error("receive task error from stream", logger.Error(err))
			continue
		}

		err = f.taskReceiver.Receive(resp)
		if err != nil {
			log.Error("receive task response", logger.Any("rep", resp), logger.Error(err))
		}
	}
}

// TaskServerFactory represents a factory to get server stream.
type TaskServerFactory interface {
	// GetStream returns a ServerStream for a node.
	GetStream(node string) common.TaskService_HandleServer
	// Register registers a stream for a node.
	Register(node string, stream common.TaskService_HandleServer) (epoch int64)
	// Deregister unregisters a stream for node, if returns true, unregister successfully.
	Deregister(epoch int64, node string) bool
	// Nodes returns all registered nodes.
	Nodes() []models.Node
}

type taskService struct {
	handle common.TaskService_HandleServer
	epoch  int64
}

// taskServerFactory implements TaskServerFactory interface
type taskServerFactory struct {
	nodeMap map[string]*taskService
	epoch   atomic.Int64
	lock    sync.RWMutex
}

// NewTaskServerFactory returns the singleton server stream factory
func NewTaskServerFactory() TaskServerFactory {
	return &taskServerFactory{
		nodeMap: make(map[string]*taskService),
	}
}

// GetStream returns a ServerStream for a node.
func (fct *taskServerFactory) GetStream(node string) common.TaskService_HandleServer {
	fct.lock.RLock()
	defer fct.lock.RUnlock()

	st, ok := fct.nodeMap[node]
	if ok {
		return st.handle
	}
	return nil
}

// Register registers a stream for a node.
func (fct *taskServerFactory) Register(node string, stream common.TaskService_HandleServer) (epoch int64) {
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
			log.Warn("parse node error", logger.Error(err))
			continue
		}
		nodes = append(nodes, *node)
	}
	return nodes
}

// Deregister unregisters a stream for node.
func (fct *taskServerFactory) Deregister(epoch int64, node string) bool {
	fct.lock.Lock()
	defer fct.lock.Unlock()
	st, ok := fct.nodeMap[node]
	if ok && st.epoch == epoch {
		delete(fct.nodeMap, node)
		return true
	}
	return false
}

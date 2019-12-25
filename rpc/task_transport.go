package rpc

import (
	"context"
	"io"
	"sync"

	"go.uber.org/atomic"

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
	CloseTaskClient(targetNodeID string)
	// SetTaskReceiver set task receiver for handling task response
	SetTaskReceiver(taskReceiver TaskReceiver)
}

type taskClient struct {
	cli     common.TaskService_HandleClient
	running atomic.Bool
}

// taskClientFactory implements TaskClientFactory interface
type taskClientFactory struct {
	currentNode  models.Node
	taskReceiver TaskReceiver
	// target node ID => client stream
	taskStreams map[string]*taskClient
	mutex       sync.RWMutex

	connFct ClientConnFactory
}

// NewTaskClientFactory creates a task client factory
func NewTaskClientFactory(currentNode models.Node) TaskClientFactory {
	return &taskClientFactory{
		currentNode: currentNode,
		connFct:     GetClientConnFactory(),
		taskStreams: make(map[string]*taskClient),
	}
}

// SetTaskReceiver set task receiver for handling task response
func (f *taskClientFactory) SetTaskReceiver(taskReceiver TaskReceiver) {
	f.taskReceiver = taskReceiver
}

func (f *taskClientFactory) GetTaskClient(target string) common.TaskService_HandleClient {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.taskStreams[target].cli
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

	conn, err := f.connFct.GetClientConn(target)
	if err != nil {
		return err
	}

	//TODO handle context?????
	ctx := createOutgoingContextWithPairs(context.TODO(), metaKeyLogicNode, (&f.currentNode).Indicator())
	cli, err := common.NewTaskServiceClient(conn).Handle(ctx)
	if err != nil {
		return err
	}

	taskClient := &taskClient{
		cli: cli,
	}
	taskClient.running.Store(true)
	go f.handleTaskResponse(taskClient)
	// cache task client stream
	f.taskStreams[targetNodeID] = taskClient
	return nil
}

// CloseTaskClient closes the task client stream for target node
func (f *taskClientFactory) CloseTaskClient(targetNodeID string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	client, ok := f.taskStreams[targetNodeID]
	if ok {
		if err := client.cli.CloseSend(); err != nil {
			log.Error("close task client stream", logger.String("target", targetNodeID), logger.Error(err))
		}
		client.running.Store(false)
		delete(f.taskStreams, targetNodeID)
		log.Info("close task client stream", logger.String("target", targetNodeID))
	}
}

// handleTaskResponse handles task response loop, if stream closed exist loop
func (f *taskClientFactory) handleTaskResponse(client *taskClient) {
	for client.running.Load() {
		resp, err := client.cli.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Error("receive task error from stream", logger.Error(err))
			continue
		}

		err = f.taskReceiver.Receive(resp)
		if err != nil {
			log.Error("receive task response", logger.Any("rep", resp), logger.Error(err))
		}
	}
}

// ServerStreamFactory represents a factory to get server stream.
type TaskServerFactory interface {
	// GetStream returns a ServerStream for a node.
	GetStream(node string) common.TaskService_HandleServer
	// Register registers a stream for a node.
	Register(node string, stream common.TaskService_HandleServer)
	// Deregister unregisters a stream for node.
	Deregister(node string)
	// Nodes returns all registered nodes.
	Nodes() []models.Node
}

// taskServerFactory implements TaskServerFactory interface
type taskServerFactory struct {
	nodeMap map[string]common.TaskService_HandleServer
	lock    sync.RWMutex
}

// GetServerStreamFactory returns the singleton server stream factory
func NewTaskServerFactory() TaskServerFactory {
	return &taskServerFactory{
		nodeMap: make(map[string]common.TaskService_HandleServer),
	}
}

// GetStream returns a ServerStream for a node.
func (fct *taskServerFactory) GetStream(node string) common.TaskService_HandleServer {
	fct.lock.RLock()
	defer fct.lock.RUnlock()

	st := fct.nodeMap[node]
	return st
}

// Register registers a stream for a node.
func (fct *taskServerFactory) Register(node string, stream common.TaskService_HandleServer) {
	fct.lock.Lock()
	defer fct.lock.Unlock()

	fct.nodeMap[node] = stream
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
func (fct *taskServerFactory) Deregister(node string) {
	fct.lock.Lock()
	defer fct.lock.Unlock()
	delete(fct.nodeMap, node)
}

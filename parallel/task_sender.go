package parallel

import (
	"sync"

	pb "github.com/lindb/lindb/rpc/proto/common"
)

//go:generate mockgen -source=./task_sender.go -destination=./task_sender_mock.go -package=parallel

// TaskSenderManager represents the task sender stream both client and server side
type TaskSenderManager interface {
	// GetClientStream returns the client stream for sending task request
	GetClientStream(indicator string) pb.TaskService_HandleClient
	// AddClientStream adds the client stream when rpc connected
	AddClientStream(indicator string, clientStream pb.TaskService_HandleClient)
	// RemoveClientStream removes the client stream when rpc dis-connected
	RemoveClientStream(indicator string)

	// GetServerStream returns the server stream for sending task response
	GetServerStream(indicator string) pb.TaskService_HandleServer
	// AddServerStream adds the server stream when rpc connected
	AddServerStream(indicator string, serverStream pb.TaskService_HandleServer)
	// RemoveServerStream removes the server stream when rpc dis-connected
	RemoveServerStream(indicator string)
}

// taskSenderManager implements the task sender manager interface for caching client/server stream
type taskSenderManager struct {
	clientStreams sync.Map
	serverStreams sync.Map
}

// NewTaskSenderManager creates the task sender manager
func NewTaskSenderManager() TaskSenderManager {
	return &taskSenderManager{}
}

// AddClientStream adds the client stream when rpc connected
func (m *taskSenderManager) AddClientStream(indicator string, clientStream pb.TaskService_HandleClient) {
	m.clientStreams.Store(indicator, clientStream)
}

// RemoveClientStream removes the client stream when rpc dis-connected
func (m *taskSenderManager) RemoveClientStream(indicator string) {
	m.clientStreams.Delete(indicator)
}

// GetClientStream returns the client stream for sending task request
func (m *taskSenderManager) GetClientStream(indicator string) pb.TaskService_HandleClient {
	stream, ok := m.clientStreams.Load(indicator)
	if !ok {
		return nil
	}
	clientStream, ok := stream.(pb.TaskService_HandleClient)
	if !ok {
		return nil
	}
	return clientStream
}

// GetServerStream returns the server stream for sending task response
func (m *taskSenderManager) GetServerStream(indicator string) pb.TaskService_HandleServer {
	stream, ok := m.serverStreams.Load(indicator)
	if !ok {
		return nil
	}
	serverStream, ok := stream.(pb.TaskService_HandleServer)
	if !ok {
		return nil
	}
	return serverStream
}

// AddServerStream adds the server stream when rpc connected
func (m *taskSenderManager) AddServerStream(indicator string, serverStream pb.TaskService_HandleServer) {
	m.serverStreams.Store(indicator, serverStream)
}

// RemoveServerStream removes the server stream when rpc dis-connected
func (m *taskSenderManager) RemoveServerStream(indicator string) {
	m.serverStreams.Delete(indicator)
}

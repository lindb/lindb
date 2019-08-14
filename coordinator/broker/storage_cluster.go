package broker

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

var log = logger.GetLogger("storage/cluster/state")

type StorageClusterState struct {
	state       *models.StorageState
	taskStreams map[string]pb.TaskService_HandleClient

	streamFactory rpc.ClientStreamFactory
}

func newStorageClusterState(streamFactory rpc.ClientStreamFactory) *StorageClusterState {
	return &StorageClusterState{
		taskStreams:   make(map[string]pb.TaskService_HandleClient),
		streamFactory: streamFactory,
	}
}

func (s *StorageClusterState) SetState(state *models.StorageState) {
	log.Info("start set new storage cluster")
	var needDelete []string
	for nodeID := range s.taskStreams {
		_, ok := state.ActiveNodes[nodeID]
		if !ok {
			needDelete = append(needDelete, nodeID)
		}
	}

	for _, nodeID := range needDelete {
		s.removeTaskStream(nodeID)
	}

	for nodeID, node := range state.ActiveNodes {
		s.removeTaskStream(nodeID)
		// create a new client stream
		client, err := s.streamFactory.CreateTaskClient(node.Node)
		if err != nil {
			log.Error("create task client stream", logger.Error(err))
			continue
		}

		// cache task client stream
		s.taskStreams[nodeID] = client
	}

	s.state = state
	log.Info("set new storage cluster successfully")
}

func (s *StorageClusterState) close() {
	log.Info("start close storage cluster state")
	for nodeID := range s.taskStreams {
		s.removeTaskStream(nodeID)
	}
	log.Info("close storage cluster state successfully")
}

func (s *StorageClusterState) removeTaskStream(nodeID string) {
	client, ok := s.taskStreams[nodeID]
	if ok {
		if err := client.CloseSend(); err != nil {
			log.Error("close task client stream", logger.Error(err))
		}
		delete(s.taskStreams, nodeID)
		log.Info("close task client stream", logger.String("target", nodeID))
	}
}

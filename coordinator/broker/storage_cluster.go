package broker

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

var log = logger.GetLogger("broker", "StorageState")

const dummy = ""

type StorageClusterState struct {
	state             *models.StorageState
	taskStreams       map[string]string
	taskClientFactory rpc.TaskClientFactory
}

func newStorageClusterState(taskClientFactory rpc.TaskClientFactory) *StorageClusterState {
	return &StorageClusterState{
		taskClientFactory: taskClientFactory,
		taskStreams:       make(map[string]string),
	}
}

func (s *StorageClusterState) SetState(state *models.StorageState) {
	log.Info("set new storage cluster state")
	var needDelete []string
	for nodeID := range s.taskStreams {
		_, ok := state.ActiveNodes[nodeID]
		if !ok {
			needDelete = append(needDelete, nodeID)
		}
	}

	for _, nodeID := range needDelete {
		s.taskClientFactory.CloseTaskClient(nodeID)
		delete(s.taskStreams, nodeID)
	}

	for nodeID, node := range state.ActiveNodes {
		// create a new client stream
		if err := s.taskClientFactory.CreateTaskClient(node.Node); err != nil {
			log.Error("create task client stream",
				logger.String("target", (&node.Node).Indicator()), logger.Error(err))
			s.taskClientFactory.CloseTaskClient(nodeID)
			delete(s.taskStreams, nodeID)
			continue
		}
		s.taskStreams[nodeID] = dummy
	}

	s.state = state
	log.Info("set new storage cluster successfully")
}

func (s *StorageClusterState) close() {
	log.Info("start close storage cluster state")
	for nodeID := range s.taskStreams {
		s.taskClientFactory.CloseTaskClient(nodeID)
		delete(s.taskStreams, nodeID)
	}
	log.Info("close storage cluster state successfully")
}

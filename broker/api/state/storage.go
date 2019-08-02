package state

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
)

// StorageAPI represents query storage cluster's state api from broker state machine
type StorageAPI struct {
	stateMachine broker.StorageStateMachine
}

// NewStorageAPI creates storage state api
func NewStorageAPI(stateMachine broker.StorageStateMachine) *StorageAPI {
	return &StorageAPI{
		stateMachine: stateMachine,
	}
}

// ListStorageCluster lists state of all storage clusters
func (s *StorageAPI) ListStorageCluster(w http.ResponseWriter, r *http.Request) {
	clusters := s.stateMachine.List()
	api.OK(w, clusters)
}

package state

import (
	"net/http"

	"github.com/eleme/lindb/broker/api"
	"github.com/eleme/lindb/coordinator/broker"
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

// List lists state of all storage clusters
func (s *StorageAPI) List(w http.ResponseWriter, r *http.Request) {
	clusters := s.stateMachine.List()
	api.OK(w, clusters)
}

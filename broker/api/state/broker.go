package state

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
)

// BrokerAPI represents query broker state api from broker state machine
type BrokerAPI struct {
	stateMachine broker.NodeStateMachine
}

// NewBrokerAPI creates the broker state api
func NewBrokerAPI(stateMachine broker.NodeStateMachine) *BrokerAPI {
	return &BrokerAPI{
		stateMachine: stateMachine,
	}
}

// ListBrokerNodes lists all alive broker nodes
func (s *BrokerAPI) ListBrokerNodes(w http.ResponseWriter, r *http.Request) {
	nodes := s.stateMachine.GetActiveNodes()
	api.OK(w, nodes)
}

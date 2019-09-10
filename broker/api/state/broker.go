package state

import (
	"fmt"
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
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

// ListBrokerNodes lists all alive broker nodes, if query parameter type is provided, returns the specific node type.
func (s *BrokerAPI) ListBrokerNodes(w http.ResponseWriter, r *http.Request) {
	nodeTypeStr, err := api.GetParamsFromRequest("type", r, string(models.NodeTypeRPC), false)
	if err != nil {
		api.Error(w, err)
		return
	}

	nodeType := models.ParseNodeType(nodeTypeStr)
	if nodeType == models.NodeTypeUnknown {
		api.Error(w, fmt.Errorf("node type %s unknown", nodeTypeStr))
		return
	}

	nodes := s.stateMachine.GetActiveNodesByType(nodeType)
	api.OK(w, nodes)
}

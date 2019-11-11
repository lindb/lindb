package state

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

// BrokerAPI represents query broker state api from broker state machine
type BrokerAPI struct {
	ctx          context.Context
	repo         state.Repository
	stateMachine broker.NodeStateMachine
}

// NewBrokerAPI creates the broker state api
func NewBrokerAPI(ctx context.Context, repo state.Repository, stateMachine broker.NodeStateMachine) *BrokerAPI {
	return &BrokerAPI{
		ctx:          ctx,
		repo:         repo,
		stateMachine: stateMachine,
	}
}

func (s *BrokerAPI) ListBrokersStat(w http.ResponseWriter, r *http.Request) {
	kvs, err := s.repo.List(s.ctx, constants.StateNodesPath)
	if err != nil {
		api.Error(w, err)
		return
	}
	// get active nodes
	nodes := s.stateMachine.GetActiveNodes()
	nodeIDs := make(map[string]string)
	for _, node := range nodes {
		id := node.Node.Indicator()
		nodeIDs[id] = id
	}
	// build result
	var result []models.NodeStat
	for _, kv := range kvs {
		_, nodeID := filepath.Split(kv.Key)
		stat := models.NodeStat{}
		if err := encoding.JSONUnmarshal(kv.Value, &stat); err != nil {
			api.Error(w, err)
			return
		}
		_, ok := nodeIDs[nodeID]
		if !ok {
			stat.IsDead = true
		}
		result = append(result, stat)
	}
	api.OK(w, result)
}

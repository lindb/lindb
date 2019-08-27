package service

import (
	"context"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./replicator.go -destination=./replicator_mock.go -package service

// ReplicatorService represents the replicator state report
type ReplicatorService interface {
	// Report reports all wal replicator state under current broker
	Report(state *models.BrokerReplicaState) error
}

// replicatorService implements ReplicatorService interface
type replicatorService struct {
	node models.Node
	repo state.Repository
}

// NewReplicatorService creates replicator service for current broker node
func NewReplicatorService(node models.Node, repo state.Repository) ReplicatorService {
	return &replicatorService{
		node: node,
		repo: repo,
	}
}

// Report reports all wal replicator state under current broker
func (s *replicatorService) Report(state *models.BrokerReplicaState) error {
	data := encoding.JSONMarshal(state)
	//TODO need using timeout
	if err := s.repo.Put(context.TODO(), constants.GetReplicaStatePath((&s.node).Indicator()), data); err != nil {
		return err
	}
	return nil
}

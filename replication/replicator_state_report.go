package replication

import (
	"context"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./replicator_state_report.go -destination=./replicator_state_report_mock.go -package replication

// ReplicatorStateReport represents the replicator state report
type ReplicatorStateReport interface {
	// Report reports all wal replicator state under current broker
	Report(state *models.BrokerReplicaState) error
}

// replicatorStateReport implements ReplicatorStateReport interface
type replicatorStateReport struct {
	node models.Node
	repo state.Repository
}

// NewReplicatorStateReport creates replicator state report for current broker node
func NewReplicatorStateReport(node models.Node, repo state.Repository) ReplicatorStateReport {
	return &replicatorStateReport{
		node: node,
		repo: repo,
	}
}

// Report reports all wal replicator state under current broker
func (s *replicatorStateReport) Report(state *models.BrokerReplicaState) error {
	data := encoding.JSONMarshal(state)
	//TODO need using timeout
	if err := s.repo.Put(context.TODO(), constants.GetReplicaStatePath((&s.node).Indicator()), data); err != nil {
		return err
	}
	return nil
}

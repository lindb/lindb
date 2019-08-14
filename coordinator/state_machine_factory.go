package coordinator

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./state_machine_factory.go -destination=./state_machine_factory_mock.go -package=coordinator

// StateMachineCfg represents the state machine config
type StateMachineCfg struct {
	Ctx                 context.Context
	CurrentNode         models.Node
	DiscoveryFactory    discovery.Factory
	ClientStreamFactory rpc.ClientStreamFactory // rpc client stream create factory
}

// StateMachineFactory represents the state machine create factory
type StateMachineFactory interface {
	// CreateNodeStateMachine creates the node state machine
	CreateNodeStateMachine() (broker.NodeStateMachine, error)
	// CreateStorageStateMachine creates the storage state machine
	CreateStorageStateMachine() (broker.StorageStateMachine, error)
	// CreateReplicaStatusStateMachine creates the shard replica status state machine
	CreateReplicaStatusStateMachine() (replica.StatusStateMachine, error)
}

// stateMachineFactory implements the interface, using state machine config for creating
type stateMachineFactory struct {
	cfg *StateMachineCfg
}

// NewStateMachineFactory creates the factory using config
func NewStateMachineFactory(cfg *StateMachineCfg) StateMachineFactory {
	return &stateMachineFactory{cfg: cfg}
}

// CreateNodeStateMachine creates the node state machine, if fail returns err
func (s *stateMachineFactory) CreateNodeStateMachine() (broker.NodeStateMachine, error) {
	return broker.NewNodeStateMachine(s.cfg.Ctx, s.cfg.CurrentNode, s.cfg.DiscoveryFactory)
}

// CreateStorageStateMachine creates the storage state machine, if fail returns err
func (s *stateMachineFactory) CreateStorageStateMachine() (broker.StorageStateMachine, error) {
	return broker.NewStorageStateMachine(s.cfg.Ctx, s.cfg.DiscoveryFactory, s.cfg.ClientStreamFactory)
}

// CreateReplicaStatusStateMachine creates the shard replica status state machine, if fail returns err
func (s *stateMachineFactory) CreateReplicaStatusStateMachine() (replica.StatusStateMachine, error) {
	return replica.NewStatusStateMachine(s.cfg.Ctx, s.cfg.DiscoveryFactory)
}

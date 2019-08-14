package coordinator

import (
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/pkg/logger"
)

// BrokerStateMachines represents all state machines for broker
type BrokerStateMachines struct {
	StorageSM       broker.StorageStateMachine
	NodeSM          broker.NodeStateMachine
	ReplicaStatusSM replica.StatusStateMachine

	factory StateMachineFactory

	log *logger.Logger
}

func NewBrokerStateMachines(factory StateMachineFactory) *BrokerStateMachines {
	return &BrokerStateMachines{
		factory: factory,
		log:     logger.GetLogger("coordinator/broker/sm"),
	}
}

// Start starts related state machines for broker
func (s *BrokerStateMachines) Start() (err error) {
	s.NodeSM, err = s.factory.CreateNodeStateMachine()
	if err != nil {
		return err
	}
	s.StorageSM, err = s.factory.CreateStorageStateMachine()
	if err != nil {
		return err
	}
	s.ReplicaStatusSM, err = s.factory.CreateReplicaStatusStateMachine()
	if err != nil {
		return err
	}
	return nil
}

// Stop stops the broker's state machines
func (s *BrokerStateMachines) Stop() {
	if s.StorageSM != nil {
		if err := s.StorageSM.Close(); err != nil {
			s.log.Error("close storage state state machine error", logger.Error(err))
		}
	}
	if s.NodeSM != nil {
		if err := s.NodeSM.Close(); err != nil {
			s.log.Error("close node state state machine error", logger.Error(err))
		}
	}
	if s.ReplicaStatusSM != nil {
		if err := s.ReplicaStatusSM.Close(); err != nil {
			s.log.Error("close replica status state state machine error", logger.Error(err))
		}
	}
}

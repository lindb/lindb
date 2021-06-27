// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package coordinator

import (
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/pkg/logger"
)

// BrokerStateMachines represents all state machines for broker
type BrokerStateMachines struct {
	StorageSM       broker.StorageStateMachine
	NodeSM          broker.NodeStateMachine
	ReplicaStatusSM replica.StatusStateMachine
	ReplicatorSM    replica.ReplicatorStateMachine
	DatabaseSM      database.DBStateMachine

	factory StateMachineFactory

	log *logger.Logger
}

func NewBrokerStateMachines(factory StateMachineFactory) *BrokerStateMachines {
	return &BrokerStateMachines{
		factory: factory,
		log:     logger.GetLogger("coordinator", "BrokerStateMachines"),
	}
}

// Start starts related state machines for broker
func (s *BrokerStateMachines) Start() (err error) {
	s.log.Info("starting BrokerStateMachines")
	s.NodeSM, err = s.factory.CreateNodeStateMachine()
	if err != nil {
		return err
	}
	s.log.Debug("started NodeStateMachine")
	s.ReplicatorSM, err = s.factory.CreateReplicatorStateMachine()
	if err != nil {
		return err
	}
	s.log.Debug("started ReplicatorStateMachine")
	s.StorageSM, err = s.factory.CreateStorageStateMachine()
	if err != nil {
		return err
	}
	s.ReplicaStatusSM, err = s.factory.CreateReplicaStatusStateMachine()
	if err != nil {
		return err
	}
	s.log.Debug("started ReplicaStatusStateMachine")
	s.DatabaseSM, err = s.factory.CreateDatabaseStateMachine()
	if err != nil {
		return err
	}
	s.log.Debug("started DatabaseStateMachine")
	s.log.Info("started BrokerStateMachines")
	return nil
}

// Stop stops the broker's state machines
func (s *BrokerStateMachines) Stop() {
	if s.StorageSM != nil {
		if err := s.StorageSM.Close(); err != nil {
			s.log.Error("close storage state machine error", logger.Error(err))
		}
	}
	if s.NodeSM != nil {
		if err := s.NodeSM.Close(); err != nil {
			s.log.Error("close node state machine error", logger.Error(err))
		}
	}
	if s.ReplicaStatusSM != nil {
		if err := s.ReplicaStatusSM.Close(); err != nil {
			s.log.Error("close replica status state machine error", logger.Error(err))
		}
	}
	if s.ReplicatorSM != nil {
		if err := s.ReplicatorSM.Close(); err != nil {
			s.log.Error("close replicator state machine error", logger.Error(err))
		}
	}
	if s.DatabaseSM != nil {
		if err := s.DatabaseSM.Close(); err != nil {
			s.log.Error("close database state machine error", logger.Error(err))
		}
	}
}

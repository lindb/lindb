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

package context

import (
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/pkg/logger"
)

// StateMachine represents all state machine for master
type StateMachine struct {
	StorageCluster storage.ClusterStateMachine
	DatabaseAdmin  broker.ShardAssignmentStateMachine
}

// MasterContext represents master context, creates it after node elect master
type MasterContext struct {
	StateMachine *StateMachine
}

// NewMasterContext creates master context using state machine
func NewMasterContext(stateMachine *StateMachine) *MasterContext {
	return &MasterContext{
		StateMachine: stateMachine,
	}
}

// Close closes all state machines, releases resource that master used
func (m *MasterContext) Close() {
	log := logger.GetLogger("coordinator", "MasterContext")
	if m.StateMachine.StorageCluster != nil {
		if err := m.StateMachine.StorageCluster.Close(); err != nil {
			log.Error("close storage cluster state machine error", logger.Error(err), logger.Stack())
		}
	}
	if m.StateMachine.DatabaseAdmin != nil {
		if err := m.StateMachine.DatabaseAdmin.Close(); err != nil {
			log.Error("close database admin state machine error", logger.Error(err), logger.Stack())
		}
	}
}

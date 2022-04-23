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

package metrics

import (
	"fmt"

	"github.com/lindb/lindb/internal/linmetric"
)

// StateManagerStatistics represents state manager statistics.
type StateManagerStatistics struct {
	HandleEvents       *linmetric.DeltaCounterVec // handle event success count
	HandleEventFailure *linmetric.DeltaCounterVec // handle event failure count
	Panics             *linmetric.DeltaCounterVec // panic count when handle event
}

// NewStateManagerStatistics creates a state manager statistics.
func NewStateManagerStatistics(role string) *StateManagerStatistics {
	scope := linmetric.BrokerRegistry.NewScope(fmt.Sprintf("lindb.%s.state_manager", role))
	return &StateManagerStatistics{
		HandleEvents:       scope.NewCounterVec("handle_events", "type"),
		HandleEventFailure: scope.NewCounterVec("handle_event_failures", "type"),
		Panics:             scope.NewCounterVec("panics", "type"),
	}
}

// ShardLeaderStatistics represents shard leader elect statistics.
type ShardLeaderStatistics struct {
	LeaderElections     *linmetric.BoundCounter // shard leader elect successfully
	LeaderElectFailures *linmetric.BoundCounter // shard leader elect failure
}

// NewShardLeaderStatistics create a shard leader elect statistics.
func NewShardLeaderStatistics() *ShardLeaderStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.master.shard.leader")
	return &ShardLeaderStatistics{
		LeaderElections:     scope.NewCounter("elections"),
		LeaderElectFailures: scope.NewCounter("elect_failures"),
	}
}

// MasterStatistics represents master statistics.
type MasterStatistics struct {
	FailOvers        *linmetric.BoundCounter // master fail over successfully
	FailOverFailures *linmetric.BoundCounter // master fail over failure
	Reassigns        *linmetric.BoundCounter // master reassign successfully
	ReassignFailures *linmetric.BoundCounter // master reassign failure
}

// NewMasterStatistics creates a master statistics.
func NewMasterStatistics() *MasterStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.master.controller")
	return &MasterStatistics{
		FailOvers:        scope.NewCounter("failovers"),
		FailOverFailures: scope.NewCounter("failover_failures"),
		Reassigns:        scope.NewCounter("reassigns"),
		ReassignFailures: scope.NewCounter("reassign_failures"),
	}
}

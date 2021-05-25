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

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

package state

import (
	"context"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/http"
)

var (
	BrokerStatePath = "/broker/cluster/state"
)

// BrokerAPI represents query broker state api from broker state machine.
type BrokerAPI struct {
	ctx context.Context

	deps *deps.HTTPDeps
}

// NewBrokerAPI creates the broker state api.
func NewBrokerAPI(ctx context.Context, deps *deps.HTTPDeps) *BrokerAPI {
	return &BrokerAPI{
		ctx:  ctx,
		deps: deps,
	}
}

// Register adds broker state url route.
func (s *BrokerAPI) Register(route gin.IRoutes) {
	route.GET(BrokerStatePath, s.ListBrokersState)
}

// ListBrokersState returns brokers state.
func (s *BrokerAPI) ListBrokersState(c *gin.Context) {
	kvs, err := s.deps.Repo.List(s.ctx, constants.StateNodesPath)
	if err != nil {
		http.Error(c, err)
		return
	}
	// get active nodes
	nodes := s.deps.StateMachines.NodeSM.GetActiveNodes()
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
			http.Error(c, err)
			return
		}
		_, ok := nodeIDs[nodeID]
		if !ok {
			stat.IsDead = true
		}
		result = append(result, stat)
	}
	http.OK(c, result)
}

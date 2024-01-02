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
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/lindb/common/pkg/http"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
)

var (
	ExplorePath = "/state/machine/explore"
)

type param struct {
	Type string `form:"type" binding:"required"`
}

type StorageStateMachineAPI struct {
	stateMgr storage.StateManager
	logger   logger.Logger
}

// NewStorageStateMachineAPI creates storage state machine api instance.
func NewStorageStateMachineAPI(stateMgr storage.StateManager) *StorageStateMachineAPI {
	return &StorageStateMachineAPI{
		stateMgr: stateMgr,
		logger:   logger.GetLogger("Storage", "StateMachineAPI"),
	}
}

// Register adds state machine url route.
func (api *StorageStateMachineAPI) Register(route gin.IRoutes) {
	route.GET(ExplorePath, api.Explore)
}

// Explore explores the state from storage state machine.
func (api *StorageStateMachineAPI) Explore(c *gin.Context) {
	param := &param{}
	err := c.ShouldBindQuery(param)
	if err != nil {
		http.Error(c, err)
		return
	}
	switch param.Type {
	case constants.ShardAssignment:
		shardAssignments := api.stateMgr.GetShardAssignments()
		sort.Slice(shardAssignments, func(i, j int) bool {
			return shardAssignments[i].Name < shardAssignments[j].Name
		})
		http.OK(c, shardAssignments)
	case constants.LiveNode:
		nodes := api.stateMgr.GetLiveNodes()
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Indicator() < nodes[j].Indicator()
		})
		http.OK(c, nodes)
	default:
		http.NotFound(c)
	}
}

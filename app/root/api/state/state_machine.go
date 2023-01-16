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

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	ExplorePath = "/state/machine/explore"
)

type Param struct {
	Type string `form:"type" binding:"required"`
}

// RootStateMachineAPI represents state machine explore api.
type RootStateMachineAPI struct {
	deps   *depspkg.HTTPDeps
	logger *logger.Logger
}

// NewRootStateMachineAPI creates root state machine api instance.
func NewRootStateMachineAPI(deps *depspkg.HTTPDeps) *RootStateMachineAPI {
	return &RootStateMachineAPI{
		deps:   deps,
		logger: logger.GetLogger("Root", "StateMachineAPI"),
	}
}

// Register adds state machine url route.
func (api *RootStateMachineAPI) Register(route gin.IRoutes) {
	route.GET(ExplorePath, api.Explore)
}

// Explore explores the state from state machine of broker/live node/database.
func (api *RootStateMachineAPI) Explore(c *gin.Context) {
	param := &Param{}
	err := c.ShouldBindQuery(param)
	if err != nil {
		http.Error(c, err)
		return
	}
	switch param.Type {
	case constants.BrokerState:
		api.writeBrokerState(c, api.deps.StateMgr.GetBrokerStates())
	case constants.LiveNode:
		nodes := api.deps.StateMgr.GetLiveNodes()
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Indicator() < nodes[j].Indicator()
		})
		http.OK(c, nodes)
	case constants.DatabaseConfig:
		api.writeDatabaseState(c, api.deps.StateMgr.GetDatabases())
	default:
		http.NotFound(c)
	}
}

// writeDatabaseState writes response with database.
func (api *RootStateMachineAPI) writeDatabaseState(c *gin.Context, dbs []models.LogicDatabase) {
	sort.Slice(dbs, func(i, j int) bool {
		return dbs[i].Name < dbs[j].Name
	})
	http.OK(c, dbs)
}

// writeBrokerState writes response with broker state.
func (api *RootStateMachineAPI) writeBrokerState(c *gin.Context, storages []models.BrokerState) {
	sort.Slice(storages, func(i, j int) bool {
		return storages[i].Name < storages[j].Name
	})
	http.OK(c, storages)
}

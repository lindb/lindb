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
	"github.com/gin-gonic/gin"

	"github.com/lindb/common/pkg/http"
	"github.com/lindb/common/pkg/logger"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/constants"
)

var (
	ExplorePath = "/state/machine/explore"
)

type Param struct {
	Type       string `form:"type" binding:"required"`
	BrokerName string `form:"brokerName"`
}

// RootStateMachineAPI represents state machine explore api.
type RootStateMachineAPI struct {
	deps   *depspkg.HTTPDeps
	logger logger.Logger
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
		if param.BrokerName != "" {
			state, ok := api.deps.StateMgr.GetBrokerState(param.BrokerName)
			if ok {
				http.OK(c, state)
			} else {
				http.NotFound(c)
			}
		} else {
			http.OK(c, api.deps.StateMgr.GetBrokerStates())
		}
	case constants.LiveNode:
		http.OK(c, api.deps.StateMgr.GetLiveNodes())
	case constants.DatabaseConfig:
		http.OK(c, api.deps.StateMgr.GetDatabases())
	default:
		http.NotFound(c)
	}
}

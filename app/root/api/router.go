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

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/root/api/state"
	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	apipkg "github.com/lindb/lindb/internal/api"
	"github.com/lindb/lindb/internal/linmetric"
	httppkg "github.com/lindb/lindb/pkg/http"
)

// API represents root http api.
type API struct {
	execute          *ExecuteAPI
	rootStateMachine *state.RootStateMachineAPI
	metricExplore    *apipkg.ExploreAPI
	env              *apipkg.EnvAPI
	config           *apipkg.ConfigAPI
	log              *apipkg.LoggerAPI
	proxy            *httppkg.ReverseProxy
}

// NewAPI creates root http api.
func NewAPI(deps *depspkg.HTTPDeps) *API {
	return &API{
		execute:          NewExecuteAPI(deps),
		rootStateMachine: state.NewRootStateMachineAPI(deps),
		metricExplore:    apipkg.NewExploreAPI(deps.GlobalKeyValues, linmetric.RootRegistry),
		env:              apipkg.NewEnvAPI(config.ToEnvs(deps.Cfg, config.NewDefaultRoot())),
		log:              apipkg.NewLoggerAPI(deps.Cfg.Logging.Dir),
		config:           apipkg.NewConfigAPI(deps.Node, deps.Cfg),
		proxy:            httppkg.NewReverseProxy(),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	v1 := router.Group(constants.APIVersion1)
	// execute lin query language statement
	api.execute.Register(v1)
	// monitoring
	api.metricExplore.Register(v1)
	api.rootStateMachine.Register(v1)
	api.config.Register(v1)
	api.log.Register(v1)

	api.proxy.Register(v1)
	api.env.Register(v1)
}

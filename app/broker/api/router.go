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

	"github.com/lindb/lindb/app/broker/api/admin"
	"github.com/lindb/lindb/app/broker/api/exec"
	"github.com/lindb/lindb/app/broker/api/ingest"
	"github.com/lindb/lindb/app/broker/api/state"
	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	apipkg "github.com/lindb/lindb/internal/api"
	"github.com/lindb/lindb/internal/linmetric"
	httppkg "github.com/lindb/lindb/pkg/http"
)

// API represents broker http api.
type API struct {
	deps *depspkg.HTTPDeps

	execute            *exec.ExecuteAPI
	database           *admin.DatabaseAPI
	flusher            *admin.DatabaseFlusherAPI
	storage            *admin.StorageClusterAPI
	brokerStateMachine *state.BrokerStateMachineAPI
	request            *apipkg.RequestAPI
	metricExplore      *apipkg.ExploreAPI
	log                *apipkg.LoggerAPI
	config             *apipkg.ConfigAPI
	env                *apipkg.EnvAPI
	write              *ingest.Write
	proxy              *httppkg.ReverseProxy
}

// NewAPI creates broker http api.
func NewAPI(deps *depspkg.HTTPDeps) *API {
	return &API{
		deps:               deps,
		execute:            exec.NewExecuteAPI(deps),
		database:           admin.NewDatabaseAPI(deps),
		flusher:            admin.NewDatabaseFlusherAPI(deps),
		storage:            admin.NewStorageClusterAPI(deps),
		brokerStateMachine: state.NewBrokerStateMachineAPI(deps),
		request:            apipkg.NewRequestAPI(),
		metricExplore:      apipkg.NewExploreAPI(deps.GlobalKeyValues, linmetric.BrokerRegistry),
		log:                apipkg.NewLoggerAPI(deps.BrokerCfg.Logging.Dir),
		config:             apipkg.NewConfigAPI(deps.Node, deps.BrokerCfg),
		env:                apipkg.NewEnvAPI(deps.BrokerCfg.Monitor, constants.BrokerRole),
		write:              ingest.NewWrite(deps),
		proxy:              httppkg.NewReverseProxy(),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	router.Use(SlowSQLLog(api.deps))
	v1 := router.Group(constants.APIVersion1)
	// execute lin query language statement
	api.execute.Register(v1)

	api.database.Register(v1)
	api.flusher.Register(v1)
	api.storage.Register(v1)

	// state
	api.brokerStateMachine.Register(v1)
	api.request.Register(v1)

	// write metric data
	api.write.Register(v1)

	// monitoring
	api.metricExplore.Register(v1)
	api.log.Register(v1)
	api.config.Register(v1)

	api.env.Register(v1)
	api.proxy.Register(v1)
}

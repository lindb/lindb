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
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
)

// API represents broker http api.
type API struct {
	execute *exec.ExecuteAPI

	database           *admin.DatabaseAPI
	flusher            *admin.DatabaseFlusherAPI
	storage            *admin.StorageClusterAPI
	brokerStateMachine *state.BrokerStateMachineAPI
	metricExplore      *monitoring.ExploreAPI
	log                *monitoring.LoggerAPI
	config             *monitoring.ConfigAPI
	influxIngestion    *ingest.InfluxWriter
	protoIngestion     *ingest.ProtoWriter
	flatIngestion      *ingest.FlatWriter
	proxy              *ReverseProxy
}

// NewAPI creates broker http api.
func NewAPI(deps *depspkg.HTTPDeps) *API {
	return &API{
		execute:            exec.NewExecuteAPI(deps),
		database:           admin.NewDatabaseAPI(deps),
		flusher:            admin.NewDatabaseFlusherAPI(deps),
		storage:            admin.NewStorageClusterAPI(deps),
		brokerStateMachine: state.NewBrokerStateMachineAPI(deps),
		metricExplore:      monitoring.NewExploreAPI(deps.GlobalKeyValues, linmetric.BrokerRegistry),
		log:                monitoring.NewLoggerAPI(deps.BrokerCfg.Logging.Dir),
		config:             monitoring.NewConfigAPI(deps.Node, deps.BrokerCfg),
		influxIngestion:    ingest.NewInfluxWriter(deps),
		protoIngestion:     ingest.NewProtoWriter(deps),
		flatIngestion:      ingest.NewFlatWriter(deps),
		proxy:              NewReverseProxy(),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	v1 := router.Group(constants.APIVersion1)
	api.execute.Register(v1)

	api.database.Register(v1)
	api.flusher.Register(v1)
	api.storage.Register(v1)

	api.brokerStateMachine.Register(v1)

	api.influxIngestion.Register(v1)
	api.protoIngestion.Register(v1)
	api.flatIngestion.Register(v1)

	// monitoring
	api.metricExplore.Register(v1)
	api.log.Register(v1)
	api.config.Register(v1)

	api.proxy.Register(v1)
}

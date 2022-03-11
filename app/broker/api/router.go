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
	"github.com/lindb/lindb/app/broker/api/metadata"
	"github.com/lindb/lindb/app/broker/api/state"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/monitoring"
)

// API represents broker http api.
type API struct {
	execute *exec.ExecuteAPI

	database        *admin.DatabaseAPI
	flusher         *admin.DatabaseFlusherAPI
	storage         *admin.StorageClusterAPI
	explore         *metadata.ExploreAPI
	stateExplore    *state.ExploreAPI
	replica         *state.ReplicaAPI
	metricExplore   *monitoring.ExploreAPI
	log             *monitoring.LoggerAPI
	config          *monitoring.ConfigAPI
	influxIngestion *ingest.InfluxWriter
	protoIngestion  *ingest.ProtoWriter
	flatIngestion   *ingest.FlatWriter
	proxy           *ReverseProxy
}

// NewAPI creates broker http api.
func NewAPI(deps *deps.HTTPDeps) *API {
	return &API{
		execute:         exec.NewExecuteAPI(deps),
		database:        admin.NewDatabaseAPI(deps),
		flusher:         admin.NewDatabaseFlusherAPI(deps),
		storage:         admin.NewStorageClusterAPI(deps),
		explore:         metadata.NewExploreAPI(deps),
		stateExplore:    state.NewExploreAPI(deps),
		replica:         state.NewReplicaAPI(deps),
		metricExplore:   monitoring.NewExploreAPI(deps.GlobalKeyValues),
		log:             monitoring.NewLoggerAPI(deps.BrokerCfg.Logging.Dir),
		config:          monitoring.NewConfigAPI(deps.Node, deps.BrokerCfg),
		influxIngestion: ingest.NewInfluxWriter(deps),
		protoIngestion:  ingest.NewProtoWriter(deps),
		flatIngestion:   ingest.NewFlatWriter(deps),
		proxy:           NewReverseProxy(),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	api.execute.Register(router)

	api.database.Register(router)
	api.flusher.Register(router)
	api.storage.Register(router)
	api.explore.Register(router)

	api.stateExplore.Register(router)
	api.replica.Register(router)

	api.influxIngestion.Register(router)
	api.protoIngestion.Register(router)
	api.flatIngestion.Register(router)
	// monitoring
	api.metricExplore.Register(router)
	api.log.Register(router)
	api.config.Register(router)

	api.proxy.Register(router)
}

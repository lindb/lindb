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
	"github.com/lindb/lindb/app/broker/api/cluster"
	"github.com/lindb/lindb/app/broker/api/ingest"
	"github.com/lindb/lindb/app/broker/api/metadata"
	"github.com/lindb/lindb/app/broker/api/query"
	"github.com/lindb/lindb/app/broker/api/state"
	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/monitoring"
)

// API represents broker http api.
type API struct {
	master          *cluster.MasterAPI
	database        *admin.DatabaseAPI
	flusher         *admin.DatabaseFlusherAPI
	storage         *admin.StorageClusterAPI
	explore         *metadata.ExploreAPI
	brokerState     *state.BrokerAPI
	storageState    *state.StorageAPI
	stateExplore    *state.ExploreAPI
	metricExplore   *monitoring.ExploreAPI
	influxIngestion *ingest.InfluxWriter
	protoIngestion  *ingest.ProtoWriter
	flatIngestion   *ingest.FlatWriter
	metric          *query.MetricAPI
	metadata        *query.MetadataAPI
}

// NewAPI creates broker http api.
func NewAPI(deps *deps.HTTPDeps) *API {
	return &API{
		master:          cluster.NewMasterAPI(deps),
		database:        admin.NewDatabaseAPI(deps),
		flusher:         admin.NewDatabaseFlusherAPI(deps),
		storage:         admin.NewStorageClusterAPI(deps),
		explore:         metadata.NewExploreAPI(deps),
		brokerState:     state.NewBrokerAPI(deps),
		storageState:    state.NewStorageAPI(deps),
		stateExplore:    state.NewExploreAPI(deps),
		metricExplore:   monitoring.NewExploreAPI(deps.GlobalKeyValues),
		influxIngestion: ingest.NewInfluxWriter(deps),
		protoIngestion:  ingest.NewProtoWriter(deps),
		flatIngestion:   ingest.NewFlatWriter(deps),
		metric:          query.NewMetricAPI(deps),
		metadata:        query.NewMetadataAPI(deps),
	}
}

// RegisterRouter registers http api router.
func (api *API) RegisterRouter(router *gin.RouterGroup) {
	api.master.Register(router)
	api.database.Register(router)
	api.flusher.Register(router)
	api.storage.Register(router)
	api.explore.Register(router)

	api.brokerState.Register(router)
	api.storageState.Register(router)
	api.stateExplore.Register(router)
	api.metricExplore.Register(router)

	api.metadata.Register(router)
	api.metric.Register(router)
	api.influxIngestion.Register(router)
	api.protoIngestion.Register(router)
	api.flatIngestion.Register(router)
}

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

package write

import (
	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/ingestion/influx"
	"github.com/lindb/lindb/replication"

	"net/http"
)

// InfluxWriter processes Influxdb line protocol
type InfluxWriter struct {
	cm replication.ChannelManager
}

// NewInfluxWriter creates influx writer
func NewInfluxWriter(cm replication.ChannelManager) *InfluxWriter {
	return &InfluxWriter{cm: cm}
}

func (iw *InfluxWriter) Write(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	namespace, _ := api.GetParamsFromRequest("ns", r, constants.DefaultNamespace, false)
	enrichedTags, err := ingestCommon.ExtractEnrichTags(r)
	if err != nil {
		api.Error(w, err)
		return
	}
	metricList, err := influx.Parse(r, enrichedTags, namespace)
	if err != nil {
		api.Error(w, err)
		return
	}
	if err := iw.cm.Write(databaseName, metricList); err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

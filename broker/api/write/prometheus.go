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
	"io/ioutil"
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	ingestCommon "github.com/lindb/lindb/ingestion/common"
	"github.com/lindb/lindb/ingestion/prometheus"
	"github.com/lindb/lindb/replication"
)

// for testing
var (
	readAllFunc = ioutil.ReadAll
)

// PrometheusWrite represents support prometheus text protocol
type PrometheusWrite struct {
	cm replication.ChannelManager
}

// NewPrometheusWrite creates prometheus write
func NewPrometheusWrite(cm replication.ChannelManager) *PrometheusWrite {
	return &PrometheusWrite{
		cm: cm,
	}
}

// Write parses prometheus text protocol then writes data into wal
func (m *PrometheusWrite) Write(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	_, _ = api.GetParamsFromRequest("ns", r, constants.DefaultNamespace, false)
	s, err := readAllFunc(r.Body)
	if err != nil {
		api.Error(w, err)
		return
	}

	enrichedTags, err := ingestCommon.ExtractEnrichTags(r)
	if err != nil {
		api.Error(w, err)
		return
	}
	metricList, err := prometheus.PromParse(s, enrichedTags)
	if err != nil {
		api.Error(w, err)
		return
	}

	if err := m.cm.Write(databaseName, metricList); err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, "success")
}

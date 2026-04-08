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

// Package metrics provides the Arrow-based metric ingest HTTP handler.
package metrics

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/api/ingest"
	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
)

// MetricPath is the URL path for the Arrow metric write endpoint.
const MetricPath = "/metrics"

// Metric handles Arrow IPC metric write requests.
// Clients must set Content-Type: application/vnd.apache.arrow.stream and
// X-LinDB-Database: <database-name>.  The request body must be
// the raw IPC stream bytes produced by MetricBuilder.Bytes().
type Metric struct {
	ingest.BaseWriter // inherits Write handler
}

// NewMetric creates a new Metric ingest handler.
func NewMetric(deps *depspkg.HTTPDeps) *Metric {
	return &Metric{
		BaseWriter: ingest.NewBaseWriter(deps, map[string]ingest.ProcessFunc{
			constants.ContentTypeArrow: ingest.StandardProcess(deps, constants.EncodingArrow),
		}),
	}
}

// Register registers POST and PUT routes for MetricPath.
func (m *Metric) Register(route gin.IRoutes) {
	route.POST(MetricPath, m.Write)
	route.PUT(MetricPath, m.Write)
}

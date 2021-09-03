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

package ingest

import (
	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/ingestion/prometheus"
)

var (
	PrometheusWritePath = "/prometheus/write"
)

// PrometheusWriter processes prometheus text protocol.
type PrometheusWriter struct {
	commonWriter
}

// NewPrometheusWriter creates prometheus writer.
func NewPrometheusWriter(deps *deps.HTTPDeps) *PrometheusWriter {
	return &PrometheusWriter{
		commonWriter: commonWriter{
			deps:   deps,
			parser: prometheus.Parse,
		},
	}
}

// Register adds prometheus write url route.
func (m *PrometheusWriter) Register(route gin.IRoutes) {
	route.PUT(PrometheusWritePath, m.Write)
	route.POST(PrometheusWritePath, m.Write)
}

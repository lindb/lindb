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
	"github.com/lindb/lindb/ingestion/influx"
	"github.com/lindb/lindb/pkg/http"
)

var (
	InfluxWritePath = "/influx/write"
	InfluxQueryPath = "/influx/query"
)

// InfluxWriter processes Influxdb line protocol.
type InfluxWriter struct {
	commonWriter
}

// NewInfluxWriter creates influx writer.
func NewInfluxWriter(deps *deps.HTTPDeps) *InfluxWriter {
	return &InfluxWriter{
		commonWriter: commonWriter{
			deps:   deps,
			parser: influx.Parse,
		},
	}
}

// Register adds influx write url route.
func (iw *InfluxWriter) Register(route gin.IRoutes) {
	route.PUT(InfluxWritePath, iw.Write)
	route.POST(InfluxWritePath, iw.Write)

	route.POST(InfluxQueryPath, iw.FakeQuery)
}

// FakeQuery answers influxdb agent such as telegraf
func (iw *InfluxWriter) FakeQuery(c *gin.Context) {
	http.OK(c, "ok")
}
